# SPECS.md - tree-avl

## Goal
Implement a set-like AVL tree for values of type `T`.

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics and comparator
- [ ] Use `T any`.
- [ ] Require comparator `cmp(a, b T) int`.
- [ ] Duplicate value policy: do not store duplicates.

## Required API
Each public API below is part of contract. Use each checklist as implementation target, not loose guidance.

### `New(cmp func(a, b T) int) *Tree[T]`
Purpose
- [ ] Create empty AVL tree with comparator-defined ordering.

Behavior expectations
- [ ] Returned tree is non-nil and empty.
- [ ] Comparator becomes persistent ordering rule for all future operations.
- [ ] `Len()` is `0` immediately after construction.

Performance expectations
- [ ] `O(1)` time.

### `Insert(v T) bool`
Purpose
- [ ] Insert new value into set-like ordered tree.

Behavior expectations
- [ ] Return `true` when value was not previously present.
- [ ] Return `false` on duplicate value and leave tree unchanged.
- [ ] BST ordering remains valid after insertion.
- [ ] AVL rotations and height updates restore balance before operation completes.

Performance expectations
- [ ] `O(log n)` time.

### `Delete(v T) bool`
Purpose
- [ ] Remove existing value from tree.

Behavior expectations
- [ ] Return `true` when value existed and was removed.
- [ ] Return `false` when value is missing.
- [ ] Delete of root with `0`, `1`, or `2` children is supported.
- [ ] BST ordering and AVL balance remain valid after deletion.

Performance expectations
- [ ] `O(log n)` time.

### `Has(v T) bool`
Purpose
- [ ] Report whether value currently exists.

Behavior expectations
- [ ] Empty tree safely returns `false`.
- [ ] Return `true` only for live stored values.
- [ ] Duplicate policy means at most one equal value can exist.

Performance expectations
- [ ] `O(log n)` time.

### `Min() (T, bool)`
Purpose
- [ ] Read smallest stored value.

Behavior expectations
- [ ] Empty tree returns `(zero, false)`.
- [ ] Non-empty tree returns leftmost live value under comparator ordering.
- [ ] `Min()` never mutates tree state.

Performance expectations
- [ ] `O(log n)` time.

### `Max() (T, bool)`
Purpose
- [ ] Read largest stored value.

Behavior expectations
- [ ] Empty tree returns `(zero, false)`.
- [ ] Non-empty tree returns rightmost live value under comparator ordering.
- [ ] `Max()` never mutates tree state.

Performance expectations
- [ ] `O(log n)` time.

### `Len() int`
Purpose
- [ ] Report number of live stored values.

Behavior expectations
- [ ] Count increases only for successful insertions.
- [ ] Count decreases only for successful deletions.
- [ ] Count becomes `0` after `Clear`.

Performance expectations
- [ ] `O(1)` time.

### `Clear()`
Purpose
- [ ] Reset tree to empty logical state.

Behavior expectations
- [ ] `Len()` becomes `0`.
- [ ] Root resets to empty-state sentinel.
- [ ] Future inserts, deletes, and queries still work correctly.
- [ ] Clear on empty tree is safe.

Performance expectations
- [ ] Target `O(1)` logical reset.
- [ ] `O(n)` node reclamation bookkeeping is allowed if implementation needs it.

### `Clone() *Tree[T]`
Purpose
- [ ] Create independent AVL tree with same sorted contents and comparator.

Behavior expectations
- [ ] Clone preserves `Len()`, comparator, and ascending in-order value sequence.
- [ ] Elements are copied with normal Go assignment.
- [ ] Clone is container-independent from source.
- [ ] Clone remains AVL-valid.

Performance expectations
- [ ] `O(n)` time for `n = Len()`.
- [ ] `O(n)` extra storage.

### `CloneWith(cloneValue func(T) T) *Tree[T]`
Purpose
- [ ] Create independent AVL tree while optionally transforming each live value.

Behavior expectations
- [ ] Preserves `Len()`, comparator, and ascending in-order sequence.
- [ ] Nil hook is equivalent to `Clone()`.
- [ ] Non-nil hook is called exactly once per live value.
- [ ] Hook call order is ascending in-order traversal.
- [ ] Hook must preserve comparator compatibility for cloned values.

Performance expectations
- [ ] `O(n)` container work plus hook cost.

### `RootNode() (NodeAPI[T], bool)`
Purpose
- [ ] Read root node through structural node-view API.

Behavior expectations
- [ ] Empty tree returns `(zero, false)` where zero is nil `NodeAPI[T]`.
- [ ] Non-empty tree returns `(rootNode, true)`.
- [ ] Returned node view is read-only.
- [ ] Node view becomes invalid after any tree mutation (`Insert`, `Delete`, `Clear`).

Performance expectations
- [ ] `O(1)` time.

### `NodeAPI[T]`
Purpose
- [ ] Provide read-only structural traversal power without exposing internal storage.

Behavior expectations
- [ ] `Value() T` returns node value.
- [ ] `ChildCount() int` returns direct child count in range `[0, 2]`.
- [ ] `Children() iter.Seq[NodeAPI[T]]` yields direct children in deterministic left-then-right order.
- [ ] `Children()` supports early stop when consumer returns `false`.
- [ ] Mutation during node traversal is not safe.

Performance expectations
- [ ] `Value()` is `O(1)`.
- [ ] `ChildCount()` is `O(1)`.
- [ ] Full `Children()` walk is `O(number of direct children)`.

### `InOrder() iter.Seq[T]`
Purpose
- [ ] Iterate live values in ascending sorted order.

