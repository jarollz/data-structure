#!/usr/bin/env python3
"""Assess data-structure folders against repo contracts.

Persistent helper for `.agents/skills/assess-implementation` skill.
Strict-hard scoring: command evidence failures apply hard caps.
"""

from __future__ import annotations

import argparse
import datetime as dt
import re
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path


CHECKLIST_RE = re.compile(r"^\s*-\s*\[[ xX]\]\s+(.*)$")
BT_RE = re.compile(r"`([^`]+)`")
FUNC_RE = re.compile(
    r"func\s+(?:\([^)]*\)\s*)?([A-Z][A-Za-z0-9_]*)\s*(?:\[[^\]]+\])?\s*\(",
    re.MULTILINE,
)
TEST_RE = re.compile(r"func\s+Test[A-Za-z0-9_]*\s*\(", re.MULTILINE)
BENCH_RE = re.compile(r"func\s+Benchmark[A-Za-z0-9_]*\s*\(", re.MULTILINE)
MAP_RE = re.compile(r"\bmap\s*\[")
SLICE_RE = re.compile(r"\[\]\s*[*A-Za-z_\[]")
SIZE_RE = re.compile(r"\b(1e3|1e4|1e5|1000|10000|100000)\b")


@dataclass
class Category:
    name: str
    score: int
    max_score: int
    evidence: str


@dataclass
class CommandStatus:
    test_required: bool
    test_ran: bool
    test_ok: bool
    bench_required: bool
    bench_ran: bool
    bench_ok: bool
    fallback_used: bool
    fallback_reasons: list[str]


@dataclass
class FolderResult:
    folder: str
    score: int
    grade: str
    specs_passed: int
    specs_total: int
    specs_pct: int
    tests_score: int
    bench_score: int
    notes: str
    hard_fail: str
    categories: list[Category]
    findings: list[str]
    command_evidence: list[str]
    weak_evidence: list[str]
    command_status: CommandStatus
    strict_cap_reasons: list[str]


@dataclass
class RunMetadata:
    mode: str
    targets: list[str]
    command_policy: str
    command_execution: str


def run(cmd: list[str], cwd: Path) -> tuple[bool, str]:
    proc = subprocess.run(cmd, cwd=str(cwd), text=True, capture_output=True)
    text = (proc.stdout or "") + (proc.stderr or "")
    text = text.strip().replace("\n", " ")
    if len(text) > 260:
        text = text[:257] + "..."
    return proc.returncode == 0, text


def normalize_targets(targets: list[str]) -> list[str]:
    out: list[str] = []
    for t in targets:
        t = t.strip().strip("`\"")
        if t.startswith("./"):
            t = t[2:]
        if t and t not in out:
            out.append(t)
    return out


def discover_folders(repo: Path) -> list[str]:
    names: list[str] = []
    for child in sorted(repo.iterdir(), key=lambda p: p.name):
        if child.is_dir() and (child / "SPECS.md").exists():
            names.append(child.name)
    return names


def required_api_names(specs_text: str) -> list[str]:
    names: list[str] = []
    in_required = False
    for line in specs_text.splitlines():
        if line.startswith("## "):
            in_required = line.strip().lower() == "## required api"
            continue
        if not in_required:
            continue
        m = CHECKLIST_RE.match(line)
        if not m:
            continue
        body = m.group(1)
        for token in BT_RE.findall(body):
            if "(" not in token:
                continue
            head = token.split("(", 1)[0].strip()
            k = re.match(r"^[A-Za-z_][A-Za-z0-9_]*$", head)
            if k:
                names.append(k.group(0))
                break
    dedup: list[str] = []
    for n in names:
        if n not in dedup:
            dedup.append(n)
    return dedup


def expected_iterator(folder: str) -> str:
    if folder.startswith("map-"):
        return "All"
    if folder in {"tree-avl", "tree-red-black"}:
        return "InOrder"
    if folder == "tree-general":
        return "PreOrder"
    return "Values"


