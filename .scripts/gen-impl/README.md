# gen-impl

Canonical contract: `.scripts/gen-impl/SPECS.md`

`gen.sh` is helper script for spawning an AI agent to generate or regenerate implementation code for this repository's data structure folders.

## Usage

Run one folder:

```bash
./.scripts/gen-impl/gen.sh list-array
```

Run all discovered folders sequentially:

```bash
./.scripts/gen-impl/gen.sh all
```

Supported folders come from root `go.work`.

Required target-folder files:

- `go.mod`
- `README.md`
- `SPECS.md`
- `api_contract.go`

Required module path in `go.mod`:

```text
module github.com/jarollz/data-structure/<folder>
```

Recommended established-folder conventions:

- `helpers_test.go`
- `bench_policy_test.go`
- exactly one `*_api.go`
- exactly one non-test implementation `.go`
- exactly one `*_test.go`
- exactly one `*_bench_test.go`

## AI Spawner Command

By default, the script prompts for the full AI spawner command at runtime:

```text
AI agent spawner command (must have [prompt] in it):
```

Rules:

- The command must contain literal `[prompt]` exactly once.
- `[prompt]` must be unquoted.
- The script probes the command by replacing `[prompt]` with `Reply with exactly this single word and nothing else: Hello`.
- The probe is accepted only if the final non-empty output line is exactly `Hello`.

Example accepted input:

```text
opencode run --dangerously-skip-permissions [prompt]
```

Non-interactive mode:

- Set `AI_SPAWNER_COMMAND` to skip interactive input.
- The command in `AI_SPAWNER_COMMAND` is validated and probed before folder processing starts.
- If validation or probe fails, the script exits immediately (no interactive fallback).

Example:

```bash
AI_SPAWNER_COMMAND='opencode run --dangerously-skip-permissions [prompt]' ./.scripts/gen-impl/gen.sh list-array
```

## Single-Folder Flow

For one folder, the script:

1. Reads `tmp/gen-impl/reports/<folder>/IMPLEMENTATION_REPORT.md` when present.
2. Applies `FORCE` rerun policy (default `FORCE=0`): skip only when report says `SUCCESS` and report file is newer than every non-test implementation `.go` file.
3. When rerun policy selects clean-first mode and implementation `.go` files already exist, runs a reset phase first so next implementation attempt starts fresh.
4. Runs up to `MAX_ATTEMPTS` implementation attempts.
5. Verifies each attempt with doc-comment audit and unit tests.
6. Runs one benchmark validation pass after implementation attempts finish.
7. Runs a final AI report phase to write `tmp/gen-impl/reports/<folder>/IMPLEMENTATION_REPORT.md`.

## all Mode

`all` mode runs the same flow for every folder listed in root `go.work`, in root `go.work` order.

- Folders run sequentially.
- Default behavior is continue on failure.
- Set `STOP_ON_FAILURE=1` to stop early after the first failed folder.

## Rerun Logic

Previous report handling by `FORCE` mode:

- `FORCE=0` (default):
  - `SUCCESS` and fresh: skip folder.
  - otherwise: rerun folder with clean-first reset.
- `FORCE=1`:
  - always rerun folder with clean-first reset.
- `FORCE=2`:
  - if prior report has `Performance Status: FAIL`: rerun folder without clean-first reset (improve in place).
  - otherwise: fallback to `FORCE=0` behavior.
  - missing or malformed report fields: fallback to `FORCE=0` behavior.

## Reset Logic

When rerun policy selects clean-first mode and non-test implementation files already exist, the script first runs a reset prompt that asks the AI to forget old implementation details by replacing implementation bodies with empty stubs or placeholders.

The reset phase must still keep:

- public API signatures unchanged
- protected files untouched
- truthful exported function doc comments

## Protected Files

Implementation attempts must not modify:

- tests
- benchmark tests
- `SPECS.md`
- markdown docs
- `api_contract.go`
- `go.mod`
- `go.sum`
- files outside the target folder

