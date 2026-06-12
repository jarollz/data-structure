from __future__ import annotations

import hashlib
import shutil
from pathlib import Path


PROTECTED_BASENAMES = {
    "api_contract.go",
    "go.mod",
    "go.sum",
    "SPECS.md",
    "helpers_test.go",
    "bench_policy_test.go",
}


def _is_excluded(rel_path: str) -> bool:
    return rel_path.startswith(".git/") or rel_path.startswith("tmp/gen-impl/runs/")


def create_phase_snapshot(snapshot_dir: Path, repo_root: Path) -> None:
    if snapshot_dir.exists():
        shutil.rmtree(snapshot_dir)
    snapshot_repo = snapshot_dir / "repo"
    snapshot_repo.mkdir(parents=True, exist_ok=True)

    for source in repo_root.rglob("*"):
        rel = source.relative_to(repo_root).as_posix()
        if rel == ".":
            continue
        if _is_excluded(rel):
            continue
        target = snapshot_repo / rel
        if source.is_dir():
            target.mkdir(parents=True, exist_ok=True)
        elif source.is_file():
            target.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(source, target)


def build_manifest(root_dir: Path) -> dict[str, str]:
    result: dict[str, str] = {}
    for file_path in root_dir.rglob("*"):
        if not file_path.is_file():
            continue
        rel = file_path.relative_to(root_dir).as_posix()
        if _is_excluded(rel):
            continue
        digest = hashlib.sha256(file_path.read_bytes()).hexdigest()
        result[rel] = digest
    return result


def diff_manifests(before: dict[str, str], after: dict[str, str]) -> list[tuple[str, str]]:
    paths = sorted(set(before) | set(after))
    out: list[tuple[str, str]] = []
    for path in paths:
        if path not in after:
            out.append(("D", path))
        elif path not in before:
            out.append(("A", path))
        elif before[path] != after[path]:
            out.append(("M", path))
    return out


def _path_is_protected(path: str, folder: str) -> bool:
    if not path.startswith(f"{folder}/"):
        return False
    base = path.rsplit("/", 1)[-1]
    if base in PROTECTED_BASENAMES:
        return True
    if base.endswith("_test.go") or base.endswith("_bench_test.go"):
        return True
    if base.endswith(".md"):
        return True
    return False


def _restore_path_from_snapshot(snapshot_dir: Path, repo_root: Path, relative_path: str) -> None:
    snapshot_file = snapshot_dir / "repo" / relative_path
    repo_file = repo_root / relative_path
    if snapshot_file.is_file():
        repo_file.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(snapshot_file, repo_file)
    else:
        if repo_file.exists():
            repo_file.unlink()


def _write_tsv(path: Path, rows: list[tuple[str, str]]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", encoding="utf-8") as fh:
        for change_type, rel in rows:
            fh.write(f"{change_type}\t{rel}\n")


def enforce_phase_scope(
    folder: str,
    snapshot_dir: Path,
    allow_report_write: bool,
    result_dir: Path,
    repo_root: Path,
    allowed_report_relative_path: str | None = None,
) -> bool:
    before = build_manifest(snapshot_dir / "repo")
    after = build_manifest(repo_root)
    changes = diff_manifests(before, after)

    kept_changes: list[tuple[str, str]] = []
    protected_violations: list[tuple[str, str]] = []
    scope_violations: list[tuple[str, str]] = []
    phase_ok = True

    for change_type, path in changes:
        if path.startswith("tmp/gen-impl/runs/"):
            continue

        if path.startswith(f"{folder}/"):
            if _path_is_protected(path, folder):
                protected_violations.append((change_type, path))
                _restore_path_from_snapshot(snapshot_dir, repo_root, path)
                phase_ok = False
            else:
                kept_changes.append((change_type, path))
            continue

        if allow_report_write and allowed_report_relative_path and path == allowed_report_relative_path:
            kept_changes.append((change_type, path))
            continue

        scope_violations.append((change_type, path))
        _restore_path_from_snapshot(snapshot_dir, repo_root, path)
        phase_ok = False

    _write_tsv(result_dir / "kept_changes.tsv", kept_changes)
    _write_tsv(result_dir / "protected_violations.tsv", protected_violations)
    _write_tsv(result_dir / "scope_violations.tsv", scope_violations)
    _write_tsv(result_dir / "changes.tsv", changes)
    return phase_ok


def write_change_table_markdown(folder: str, start_snapshot: Path, repo_root: Path, output_file: Path) -> None:
    before = build_manifest(start_snapshot / "repo")
    after = build_manifest(repo_root)
    changes = diff_manifests(before, after)

    lines = ["| File | Change |", "|---|---|"]
    count = 0
    for change_type, path in changes:
        if not path.startswith(f"{folder}/"):
            continue
        if change_type == "A":
            label = "added"
        elif change_type == "D":
            label = "deleted"
        else:
            label = "modified"
        lines.append(f"| `{path}` | {label} |")
        count += 1

    if count == 0:
        lines.append("| `(none)` | unchanged |")

    output_file.parent.mkdir(parents=True, exist_ok=True)
    output_file.write_text("\n".join(lines) + "\n", encoding="utf-8")
