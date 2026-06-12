# gen-impl

`gen.sh` is helper script for spawning an AI agent to generate or regenerate implementation code for this repository's data structure folders.

## Usage

Run one folder:

```bash
./.scripts/gen-impl/gen.sh list-array
```

Run all supported folders sequentially:

```bash
./.scripts/gen-impl/gen.sh all
```

Supported folders:

- `list-array`
- `list-linked-singly`
- `list-linked-doubly`
- `list-skip`
- `queue`
- `stack`
- `heap`
- `tree-general`
- `tree-avl`
- `tree-red-black`
- `map-hash`
- `map-trie`
- `map-tree-avl`
- `map-tree-red-black`

## AI Spawner Command

The script prompts for the full AI spawner command at runtime:

```text
type in the AI agent spawner command (must have "[prompt]" in it):
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

## Single-Folder Flow

For one folder, the script:

1. Reads `tmp/gen-impl/reports/<folder>/IMPLEMENTATION_REPORT.md` when present.
2. Skips the folder only when the report says `SUCCESS` and the report file is newer than every non-test implementation `.go` file.
3. If the report is missing or says `FAILURE` and implementation `.go` files already exist, runs a reset phase first so the next implementation attempt starts fresh.
4. Runs up to `MAX_ATTEMPTS` implementation attempts.
5. Verifies each attempt with doc-comment audit, unit tests, and benchmark tests.
6. Runs a final AI report phase to write `tmp/gen-impl/reports/<folder>/IMPLEMENTATION_REPORT.md`.

## all Mode

`all` mode runs the same flow for every supported folder in fixed order.

- Folders run sequentially.
- Default behavior is continue on failure.
- Set `STOP_ON_FAILURE=1` to stop early after the first failed folder.

## Rerun Logic

Previous report handling:

- `SUCCESS` and fresh: skip folder.
- `SUCCESS` but stale: rerun folder.
- `FAILURE`: rerun folder and feed prior failure suggestions into the next prompt.
- Missing or malformed report: rerun folder.

Set `FORCE=1` to ignore a fresh `SUCCESS` report and rerun anyway.

## Reset Logic

When the report is missing or says `FAILURE`, and non-test implementation files already exist, the script first runs a reset prompt that asks the AI to forget the old implementation by replacing implementation bodies with empty stubs or placeholders.

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
3. `make test-folder FOLDER=<folder>`
4. `make bench-folder FOLDER=<folder>`

## Report Contract

The final report must be written to `tmp/gen-impl/reports/<folder>/IMPLEMENTATION_REPORT.md` and must contain:

- `Operation Status: SUCCESS` or `Operation Status: FAILURE`
- `Folder: <folder>`
- `Attempts Used: <n>`
- `## Files Changed` table
- `## Failure Causes` table when status is `FAILURE`

Each failure cause row must include a suggestion.

## Output and Logs

Run artifacts are written under:

```text
tmp/gen-impl/runs/<timestamp>/
```

Artifacts include prompts, AI output logs, validation logs, attempt summaries, per-folder summaries, and an aggregate run summary.

## Environment Variables

- `FORCE=1`: ignore fresh `SUCCESS` report and rerun.
- `STOP_ON_FAILURE=1`: stop early in `all` mode after first failure.
- `MAX_ATTEMPTS=5`: override implementation retry count.
- `PROBE_TIMEOUT_SECONDS`: override spawner probe timeout.
- `SPAWNER_TIMEOUT_SECONDS`: override AI run timeout.
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
