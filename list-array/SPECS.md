# SPECS.md - list-array

## Goal
Implement a generic array-backed list with index operations.

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics
- [ ] Use `T any` for element type.

## Required API
Each public API below is part of contract. Use each checklist as implementation target, not loose guidance.

### `New(capacity int) *List[T]`
Purpose
- [ ] Create empty array-backed list with normalized starting capacity.

Behavior expectations
- [ ] Normalize `capacity <= 0` to `16`.
- [ ] Effective starting capacity is `startCap = max(16, capacity)`.
- [ ] Returned list is non-nil and empty.
- [ ] `Len()` is `0` immediately after construction.
- [ ] `Cap()` is normalized starting capacity immediately after construction.

Performance expectations
- [ ] `O(startCap)` time.
- [ ] `O(startCap)` backing storage.

### `Append(v T) bool`
Purpose
- [ ] Add value at logical tail.

Behavior expectations
- [ ] New value appears at index `Len()-1` after append completes.
- [ ] Existing element order is unchanged.
- [ ] Full list grows before write and still succeeds.
- [ ] Method always returns `true`.

Performance expectations
- [ ] Amortized `O(1)` time.
- [ ] Worst-case single call may be `O(n)` during resize copy.

### `Get(i int) (T, bool)`
Purpose
- [ ] Read value at index without mutating list.

Behavior expectations
- [ ] Valid indexes are `[0, Len())`.
- [ ] Invalid indexes return `(zero, false)`.
- [ ] `Get` never changes list contents, length, or capacity.

Performance expectations
- [ ] `O(1)` time.

### `Set(i int, v T) bool`
Purpose
- [ ] Replace existing live element at index.

Behavior expectations
- [ ] Valid indexes are `[0, Len())`.
- [ ] Invalid indexes return `false` and leave list unchanged.
- [ ] Replacement does not change `Len()` or `Cap()`.
- [ ] Element order before and after target index is unchanged.

Performance expectations
- [ ] `O(1)` time.

### `Insert(i int, v T) bool`
Purpose
- [ ] Insert value before index `i`.

Behavior expectations
- [ ] Valid indexes are `[0, Len()]`.
- [ ] `i == 0` inserts at head.
- [ ] `i == Len()` appends at tail.
- [ ] Invalid indexes return `false` and leave list unchanged.
- [ ] Existing elements at and after `i` shift one slot right.
- [ ] Relative order of existing elements is preserved.

Performance expectations
- [ ] `O(n)` time because middle and head inserts may shift elements.
- [ ] Worst-case single call may include resize copy plus shift work.

### `Delete(i int) (T, bool)`
Purpose
- [ ] Remove and return element at index `i`.

Behavior expectations
- [ ] Valid indexes are `[0, Len())`.
- [ ] Invalid indexes return `(zero, false)`.
- [ ] Elements after deleted index shift one slot left.
- [ ] Relative order of remaining elements is preserved.
- [ ] Delete of head and tail both work correctly.

Performance expectations
- [ ] `O(n)` time because middle and head deletes may shift elements.
- [ ] Worst-case single call may include shrink copy plus shift work.

### `Len() int`
Purpose
- [ ] Report number of live elements.

Behavior expectations
- [ ] Count increases after append or successful insert.
- [ ] Count decreases after successful delete.
- [ ] Count becomes `0` after clear.

Performance expectations
- [ ] `O(1)` time.

### `Cap() int`
Purpose
- [ ] Report current backing-array capacity.

Behavior expectations
- [ ] Capacity is storage size, not live element count.
- [ ] Capacity never drops below normalized starting capacity.
- [ ] Capacity reflects grow and shrink operations.

Performance expectations
- [ ] `O(1)` time.

### `Clear()`
Purpose
- [ ] Reset list to empty logical state.

Behavior expectations
- [ ] `Len()` becomes `0`.
- [ ] Future append, insert, get, set, and delete operations still work correctly.
- [ ] Clear on empty list is safe.

Performance expectations
- [ ] Target `O(1)` logical reset.
- [ ] `O(n)` slot reset is allowed only if implementation needs it.

### `Clone() *List[T]`
Purpose
- [ ] Create independent list container with same visible contents and capacity.

Behavior expectations
- [ ] Clone preserves `Len()`, `Cap()`, and index order.
- [ ] Elements are copied with normal Go assignment.
- [ ] Clone is container-independent from source.
- [ ] Empty source produces empty independent clone with same `Cap()`.

Performance expectations
- [ ] `O(Cap())` time.
- [ ] `O(Cap())` extra storage.

### `CloneWith(cloneValue func(T) T) *List[T]`
Purpose
- [ ] Create independent list container while optionally transforming each live element.

