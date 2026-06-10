# RULES.md - queue

## Goal
Implement a generic FIFO queue.

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics
- [ ] Use `T any`.

## Required API
- [ ] `New(capacity int) *Queue[T]`
- [ ] `Enqueue(v T) bool`
- [ ] `Dequeue() (T, bool)`
- [ ] `PeekFront() (T, bool)`
- [ ] `Len() int`
- [ ] `Cap() int`
- [ ] `Clear()`
- [ ] `Values() iter.Seq[T]`

## Internal representation
- [ ] Recommended: circular buffer using array.
- [ ] Track head index, tail index, and size.
- [ ] Wrap indices correctly.

## Auto-resize policy
- [ ] Grow when `Len() == Cap()`.
- [ ] New capacity on grow: `2x` when `Cap() < 1024`, otherwise `Cap() + Cap()/2`.
- [ ] Shrink when `Len() <= Cap()/4` and `Cap() > minCap`.
- [ ] New capacity on shrink: `max(minCap, Cap()/2, 2*Len())`.
- [ ] `minCap` is `max(16, initial capacity)`.
- [ ] On resize, re-pack logical order front-to-back into new array.
- [ ] Use hysteresis; do not resize on every dequeue.

## Invariants
- [ ] FIFO order is preserved.
- [ ] `Len()` equals number of queued elements.
- [ ] Head and tail positions are valid modulo capacity.

## Iterator contract
- [ ] `Values()` yields front to back.
- [ ] Each live element is yielded exactly once.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty queue yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Dequeue/peek on empty returns `(zero, false)`.
- [ ] Full-capacity behavior is clearly defined.
- [ ] Wrap-around enqueue/dequeue remains correct.

## Test checklist
- [ ] FIFO order tests.
- [ ] Wrap-around index tests.
- [ ] Random operation sequence tests.

## Benchmark checklist
- [ ] Enqueue benchmark.
- [ ] Dequeue benchmark.
- [ ] Mixed enqueue/dequeue benchmark.

## Test Generator Hints
- Stress wrap-around behavior for circular-buffer implementations.
- Validate FIFO sequence against oracle under randomized operations.
- Check empty/full transitions and length correctness.
- Iterator tests must verify front-to-back order and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for queue API including enqueue/dequeue/peek on empty/full and wrap-around scenarios."
- Property tests: "Generate randomized queue operations with fixed seed and compare dequeue order to oracle FIFO model."
- Iterator tests: "Generate tests for `Values() iter.Seq[T]` checking front-to-back order, exact count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate queue benchmarks for enqueue-only, dequeue-only, and mixed workloads across 1e3/1e4/1e5 sizes."
