from __future__ import annotations

import json
import re
from pathlib import Path


def folder_report_path(repo_root: Path, folder: str) -> Path:
    return repo_root / "tmp" / "gen-impl" / "reports" / folder / "IMPLEMENTATION_REPORT.md"


def report_field_from_file(report_file: Path, field_name: str) -> str | None:
    if not report_file.is_file():
        return None
    prefix = f"{field_name}: "
    for line in report_file.read_text(encoding="utf-8", errors="replace").splitlines():
        if line.startswith(prefix):
            return line[len(prefix) :].strip()
    return None


def report_status_from_file(report_file: Path) -> str:
    if not report_file.is_file():
        return "MISSING"
    status = report_field_from_file(report_file, "Operation Status")
    if status in {"SUCCESS", "FAILURE"}:
        return status
    return "INVALID"


def report_performance_status_from_file(report_file: Path) -> str:
    if not report_file.is_file():
        return "MISSING"
    status = report_field_from_file(report_file, "Performance Status")
    if status in {"PASS", "FAIL", "NOT_RUN"}:
        return status
    return "INVALID"


def report_attempts_from_file(report_file: Path) -> int:
    if not report_file.is_file():
        return 0
    value = report_field_from_file(report_file, "Attempts Used")
    if value is None or not value.isdigit():
        return 0
    return int(value)


def list_non_test_impl_files(repo_root: Path, folder: str) -> list[Path]:
    folder_dir = repo_root / folder
    result: list[Path] = []
    for file in sorted(folder_dir.glob("*.go")):
        base = file.name
        if base in {"api_contract.go", "helpers_test.go", "bench_policy_test.go"}:
            continue
        if base.endswith("_test.go") or base.endswith("_bench_test.go"):
            continue
        result.append(file)
    return result


def folder_has_non_test_impl_files(repo_root: Path, folder: str) -> bool:
    return len(list_non_test_impl_files(repo_root, folder)) > 0


def report_is_fresh(repo_root: Path, folder: str, report_file: Path) -> bool:
    if not report_file.is_file():
        return False
    report_mtime = report_file.stat().st_mtime
    for file in list_non_test_impl_files(repo_root, folder):
        if file.stat().st_mtime > report_mtime:
            return False
    return True


def extract_failure_suggestions(report_file: Path, output_file: Path) -> None:
    if not report_file.is_file():
        output_file.write_text("No prior implementation report exists.\n", encoding="utf-8")
        return

    lines = report_file.read_text(encoding="utf-8", errors="replace").splitlines()
    capture = False
    out_lines: list[str] = []
    section_titles = {"## Operational Failure Causes", "## Benchmark Failure Causes", "## Failure Causes"}
    for line in lines:
        if line in section_titles:
            capture = True
        if capture and line.startswith("## ") and line not in section_titles:
            break
        if capture:
            out_lines.append(line)

    if not out_lines:
        output_file.write_text(
            f"No prior failure suggestions were found in {report_file}.\n",
            encoding="utf-8",
        )
        return
    output_file.write_text("\n".join(out_lines) + "\n", encoding="utf-8")


def validate_generated_report(report_file: Path, expected_status: str, expected_folder: str, expected_attempts: int) -> tuple[bool, str]:
    if not report_file.is_file():
        return False, f"missing report file: {report_file}"
    if report_status_from_file(report_file) != expected_status:
        return False, f"unexpected report status in {report_file}"
    folder_value = report_field_from_file(report_file, "Folder")
    if folder_value != expected_folder:
        return False, f"unexpected folder field in {report_file}"
    attempts = report_attempts_from_file(report_file)
    if attempts != expected_attempts:
        return False, f"unexpected attempts field in {report_file}"
    text = report_file.read_text(encoding="utf-8", errors="replace")
    if "| File | Change |" not in text:
        return False, f"missing file change table in {report_file}"
    required_headers = (
        "## Unit Test Summary",
        "| no. | scenario | pass / fail |",
        "## Benchmark Summary",
        "| no. | scenario | budget-ratio | good/bad | pass / fail |",
        "## Operational Failure Causes",
        "### Improvement Suggestions",
        "| no. | cause | failed scenario | suggestion |",
        "| no. | cause | scenario | evidence | suggestion |",
        "## Benchmark Failure Causes",
    )
    for header in required_headers:
        if header not in text:
            return False, f"missing required report section/header '{header}' in {report_file}"
    if text.count("### Improvement Suggestions") < 2:
        return False, f"missing failure-cause improvement suggestions subsections in {report_file}"
    return True, ""


