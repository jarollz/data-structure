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
- [ ] `New(capacity int) *List[T]` creates an empty list. Normalize `capacity <= 0` to `16`. Effective starting capacity is `max(16, capacity)`.
- [ ] `Append(v T) bool` appends at tail and always returns `true`.
- [ ] `Get(i int) (T, bool)`
- [ ] `Set(i int, v T) bool` updates existing element and returns `false` when `i` is outside `[0, Len())`.
- [ ] `Insert(i int, v T) bool` inserts before index `i`. Valid indexes are `[0, Len()]`. Return `false` when `i` is outside that range.
- [ ] `Delete(i int) (T, bool)` removes and returns element at index `i`. Return `(zero, false)` when `i` is outside `[0, Len())`.
- [ ] `Len() int`
- [ ] `Cap() int`
- [ ] `Clear()`
- [ ] `Values() iter.Seq[T]`

## Internal representation
- [ ] Contiguous backing storage plus explicit length and capacity fields.
- [ ] Backing storage always keeps live elements in indexes `[0, Len())` with no gaps.
- [ ] On grow or shrink, allocate new backing storage and copy live elements manually in index order.
- [ ] Shift elements for middle insert/delete.

## Auto-resize policy
- [ ] `Cap()` immediately after `New` equals normalized starting capacity.
- [ ] Grow when `Len() == Cap()`.
- [ ] New capacity on grow: `2x` when `Cap() < 1024`, otherwise `Cap() + Cap()/2`.
- [ ] Shrink when `Len() <= Cap()/4` and `Cap() > minCap`.
- [ ] New capacity on shrink: `max(minCap, Cap()/2, 2*Len())`.
- [ ] `minCap` is normalized starting capacity.
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
- [ ] Insert at `i == Len()` appends.
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
