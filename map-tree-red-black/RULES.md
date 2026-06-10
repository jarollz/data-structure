# RULES.md - map-tree-red-black

## Goal
Implement an ordered map `K -> V` using a red-black tree.

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics and comparator
- [ ] Use `K` for keys and `V any` for values.
- [ ] Require explicit key comparator.
- [ ] Duplicate keys overwrite values.

## Required API
- [ ] `New(cmp func(a, b K) int) *Map[K, V]`
- [ ] `Put(key K, value V)` inserts or overwrites. Overwrite does not change `Len()`.
- [ ] `Get(key K) (V, bool)`
- [ ] `Delete(key K) bool`
- [ ] `Has(key K) bool`
- [ ] `Min() (K, V, bool)` returns `(zeroK, zeroV, false)` on empty map.
- [ ] `Max() (K, V, bool)` returns `(zeroK, zeroV, false)` on empty map.
- [ ] `Len() int`
- [ ] `Clear()`
- [ ] `All() iter.Seq2[K, V]`

## Internal representation
- [ ] Array-backed nodes with indexes and `-1` as nil/sentinel.
- [ ] Track root index explicitly.
- [ ] Store `left`, `right`, `parent`, and `color` per node.
- [ ] Store keys and values in arrays keyed by node index.
- [ ] Use free-list for node reuse.
- [ ] Apply red-black insert and delete fix-up rules.

## Auto-resize policy
- [ ] No capacity-based `Grow()` or `Shrink()` API.
- [ ] When free-list is empty, allocate larger node-pool arrays and copy node fields by index.
- [ ] Allocate node storage on insert operations.
- [ ] Reclaim node storage on delete/clear operations.
- [ ] Free-list reuse is recommended to reduce allocation churn.

## Invariants
- [ ] BST ordering holds.
- [ ] Root is black.
- [ ] No red node has a red child.
- [ ] Every root-to-leaf path has equal black height.
- [ ] `Len()` equals live node count.

## Iterator contract
- [ ] `All()` yields each key-value pair exactly once.
- [ ] Yield order is ascending key order.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty map yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Insert into empty tree.
- [ ] Delete root in all child-count cases.
- [ ] Delete node with black sibling scenarios.
- [ ] Insert duplicate key updates value only.
- [ ] Delete missing key returns `false`.

## Test checklist
- [ ] API behavior tests.
- [ ] Color-property and black-height property checks.
- [ ] Randomized operation tests versus reference model.
- [ ] In-order traversal is sorted.

## Benchmark checklist
- [ ] `Put/Get/Delete` benchmarks across sizes.
- [ ] Ordered iteration benchmark.
- [ ] Mixed read/write benchmark.

## Test Generator Hints
- Cover insert and delete fix-up branches, especially recolor and rotations.
- Use randomized operations with fixed seed against an ordered oracle model.
- Validate root-black, red-red exclusion, and equal black-height properties.
- Iterator tests must verify ascending key order and early stop behavior.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for `map-tree-red-black` API including overwrite, min/max, delete-root, and empty-structure behavior."
- Property tests: "Generate randomized operations for red-black map with oracle comparison and full red-black invariant checks after each mutation."
- Iterator tests: "Generate tests for `All() iter.Seq2[K,V]` that enforce ascending key order, exact count, and early-stop correctness."
- Benchmarks: "Generate red-black map benchmarks for 1e3/1e4/1e5 with read-heavy, write-heavy, and mixed workloads."
