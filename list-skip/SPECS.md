# SPECS.md - list-skip

## Goal
Implement a generic ordered skip list.

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics and ordering
- [ ] Use `T any`.
- [ ] Require comparator for ordering.
- [ ] Duplicate policy: do not store duplicates. `Insert(v)` returns `false` when equal value already exists.

## Required API
Each public API below is part of contract. Use each checklist as implementation target, not loose guidance.

### `New(maxLevel int, cmp func(a, b T) int) *SkipList[T]`
Purpose
- [ ] Create empty ordered skip list with deterministic level generator.

Behavior expectations
- [ ] Normalize `maxLevel < 1` to `1`.
- [ ] Effective configured level count is `levelCap = max(1, maxLevel)`.
- [ ] Returned skip list is non-nil and empty.
- [ ] Comparator becomes persistent ordering rule for all future operations.
- [ ] Deterministic RNG state starts from documented initial value.
- [ ] `Len()` is `0` immediately after construction.

Performance expectations
- [ ] `O(levelCap)` time.
- [ ] `O(levelCap)` initial level-head storage.

### `Insert(v T) bool`
Purpose
- [ ] Insert new value into sorted set-like structure.

Behavior expectations
- [ ] Return `true` when value was not previously present.
- [ ] Return `false` on duplicate value and leave structure unchanged.
- [ ] Level `0` remains fully sorted after insertion.
- [ ] Higher-level links remain sorted after insertion.
- [ ] Deterministic RNG state advances exactly once per attempted insertion that reaches level generation.

Performance expectations
- [ ] Expected `O(log n)` time when `levelCap` is configured to scale with population size.
- [ ] If `levelCap` is too small for current `n`, performance degrades toward `O(n)`.
- [ ] Worst-case single call may include node-pool growth copy.

### `Delete(v T) bool`
Purpose
- [ ] Remove existing value from skip list.

Behavior expectations
- [ ] Return `true` when value existed and was removed.
- [ ] Return `false` when value is missing.
- [ ] Sorted order remains valid after deletion.
- [ ] `currentLevel` shrinks when highest levels become empty.

Performance expectations
- [ ] Expected `O(log n)` time when `levelCap` is configured to scale with population size.
- [ ] If `levelCap` is too small for current `n`, performance degrades toward `O(n)`.

### `Has(v T) bool`
Purpose
- [ ] Report whether value currently exists.

Behavior expectations
- [ ] Return `true` only for live stored values.
- [ ] Duplicate policy means at most one matching value can exist.
- [ ] Empty list safely returns `false`.

Performance expectations
- [ ] Expected `O(log n)` time when `levelCap` is configured to scale with population size.
- [ ] If `levelCap` is too small for current `n`, performance degrades toward `O(n)`.

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
- [ ] Reset skip list to empty logical state.

Behavior expectations
- [ ] `Len()` becomes `0`.
- [ ] Head sentinel remains valid for future operations.
- [ ] `currentLevel` resets to empty-state value.
- [ ] Future insert, delete, and has operations still work correctly.

Performance expectations
- [ ] Target `O(1)` logical reset.
- [ ] `O(n)` node reclamation bookkeeping is allowed if implementation needs it.

### `Clone() *SkipList[T]`
Purpose
- [ ] Create independent skip list with same sorted contents and deterministic future behavior.

Behavior expectations
- [ ] Clone preserves `Len()`, sorted order, comparator, `maxLevel`, `currentLevel`, and deterministic RNG state.
- [ ] Elements are copied with normal Go assignment.
- [ ] Clone is container-independent from source.
- [ ] Future insert sequence on original and clone should choose same levels when operation sequence stays same.

Performance expectations
- [ ] `O(levelCap * n)` time for `n = Len()`.
- [ ] `O(levelCap * n)` extra storage.

### `CloneWith(cloneValue func(T) T) *SkipList[T]`
Purpose
- [ ] Create independent skip list while optionally transforming each live value.

Behavior expectations
- [ ] Preserves `Len()`, sorted order, comparator, `maxLevel`, `currentLevel`, and deterministic RNG state.
- [ ] Nil hook is equivalent to `Clone()`.
- [ ] Non-nil hook is called exactly once per live value.
- [ ] Hook call order is sorted order from level `0` traversal.
- [ ] Hook is never called for reclaimed or free-list nodes.

Performance expectations
- [ ] `O(levelCap * n)` container work plus hook cost.

### `Values() iter.Seq[T]`
Purpose
- [ ] Iterate live values in sorted order using level `0`.

Behavior expectations
- [ ] Yield each live value exactly once.
- [ ] Empty skip list yields nothing.
- [ ] Early stop works when consumer returns `false`.
- [ ] Mutation during iteration is not safe.

Performance expectations
- [ ] Full traversal is `O(n)`.
- [ ] Early-stop traversal is `O(k)` for yielded prefix size `k`.
- [ ] Iterator setup is `O(1)`.

