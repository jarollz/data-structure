from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
import sys

from .config import Config, SUPPORTED_FOLDERS
from .io_safe import can_write_path, ensure_dir_writable, log, warn, write_text, timestamp_utc
from .reporting import (
    extract_failure_suggestions,
    folder_has_non_test_impl_files,
    folder_report_path,
    render_report_content,
    report_attempts_from_file,
    report_is_fresh,
    report_status_from_file,
    validate_generated_report,
)
from .progress import LiveProgress
from .scope_guard import create_phase_snapshot, enforce_phase_scope, write_change_table_markdown
from .spawner import prompt_for_spawner_command, run_spawner_command
from .verify import run_doc_comment_audit, run_make_with_timeout


@dataclass
class FolderResult:
    status: str
    attempts_used: int
    report_path: Path
    note: str
    fatal: bool = False


def _resolve_targets(target: str) -> list[str]:
    if target == "all":
        return list(SUPPORTED_FOLDERS)
    if target not in SUPPORTED_FOLDERS:
        raise RuntimeError(f"unsupported folder '{target}'")
    return [target]


def _require_folder_layout(repo_root: Path, folder: str) -> None:
    folder_dir = repo_root / folder
    if not folder_dir.is_dir():
        raise RuntimeError(f"folder does not exist: {folder}")
    if not (folder_dir / "go.mod").is_file():
        raise RuntimeError(f"missing go.mod in folder: {folder}")


def _count_lines(path: Path) -> int:
    if not path.is_file() or path.stat().st_size == 0:
        return 0
    return len(path.read_text(encoding="utf-8", errors="replace").splitlines())


def _render_template(template_file: Path, output_file: Path, replacements: dict[str, str]) -> None:
    text = template_file.read_text(encoding="utf-8")
    for key, value in replacements.items():
        text = text.replace(key, value)
    write_text(output_file, text)