def extract_public_funcs(texts: list[str]) -> set[str]:
    out: set[str] = set()
    for t in texts:
        for m in FUNC_RE.finditer(t):
            out.add(m.group(1))
    return out


def is_placeholder_test(content: str) -> bool:
    names = re.findall(r"func\s+(Test[A-Za-z0-9_]*)\s*\(", content)
    if not names:
        return True
    if all(n in {"TestNothing", "TestPlaceholder", "TestDummy"} for n in names):
        return True
    return False


def parse_checklist(specs_text: str) -> list[str]:
    items: list[str] = []
    for line in specs_text.splitlines():
        m = CHECKLIST_RE.match(line)
        if m:
            items.append(m.group(1).strip())
    return items


def checklist_passes(
    items: list[str],
    api_found: set[str],
    readme_ok: bool,
    specs_ok: bool,
    has_impl: bool,
    has_tests: bool,
    has_bench: bool,
    no_map: bool,
    no_slice: bool,
) -> int:
    passed = 0
    for item in items:
        low = item.lower()
        done = False
        if "readme.md" in low:
            done = readme_ok
        elif "specs.md" in low:
            done = specs_ok
        elif "do not use `map`" in low:
            done = no_map
        elif "do not use `slice`" in low:
            done = no_slice
        elif "benchmark" in low:
            done = has_bench
        elif "test" in low:
            done = has_tests
        elif "implementation" in low or "internal representation" in low:
            done = has_impl
        else:
            for token in BT_RE.findall(item):
                if "(" in token:
                    name = token.split("(", 1)[0].strip()
                    if name in api_found:
                        done = True
                        break
        if done:
            passed += 1
    return passed


def grade(score: int, hard_fail: bool) -> str:
    if hard_fail or score < 60:
        return "F"
    if score >= 90:
        return "A"
    if score >= 80:
        return "B"
    if score >= 70:
        return "C"
    return "D"


def method_coverage_ratio(required_methods: list[str], test_blob: str) -> float:
    if not required_methods:
        return 1.0
    hit = 0
    for method in required_methods:
        pattern = re.compile(rf"\b{re.escape(method.lower())}\s*\(")
        if pattern.search(test_blob):
            hit += 1
    return hit / len(required_methods)


def should_fallback_to_go(output: str) -> bool:
    low = output.lower()
    indicators = [
        "no rule to make target",
        "command not found",
        "not found",
        "invalid folder",
        "missing",
    ]
    return any(k in low for k in indicators)


def run_test_command(repo: Path, folder: str) -> tuple[bool, list[str], bool, list[str]]:
    evidence: list[str] = []
    fallback_reasons: list[str] = []
    ok, out = run(["make", "test-folder", f"FOLDER={folder}"], repo)
    evidence.append(f"`make test-folder FOLDER={folder}` -> {'ok' if ok else 'failed'} - {out or 'no output'}")
    if ok:
        return True, evidence, False, fallback_reasons
    if not should_fallback_to_go(out):
        return False, evidence, False, fallback_reasons
    fallback_reasons.append("make test-folder unavailable/broken; fallback to go test")
    fb_ok, fb_out = run(["go", "test", f"./{folder}"], repo)
    evidence.append(f"`go test ./{folder}` -> {'ok' if fb_ok else 'failed'} - {fb_out or 'no output'}")
    return fb_ok, evidence, True, fallback_reasons


def run_bench_command(repo: Path, folder: str) -> tuple[bool, list[str], bool, list[str]]:
    evidence: list[str] = []
    fallback_reasons: list[str] = []
    ok, out = run(["make", "bench-folder", f"FOLDER={folder}"], repo)
    evidence.append(
        f"`make bench-folder FOLDER={folder}` -> {'ok' if ok else 'failed'} - {out or 'no output'}"
    )
    if ok:
        return True, evidence, False, fallback_reasons
    if not should_fallback_to_go(out):
        return False, evidence, False, fallback_reasons
    fallback_reasons.append("make bench-folder unavailable/broken; fallback to go test -bench")
    fb_ok, fb_out = run(["go", "test", f"./{folder}", "-run", "^$", "-bench", ".", "-benchmem"], repo)
    evidence.append(
        f"`go test ./{folder} -run ^$ -bench . -benchmem` -> {'ok' if fb_ok else 'failed'} - {fb_out or 'no output'}"
    )
    return fb_ok, evidence, True, fallback_reasons


