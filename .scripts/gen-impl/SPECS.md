# gen-impl SPECS

This file is canonical contract for `.scripts/gen-impl`.

If implementation, `README.md`, prompt templates, comments, or generated examples conflict with this file, this file wins.

## Purpose

`gen-impl` generates or regenerates one data structure folder, or all valid data structure folders, by prompting an external AI spawner command, then verifying and reporting results.

## Scope

The package scope is only `.scripts/gen-impl/` runtime code plus its output under `tmp/gen-impl/`.

The package must not maintain alternate legacy shell-helper codepaths. Shell entrypoint is kept only as thin wrapper `gen.sh`.

## Required Package Layout

These files must exist:

- `.scripts/gen-impl/gen.sh`
- `.scripts/gen-impl/main.py`
- `.scripts/gen-impl/README.md`
- `.scripts/gen-impl/SPECS.md`
- `.scripts/gen-impl/core/__init__.py`
- `.scripts/gen-impl/core/config.py`
- `.scripts/gen-impl/core/io_safe.py`
- `.scripts/gen-impl/core/process_control.py`
- `.scripts/gen-impl/core/progress.py`
- `.scripts/gen-impl/core/reporting.py`
- `.scripts/gen-impl/core/runner.py`
- `.scripts/gen-impl/core/scope_guard.py`
- `.scripts/gen-impl/core/spawner.py`
- `.scripts/gen-impl/core/verify.py`
- `.scripts/gen-impl/helpers/audit_doc_comments.go`
- `.scripts/gen-impl/helpers/prompt_impl.md.tmpl`
- `.scripts/gen-impl/helpers/prompt_report.md.tmpl`
- `.scripts/gen-impl/helpers/prompt_reset.md.tmpl`

These files must not exist:

- `.scripts/gen-impl/helpers/common.sh`
- `.scripts/gen-impl/helpers/folders.sh`
- `.scripts/gen-impl/helpers/report.sh`
- `.scripts/gen-impl/helpers/spawner.sh`
- `.scripts/gen-impl/helpers/verify.sh`

## CLI Contract

Exact command:

```text
./.scripts/gen-impl/gen.sh <folder|all>
```

`gen.sh` must stay a thin shell wrapper that executes `main.py` with `python3` and forwards args unchanged.

If argument count is not exactly one target arg, the program must print this exact usage text to stderr and exit status `1`:

```text
Usage: ./.scripts/gen-impl/gen.sh <folder|all>

Examples:
  ./.scripts/gen-impl/gen.sh list-array
  ./.scripts/gen-impl/gen.sh all

Environment:
  FORCE=0             Default. Skip fresh SUCCESS; otherwise rerun clean-first.
  FORCE=1             Always rerun and clean-first reset.
  FORCE=2             Rerun in-place when prior Performance Status is FAIL.
  STOP_ON_FAILURE=1   Stop after first failed folder in all mode.
  MAX_ATTEMPTS=5      Maximum implementation attempts per folder.
  SPAWNER_TIMEOUT_SECONDS
                      Hard timeout for each AI spawner run.
  SPAWNER_IDLE_TIMEOUT_SECONDS
                      Fail AI spawner run when output stays idle.
  AI_SPAWNER_COMMAND  Full spawner command template with [prompt].
```

## Environment Variable Contract

Defaults:

- `MAX_ATTEMPTS=5`
- `REPORT_ATTEMPTS=5`
- `STOP_ON_FAILURE=0`
- `FORCE=0`
- `PROBE_TIMEOUT_SECONDS=120`
- `SPAWNER_TIMEOUT_SECONDS=1800`
- `SPAWNER_IDLE_TIMEOUT_SECONDS=180`
- `DOC_AUDIT_TIMEOUT_SECONDS=300`
- `TEST_TIMEOUT_SECONDS=900`
- `BENCH_TIMEOUT_SECONDS=1800`
- `AI_SPAWNER_COMMAND=""`

`FORCE` accepts only `0`, `1`, or `2`.

Invalid `FORCE` must raise exact message template:

```text
invalid {name} value '{raw}'; expected one of: 0, 1, 2
```

Program-level fatal errors must be printed through exact prefix template:

```text
[gen-impl] Error: {message}
```

## Folder Discovery Contract

Source of truth is root `go.work` only.

Discovery mechanism must run exact command:

```text
go work edit -json
```

Execution rules:

- run command with cwd at repo root
- parse only `Use[].DiskPath`
- preserve `go.work` order exactly
- normalize `./folder` to `folder`
- reject duplicate normalized folders

