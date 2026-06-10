# RULES.md - stack

## Goal
Implement a generic LIFO stack.

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics
- [ ] Use `T any`.

## Required API
- [ ] `New(capacity int) *Stack[T]`
- [ ] `Push(v T) bool`
- [ ] `Pop() (T, bool)`
- [ ] `PeekTop() (T, bool)`
- [ ] `Len() int`
- [ ] `Cap() int`
- [ ] `Clear()`
- [ ] `Values() iter.Seq[T]`

## Internal representation
- [ ] Array-backed storage with top index.
- [ ] Define overflow behavior (fail or grow by new array copy).

## Auto-resize policy
- [ ] Grow when `Len() == Cap()`.
- [ ] New capacity on grow: `2x` when `Cap() < 1024`, otherwise `Cap() + Cap()/2`.
- [ ] Shrink when `Len() <= Cap()/4` and `Cap() > minCap`.
- [ ] New capacity on shrink: `max(minCap, Cap()/2, 2*Len())`.
- [ ] `minCap` is `max(16, initial capacity)`.
- [ ] Use hysteresis; do not resize on every pop.

## Invariants
- [ ] Top element is most recently pushed not yet popped.
- [ ] `Len()` equals number of live elements.
- [ ] No out-of-range access for top index.

## Iterator contract
- [ ] `Values()` yields top to bottom.
- [ ] Each live element is yielded exactly once.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty stack yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Pop/peek on empty returns `(zero, false)`.
- [ ] Push at full capacity follows documented policy.
- [ ] Clear resets stack state.

## Test checklist
- [ ] Push/pop ordering tests.
- [ ] Empty and single-element transition tests.
- [ ] Random operation tests against reference model.

## Benchmark checklist
- [ ] Push benchmark.
- [ ] Pop benchmark.
- [ ] Mixed push/pop benchmark.

## Test Generator Hints
- Validate strict LIFO order under interleaved push/pop operations.
- Cover empty and single-element transitions for pop/peek.
- Use randomized operation streams with oracle stack comparison.
- Iterator tests must verify top-to-bottom order and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for stack API including push/pop/peek behavior for empty, single, and multi-element states."
- Property tests: "Generate randomized push/pop with fixed seed and compare to oracle LIFO model while checking invariants."
- Iterator tests: "Generate tests for `Values() iter.Seq[T]` enforcing top-to-bottom order, exact count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate stack benchmarks for push-only, pop-only, and mixed workloads at 1e3/1e4/1e5 sizes."
