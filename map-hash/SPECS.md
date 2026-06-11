# SPECS.md - map-hash

## Goal
Implement a generic hash map with:
- Open addressing
- Linear probing
- Tombstones for delete
- Array-backed storage only

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.
- [ ] Do not use helper containers that replace core hashmap logic.

## Generics and key contract
- [ ] Use `K any` for keys and `V any` for values.
- [ ] Key contract is explicit hooks passed to constructor: `hash func(K) uint64` and `equal func(a, b K) bool`.
- [ ] `hash` and `equal` must both be non-nil; `New` may panic if either is nil.
- [ ] Equal keys must hash consistently.
- [ ] `Put` on an existing key overwrites value.

## Required API
Each public API below is part of contract. Use each checklist as implementation target, not loose guidance.

### `New(capacity int, hash func(K) uint64, equal func(a, b K) bool) *Map[K, V]`
Purpose
- [ ] Create empty hash map with explicit hash and equality hooks.

Behavior expectations
- [ ] Normalize `capacity <= 0` to `16`.
- [ ] Effective starting capacity is `startCap = max(16, capacity)`.
- [ ] `hash` and `equal` must both be non-nil.
- [ ] Returned map is non-nil and empty.
- [ ] `Len()` is `0` immediately after construction.
- [ ] `Cap()` is normalized starting capacity immediately after construction.

Performance expectations
- [ ] `O(startCap)` time.
- [ ] `O(startCap)` backing storage.

### `Put(key K, value V)`
Purpose
- [ ] Insert new entry or overwrite existing entry.

Behavior expectations
- [ ] Missing key is inserted.
- [ ] Existing key has its value overwritten in place logically.
- [ ] Overwrite does not change `Len()`.
- [ ] Probe chain remains correct across collisions and tombstones.
- [ ] If resize or cleanup rehash is needed, all live keys remain retrievable afterward.

Performance expectations
- [ ] Expected `O(1)` average time under documented load policy.
- [ ] Worst-case single call may include `O(n)` rehash.

### `Get(key K) (V, bool)`
Purpose
- [ ] Read value currently stored for key.

Behavior expectations
- [ ] Existing key returns `(value, true)`.
- [ ] Missing key returns `(zero, false)`.
- [ ] Tombstones do not break lookup reachability for later keys in same probe chain.
- [ ] `Get` never mutates map state.

Performance expectations
- [ ] Expected `O(1)` average time.

### `Delete(key K) bool`
Purpose
- [ ] Remove live entry for key.

Behavior expectations
- [ ] Existing key returns `true` and becomes unreachable by `Get` or `Has`.
- [ ] Missing key returns `false`.
- [ ] Delete creates tombstone-safe state until any later cleanup rehash.
- [ ] Keys behind deleted slot remain retrievable.

Performance expectations
- [ ] Expected `O(1)` average time.
- [ ] Worst-case single call may include `O(n)` rehash when cleanup or shrink triggers.

### `Has(key K) bool`
Purpose
- [ ] Report whether key currently exists.

Behavior expectations
- [ ] Returns `true` only for live entries.
- [ ] Deleted and never-inserted keys return `false`.
- [ ] `Has` follows same reachability rules as `Get`.

Performance expectations
- [ ] Expected `O(1)` average time.

### `Len() int`
Purpose
- [ ] Report number of live entries.

Behavior expectations
- [ ] Count increases only for inserts of previously missing keys.
- [ ] Count does not change on overwrite.
- [ ] Count decreases only for successful deletes.
- [ ] Count becomes `0` after `Clear`.

Performance expectations
- [ ] `O(1)` time.

### `Cap() int`
Purpose
- [ ] Report current table capacity.

Behavior expectations
- [ ] Capacity is slot count, not live entry count.
- [ ] Capacity never drops below normalized starting capacity.
- [ ] Capacity reflects grow, shrink, and cleanup rehash decisions.

Performance expectations
- [ ] `O(1)` time.

### `Clear()`
Purpose
- [ ] Reset map to empty state without stale live entries or tombstones.

Behavior expectations
- [ ] `Len()` becomes `0`.
- [ ] All old keys become unreachable.
- [ ] Tombstones are removed.
- [ ] `Cap()` becomes `minCap`.
- [ ] Future put/get/delete operations still work correctly.

Performance expectations
- [ ] Target `O(capacity)` or `O(minCap)` rebuild cost is acceptable.

### `LoadFactor() float64`
Purpose
- [ ] Report live-entry load factor.

