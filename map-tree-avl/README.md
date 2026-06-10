# map-tree-avl

## Overview
An AVL tree map stores key-value pairs in ascending key order.

Compared with a hash map, it gives ordered behavior and predictable logarithmic operations at the cost of more balancing work.

## Project contract
- `New(cmp)` creates an empty ordered map.
- The comparator defines key order.
- `Put(key, value)` inserts a new key or overwrites an existing key without changing `Len()` on overwrite.
- `Delete(key)` returns `false` when the key is missing.
- `Min()` and `Max()` return `(zeroK, zeroV, false)` when the map is empty.
- `All()` yields key-value pairs in ascending key order.
- Mutation during iteration is not safe.

## When to use
- You need ordered keys.
- You need `Min`, `Max`, or sorted iteration.

## When not to use
- You only need the fastest average exact-key lookup.
- You do not care about key order.

## Complexity
- `Get(key)`: `O(log n)`
- `Put(key, value)`: `O(log n)`
- `Delete(key)`: `O(log n)`
- `Min()` and `Max()`: `O(log n)`
- `All()`: `O(n)`
- Space: `O(n)`

## Implementation notes
- Use an index-based node pool backed by arrays.
- Use `-1` as the nil sentinel index and track the root explicitly.
- Maintain BST ordering plus AVL height and balance invariants.

## Implementation Rules
- Read and follow `map-tree-avl/RULES.md` before writing code.
