# SPECS.md - map-trie

## Goal
Implement a string-keyed trie map with exact lookup and prefix queries.

## Hard constraints
- [ ] Use only simple Go types and arrays for internal storage.
- [ ] Do not use `slice` in implementation code.
- [ ] Do not use `map` in implementation code.

## Generics and key contract
- [ ] Use `V any` for stored values.
- [ ] Keys are `string` only.
- [ ] Key traversal semantics are byte-wise over `string` indexes.
- [ ] Empty string key `""` is allowed.
- [ ] Performance notation: `m = len(key)`, `p = len(prefix)`, and `childLookupCost` is cost to locate next child in one trie node under chosen array-backed representation.

## Required API
Each public API below is part of contract. Use each checklist as implementation target, not loose guidance.

### `New() *MapTrie[V]`
Purpose
- [ ] Create empty trie map.

Behavior expectations
- [ ] Returned trie is non-nil and empty.
- [ ] `Len()` is `0` immediately after construction.
- [ ] Empty string key is absent immediately after construction.

Performance expectations
- [ ] `O(1)` time.

### `Put(key string, value V)`
Purpose
- [ ] Insert new key-value entry or overwrite existing key.

Behavior expectations
- [ ] Missing key inserts new entry.
- [ ] Existing key overwrites stored value and does not change `Len()`.
- [ ] Shared prefixes are reused instead of duplicated logically.
- [ ] Empty string key is stored at root terminal state.

Performance expectations
- [ ] `O(m * childLookupCost)` time for `m = len(key)`.

### `Get(key string) (V, bool)`
Purpose
- [ ] Read value currently stored for key.

Behavior expectations
- [ ] Existing key returns `(value, true)`.
- [ ] Missing key returns `(zero, false)`.
- [ ] Prefix-only path without terminal value returns `(zero, false)`.
- [ ] `Get` never mutates trie state.

Performance expectations
- [ ] `O(m * childLookupCost)` time for `m = len(key)`.

### `Delete(key string) bool`
Purpose
- [ ] Remove live entry for key.

Behavior expectations
- [ ] Existing key returns `true` and is removed.
- [ ] Missing key returns `false`.
- [ ] Deleting key that is prefix of longer key clears only terminal value and keeps longer key reachable.
- [ ] Deleting leaf key prunes dead suffix nodes no longer needed by any other key.

Performance expectations
- [ ] `O(m * childLookupCost)` time for `m = len(key)`.

### `Has(key string) bool`
Purpose
- [ ] Report whether key currently exists.

Behavior expectations
- [ ] Existing key returns `true`.
- [ ] Missing key returns `false`.
- [ ] Prefix-only path without terminal value returns `false`.

Performance expectations
- [ ] `O(m * childLookupCost)` time for `m = len(key)`.

### `HasPrefix(prefix string) bool`
Purpose
- [ ] Report whether any stored key starts with `prefix`.

Behavior expectations
- [ ] Existing stored key implies `HasPrefix(key) == true`.
- [ ] Existing longer key implies `HasPrefix(shorterPrefix) == true`.
- [ ] Missing prefix returns `false`.
- [ ] Empty prefix returns `true` when trie is non-empty and `false` when trie is empty.

Performance expectations
- [ ] `O(p * childLookupCost)` time for `p = len(prefix)`.

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
- [ ] Reset trie to empty logical state.

Behavior expectations
- [ ] `Len()` becomes `0`.
- [ ] All old keys become unreachable.
- [ ] Future puts, deletes, and queries still work correctly.
- [ ] Clear on empty trie is safe.

Performance expectations
- [ ] Target `O(1)` logical reset.

### `Clone() *MapTrie[V]`
Purpose
- [ ] Create independent trie copy with same keys, values, and iteration order.

Behavior expectations
- [ ] Clone preserves `Len()`, exact lookup results, prefix behavior, and ascending byte-lex iteration order.
- [ ] Values are copied with normal Go assignment.
- [ ] Clone is container-independent from source.
- [ ] Empty source produces empty independent clone.