Each `DiskPath` must satisfy all rules:

- non-empty string
- relative path only
- must not contain `..`
- must resolve inside repo root
- must resolve to a top-level folder exactly one path segment deep
- resolved directory must exist

Exact discovery failure templates:

```text
missing root go.work: {go_work_path}
failed to run 'go work edit -json': {exc}
'go work edit -json' failed: {reason}
invalid JSON from 'go work edit -json': {exc}
root go.work contains no use entries
root go.work contains an invalid use entry
root go.work contains a use entry without DiskPath
root go.work contains an empty use entry
root go.work entry '{disk_path}' must be relative
root go.work entry '{disk_path}' must not contain '..'
root go.work entry '{disk_path}' points outside repo root
root go.work entry '{disk_path}' must point to a top-level folder
root go.work entry '{disk_path}' points to missing directory: {folder}
root go.work lists duplicate folder '{folder}'
```

Raw error passthrough in `{exc}` and `{reason}` is allowed.

## Target Resolution Contract

Exact target rules:

- `all` means every discovered folder in root `go.work` order
- any other target must exactly match one discovered folder

Exact invalid target error template:

```text
target '{target}' is not listed in root go.work
```

## Folder Validation Contract

### Hard Fail Rules

Each discovered target folder must contain these exact files:

- `go.mod`
- `SPECS.md`
- `api_contract.go`
- `README.md`

Exact missing-file error template:

```text
missing required file {folder}/{required_file}
```

Each folder `go.mod` must contain exactly one module declaration matching:

```text
module github.com/jarollz/data-structure/{folder}
```

Exact module validation failure templates:

```text
missing module declaration in {folder}/go.mod
invalid module path in {folder}/go.mod: expected 'module github.com/jarollz/data-structure/{folder}'
```

### Soft Warning Rules

Warnings must not stop execution.

Recommended files and counts, derived from all established structure folders at time this spec was written:

- `helpers_test.go`
- `bench_policy_test.go`
- exactly one `*_api.go`
- exactly one non-test implementation `.go`
- exactly one `*_test.go`
- exactly one `*_bench_test.go`

Exact warning templates:

```text
[gen-impl] Warning: folder {folder} missing recommended file {required_name}
[gen-impl] Warning: folder {folder} expected exactly one *_api.go file, found {count}
[gen-impl] Warning: folder {folder} expected exactly one non-test implementation .go file, found {count}
[gen-impl] Warning: folder {folder} expected exactly one *_test.go file, found {count}
[gen-impl] Warning: folder {folder} expected exactly one *_bench_test.go file, found {count}
```

## Output Prefix Contract

Exact stdout prefix template:

```text
[gen-impl] {message}
```

Exact stderr warning prefix template:

```text
[gen-impl] Warning: {message}
```

Exact stderr error prefix template:

```text
[gen-impl] Error: {message}
```

## Run Artifact Contract

Timestamp format must be UTC and exact pattern:

```text
%Y%m%d_%H%M%S
```

Run root path:

```text
tmp/gen-impl/runs/{timestamp}/
```

Report path:

```text
tmp/gen-impl/reports/{folder}/IMPLEMENTATION_REPORT.md
```

Required aggregate run files:

- `summary.tsv`
- `summary.md`

Required per-folder artifacts:

- `summary.md`
- `prior_IMPLEMENTATION_REPORT.md`
- `prior_failure_suggestions.md`
- `no_previous_failure_summary.md`
- `change_table.md`
- `operational_failure_causes.tsv`
- `benchmark_failure_causes.tsv`
- `tests.log`
- `benchmarks.log`
- `report_input.md`

Probe artifacts:

- interactive: `spawner_probe_{attempt}/probe_raw.log`, `probe_clean.log`
- env preset: `spawner_probe_env/probe_raw.log`, `probe_clean.log`

## Writable Path Contract

Exact writable-dir failure template:

```text
{label} not writable: {path}: {exc}
```

Exact report-path write failure stderr template:

```text
[gen-impl] Error: cannot write report path {report_path}: {reason}
```

If report path is not writable, program must:

1. print synthesized report content to stdout first
2. print exact stderr error line above second
3. mark folder fatal with note `report path not writable`

Raw OS error passthrough in `{exc}` and `{reason}` is allowed.

## Spawner Command Contract

The spawner command template must contain literal `[prompt]` exactly once and unquoted.

Exact validation reason strings:

```text
empty command
[prompt] must be unquoted
missing [prompt] placeholder
[prompt] must appear exactly once
```

Interactive prompt text must be exact:

