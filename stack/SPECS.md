# SPECS.md - stack

## Goal
Implement a generic LIFO stack.

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics
- [ ] Use `T any`.

## Required API
Each public API below is part of contract. Use each checklist as implementation target, not loose guidance.

### `New(capacity int) *Stack[T]`
Purpose
- [ ] Create empty stack with normalized starting capacity.

Behavior expectations
- [ ] Normalize `capacity <= 0` to `16`.
- [ ] Effective starting capacity is `max(16, capacity)`.
- [ ] Returned stack is non-nil and empty.
- [ ] `Len()` is `0` immediately after construction.
- [ ] `Cap()` is normalized starting capacity immediately after construction.

Performance expectations
- [ ] `O(1)` time.
- [ ] `O(capacity)` backing storage.

Implementation checklist
- [ ] Store explicit logical length.
- [ ] Keep normalized start capacity as `minCap` for shrink policy.

### `Push(v T) bool`
Purpose
- [ ] Add value on logical top of stack.

Behavior expectations
- [ ] New value becomes next value returned by `PeekTop()`.
- [ ] New value becomes first value returned by later `Pop()`.
- [ ] Grow backing storage before write when full.
- [ ] Method always returns `true`.

Performance expectations
- [ ] Amortized `O(1)` time.
- [ ] Worst-case single call may be `O(n)` during resize.

Implementation checklist
- [ ] Top value is stored at index `Len()-1` after push.
- [ ] Resize copy preserves existing bottom-to-top storage order.

### `Pop() (T, bool)`
Purpose
- [ ] Remove and return current top value.

Behavior expectations
- [ ] Non-empty stack returns most recently pushed value not yet popped.
- [ ] Empty stack returns `(zero, false)` and does not panic.
- [ ] Single-element pop leaves stack empty.
- [ ] Shrink may happen after removal if resize policy says so.

Performance expectations
- [ ] Amortized `O(1)` time.
- [ ] Worst-case single call may be `O(n)` during shrink copy.

Implementation checklist
- [ ] Read top value before decrementing length.
- [ ] Preserve LIFO order after shrink.

### `PeekTop() (T, bool)`
Purpose
- [ ] Read current top value without removing it.

Behavior expectations
- [ ] Non-empty stack returns current top value and `true`.
- [ ] Empty stack returns `(zero, false)` and does not panic.
- [ ] `PeekTop()` never changes `Len()`, `Cap()`, or stack order.

Performance expectations
- [ ] `O(1)` time.

### `Len() int`
Purpose
- [ ] Report number of live stack elements.

Behavior expectations
- [ ] Count increases after `Push`.
- [ ] Count decreases after successful `Pop`.
- [ ] Count becomes `0` after `Clear`.

Performance expectations
- [ ] `O(1)` time.

### `Cap() int`
Purpose
- [ ] Report current backing-storage capacity.

Behavior expectations
- [ ] Capacity is storage size, not live element count.
- [ ] Capacity never drops below normalized starting capacity.
- [ ] Capacity reflects grow and shrink operations.

Performance expectations
- [ ] `O(1)` time.

### `Clear()`
Purpose
- [ ] Reset stack to empty logical state.

Behavior expectations
- [ ] `Len()` becomes `0`.
- [ ] Later `Pop()` and `PeekTop()` behave like empty stack operations.
- [ ] Future `Push()` calls still work correctly.
- [ ] Clear on empty stack is safe.

Performance expectations
- [ ] Target `O(1)` logical reset.
- [ ] `O(n)` slot reset is allowed only if implementation needs it.

### `Clone() *Stack[T]`
Purpose
- [ ] Create independent stack container with same visible contents.

Behavior expectations
- [ ] Clone preserves `Len()`, `Cap()`, and top-to-bottom order.
- [ ] Elements are copied with normal Go assignment.
- [ ] Clone is container-independent from source.
- [ ] Empty source produces empty independent clone with same `Cap()`.

Performance expectations
- [ ] `O(n)` time for `n = Len()`.
- [ ] `O(n)` extra storage.

Implementation checklist
- [ ] Copy only live elements.
- [ ] Preserve observable top-to-bottom order in clone.

### `CloneWith(cloneValue func(T) T) *Stack[T]`
Purpose
- [ ] Create independent stack container while optionally transforming each live element.

Behavior expectations
- [ ] Preserves `Len()`, `Cap()`, and top-to-bottom order.
- [ ] Nil hook is equivalent to `Clone()`.
- [ ] Non-nil hook is called exactly once per live element.
- [ ] Hook call order is top to bottom.
- [ ] Hook is never called for unused capacity slots.

