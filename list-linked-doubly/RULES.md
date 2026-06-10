# RULES.md - list-linked-doubly

## Goal
Implement a generic doubly linked list.

## Hard constraints
- [ ] Use only simple Go types and arrays for storage models you choose.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics
- [ ] Use `T any`.

## Required API
- [ ] `New() *List[T]`
- [ ] `PushFront(v T)`
- [ ] `PushBack(v T)`
- [ ] `PopFront() (T, bool)`
- [ ] `PopBack() (T, bool)`
- [ ] `Len() int`
- [ ] `Clear()`
- [ ] `Values() iter.Seq[T]`

## Internal representation
- [ ] Node contains value, prev, and next.
- [ ] Track head and tail.
- [ ] Keep length counter.

## Auto-resize policy
- [ ] No capacity-based `Grow()` or `Shrink()` API.
- [ ] Allocate node storage on insert operations.
- [ ] Reclaim node storage on delete/clear operations.
- [ ] Optional free-list reuse is allowed if list invariants remain valid.

## Invariants
- [ ] `head.prev` is nil/sentinel.
- [ ] `tail.next` is nil/sentinel.
- [ ] For adjacent nodes, `a.next == b` implies `b.prev == a`.
- [ ] Forward traversal count equals `Len()` and matches backward traversal.

## Iterator contract
- [ ] `Values()` yields from head to tail.
- [ ] Each element is yielded exactly once.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty list yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Empty and single-node transitions are correct.
- [ ] Pop on empty returns `(zero, false)`.
- [ ] Clear resets all structural fields.

## Test checklist
- [ ] Prev/next consistency tests.
- [ ] Head/tail update tests.
- [ ] Forward/backward traversal consistency tests.

## Benchmark checklist
- [ ] Push/pop both ends benchmark.
- [ ] Iteration benchmark.

## Test Generator Hints
- Stress empty/single/multi-node transitions on both ends.
- Validate `prev/next` consistency and head/tail boundary invariants.
- Use randomized sequences for push/pop combinations.
- Iterator tests must verify head-to-tail order and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for doubly linked list API including pop-front/pop-back edge cases and clear resets."
- Property tests: "Generate randomized end-operations with fixed seed and assert prev/next consistency invariants after each batch."
- Iterator tests: "Generate tests for `Values() iter.Seq[T]` checking exact count, traversal order, early stop, and mutation-unsafety note."
- Benchmarks: "Generate doubly linked list benchmarks for both-end updates and full iteration at 1e3/1e4/1e5 sizes."
