# assess-implementation resources

Persistent helper assets for the `assess-implementation` skill.

## Files

- `assess_impl.py`: main assessment runner.

## Usage

From repository root:

```sh
python3 .agents/skills/assess-implementation/resources/assess_impl.py
```

Scoped run:

```sh
python3 .agents/skills/assess-implementation/resources/assess_impl.py --folders list-array stack
```

Skip command execution:

```sh
python3 .agents/skills/assess-implementation/resources/assess_impl.py --skip-commands
```

## Command policy

- Strict-hard mode uses per-folder commands in both whole-repo and scoped runs:
  - `make test-folder FOLDER=<folder>`
  - `make bench-folder FOLDER=<folder>` when benchmark funcs exist
- Fallback to direct `go test` commands only when make infrastructure is unavailable (missing target/tooling), not when tests/benchmarks fail.
- Command evidence is required for implementation folders. Missing or failed command evidence applies strict score caps.

## Strict-hard caps

- Test command failed or missing -> cap `Test evidence` to `4/20` and `Invariants and behavior` to `6/20`.
- Benchmark command failed or missing (when benchmarks exist) -> cap `Benchmark and delivery evidence` to `5/10`.
- Existing hard-fail rules still force final grade `F`.

## Output

- Prints run metadata and `## Summary Table` markdown block to terminal.
- Terminal output must include table header and data rows (not prose-only summary).
- Writes full report to `tmp/assessment_YYYYMMDD_HHMMSS.md`.

## Important

- Reuse this helper. Do not generate ad-hoc Python scripts for routine assessment/report output.
- Only create throwaway debug scripts for parser investigation, then delete them immediately.
