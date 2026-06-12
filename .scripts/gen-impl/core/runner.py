from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
import re
import sys

from .config import Config
from .io_safe import can_write_path, ensure_dir_writable, log, warn, write_text, timestamp_utc
from .reporting import (
    benchmark_failure_causes,
    extract_failure_suggestions,
    failure_improvement_suggestions_from_causes,
    folder_has_non_test_impl_files,
    list_non_test_impl_files,
    report_performance_status_from_file,
    folder_report_path,
    parse_benchmark_rows_from_log,
    parse_unit_test_scenarios_from_json_log,
    render_report_content,
    render_benchmark_table,
    render_failure_causes_table,
    render_failure_improvement_table,
    render_scenario_table,
    report_attempts_from_file,
    report_is_fresh,
    report_status_from_file,
    runtime_failure_cause_from_log,
    validate_generated_report,
)
from .progress import LiveProgress
from .scope_guard import create_phase_snapshot, enforce_phase_scope, write_change_table_markdown
from .spawner import prompt_for_spawner_command, run_spawner_command
from .verify import run_doc_comment_audit, run_go_test_json_with_timeout, run_make_with_timeout


@dataclass
class FolderResult:
    status: str
    attempts_used: int
    report_path: Path
    note: str
    fatal: bool = False


@dataclass(frozen=True)
class FolderRunPolicy:
    skip: bool
    note: str
    clean_first: bool


def _resolve_targets(cfg: Config, target: str) -> list[str]:
    if target == "all":
        return list(cfg.supported_folders)
    if target not in cfg.supported_folders:
        raise RuntimeError(f"target '{target}' is not listed in root go.work")
    return [target]


def _policy_for_folder(
    force_mode: int,
    report_status: str,
    performance_status: str,
    report_fresh: bool,
) -> FolderRunPolicy:
    if force_mode == 1:
        return FolderRunPolicy(skip=False, note="FORCE=1 always reruns with clean-first reset", clean_first=True)

    if force_mode == 2 and performance_status == "FAIL":
        return FolderRunPolicy(
            skip=False,
            note="FORCE=2 reruns because prior performance status is FAIL (improve in place)",
            clean_first=False,
        )

    if report_status == "SUCCESS" and report_fresh:
        return FolderRunPolicy(skip=True, note="fresh SUCCESS report already exists", clean_first=False)

    return FolderRunPolicy(skip=False, note="report requires rerun with clean-first reset", clean_first=True)


MODULE_PATTERN = re.compile(r"^module\s+(\S+)\s*$")
EXPECTED_MODULE_PREFIX = "github.com/jarollz/data-structure/"


def _glob_count(folder_dir: Path, pattern: str) -> int:
    return len(list(folder_dir.glob(pattern)))


def _require_folder_layout(repo_root: Path, folder: str) -> None:
    folder_dir = repo_root / folder
    if not folder_dir.is_dir():
        raise RuntimeError(f"folder does not exist: {folder}")

    for required_file in ("go.mod", "SPECS.md", "api_contract.go", "README.md"):
        required_path = folder_dir / required_file
        if not required_path.is_file():
            raise RuntimeError(f"missing required file {folder}/{required_file}")

    expected_module = EXPECTED_MODULE_PREFIX + folder
    module_line = None
    for line in (folder_dir / "go.mod").read_text(encoding="utf-8", errors="replace").splitlines():
        match = MODULE_PATTERN.match(line.strip())
        if match:
            module_line = match.group(1)
            break

    if module_line is None:
        raise RuntimeError(f"missing module declaration in {folder}/go.mod")
    if module_line != expected_module:
        raise RuntimeError(
            f"invalid module path in {folder}/go.mod: expected 'module {expected_module}'"
        )