def render_report_content(
    status: str,
    folder: str,
    attempts_used: int,
    change_table_file: Path,
    operational_failure_causes_file: Path,
    benchmark_failure_causes_file: Path,
) -> str:
    parts = [
        "# IMPLEMENTATION_REPORT",
        "",
        f"Operation Status: {status}",
        "Performance Status: NOT_RUN",
        f"Folder: {folder}",
        f"Attempts Used: {attempts_used}",
        "",
        "## Files Changed",
    ]
    parts.extend(change_table_file.read_text(encoding="utf-8", errors="replace").strip().splitlines())
    parts.append("")

    parts.extend(
        [
            "## Unit Test Summary",
            "",
            "| no. | scenario | pass / fail |",
            "|---:|---|---|",
            "| 1 | (none) | PASS |",
            "",
            "## Benchmark Summary",
            "",
            "| no. | scenario | budget-ratio | good/bad | pass / fail |",
            "|---:|---|---:|---|---|",
            "| 1 | (not run) | N/A | BAD | FAIL |",
            "",
            "## Operational Failure Causes",
            "",
            "| no. | cause | scenario | evidence | suggestion |",
            "|---:|---|---|---|---|",
        ]
    )
    parts.extend(_failure_causes_rows(operational_failure_causes_file))
    parts.extend(
        [
            "",
            "### Improvement Suggestions",
            "",
            "| no. | cause | failed scenario | suggestion |",
            "|---:|---|---|---|",
        ]
    )
    parts.extend(_failure_improvement_rows(operational_failure_causes_file))
    parts.extend(
        [
            "",
            "## Benchmark Failure Causes",
            "",
            "| no. | cause | scenario | evidence | suggestion |",
            "|---:|---|---|---|---|",
        ]
    )
    parts.extend(_failure_causes_rows(benchmark_failure_causes_file))
    parts.extend(
        [
            "",
            "### Improvement Suggestions",
            "",
            "| no. | cause | failed scenario | suggestion |",
            "|---:|---|---|---|",
        ]
    )
    parts.extend(_failure_improvement_rows(benchmark_failure_causes_file))
    parts.append("")

    return "\n".join(parts).rstrip() + "\n"


def _md_cell(value: str) -> str:
    return value.replace("|", "\\|").replace("\n", " ").strip()


def render_scenario_table(rows: list[tuple[str, str]]) -> str:
    out = ["| no. | scenario | pass / fail |", "|---:|---|---|"]
    if not rows:
        rows = [("(none)", "PASS")]
    for i, (scenario, status) in enumerate(rows, start=1):
        out.append(f"| {i} | {_md_cell(scenario)} | {_md_cell(status.upper())} |")
    return "\n".join(out)


def render_suggestion_table(rows: list[tuple[str, str]], none_text: str) -> str:
    out = ["| no. | failed scenario | suggestion |", "|---:|---|---|"]
    if not rows:
        rows = [("(none)", none_text)]
    for i, (scenario, suggestion) in enumerate(rows, start=1):
        out.append(f"| {i} | {_md_cell(scenario)} | {_md_cell(suggestion)} |")
    return "\n".join(out)


def render_failure_improvement_table(rows: list[tuple[str, str, str]]) -> str:
    out = ["| no. | cause | failed scenario | suggestion |", "|---:|---|---|---|"]
    if not rows:
        rows = [("(none)", "(none)", "No failures recorded.")]
    for i, (cause, scenario, suggestion) in enumerate(rows, start=1):
        out.append(f"| {i} | {_md_cell(cause)} | {_md_cell(scenario)} | {_md_cell(suggestion)} |")
    return "\n".join(out)


