# RULES.md - tree-red-black

## Goal
Implement a set-like red-black tree for values of type `T`.

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
- [ ] `Insert(v T) bool` inserts new value and returns `true`; return `false` when equal value already exists.
- [ ] `Delete(v T) bool` removes existing value and returns `true`; return `false` when value is missing.
- [ ] `Has(v T) bool`
- [ ] `Min() (T, bool)` returns `(zero, false)` on empty tree.
- [ ] `Max() (T, bool)` returns `(zero, false)` on empty tree.
- [ ] `Len() int`
- [ ] `Clear()`
- [ ] `InOrder() iter.Seq[T]`

## Internal representation
- [ ] Use an index-based node pool; do not use pointer-linked nodes.
- [ ] Use `-1` as the nil/sentinel index and track root index explicitly.
- [ ] Store `left`, `right`, `parent`, and `color` in arrays keyed by node index.
- [ ] Free-list for reusable nodes.
- [ ] Insert and delete fix-up routines maintain color properties.

## Auto-resize policy
- [ ] No capacity-based `Grow()` or `Shrink()` API.
- [ ] When free-list is empty, allocate larger node-pool arrays and copy node fields by index.
- [ ] Allocate node storage on insert operations.
- [ ] Reclaim node storage on delete/clear operations.
- [ ] Free-list reuse is recommended to reduce allocation churn.

## Invariants
- [ ] BST ordering holds.
- [ ] Root is black.
- [ ] Red node cannot have red child.
- [ ] All root-to-leaf paths have equal black height.
- [ ] `Len()` equals live node count.

## Iterator contract
- [ ] `InOrder()` yields ascending sorted values.
- [ ] Each value is yielded exactly once.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty tree yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Insert/delete on empty tree.
- [ ] Duplicate insert returns `false` and leaves tree unchanged.
- [ ] Delete root and near-root nodes.
- [ ] Recolor and rotation branches are fully covered.

## Test checklist
- [ ] Property checks for all red-black rules.
- [ ] Sorted-order traversal tests.
- [ ] Random operations against reference set model.

## Benchmark checklist
- [ ] Insert benchmark.
- [ ] Has benchmark.
- [ ] Delete benchmark.
- [ ] In-order traversal benchmark.

## Test Generator Hints
- Cover recolor and rotation paths for insert and delete fix-up.
- Use randomized set operations with fixed seed and sorted oracle.
- Validate root-black, no-red-red, and equal black-height properties.
- Iterator tests must verify sorted in-order output and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for set-like red-black tree API including duplicate insert policy, delete-root, and min/max behavior."
- Property tests: "Generate randomized insert/delete/has with fixed seed and oracle set comparison, validating all red-black invariants after each mutation."
- Iterator tests: "Generate tests for `InOrder() iter.Seq[T]` enforcing sorted order, exact count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate red-black set benchmarks for insert/has/delete and in-order traversal at 1e3/1e4/1e5 sizes."
