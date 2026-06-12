---
name: assess-implementation
description: Use when grading data structure implementations against AGENTS.md, STRUCTURE-OVERVIEW.md, and per-folder SPECS.md. Supports whole-repo grading or exact scoped requests like `invoke the skill assess-implementation only for folder list-array` and `invoke the skill assess-implementation only for folders list-array, stack`. Prints compact markdown grading table to terminal and writes full markdown report to tmp/assessment_YYYYMMDD_HHMMSS.md.
---

# Assess Implementation

## Purpose

Perform strict evidence-based grading for every data structure folder in this repository, or for selected target folders when user requests scoped assessment.
Compare implementation, tests, benchmarks, and required files against:

- root `AGENTS.md`
- root `STRUCTURE-OVERVIEW.md`
- each structure folder `SPECS.md`

No benefit of doubt. Missing proof counts as failure.

## Scope

Default mode: assess every top-level structure folder discovered from root `go.work` that also contains `SPECS.md`.

Scoped mode: if user says `invoke the skill assess-implementation only for folder xxxx` or `invoke the skill assess-implementation only for folders xxxx, yyyy`, assess only requested folders.

Scoped mode rules:

- Each requested target may be exact top-level structure folder name, for example `list-array` or `tree-avl`.
- Each requested target may also be repo-relative folder path pointing to top-level structure folder, for example `./list-array`.
- Normalize repo-relative forms before validation. Example: `./list-array` becomes `list-array`.
- Accept multiple requested targets in comma-separated, space-separated, or individually quoted/backticked forms when intent is clear.
- Deduplicate normalized target folders before assessment.
- Reject absolute paths, parent traversal such as `../`, and non-top-level paths.
- Assessable target must be registered in root `go.work` and contain its own `SPECS.md`.
- If requested target is listed in root `go.work` but missing `SPECS.md`, skip it with warning.
- If requested target contains `SPECS.md` but is not registered in root `go.work`, skip it with warning.
- If every requested target is skipped, print warnings, print `No assessable folders.`, write report, and exit nonzero.
- In scoped mode, do not assess any structure folders outside selected target set.
- In scoped mode, terminal table contains data rows for selected targets only.
- In scoped mode, full markdown report contains summary and findings for selected targets only.
- If requested target is missing or invalid, stop and report invalid targets instead of silently falling back to whole-repo assessment.

Discovery rules:

- Run from repo root.
- Read root `go.work`.
- Assess folders in root `go.work` order only.
- Assess folder only when top-level folder is registered in root `go.work` and also contains `SPECS.md`.
- If top-level folder is listed in root `go.work` but missing `SPECS.md`, skip it with warning.
- If top-level folder contains `SPECS.md` but is not listed in root `go.work`, skip it with warning.

## Invocation Examples

- Whole repo: `invoke the skill assess-implementation`
- Single folder: `invoke the skill assess-implementation only for folder list-array`
- Single folder: `invoke the skill assess-implementation only for folder ./list-array`
- Single folder: `invoke the skill assess-implementation only for folder map-hash`
- Multiple folders: `invoke the skill assess-implementation only for folders list-array, stack`
- Multiple folders: `invoke the skill assess-implementation only for folders ./list-array, ./queue, tree-avl`

## Bundled Resources And Reuse Policy

- Persistent helper script path: `.agents/skills/assess-implementation/resources/assess_impl.py`
- Helper docs path: `.agents/skills/assess-implementation/resources/README.md`
- Reuse bundled helper for normal assessment and report generation.
- Do not create ad-hoc Python scripts in random paths or `tmp/` for routine runs.
- Temporary debug scripts are allowed only for parser debugging. Delete them immediately after debugging and state why they were needed.

Preferred helper commands:

- Whole repo: `python3 .agents/skills/assess-implementation/resources/assess_impl.py`
- Scoped: `python3 .agents/skills/assess-implementation/resources/assess_impl.py --folders list-array stack`
- No local command execution: `python3 .agents/skills/assess-implementation/resources/assess_impl.py --skip-commands`

## Evidence Rules

- Read root `AGENTS.md` first.
- Read root `STRUCTURE-OVERVIEW.md` second.
- Read each folder `SPECS.md` before grading that folder.
- Read each folder `api_contract.go` and use it to cross-check iterator contract.
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
- `SPECS.md` exists.
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

- Public constructor and methods required by folder `SPECS.md` exist.
- Naming matches contract exactly.
- Iterator methods come from folder contracts, not hardcoded folder-name guesses.
- Iterator naming matches root `AGENTS.md` and primary iterator in `STRUCTURE-OVERVIEW.md`.
- Return shapes and edge-case behavior are supported by code and tests.

### 5. Invariants And Behavior

- Internal representation matches folder rules.
- Resize policy or node-pool policy matches folder rules.
- Invariant-preserving logic exists in implementation.
- Iterator order, count, early-stop, and mutation-unsafety are covered by code and tests.
- Edge cases from `SPECS.md` are handled and tested.

### 6. Test Evidence

- Unit tests cover every public API.
- Tests validate invariants.
- Tests validate iterator contract.
- Tests include checklist items from `SPECS.md`.
- Randomized or oracle-based tests exist when required.
- Placeholder tests score zero.

### 7. Benchmark Evidence

- Benchmarks exist.
- Benchmarks cover workloads required by `SPECS.md`.
- Benchmarks cover multiple sizes when required by root `AGENTS.md`.

### 8. Command Evidence

