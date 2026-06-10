# map-tree-avl

## Overview
An AVL tree map stores key-value pairs in sorted key order. It is a binary search tree that rebalances after updates.

Compared with a hash map, this structure is slower on average lookup but supports ordered operations.

## When to use
- You need ordered keys.
- You need range queries, min/max, predecessor/successor, or sorted iteration.
- You want predictable `O(log n)` operations.

## When not to use
- You only need fastest average exact-key lookup.
- You do not care about key order.

## Pros and cons
- Pros: sorted iteration, stable `O(log n)` bounds, good for range operations.
- Cons: more complex rebalancing logic, higher constant factors than hash maps.

## Complexity
- `Get`: `O(log n)`
- `Put`: `O(log n)`
- `Delete`: `O(log n)`
- Ordered iteration: `O(n)`
- Space: `O(n)`

## Popular Go Libraries
- `github.com/zyedidia/generic/avl` - generic AVL tree.
- `github.com/Workiva/go-datastructures/tree/avl` - immutable AVL implementation.
- `github.com/emirpasic/gods/maps/treemap` - ordered map API (tree-based).

## Stdlib (Go 1.25+)
- No AVL tree map in stdlib.

## Language Built-ins
- No built-in ordered map type in Go.

## Implementation Rules
- Read and follow `map-tree-avl/RULES.md` before writing code.