def _warn_on_common_folder_conventions(repo_root: Path, folder: str) -> None:
    folder_dir = repo_root / folder
    for required_name in ("helpers_test.go", "bench_policy_test.go"):
        if not (folder_dir / required_name).is_file():
            warn(f"folder {folder} missing recommended file {required_name}")

    api_count = _glob_count(folder_dir, "*_api.go")
    if api_count != 1:
        warn(f"folder {folder} expected exactly one *_api.go file, found {api_count}")

    impl_count = len(list_non_test_impl_files(repo_root, folder))
    if impl_count != 1:
        warn(f"folder {folder} expected exactly one non-test implementation .go file, found {impl_count}")

    test_count = _glob_count(folder_dir, "*_test.go") - _glob_count(folder_dir, "*_bench_test.go")
    if test_count != 1:
        warn(f"folder {folder} expected exactly one *_test.go file, found {test_count}")

    bench_count = _glob_count(folder_dir, "*_bench_test.go")
    if bench_count != 1:
        warn(f"folder {folder} expected exactly one *_bench_test.go file, found {bench_count}")


def _count_lines(path: Path) -> int:
    if not path.is_file() or path.stat().st_size == 0:
        return 0
    return len(path.read_text(encoding="utf-8", errors="replace").splitlines())


def _render_template(template_file: Path, output_file: Path, replacements: dict[str, str]) -> None:
    text = template_file.read_text(encoding="utf-8")
    for key, value in replacements.items():
        text = text.replace(key, value)
    write_text(output_file, text)