During implementation and reset phases, the report path `tmp/gen-impl/reports/<folder>/IMPLEMENTATION_REPORT.md` is also protected.

If protected or out-of-scope files are edited, the script restores them from a phase snapshot and marks the attempt as failed.

## Doc Comment Rules

The script runs `.scripts/gen-impl/helpers/audit_doc_comments.go` before tests.

Rules enforced by the audit:

- every exported implemented function must have a doc comment
- API interface methods must start with `FuncName implements the API interface.`
- exported non-interface functions such as `New` must use normal truthful Go doc comments
- every exported implemented function comment must include `Example:`

## Verification Flow

Each implementation attempt must pass, in order:

1. scope and protected-file enforcement
2. doc-comment audit
3. `go test -race -json ./<folder>/...`

After implementation attempts finish, the script runs full benchmark validation once for reporting:

4. `make bench-folder FOLDER=<folder>`

Benchmark threshold failures do not trigger implementation retries.

## Report Contract

The final report must be written to `tmp/gen-impl/reports/<folder>/IMPLEMENTATION_REPORT.md` and must contain:

- `Operation Status: SUCCESS` or `Operation Status: FAILURE`
- `Performance Status: PASS`, `FAIL`, or `NOT_RUN`
- `Folder: <folder>`
- `Attempts Used: <n>`
- `## Files Changed` table
- `## Unit Test Summary` table with columns `no. | scenario | pass / fail`
- `## Benchmark Summary` table with columns `no. | scenario | budget-ratio | good/bad | pass / fail`
- `## Operational Failure Causes` table with columns `no. | cause | scenario | evidence | suggestion`
- `### Improvement Suggestions` subsection under Operational Failure Causes with table `no. | cause | failed scenario | suggestion`
- `## Benchmark Failure Causes` table with columns `no. | cause | scenario | evidence | suggestion`
- `### Improvement Suggestions` subsection under Benchmark Failure Causes with table `no. | cause | failed scenario | suggestion`

`budget-ratio > 1.0` is treated as `BAD` and `FAIL` in benchmark summary.

Failure-cause tables are always present; when no failures are recorded, they contain a `(none)` row.

## Output and Logs

Run artifacts are written under:

```text
tmp/gen-impl/runs/<timestamp>/
```

Artifacts include prompts, AI output logs, validation logs, attempt summaries, per-folder summaries, and an aggregate run summary.

## Environment Variables

- `FORCE=0`: default behavior; skip only fresh `SUCCESS`, otherwise rerun clean-first.
- `FORCE=1`: always rerun clean-first.
- `FORCE=2`: rerun in place when prior report performance status is `FAIL`; otherwise fallback to `FORCE=0`.
- `STOP_ON_FAILURE=1`: stop early in `all` mode after first failure.
- `MAX_ATTEMPTS=5`: override implementation retry count.
- `AI_SPAWNER_COMMAND`: set full AI spawner command template and run non-interactively.
- `PROBE_TIMEOUT_SECONDS`: override spawner probe timeout.
- `SPAWNER_TIMEOUT_SECONDS`: override AI run timeout.
- `SPAWNER_IDLE_TIMEOUT_SECONDS`: fail AI run if no new output appears for this many seconds.
- `DOC_AUDIT_TIMEOUT_SECONDS`: override doc-comment audit timeout.
- `TEST_TIMEOUT_SECONDS`: override unit test timeout.
- `BENCH_TIMEOUT_SECONDS`: override benchmark timeout.

## Troubleshooting

If the spawner command probe fails:

- make sure `[prompt]` is present exactly once and unquoted
- make sure the command really runs the AI agent
- make sure the final output line is exactly `Hello`

If an attempt fails repeatedly:

- read the latest attempt summary under `tmp/gen-impl/runs/.../<folder>/attempt_<n>/summary.md`
- read the corresponding test or benchmark logs
- review the previous report under `tmp/gen-impl/reports/<folder>/IMPLEMENTATION_REPORT.md` before rerunning

For exact behavior, exact output text, and full reconstruction contract, read `.scripts/gen-impl/SPECS.md`.
