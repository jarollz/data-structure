# RULES.md - tree-avl

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
- [ ] `New(cmp func(a, b T) int) *Tree[T]`
- [ ] `Insert(v T) bool`
- [ ] `Delete(v T) bool`
- [ ] `Has(v T) bool`
- [ ] `Min() (T, bool)`
- [ ] `Max() (T, bool)`
- [ ] `Len() int`
- [ ] `Clear()`
- [ ] `InOrder() iter.Seq[T]`

## Internal representation
- [ ] Node arrays with child references and height/balance data.
- [ ] Free-list for reusable node slots.
- [ ] AVL rotations for rebalance.

## Auto-resize policy
- [ ] No capacity-based `Grow()` or `Shrink()` API.
- [ ] Allocate node storage on insert operations.
- [ ] Reclaim node storage on delete/clear operations.
- [ ] Free-list reuse is recommended to reduce allocation churn.

## Invariants
- [ ] BST ordering holds.
- [ ] Balance factor each node is in `[-1, 1]`.
- [ ] Height metadata is correct.
- [ ] `Len()` equals live node count.

## Iterator contract
- [ ] `InOrder()` yields ascending sorted values.
- [ ] Each value is yielded exactly once.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty tree yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Insert/delete on empty tree.
- [ ] Delete root with 0, 1, 2 children.
- [ ] All four rotation patterns are covered.

## Test checklist
- [ ] Sorted order tests.
- [ ] Rotation scenario tests.
- [ ] Random operations against reference set model.

## Benchmark checklist
- [ ] Insert benchmark.
- [ ] Has benchmark.
- [ ] Delete benchmark.
- [ ] Full in-order traversal benchmark.

## Test Generator Hints
- Cover all AVL rebalancing patterns with deterministic cases.
- Use randomized set operations with fixed seed and sorted oracle.
- Validate BST ordering, height metadata, and balance factor bounds.
- Iterator tests must verify sorted in-order output and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for set-like AVL API including duplicate inserts, delete-root cases, and min/max behavior."
- Property tests: "Generate randomized insert/delete/has tests with fixed seed and oracle set comparison plus AVL invariant checks each step."
- Iterator tests: "Generate tests for `InOrder() iter.Seq[T]` enforcing sorted order, exact count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate AVL set benchmarks for insert/has/delete and full in-order traversal at 1e3/1e4/1e5 sizes."
