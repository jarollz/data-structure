# SPECS.md - heap

## Goal
Implement a generic binary heap using arrays only.

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics and ordering
- [ ] Use `T any`.
- [ ] Require comparator to define min-heap or max-heap behavior.

## Required API
Each public API below is part of contract. Use each checklist as implementation target, not loose guidance.

### `New(capacity int, cmp func(a, b T) int) *Heap[T]`
Purpose
- [ ] Create empty heap with normalized starting capacity and comparator-defined ordering.

Behavior expectations
- [ ] Normalize `capacity <= 0` to `16`.
- [ ] Effective starting capacity is `max(16, capacity)`.
- [ ] Returned heap is non-nil and empty.
- [ ] `Len()` is `0` immediately after construction.
- [ ] `Cap()` is normalized starting capacity immediately after construction.
- [ ] Comparator becomes persistent ordering rule for all future operations.

Performance expectations
- [ ] `O(1)` time.
- [ ] `O(capacity)` backing storage.

### `Push(v T) bool`
Purpose
- [ ] Insert value and restore heap order.

Behavior expectations
- [ ] New value becomes part of occupied heap prefix.
- [ ] Heap property holds after insertion finishes.
- [ ] Full heap grows before write and still succeeds.
- [ ] Method always returns `true`.

Performance expectations
- [ ] `O(log n)` time after any needed resize.
- [ ] Worst-case single call may include `O(n)` copy during resize plus sift-up work.

Implementation checklist
- [ ] Insert at end of occupied prefix first.
- [ ] Use sift-up to restore comparator order.

### `PopTop() (T, bool)`
Purpose
- [ ] Remove and return current heap top value.

Behavior expectations
- [ ] Non-empty heap returns comparator-best value currently stored.
- [ ] Empty heap returns `(zero, false)` and does not panic.
- [ ] Last live element moves to root before sift-down when needed.
- [ ] Heap property still holds after removal and any shrink.

Performance expectations
- [ ] `O(log n)` time after any needed resize.
- [ ] Worst-case single call may include `O(n)` copy during shrink plus sift-down work.

### `PeekTop() (T, bool)`
Purpose
- [ ] Read current heap top without removing it.

Behavior expectations
- [ ] Non-empty heap returns comparator-best value and `true`.
- [ ] Empty heap returns `(zero, false)` and does not panic.
- [ ] `PeekTop()` never changes heap shape, order, length, or capacity.

Performance expectations
- [ ] `O(1)` time.

### `Len() int`
Purpose
- [ ] Report number of live heap elements.

Behavior expectations
- [ ] Count increases after `Push`.
- [ ] Count decreases after successful `PopTop`.
- [ ] Count becomes `0` after `Clear`.

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
- [ ] Reset heap to empty logical state.

Behavior expectations
- [ ] `Len()` becomes `0`.
- [ ] Later `PopTop()` and `PeekTop()` behave like empty heap operations.
- [ ] Future `Push()` still works correctly.
- [ ] Clear on empty heap is safe.

Performance expectations
- [ ] Target `O(1)` logical reset.
- [ ] `O(n)` slot reset is allowed only if implementation needs it.

### `Clone() *Heap[T]`
Purpose
- [ ] Create independent heap container with same comparator and same internal heap state.

Behavior expectations
- [ ] Clone preserves `Len()`, `Cap()`, comparator, and internal heap-array order.
- [ ] Elements are copied with normal Go assignment.
- [ ] Clone remains heap-valid and container-independent from source.
- [ ] Empty source produces empty independent clone with same `Cap()` and comparator.

Performance expectations
- [ ] `O(n)` time for `n = Len()`.
- [ ] `O(n)` extra storage.

### `CloneWith(cloneValue func(T) T) *Heap[T]`
Purpose
- [ ] Create independent heap container while optionally transforming each live element.

Behavior expectations
- [ ] Preserves `Len()`, `Cap()`, comparator, and internal heap-array order.
- [ ] Nil hook is equivalent to `Clone()`.
- [ ] Non-nil hook is called exactly once per live element.
- [ ] Hook call order is current internal array order, not sorted order.
- [ ] Hook is never called for unused capacity slots.

Performance expectations
- [ ] `O(n)` container work plus hook cost.

### `Values() iter.Seq[T]`
Purpose
- [ ] Iterate live heap contents in documented internal-array order.

Behavior expectations
- [ ] Yield each live element exactly once.
- [ ] Yield order is internal array order, not sorted order.
- [ ] Empty heap yields nothing.
- [ ] Early stop works when consumer returns `false`.
- [ ] Mutation during iteration is not safe.

Performance expectations
- [ ] Full traversal is `O(n)`.
- [ ] Early-stop traversal is `O(k)` for yielded prefix size `k`.
- [ ] Iterator setup is `O(1)`.

