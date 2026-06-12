from __future__ import annotations

import os
from dataclasses import dataclass
from pathlib import Path


SUPPORTED_FOLDERS: tuple[str, ...] = (
    "list-array",
    "list-linked-singly",
    "list-linked-doubly",
    "list-skip",
    "queue",
    "stack",
    "heap",
    "tree-general",
    "tree-avl",
    "tree-red-black",
    "map-hash",
    "map-trie",
    "map-tree-avl",
    "map-tree-red-black",
)


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


@dataclass(frozen=True)
class Config:
    repo_root: Path
    script_root: Path
    max_attempts: int
    report_attempts: int
    stop_on_failure: bool
    force: bool
    probe_timeout_seconds: int
    spawner_timeout_seconds: int
    doc_audit_timeout_seconds: int
    test_timeout_seconds: int
    bench_timeout_seconds: int

    @classmethod
    def from_env(cls) -> "Config":
        script_root = Path(__file__).resolve().parents[1]
        repo_root = script_root.parents[1]
        return cls(
            repo_root=repo_root,
            script_root=script_root,
            max_attempts=_env_int("MAX_ATTEMPTS", 5),
            report_attempts=_env_int("REPORT_ATTEMPTS", 5),
            stop_on_failure=_env_bool("STOP_ON_FAILURE", False),
            force=_env_bool("FORCE", False),
            probe_timeout_seconds=_env_int("PROBE_TIMEOUT_SECONDS", 120),
            spawner_timeout_seconds=_env_int("SPAWNER_TIMEOUT_SECONDS", 1800),
            doc_audit_timeout_seconds=_env_int("DOC_AUDIT_TIMEOUT_SECONDS", 300),
            test_timeout_seconds=_env_int("TEST_TIMEOUT_SECONDS", 900),
            bench_timeout_seconds=_env_int("BENCH_TIMEOUT_SECONDS", 1800),
        )