def render_benchmark_table(rows: list[tuple[str, str, str, str]]) -> str:
    out = [
        "| no. | scenario | budget-ratio | good/bad | pass / fail |",
        "|---:|---|---:|---|---|",
    ]
    if not rows:
        rows = [("(none)", "N/A", "GOOD", "PASS")]
    for i, (scenario, ratio, good_bad, status) in enumerate(rows, start=1):
        out.append(
            f"| {i} | {_md_cell(scenario)} | {_md_cell(ratio)} | {_md_cell(good_bad)} | {_md_cell(status)} |"
        )
    return "\n".join(out)


def render_failure_causes_table(rows: list[tuple[str, str, str, str]]) -> str:
    out = ["| no. | cause | scenario | evidence | suggestion |", "|---:|---|---|---|---|"]
    if not rows:
        rows = [("(none)", "(none)", "(none)", "No failures recorded.")]
    for i, (cause, scenario, evidence, suggestion) in enumerate(rows, start=1):
        out.append(
            f"| {i} | {_md_cell(cause)} | {_md_cell(scenario)} | {_md_cell(evidence)} | {_md_cell(suggestion)} |"
        )
    return "\n".join(out)


def parse_unit_test_scenarios_from_json_log(log_file: Path) -> list[tuple[str, str]]:
    scenarios: dict[str, str] = {}
    if not log_file.is_file():
        return []
    for line in log_file.read_text(encoding="utf-8", errors="replace").splitlines():
        text = line.strip()
        if not text.startswith("{"):
            continue
        try:
            payload = json.loads(text)
        except json.JSONDecodeError:
            continue
        scenario = payload.get("Test")
        action = payload.get("Action")
        if not isinstance(scenario, str):
            continue
        if action == "pass":
            scenarios[scenario] = "PASS"
        elif action == "fail":
            scenarios[scenario] = "FAIL"

    rows = sorted(scenarios.items(), key=lambda item: item[0])
    return rows


BENCHMARK_RATIO_RE = re.compile(
    r"^(Benchmark\S+)\s+.*?([0-9]+(?:\.[0-9]+)?(?:[eE][+-]?[0-9]+)?)\s+budget-ratio\b"
)


def parse_benchmark_rows_from_log(log_file: Path) -> list[tuple[str, str, str, str]]:
    rows: list[tuple[str, str, str, str]] = []
    if not log_file.is_file():
        return rows
    for line in log_file.read_text(encoding="utf-8", errors="replace").splitlines():
        match = BENCHMARK_RATIO_RE.match(line.strip())
        if not match:
            continue
        scenario = match.group(1)
        try:
            ratio_value = float(match.group(2))
        except ValueError:
            continue
        ratio = f"{ratio_value:.4f}"
        good_bad = "GOOD" if ratio_value <= 1.0 else "BAD"
        status = "PASS" if ratio_value <= 1.0 else "FAIL"
        rows.append((scenario, ratio, good_bad, status))
    return rows


def unit_test_improvement_suggestions(unit_rows: list[tuple[str, str]]) -> list[tuple[str, str]]:
    failed = [scenario for scenario, status in unit_rows if status == "FAIL"]
    return [(scenario, _unit_suggestion_for_scenario(scenario)) for scenario in failed]


def benchmark_improvement_suggestions(bench_rows: list[tuple[str, str, str, str]]) -> list[tuple[str, str]]:
    failed = [scenario for scenario, _, _, status in bench_rows if status == "FAIL"]
    return [(scenario, _benchmark_suggestion_for_scenario(scenario)) for scenario in failed]