## Internal representation
- [ ] Use an index-based node pool with `-1` as nil/sentinel index.
- [ ] Keep head sentinel node with forward references for all levels `0..maxLevel-1`.
- [ ] Node stores value and forward references for each level.
- [ ] Level references use arrays/indexes only.
- [ ] Track `currentLevel`, length, free-list head, and deterministic RNG state.
- [ ] Level generation policy is fixed: initialize internal xorshift32 state to `1` in `New`; on each insert, advance state and promote while low bit is `1`, stopping at first `0` or `maxLevel`. Same operation sequence must produce same levels.

## Auto-resize policy
- [ ] No capacity-based `Grow()` or `Shrink()` API.
- [ ] When free-list is empty, allocate larger node-pool arrays and copy node/link fields by index.
- [ ] Allocate node storage on insert operations.
- [ ] Reclaim node storage on delete/clear operations.
- [ ] Optional free-list reuse is allowed if skip-list invariants remain valid.
- [ ] After deletions, reduce `currentLevel` while highest level has no live nodes.

## Invariants
- [ ] Level 0 contains all live elements in sorted order.
- [ ] Higher-level links preserve sorted order.
- [ ] Search path from top level reaches correct target/miss.
- [ ] `Len()` equals number of live elements.
- [ ] `Clone()` and `CloneWith(...)` preserve sorted order, length, comparator, `maxLevel`, `currentLevel`, and deterministic future RNG behavior.

## Iterator contract
- [ ] `Values()` yields sorted order using level 0 traversal.
- [ ] Each element is yielded exactly once.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty list yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Empty list operations are safe.
- [ ] Duplicate insert behavior is consistent with contract.
- [ ] Deleting missing value is safe.
- [ ] `Clone()` on empty list returns empty independent list with same `maxLevel`, comparator, `currentLevel`, and deterministic RNG state.
- [ ] `CloneWith(nil)` is equivalent to `Clone()`.
- [ ] `CloneWith(...)` calls `cloneValue` only for live elements on level `0`, never for reclaimed or free-list nodes.

## Test checklist
- [ ] Sorted-order verification tests.
- [ ] Random insert/delete/has tests.
- [ ] Level-link consistency checks.
- [ ] `Clone/CloneWith` tests for independence, nil-hook equivalence, shallow-copy default, custom hook behavior, live-node-only hook calls, and preserved deterministic future behavior.

## Benchmark checklist
- [ ] Benchmark only non-trivial public APIs.
- [ ] Do not benchmark trivial metadata accessor `Len()`.
- [ ] `Insert` benchmark.
- [ ] `Delete` benchmark.
- [ ] `Has` benchmark.
- [ ] `Clear` benchmark.
- [ ] `Clone` benchmark.
- [ ] `CloneWith` benchmark.
- [ ] `Values` benchmark.
- [ ] Mixed ordered workload benchmark.
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
- [ ] Expected `O(log n)` APIs must keep `ns/op(1e5) <= 2.5 * ns/op(1e3)` when `levelCap` is scaled appropriately for each benchmark size.
- [ ] Expected `O(n)` APIs must compare normalized `(ns/op)/n` and keep normalized `1e5` result within `3.0x` of normalized `1e3` result.
- [ ] Expected `O(levelCap * n)` APIs must compare normalized `(ns/op)/(levelCap*n)` and keep normalized `1e5` result within `3.0x` of normalized `1e3` result.
- [ ] Mixed workloads are judged by dominant documented complexity of included operations.

### Layer 2: Absolute timing thresholds
- [ ] Absolute thresholds are per API, per payload class, and per benchmark size.
- [ ] Provisional thresholds should be derived from expected complexity and payload size.
- [ ] Calibrated thresholds should be derived from repeated local runs, median `ns/op`, and safety factor.
- [ ] Suggested safety factor for stable simple skip-list ops is `1.75x`.
- [ ] Suggested safety factor for iteration and mixed workloads is `2.00x`.
- [ ] Suggested safety factor for copy-heavy clone or large-payload workloads is `2.25x`.

### Payload classes
- [ ] Tiny payload benchmark uses scalar-like elements such as `int`.
- [ ] Large payload benchmark uses fixed-size struct or fixed-size array payload.
- [ ] Large-payload thresholds must be tracked independently from tiny-payload thresholds.

## Test Generator Hints
- Validate sorted order on level 0 after every mutation batch.
- Cover duplicate-policy behavior and delete-missing behavior.
- Validate `Clone/CloneWith` preserve sorted order, RNG state, and hook-call count.
- Use randomized insert/delete/has with fixed seed and sorted oracle.
- Iterator tests must verify sorted output and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for skip list API covering insert/delete/has semantics and duplicate policy."
- Clone tests: "Generate tests for `Clone()` and `CloneWith(...)` verifying independent container state, shallow-copy default, nil-hook equivalence, custom hook behavior, live-node-only hook calls, and preserved deterministic future behavior."
- Property tests: "Generate randomized ordered operations with fixed seed and compare with sorted oracle while checking skip-list invariants."
- Iterator tests: "Generate tests for `Values() iter.Seq[T]` enforcing sorted order, exact count, early stop, and mutation-unsafety note."
- Benchmarks: "Generate skip-list benchmarks for search, insert, delete, and mixed workloads at 1e3/1e4/1e5 sizes."