```text
AI agent spawner command (must have [prompt] in it):
```

Probe prompt must be exact:

```text
Reply with exactly this single word and nothing else: Hello
```

Probe success rule:

- strip ANSI escape sequences
- replace carriage returns with newlines
- final non-empty output line must equal exact literal `Hello`

Exact spawner-related fatal templates:

```text
failed to read spawner command from stdin
AI_SPAWNER_COMMAND invalid: {reason}
AI_SPAWNER_COMMAND probe command failed or final output was not Hello. Review {probe_clean_log}
```

Exact interactive warning template after failed probe:

```text
[gen-impl] Warning: probe command failed or final output was not Hello. Review {probe_clean_log} and try again.
```

## Phase Flow Contract

Per run:

1. create run root and reports root
2. resolve targets
3. acquire valid spawner command
4. process targets sequentially
5. write aggregate `summary.md`
6. print aggregate summary table
7. print run-artifacts log line

Per folder:

1. hard-validate layout
2. emit soft warnings for convention drift
3. load prior report metadata
4. apply rerun policy
5. if skipped, write folder summary and summary row
6. snapshot folder start
7. optional reset phase when clean-first and implementation files already exist
8. run up to `MAX_ATTEMPTS` implementation attempts
9. after operational success, run final benchmark validation once
10. write `report_input.md`
11. run report-generation phase up to `REPORT_ATTEMPTS`
12. write folder summary and append summary row

Folders run sequentially only.

`STOP_ON_FAILURE=1` only changes `all` mode; stop after first failed folder.

## Rerun Policy Contract

Exact policy notes:

```text
FORCE=1 always reruns with clean-first reset
FORCE=2 reruns because prior performance status is FAIL (improve in place)
fresh SUCCESS report already exists
report requires rerun with clean-first reset
```

Rules:

- `FORCE=1`: always rerun, clean-first
- `FORCE=2`: rerun in place only when prior `Performance Status` is `FAIL`
- default or fallback: skip only when prior report has `Operation Status: SUCCESS` and report is fresher than every non-test implementation `.go` file

## Protected File And Scope Contract

During reset and implementation phases, allowed edits are only non-test implementation `.go` files inside target folder.

Protected files include:

- `api_contract.go`
- `go.mod`
- `go.sum`
- `SPECS.md`
- any `*.md`
- any `*_test.go`
- any `*_bench_test.go`
- `helpers_test.go`
- `bench_policy_test.go`

Files outside target folder are out of scope.

During report phase only, target report file is additionally allowed.

If protected or out-of-scope files are modified, runtime must restore snapshot content and fail phase.

Exact warning template for report-phase restoration:

```text
[gen-impl] Warning: report phase {report_attempt} for {folder} edited protected or out-of-scope files; restored snapshot content
```

## Verification Contract

Doc comment audit command:

```text
go run .scripts/gen-impl/helpers/audit_doc_comments.go --folder {folder}
```

Unit test command:

```text
go test -json ./{folder}/...
```

Benchmark command:

```text
make -C {repo_root} bench-folder FOLDER={folder}
```

Benchmark threshold failures do not trigger implementation retries.

## Timeout Marker Contract

Exact timeout lines written into logs:

```text
[gen-impl] command timeout
[gen-impl] spawner timeout
[gen-impl] spawner idle timeout ({idle_timeout_seconds}s without output)
```

## Prompt Template Contract

Prompt template files must remain byte-for-byte equivalent to these bodies after placeholder substitution rules are applied.

### `helpers/prompt_impl.md.tmpl`

