# RULES.md - list-linked-singly

## Goal
Implement a generic singly linked list.

## Hard constraints
- [ ] Use only simple Go types and arrays for storage models you choose.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics
- [ ] Use `T any` for element type.

## Required API
- [ ] `New() *List[T]`
- [ ] `PushFront(v T)`
- [ ] `PopFront() (T, bool)`
- [ ] `Append(v T)`
- [ ] `DeleteFirst(match func(T) bool) bool`
- [ ] `Len() int`
- [ ] `Clear()`
- [ ] `Values() iter.Seq[T]`

## Internal representation
- [ ] Each node stores value and next reference/index.
- [ ] Track head (and tail if append should be `O(1)`).
- [ ] Keep length counter.

## Invariants
- [ ] Traversing from head visits exactly `Len()` nodes.
- [ ] Last node points to nil/sentinel.
- [ ] No cycles in normal operations.

## Iterator contract
- [ ] `Values()` yields from head to tail.
- [ ] Each element is yielded exactly once.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty list yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Pop on empty list returns `(zero, false)`.
- [ ] Delete from empty list is safe.
- [ ] Delete head node works.
- [ ] Single-element list transitions are correct.

## Test checklist
- [ ] Head/tail transition tests.
- [ ] Delete-first behavior tests.
- [ ] Length consistency checks.

## Benchmark checklist
- [ ] PushFront/PopFront throughput.
- [ ] Append throughput (with and without tail tracking).
- [ ] Full iteration benchmark.

## Test Generator Hints
- Focus head transitions, single-element transitions, and delete-first behavior.
- Validate traversal length equals `Len()` and no cycles are introduced.
- Use randomized operation streams and oracle comparison for final order.
- Iterator tests must verify head-to-tail order and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for singly linked list API covering empty operations, delete-first cases, and append behavior."
- Property tests: "Generate randomized list mutations with fixed seed and verify no-cycle + length invariants after each batch."
- Iterator tests: "Generate tests for `Values() iter.Seq[T]` enforcing head-to-tail order, count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate singly linked list benchmarks for push/pop-front, append, and full traversal at 1e3/1e4/1e5."