## Internal representation
- [ ] Array-backed complete tree by index.
- [ ] Parent/child index formulas are: parent `(i-1)/2`, left child `2*i+1`, right child `2*i+2`.
- [ ] Use sift-up and sift-down to restore heap property.

## Auto-resize policy
- [ ] `Cap()` immediately after `New` equals normalized starting capacity.
- [ ] Grow when `Len() == Cap()`.
- [ ] New capacity on grow: `2x` when `Cap() < 1024`, otherwise `Cap() + Cap()/2`.
- [ ] Shrink when `Len() <= Cap()/4` and `Cap() > minCap`.
- [ ] New capacity on shrink: `max(minCap, Cap()/2, 2*Len())`.
- [ ] `minCap` is normalized starting capacity.
- [ ] Resize only backing array; heap-order invariants must remain valid.
- [ ] Use hysteresis; do not resize on every `PopTop`.

## Invariants
- [ ] Complete tree shape in occupied prefix.
- [ ] Heap property holds for all parent-child pairs.
- [ ] `Len()` equals number of stored elements.
- [ ] `Clone()` and `CloneWith(...)` preserve `Len()`, `Cap()`, comparator, and internal heap-array order in clone.

## Iterator contract
- [ ] `Values()` yields each stored element exactly once.
- [ ] Yield order is internal array order, not sorted order.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty heap yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Push into full capacity grows backing storage and still returns `true`.
- [ ] Pop/peek on empty heap returns `(zero, false)`.
- [ ] Repeated equal-priority values handled correctly.
- [ ] `Clone()` on empty heap returns empty independent heap with same `Cap()` and comparator.
- [ ] `CloneWith(nil)` is equivalent to `Clone()`.
- [ ] `CloneWith(...)` calls `cloneValue` only for live heap elements in internal array order.

## Test checklist
- [ ] Push/pop correctness tests.
- [ ] Heap-property check after each mutation batch.
- [ ] Randomized compare against sorted reference behavior.
- [ ] `Clone/CloneWith` tests for independence, nil-hook equivalence, shallow-copy default, custom hook behavior, preserved capacity, preserved internal array order, and live-element-only hook calls.

## Benchmark checklist
- [ ] Benchmark only non-trivial public APIs.
- [ ] Do not benchmark trivial metadata accessors such as `Len()` and `Cap()`.
- [ ] `Push` benchmark.
- [ ] `PopTop` benchmark.
- [ ] `PeekTop` benchmark.
- [ ] `Clear` benchmark.
- [ ] `Clone` benchmark.
- [ ] `CloneWith` benchmark.
- [ ] `Values` benchmark.
- [ ] Mixed priority-queue workload benchmark.
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
- [ ] Expected `O(log n)` APIs must keep `ns/op(1e5) <= 2.5 * ns/op(1e3)`.
- [ ] Expected `O(n)` APIs must compare normalized `(ns/op)/n` and keep normalized `1e5` result within `3.0x` of normalized `1e3` result.

### Layer 2: Absolute timing thresholds
- [ ] Absolute thresholds are per API, per payload class, and per benchmark size.
- [ ] Provisional thresholds should be derived from expected complexity and payload size.
- [ ] Calibrated thresholds should be derived from repeated local runs, median `ns/op`, and safety factor.
- [ ] Suggested safety factor for stable simple heap ops is `1.75x`.
- [ ] Suggested safety factor for iteration and mixed workloads is `2.00x`.
- [ ] Suggested safety factor for copy-heavy clone or large-payload workloads is `2.25x`.

### Payload classes
- [ ] Tiny payload benchmark uses scalar-like elements such as `int`.
- [ ] Large payload benchmark uses fixed-size struct or fixed-size array payload.
- [ ] Large-payload thresholds must be tracked independently from tiny-payload thresholds.

## Test Generator Hints
- Verify heap property after every mutating operation sequence.
- Include duplicate priorities and comparator edge cases.
- Validate `Clone/CloneWith` preserve `Cap()`, internal array order, and hook-call count without breaking heap property.
- Check repeated `PopTop` yields monotonic priority order.
- Iterator tests must validate count, internal-order expectation, and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for heap API (`Push`, `PopTop`, `PeekTop`) including empty and full-capacity behavior."
- Clone tests: "Generate tests for `Clone()` and `CloneWith(...)` verifying independent container state, shallow-copy default, nil-hook equivalence, custom hook behavior, preserved capacity, preserved internal array order, and live-element-only hook calls."
- Property tests: "Generate randomized push/pop sequences with fixed seed and assert heap property after each batch."
- Iterator tests: "Generate tests for `Values() iter.Seq[T]` ensuring exact count, early stop, and documented non-sorted internal order."
- Benchmarks: "Generate heap benchmarks for push-only, pop-only, and mixed workloads at 1e3/1e4/1e5 sizes."