def assess_folder(repo: Path, folder: str, run_commands: bool) -> FolderResult:
    base = repo / folder
    specs_path = base / "SPECS.md"
    readme_ok = (base / "README.md").exists()
    specs_ok = specs_path.exists()
    specs_text = specs_path.read_text(encoding="utf-8") if specs_ok else ""
    checklist = parse_checklist(specs_text)

    go_files = sorted(base.glob("*.go"))
    impl_files = [p for p in go_files if not p.name.endswith("_test.go")]
    test_files = [p for p in go_files if p.name.endswith("_test.go")]
    impl_texts = [p.read_text(encoding="utf-8") for p in impl_files]
    test_texts = [p.read_text(encoding="utf-8") for p in test_files]
    test_blob = "\n".join(test_texts)
    test_blob_low = test_blob.lower()

    has_impl = len(impl_files) > 0
    api_required = required_api_names(specs_text)
    api_found = extract_public_funcs(impl_texts)
    missing_api = [n for n in api_required if n not in api_found]

    has_tests = False
    if test_texts:
        has_tests = bool(TEST_RE.search(test_blob)) and not is_placeholder_test(test_blob)
    has_bench = bool(BENCH_RE.search(test_blob)) if test_texts else False
    has_bench_sizes = bool(SIZE_RE.search(test_blob)) if test_texts else False

    no_map = True
    no_slice = True
    for content in impl_texts:
        if MAP_RE.search(content):
            no_map = False
        if SLICE_RE.search(content):
            no_slice = False

    iterator_name = expected_iterator(folder)
    iterator_impl_ok = iterator_name in api_found
    iterator_test_ok = has_tests and (
        "iterator" in test_blob_low
        or "values(" in test_blob_low
        or "all(" in test_blob_low
        or "inorder(" in test_blob_low
        or "preorder(" in test_blob_low
    )
    invariant_test_ok = has_tests and (
        "invariant" in test_blob_low
        or "len(" in test_blob_low
        or "cycle" in test_blob_low
        or "balance" in test_blob_low
    )
    random_oracle_ok = has_tests and ("random" in test_blob_low or "oracle" in test_blob_low)
    api_cov = method_coverage_ratio(api_required, test_blob_low)

    command_evidence: list[str] = []
    fallback_reasons: list[str] = []
    test_required = has_impl
    test_ran = False
    test_ok = False
    bench_required = has_impl and has_bench
    bench_ran = False
    bench_ok = False
    fallback_used = False

    if run_commands and has_impl:
        test_ran = True
        test_ok, ev, used_fb, fb = run_test_command(repo, folder)
        command_evidence.extend(ev)
        fallback_used = fallback_used or used_fb
        fallback_reasons.extend(fb)

        if has_bench:
            bench_ran = True
            bench_ok, ev, used_fb, fb = run_bench_command(repo, folder)
            command_evidence.extend(ev)
            fallback_used = fallback_used or used_fb
            fallback_reasons.extend(fb)
        else:
            command_evidence.append(f"`make bench-folder FOLDER={folder}` -> skipped (no benchmarks)")

    command_status = CommandStatus(
        test_required=test_required,
        test_ran=test_ran,
        test_ok=test_ok,
        bench_required=bench_required,
        bench_ran=bench_ran,
        bench_ok=bench_ok,
        fallback_used=fallback_used,
        fallback_reasons=fallback_reasons,
    )

    api_score = 0
    if api_required:
        api_score = int(round(25 * (len(api_required) - len(missing_api)) / len(api_required)))

    storage_score = 0
    if has_impl and no_map and no_slice:
        storage_score = 25
    elif has_impl and (no_map or no_slice):
        storage_score = 8

    behavior_score = 0
    if has_impl:
        behavior_score += 4
    if iterator_impl_ok:
        behavior_score += 4
    if iterator_impl_ok and iterator_test_ok:
        behavior_score += 6
    if invariant_test_ok:
        behavior_score += 6

    tests_score = 0
    if has_tests:
        tests_score += 6
        tests_score += int(round(6 * api_cov))
        if iterator_test_ok:
            tests_score += 4
        if random_oracle_ok:
            tests_score += 4

    bench_delivery_score = 0
    if readme_ok:
        bench_delivery_score += 2
    if specs_ok:
        bench_delivery_score += 1
    if has_tests:
        bench_delivery_score += 2
    if has_bench:
        bench_delivery_score += 3
    if has_bench and has_bench_sizes:
        bench_delivery_score += 2

    strict_cap_reasons: list[str] = []
    if test_required and (not test_ran or not test_ok):
        if tests_score > 4:
            strict_cap_reasons.append("test command failed or missing; cap Test evidence to 4/20")
        if behavior_score > 6:
            strict_cap_reasons.append("test command failed or missing; cap Invariants and behavior to 6/20")
        tests_score = min(tests_score, 4)
        behavior_score = min(behavior_score, 6)

    if bench_required and (not bench_ran or not bench_ok):
        if bench_delivery_score > 5:
            strict_cap_reasons.append("benchmark command failed or missing; cap Benchmark and delivery to 5/10")
        bench_delivery_score = min(bench_delivery_score, 5)

    specs_total = len(checklist)
    specs_passed = checklist_passes(
        checklist,
        api_found,
        readme_ok,
        specs_ok,
        has_impl,
        has_tests,
        has_bench,
        no_map,
        no_slice,
    )
    specs_pct = int(round((specs_passed / specs_total) * 100)) if specs_total else 0

    findings: list[str] = []
    for method in missing_api:
        findings.append(f"`{folder}/SPECS.md` required API `{method}` not found in implementation.")
    if not has_tests:
        findings.append(f"`{folder}` no non-placeholder tests proved.")
    if not has_bench:
        findings.append(f"`{folder}` no `Benchmark...` functions found.")
    if not no_map:
        findings.append(f"`{folder}` implementation uses forbidden built-in `map`.")
    if not no_slice:
        findings.append(f"`{folder}` implementation uses forbidden slice syntax.")
    if not has_impl:
        findings.append(f"`{folder}` has no non-test implementation `.go` files.")
    if test_required and (not test_ran or not test_ok):
        findings.append(f"`{folder}` required test command evidence failed or missing.")
    if bench_required and (not bench_ran or not bench_ok):
        findings.append(f"`{folder}` required benchmark command evidence failed or missing.")

    weak_evidence = [
        "Invariant points require direct invariant test evidence.",
        "Iterator behavior points require iterator test evidence.",
    ]
    if not has_bench_sizes:
        weak_evidence.append("Benchmark size coverage missing: 1e3, 1e4, 1e5")
    if not run_commands and has_impl:
        weak_evidence.append("Command execution disabled; strict-hard caps applied for missing command evidence")
    if fallback_reasons:
        weak_evidence.append("Fallback used: " + "; ".join(fallback_reasons))

    hard_fail = False
    hard_fail_reason = "no"
    if not has_impl:
        hard_fail = True
        hard_fail_reason = "yes - no implementation files"
    elif not no_map:
        hard_fail = True
        hard_fail_reason = "yes - forbidden built-in map in implementation"
    elif not no_slice:
        hard_fail = True
        hard_fail_reason = "yes - forbidden slice use in implementation"
    elif api_required and (len(missing_api) / len(api_required) >= 0.5):
        hard_fail = True
        hard_fail_reason = "yes - public API fundamentally incompatible with contract"

    categories = [
        Category(
            "API compliance",
            api_score,
            25,
            "all required API names found" if not missing_api else f"missing methods: {', '.join(missing_api)}",
        ),
        Category(
            "Storage and forbidden-feature compliance",
            storage_score,
            25,
            "no forbidden token match in implementation files"
            if (has_impl and no_map and no_slice)
            else "implementation missing or forbidden features detected",
        ),
        Category(
            "Invariants and behavior",
            behavior_score,
            20,
            "iterator+invariant behavior proven by tests"
            if iterator_test_ok and invariant_test_ok
            else "partial behavior proof from implementation/tests",
        ),
        Category(
            "Test evidence",
            tests_score,
            20,
            "test coverage + iterator/random/invariant evidence"
            if has_tests
            else "no real tests detected",
        ),
        Category(
            "Benchmark and delivery evidence",
            bench_delivery_score,
            10,
            f"benchmark funcs={'yes' if has_bench else 'no'}; benchmark sizes={'yes' if has_bench_sizes else 'no'}",
        ),
    ]

    score = sum(c.score for c in categories)
    letter = grade(score, hard_fail)

    notes: list[str] = []
    if hard_fail:
        notes.append("hard fail")
    if test_required and (not test_ran or not test_ok):
        notes.append("tests failing")
    elif not has_tests:
        notes.append("tests weak")
    if missing_api:
        notes.append("API missing")
    if not has_bench:
        notes.append("no benchmarks")
    if not notes:
        notes.append("evidence consistent")

    return FolderResult(
        folder=folder,
        score=score,
        grade=letter,
        specs_passed=specs_passed,
        specs_total=specs_total,
        specs_pct=specs_pct,
        tests_score=tests_score,
        bench_score=bench_delivery_score,
        notes=", ".join(notes),
        hard_fail=hard_fail_reason,
        categories=categories,
        findings=findings,
        command_evidence=command_evidence,
        weak_evidence=weak_evidence,
        command_status=command_status,
        strict_cap_reasons=strict_cap_reasons,
    )