```text
You are implementing `@@FOLDER@@` inside repository `@@REPO_ROOT@@`.

Read these files before editing:
1. `AGENTS.md`
2. `STRUCTURE-OVERVIEW.md`
3. `@@FOLDER@@/README.md`
4. `@@FOLDER@@/SPECS.md`
5. `@@FOLDER@@/api_contract.go`
6. All existing Go files in `@@FOLDER@@`
7. All tests and benchmarks in `@@FOLDER@@`
8. Prior implementation report copy at `@@PRIOR_REPORT_PATH@@`
9. Prior failure suggestions at `@@PRIOR_SUGGESTIONS_PATH@@`
10. Previous attempt summary at `@@FAILURE_SUMMARY_PATH@@`

Attempt context:
- Attempt `@@ATTEMPT@@` of `@@MAX_ATTEMPTS@@`
- Run artifacts directory: `@@RUN_DIR@@`

Goal:
- Implement or fix `@@FOLDER@@` so behavior matches the specs, API contract, tests, and benchmark policy.
- Ensure unit tests pass.
- Keep implementation benchmark-ready for the script's benchmark verification phase.

Hard rules:
- Modify only non-test implementation `.go` files inside `@@FOLDER@@`.
- Do not modify tests.
- Do not modify `SPECS.md`.
- Do not modify any markdown docs.
- Do not modify `api_contract.go`.
- Do not modify `go.mod` or `go.sum`.
- Do not modify files outside `@@FOLDER@@`.
- Do not create or edit `IMPLEMENTATION_REPORT.md` in this phase.

Doc comment rules:
- Every exported implemented function must have a proper Go doc comment.
- API interface methods must start with `FuncName implements the API interface.`.
- Non-interface exported functions such as `New` must use normal truthful Go doc comments.
- Every exported implemented function comment must include `Example:`.
- Doc comments must stay consistent with both `SPECS.md` and the implementation itself.

Implementation rules:
- Follow repo-wide forbidden-feature constraints exactly.
- Keep changes minimal and focused.
- Treat tests, benchmarks, markdown docs, and API contracts as fixed input.
- If prior failure suggestions exist, use them as concrete guidance.
- If current implementation was reset to stubs, replace the stubs completely.
- Do not run full benchmark suites in this implementation phase.
- If you need a benchmark sanity check, run only a narrow quick check (for example one benchmark with `-benchtime=1x`).

Before finishing:
- Review your own changes carefully.
- Make sure the implementation is consistent with the tests and benchmarks.
```

### `helpers/prompt_reset.md.tmpl`

```text
You are resetting the implementation state for `@@FOLDER@@` inside repository `@@REPO_ROOT@@`.

Read these files before editing:
1. `AGENTS.md`
2. `STRUCTURE-OVERVIEW.md`
3. `@@FOLDER@@/README.md`
4. `@@FOLDER@@/SPECS.md`
5. `@@FOLDER@@/api_contract.go`
6. All existing Go files in `@@FOLDER@@`
7. All tests and benchmarks in `@@FOLDER@@`
8. Prior report copy at `@@PRIOR_REPORT_PATH@@`
9. Prior failure suggestions at `@@PRIOR_SUGGESTIONS_PATH@@`

Reset-phase goal:
- Forget the current non-test implementation for `@@FOLDER@@`.
- Replace existing implementation bodies with empty stubs or obvious placeholders so the next implementation attempt starts fresh.
- Keep the public API surface and function signatures unchanged.
- Keep doc comments truthful to the current stub behavior.

Hard rules:
- Modify only non-test implementation `.go` files inside `@@FOLDER@@`.
- Do not modify tests.
- Do not modify `SPECS.md`.
- Do not modify any markdown docs.
- Do not modify `api_contract.go`.
- Do not modify `go.mod` or `go.sum`.
- Do not modify files outside `@@FOLDER@@`.
- Do not create or edit `IMPLEMENTATION_REPORT.md` in this reset phase.

Doc comment rules:
- Every exported implemented function must have a proper Go doc comment.
- API interface methods must start with `FuncName implements the API interface.`.
- Non-interface exported functions such as `New` must use normal truthful Go doc comments.
- Every exported implemented function comment must include `Example:`.

Keep this phase minimal. Reset only. Full implementation happens in the next phase.
```

### `helpers/prompt_report.md.tmpl`

```text
You are writing the final implementation report for `@@FOLDER@@` inside repository `@@REPO_ROOT@@`.

Read this input file first:
- `@@REPORT_INPUT_PATH@@`

Write exactly one output file:
- `@@OUTPUT_REPORT_PATH@@`

Hard rules:
- Modify no other file.
- Keep the report factual and concise.
- Use the exact status, folder name, attempt count, and change facts from the input file.

Required report structure:
1. `# IMPLEMENTATION_REPORT`
2. `Operation Status: ...`
3. `Performance Status: ...`
4. `Folder: ...`
5. `Attempts Used: ...`
6. `## Files Changed` with markdown table
7. `## Unit Test Summary` with markdown table exactly shaped as `no. | scenario | pass / fail`
8. `## Benchmark Summary` with markdown table exactly shaped as `no. | scenario | budget-ratio | good/bad | pass / fail`
9. `## Operational Failure Causes` with markdown table exactly shaped as `no. | cause | scenario | evidence | suggestion`
10. Under operational failure causes, `### Improvement Suggestions` with markdown table exactly shaped as `no. | cause | failed scenario | suggestion`
11. `## Benchmark Failure Causes` with markdown table exactly shaped as `no. | cause | scenario | evidence | suggestion`
12. Under benchmark failure causes, `### Improvement Suggestions` with markdown table exactly shaped as `no. | cause | failed scenario | suggestion`

