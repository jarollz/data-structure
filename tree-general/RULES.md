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
- [ ] `New(rootValue T) *Tree[T]`
- [ ] `AddChild(parentID int, value T) (childID int, ok bool)`
- [ ] `RemoveSubtree(nodeID int) bool`
- [ ] `Get(nodeID int) (T, bool)`
- [ ] `Parent(nodeID int) (int, bool)`
- [ ] `ChildCount(nodeID int) int`
- [ ] `Len() int`
- [ ] `PreOrder() iter.Seq[T]`

## Internal representation
- [ ] Define stable node identity (index or ID).
- [ ] Store parent reference/index.
- [ ] Store child links with array-only strategy.

## Auto-resize policy
- [ ] No capacity-based `Grow()` or `Shrink()` API.
- [ ] Allocate node storage on `AddChild`.
- [ ] Reclaim node storage on `RemoveSubtree`.
- [ ] Optional free-list reuse is allowed if tree invariants remain valid.

## Invariants
- [ ] Exactly one root exists while tree is non-empty.
- [ ] Non-root nodes have exactly one parent.
- [ ] No cycles.
- [ ] `Len()` equals number of live nodes.

## Iterator contract
- [ ] `PreOrder()` yields parent before its children.
- [ ] Each live node is yielded exactly once.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty tree yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Remove root behavior is explicitly defined.
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