def _append_failure_cause(path: Path, cause: str, scenario: str, evidence: str, suggestion: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    cause = cause.replace("\n", " ")
    scenario = scenario.replace("\n", " ")
    evidence = evidence.replace("\n", " ")
    suggestion = suggestion.replace("\n", " ")
    with path.open("a", encoding="utf-8") as fh:
        fh.write(f"{cause}\t{scenario}\t{evidence}\t{suggestion}\n")


def _read_failure_causes(path: Path) -> list[tuple[str, str, str, str]]:
    rows: list[tuple[str, str, str, str]] = []
    if not path.is_file() or path.stat().st_size == 0:
        return rows
    for line in path.read_text(encoding="utf-8", errors="replace").splitlines():
        cols = line.split("\t")
        if len(cols) != 4:
            continue
        rows.append((cols[0], cols[1], cols[2], cols[3]))
    return rows


def _write_failure_causes(path: Path, rows: list[tuple[str, str, str, str]]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    lines = [f"{cause}\t{scenario}\t{evidence}\t{suggestion}" for cause, scenario, evidence, suggestion in rows]
    write_text(path, "\n".join(lines) + ("\n" if lines else ""))


def _append_unique_failure_cause(
    path: Path,
    cause: str,
    scenario: str,
    evidence: str,
    suggestion: str,
) -> None:
    existing = _read_failure_causes(path)
    candidate = (cause, scenario, evidence, suggestion)
    if candidate in existing:
        return
    _append_failure_cause(path, cause, scenario, evidence, suggestion)


def _write_attempt_summary(
    output: Path,
    folder: str,
    attempt: int,
    max_attempts: int,
    command_exit: int,
    attempt_dir: Path,
    doc_status: str,
    test_status: str,
    bench_status: str,
) -> None:
    content = f"""# Attempt Summary

- Folder: {folder}
- Attempt: {attempt} of {max_attempts}
- AI command exit status: {command_exit}
- Protected file violations: {_count_lines(attempt_dir / 'protected_violations.tsv')}
- Out-of-scope violations: {_count_lines(attempt_dir / 'scope_violations.tsv')}
- Doc comment audit: {doc_status}
- Unit tests: {test_status}
- Benchmarks: {bench_status}

## Evidence Files

- AI output log: {attempt_dir / 'ai_output.log'}
- Protected file violations: {attempt_dir / 'protected_violations.tsv'}
- Out-of-scope violations: {attempt_dir / 'scope_violations.tsv'}
- Doc comment audit log: {attempt_dir / 'doc_comment_audit.log'}
- Unit test log: {attempt_dir / 'tests.log'}
- Benchmark log: {attempt_dir / 'benchmarks.log'}

## Retry Guidance

- Read the failing evidence files above before changing code again.
- Re-read AGENTS.md, STRUCTURE-OVERVIEW.md, {folder}/SPECS.md, {folder}/api_contract.go, and all tests.
- Keep doc comments aligned with the implementation and specs.
"""
    write_text(output, content)


def _write_report_input(
    output: Path,
    status: str,
    performance_status: str,
    folder: str,
    attempts_used: int,
    change_table_file: Path,
    unit_test_table: str,
    benchmark_table: str,
    operational_failure_causes_table: str,
    operational_failure_suggestions_table: str,
    benchmark_failure_causes_table: str,
    benchmark_failure_suggestions_table: str,
    folder_run_dir: Path,
    unit_test_log: Path,
    benchmark_log: Path,
) -> None:
    lines = [
        "# Report Facts",
        "",
        f"Operation Status: {status}",
        f"Performance Status: {performance_status}",
        f"Folder: {folder}",
        f"Attempts Used: {attempts_used}",
        "",
        "## Files Changed",
    ]
    lines.extend(change_table_file.read_text(encoding="utf-8", errors="replace").rstrip().splitlines())
    lines.append("")

    lines.extend(["## Unit Test Summary", ""])
    lines.extend(unit_test_table.splitlines())
    lines.extend(["", "## Benchmark Summary", ""])
    lines.extend(benchmark_table.splitlines())
    lines.extend(["", "## Operational Failure Causes", ""])
    lines.extend(operational_failure_causes_table.splitlines())
    lines.extend(["", "### Improvement Suggestions", ""])
    lines.extend(operational_failure_suggestions_table.splitlines())
    lines.extend(["", "## Benchmark Failure Causes", ""])
    lines.extend(benchmark_failure_causes_table.splitlines())
    lines.extend(["", "### Improvement Suggestions", ""])
    lines.extend(benchmark_failure_suggestions_table.splitlines())
    lines.append("")

    lines.extend(
        [
            "## Run Artifacts",
            "",
            f"- Folder run directory: {folder_run_dir}",
            f"- Folder summary: {folder_run_dir / 'summary.md'}",
            f"- Unit test log: {unit_test_log}",
            f"- Benchmark log: {benchmark_log}",
            "",
        ]
    )
    write_text(output, "\n".join(lines))


def _write_folder_summary(output: Path, folder: str, status: str, attempts_used: int, report_path: Path, note: str) -> None:
    content = f"""# Folder Summary

- Folder: {folder}
- Final status: {status}
- Attempts used: {attempts_used}
- Report path: {report_path}
- Note: {note}
"""
    write_text(output, content)


def _write_run_summary_markdown(summary_tsv: Path, output: Path) -> None:
    lines = ["# gen-impl Summary", "", "| Folder | Status | Attempts | Report |", "|---|---|---:|---|"]
    for row in summary_tsv.read_text(encoding="utf-8", errors="replace").splitlines():
        if not row.strip():
            continue
        folder, status, attempts, report = row.split("\t", 3)
        lines.append(f"| `{folder}` | {status} | {attempts} | `{report}` |")
    write_text(output, "\n".join(lines) + "\n")


def _print_run_summary(summary_tsv: Path) -> None:
    print("| Folder | Status | Attempts | Report |")
    print("|---|---|---:|---|")
    for row in summary_tsv.read_text(encoding="utf-8", errors="replace").splitlines():
        if not row.strip():
            continue
        folder, status, attempts, report = row.split("\t", 3)
        print(f"| `{folder}` | {status} | {attempts} | `{report}` |")


def _generate_report(
    cfg: Config,
    folder: str,
    folder_run_dir: Path,
    spawner_command: str,
    expected_status: str,
    attempts_used: int,
    report_input_file: Path,
    report_path: Path,
    progress: LiveProgress,
) -> tuple[bool, bool]:
    # returns (success, fatal)
    can_write, reason = can_write_path(report_path)
    if not can_write:
        print(
            render_report_content(
                expected_status,
                folder,
                attempts_used,
                folder_run_dir / "change_table.md",
                folder_run_dir / "operational_failure_causes.tsv",
                folder_run_dir / "benchmark_failure_causes.tsv",
            )
        )
        print(f"[gen-impl] Error: cannot write report path {report_path}: {reason}", file=sys.stderr)
        return False, True

    prompt_template = cfg.script_root / "helpers" / "prompt_report.md.tmpl"
    allowed_report_relative_path = report_path.relative_to(cfg.repo_root).as_posix()

    for report_attempt in range(1, cfg.report_attempts + 1):
        report_dir = folder_run_dir / f"report_phase_{report_attempt}"
        report_prompt = report_dir / "prompt.md"
        report_output_log = report_dir / "ai_output.log"
        report_snapshot = report_dir / "snapshot"
        validation_log = report_dir / "validation.log"

        report_dir.mkdir(parents=True, exist_ok=True)
        _render_template(
            prompt_template,
            report_prompt,
            {
                "@@FOLDER@@": folder,
                "@@REPO_ROOT@@": str(cfg.repo_root),
                "@@REPORT_INPUT_PATH@@": str(report_input_file),
                "@@OUTPUT_REPORT_PATH@@": str(report_path),
            },
        )

        create_phase_snapshot(report_snapshot, cfg.repo_root)
        progress.begin_phase(f"report ai attempt {report_attempt}/{cfg.report_attempts}", "ai")
        run_spawner_command(
            spawner_command,
            report_prompt.read_text(encoding="utf-8"),
            report_output_log,
            cfg.spawner_timeout_seconds,
            cfg.repo_root,
            idle_timeout_seconds=cfg.spawner_idle_timeout_seconds,
            progress_callback=progress.update,
        )

        phase_ok = enforce_phase_scope(
            folder,
            report_snapshot,
            allow_report_write=True,
            result_dir=report_dir,
            repo_root=cfg.repo_root,
            allowed_report_relative_path=allowed_report_relative_path,
        )
        if not phase_ok:
            warn(f"report phase {report_attempt} for {folder} edited protected or out-of-scope files; restored snapshot content")

        valid, reason = validate_generated_report(report_path, expected_status, folder, attempts_used)
        write_text(validation_log, ("ok\n" if valid else reason + "\n"))
        if valid:
            return True, False

    return False, False


def _process_folder(
    cfg: Config,
    folder: str,
    run_root: Path,
    spawner_command: str,
    summary_tsv: Path,
    progress: LiveProgress,
) -> FolderResult:
    _require_folder_layout(cfg.repo_root, folder)
    _warn_on_common_folder_conventions(cfg.repo_root, folder)
    folder_run_dir = run_root / folder
    ensure_dir_writable(folder_run_dir, f"folder run directory for {folder}")

    summary_file = folder_run_dir / "summary.md"
    prior_report_copy = folder_run_dir / "prior_IMPLEMENTATION_REPORT.md"
    prior_suggestions_file = folder_run_dir / "prior_failure_suggestions.md"
    empty_previous_summary = folder_run_dir / "no_previous_failure_summary.md"
    folder_start_snapshot = folder_run_dir / "folder_start_snapshot"
    change_table_file = folder_run_dir / "change_table.md"
    operational_failure_causes_file = folder_run_dir / "operational_failure_causes.tsv"
    benchmark_failure_causes_file = folder_run_dir / "benchmark_failure_causes.tsv"
    final_benchmark_log = folder_run_dir / "benchmarks.log"
    final_unit_test_log = folder_run_dir / "tests.log"
    report_input_file = folder_run_dir / "report_input.md"

    write_text(operational_failure_causes_file, "")
    write_text(benchmark_failure_causes_file, "")
    write_text(empty_previous_summary, "No previous attempt summary exists for this run.\n")

    last_failure_summary = empty_previous_summary
    attempts_used = 0
    status = "FAILURE"
    performance_status = "NOT_RUN"
    note = "implementation attempts exhausted"
    latest_unit_rows: list[tuple[str, str]] = []

    report_path = folder_report_path(cfg.repo_root, folder)
    report_status_current = report_status_from_file(report_path)
    report_performance_current = report_performance_status_from_file(report_path)
    report_fresh = report_is_fresh(cfg.repo_root, folder, report_path)
    prior_attempts = report_attempts_from_file(report_path)

    if report_path.is_file():
        write_text(prior_report_copy, report_path.read_text(encoding="utf-8", errors="replace"))
    else:
        write_text(prior_report_copy, "No prior implementation report found.\n")
    extract_failure_suggestions(report_path, prior_suggestions_file)

    run_policy = _policy_for_folder(
        cfg.force_mode,
        report_status_current,
        report_performance_current,
        report_fresh,
    )

    if run_policy.skip:
        status = "SKIPPED"
        attempts_used = prior_attempts
        note = run_policy.note
        _write_folder_summary(summary_file, folder, status, attempts_used, report_path, note)
        with summary_tsv.open("a", encoding="utf-8") as fh:
            fh.write(f"{folder}\t{status}\t{attempts_used}\t{report_path}\n")
        return FolderResult(status, attempts_used, report_path, note)

    create_phase_snapshot(folder_start_snapshot, cfg.repo_root)

    need_reset = run_policy.clean_first and folder_has_non_test_impl_files(cfg.repo_root, folder)
    if need_reset:
        reset_dir = folder_run_dir / "reset_phase"
        reset_prompt = reset_dir / "prompt.md"
        reset_output = reset_dir / "ai_output.log"
        reset_snapshot = reset_dir / "snapshot"

        reset_dir.mkdir(parents=True, exist_ok=True)
        _render_template(
            cfg.script_root / "helpers" / "prompt_reset.md.tmpl",
            reset_prompt,
            {
                "@@FOLDER@@": folder,
                "@@REPO_ROOT@@": str(cfg.repo_root),
                "@@PRIOR_REPORT_PATH@@": str(prior_report_copy),
                "@@PRIOR_SUGGESTIONS_PATH@@": str(prior_suggestions_file),
            },
        )

        create_phase_snapshot(reset_snapshot, cfg.repo_root)
        progress.begin_phase("reset ai run", "ai")
        if not run_spawner_command(
            spawner_command,
            reset_prompt.read_text(encoding="utf-8"),
            reset_output,
            cfg.spawner_timeout_seconds,
            cfg.repo_root,
            idle_timeout_seconds=cfg.spawner_idle_timeout_seconds,
            progress_callback=progress.update,
        ):
            _append_failure_cause(
                operational_failure_causes_file,
                "Reset phase failed",
                "(reset)",
                str(reset_output),
                "Check the reset prompt and spawner command, then rerun the script.",
            )

        if not enforce_phase_scope(
            folder,
            reset_snapshot,
            allow_report_write=False,
            result_dir=reset_dir,
            repo_root=cfg.repo_root,
        ):
            _append_failure_cause(
                operational_failure_causes_file,
                "Reset phase touched protected or out-of-scope files",
                "(reset)",
                f"{reset_dir / 'protected_violations.tsv'} and {reset_dir / 'scope_violations.tsv'}",
                "Keep reset changes inside target implementation files only.",
            )

    for attempt in range(1, cfg.max_attempts + 1):
        attempt_dir = folder_run_dir / f"attempt_{attempt}"
        prompt_file = attempt_dir / "prompt.md"
        output_log = attempt_dir / "ai_output.log"
        phase_snapshot = attempt_dir / "snapshot"
        doc_log = attempt_dir / "doc_comment_audit.log"
        test_log = attempt_dir / "tests.log"
        command_exit = 0
        phase_ok = True
        doc_status = "SKIPPED"
        test_status = "SKIPPED"
        bench_status = "NOT_RUN"

        attempt_dir.mkdir(parents=True, exist_ok=True)
        _render_template(
            cfg.script_root / "helpers" / "prompt_impl.md.tmpl",
            prompt_file,
            {
                "@@FOLDER@@": folder,
                "@@REPO_ROOT@@": str(cfg.repo_root),
                "@@ATTEMPT@@": str(attempt),
                "@@MAX_ATTEMPTS@@": str(cfg.max_attempts),
                "@@PRIOR_REPORT_PATH@@": str(prior_report_copy),
                "@@PRIOR_SUGGESTIONS_PATH@@": str(prior_suggestions_file),
                "@@FAILURE_SUMMARY_PATH@@": str(last_failure_summary),
                "@@RUN_DIR@@": str(folder_run_dir),
            },
        )

        create_phase_snapshot(phase_snapshot, cfg.repo_root)
        progress.begin_phase(f"impl ai attempt {attempt}/{cfg.max_attempts}", "ai")
        if run_spawner_command(
            spawner_command,
            prompt_file.read_text(encoding="utf-8"),
            output_log,
            cfg.spawner_timeout_seconds,
            cfg.repo_root,
            idle_timeout_seconds=cfg.spawner_idle_timeout_seconds,
            progress_callback=progress.update,
        ):
            command_exit = 0
        else:
            command_exit = 1
            phase_ok = False
            _append_failure_cause(
                operational_failure_causes_file,
                "AI command exited non-zero",
                "(implementation)",
                str(output_log),
                "Fix the spawner command or prompt so the AI run finishes successfully.",
            )

        if not enforce_phase_scope(
            folder,
            phase_snapshot,
            allow_report_write=False,
            result_dir=attempt_dir,
            repo_root=cfg.repo_root,
        ):
            phase_ok = False
            if _count_lines(attempt_dir / "protected_violations.tsv") > 0:
                _append_failure_cause(
                    operational_failure_causes_file,
                    "Protected files were modified",
                    "(scope)",
                    str(attempt_dir / "protected_violations.tsv"),
                    "Do not edit tests, specs, docs, api_contract.go, go.mod, or go.sum.",
                )
            if _count_lines(attempt_dir / "scope_violations.tsv") > 0:
                _append_failure_cause(
                    operational_failure_causes_file,
                    "Out-of-scope files were modified",
                    "(scope)",
                    str(attempt_dir / "scope_violations.tsv"),
                    "Keep all edits inside the target folder implementation files only.",
                )

        if phase_ok:
            progress.begin_phase(f"doc audit attempt {attempt}/{cfg.max_attempts}", "log")
            if run_doc_comment_audit(
                cfg.repo_root,
                cfg.script_root,
                folder,
                doc_log,
                cfg.doc_audit_timeout_seconds,
                progress_callback=progress.update,
            ):
                doc_status = "PASS"
            else:
                phase_ok = False
                doc_status = "FAIL"
                _append_failure_cause(
                    operational_failure_causes_file,
                    "Doc comment audit failed",
                    "(doc audit)",
                    str(doc_log),
                    "Align exported function comments with SPECS.md, API header rules, and actual implementation behavior.",
                )

        if phase_ok:
            progress.begin_phase(f"unit tests attempt {attempt}/{cfg.max_attempts}", "log")
            unit_ok = run_go_test_json_with_timeout(
                cfg.repo_root,
                folder,
                cfg.test_timeout_seconds,
                test_log,
                progress_callback=progress.update,
            )
            latest_unit_rows = parse_unit_test_scenarios_from_json_log(test_log)
            if test_log.is_file():
                write_text(final_unit_test_log, test_log.read_text(encoding="utf-8", errors="replace"))
            if unit_ok:
                test_status = "PASS"
            else:
                phase_ok = False
                test_status = "FAIL"
                runtime_row = runtime_failure_cause_from_log(test_log, "operational")
                if runtime_row is not None:
                    _append_unique_failure_cause(operational_failure_causes_file, *runtime_row)
                else:
                    _append_failure_cause(
                        operational_failure_causes_file,
                        "Unit tests failed",
                        "(unit tests)",
                        str(test_log),
                        "Read the failing tests carefully and align implementation behavior with the fixed test contract.",
                    )

        _write_attempt_summary(
            attempt_dir / "summary.md",
            folder,
            attempt,
            cfg.max_attempts,
            command_exit,
            attempt_dir,
            doc_status,
            test_status,
            bench_status,
        )
        last_failure_summary = attempt_dir / "summary.md"
        attempts_used = attempt

        if phase_ok:
            status = "SUCCESS"
            note = "implementation and verification succeeded"
            break
        note = "last attempt failed"

    benchmark_run_ok = False
    benchmark_rows: list[tuple[str, str, str, str]] = []
    if status == "SUCCESS":
        progress.begin_phase("bench checks final", "log")
        benchmark_run_ok = run_make_with_timeout(
            cfg.repo_root,
            cfg.bench_timeout_seconds,
            final_benchmark_log,
            "bench-folder",
            folder,
            progress_callback=progress.update,
        )
        benchmark_rows = parse_benchmark_rows_from_log(final_benchmark_log)
        if not benchmark_rows:
            benchmark_rows = [("(benchmark command)", "N/A", "BAD", "FAIL")]
        if benchmark_run_ok and all(status_row == "PASS" for _, _, _, status_row in benchmark_rows):
            performance_status = "PASS"
        else:
            performance_status = "FAIL"
    else:
        write_text(final_benchmark_log, "Benchmark run skipped because operational status is FAILURE.\n")
        benchmark_rows = [("(not run)", "N/A", "BAD", "FAIL")]
        performance_status = "NOT_RUN"

    if not latest_unit_rows:
        if status == "SUCCESS":
            latest_unit_rows = [("(none)", "PASS")]
        else:
            latest_unit_rows = [("(not run)", "FAIL")]
    if not final_unit_test_log.is_file():
        write_text(final_unit_test_log, "Unit tests did not run in this folder processing flow.\n")

    benchmark_causes = benchmark_failure_causes(
        benchmark_rows,
        final_benchmark_log,
        benchmark_run_ok,
        performance_status,
    )
    _write_failure_causes(benchmark_failure_causes_file, benchmark_causes)

    operational_rows = _read_failure_causes(operational_failure_causes_file)
    benchmark_rows_for_causes = _read_failure_causes(benchmark_failure_causes_file)

    unit_test_table = render_scenario_table(latest_unit_rows)
    benchmark_table = render_benchmark_table(benchmark_rows)
    operational_failure_causes_table = render_failure_causes_table(operational_rows)
    operational_failure_suggestions_table = render_failure_improvement_table(
        failure_improvement_suggestions_from_causes(operational_rows)
    )
    benchmark_failure_causes_table = render_failure_causes_table(benchmark_rows_for_causes)
    benchmark_failure_suggestions_table = render_failure_improvement_table(
        failure_improvement_suggestions_from_causes(benchmark_rows_for_causes)
    )

    write_change_table_markdown(folder, folder_start_snapshot, cfg.repo_root, change_table_file)
    _write_report_input(
        report_input_file,
        status,
        performance_status,
        folder,
        attempts_used,
        change_table_file,
        unit_test_table,
        benchmark_table,
        operational_failure_causes_table,
        operational_failure_suggestions_table,
        benchmark_failure_causes_table,
        benchmark_failure_suggestions_table,
        folder_run_dir,
        final_unit_test_log,
        final_benchmark_log,
    )

    report_ok, report_fatal = _generate_report(
        cfg,
        folder,
        folder_run_dir,
        spawner_command,
        status,
        attempts_used,
        report_input_file,
        report_path,
        progress,
    )
    if report_fatal:
        return FolderResult("FAILURE", attempts_used, report_path, "report path not writable", fatal=True)

    if not report_ok:
        status = "FAILURE"
        note = "report generation failed"
        _append_failure_cause(
            operational_failure_causes_file,
            "Implementation report generation failed",
            "(report)",
            str(folder_run_dir / "report_phase_*/validation.log"),
            "Fix the report prompt or spawner command, then rerun the script.",
        )

    _write_folder_summary(summary_file, folder, status, attempts_used, report_path, note)
    with summary_tsv.open("a", encoding="utf-8") as fh:
        fh.write(f"{folder}\t{status}\t{attempts_used}\t{report_path}\n")
    return FolderResult(status, attempts_used, report_path, note)


def run(cfg: Config, target: str) -> int:
    run_root = cfg.repo_root / "tmp" / "gen-impl" / "runs" / timestamp_utc()
    reports_root = cfg.repo_root / "tmp" / "gen-impl" / "reports"
    ensure_dir_writable(run_root, "run root")
    ensure_dir_writable(reports_root, "reports root")

    targets = _resolve_targets(cfg, target)
    spawner_command = prompt_for_spawner_command(
        run_root,
        cfg.repo_root,
        cfg.probe_timeout_seconds,
        cfg.spawner_idle_timeout_seconds,
        cfg.ai_spawner_command,
    )

    summary_tsv = run_root / "summary.tsv"
    write_text(summary_tsv, "")

    overall_status = 0
    for folder in targets:
        log(f"Processing {folder}")
        progress = LiveProgress()
        try:
            result = _process_folder(cfg, folder, run_root, spawner_command, summary_tsv, progress)
        finally:
            progress.clear()
        log(f"Done {folder} ({result.status}, attempts={result.attempts_used})")
        if result.fatal:
            overall_status = 1
            break
        if result.status == "FAILURE":
            overall_status = 1
            if target == "all" and cfg.stop_on_failure:
                break

    summary_md = run_root / "summary.md"
    _write_run_summary_markdown(summary_tsv, summary_md)
    _print_run_summary(summary_tsv)
    log(f"Run artifacts: {run_root}")
    return overall_status
