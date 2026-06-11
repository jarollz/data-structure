# SPECS.md - tree-general

## Goal
Implement a generic n-ary tree for hierarchical data.

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics
- [ ] Use `T any` for node value.

## Required API
Each public API below is part of contract. Use each checklist as implementation target, not loose guidance.

### `New(rootValue T) *Tree[T]`
Purpose
- [ ] Create non-empty general tree with one root node.

Behavior expectations
- [ ] Root node ID is always `0`.
- [ ] Returned tree is non-nil.
- [ ] `Len()` is `1` immediately after construction.
- [ ] `Get(0)` returns `rootValue`.
- [ ] Initial next allocated ID after root is `1`.

Performance expectations
- [ ] `O(1)` time.

### `AddChild(parentID int, value T) (childID int, ok bool)`
Purpose
- [ ] Add new last child under existing live parent node.

Behavior expectations
- [ ] Invalid or removed `parentID` returns `(-1, false)`.
- [ ] Empty tree returns `(-1, false)`.
- [ ] New child is appended after any existing children of parent.
- [ ] Child IDs are stable, monotonically increasing, and never reused.
- [ ] Parent relationship and sibling order are updated consistently.

Performance expectations
- [ ] `O(1)` target for direct-link updates, excluding any node-array growth copy.

### `RemoveSubtree(nodeID int) bool`
Purpose
- [ ] Remove target node and all live descendants.

Behavior expectations
- [ ] Invalid or removed ID returns `false`.
- [ ] Removing root empties tree.
- [ ] Removing non-root updates parent child links and sibling links consistently.
- [ ] Public IDs of removed nodes are never reused.
- [ ] Removed-ID holes remain observable through later clone behavior and future ID progression.

Performance expectations
- [ ] `O(size of removed subtree)` time.

### `Get(nodeID int) (T, bool)`
Purpose
- [ ] Read value stored at live node ID.

Behavior expectations
- [ ] Invalid or removed ID returns `(zero, false)`.
- [ ] Live ID returns stored value and `true`.
- [ ] `Get` never mutates tree structure.

Performance expectations
- [ ] `O(1)` time.

### `Parent(nodeID int) (int, bool)`
Purpose
- [ ] Report parent ID for live non-root node.

Behavior expectations
- [ ] Root, invalid, or removed ID returns `(-1, false)`.
- [ ] Live non-root node returns direct parent ID and `true`.
- [ ] `Parent` never mutates tree structure.

Performance expectations
- [ ] `O(1)` time.

### `ChildCount(nodeID int) int`
Purpose
- [ ] Report number of live direct children for node.

Behavior expectations
- [ ] Invalid or removed ID returns `-1`.
- [ ] Live node returns count of current direct children only.
- [ ] Count changes after child insertion and subtree removal.

Performance expectations
- [ ] `O(number of direct children)` if counted by traversal.
- [ ] `O(1)` is allowed if implementation tracks child count explicitly.

### `Len() int`
Purpose
- [ ] Report number of live nodes.

Behavior expectations
- [ ] Count increases after successful `AddChild`.
- [ ] Count decreases by full removed subtree size after successful `RemoveSubtree`.
- [ ] Empty tree reports `0`.

Performance expectations
- [ ] `O(1)` time.

### `Clone() *Tree[T]`
Purpose
- [ ] Create independent general tree with same public ID behavior and hierarchy.

Behavior expectations
- [ ] Clone preserves live IDs, removed-ID holes, child order, parent relationships, root ID, and next-ID progression.
- [ ] Node values are copied with normal Go assignment.
- [ ] Clone is container-independent from source.
- [ ] Empty source remains empty in clone and preserves future ID progression.

Performance expectations
- [ ] Let `idRange` be next allocated ID, covering live IDs plus removed-ID holes in `[0, idRange)`.
- [ ] `O(idRange)` time.
- [ ] `O(idRange)` extra storage.

### `CloneWith(cloneValue func(T) T) *Tree[T]`
Purpose
- [ ] Create independent general tree while optionally transforming each live node value.

Behavior expectations
- [ ] Preserves live IDs, removed-ID holes, child order, parent relationships, root ID, and next-ID progression.
- [ ] Nil hook is equivalent to `Clone()`.
- [ ] Non-nil hook is called exactly once per live node.
- [ ] Hook call order is pre-order traversal.
- [ ] Hook is never called for removed-ID holes.

Performance expectations
- [ ] Let `idRange` be next allocated ID, covering live IDs plus removed-ID holes in `[0, idRange)`.
- [ ] `O(idRange)` container work plus hook cost.

### `PreOrder() iter.Seq[T]`
Purpose
- [ ] Iterate live node values in parent-before-children order.

Behavior expectations
- [ ] Parent value is yielded before any of its children.
- [ ] Sibling order follows stored child order.
- [ ] Empty tree yields nothing.
- [ ] Early stop works when consumer returns `false`.
- [ ] Mutation during iteration is not safe.

