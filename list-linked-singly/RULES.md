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
- [ ] `Append(v T)` appends at tail in `O(1)` time.
- [ ] `DeleteFirst(match func(T) bool) bool` removes first node for which `match(value)` returns `true`. Return `false` when no such node exists.
- [ ] `Len() int`
- [ ] `Clear()`
- [ ] `Values() iter.Seq[T]`

## Internal representation
- [ ] Use an index-based node pool. Each live node has a value and `next` index.
- [ ] Use `-1` as the nil/sentinel index.
- [ ] Track head, tail, free-list head, and length.
- [ ] `Append` must use tracked tail; do not re-scan from head.
- [ ] Keep length counter.

## Auto-resize policy
- [ ] No capacity-based `Grow()` or `Shrink()` API.
- [ ] When free-list is empty, allocate larger node-pool arrays and copy node fields by index.
- [ ] Allocate node storage on insert operations.
- [ ] Reclaim node storage on delete/clear operations.
- [ ] Optional free-list reuse is allowed if list invariants remain valid.

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
- [ ] `Clear()` resets head, tail, free-list state, and length.

## Test checklist
- [ ] Head/tail transition tests.
- [ ] Delete-first behavior tests.
- [ ] Length consistency checks.

## Benchmark checklist
- [ ] PushFront/PopFront throughput.
- [ ] Append throughput for contracted O(1) tail-tracked Append.
- [ ] Full iteration benchmark.

## Test Generator Hints
- Focus head transitions, single-element transitions, and delete-first behavior.
- Validate traversal length equals `Len()` and no cycles are introduced.
- Use randomized operation streams and oracle comparison for final order.
- Iterator tests must verify head-to-tail order and early stop.
- Benchmarks must target contracted public API behavior.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for singly linked list API covering empty operations, delete-first cases, and append behavior."
- Property tests: "Generate randomized list mutations with fixed seed and verify no-cycle + length invariants after each batch."
- Iterator tests: "Generate tests for `Values() iter.Seq[T]` enforcing head-to-tail order, count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate singly linked list benchmarks for push/pop-front, contracted O(1) append, and full traversal at 1e3/1e4/1e5."
