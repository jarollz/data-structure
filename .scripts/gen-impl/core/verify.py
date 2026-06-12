from __future__ import annotations

import subprocess
from pathlib import Path


def _run_to_log(command: list[str], cwd: Path, timeout_seconds: int, output_log: Path) -> bool:
    output_log.parent.mkdir(parents=True, exist_ok=True)
    try:
        with output_log.open("w", encoding="utf-8") as fh:
            completed = subprocess.run(
                command,
                cwd=cwd,
                stdout=fh,
                stderr=subprocess.STDOUT,
                timeout=timeout_seconds,
                check=False,
                text=True,
            )
        return completed.returncode == 0
    except subprocess.TimeoutExpired:
        with output_log.open("a", encoding="utf-8") as fh:
            fh.write("\n[gen-impl] command timeout\n")
        return False


def run_doc_comment_audit(repo_root: Path, script_root: Path, folder: str, output_log: Path, timeout_seconds: int) -> bool:
    audit_tool = script_root / "helpers" / "audit_doc_comments.go"
    return _run_to_log(
        ["go", "run", str(audit_tool), "--folder", folder],
        cwd=repo_root,
        timeout_seconds=timeout_seconds,
        output_log=output_log,
    )


def run_make_with_timeout(repo_root: Path, timeout_seconds: int, output_log: Path, make_target: str, folder: str) -> bool:
    return _run_to_log(
        ["make", "-C", str(repo_root), make_target, f"FOLDER={folder}"],
        cwd=repo_root,
        timeout_seconds=timeout_seconds,
        output_log=output_log,
    )