- If implementation exists, run local verification commands when available.
- Prefer Makefile tools over raw `go test` when available.
- Strict-hard command mode uses per-folder commands for both whole-repo and scoped assessment: `make test-folder FOLDER=<folder>`, then optional `make bench-folder FOLDER=<folder>`.
- If Makefile target is unavailable or broken, fallback to direct `go test` command for equivalent scope.
- If benchmarks exist and runtime is reasonable, run benchmark smoke check.
- Record exact command and result in report.
- Command success does not override rule failures.
- Command failure is evidence, not a blocker for finishing assessment.

Strict-hard caps tied to command evidence:

- If required test command fails or is missing, cap `Test evidence` at `4/20` and `Invariants and behavior` at `6/20`.
- If required benchmark command fails or is missing, cap `Benchmark and delivery evidence` at `5/10`.
- Missing command execution when implementation exists is treated as failed command evidence, not neutral.

## SPECS Parsing Discipline

- Parse checklist items only from markdown checkbox bullets: `- [ ] ...`.
- Parse required API names only from backticked signatures in `## Required API` section.
- Do not treat prose tokens as API names. Ignore complexity markers and type terms such as `O(1)` and `float64(...)` unless they are part of explicit backticked API signature.
- If parser confidence is low, mark evidence missing instead of guessing.

## Scoring Rubric

Score each structure from `0` to `100`.

| Category | Weight | What to check |
|---|---:|---|
| API compliance | 25 | constructor, methods, iterator naming, signature shape, contract behavior |
| Storage and forbidden-feature compliance | 25 | arrays-only rule, no `map`, no slice use in implementation code, representation contract |
| Invariants and behavior | 20 | correctness logic, resize policy, iterator contract, edge cases |
| Test evidence | 20 | API coverage, invariant tests, iterator tests, randomized/oracle tests |
| Benchmark and delivery evidence | 10 | benchmarks, multiple sizes, README/SPECS/tests/bench presence |

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

- `Specs passed`: number of satisfied checklist items found in that folder `SPECS.md`
- `Specs total`: total checklist items in that folder `SPECS.md`
- `Compliance % = rules passed / rules total * 100`

## Required Output

Produce both outputs in one run:

1. Print compact markdown table to terminal.
2. Write full markdown report to `tmp/assessment_YYYYMMDD_HHMMSS.md`.

Terminal table output is mandatory. Do not consider assessment complete unless terminal shows the summary table rows.

If command streaming hides or drops stdout table in agent view, recover by reading generated report and printing the `## Summary Table` block to terminal before finishing.

In scoped mode, both outputs must contain only requested folders.

If any folder is skipped, both terminal output and markdown report must include warning bullets naming every skipped folder and reason.

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

| Folder | Grade | Score | Specs | Tests | Bench | Notes |
|---|---:|---:|---:|---:|---:|---|

Column meaning:

- `Score`: numeric score out of `100`
- `Specs`: compliance percent rounded to whole number
- `Tests`: test evidence score out of `20`
- `Bench`: benchmark and delivery score out of `10`
- `Notes`: short failure summary only

Keep `Notes` brief, for example:

- `no impl`
- `uses slice`
- `API missing, no benchmarks`
- `tests weak, iterator unproven`

When responding after assessment run, include this exact markdown table block in terminal-visible output. Do not replace with prose summary.

## Full Report Format

Write markdown file with this structure:

```md
# Implementation Assessment

- Repository: <repo path or name>
- Timestamp: <ISO-like timestamp>
- Inputs: `AGENTS.md`, `STRUCTURE-OVERVIEW.md`, each folder `SPECS.md`
- Policy: strict evidence only, no benefit of doubt

## Skipped Folders

- `folder`: skipped reason

## Summary Table

| Folder | Grade | Score | Specs | Tests | Bench | Notes |
|---|---:|---:|---:|---:|---:|---|
| ... |

## Folder Findings

### <folder>

- Grade: `F`
- Score: `2/100`
- Specs passed: `2/31 (6%)`
- Hard fail: `yes - no implementation files`

#### Category Breakdown

| Category | Score | Evidence |
|---|---:|---|
| API compliance | 0/25 | `no implementation files` |
| Storage and forbidden-feature compliance | 0/25 | `not assessable because implementation missing` |
| Invariants and behavior | 0/20 | `no code, no proofs` |
| Test evidence | 0/20 | `nothing_test.go` only; placeholder |
| Benchmark and delivery evidence | 2/10 | `README.md` and `SPECS.md` exist; benchmarks missing |

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
5. Read root `go.work` and discover top-level registered folders in root `go.work` order.
6. Discover top-level folders containing `SPECS.md`.
7. Compute assessable folders as intersection of root `go.work` folders and top-level `SPECS.md` folders.
8. Record skipped-folder warnings for mismatches from either side.
9. If scoped mode requested, keep requested assessable targets only; requested mismatches become warnings, invalid targets still fail.
7. Run bundled helper script from `.agents/skills/assess-implementation/resources/assess_impl.py`.
8. Run per-folder strict-hard command checks: `make test-folder FOLDER=<folder>` and optional `make bench-folder FOLDER=<folder>`.
9. Apply strict-hard command-evidence score caps when required commands fail or are missing.
10. Ensure parser uses strict checkbox and backticked-signature extraction rules.
11. Print run metadata, skipped-folder warnings when present, and markdown summary table to terminal.
12. Write full markdown report to `tmp/assessment_YYYYMMDD_HHMMSS.md` including run metadata, skipped folders when present, and cap reasons.
13. Verify terminal output contains summary table header; if missing, print table recovered from report.
14. End with short list of human-only extra assessments.