Behavior expectations
- [ ] Return `float64(Len()) / float64(Cap())`.
- [ ] Tombstones are excluded from numerator.
- [ ] Empty live table returns `0` when `Len()` is `0`.
- [ ] `LoadFactor()` never mutates map state.

Performance expectations
- [ ] `O(1)` time.

### `Clone() *Map[K, V]`
Purpose
- [ ] Create independent hash map with same live lookup results and operational hooks.

Behavior expectations
- [ ] Clone preserves live entries, `Len()`, `Cap()`, `LoadFactor()`, `hash`, and `equal`.
- [ ] Keys and values are copied with normal Go assignment.
- [ ] Clone is container-independent from source.
- [ ] Tombstone-only slots are not treated as live entries.

Performance expectations
- [ ] `O(Cap())` time.
- [ ] `O(Cap())` extra storage.

### `CloneWith(cloneKey func(K) K, cloneValue func(V) V) *Map[K, V]`
Purpose
- [ ] Create independent hash map while optionally transforming live keys and values.

Behavior expectations
- [ ] Preserves live entries, `Len()`, `Cap()`, `LoadFactor()`, `hash`, and `equal`.
- [ ] Nil hook for one side means normal Go assignment for that payload type.
- [ ] Non-nil hooks are called exactly once per live entry.
- [ ] Hooks are never called for empty or tombstone slots.
- [ ] Cloned keys must remain compatible with `hash` and `equal`.

Performance expectations
- [ ] `O(Cap())` container work plus hook cost.

### `All() iter.Seq2[K, V]`
Purpose
- [ ] Iterate all live key-value entries exactly once.

Behavior expectations
- [ ] Yields each live entry exactly once.
- [ ] Empty map yields nothing.
- [ ] Iteration order is unspecified and may change after mutations or resize.
- [ ] Early stop works when consumer returns `false`.
- [ ] Mutation during iteration is not safe.

Performance expectations
- [ ] Full traversal is `O(capacity)` or `O(n + tombstones)` depending on scan strategy.
- [ ] Early-stop traversal is proportional to visited slots.

## Internal representation
- [ ] Main table is an array of entry slots.
- [ ] Slot state is explicit: empty, occupied, deleted.
- [ ] Collision policy is linear probing only.
- [ ] Resize by allocating a larger array and rehashing occupied entries.
- [ ] Tombstones must not break probe-chain correctness.
- [ ] `Put` should reuse first tombstone found during probe if key is not found later in same probe chain.

## Auto-resize policy
- [ ] Grow by rehashing to larger table when probe occupancy `(live + tombstones) / Cap() >= 0.70`.
- [ ] Shrink by rehashing to smaller table when `LoadFactor() <= 0.15` and `Cap() > minCap`.
- [ ] Cleanup rehash at same capacity when tombstones are high (for example `tombstones > live/2`).
- [ ] Maintain tombstone-safe probing after any grow/shrink/cleanup rehash.
- [ ] `minCap` is normalized starting capacity.
- [ ] Resizing is internal policy; do not expose public `Grow()` or `Shrink()` APIs.
- [ ] Use hysteresis to avoid resize thrash around thresholds.

## Invariants
- [ ] `Len()` equals number of occupied slots.
- [ ] No duplicate keys among occupied slots.
- [ ] Existing keys are always retrievable.
- [ ] Deleted keys are not retrievable.
- [ ] Keys behind tombstones remain retrievable.
- [ ] After resize, all live entries are still retrievable.
- [ ] `Clone()` and `CloneWith(...)` preserve lookup results for all live keys, plus `Len()`, `Cap()`, `LoadFactor()`, `hash`, and `equal` in clone.

## Iterator contract
- [ ] `All()` yields each live key-value pair exactly once.
- [ ] Iteration order is unspecified and can change after inserts/deletes/resizes.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty map yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Empty map operations do not panic.
- [ ] Repeated `Put` on same key does not increase length.
- [ ] Delete missing key returns `false`.
- [ ] Delete then reinsert same key works correctly.
- [ ] Small initial capacity still works.
- [ ] `Clear()` removes all live entries and tombstones and keeps `Cap() == minCap`.
- [ ] `Clone()` on empty map returns empty independent map with same `Cap()`, `hash`, and `equal`.
- [ ] `CloneWith(nil, nil)` is equivalent to `Clone()`.
- [ ] `CloneWith(...)` calls hooks only for live entries, never for empty or tombstone slots.
- [ ] `cloneKey` must preserve `hash` and `equal` compatibility for cloned keys.

