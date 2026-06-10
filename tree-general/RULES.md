# RULES.md - tree-general

## Goal
Implement a generic n-ary tree for hierarchical data.

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics
- [ ] Use `T any` for node value.

## Required API
- [ ] `New(rootValue T) *Tree[T]` creates a tree with one root node. Root node ID is always `0`.
- [ ] `AddChild(parentID int, value T) (childID int, ok bool)` adds new last child of `parentID`. Return `(-1, false)` when `parentID` is invalid or tree is empty.
- [ ] `RemoveSubtree(nodeID int) bool` removes `nodeID` and all descendants. Return `false` for invalid IDs. If `nodeID == 0`, whole tree becomes empty.
- [ ] `Get(nodeID int) (T, bool)` returns `(zero, false)` for invalid or removed IDs.
- [ ] `Parent(nodeID int) (int, bool)` returns `(-1, false)` for root, invalid, or removed IDs.
- [ ] `ChildCount(nodeID int) int` returns number of live direct children, or `-1` for invalid or removed IDs.
- [ ] `Len() int`
- [ ] `PreOrder() iter.Seq[T]`

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

## Test checklist
- [ ] Parent-child consistency tests.
- [ ] No-cycle verification.
- [ ] Traversal order tests.

## Benchmark checklist
- [ ] Add-child benchmark.
- [ ] Remove-subtree benchmark.
- [ ] Full traversal benchmark.

## Test Generator Hints
- Verify parent-child consistency and acyclic property after subtree edits.
- Cover root removal policy and invalid ID handling.
- Compare traversal output against expected hierarchy snapshots.
- Iterator tests must verify preorder semantics and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for n-ary tree API including add-child, remove-subtree, parent lookup, and invalid ID behavior."
- Property tests: "Generate randomized tree mutations with fixed seed and assert no-cycle, single-parent, and length invariants after each batch."
- Iterator tests: "Generate tests for `PreOrder() iter.Seq[T]` verifying parent-before-children ordering, exact count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate tree-general benchmarks for add-child, remove-subtree, and full preorder traversal at 1e3/1e4/1e5 nodes."