def render_table(results: list[FolderResult]) -> str:
    lines = [
        "| Folder | Grade | Score | Specs | Tests | Bench | Notes |",
        "|---|---:|---:|---:|---:|---:|---|",
    ]
    for r in results:
        lines.append(
            f"| {r.folder} | {r.grade} | {r.score} | {r.specs_pct} | {r.tests_score} | {r.bench_score} | {r.notes} |"
        )
    return "\n".join(lines)


def render_run_metadata(meta: RunMetadata) -> str:
    lines = [
        "## Run Metadata",
        "",
        f"- Mode: `{meta.mode}`",
        f"- Targets: `{', '.join(meta.targets)}`",
        f"- Command policy: `{meta.command_policy}`",
        f"- Command execution: `{meta.command_execution}`",
    ]
    return "\n".join(lines)


def render_report(repo: Path, results: list[FolderResult], timestamp: str, meta: RunMetadata) -> str:
    lines: list[str] = []
    lines.append("# Implementation Assessment")
    lines.append("")
    lines.append(f"- Repository: `{repo}`")
    lines.append(f"- Timestamp: `{timestamp}`")
    lines.append("- Inputs: `AGENTS.md`, `STRUCTURE-OVERVIEW.md`, each folder `SPECS.md`")
    lines.append("- Policy: strict evidence only, no benefit of doubt")
    lines.append("")
    lines.append(render_run_metadata(meta))
    lines.append("")
    lines.append("## Summary Table")
    lines.append("")
    lines.append(render_table(results))
    lines.append("")
    lines.append("## Folder Findings")
    lines.append("")

    for r in results:
        lines.append(f"### {r.folder}")
        lines.append("")
        lines.append(f"- Grade: `{r.grade}`")
        lines.append(f"- Score: `{r.score}/100`")
        lines.append(f"- Specs passed: `{r.specs_passed}/{r.specs_total} ({r.specs_pct}%)`")
        lines.append(f"- Hard fail: `{r.hard_fail}`")
        lines.append("")
        lines.append("#### Category Breakdown")
        lines.append("")
        lines.append("| Category | Score | Evidence |")
        lines.append("|---|---:|---|")
        for c in r.categories:
            lines.append(f"| {c.name} | {c.score}/{c.max_score} | `{c.evidence}` |")
        lines.append("")
        lines.append("#### Findings")
        lines.append("")
        if r.findings:
            for f in r.findings:
                lines.append(f"- {f}")
        else:
            lines.append("- no major failures detected")
        lines.append("")
        if r.strict_cap_reasons:
            lines.append("#### Strict-Hard Caps Applied")
            lines.append("")
            for reason in r.strict_cap_reasons:
                lines.append(f"- {reason}")
            lines.append("")
        lines.append("#### Command Evidence")
        lines.append("")
        if r.command_evidence:
            for c in r.command_evidence:
                lines.append(f"- {c}")
        else:
            lines.append("- commands skipped")
        lines.append("")
        lines.append("#### Missing Or Weak Evidence")
        lines.append("")
        for w in r.weak_evidence:
            lines.append(f"- {w}")
        lines.append("")

    lines.append("## Human-Only Extra Assessments")
    lines.append("")
    lines.append("- Real developer comprehension study with timed modification tasks.")
    lines.append("- Real workload fitness benchmark using production-like traces.")
    lines.append("- Long-term defect escape rate tracking after merge.")
    lines.append("- Maintenance cost over time across refactors and bug fixes.")
    lines.append("- Learning-value survey with human learners.")
    return "\n".join(lines)