Behavior expectations
- [ ] Preserves `Len()`, `Cap()`, and index order.
- [ ] Nil hook is equivalent to `Clone()`.
- [ ] Non-nil hook is called exactly once per live element.
- [ ] Hook call order is ascending index order.
- [ ] Hook is never called for unused capacity slots.

Performance expectations
- [ ] `O(Cap())` container work plus hook cost.

### `Values() iter.Seq[T]`
Purpose
- [ ] Iterate live elements in ascending index order.

Behavior expectations
- [ ] Yield indexes `0` through `Len()-1` in order.
- [ ] Yield each live element exactly once.
- [ ] Empty list yields nothing.
- [ ] Early stop works when consumer returns `false`.
- [ ] Mutation during iteration is not safe.

Performance expectations
- [ ] Full traversal is `O(n)`.
- [ ] Early-stop traversal is `O(k)` for yielded prefix size `k`.
- [ ] Iterator setup is `O(1)`.

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
- [ ] `Clone()` and `CloneWith(...)` preserve `Len()`, `Cap()`, and index order in clone.

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
- [ ] `Clone()` on empty list returns empty independent list with same `Cap()`.
- [ ] `CloneWith(nil)` is equivalent to `Clone()`.
- [ ] `CloneWith(...)` calls `cloneValue` only for live elements in index order.

## Test checklist
- [ ] Index boundary tests.
- [ ] Insert/delete shift correctness tests.
- [ ] Random operation tests against reference sequence.
- [ ] `Clone/CloneWith` tests for independence, nil-hook equivalence, shallow-copy default, custom hook behavior, and live-element-only hook calls.

## Benchmark checklist
- [ ] Benchmark only non-trivial public APIs.
- [ ] Do not benchmark trivial metadata accessors such as `Len()` and `Cap()`.
- [ ] `Append` benchmark.
- [ ] `Get` benchmark.
- [ ] `Set` benchmark.
- [ ] `Insert` benchmark.
- [ ] `Delete` benchmark.
- [ ] `Clear` benchmark.
- [ ] `Clone` benchmark.
- [ ] `CloneWith` benchmark.
- [ ] `Values` benchmark.
- [ ] Run benchmarks at `1e3`, `1e4`, and `1e5` sizes.
- [ ] Use tiny payload benchmarks for structure overhead.
- [ ] Use large payload benchmarks for copy-sensitive APIs.

## Benchmark Validation Policy
- [ ] Use two validation layers for every benchmarked non-trivial API.
- [ ] Layer 1 validates complexity growth.
- [ ] Layer 2 validates absolute timing budget.
- [ ] Before valid implementation exists, absolute thresholds are provisional engineering targets, not calibrated facts.
- [ ] Implementer should aim to get as close as practical to provisional targets while keeping behavior correct.
- [ ] After first correct implementation exists, rerun calibration on target machine and replace provisional targets with measured median-based thresholds.

### Layer 1: Complexity-growth thresholds
- [ ] Expected `O(1)` APIs must keep `ns/op(1e5) <= 3.0 * ns/op(1e3)`.
- [ ] Expected `O(n)` APIs must compare normalized `(ns/op)/n` and keep normalized `1e5` result within `3.0x` of normalized `1e3` result.

### Layer 2: Absolute timing thresholds
- [ ] Absolute thresholds are per API, per payload class, and per benchmark size.
- [ ] Provisional thresholds should be derived from expected complexity and payload size.
- [ ] Calibrated thresholds should be derived from repeated local runs, median `ns/op`, and safety factor.
- [ ] Suggested safety factor for stable simple list ops is `1.50x`.
- [ ] Suggested safety factor for traversal and mixed edit workloads is `2.00x`.
- [ ] Suggested safety factor for copy-heavy clone or large-payload workloads is `2.25x`.

### Payload classes
- [ ] Tiny payload benchmark uses scalar-like elements such as `int`.
- [ ] Large payload benchmark uses fixed-size struct or fixed-size array payload.
- [ ] Large-payload thresholds must be tracked independently from tiny-payload thresholds.

## Test Generator Hints
- Emphasize index boundaries, head/middle/tail edits, and order preservation.
- Compare sequence state with oracle after randomized operations.
- Validate `Len/Cap` consistency during growth/clear paths.
- Validate `Clone/CloneWith` preserve `Len/Cap`, order, and hook-call count.
- Iterator tests must verify index-order traversal and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for array-list API including bounds checks, insert/delete shifts, and clear behavior."
- Clone tests: "Generate tests for `Clone()` and `CloneWith(...)` verifying independent container state, shallow-copy default, nil-hook equivalence, custom hook behavior, and live-element-only hook calls."
- Property tests: "Generate randomized list operations with fixed seed and compare final/stepwise state to oracle sequence."
- Iterator tests: "Generate tests for `Values() iter.Seq[T]` verifying index order, exact count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate list-array benchmarks for append-heavy, middle-edit-heavy, and iteration workloads at 1e3/1e4/1e5."