Performance expectations
- [ ] `O(total live key bytes + entry count)` time.
- [ ] `O(total live key bytes + entry count)` extra storage.

### `CloneWith(cloneValue func(V) V) *MapTrie[V]`
Purpose
- [ ] Create independent trie copy while optionally transforming each live value.

Behavior expectations
- [ ] Preserves `Len()`, exact lookup results, prefix behavior, and ascending byte-lex iteration order.
- [ ] Nil hook is equivalent to `Clone()`.
- [ ] Non-nil hook is called exactly once per live entry.
- [ ] Hook call order is ascending byte-lex key order.
- [ ] Hook is never called for non-terminal prefix nodes.

Performance expectations
- [ ] `O(total live key bytes + entry count)` container work plus hook cost.

### `All() iter.Seq2[string, V]`
Purpose
- [ ] Iterate all live entries in ascending byte-lex key order.

Behavior expectations
- [ ] Yields each live key-value pair exactly once.
- [ ] Empty trie yields nothing.
- [ ] Keys are yielded in ascending byte-lex key order.
- [ ] Shorter key is yielded before longer key when shorter key is prefix of longer key.
- [ ] Early stop works when consumer returns `false`.
- [ ] Mutation during iteration is not safe.

Performance expectations
- [ ] Full traversal is `O(total live key bytes + entry count)`.
- [ ] Early-stop traversal is proportional to visited output prefix.
- [ ] Iterator setup is `O(1)`.

### `WithPrefix(prefix string) iter.Seq2[string, V]`
Purpose
- [ ] Iterate live entries whose keys start with `prefix`.

Behavior expectations
- [ ] Missing prefix yields nothing.
- [ ] Empty prefix yields same entries and same order as `All()`.
- [ ] Matching entries are yielded in ascending byte-lex key order.
- [ ] Exact key equal to prefix is yielded before longer descendants when that key exists.
- [ ] Early stop works when consumer returns `false`.
- [ ] Mutation during iteration is not safe.

Performance expectations
- [ ] `O(p * childLookupCost + visitedSubtreeWork)` time for `p = len(prefix)`.
- [ ] `visitedSubtreeWork` is proportional to visited subtree needed to produce yielded matches.

## Internal representation
- [ ] Use array-backed trie storage only; do not use Go slices or maps.
- [ ] Root represents empty-prefix state.
- [ ] Keys are traversed byte by byte.
- [ ] Representation must distinguish live terminal keys from internal prefix-only nodes.
- [ ] Representation must support exact-key lookup, prefix lookup, delete cleanup, and ascending byte-lex iteration.
- [ ] Integer indexes and sentinel values are allowed and recommended when using node pools.

## Auto-resize policy
- [ ] No public `Grow()` or `Shrink()` API.
- [ ] When storage is exhausted, allocate larger array-backed storage and copy or rebuild internal references as needed.
- [ ] Allocate storage on insert operations.
- [ ] Reclaim or reuse storage on delete and clear operations.
- [ ] Reuse strategy is implementation-defined as long as trie invariants remain valid.
- [ ] `Clear()` may reset trie by rebuilding fresh empty root state.

## Invariants
- [ ] `Len()` equals number of terminal nodes with stored live values.
- [ ] Every live non-root node is reachable from root by parent and child links.
- [ ] Child sibling order is ascending by edge byte.
- [ ] No dead suffix nodes remain after successful delete pruning.
- [ ] `Clone()` and `CloneWith(...)` preserve `Len()`, exact lookup results, prefix behavior, and byte-lex iteration order in clone.

## Iterator contract
- [ ] `All()` yields each live key-value pair exactly once.
- [ ] `WithPrefix(prefix)` yields each matching key-value pair exactly once.
- [ ] Yield order is ascending byte-lex key order.
- [ ] Early stop works when `yield` returns `false`.
- [ ] Empty result yields nothing and does not panic.
- [ ] Mutation during iteration is not safe.