def main() -> int:
    parser = argparse.ArgumentParser(description="Assess implementation folders.")
    parser.add_argument("--repo", default=".", help="Repository root path.")
    parser.add_argument(
        "--folders",
        nargs="*",
        default=[],
        help="Scoped folders. Default: auto-discover from */SPECS.md",
    )
    parser.add_argument(
        "--skip-commands",
        action="store_true",
        help="Do not run make/go verification commands.",
    )
    args = parser.parse_args()

    repo = Path(args.repo).resolve()
    all_folders = discover_folders(repo)
    targets = normalize_targets(args.folders)
    scoped = len(targets) > 0

    if not scoped:
        targets = all_folders

    invalid: list[str] = []
    for t in targets:
        if t.startswith("/") or ".." in t or "/" in t:
            invalid.append(t)
            continue
        if t not in all_folders:
            invalid.append(t)
            continue

    if invalid:
        print("Invalid targets (must be top-level folders with SPECS.md):")
        for t in invalid:
            print(f"- {t}")
        return 2

    run_commands = not args.skip_commands
    mode = "scoped" if scoped else "whole-repo"
    metadata = RunMetadata(
        mode=mode,
        targets=sorted(targets),
        command_policy="strict-hard; make test-folder/bench-folder per folder; fallback only on make infra failure",
        command_execution="enabled" if run_commands else "disabled",
    )

    results = [assess_folder(repo, f, run_commands=run_commands) for f in targets]
    results.sort(key=lambda x: x.folder)

    table = render_table(results)

    print(f"Run mode: {metadata.mode}")
    print(f"Targets: {', '.join(metadata.targets)}")
    print(f"Command policy: {metadata.command_policy}")
    print(f"Command execution: {metadata.command_execution}\n")
    print("## Summary Table\n")
    print(table)

    ts_file = dt.datetime.now().strftime("%Y%m%d_%H%M%S")
    ts_text = dt.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    out_dir = repo / "tmp"
    out_dir.mkdir(parents=True, exist_ok=True)
    out_path = out_dir / f"assessment_{ts_file}.md"
    out_path.write_text(render_report(repo, results, ts_text, metadata), encoding="utf-8")
    print(f"\nReport: {out_path}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
