---
name: assess-implementation
description: Use when grading data structure implementations against AGENTS.md, STRUCTURE-CONTRACTS.md, and per-folder RULES.md. Supports whole-repo grading or exact scoped requests like `invoke the skill assess-implementation only for folder list-array` and `invoke the skill assess-implementation only for folders list-array, stack`. Prints compact markdown grading table to terminal and writes full markdown report to tmp/assessment_YYYYMMDD_HHMMSS.md.
---

# Assess Implementation

## Purpose

Perform strict evidence-based grading for every data structure folder in this repository, or for selected target folders when user requests scoped assessment.
Compare implementation, tests, benchmarks, and required files against:

- root `AGENTS.md`
- root `STRUCTURE-CONTRACTS.md`
- each structure folder `RULES.md`

No benefit of doubt. Missing proof counts as failure.

## Scope

Default mode: assess every top-level structure folder discovered from `*/RULES.md`.

Scoped mode: if user says `invoke the skill assess-implementation only for folder xxxx` or `invoke the skill assess-implementation only for folders xxxx, yyyy`, assess only requested folders.

Scoped mode rules:

- Each requested target may be exact top-level structure folder name, for example `list-array` or `tree-avl`.
- Each requested target may also be repo-relative folder path pointing to top-level structure folder, for example `./list-array`.
- Normalize repo-relative forms before validation. Example: `./list-array` becomes `list-array`.
- Accept multiple requested targets in comma-separated, space-separated, or individually quoted/backticked forms when intent is clear.
- Deduplicate normalized target folders before assessment.
- Reject absolute paths, parent traversal such as `../`, and non-top-level paths.
- Every requested target folder must exist and contain its own `RULES.md`.
- In scoped mode, do not assess any structure folders outside selected target set.
- In scoped mode, terminal table contains data rows for selected targets only.
- In scoped mode, full markdown report contains summary and findings for selected targets only.
- If any requested target is missing, ambiguous, or lacks `RULES.md`, stop and report invalid targets instead of silently falling back to whole-repo assessment.

Assess every top-level structure folder discovered from `*/RULES.md`, including:

- `map-hash`
- `map-tree-avl`
- `map-tree-red-black`
- `heap`
- `list-array`
- `list-linked-singly`
- `list-linked-doubly`
- `list-skip`
- `queue`
- `stack`
- `tree-general`
- `tree-avl`
- `tree-red-black`

If more structure folders exist later, include them automatically when they contain `RULES.md`.

## Invocation Examples

- Whole repo: `invoke the skill assess-implementation`
- Single folder: `invoke the skill assess-implementation only for folder list-array`
- Single folder: `invoke the skill assess-implementation only for folder ./list-array`
- Single folder: `invoke the skill assess-implementation only for folder map-hash`
- Multiple folders: `invoke the skill assess-implementation only for folders list-array, stack`
- Multiple folders: `invoke the skill assess-implementation only for folders ./list-array, ./queue, tree-avl`

## Evidence Rules

- Read root `AGENTS.md` first.
- Read root `STRUCTURE-CONTRACTS.md` second.
- Read each folder `RULES.md` before grading that folder.
- Inspect implementation `.go` files separately from tests and benchmarks.
- Treat files ending in `_test.go` as tests only.
- Treat benchmark functions as real only when they use `func Benchmark...`.
- Placeholder files such as `nothing_test.go`, trivial `TestNothing`, or equivalent dummy checks do not satisfy test requirements.
- Missing implementation files are hard failures.
- Missing tests are hard failures for validation categories.
- Missing benchmarks are hard failures for benchmark category.
- If behavior is claimed in README or comments but not proved by implementation or tests, mark failed.
- When local commands are available, run cheap verification commands and record exact results as additional evidence.
- Cite concrete evidence with file paths and line numbers in detailed findings.

## Hard Fail Rules

Assign overall grade `F` immediately when any of these is true:

- No implementation code exists for structure.
- Implementation uses forbidden built-in `map`.
- Implementation uses forbidden slice syntax or slice-backed logic in implementation code.
- Public API is fundamentally absent or incompatible with folder contract.

Still compute numeric category scores. Grade remains `F`.

## Required Checks

For each structure folder, verify all of these.

### 1. Required Delivery

- `README.md` exists.
- `RULES.md` exists.
- Real tests exist.
- Real benchmarks exist.

### 2. Implementation Presence

- Non-test `.go` files exist.
- Files define actual data structure implementation, not stubs only.

### 3. Forbidden Feature Compliance

- Implementation code does not use built-in `map`.
- Implementation code does not use slice syntax or slice-backed logic anywhere.
- Implementation code does not hide forbidden logic behind helper containers.
- Storage model matches repo and folder contract.

Check implementation files only. Ignore test-only use of slices or maps.

### 4. API Compliance

- Public constructor and methods required by folder `RULES.md` exist.
- Naming matches contract exactly.
- Iterator naming matches root `AGENTS.md` and `STRUCTURE-CONTRACTS.md`.
- Return shapes and edge-case behavior are supported by code and tests.

### 5. Invariants And Behavior

- Internal representation matches folder rules.
- Resize policy or node-pool policy matches folder rules.
- Invariant-preserving logic exists in implementation.
- Iterator order, count, early-stop, and mutation-unsafety are covered by code and tests.
- Edge cases from `RULES.md` are handled and tested.

### 6. Test Evidence

- Unit tests cover every public API.
- Tests validate invariants.
- Tests validate iterator contract.
- Tests include checklist items from `RULES.md`.
- Randomized or oracle-based tests exist when required.
- Placeholder tests score zero.

### 7. Benchmark Evidence

- Benchmarks exist.
- Benchmarks cover workloads required by `RULES.md`.
- Benchmarks cover multiple sizes when required by root `AGENTS.md`.

