from __future__ import annotations

import os
import signal
import subprocess
import threading
from typing import TextIO


def _wait(proc: subprocess.Popen[str], timeout_seconds: float) -> bool:
    try:
        proc.wait(timeout=timeout_seconds)
        return True
    except subprocess.TimeoutExpired:
        return False


def _terminate_single_process(proc: subprocess.Popen[str], grace_seconds: float) -> None:
    try:
        proc.terminate()
    except ProcessLookupError:
        return

    if _wait(proc, grace_seconds):
        return

    try:
        proc.kill()
    except ProcessLookupError:
        return
    _wait(proc, 1.0)


def _terminate_process_group(proc: subprocess.Popen[str], grace_seconds: float) -> bool:
    try:
        pgid = os.getpgid(proc.pid)
    except ProcessLookupError:
        return True

    try:
        os.killpg(pgid, signal.SIGTERM)
    except ProcessLookupError:
        return True
    except PermissionError:
        return False

    if _wait(proc, grace_seconds):
        return True

    try:
        os.killpg(pgid, signal.SIGKILL)
    except ProcessLookupError:
        return True
    except PermissionError:
        return False

    _wait(proc, 1.0)
    return True


def terminate_timed_out_process(
    proc: subprocess.Popen[str],
    reader: threading.Thread,
    output: TextIO,
    timeout_message: str,
    grace_seconds: float = 5.0,
) -> None:
    output.write(f"\n{timeout_message}\n")
    output.flush()

    if not _terminate_process_group(proc, grace_seconds):
        _terminate_single_process(proc, grace_seconds)

    reader.join(timeout=1.0)