Do not add speculative claims that are not supported by the input file.
```

## Report Contract

Final report path must be exactly:

```text
tmp/gen-impl/reports/{folder}/IMPLEMENTATION_REPORT.md
```

Required exact top-level structure:

```text
# IMPLEMENTATION_REPORT

Operation Status: {SUCCESS|FAILURE}
Performance Status: {PASS|FAIL|NOT_RUN}
Folder: {folder}
Attempts Used: {attempts_used}
```

Required sections and headers:

- `## Files Changed`
- `| File | Change |`
- `## Unit Test Summary`
- `| no. | scenario | pass / fail |`
- `## Benchmark Summary`
- `| no. | scenario | budget-ratio | good/bad | pass / fail |`
- `## Operational Failure Causes`
- `| no. | cause | scenario | evidence | suggestion |`
- `### Improvement Suggestions`
- `| no. | cause | failed scenario | suggestion |`
- `## Benchmark Failure Causes`
- second `### Improvement Suggestions`

Fallback rows:

- unit scenarios absent: `| 1 | (none) | PASS |`
- benchmark scenarios absent after operational success: `| 1 | (none) | N/A | GOOD | PASS |`
- benchmark skipped because operational failure: `| 1 | (not run) | N/A | BAD | FAIL |`
- failure-causes absent: `(none)` rows

`budget-ratio > 1.0` means `BAD` and `FAIL`.

## Summary Table Contract

Exact stdout summary table header:

```text
| Folder | Status | Attempts | Report |
|---|---|---:|---|
```

Exact row template:

```text
| `{folder}` | {status} | {attempts} | `{report_path}` |
```

## User-Facing Text Contract

Exact stdout log templates:

```text
[gen-impl] Processing {folder}
[gen-impl] Done {folder} ({status}, attempts={attempts_used})
[gen-impl] Run artifacts: {run_root}
```

Exact message template when report phase violates scope:

```text
[gen-impl] Warning: report phase {report_attempt} for {folder} edited protected or out-of-scope files; restored snapshot content
```

The program must not invent extra user-facing prose outside the templates in this spec, except for raw error passthrough placeholders explicitly allowed here.

## Live TTY Progress Contract

Live progress is enabled only when `sys.stdout.isatty()` is true.

Terminal width rules:

- detect width with `shutil.get_terminal_size(fallback=(120, 24)).columns`
- if detected width is below `50`, use `50`

Exact spinner sequence:

- `.`
- `..`
- `...`

Exact rendered lines before truncation:

```text
{phase} | {spinner} | {elapsed_seconds:.2f}s
{label}: {tail_text}
```

If tail text is empty after sanitization, exact fallback is:

```text
(waiting output...)
```

Tail sanitization rules, in order:

1. strip ANSI escape sequences
2. replace `\r` with `\n`
3. collapse all whitespace runs to single spaces
4. strip ASCII control characters except `\n` handling above
5. trim surrounding whitespace

Truncation rules:

- if line length fits width, keep it unchanged
- if width is `<= 3`, hard-cut to width
- otherwise truncate to `width - 3` and append `...`

Redraw rules:

- update prints exactly two lines
- when previous dynamic lines exist, clear them first with carriage-return and ANSI clear-line sequence
- `clear()` removes active dynamic lines only when TTY mode was active

## Determinism Contract

These behaviors must remain deterministic for same repository state and same external command outputs:

- folder order comes from root `go.work`
- all folder-local glob results are sorted by Python `Path.glob()` iteration wrapped in `sorted(...)` where implementation depends on file order
- report and summary table ordering must remain stable
- prompt template text must remain exact
- user-facing strings must remain exact templates from this spec

Dynamic placeholders remain variable by input or environment:

- `{folder}`
- `{required_file}`
- `{required_name}`
- `{count}`
- `{target}`
- `{reason}`
- `{exc}`
- `{probe_clean_log}`
- `{report_path}`
- `{run_root}`
- `{elapsed_seconds:.2f}`
- timestamp value matching `%Y%m%d_%H%M%S`

## README Contract

`.scripts/gen-impl/README.md` is human-facing summary only.

It must contain exact line near top:

```text
Canonical contract: `.scripts/gen-impl/SPECS.md`
```

Root `README.md` must only mention that `gen-impl` exists, what it is for, and point to `.scripts/gen-impl/README.md` for full documentation. Root `README.md` must not describe detailed `gen-impl` implementation behavior.
