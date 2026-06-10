# tree-red-black

## Overview
A red-black tree is a self-balancing binary search tree based on node color rules.

This folder is set-like tree behavior (single value `T`), not key-value map behavior.

## When to use
- You need sorted set operations with reliable logarithmic complexity.
- You expect frequent inserts and deletes.

## When not to use
- You need key-value map semantics (use `map-tree-red-black`).
- You only need fastest average key lookup without order.

## Pros and cons
- Pros: balanced tree with practical update performance, ordered traversal.
- Cons: color-fix logic is harder to debug than simple lists/arrays.

## Complexity
- `Insert`: `O(log n)`
- `Delete`: `O(log n)`
- `Has`: `O(log n)`
- In-order traversal: `O(n)`
- Space: `O(n)`

## Popular Go Libraries
- `github.com/petar/GoLLRB/llrb`
- `github.com/emirpasic/gods/trees/redblacktree`

## Stdlib (Go 1.25+)
- No red-black tree in stdlib.

## Language Built-ins
- No built-in balanced tree type.

## Implementation Rules
- Read and follow `tree-red-black/RULES.md` before writing code.
