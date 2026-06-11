# SPECS.md - queue

## Goal
Implement a generic FIFO queue.

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics
- [ ] Use `T any`.

## Required API
Each public API below is part of contract. Use each checklist as implementation target, not loose guidance.

### `New(capacity int) *Queue[T]`
Purpose
- [ ] Create empty FIFO queue with normalized starting capacity.

Behavior expectations
- [ ] Normalize `capacity <= 0` to `16`.
- [ ] Effective starting capacity is `max(16, capacity)`.
- [ ] Returned queue is non-nil and empty.
- [ ] `Len()` is `0` immediately after construction.
- [ ] `Cap()` is normalized starting capacity immediately after construction.

Performance expectations
- [ ] `O(1)` time.
- [ ] `O(capacity)` backing storage.

Implementation checklist
- [ ] Keep normalized start capacity as `minCap` for shrink policy.
- [ ] Initialize circular-buffer state so future wrap-around is correct.

### `Enqueue(v T) bool`
Purpose
- [ ] Add value at logical back of queue.

Behavior expectations
- [ ] New value becomes last value yielded by queue iteration and dequeue order.
- [ ] Full queue grows before write and still succeeds.
- [ ] Method always returns `true`.
- [ ] Wrap-around writes preserve FIFO order.

Performance expectations
- [ ] Amortized `O(1)` time.
- [ ] Worst-case single call may be `O(n)` during resize and re-pack.

Implementation checklist
- [ ] Write at logical tail position.
- [ ] Preserve front-to-back order when re-packing on resize.

### `Dequeue() (T, bool)`
Purpose
- [ ] Remove and return logical front value.

Behavior expectations
- [ ] Non-empty queue returns oldest enqueued value not yet dequeued.
- [ ] Empty queue returns `(zero, false)` and does not panic.
- [ ] Single-element dequeue leaves queue empty.
- [ ] Wrap-around and shrink paths preserve FIFO order.

Performance expectations
- [ ] Amortized `O(1)` time.
- [ ] Worst-case single call may be `O(n)` during shrink and re-pack.

### `PeekFront() (T, bool)`
Purpose
- [ ] Read logical front value without removing it.

Behavior expectations
- [ ] Empty queue returns `(zero, false)` and does not panic.
- [ ] Non-empty queue returns same value next `Dequeue()` would remove.
- [ ] `PeekFront()` never changes `Len()`, `Cap()`, or order.

Performance expectations
- [ ] `O(1)` time.

### `Len() int`
Purpose
- [ ] Report number of queued live elements.

Behavior expectations
- [ ] Count increases after enqueue.
- [ ] Count decreases after successful dequeue.
- [ ] Count becomes `0` after clear.

Performance expectations
- [ ] `O(1)` time.

### `Cap() int`
Purpose
- [ ] Report current backing-storage capacity.

Behavior expectations
- [ ] Capacity is storage size, not element count.
- [ ] Capacity never drops below normalized starting capacity.
- [ ] Capacity reflects grow and shrink operations.

Performance expectations
- [ ] `O(1)` time.

### `Clear()`
Purpose
- [ ] Reset queue to empty logical state.

Behavior expectations
- [ ] `Len()` becomes `0`.
- [ ] Later `Dequeue()` and `PeekFront()` behave like empty queue operations.
- [ ] Future `Enqueue()` still works from clean state.
- [ ] Clear on empty queue is safe.

Performance expectations
- [ ] Target `O(1)` logical reset.
- [ ] `O(n)` slot reset is allowed only if implementation needs it.

### `Clone() *Queue[T]`
Purpose
- [ ] Create independent queue container with same visible contents.

Behavior expectations
- [ ] Clone preserves `Len()`, `Cap()`, and front-to-back order.
- [ ] Elements are copied with normal Go assignment.
- [ ] Clone is container-independent from source.
- [ ] Empty source produces empty independent clone with same `Cap()`.

Performance expectations
- [ ] `O(n)` time for `n = Len()`.
- [ ] `O(n)` extra storage.

Implementation checklist
- [ ] Copy only live queued elements.
- [ ] Preserve logical front-to-back order even if source is wrapped.

### `CloneWith(cloneValue func(T) T) *Queue[T]`
Purpose
- [ ] Create independent queue container while optionally transforming each live element.

Behavior expectations
- [ ] Preserves `Len()`, `Cap()`, and front-to-back order.
- [ ] Nil hook is equivalent to `Clone()`.
- [ ] Non-nil hook is called exactly once per live queued element.
- [ ] Hook call order is front to back.
- [ ] Hook is never called for unused slots outside logical queue contents.

Performance expectations
- [ ] `O(n)` container work plus hook cost.

### `Values() iter.Seq[T]`
Purpose
- [ ] Iterate live queue values from front to back.

