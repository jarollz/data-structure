# RULES.md - list-skip

## Goal
Implement a generic ordered skip list.

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics and ordering
- [ ] Use `T any`.
- [ ] Require comparator for ordering.
- [ ] Duplicate policy: do not store duplicates. `Insert(v)` returns `false` when equal value already exists.

## Required API
- [ ] `New(maxLevel int, cmp func(a, b T) int) *SkipList[T]` creates an empty skip list. Normalize `maxLevel < 1` to `1`.
- [ ] `Insert(v T) bool` inserts new value and returns `true`; return `false` on duplicate.
- [ ] `Delete(v T) bool` removes existing value and returns `true`; return `false` when value is missing.
- [ ] `Has(v T) bool`
- [ ] `Len() int`
- [ ] `Clear()`
- [ ] `Values() iter.Seq[T]`

## Internal representation
- [ ] Use an index-based node pool with `-1` as nil/sentinel index.
- [ ] Keep head sentinel node with forward references for all levels `0..maxLevel-1`.
- [ ] Node stores value and forward references for each level.
- [ ] Level references use arrays/indexes only.
- [ ] Track `currentLevel`, length, free-list head, and deterministic RNG state.
- [ ] Level generation policy is fixed: initialize internal xorshift32 state to `1` in `New`; on each insert, advance state and promote while low bit is `1`, stopping at first `0` or `maxLevel`. Same operation sequence must produce same levels.

## Auto-resize policy
- [ ] No capacity-based `Grow()` or `Shrink()` API.
- [ ] When free-list is empty, allocate larger node-pool arrays and copy node/link fields by index.
- [ ] Allocate node storage on insert operations.
- [ ] Reclaim node storage on delete/clear operations.
- [ ] Optional free-list reuse is allowed if skip-list invariants remain valid.
- [ ] After deletions, reduce `currentLevel` while highest level has no live nodes.

## Invariants
- [ ] Level 0 contains all live elements in sorted order.
- [ ] Higher-level links preserve sorted order.
- [ ] Search path from top level reaches correct target/miss.
- [ ] `Len()` equals number of live elements.

## Iterator contract
- [ ] `Values()` yields sorted order using level 0 traversal.
- [ ] Each element is yielded exactly once.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty list yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Empty list operations are safe.
- [ ] Duplicate insert behavior is consistent with contract.
- [ ] Deleting missing value is safe.

## Test checklist
- [ ] Sorted-order verification tests.
- [ ] Random insert/delete/has tests.
- [ ] Level-link consistency checks.

## Benchmark checklist
- [ ] Search benchmark.
- [ ] Insert benchmark.
- [ ] Delete benchmark.
- [ ] Mixed ordered workload benchmark.

## Test Generator Hints
- Validate sorted order on level 0 after every mutation batch.
- Cover duplicate-policy behavior and delete-missing behavior.
- Use randomized insert/delete/has with fixed seed and sorted oracle.
- Iterator tests must verify sorted output and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for skip list API covering insert/delete/has semantics and duplicate policy."
- Property tests: "Generate randomized ordered operations with fixed seed and compare with sorted oracle while checking skip-list invariants."
- Iterator tests: "Generate tests for `Values() iter.Seq[T]` enforcing sorted order, exact count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate skip-list benchmarks for search, insert, delete, and mixed workloads at 1e3/1e4/1e5 sizes."