def benchmark_failure_causes(
    bench_rows: list[tuple[str, str, str, str]],
    bench_log: Path,
    bench_run_ok: bool,
    performance_status: str,
) -> list[tuple[str, str, str, str]]:
    causes: list[tuple[str, str, str, str]] = []
    folder = _infer_folder_from_path(bench_log)
    for scenario, ratio, _, status in bench_rows:
        if status != "FAIL":
            continue
        if ratio.strip().upper() == "N/A":
            continue
        causes.append(
            (
                "Benchmark threshold exceeded",
                scenario,
                f"{bench_log} (budget-ratio={ratio})",
                _benchmark_suggestion_for_scenario(scenario, folder),
            )
        )

    if not bench_run_ok:
        runtime_row = runtime_failure_cause_from_log(bench_log, "benchmark")
        if runtime_row is not None:
            causes.append(runtime_row)
        else:
            causes.append(
                (
                    "Benchmark command failed",
                    "(benchmark command)",
                    str(bench_log),
                    _format_retry_suggestion(
                        "Benchmark command exited non-zero",
                        "benchmark command path",
                        "inspect first failing line in benchmarks.log and patch referenced implementation path",
                        _benchmark_verify_command(folder),
                        "benchmark command exits 0 without panic",
                    ),
                )
            )
    elif performance_status == "NOT_RUN":
        causes.append(
            (
                "Benchmark not run",
                "(benchmark command)",
                str(bench_log),
                _format_retry_suggestion(
                    "Benchmark phase skipped",
                    "operational verification flow",
                    "fix operational failures first, then execute benchmark phase",
                    _benchmark_verify_command(folder),
                    "benchmark command exits 0 and emits budget-ratio rows",
                ),
            )
        )

    return causes


def _unit_suggestion_for_scenario(scenario: str) -> str:
    lower = scenario.lower()
    if "clone" in lower:
        return "Recheck clone contract: copy live elements only, preserve container independence, and apply clone hooks consistently."
    if "iter" in lower or "all" in lower or "inorder" in lower or "preorder" in lower or "values" in lower:
        return "Recheck iterator contract: order, count, early-stop behavior, and mutation-unsafety note."
    if "put" in lower or "set" in lower or "insert" in lower:
        return "Recheck insert/update semantics, overwrite behavior, and size accounting."
    if "get" in lower or "has" in lower or "find" in lower:
        return "Recheck lookup semantics for hit/miss paths and key-equality checks."
    if "delete" in lower or "remove" in lower:
        return "Recheck delete semantics, tombstone/state cleanup, and size accounting after removal."
    if "clear" in lower:
        return "Recheck clear behavior: reset internal state, length, and capacity/load-factor contracts."
    if "len" in lower or "cap" in lower or "size" in lower:
        return "Recheck Len/Cap related contract values after each mutating operation."
    return "Recheck this scenario against SPECS.md and tests; align implementation behavior with expected contract."


def _benchmark_suggestion_for_scenario(scenario: str, folder: str | None = None) -> str:
    lower = scenario.lower()
    code_target = f"{scenario} hot path"
    verify_command = _benchmark_verify_command(folder)
    if "put" in lower or "set" in lower or "insert" in lower:
        return _format_retry_suggestion(
            "Benchmark threshold exceeded",
            code_target,
            "reduce redundant probes/writes on insert path and trigger rehash only at load-factor boundary",
            verify_command,
            "budget-ratio <= 1.0 and benchmark exits 0",
        )
    if "get" in lower or "has" in lower or "find" in lower:
        return _format_retry_suggestion(
            "Benchmark threshold exceeded",
            code_target,
            "reduce branch/probe work on hit and miss lookup paths and avoid repeated hash computations",
            verify_command,
            "budget-ratio <= 1.0 and benchmark exits 0",
        )
    if "delete" in lower or "remove" in lower:
        return _format_retry_suggestion(
            "Benchmark threshold exceeded",
            code_target,
            "reduce probe-chain cleanup cost and avoid redundant tombstone/state writes on delete",
            verify_command,
            "budget-ratio <= 1.0 and benchmark exits 0",
        )
    if "clone" in lower:
        return _format_retry_suggestion(
            "Benchmark threshold exceeded",
            code_target,
            "pre-size destination storage and copy only live entries with minimal per-entry allocations",
            verify_command,
            "budget-ratio <= 1.0 and benchmark exits 0",
        )
    if "clear" in lower:
        return _format_retry_suggestion(
            "Benchmark threshold exceeded",
            code_target,
            "reset internal state in bulk and avoid per-entry heavy work during clear",
            verify_command,
            "budget-ratio <= 1.0 and benchmark exits 0",
        )
    if "mixed" in lower:
        return _format_retry_suggestion(
            "Benchmark threshold exceeded",
            code_target,
            "reduce resize/probe churn across alternating operations and reuse computed probe state where safe",
            verify_command,
            "budget-ratio <= 1.0 and benchmark exits 0",
        )
    if "all" in lower or "inorder" in lower or "preorder" in lower or "values" in lower:
        return _format_retry_suggestion(
            "Benchmark threshold exceeded",
            code_target,
            "skip empty/tombstone slots early and reduce per-item iterator overhead",
            verify_command,
            "budget-ratio <= 1.0 and benchmark exits 0",
        )
    return _format_retry_suggestion(
        "Benchmark threshold exceeded",
        code_target,
        "profile scenario and remove dominant hot-path operations that drive budget-ratio above threshold",
        verify_command,
        "budget-ratio <= 1.0 and benchmark exits 0",
    )


