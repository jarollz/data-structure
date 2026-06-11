# SPECS.md - list-linked-singly

## Goal
Implement a generic singly linked list.

## Hard constraints
- [ ] Use only simple Go types and arrays for storage models you choose.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics
- [ ] Use `T any` for element type.

## Required API
Each public API below is part of contract. Use each checklist as implementation target, not loose guidance.

### `New() *List[T]`
Purpose
- [ ] Create empty singly linked list.

Behavior expectations
- [ ] Returned list is non-nil and empty.
- [ ] `Len()` is `0` immediately after construction.
- [ ] Head and tail both represent empty state.

Performance expectations
- [ ] `O(1)` time.

### `PushFront(v T)`
Purpose
- [ ] Insert value as new head node.

Behavior expectations
- [ ] New value becomes first value yielded by `Values()`.
- [ ] Pushing into empty list updates both head and tail to same live node.
- [ ] Existing relative order after new head is preserved.

Performance expectations
- [ ] `O(1)` time.

### `PopFront() (T, bool)`
Purpose
- [ ] Remove and return current head value.

Behavior expectations
- [ ] Empty list returns `(zero, false)` and does not panic.
- [ ] Non-empty list returns current head value.
- [ ] Single-node pop leaves list empty and resets tail correctly.

Performance expectations
- [ ] `O(1)` time.

### `Append(v T)`
Purpose
- [ ] Insert value as new tail node.

Behavior expectations
- [ ] New value becomes last value yielded by `Values()`.
- [ ] Appending into empty list updates both head and tail to same live node.
- [ ] Append must use tracked tail, not head scan.

Performance expectations
- [ ] `O(1)` time by contract.

### `DeleteFirst(match func(T) bool) bool`
Purpose
- [ ] Remove first node whose value matches predicate.

Behavior expectations
- [ ] Visit nodes from head to tail.
- [ ] Remove only first matching node.
- [ ] Return `false` when no match exists.
- [ ] Deleting head and deleting tail both update structure correctly.
- [ ] List order of remaining nodes is preserved.

Performance expectations
- [ ] `O(n)` time in number of live nodes.

### `Len() int`
Purpose
- [ ] Report number of live nodes.

Behavior expectations
- [ ] Count increases after `PushFront` and `Append`.
- [ ] Count decreases after successful `PopFront` and `DeleteFirst`.
- [ ] Count becomes `0` after `Clear`.

Performance expectations
- [ ] `O(1)` time.

### `Clear()`
Purpose
- [ ] Reset list to empty logical state.

Behavior expectations
- [ ] `Len()` becomes `0`.
- [ ] Head, tail, and free-list state reset consistently.
- [ ] Future `PushFront`, `Append`, and `PopFront` still work correctly.
- [ ] Clear on empty list is safe.

Performance expectations
- [ ] Target `O(1)` logical reset.
- [ ] `O(n)` node reclamation bookkeeping is allowed if implementation needs it.

### `Clone() *List[T]`
Purpose
- [ ] Create independent singly linked list with same visible contents.

Behavior expectations
- [ ] Clone preserves `Len()` and head-to-tail order.
- [ ] Elements are copied with normal Go assignment.
- [ ] Clone is container-independent from source.
- [ ] Empty source produces empty independent clone.

Performance expectations
- [ ] `O(n)` time for `n = Len()`.
- [ ] `O(n)` extra storage.

### `CloneWith(cloneValue func(T) T) *List[T]`
Purpose
- [ ] Create independent singly linked list while optionally transforming each live node value.

Behavior expectations
- [ ] Preserves `Len()` and head-to-tail order.
- [ ] Nil hook is equivalent to `Clone()`.
- [ ] Non-nil hook is called exactly once per live node.
- [ ] Hook call order is head to tail.
- [ ] Hook is never called for reclaimed or free-list nodes.

Performance expectations
- [ ] `O(n)` container work plus hook cost.

### `Values() iter.Seq[T]`
Purpose
- [ ] Iterate live nodes from head to tail.

Behavior expectations
- [ ] Yield each live node exactly once.
- [ ] Empty list yields nothing.
- [ ] Early stop works when consumer returns `false`.
- [ ] Mutation during iteration is not safe.

Performance expectations
- [ ] Full traversal is `O(n)`.
- [ ] Early-stop traversal is `O(k)` for yielded prefix size `k`.
- [ ] Iterator setup is `O(1)`.

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
- [ ] `Clone()` and `CloneWith(...)` preserve `Len()` and head-to-tail order in clone.

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
- [ ] `Clone()` on empty list returns empty independent list.
- [ ] `CloneWith(nil)` is equivalent to `Clone()`.
- [ ] `CloneWith(...)` calls `cloneValue` only for live nodes, never for reclaimed or free-list nodes.

## Test checklist
- [ ] Head/tail transition tests.
- [ ] Delete-first behavior tests.
- [ ] Length consistency checks.
- [ ] `Clone/CloneWith` tests for independence, nil-hook equivalence, shallow-copy default, custom hook behavior, and live-node-only hook calls.

## Benchmark checklist
- [ ] Benchmark only non-trivial public APIs.
- [ ] Do not benchmark trivial metadata accessor `Len()`.
- [ ] `PushFront` benchmark.
- [ ] `PopFront` benchmark.
- [ ] `Append` benchmark.
- [ ] `DeleteFirst` benchmark.
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
- Focus head transitions, single-element transitions, and delete-first behavior.
- Validate traversal length equals `Len()` and no cycles are introduced.
- Validate `Clone/CloneWith` preserve head-to-tail order and apply hook once per live node.
- Use randomized operation streams and oracle comparison for final order.
- Iterator tests must verify head-to-tail order and early stop.
- Benchmarks must target contracted public API behavior.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for singly linked list API covering empty operations, delete-first cases, and append behavior."
- Clone tests: "Generate tests for `Clone()` and `CloneWith(...)` verifying independent container state, shallow-copy default, nil-hook equivalence, custom hook behavior, and live-node-only hook calls."
- Property tests: "Generate randomized list mutations with fixed seed and verify no-cycle + length invariants after each batch."
- Iterator tests: "Generate tests for `Values() iter.Seq[T]` enforcing head-to-tail order, count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate singly linked list benchmarks for push/pop-front, contracted O(1) append, and full traversal at 1e3/1e4/1e5."