## Edge cases
- [ ] Empty string key can be inserted, read, deleted, cloned, and iterated.
- [ ] Prefix-only path without terminal value is not treated as live key.
- [ ] Deleting key that has descendants keeps descendants reachable.
- [ ] Deleting empty string key removes only root terminal value and preserves longer keys.
- [ ] Deleting last key under prefix removes that dead suffix path.
- [ ] `HasPrefix("")` reflects whether trie is empty.
- [ ] `WithPrefix("")` is equivalent to `All()`.
- [ ] `Clone()` on empty trie returns empty independent trie.
- [ ] `CloneWith(nil)` is equivalent to `Clone()`.
- [ ] `CloneWith(...)` calls hook only for live entries in ascending byte-lex order.

## Test checklist
- [ ] API behavior tests.
- [ ] Shared-prefix insert/delete tests.
- [ ] Empty string key tests.
- [ ] Exact lookup versus prefix-only path tests.
- [ ] Randomized operations against reference model.
- [ ] `Clone/CloneWith` tests for independence, nil-hook equivalence, shallow-copy default, custom hook behavior, and live-entry-only hook calls.
- [ ] Iterator tests for order, exact count, empty behavior, prefix filtering, and early stop.

## Benchmark checklist
- [ ] Benchmark only non-trivial public APIs.
- [ ] Do not benchmark trivial metadata accessor `Len()`.
- [ ] `Put` benchmark.
- [ ] `Get` hit benchmark.
- [ ] `Get` miss benchmark.
- [ ] `Delete` benchmark.
- [ ] `Has` benchmark.
- [ ] `HasPrefix` benchmark.
- [ ] `Clear` benchmark.
- [ ] `Clone` benchmark.
- [ ] `CloneWith` benchmark.
- [ ] `All` benchmark.
- [ ] `WithPrefix` benchmark.
- [ ] Mixed `Put/Get/Delete` benchmark.
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
- [ ] Exact lookup and prefix-probe APIs should keep `ns/op(1e5) <= 3.0 * ns/op(1e3)` for similar key-shape workloads.
- [ ] Iteration-style APIs must compare normalized visited-output cost and keep normalized `1e5` result within `3.0x` of normalized `1e3` result.

### Layer 2: Absolute timing thresholds
- [ ] Absolute thresholds are per API, per payload class, and per benchmark size.
- [ ] Provisional thresholds should be derived from expected complexity and payload size.
- [ ] Calibrated thresholds should be derived from repeated local runs, median `ns/op`, and safety factor.
- [ ] Suggested safety factor for stable simple trie-map ops is `1.75x`.
- [ ] Suggested safety factor for subtree iteration and prefix workloads is `2.00x`.
- [ ] Suggested safety factor for copy-heavy clone or large-payload workloads is `2.25x`.

### Payload classes
- [ ] Tiny payload benchmark uses scalar-like values such as `int`.
- [ ] Large payload benchmark uses fixed-size struct or fixed-size array payload.
- [ ] Large-payload thresholds must be tracked independently from tiny-payload thresholds.

## Test Generator Hints
- Emphasize shared prefixes, delete pruning, and exact-key versus prefix-only distinctions.
- Compare `HasPrefix` and `WithPrefix` against oracle scan over reference model.
- Validate `Clone/CloneWith` preserve lookup results and byte-lex iteration order.
- Iterator tests must verify shorter-prefix-before-descendant ordering and early stop.

## AI Prompt Snippets
- Unit tests: "Generate table-driven tests for trie-map API including shared prefixes, empty string key, delete pruning, and prefix lookup behavior."
- Clone tests: "Generate tests for `Clone()` and `CloneWith(...)` verifying independent container state, shallow-copy default, nil-hook equivalence, custom hook behavior, and live-entry-only hook calls in byte-lex order."
- Property tests: "Generate randomized string-key operations with fixed seed for trie map, compare against oracle map, and assert exact lookup plus prefix behavior after each batch."
- Iterator tests: "Generate tests for `All()` and `WithPrefix(prefix)` verifying byte-lex order, exact count, prefix filtering, early stop, and mutation-unsafety note."
- Benchmarks: "Generate trie-map benchmarks for exact lookup, prefix queries, clone, and traversal workloads at 1e3/1e4/1e5."
