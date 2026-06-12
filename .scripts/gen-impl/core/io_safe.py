from __future__ import annotations

import tempfile
import sys
from datetime import datetime, timezone
from pathlib import Path


def log(message: str) -> None:
    print(f"[gen-impl] {message}")


def warn(message: str) -> None:
    print(f"[gen-impl] Warning: {message}", file=sys.stderr)


def die(message: str) -> None:
    print(f"[gen-impl] Error: {message}", file=sys.stderr)


def timestamp_utc() -> str:
    return datetime.now(tz=timezone.utc).strftime("%Y%m%d_%H%M%S")


def ensure_dir(path: Path) -> None:
    path.mkdir(parents=True, exist_ok=True)


def ensure_dir_writable(path: Path, label: str) -> None:
    ensure_dir(path)
    try:
        with tempfile.NamedTemporaryFile(prefix=".writable_", dir=path, delete=True):
            pass
    except PermissionError as exc:
        raise RuntimeError(f"{label} not writable: {path}: {exc}") from exc


def can_write_path(path: Path) -> tuple[bool, str]:
    parent = path.parent
    try:
        parent.mkdir(parents=True, exist_ok=True)
    except PermissionError as exc:
        return False, f"cannot create parent directory {parent}: {exc}"

    try:
        if path.exists():
            with path.open("a", encoding="utf-8"):
                pass
        else:
            with tempfile.NamedTemporaryFile(prefix=".probe_", dir=parent, delete=True):
                pass
    except PermissionError as exc:
        return False, str(exc)
    return True, ""


def write_text(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content, encoding="utf-8")
