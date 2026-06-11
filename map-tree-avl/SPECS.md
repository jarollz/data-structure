# SPECS.md - map-tree-avl

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
Each public API below is part of contract. Use each checklist as implementation target, not loose guidance.

### `New(cmp func(a, b K) int) *Map[K, V]`
Purpose
- [ ] Create empty ordered map with AVL-backed key ordering.

Behavior expectations
- [ ] Returned map is non-nil and empty.
- [ ] Comparator becomes persistent ordering rule for all future operations.
- [ ] `Len()` is `0` immediately after construction.

Performance expectations
- [ ] `O(1)` time.

### `Put(key K, value V)`
Purpose
- [ ] Insert new key-value entry or overwrite existing key.

Behavior expectations
- [ ] Missing key inserts new entry.
- [ ] Existing key overwrites stored value and does not change `Len()`.
- [ ] Key order remains valid under comparator.
- [ ] AVL balance and metadata remain valid after operation.

Performance expectations
- [ ] `O(log n)` time.

### `Get(key K) (V, bool)`
Purpose
- [ ] Read value currently stored for key.

Behavior expectations
- [ ] Existing key returns `(value, true)`.
- [ ] Missing key returns `(zero, false)`.
- [ ] `Get` never mutates map state.

Performance expectations
- [ ] `O(log n)` time.

### `Delete(key K) bool`
Purpose
- [ ] Remove live entry for key.

Behavior expectations
- [ ] Existing key returns `true` and is removed.
- [ ] Missing key returns `false`.
- [ ] Delete root with `0`, `1`, or `2` children works correctly.
- [ ] Comparator order and AVL invariants remain valid after deletion.

Performance expectations
- [ ] `O(log n)` time.

### `Has(key K) bool`
Purpose
- [ ] Report whether key currently exists.

Behavior expectations
- [ ] Existing key returns `true`.
- [ ] Missing key returns `false`.
- [ ] `Has` never mutates map state.

Performance expectations
- [ ] `O(log n)` time.

### `Min() (K, V, bool)`
Purpose
- [ ] Read smallest key and its value.

Behavior expectations
- [ ] Empty map returns `(zeroK, zeroV, false)`.
- [ ] Non-empty map returns leftmost key-value pair under comparator ordering.
- [ ] `Min()` never mutates map state.

Performance expectations
- [ ] `O(log n)` time.

### `Max() (K, V, bool)`
Purpose
- [ ] Read largest key and its value.

Behavior expectations
- [ ] Empty map returns `(zeroK, zeroV, false)`.
- [ ] Non-empty map returns rightmost key-value pair under comparator ordering.
- [ ] `Max()` never mutates map state.

Performance expectations
- [ ] `O(log n)` time.

### `Len() int`
Purpose
- [ ] Report number of live entries.

Behavior expectations
- [ ] Count increases only for new keys.
- [ ] Count does not change on overwrite.
- [ ] Count decreases only for successful deletes.
- [ ] Count becomes `0` after `Clear`.

Performance expectations
- [ ] `O(1)` time.

### `Clear()`
Purpose
- [ ] Reset map to empty logical state.

Behavior expectations
- [ ] `Len()` becomes `0`.
- [ ] Root resets to empty-state sentinel.
- [ ] Future puts, deletes, and queries still work correctly.
- [ ] Clear on empty map is safe.

Performance expectations
- [ ] Target `O(1)` logical reset.
- [ ] `O(n)` node reclamation bookkeeping is allowed if implementation needs it.

### `Clone() *Map[K, V]`
Purpose
- [ ] Create independent ordered map with same sorted entries and comparator.

Behavior expectations
- [ ] Clone preserves `Len()`, comparator, ascending key order, and lookup results.
- [ ] Keys and values are copied with normal Go assignment.
- [ ] Clone is container-independent from source.
- [ ] Clone remains AVL-valid.

Performance expectations
- [ ] `O(n)` time for `n = Len()`.
- [ ] `O(n)` extra storage.

### `CloneWith(cloneKey func(K) K, cloneValue func(V) V) *Map[K, V]`
Purpose
- [ ] Create independent ordered map while optionally transforming each live key and value.

Behavior expectations
- [ ] Preserves `Len()`, comparator, ascending key order, and lookup results under transformed keys.
- [ ] Nil hook for one side means normal Go assignment for that payload type.
- [ ] Non-nil hooks are called exactly once per live entry.
- [ ] Hook call order is ascending key order.
- [ ] Cloned keys must remain comparator-compatible.

Performance expectations
- [ ] `O(n)` container work plus hook cost.

### `All() iter.Seq2[K, V]`
Purpose
- [ ] Iterate all live entries in ascending key order.

Behavior expectations
- [ ] Yield each live key-value pair exactly once.
- [ ] Empty map yields nothing.
- [ ] Early stop works when consumer returns `false`.
- [ ] Mutation during iteration is not safe.

