# RULES.md - map-tree-avl

## Goal
Implement an ordered map `K -> V` on top of an AVL tree.

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics and comparator
- [ ] Use `K` for keys and `V any` for values.
- [ ] Require explicit comparator for keys.
- [ ] Comparator must define strict weak ordering.
- [ ] Duplicate keys are not stored twice; `Put` overwrites value.

## Required API
- [ ] `New(cmp func(a, b K) int) *Map[K, V]`
- [ ] `Put(key K, value V)`
- [ ] `Get(key K) (V, bool)`
- [ ] `Delete(key K) bool`
- [ ] `Has(key K) bool`
- [ ] `Min() (K, V, bool)`
- [ ] `Max() (K, V, bool)`
- [ ] `Len() int`
- [ ] `Clear()`
- [ ] `All() iter.Seq2[K, V]`

## Internal representation
- [ ] Store nodes in arrays by index (no slices/maps).
- [ ] Track `left`, `right`, and `height` (or balance factor) per node.
- [ ] Keep free-list for node reuse.
- [ ] Rebalance with AVL rotations after insert/delete.

## Auto-resize policy
- [ ] No capacity-based `Grow()` or `Shrink()` API.
- [ ] Allocate node storage on insert operations.
- [ ] Reclaim node storage on delete/clear operations.
- [ ] Free-list reuse is recommended to reduce allocation churn.

## Invariants
- [ ] BST ordering holds for every node.
- [ ] For each node, height difference between children is at most 1.
- [ ] Stored height (or balance factor) matches subtree reality.
- [ ] `Len()` equals number of live nodes.

## Iterator contract
- [ ] `All()` yields each key-value pair exactly once.
- [ ] Yield order is ascending key order.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty map yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Delete from empty map.
- [ ] Insert duplicate key updates value.
- [ ] Delete root with 0, 1, and 2 children.
- [ ] Rebalance on LL, RR, LR, RL cases.

## Test checklist
- [ ] API behavior tests.
- [ ] Comparator-driven ordering tests.
- [ ] Rotation case tests.
- [ ] Randomized operations against reference model.
- [ ] Invariant checks after every mutation batch.

## Benchmark checklist
- [ ] `Put/Get/Delete` benchmarks across sizes.
- [ ] Ordered full-iteration benchmark.
- [ ] Range-like workload benchmark (if implemented).

## Test Generator Hints
- Cover LL/RR/LR/RL rebalance paths explicitly.
- Use randomized map operations against a sorted oracle model.
- Validate BST ordering, balance bounds, and metadata consistency after mutations.
- Iterator tests must verify ascending key order and early stop behavior.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for `map-tree-avl` API including duplicate-key overwrite, min/max, and delete-root cases."
- Property tests: "Generate randomized `Put/Get/Delete` tests for AVL map with fixed seed and oracle comparison, plus invariant checks each step."
- Iterator tests: "Generate tests for `All() iter.Seq2[K,V]` verifying strict ascending keys, exact count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate AVL map benchmarks for 1e3/1e4/1e5 keys with read-heavy/write-heavy/mixed workloads and full ordered iteration."