Behavior expectations
- [ ] Yield each live value exactly once.
- [ ] Empty tree yields nothing.
- [ ] Early stop works when consumer returns `false`.
- [ ] Mutation during iteration is not safe.

Performance expectations
- [ ] Full traversal is `O(n)`.
- [ ] Early-stop traversal is `O(k)` for yielded prefix size `k`.
- [ ] Iterator setup is `O(1)` or `O(height)` if stack state is prepared lazily from tree root.

## Internal representation
- [ ] Use an index-based node pool; do not use pointer-linked nodes.
- [ ] Use `-1` as the nil/sentinel index and track root index explicitly.
- [ ] Store child indexes and height/balance data in arrays keyed by node index.
- [ ] Free-list for reusable node slots.
- [ ] AVL rotations for rebalance.

## Auto-resize policy
- [ ] No capacity-based `Grow()` or `Shrink()` API.
- [ ] When free-list is empty, allocate larger node-pool arrays and copy node fields by index.
- [ ] Allocate node storage on insert operations.
- [ ] Reclaim node storage on delete/clear operations.
- [ ] Free-list reuse is recommended to reduce allocation churn.

## Invariants
- [ ] BST ordering holds.
- [ ] Balance factor each node is in `[-1, 1]`.
- [ ] Height metadata is correct.
- [ ] `Len()` equals live node count.
- [ ] `Clone()` and `CloneWith(...)` preserve `Len()`, comparator, sorted order, and AVL invariants in clone.

## Iterator contract
- [ ] `InOrder()` yields ascending sorted values.
- [ ] Each value is yielded exactly once.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty tree yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.
- [ ] Mutation during `NodeAPI` traversal is not safe.

## Edge cases
- [ ] Insert/delete on empty tree.
- [ ] Duplicate insert returns `false` and leaves tree unchanged.
- [ ] Delete root with 0, 1, 2 children.
- [ ] All four rotation patterns are covered.
- [ ] `RootNode()` on empty tree returns `(zero, false)`.
- [ ] `Clone()` on empty tree returns empty independent tree with same comparator.
- [ ] `CloneWith(nil)` is equivalent to `Clone()`.
- [ ] `CloneWith(...)` calls `cloneValue` only for live values in ascending order.
- [ ] `cloneValue` must preserve comparator ordering contract for cloned values.

## Test checklist
- [ ] Sorted order tests.
- [ ] Rotation scenario tests.
- [ ] Random operations against reference set model.
- [ ] `RootNode/NodeAPI` tests for empty/non-empty root access, left-right child order, child-count consistency, and early stop in `Children()`.
- [ ] `RootNode/NodeAPI` DFS walk visits exactly `Len()` live nodes.
- [ ] `Clone/CloneWith` tests for independence, nil-hook equivalence, shallow-copy default, custom hook behavior, live-value-only hook calls, and preserved sorted/min/max behavior.

## Benchmark checklist
- [ ] Benchmark only non-trivial public APIs.
- [ ] Do not benchmark trivial metadata accessor `Len()`.
- [ ] `Insert` benchmark.
- [ ] `Delete` benchmark.
- [ ] `Has` benchmark.
- [ ] `Min` benchmark.
- [ ] `Max` benchmark.
- [ ] `Clear` benchmark.
- [ ] `Clone` benchmark.
- [ ] `CloneWith` benchmark.
- [ ] `InOrder` benchmark.
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
- [ ] Expected `O(log n)` APIs must keep `ns/op(1e5) <= 2.5 * ns/op(1e3)`.
- [ ] Expected `O(n)` APIs must compare normalized `(ns/op)/n` and keep normalized `1e5` result within `3.0x` of normalized `1e3` result.

### Layer 2: Absolute timing thresholds
- [ ] Absolute thresholds are per API, per payload class, and per benchmark size.
- [ ] Provisional thresholds should be derived from expected complexity and payload size.
- [ ] Calibrated thresholds should be derived from repeated local runs, median `ns/op`, and safety factor.
- [ ] Suggested safety factor for stable simple AVL ops is `1.75x`.
- [ ] Suggested safety factor for traversal workloads is `2.00x`.
- [ ] Suggested safety factor for copy-heavy clone or large-payload workloads is `2.25x`.

### Payload classes
- [ ] Tiny payload benchmark uses scalar-like values such as `int`.
- [ ] Large payload benchmark uses fixed-size struct or fixed-size array payload.
- [ ] Large-payload thresholds must be tracked independently from tiny-payload thresholds.

## Test Generator Hints
- Cover all AVL rebalancing patterns with deterministic cases.
- Use randomized set operations with fixed seed and sorted oracle.
- Validate `Clone/CloneWith` preserve sorted order, min/max, and hook-call count.
- Validate BST ordering, height metadata, and balance factor bounds.
- Validate `RootNode/NodeAPI` child ordering, child-count consistency, and full-node coverage.
- Iterator tests must verify sorted in-order output and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for set-like AVL API including duplicate inserts, delete-root cases, and min/max behavior."
- Clone tests: "Generate tests for `Clone()` and `CloneWith(...)` verifying independent container state, shallow-copy default, nil-hook equivalence, custom hook behavior, live-value-only hook calls, and preserved sorted/min/max behavior."
- Property tests: "Generate randomized insert/delete/has tests with fixed seed and oracle set comparison plus AVL invariant checks each step."
- Iterator tests: "Generate tests for `InOrder() iter.Seq[T]` enforcing sorted order, exact count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate AVL set benchmarks for insert/has/delete and full in-order traversal at 1e3/1e4/1e5 sizes."