### 8. Command Evidence

- If implementation exists, run local verification commands when available.
- For Go folders, prefer `go test` for that folder.
- If benchmarks exist and runtime is reasonable, run a short benchmark smoke check.
- Record exact command and result in report.
- Command success does not override rule failures.
- Command failure is evidence, not a blocker for finishing assessment.

## Scoring Rubric

Score each structure from `0` to `100`.

| Category | Weight | What to check |
|---|---:|---|
| API compliance | 25 | constructor, methods, iterator naming, signature shape, contract behavior |
| Storage and forbidden-feature compliance | 25 | arrays-only rule, no `map`, no slice use in implementation code, representation contract |
| Invariants and behavior | 20 | correctness logic, resize policy, iterator contract, edge cases |
| Test evidence | 20 | API coverage, invariant tests, iterator tests, randomized/oracle tests |
| Benchmark and delivery evidence | 10 | benchmarks, multiple sizes, README/RULES/tests/bench presence |

Category scoring must be harsh:

- Full points only when code or tests directly prove compliance.
- Partial points only when some checklist items are satisfied and evidence is direct.
- Zero points when evidence is missing, dummy, or contradicted.

Letter grades:

- `A`: `90-100`
- `B`: `80-89`
- `C`: `70-79`
- `D`: `60-69`
- `F`: `<60` or any hard fail rule triggered

Also calculate:

- `Rules passed`: number of satisfied checklist items found in that folder `RULES.md`
- `Rules total`: total checklist items in that folder `RULES.md`
- `Compliance % = rules passed / rules total * 100`

## Required Output

Produce both outputs in one run:

1. Print compact markdown table to terminal.
2. Write full markdown report to `tmp/assessment_YYYYMMDD_HHMMSS.md`.

In scoped mode, both outputs must contain only requested folders.

Create `tmp/` first if missing.

Filename pattern:

- `assessment_YYYYMMDD_HHMMSS.md`

Example timestamp command:

```sh
date +"%Y%m%d_%H%M%S"
```

## Terminal Table Format

Keep terminal table narrow so plain terminals can still read it.

Use this exact column set:

| Folder | Grade | Score | Rules | Tests | Bench | Notes |
|---|---:|---:|---:|---:|---:|---|

Column meaning:

- `Score`: numeric score out of `100`
- `Rules`: compliance percent rounded to whole number
- `Tests`: test evidence score out of `20`
- `Bench`: benchmark and delivery score out of `10`
- `Notes`: short failure summary only

Keep `Notes` brief, for example:

- `no impl`
- `uses slice`
- `API missing, no benchmarks`
- `tests weak, iterator unproven`

## Full Report Format

Write markdown file with this structure:

```md
# Implementation Assessment

- Repository: <repo path or name>
- Timestamp: <ISO-like timestamp>
- Inputs: `AGENTS.md`, `STRUCTURE-CONTRACTS.md`, each folder `RULES.md`
- Policy: strict evidence only, no benefit of doubt

## Summary Table

| Folder | Grade | Score | Rules | Tests | Bench | Notes |
|---|---:|---:|---:|---:|---:|---|
| ... |

## Folder Findings

### <folder>

- Grade: `F`
- Score: `2/100`
- Rules passed: `2/31 (6%)`
- Hard fail: `yes - no implementation files`

#### Category Breakdown

| Category | Score | Evidence |
|---|---:|---|
| API compliance | 0/25 | `no implementation files` |
| Storage and forbidden-feature compliance | 0/25 | `not assessable because implementation missing` |
| Invariants and behavior | 0/20 | `no code, no proofs` |
| Test evidence | 0/20 | `nothing_test.go` only; placeholder |
| Benchmark and delivery evidence | 2/10 | `README.md` and `RULES.md` exist; benchmarks missing |

#### Findings

- `path:line` failed rule summary.
- `path:line` failed rule summary.

#### Command Evidence

- `go test ./...` -> failed or skipped with reason.

#### Missing Or Weak Evidence

- specific missing tests
- specific missing benchmarks
- specific missing invariants or iterator proofs
```

Fix broken example details if they do not match actual findings. Never leave contradictory sample text in final report.

## Review Tone

- Be blunt.
- No praise unless a category genuinely earns full credit.
- Prefer short factual statements.
- Findings first. No motivational filler.

## Human-Only Extra Assessments

Suggest these only after completing repo-based grading. These are outside direct AI proof unless humans provide extra data.

- Real developer comprehension study: measure how quickly humans understand and modify each implementation.
- Real workload fitness: benchmark against actual usage traces or production-like traffic supplied by humans.
- Long-term defect escape rate: track bugs discovered after merge or during later use.
- Maintenance cost over time: track refactor effort, review time, and bug-fix effort across months.
- Learning-value assessment: ask human learners whether each implementation teaches intended concepts clearly.

## Execution Checklist

1. Read root contracts.
2. Detect whether user requested scoped mode with one or more exact folder names or repo-relative folder paths.
3. If scoped mode requested, normalize allowed repo-relative forms to canonical top-level folder names.
4. If scoped mode requested, deduplicate targets.
5. If scoped mode requested, validate every target folder has `RULES.md`; otherwise stop with invalid-target message listing bad targets.
6. If scoped mode not requested, enumerate structure folders from `*/RULES.md`.
7. Inspect implementation, tests, benchmarks, README, and `RULES.md` for each selected folder only.
8. Run cheap local verification commands when implementation exists.
9. Compute checklist pass counts and rubric scores.
10. Print markdown summary table to terminal.
11. Write full markdown report to `tmp/assessment_YYYYMMDD_HHMMSS.md`.
12. End with short list of human-only extra assessments.
