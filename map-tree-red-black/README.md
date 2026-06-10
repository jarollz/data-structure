# map-tree-red-black

## Overview
A red-black tree map stores key-value pairs in sorted key order using color rules to keep height near logarithmic.

Compared with AVL, red-black trees are usually less strictly balanced but often need fewer rotations on updates.

## When to use
- You need ordered key-value operations.
- You need reliable `O(log n)` behavior.
- You do frequent inserts and deletes.

## When not to use
- You only need fastest average exact-key lookup.
- You do not need key ordering.

## Pros and cons
- Pros: ordered operations, robust balancing, practical update performance.
- Cons: color-fix logic is complex, higher constants than hash maps.

## Complexity
- `Get`: `O(log n)`
- `Put`: `O(log n)`
- `Delete`: `O(log n)`
- Ordered iteration: `O(n)`
- Space: `O(n)`

## Popular Go Libraries
- `github.com/petar/GoLLRB/llrb` - well-known left-leaning red-black tree.
- `github.com/emirpasic/gods/maps/treemap` - ordered map API backed by a tree.
- `github.com/zyedidia/generic/btree` - ordered map/set alternative (B-tree family).

## Stdlib (Go 1.25+)
- No red-black tree map in stdlib.

## Language Built-ins
- No built-in ordered map type in Go.

## Implementation Rules
- Read and follow `map-tree-red-black/RULES.md` before writing code.