Performance expectations
- [ ] Full traversal is `O(n)`.
- [ ] Early-stop traversal is `O(k)` for yielded prefix size `k`.
- [ ] Iterator setup is `O(1)`.

## Internal representation
- [ ] Node IDs are stable, start at `0`, increase monotonically, and are never reused.
- [ ] Use index/ID-based storage only; do not use pointer-linked nodes.
- [ ] Store `parent`, `firstChild`, `nextSibling`, and `prevSibling` indexes.
- [ ] Use `-1` as the nil/sentinel index.

## Auto-resize policy
- [ ] No capacity-based `Grow()` or `Shrink()` API.
- [ ] When node arrays are full, allocate larger arrays and copy node fields by ID.
- [ ] Allocate node storage on `AddChild`.
- [ ] `RemoveSubtree` marks removed nodes dead but does not reuse their public IDs.

## Invariants
- [ ] Root ID is `0` whenever tree is non-empty.
- [ ] Exactly one root exists while tree is non-empty.
- [ ] Non-root nodes have exactly one parent.
- [ ] No cycles.
- [ ] `Len()` equals number of live nodes.
- [ ] `Clone()` and `CloneWith(...)` preserve root ID, live IDs, removed-ID holes, parent-child links, sibling order, and `Len()` in clone.

## Iterator contract
- [ ] `PreOrder()` yields parent before its children.
- [ ] Each live node is yielded exactly once.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty tree yields nothing and does not panic. This includes tree after removing root.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] `RemoveSubtree(0)` removes entire tree and leaves it empty.
- [ ] Invalid node IDs are handled safely.
- [ ] Add/remove around leaves works correctly.
- [ ] `Clone()` of empty tree remains empty and keeps future root-child ID progression identical to original empty tree state.
- [ ] `CloneWith(nil)` is equivalent to `Clone()`.
- [ ] `CloneWith(...)` calls `cloneValue` only for live nodes, never for removed-ID holes.

## Test checklist
- [ ] Parent-child consistency tests.
- [ ] No-cycle verification.
- [ ] Traversal order tests.
- [ ] `Clone/CloneWith` tests for independence, nil-hook equivalence, shallow-copy default, custom hook behavior, live-node-only hook calls, preserved IDs/holes, and preserved next-ID progression.

## Benchmark checklist
- [ ] Benchmark only non-trivial public APIs.
- [ ] Do not benchmark trivial metadata accessor `Len()`.
- [ ] `AddChild` benchmark.
- [ ] `RemoveSubtree` benchmark.
- [ ] `Get` benchmark.
- [ ] `Parent` benchmark.
- [ ] `ChildCount` benchmark.
- [ ] `Clone` benchmark.
- [ ] `CloneWith` benchmark.
- [ ] `PreOrder` benchmark.
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
- [ ] Expected `O(idRange)` APIs must compare normalized `(ns/op)/idRange` and keep normalized `1e5` result within `3.0x` of normalized `1e3` result.

### Layer 2: Absolute timing thresholds
- [ ] Absolute thresholds are per API, per payload class, and per benchmark size.
- [ ] Provisional thresholds should be derived from expected complexity and payload size.
- [ ] Calibrated thresholds should be derived from repeated local runs, median `ns/op`, and safety factor.
- [ ] Suggested safety factor for stable simple tree-general ops is `1.75x`.
- [ ] Suggested safety factor for traversal and subtree workloads is `2.00x`.
- [ ] Suggested safety factor for copy-heavy clone or large-payload workloads is `2.25x`.

### Payload classes
- [ ] Tiny payload benchmark uses scalar-like values such as `int`.
- [ ] Large payload benchmark uses fixed-size struct or fixed-size array payload.
- [ ] Large-payload thresholds must be tracked independently from tiny-payload thresholds.

## Test Generator Hints
- Verify parent-child consistency and acyclic property after subtree edits.
- Cover root removal policy and invalid ID handling.
- Validate `Clone/CloneWith` preserve IDs, holes, pre-order traversal, and next inserted ID.
- Compare traversal output against expected hierarchy snapshots.
- Iterator tests must verify preorder semantics and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for n-ary tree API including add-child, remove-subtree, parent lookup, and invalid ID behavior."
- Clone tests: "Generate tests for `Clone()` and `CloneWith(...)` verifying independent container state, shallow-copy default, nil-hook equivalence, custom hook behavior, live-node-only hook calls, preserved IDs/holes, and preserved next-ID progression."
- Property tests: "Generate randomized tree mutations with fixed seed and assert no-cycle, single-parent, and length invariants after each batch."
- Iterator tests: "Generate tests for `PreOrder() iter.Seq[T]` verifying parent-before-children ordering, exact count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate tree-general benchmarks for add-child, remove-subtree, and full preorder traversal at 1e3/1e4/1e5 nodes."