Performance expectations
- [ ] `O(n)` container work plus hook cost.

Implementation checklist
- [ ] Apply hook only to live elements.
- [ ] Preserve same `Cap()` as source stack.

### `Values() iter.Seq[T]`
Purpose
- [ ] Iterate live stack values from top to bottom.

Behavior expectations
- [ ] Yield current top element first.
- [ ] Yield each live element exactly once.
- [ ] Empty stack yields nothing.
- [ ] Early stop works when consumer returns `false`.
- [ ] Mutation during iteration is not safe.

Performance expectations
- [ ] Full traversal is `O(n)`.
- [ ] Early-stop traversal is `O(k)` for yielded prefix size `k`.
- [ ] Iterator setup is `O(1)`.

## Internal representation
- [ ] Array-backed storage with explicit length. Top element is at index `Len()-1`.
- [ ] When full, grow by allocating new backing storage and copying live elements in order.

## Auto-resize policy
- [ ] `Cap()` immediately after `New` equals normalized starting capacity.
- [ ] Grow when `Len() == Cap()`.
- [ ] New capacity on grow: `2x` when `Cap() < 1024`, otherwise `Cap() + Cap()/2`.
- [ ] Shrink when `Len() <= Cap()/4` and `Cap() > minCap`.
- [ ] New capacity on shrink: `max(minCap, Cap()/2, 2*Len())`.
- [ ] `minCap` is normalized starting capacity.
- [ ] Use hysteresis; do not resize on every pop.

## Invariants
- [ ] Top element is most recently pushed not yet popped.
- [ ] `Len()` equals number of live elements.
- [ ] No out-of-range access for top index.
- [ ] `Clone()` and `CloneWith(...)` preserve `Len()`, `Cap()`, and top-to-bottom order in clone.

## Iterator contract
- [ ] `Values()` yields top to bottom.
- [ ] Each live element is yielded exactly once.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty stack yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Pop/peek on empty returns `(zero, false)`.
- [ ] Push at full capacity grows backing storage and still returns `true`.
- [ ] Clear resets stack state.
- [ ] `Clone()` on empty stack returns empty independent stack with same `Cap()`.
- [ ] `CloneWith(nil)` is equivalent to `Clone()`.
- [ ] `CloneWith(...)` calls `cloneValue` only for live stack elements from top to bottom.

## Test checklist
- [ ] Push/pop ordering tests.
- [ ] Empty and single-element transition tests.
- [ ] Random operation tests against reference model.
- [ ] `Clone/CloneWith` tests for independence, nil-hook equivalence, shallow-copy default, custom hook behavior, preserved capacity, and live-element-only hook calls.

## Benchmark checklist
- [ ] Benchmark only non-trivial public APIs.
- [ ] Do not benchmark trivial metadata accessors such as `Len()` and `Cap()`.
- [ ] `Push` benchmark.
- [ ] `Pop` benchmark.
- [ ] `PeekTop` benchmark.
- [ ] `Clear` benchmark.
- [ ] `Clone` benchmark.
- [ ] `CloneWith` benchmark.
- [ ] `Values` benchmark.
- [ ] Mixed push/pop benchmark.
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
- [ ] Suggested safety factor for stable simple stack ops is `1.50x`.
- [ ] Suggested safety factor for iteration and mixed workloads is `2.00x`.
- [ ] Suggested safety factor for copy-heavy clone or large-payload workloads is `2.25x`.

### Payload classes
- [ ] Tiny payload benchmark uses scalar-like elements such as `int`.
- [ ] Large payload benchmark uses fixed-size struct or fixed-size array payload.
- [ ] Large-payload thresholds must be tracked independently from tiny-payload thresholds.

## Test Generator Hints
- Validate strict LIFO order under interleaved push/pop operations.
- Cover empty and single-element transitions for pop/peek.
- Validate `Clone/CloneWith` preserve `Cap()`, top-to-bottom order, and hook-call count.
- Use randomized operation streams with oracle stack comparison.
- Iterator tests must verify top-to-bottom order and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for stack API including push/pop/peek behavior for empty, single, and multi-element states."
- Clone tests: "Generate tests for `Clone()` and `CloneWith(...)` verifying independent container state, shallow-copy default, nil-hook equivalence, custom hook behavior, preserved capacity, and live-element-only hook calls."
- Property tests: "Generate randomized push/pop with fixed seed and compare to oracle LIFO model while checking invariants."
- Iterator tests: "Generate tests for `Values() iter.Seq[T]` enforcing top-to-bottom order, exact count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate stack benchmarks for push-only, pop-only, and mixed workloads at 1e3/1e4/1e5 sizes."
