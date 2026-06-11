# SPECS.md - list-linked-doubly

## Goal
Implement a generic doubly linked list.

## Hard constraints
- [ ] Use only simple Go types and arrays for storage models you choose.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics
- [ ] Use `T any`.

## Required API
Each public API below is part of contract. Use each checklist as implementation target, not loose guidance.

### `New() *List[T]`
Purpose
- [ ] Create empty doubly linked list.

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
- [ ] New value becomes first yielded value.
- [ ] Empty-list push updates both head and tail.
- [ ] Old head becomes second node with correct `prev` linkage.

Performance expectations
- [ ] `O(1)` time by contract.

### `PushBack(v T)`
Purpose
- [ ] Insert value as new tail node.

Behavior expectations
- [ ] New value becomes last yielded value.
- [ ] Empty-list push updates both head and tail.
- [ ] Old tail becomes previous node with correct `next` linkage.

Performance expectations
- [ ] `O(1)` time by contract.

### `PopFront() (T, bool)`
Purpose
- [ ] Remove and return current head value.

Behavior expectations
- [ ] Empty list returns `(zero, false)` and does not panic.
- [ ] Non-empty list returns current head value.
- [ ] Single-node pop leaves list empty with both ends reset.

Performance expectations
- [ ] `O(1)` time.

### `PopBack() (T, bool)`
Purpose
- [ ] Remove and return current tail value.

Behavior expectations
- [ ] Empty list returns `(zero, false)` and does not panic.
- [ ] Non-empty list returns current tail value.
- [ ] Single-node pop leaves list empty with both ends reset.

Performance expectations
- [ ] `O(1)` time.

### `Len() int`
Purpose
- [ ] Report number of live nodes.

Behavior expectations
- [ ] Count increases after pushes.
- [ ] Count decreases after successful pops.
- [ ] Count becomes `0` after `Clear`.

Performance expectations
- [ ] `O(1)` time.

### `Clear()`
Purpose
- [ ] Reset list to empty logical state.

Behavior expectations
- [ ] `Len()` becomes `0`.
- [ ] Head, tail, and free-list state reset consistently.
- [ ] Future push and pop operations still work correctly.
- [ ] Clear on empty list is safe.

Performance expectations
- [ ] Target `O(1)` logical reset.
- [ ] `O(n)` node reclamation bookkeeping is allowed if implementation needs it.

### `Clone() *List[T]`
Purpose
- [ ] Create independent doubly linked list with same visible contents.

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
- [ ] Create independent doubly linked list while optionally transforming each live node value.

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
- [ ] Use an index-based node pool. Each live node has value, `prev`, and `next` indexes.
- [ ] Use `-1` as the nil/sentinel index.
- [ ] Track head, tail, free-list head, and length.
- [ ] Keep length counter.

## Auto-resize policy
- [ ] No capacity-based `Grow()` or `Shrink()` API.
- [ ] When free-list is empty, allocate larger node-pool arrays and copy node fields by index.
- [ ] Allocate node storage on insert operations.
- [ ] Reclaim node storage on delete/clear operations.
- [ ] Optional free-list reuse is allowed if list invariants remain valid.

## Invariants
- [ ] `head.prev` is nil/sentinel.
- [ ] `tail.next` is nil/sentinel.
- [ ] For adjacent nodes, `a.next == b` implies `b.prev == a`.
- [ ] Forward traversal count equals `Len()` and matches backward traversal.
- [ ] `Clone()` and `CloneWith(...)` preserve `Len()` and head-to-tail order in clone.

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
- [ ] After removing last node, both head and tail become `-1`.
- [ ] `Clone()` on empty list returns empty independent list.
- [ ] `CloneWith(nil)` is equivalent to `Clone()`.
- [ ] `CloneWith(...)` calls `cloneValue` only for live nodes, never for reclaimed or free-list nodes.

## Test checklist
- [ ] Prev/next consistency tests.
- [ ] Head/tail update tests.
- [ ] Forward/backward traversal consistency tests.
- [ ] `Clone/CloneWith` tests for independence, nil-hook equivalence, shallow-copy default, custom hook behavior, and live-node-only hook calls.

## Benchmark checklist
- [ ] Benchmark only non-trivial public APIs.
- [ ] Do not benchmark trivial metadata accessor `Len()`.
- [ ] `PushFront` benchmark.
- [ ] `PushBack` benchmark.
- [ ] `PopFront` benchmark.
- [ ] `PopBack` benchmark.
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
- Stress empty/single/multi-node transitions on both ends.
- Validate `prev/next` consistency and head/tail boundary invariants.
- Validate `Clone/CloneWith` preserve traversal order and apply hook once per live node.
- Use randomized sequences for push/pop combinations.
- Iterator tests must verify head-to-tail order and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for doubly linked list API including pop-front/pop-back edge cases and clear resets."
- Clone tests: "Generate tests for `Clone()` and `CloneWith(...)` verifying independent container state, shallow-copy default, nil-hook equivalence, custom hook behavior, and live-node-only hook calls."
- Property tests: "Generate randomized end-operations with fixed seed and assert prev/next consistency invariants after each batch."
- Iterator tests: "Generate tests for `Values() iter.Seq[T]` checking exact count, traversal order, early stop, and mutation-unsafety note."
- Benchmarks: "Generate doubly linked list benchmarks for both-end updates and full iteration at 1e3/1e4/1e5 sizes."
