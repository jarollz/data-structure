# RULES.md - map-hash

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
- [ ] `New(capacity int, hash func(K) uint64, equal func(a, b K) bool) *Map[K, V]` creates an empty map. Normalize `capacity <= 0` to `16`. Effective starting capacity is `max(16, capacity)`.
- [ ] `Put(key K, value V)` inserts or overwrites. Overwrite does not change `Len()`.
- [ ] `Get(key K) (V, bool)`
- [ ] `Delete(key K) bool`
- [ ] `Has(key K) bool`
- [ ] `Len() int`
- [ ] `Cap() int`
- [ ] `Clear()`
- [ ] `LoadFactor() float64` returns `float64(Len()) / float64(Cap())` and excludes tombstones.
- [ ] `All() iter.Seq2[K, V]`

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

## Test checklist
- [ ] Unit tests for every API method.
- [ ] Collision tests for multiple keys into same probe region.
- [ ] Tombstone tests with middle-of-chain deletes.
- [ ] Resize tests that preserve all live entries.
- [ ] Random operation tests against a reference model.
- [ ] Invariant check after operation sequences.

## Benchmark checklist
- [ ] `Put` benchmark with unique keys.
- [ ] `Get` benchmark for hit-heavy workload.
- [ ] `Get` benchmark for miss-heavy workload.
- [ ] `Delete` benchmark for mixed workload.
- [ ] Mixed benchmark (`Put/Get/Delete`) across multiple sizes.
- [ ] Record load factor during benchmark runs.

## Test Generator Hints
- Focus collision-heavy cases, tombstones, and resize/rehash correctness.
- Run randomized `Put/Get/Delete` with fixed seed and compare against test oracle.
- Validate invariants after each mutation batch, especially probe-chain reachability.
- Iterator tests must check count, early stop, and unspecified iteration order.

## AI Prompt Snippets
- Unit tests: "Generate table-driven Go tests for `map-hash` API from this RULES file, including collisions, overwrite, delete-missing, and resize behavior."
- Property tests: "Generate randomized operation tests (fixed seed) for `map-hash`, compare with oracle model, and assert invariants after each batch."
- Iterator tests: "Generate tests for `All() iter.Seq2[K,V]` verifying exact count, early stop, empty behavior, and mutation-during-iteration is not safe."
- Benchmarks: "Generate benchmarks for `map-hash` at sizes 1e3/1e4/1e5 with read-heavy, write-heavy, and mixed workloads; report ns/op and allocs/op."
