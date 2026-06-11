# AGENTS.md

## Scope
This file defines strict rules for AI/code agents working in this repository.

## Non-negotiable rules
- MUST implement data structures using only simple Go types and arrays.
- MUST NOT use `slice` in implementation code.
- MUST NOT use `map` in implementation code.
- MUST keep implementation and tests separated.
- MUST keep iterator contracts aligned with folder `RULES.md`.
- MUST treat mutation during iteration as not safe.

## Allowed exception for tests
- Tests MAY use built-in `map` and slices as oracle/reference models.
- Tests MUST NOT rely on undefined behavior from implementation internals.

## Iterator naming standard
- Linear structures (`list-*`, `queue`, `stack`, `heap`): `Values() iter.Seq[T]`
- Map structures (`map-*`): `All() iter.Seq2[K, V]`
- Set-like trees (`tree-avl`, `tree-red-black`): `InOrder() iter.Seq[T]`
- General tree (`tree-general`): `PreOrder() iter.Seq[T]`

## Required delivery for each structure folder
- `README.md` (human explanation)
- `RULES.md` (implementation contract)
- Tests validating API, invariants, and iterator contract
- Benchmarks with at least read-heavy, write-heavy, and mixed workloads (where applicable)

## Quality gates before marking task done
- Implementation follows forbidden-feature rules (`slice`, `map`).
- API matches folder `RULES.md` required contract.
- Invariants are validated by tests.
- Iterator tests include order, count, early-stop, and mutation-unsafety note.
- Benchmark suite runs for multiple sizes (suggested: `1e3`, `1e4`, `1e5`).

## Change discipline
- Do not weaken constraints in `RULES.md` unless explicitly requested.
- If behavior is ambiguous, prefer stricter contracts and document assumptions.
- Keep language plain and direct in all generated docs.

## Tagging and versioning rules
- MUST follow `RELEASING.md` for all release and tag decisions.
- MUST use `<folder>/vX.Y.Z` for structure folder releases.
- MUST treat root tag `vX.Y.Z` as invalid unless root `go.mod` exists.
- MUST keep published release tags immutable (no move, rewrite, or force-update).
- MUST publish a new version when fixing a bad release tag.
- MUST ensure tag major version matches module path major in `go.mod` (`/vN` for `v2+`).