## Test checklist
- [ ] Unit tests for every API method.
- [ ] Collision tests for multiple keys into same probe region.
- [ ] Tombstone tests with middle-of-chain deletes.
- [ ] Resize tests that preserve all live entries.
- [ ] Random operation tests against a reference model.
- [ ] Invariant check after operation sequences.
- [ ] `Clone/CloneWith` tests for independence, nil-hook equivalence, shallow-copy default, custom key/value hook behavior, preserved capacity/load factor, and live-entry-only hook calls.

## Benchmark checklist
- [ ] Benchmark only non-trivial public APIs.
- [ ] Do not benchmark trivial metadata accessors such as `Len()`, `Cap()`, and `LoadFactor()`.
- [ ] `Put` benchmark.
- [ ] `Get` hit benchmark.
- [ ] `Get` miss benchmark.
- [ ] `Delete` benchmark.
- [ ] `Has` benchmark.
- [ ] `Clear` benchmark.
- [ ] `Clone` benchmark.
- [ ] `CloneWith` benchmark.
- [ ] `All` benchmark.
- [ ] Mixed `Put/Get/Delete` benchmark.
- [ ] Run benchmarks at `1e3`, `1e4`, and `1e5` sizes.
- [ ] Use tiny payload benchmarks for structure overhead.
- [ ] Use large payload benchmarks for copy-sensitive APIs.
- [ ] Record load factor during benchmark runs.

## Benchmark Validation Policy
- [ ] Use two validation layers for every benchmarked non-trivial API.
- [ ] Layer 1 validates complexity growth.
- [ ] Layer 2 validates absolute timing budget.
- [ ] Before valid implementation exists, absolute thresholds are provisional engineering targets, not calibrated facts.
- [ ] Implementer should aim to get as close as practical to provisional targets while keeping behavior correct.
- [ ] After first correct implementation exists, rerun calibration on target machine and replace provisional targets with measured median-based thresholds.

### Layer 1: Complexity-growth thresholds
- [ ] Expected `O(1)` average APIs must keep `ns/op(1e5) <= 3.0 * ns/op(1e3)` under comparable load-factor setup.
- [ ] Expected `O(n)` iteration-style APIs must compare normalized `(ns/op)/n` and keep normalized `1e5` result within `3.0x` of normalized `1e3` result.
- [ ] Mixed workloads are judged by dominant documented complexity of included operations.

### Layer 2: Absolute timing thresholds
- [ ] Absolute thresholds are per API, per payload class, and per benchmark size.
- [ ] Provisional thresholds should be derived from expected complexity, payload size, and target load factor.
- [ ] Calibrated thresholds should be derived from repeated local runs, median `ns/op`, and safety factor.
- [ ] Suggested safety factor for stable simple hash-map ops is `1.75x`.
- [ ] Suggested safety factor for iteration and mixed workloads is `2.00x`.
- [ ] Suggested safety factor for copy-heavy clone or large-payload workloads is `2.25x`.

### Payload classes
- [ ] Tiny payload benchmark uses scalar-like keys and values such as `int`.
- [ ] Large payload benchmark uses fixed-size struct or fixed-size array payload.
- [ ] Large-payload thresholds must be tracked independently from tiny-payload thresholds.

## Test Generator Hints
- Focus collision-heavy cases, tombstones, and resize/rehash correctness.
- Validate `Clone/CloneWith` preserve `Cap()`, `LoadFactor()`, and tombstone-safe lookup while only cloning live entries.
- Run randomized `Put/Get/Delete` with fixed seed and compare against test oracle.
- Validate invariants after each mutation batch, especially probe-chain reachability.
- Iterator tests must check count, early stop, and unspecified iteration order.

## AI Prompt Snippets
- Unit tests: "Generate table-driven Go tests for `map-hash` API from this SPECS file, including collisions, overwrite, delete-missing, and resize behavior."
- Clone tests: "Generate tests for `Clone()` and `CloneWith(...)` verifying independent container state, shallow-copy default, nil-hook equivalence, custom key/value hook behavior, preserved capacity/load factor, tombstone-safe lookup, and live-entry-only hook calls."
- Property tests: "Generate randomized operation tests (fixed seed) for `map-hash`, compare with oracle model, and assert invariants after each batch."
- Iterator tests: "Generate tests for `All() iter.Seq2[K,V]` verifying exact count, early stop, empty behavior, and mutation-during-iteration is not safe."
- Benchmarks: "Generate benchmarks for `map-hash` at sizes 1e3/1e4/1e5 with read-heavy, write-heavy, and mixed workloads; report ns/op and allocs/op."
