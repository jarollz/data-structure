from __future__ import annotations

import json
import os
from pathlib import PurePosixPath
import subprocess
from dataclasses import dataclass
from pathlib import Path


def _normalize_workspace_folder(repo_root: Path, disk_path: str) -> str:
    if not disk_path:
        raise RuntimeError("root go.work contains an empty use entry")

    raw_path = PurePosixPath(disk_path)
    if raw_path.is_absolute():
        raise RuntimeError(f"root go.work entry '{disk_path}' must be relative")
    if ".." in raw_path.parts:
        raise RuntimeError(f"root go.work entry '{disk_path}' must not contain '..'")

    resolved_path = (repo_root / Path(disk_path)).resolve()
    resolved_root = repo_root.resolve()
    try:
        relative_path = resolved_path.relative_to(resolved_root)
    except ValueError as exc:
        raise RuntimeError(f"root go.work entry '{disk_path}' points outside repo root") from exc

    if len(relative_path.parts) != 1:
        raise RuntimeError(f"root go.work entry '{disk_path}' must point to a top-level folder")
    if not resolved_path.is_dir():
        raise RuntimeError(f"root go.work entry '{disk_path}' points to missing directory: {relative_path.as_posix()}")
    return relative_path.as_posix()


def discover_supported_folders(repo_root: Path) -> tuple[str, ...]:
    go_work = repo_root / "go.work"
    if not go_work.is_file():
        raise RuntimeError(f"missing root go.work: {go_work}")

    try:
        proc = subprocess.run(
            ["go", "work", "edit", "-json"],
            cwd=repo_root,
            check=False,
            capture_output=True,
            text=True,
        )
    except OSError as exc:
        raise RuntimeError(f"failed to run 'go work edit -json': {exc}") from exc

    if proc.returncode != 0:
        reason = proc.stderr.strip() or proc.stdout.strip() or f"exit status {proc.returncode}"
        raise RuntimeError(f"'go work edit -json' failed: {reason}")

    try:
        payload = json.loads(proc.stdout)
    except json.JSONDecodeError as exc:
        raise RuntimeError(f"invalid JSON from 'go work edit -json': {exc}") from exc

    use_entries = payload.get("Use")
    if not isinstance(use_entries, list) or not use_entries:
        raise RuntimeError("root go.work contains no use entries")

    folders: list[str] = []
    seen: set[str] = set()
    for entry in use_entries:
        if not isinstance(entry, dict):
            raise RuntimeError("root go.work contains an invalid use entry")
        disk_path = entry.get("DiskPath")
        if not isinstance(disk_path, str):
            raise RuntimeError("root go.work contains a use entry without DiskPath")
        folder = _normalize_workspace_folder(repo_root, disk_path)
        if folder in seen:
            raise RuntimeError(f"root go.work lists duplicate folder '{folder}'")
        seen.add(folder)
        folders.append(folder)

    return tuple(folders)


def _env_int(name: str, default: int) -> int:
    raw = os.environ.get(name)
    if raw is None or raw == "":
        return default
    try:
        value = int(raw)
    except ValueError:
        return default
    return value if value > 0 else default


def _env_bool(name: str, default: bool = False) -> bool:
    raw = os.environ.get(name)
    if raw is None:
        return default
    return raw == "1"


def _env_str(name: str, default: str = "") -> str:
    raw = os.environ.get(name)
    if raw is None:
        return default
    return raw.strip()


def _env_force_mode(name: str = "FORCE", default: int = 0) -> int:
    raw = os.environ.get(name)
    if raw is None or raw.strip() == "":
        return default
    value = raw.strip()
    if value not in {"0", "1", "2"}:
        raise RuntimeError(f"invalid {name} value '{raw}'; expected one of: 0, 1, 2")
    return int(value)


@dataclass(frozen=True)
class Config:
    repo_root: Path
    script_root: Path
    supported_folders: tuple[str, ...]
    max_attempts: int
    report_attempts: int
    stop_on_failure: bool
    force_mode: int
    probe_timeout_seconds: int
    spawner_timeout_seconds: int
    spawner_idle_timeout_seconds: int
    doc_audit_timeout_seconds: int
    test_timeout_seconds: int
    bench_timeout_seconds: int
    ai_spawner_command: str

    @classmethod
    def from_env(cls) -> "Config":
        script_root = Path(__file__).resolve().parents[1]
        repo_root = script_root.parents[1]
        return cls(
            repo_root=repo_root,
            script_root=script_root,
            supported_folders=discover_supported_folders(repo_root),
            max_attempts=_env_int("MAX_ATTEMPTS", 5),
            report_attempts=_env_int("REPORT_ATTEMPTS", 5),
            stop_on_failure=_env_bool("STOP_ON_FAILURE", False),
            force_mode=_env_force_mode("FORCE", 0),
            probe_timeout_seconds=_env_int("PROBE_TIMEOUT_SECONDS", 120),
            spawner_timeout_seconds=_env_int("SPAWNER_TIMEOUT_SECONDS", 1800),
            spawner_idle_timeout_seconds=_env_int("SPAWNER_IDLE_TIMEOUT_SECONDS", 180),
            doc_audit_timeout_seconds=_env_int("DOC_AUDIT_TIMEOUT_SECONDS", 300),
            test_timeout_seconds=_env_int("TEST_TIMEOUT_SECONDS", 900),
            bench_timeout_seconds=_env_int("BENCH_TIMEOUT_SECONDS", 1800),
            ai_spawner_command=_env_str("AI_SPAWNER_COMMAND", ""),
        )