def runtime_failure_cause_from_log(log_file: Path, domain: str) -> tuple[str, str, str, str] | None:
    if not log_file.is_file():
        return None
    lines = _normalized_log_lines(log_file)
    if not lines:
        return None

    scenario = "(benchmark command)" if domain == "benchmark" else "(unit tests)"
    marker = _latest_scenario_marker(lines, domain)
    if marker is not None:
        scenario = marker

    detail = _data_race_detail(lines)
    if detail is not None:
        cause = "Data race detected"
    else:
        cause = ""
        for line in lines:
            panic_match = re.search(r"\bpanic:\s*(.+)$", line)
            if panic_match:
                detail = panic_match.group(1).strip()
                cause = "Runtime panic"
                break
            fatal_match = re.search(r"\bfatal error:\s*(.+)$", line)
            if fatal_match:
                detail = fatal_match.group(1).strip()
                cause = "Fatal runtime error"
                break
            signal_match = re.search(r"\bsignal:\s*([^\s].*)$", line)
            if signal_match:
                detail = signal_match.group(1).strip()
                cause = "Process terminated by signal"
                break

    if detail is None:
        for line in lines:
            if "[gen-impl] command timeout" in line:
                detail = "command timeout"
                cause = "Command timed out"
                break
        if detail is None:
            for line in lines:
                exit_match = re.search(r"\bexit status\s+([0-9]+)\b", line)
                if exit_match:
                    detail = f"exit status {exit_match.group(1)}"
                    cause = "Command exited non-zero"
                    break
        if detail is None:
            for line in lines:
                make_match = re.search(r"make:\s+\*\*\*\s+.+?\s+Error\s+([0-9]+)\b", line)
                if make_match:
                    detail = f"make error {make_match.group(1)}"
                    cause = "Command exited non-zero"
                    break

    if detail is None:
        return None

    evidence = f"{log_file} ({_compact_line(detail)})"
    folder = _infer_folder_from_path(log_file)
    code_target = _extract_code_target(lines, scenario)
    return (cause, scenario, evidence, _runtime_failure_suggestion(cause, detail, scenario, domain, code_target, folder))


def _data_race_detail(lines: list[str]) -> str | None:
    for idx, line in enumerate(lines):
        if "WARNING: DATA RACE" not in line and line != "DATA RACE":
            continue
        for follow in lines[idx + 1 :]:
            lower = follow.lower()
            if lower.startswith("read at ") or lower.startswith("write at "):
                return f"{line} -> {follow}"
            if lower.startswith("previous read at ") or lower.startswith("previous write at "):
                return f"{line} -> {follow}"
        return line

    return None


def failure_improvement_suggestions_from_causes(
    rows: list[tuple[str, str, str, str]],
) -> list[tuple[str, str, str]]:
    suggestions: list[tuple[str, str, str]] = []
    seen: set[tuple[str, str, str]] = set()
    for cause, scenario, _, suggestion in rows:
        key = (cause, scenario, suggestion)
        if key in seen:
            continue
        seen.add(key)
        suggestions.append((cause, scenario, suggestion))
    return suggestions


