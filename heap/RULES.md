# RULES.md - heap

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
- [ ] `New(capacity int, cmp func(a, b T) int) *Heap[T]` creates an empty heap. Normalize `capacity <= 0` to `16`. Effective starting capacity is `max(16, capacity)`.
- [ ] `Push(v T) bool` inserts `v`, restores heap order, and always returns `true`.
- [ ] `PopTop() (T, bool)`
- [ ] `PeekTop() (T, bool)`
- [ ] `Len() int`
- [ ] `Cap() int`
- [ ] `Clear()`
- [ ] `Values() iter.Seq[T]`

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

## Test checklist
- [ ] Push/pop correctness tests.
- [ ] Heap-property check after each mutation batch.
- [ ] Randomized compare against sorted reference behavior.

## Benchmark checklist
- [ ] Push-only benchmark.
- [ ] Pop-only benchmark after prefill.
- [ ] Mixed priority-queue workload benchmark.

## Test Generator Hints
- Verify heap property after every mutating operation sequence.
- Include duplicate priorities and comparator edge cases.
- Check repeated `PopTop` yields monotonic priority order.
- Iterator tests must validate count, internal-order expectation, and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for heap API (`Push`, `PopTop`, `PeekTop`) including empty and full-capacity behavior."
- Property tests: "Generate randomized push/pop sequences with fixed seed and assert heap property after each batch."
- Iterator tests: "Generate tests for `Values() iter.Seq[T]` ensuring exact count, early stop, and documented non-sorted internal order."
- Benchmarks: "Generate heap benchmarks for push-only, pop-only, and mixed workloads at 1e3/1e4/1e5 sizes."