def _append_failure_cause(path: Path, cause: str, evidence: str, suggestion: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    cause = cause.replace("\n", " ")
    evidence = evidence.replace("\n", " ")
    suggestion = suggestion.replace("\n", " ")
    with path.open("a", encoding="utf-8") as fh:
        fh.write(f"{cause}\t{evidence}\t{suggestion}\n")


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
    folder: str,
    attempts_used: int,
    change_table_file: Path,
    failure_causes_file: Path,
    folder_run_dir: Path,
) -> None:
    lines = [
        "# Report Facts",
        "",
        f"Operation Status: {status}",
        f"Folder: {folder}",
        f"Attempts Used: {attempts_used}",
        "",
        "## Files Changed",
    ]
    lines.extend(change_table_file.read_text(encoding="utf-8", errors="replace").rstrip().splitlines())
    lines.append("")

    if status == "FAILURE":
        lines.extend(["## Failure Causes", "", "| Cause | Evidence | Suggestion |", "|---|---|---|"])
        if failure_causes_file.is_file() and failure_causes_file.stat().st_size > 0:
            for row in failure_causes_file.read_text(encoding="utf-8", errors="replace").splitlines():
                cols = row.split("\t")
                if len(cols) != 3:
                    continue
                lines.append(f"| {cols[0]} | {cols[1]} | {cols[2]} |")
        else:
            lines.append(
                f"| Unknown failure | {folder_run_dir / 'summary.md'} | Review the latest attempt logs and rerun after fixing the root cause. |"
            )
        lines.append("")

    lines.extend(
        [
            "## Run Artifacts",
            "",
            f"- Folder run directory: {folder_run_dir}",
            f"- Folder summary: {folder_run_dir / 'summary.md'}",
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
                folder_run_dir / "failure_causes.tsv",
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
    folder_run_dir = run_root / folder
    ensure_dir_writable(folder_run_dir, f"folder run directory for {folder}")

    summary_file = folder_run_dir / "summary.md"
    prior_report_copy = folder_run_dir / "prior_IMPLEMENTATION_REPORT.md"
    prior_suggestions_file = folder_run_dir / "prior_failure_suggestions.md"
    empty_previous_summary = folder_run_dir / "no_previous_failure_summary.md"
    folder_start_snapshot = folder_run_dir / "folder_start_snapshot"
    change_table_file = folder_run_dir / "change_table.md"
    failure_causes_file = folder_run_dir / "failure_causes.tsv"
    report_input_file = folder_run_dir / "report_input.md"
    last_failure_summary = empty_previous_summary
    attempts_used = 0
    status = "FAILURE"
    note = "implementation attempts exhausted"

    write_text(failure_causes_file, "")
    write_text(empty_previous_summary, "No previous attempt summary exists for this run.\n")

    report_path = folder_report_path(cfg.repo_root, folder)
    report_status_current = report_status_from_file(report_path)
    prior_attempts = report_attempts_from_file(report_path)

    if report_path.is_file():
        write_text(prior_report_copy, report_path.read_text(encoding="utf-8", errors="replace"))
    else:
        write_text(prior_report_copy, "No prior implementation report found.\n")
    extract_failure_suggestions(report_path, prior_suggestions_file)

    if not cfg.force and report_status_current == "SUCCESS" and report_is_fresh(cfg.repo_root, folder, report_path):
        status = "SKIPPED"
        attempts_used = prior_attempts
        note = "fresh SUCCESS report already exists"
        _write_folder_summary(summary_file, folder, status, attempts_used, report_path, note)
        with summary_tsv.open("a", encoding="utf-8") as fh:
            fh.write(f"{folder}\t{status}\t{attempts_used}\t{report_path}\n")
        return FolderResult(status, attempts_used, report_path, note)

    create_phase_snapshot(folder_start_snapshot, cfg.repo_root)

    need_reset = False
    if report_status_current in {"FAILURE", "INVALID", "MISSING"} and folder_has_non_test_impl_files(cfg.repo_root, folder):
        need_reset = True

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
            progress_callback=progress.update,
        ):
            _append_failure_cause(
                failure_causes_file,
                "Reset phase failed",
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
                failure_causes_file,
                "Reset phase touched protected or out-of-scope files",
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
        bench_log = attempt_dir / "benchmarks.log"
        command_exit = 0
        phase_ok = True
        doc_status = "SKIPPED"
        test_status = "SKIPPED"
        bench_status = "SKIPPED"

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
            progress_callback=progress.update,
        ):
            command_exit = 0
        else:
            command_exit = 1
            phase_ok = False
            _append_failure_cause(
                failure_causes_file,
                "AI command exited non-zero",
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
                    failure_causes_file,
                    "Protected files were modified",
                    str(attempt_dir / "protected_violations.tsv"),
                    "Do not edit tests, specs, docs, api_contract.go, go.mod, or go.sum.",
                )
            if _count_lines(attempt_dir / "scope_violations.tsv") > 0:
                _append_failure_cause(
                    failure_causes_file,
                    "Out-of-scope files were modified",
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
                    failure_causes_file,
                    "Doc comment audit failed",
                    str(doc_log),
                    "Align exported function comments with SPECS.md, API header rules, and actual implementation behavior.",
                )

        if phase_ok:
            progress.begin_phase(f"unit tests attempt {attempt}/{cfg.max_attempts}", "log")
            if run_make_with_timeout(
                cfg.repo_root,
                cfg.test_timeout_seconds,
                test_log,
                "test-folder",
                folder,
                progress_callback=progress.update,
            ):
                test_status = "PASS"
            else:
                phase_ok = False
                test_status = "FAIL"
                _append_failure_cause(
                    failure_causes_file,
                    "Unit tests failed",
                    str(test_log),
                    "Read the failing tests carefully and align implementation behavior with the fixed test contract.",
                )

        if phase_ok:
            progress.begin_phase(f"bench checks attempt {attempt}/{cfg.max_attempts}", "log")
            if run_make_with_timeout(
                cfg.repo_root,
                cfg.bench_timeout_seconds,
                bench_log,
                "bench-folder",
                folder,
                progress_callback=progress.update,
            ):
                bench_status = "PASS"
            else:
                phase_ok = False
                bench_status = "FAIL"
                _append_failure_cause(
                    failure_causes_file,
                    "Benchmark checks failed",
                    str(bench_log),
                    "Review benchmark failures and adjust implementation until benchmark validation passes.",
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

    write_change_table_markdown(folder, folder_start_snapshot, cfg.repo_root, change_table_file)
    _write_report_input(report_input_file, status, folder, attempts_used, change_table_file, failure_causes_file, folder_run_dir)

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
            failure_causes_file,
            "Implementation report generation failed",
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

    targets = _resolve_targets(target)
    spawner_command = prompt_for_spawner_command(
        run_root,
        cfg.repo_root,
        cfg.probe_timeout_seconds,
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