def _normalized_log_lines(log_file: Path) -> list[str]:
    normalized: list[str] = []
    for raw_line in log_file.read_text(encoding="utf-8", errors="replace").splitlines():
        text = raw_line.strip()
        if not text:
            continue
        if text.startswith("{"):
            try:
                payload = json.loads(text)
            except json.JSONDecodeError:
                normalized.append(text)
                continue
            output = payload.get("Output")
            if isinstance(output, str) and output.strip():
                normalized.append(output.strip())
                continue
            action = payload.get("Action")
            test_name = payload.get("Test")
            if action == "fail" and isinstance(test_name, str) and test_name:
                normalized.append(f"--- FAIL: {test_name}")
                continue
            continue
        normalized.append(text)
    return normalized


def _latest_scenario_marker(lines: list[str], domain: str) -> str | None:
    scenario: str | None = None
    for line in lines:
        if domain == "benchmark":
            benchmark_match = re.search(r"\b(Benchmark[^\s]+)", line)
            if benchmark_match:
                scenario = benchmark_match.group(1)
        test_match = re.search(r"^--- FAIL:\s+([^\s(]+)", line)
        if test_match:
            scenario = test_match.group(1)
    return scenario


def _compact_line(text: str, limit: int = 180) -> str:
    text = re.sub(r"\s+", " ", text).strip()
    if len(text) <= limit:
        return text
    return text[: limit - 3] + "..."


def _runtime_failure_suggestion(
    cause: str,
    detail: str,
    scenario: str,
    domain: str,
    code_target: str,
    folder: str | None,
) -> str:
    verify_command = _benchmark_verify_command(folder) if domain == "benchmark" else _operational_verify_command(folder, scenario)
    pass_signal = (
        "benchmark exits 0 without panic"
        if domain == "benchmark"
        else "go test -race exits 0 with no panic or race report"
    )
    lower = detail.lower()
    if cause == "Data race detected":
        return _format_retry_suggestion(
            "Data race detected",
            code_target,
            "isolate shared mutable state, remove concurrent read/write paths, and honor iterator mutation-unsafety contract before rerun",
            verify_command,
            pass_signal,
        )
    if "capacity exhausted" in lower:
        return _format_retry_suggestion(
            "Runtime panic (capacity exhausted)",
            code_target,
            "trigger growth before insert, rehash only live entries, and enforce load-factor guard before write",
            verify_command,
            pass_signal,
        )
    if "index out of range" in lower:
        return _format_retry_suggestion(
            "Runtime panic (index out of range)",
            code_target,
            "bound probe/index loops to backing array length and guard every slot access",
            verify_command,
            pass_signal,
        )
    if "nil pointer" in lower or "invalid memory address" in lower:
        return _format_retry_suggestion(
            "Runtime panic (nil pointer)",
            code_target,
            "initialize internal state before first mutation and guard nil paths before dereference",
            verify_command,
            pass_signal,
        )
    if "concurrent" in lower:
        return _format_retry_suggestion(
            "Runtime failure (concurrency)",
            code_target,
            "remove concurrent mutation/read path and enforce mutation-unsafety contract during iteration",
            verify_command,
            pass_signal,
        )
    if "timeout" in lower:
        return _format_retry_suggestion(
            "Command timed out",
            code_target,
            "remove unbounded loops and reduce pathological probe/scan complexity causing timeout",
            verify_command,
            "command exits 0 within timeout",
        )
    if cause == "Command exited non-zero":
        return _format_retry_suggestion(
            "Command exited non-zero",
            code_target,
            "patch function referenced by first failing stack/test line in log before rerun",
            verify_command,
            "command exits 0",
        )
    if cause == "Process terminated by signal":
        return _format_retry_suggestion(
            "Process terminated by signal",
            code_target,
            "stabilize memory/runtime path that triggers signal and remove invalid access pattern",
            verify_command,
            pass_signal,
        )
    if cause == "Fatal runtime error":
        return _format_retry_suggestion(
            "Fatal runtime error",
            code_target,
            "fix invariant violation shown in fatal message and align control flow with container contract",
            verify_command,
            pass_signal,
        )
    if domain == "benchmark":
        return _format_retry_suggestion(
            f"{cause} ({_compact_line(detail, 64)})",
            code_target,
            "patch failing benchmark path based on panic/fatal detail and remove unstable branch",
            verify_command,
            pass_signal,
        )
    return _format_retry_suggestion(
        f"{cause} ({_compact_line(detail, 64)})",
        code_target,
        "patch failing test path based on panic/fatal detail and align behavior with SPECS and tests",
        verify_command,
        pass_signal,
    )


