# map-hash

## Overview
A hash map stores key-value pairs and targets very fast lookup by key.

This learning implementation does not use Go's built-in `map`. You implement hashing, probing, deletion, and resizing directly.

This folder uses open addressing with linear probing.

## When to use
- You need fast average `Get`/`Put`/`Delete`.
- You do not need sorted key order.
- You want direct control over load factor and resize behavior.

## When not to use
- You need sorted iteration or range queries.
- You need strict worst-case `O(log n)` guarantees.

## Pros and cons
- Pros: very fast average access, cache-friendly array layout, simple memory model.
- Cons: delete logic needs tombstones, clustering can hurt performance, worst case is `O(n)`.

## Complexity
- `Get`: average `O(1)`, worst `O(n)`
- `Put`: average `O(1)`, worst `O(n)`
- `Delete`: average `O(1)`, worst `O(n)`
- Space: `O(n)`

## Popular Go Libraries
- `github.com/cornelk/hashmap` - lock-free/thread-safe map focused on read performance.
- `github.com/orcaman/concurrent-map/v2` - sharded concurrent map for practical service workloads.
- `github.com/zyedidia/generic/hashmap` - generic hashmap with linear probing.

## Stdlib (Go 1.25+)
- No low-level hash table package for custom probing logic.
- `sync.Map` exists for specific concurrent access patterns.

## Language Built-ins
- Built-in `map[K]V` provides hash-map behavior directly.
- This project forbids using built-in `map` in your own implementation.

## Implementation Rules
- Read and follow `map-hash/RULES.md` before writing code.
