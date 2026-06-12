from __future__ import annotations

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
    for line in lines:
        if line == "## Failure Causes":
            capture = True
        if capture and line.startswith("## ") and line != "## Failure Causes":
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
    if expected_status == "FAILURE" and "## Failure Causes" not in text:
        return False, f"missing failure cause table in {report_file}"
    return True, ""


def render_report_content(status: str, folder: str, attempts_used: int, change_table_file: Path, failure_causes_file: Path) -> str:
    parts = [
        "# IMPLEMENTATION_REPORT",
        "",
        f"Operation Status: {status}",
        f"Folder: {folder}",
        f"Attempts Used: {attempts_used}",
        "",
        "## Files Changed",
    ]
    parts.extend(change_table_file.read_text(encoding="utf-8", errors="replace").strip().splitlines())
    parts.append("")

    if status == "FAILURE":
        parts.extend(["## Failure Causes", "", "| Cause | Evidence | Suggestion |", "|---|---|---|"])
        if failure_causes_file.is_file() and failure_causes_file.stat().st_size > 0:
            for line in failure_causes_file.read_text(encoding="utf-8", errors="replace").splitlines():
                pieces = line.split("\t")
                if len(pieces) != 3:
                    continue
                cause, evidence, suggestion = pieces
                parts.append(f"| {cause} | {evidence} | {suggestion} |")
        else:
            parts.append("| Unknown failure | `(missing)` | Review latest attempt logs and rerun after fixing root cause. |")
        parts.append("")

    return "\n".join(parts).rstrip() + "\n"
