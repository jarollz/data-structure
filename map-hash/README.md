# map-hash

## Overview
A hash map stores key-value pairs and targets very fast average lookup by key.

In this repository, the implementation uses open addressing with linear probing and tombstones. Do not use Go's built-in `map` in the implementation.

## Project contract
- `New(capacity, hash, equal)` creates an empty map. A capacity less than or equal to `0` is normalized to `16`.
- The caller provides `hash func(K) uint64` and `equal func(a, b K) bool`.
- `Put(key, value)` inserts a new key or overwrites an existing key.
- `Delete(key)` returns `false` when the key is missing.
- `LoadFactor()` is `Len()/Cap()` and does not count tombstones.
- Iteration order from `All()` is unspecified.
- Mutation during iteration is not safe.

## When to use
- You need fast average `Get`, `Put`, and `Delete` operations.
- You do not need sorted key order.

## When not to use
- You need ordered iteration or range-style queries.
- You need strict worst-case logarithmic behavior.

## Complexity
- `Get(key)`: average `O(1)`, worst `O(n)`
- `Put(key, value)`: average `O(1)`, worst `O(n)`
- `Delete(key)`: average `O(1)`, worst `O(n)`
- Space: `O(n)`

## Implementation notes
- Use an array of slots with explicit states: empty, occupied, and deleted.
- Reuse tombstones carefully so probe chains remain correct.
- Grow based on probe occupancy, not only live load factor.
- `Clear()` removes live entries and tombstones and resets capacity to the minimum configured capacity.

## Implementation Rules
- Read and follow `map-hash/RULES.md` before writing code.
