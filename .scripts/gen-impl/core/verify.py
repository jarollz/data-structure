from __future__ import annotations

import subprocess
import threading
import time
from collections.abc import Callable
from pathlib import Path

from .process_control import terminate_timed_out_process


def _run_to_log(
    command: list[str],
    cwd: Path,
    timeout_seconds: int,
    output_log: Path,
    progress_callback: Callable[[float, str | None], None] | None = None,
) -> bool:
    output_log.parent.mkdir(parents=True, exist_ok=True)
    with output_log.open("w", encoding="utf-8") as fh:
        proc = subprocess.Popen(
            command,
            cwd=cwd,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            bufsize=1,
            start_new_session=True,
        )

        tail: dict[str, str | None] = {"value": None}

        def _reader() -> None:
            if proc.stdout is None:
                return
            try:
                for chunk in proc.stdout:
                    fh.write(chunk)
                    fh.flush()
                    for piece in chunk.replace("\r", "\n").split("\n"):
                        if piece.strip():
                            tail["value"] = piece.strip()
            finally:
                if proc.stdout is not None:
                    proc.stdout.close()

        reader = threading.Thread(target=_reader, daemon=True)
        reader.start()

        start = time.monotonic()
        while True:
            elapsed = time.monotonic() - start
            if progress_callback is not None:
                progress_callback(elapsed, tail["value"])

            code = proc.poll()
            if code is not None:
                reader.join()
                return code == 0

            if elapsed >= timeout_seconds:
                terminate_timed_out_process(proc, reader, fh, "[gen-impl] command timeout")
                return False

            time.sleep(0.2)


def run_doc_comment_audit(
    repo_root: Path,
    script_root: Path,
    folder: str,
    output_log: Path,
    timeout_seconds: int,
    progress_callback: Callable[[float, str | None], None] | None = None,
) -> bool:
    audit_tool = script_root / "helpers" / "audit_doc_comments.go"
    return _run_to_log(
        ["go", "run", str(audit_tool), "--folder", folder],
        cwd=repo_root,
        timeout_seconds=timeout_seconds,
        output_log=output_log,
        progress_callback=progress_callback,
    )


def run_make_with_timeout(
    repo_root: Path,
    timeout_seconds: int,
    output_log: Path,
    make_target: str,
    folder: str,
    progress_callback: Callable[[float, str | None], None] | None = None,
) -> bool:
    return _run_to_log(
        ["make", "-C", str(repo_root), make_target, f"FOLDER={folder}"],
        cwd=repo_root,
        timeout_seconds=timeout_seconds,
        output_log=output_log,
        progress_callback=progress_callback,
    )


def run_go_test_json_with_timeout(
    repo_root: Path,
    folder: str,
    timeout_seconds: int,
    output_log: Path,
    progress_callback: Callable[[float, str | None], None] | None = None,
) -> bool:
    return _run_to_log(
        ["go", "test", "-json", f"./{folder}/..."],
        cwd=repo_root,
        timeout_seconds=timeout_seconds,
        output_log=output_log,
        progress_callback=progress_callback,
    )