Behavior expectations
- [ ] Yield logical front first and logical back last.
- [ ] Yield each live element exactly once.
- [ ] Empty queue yields nothing.
- [ ] Early stop works when consumer returns `false`.
- [ ] Mutation during iteration is not safe.

Performance expectations
- [ ] Full traversal is `O(n)`.
- [ ] Early-stop traversal is `O(k)` for yielded prefix size `k`.
- [ ] Iterator setup is `O(1)`.

## Internal representation
- [ ] Recommended: circular buffer using array.
- [ ] Track head index and size. Tail may be stored explicitly or derived as `(head + size) % capacity`.
- [ ] Wrap indices correctly.

## Auto-resize policy
- [ ] `Cap()` immediately after `New` equals normalized starting capacity.
- [ ] Grow when `Len() == Cap()`.
- [ ] New capacity on grow: `2x` when `Cap() < 1024`, otherwise `Cap() + Cap()/2`.
- [ ] Shrink when `Len() <= Cap()/4` and `Cap() > minCap`.
- [ ] New capacity on shrink: `max(minCap, Cap()/2, 2*Len())`.
- [ ] `minCap` is normalized starting capacity.
- [ ] On resize, re-pack logical order front-to-back into new array.
- [ ] Use hysteresis; do not resize on every dequeue.

## Invariants
- [ ] FIFO order is preserved.
- [ ] `Len()` equals number of queued elements.
- [ ] Head and tail positions are valid modulo capacity.
- [ ] `Clone()` and `CloneWith(...)` preserve `Len()`, `Cap()`, and front-to-back order in clone.

## Iterator contract
- [ ] `Values()` yields front to back.
- [ ] Each live element is yielded exactly once.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty queue yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Dequeue/peek on empty returns `(zero, false)`.
- [ ] Full-capacity behavior: grow first, then write. `Enqueue` does not fail for capacity reasons.
- [ ] Wrap-around enqueue/dequeue remains correct.
- [ ] `Clone()` on empty queue returns empty independent queue with same `Cap()`.
- [ ] `CloneWith(nil)` is equivalent to `Clone()`.
- [ ] `CloneWith(...)` calls `cloneValue` only for live queued elements from front to back.

## Test checklist
- [ ] FIFO order tests.
- [ ] Wrap-around index tests.
- [ ] Random operation sequence tests.
- [ ] `Clone/CloneWith` tests for independence, nil-hook equivalence, shallow-copy default, custom hook behavior, preserved capacity, and live-element-only hook calls.

## Benchmark checklist
- [ ] Benchmark only non-trivial public APIs.
- [ ] Do not benchmark trivial metadata accessors such as `Len()` and `Cap()`.
- [ ] `Enqueue` benchmark.
- [ ] `Dequeue` benchmark.
- [ ] `PeekFront` benchmark.
- [ ] `Clear` benchmark.
- [ ] `Clone` benchmark.
- [ ] `CloneWith` benchmark.
- [ ] `Values` benchmark.
- [ ] Mixed enqueue/dequeue benchmark.
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
- [ ] Mixed workloads are judged by dominant documented complexity of included operations.

### Layer 2: Absolute timing thresholds
- [ ] Absolute thresholds are per API, per payload class, and per benchmark size.
- [ ] Provisional thresholds should be derived from expected complexity and payload size.
- [ ] Calibrated thresholds should be derived from repeated local runs, median `ns/op`, and safety factor.
- [ ] Suggested safety factor for stable simple queue ops is `1.50x`.
- [ ] Suggested safety factor for iteration and mixed workloads is `2.00x`.
- [ ] Suggested safety factor for copy-heavy clone or large-payload workloads is `2.25x`.

### Payload classes
- [ ] Tiny payload benchmark uses scalar-like elements such as `int`.
- [ ] Large payload benchmark uses fixed-size struct or fixed-size array payload.
- [ ] Large-payload thresholds must be tracked independently from tiny-payload thresholds.

## Test Generator Hints
- Stress wrap-around behavior for circular-buffer implementations.
- Validate FIFO sequence against oracle under randomized operations.
- Validate `Clone/CloneWith` preserve `Cap()`, wrap-around order, and hook-call count.
- Check empty/full transitions and length correctness.
- Iterator tests must verify front-to-back order and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for queue API including enqueue/dequeue/peek on empty/full and wrap-around scenarios."
- Clone tests: "Generate tests for `Clone()` and `CloneWith(...)` verifying independent container state, shallow-copy default, nil-hook equivalence, custom hook behavior, preserved capacity, and live-element-only hook calls."
- Property tests: "Generate randomized queue operations with fixed seed and compare dequeue order to oracle FIFO model."
- Iterator tests: "Generate tests for `Values() iter.Seq[T]` checking front-to-back order, exact count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate queue benchmarks for enqueue-only, dequeue-only, and mixed workloads across 1e3/1e4/1e5 sizes."
