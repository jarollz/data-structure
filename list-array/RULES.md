# RULES.md - list-array

## Goal
Implement a generic array-backed list with index operations.

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics
- [ ] Use `T any` for element type.

## Required API
- [ ] `New(capacity int) *List[T]`
- [ ] `Append(v T) bool`
- [ ] `Get(i int) (T, bool)`
- [ ] `Set(i int, v T) bool`
- [ ] `Insert(i int, v T) bool`
- [ ] `Delete(i int) (T, bool)`
- [ ] `Len() int`
- [ ] `Cap() int`
- [ ] `Clear()`
- [ ] `Values() iter.Seq[T]`

## Internal representation
- [ ] Fixed-capacity array plus length counter.
- [ ] If growth is implemented, allocate new larger array and copy manually.
- [ ] Shift elements for middle insert/delete.

## Auto-resize policy
- [ ] Grow when `Len() == Cap()`.
- [ ] New capacity on grow: `2x` when `Cap() < 1024`, otherwise `Cap() + Cap()/2`.
- [ ] Shrink when `Len() <= Cap()/4` and `Cap() > minCap`.
- [ ] New capacity on shrink: `max(minCap, Cap()/2, 2*Len())`.
- [ ] `minCap` is `max(16, initial capacity)`.
- [ ] Use hysteresis; do not resize on every delete.

## Invariants
- [ ] Live elements are in index range `[0, Len())`.
- [ ] `Len()` never exceeds capacity.
- [ ] Relative order is preserved after insert/delete operations.

## Iterator contract
- [ ] `Values()` yields elements from index `0` to `Len()-1`.
- [ ] Each live element is yielded exactly once.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty list yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Negative index and `i >= Len()` are handled safely.
- [ ] Insert at head and tail works.
- [ ] Delete head and tail works.
- [ ] Clear resets length state.

## Test checklist
- [ ] Index boundary tests.
- [ ] Insert/delete shift correctness tests.
- [ ] Random operation tests against reference sequence.

## Benchmark checklist
- [ ] Append benchmark.
- [ ] Middle insert/delete benchmark.
- [ ] Iteration benchmark.

## Test Generator Hints
- Emphasize index boundaries, head/middle/tail edits, and order preservation.
- Compare sequence state with oracle after randomized operations.
- Validate `Len/Cap` consistency during growth/clear paths.
- Iterator tests must verify index-order traversal and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for array-list API including bounds checks, insert/delete shifts, and clear behavior."
- Property tests: "Generate randomized list operations with fixed seed and compare final/stepwise state to oracle sequence."
- Iterator tests: "Generate tests for `Values() iter.Seq[T]` verifying index order, exact count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate list-array benchmarks for append-heavy, middle-edit-heavy, and iteration workloads at 1e3/1e4/1e5."