def _format_retry_suggestion(
    root_issue: str,
    code_target: str,
    fix_action: str,
    verify_command: str,
    pass_signal: str,
) -> str:
    return f"{root_issue} -> {code_target} -> {fix_action} -> {verify_command} -> {pass_signal}."


def _benchmark_verify_command(folder: str | None) -> str:
    resolved = folder if folder else "<folder>"
    return f"make bench-folder FOLDER={resolved}"


def _operational_verify_command(folder: str | None, scenario: str) -> str:
    resolved = folder if folder else "<folder>"
    if scenario and not scenario.startswith("("):
        return f"go test -race -json ./{resolved}/... -run ^{re.escape(scenario)}$"
    return f"go test -race -json ./{resolved}/..."


def _infer_folder_from_path(path: Path) -> str | None:
    parts = list(path.parts)
    for idx, part in enumerate(parts):
        if re.fullmatch(r"attempt_\d+", part) and idx > 0:
            return parts[idx - 1]
    if path.name in {"tests.log", "benchmarks.log"}:
        parent = path.parent.name
        if not re.fullmatch(r"attempt_\d+", parent):
            return parent
    return None


def _extract_code_target(lines: list[str], scenario: str) -> str:
    file_line_re = re.compile(r"(/[^:\s]+\.go):(\d+)(?::\d+)?")
    func_name: str | None = None
    file_hint: str | None = None
    for idx, line in enumerate(lines):
        file_match = file_line_re.search(line)
        if not file_match:
            continue
        file_hint = f"{Path(file_match.group(1)).name}:{file_match.group(2)}"
        prev = lines[idx - 1] if idx > 0 else ""
        func_match = re.search(r"\.([A-Za-z_][A-Za-z0-9_]*)\(", prev)
        if func_match:
            func_name = func_match.group(1)
        break
    if file_hint and func_name:
        return f"{file_hint} ({func_name})"
    if file_hint:
        return file_hint
    if scenario and not scenario.startswith("("):
        return scenario
    return "implementation path"


def _failure_causes_rows(tsv_file: Path) -> list[str]:
    if not tsv_file.is_file() or tsv_file.stat().st_size == 0:
        return ["| 1 | (none) | (none) | (none) | No failures recorded. |"]
    rows: list[str] = []
    for idx, line in enumerate(tsv_file.read_text(encoding="utf-8", errors="replace").splitlines(), start=1):
        pieces = line.split("\t")
        if len(pieces) != 4:
            continue
        cause, scenario, evidence, suggestion = pieces
        rows.append(
            f"| {idx} | {_md_cell(cause)} | {_md_cell(scenario)} | {_md_cell(evidence)} | {_md_cell(suggestion)} |"
        )
    if not rows:
        return ["| 1 | (none) | (none) | (none) | No failures recorded. |"]
    return rows


def _failure_improvement_rows(tsv_file: Path) -> list[str]:
    if not tsv_file.is_file() or tsv_file.stat().st_size == 0:
        return ["| 1 | (none) | (none) | No failures recorded. |"]
    rows: list[str] = []
    seen: set[tuple[str, str, str]] = set()
    for line in tsv_file.read_text(encoding="utf-8", errors="replace").splitlines():
        pieces = line.split("\t")
        if len(pieces) != 4:
            continue
        cause, scenario, _, suggestion = pieces
        key = (cause, scenario, suggestion)
        if key in seen:
            continue
        seen.add(key)
        rows.append(
            f"| {len(rows) + 1} | {_md_cell(cause)} | {_md_cell(scenario)} | {_md_cell(suggestion)} |"
        )
    if not rows:
        return ["| 1 | (none) | (none) | No failures recorded. |"]
    return rows