Performance expectations
- [ ] Full traversal is `O(n)`.
- [ ] Early-stop traversal is `O(k)` for yielded prefix size `k`.
- [ ] Iterator setup is `O(1)` or `O(height)` if traversal stack is built lazily.

## Internal representation
- [ ] Store nodes in arrays by index (no slices/maps).
- [ ] Use `-1` as the nil/sentinel index and track root index explicitly.
- [ ] Track `left`, `right`, and `height` (or balance factor) per node.
- [ ] Store keys and values in arrays keyed by node index.
- [ ] Keep free-list for node reuse.
- [ ] Rebalance with AVL rotations after insert/delete.

## Auto-resize policy
- [ ] No capacity-based `Grow()` or `Shrink()` API.
- [ ] When free-list is empty, allocate larger node-pool arrays and copy node fields by index.
- [ ] Allocate node storage on insert operations.
- [ ] Reclaim node storage on delete/clear operations.
- [ ] Free-list reuse is recommended to reduce allocation churn.

## Invariants
- [ ] BST ordering holds for every node.
- [ ] For each node, height difference between children is at most 1.
- [ ] Stored height (or balance factor) matches subtree reality.
- [ ] `Len()` equals number of live nodes.
- [ ] `Clone()` and `CloneWith(...)` preserve `Len()`, comparator, ascending key order, and AVL invariants in clone.

## Iterator contract
- [ ] `All()` yields each key-value pair exactly once.
- [ ] Yield order is ascending key order.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty map yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Delete from empty map.
- [ ] Insert duplicate key updates value.
- [ ] Delete missing key returns `false`.
- [ ] Delete root with 0, 1, and 2 children.
- [ ] Rebalance on LL, RR, LR, RL cases.
- [ ] `Clone()` on empty map returns empty independent map with same comparator.
- [ ] `CloneWith(nil, nil)` is equivalent to `Clone()`.
- [ ] `CloneWith(...)` calls hooks only for live entries in ascending key order.
- [ ] `cloneKey` must preserve comparator ordering contract for cloned keys.

## Test checklist
- [ ] API behavior tests.
- [ ] Comparator-driven ordering tests.
- [ ] Rotation case tests.
- [ ] Randomized operations against reference model.
- [ ] Invariant checks after every mutation batch.
- [ ] `Clone/CloneWith` tests for independence, nil-hook equivalence, shallow-copy default, custom key/value hook behavior, live-entry-only hook calls, and preserved sorted/min/max behavior.

## Benchmark checklist
- [ ] Benchmark only non-trivial public APIs.
- [ ] Do not benchmark trivial metadata accessor `Len()`.
- [ ] `Put` benchmark.
- [ ] `Get` benchmark.
- [ ] `Delete` benchmark.
- [ ] `Has` benchmark.
- [ ] `Min` benchmark.
- [ ] `Max` benchmark.
- [ ] `Clear` benchmark.
- [ ] `Clone` benchmark.
- [ ] `CloneWith` benchmark.
- [ ] `All` benchmark.
- [ ] Mixed read/write benchmark.
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
- [ ] Mixed workloads are judged by dominant documented complexity of included operations.

### Layer 2: Absolute timing thresholds
- [ ] Absolute thresholds are per API, per payload class, and per benchmark size.
- [ ] Provisional thresholds should be derived from expected complexity and payload size.
- [ ] Calibrated thresholds should be derived from repeated local runs, median `ns/op`, and safety factor.
- [ ] Suggested safety factor for stable simple AVL-map ops is `1.75x`.
- [ ] Suggested safety factor for traversal and mixed workloads is `2.00x`.
- [ ] Suggested safety factor for copy-heavy clone or large-payload workloads is `2.25x`.

### Payload classes
- [ ] Tiny payload benchmark uses scalar-like keys and values such as `int`.
- [ ] Large payload benchmark uses fixed-size struct or fixed-size array payload.
- [ ] Large-payload thresholds must be tracked independently from tiny-payload thresholds.

## Test Generator Hints
- Cover LL/RR/LR/RL rebalance paths explicitly.
- Validate `Clone/CloneWith` preserve ascending key order, min/max results, and hook-call count.
- Use randomized map operations against a sorted oracle model.
- Validate BST ordering, balance bounds, and metadata consistency after mutations.
- Iterator tests must verify ascending key order and early stop behavior.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for `map-tree-avl` API including duplicate-key overwrite, min/max, and delete-root cases."
- Clone tests: "Generate tests for `Clone()` and `CloneWith(...)` verifying independent container state, shallow-copy default, nil-hook equivalence, custom key/value hook behavior, live-entry-only hook calls, and preserved sorted/min/max behavior."
- Property tests: "Generate randomized `Put/Get/Delete` tests for AVL map with fixed seed and oracle comparison, plus invariant checks each step."
- Iterator tests: "Generate tests for `All() iter.Seq2[K,V]` verifying strict ascending keys, exact count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate AVL map benchmarks for 1e3/1e4/1e5 keys with read-heavy/write-heavy/mixed workloads and full ordered iteration."
