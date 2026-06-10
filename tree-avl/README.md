# tree-avl

## Overview
An AVL tree is a self-balancing binary search tree.

This folder is set-like tree behavior (single value `T`), not key-value map behavior.

## When to use
- You need sorted set operations.
- You need strict balance and strong lookup performance.

## When not to use
- You need key-value map semantics (use `map-tree-avl`).
- You need fastest average exact lookup without ordering (use hash map).

## Pros and cons
- Pros: strict height balance, predictable `O(log n)`.
- Cons: more rotations than red-black in some update workloads.

## Complexity
- `Insert`: `O(log n)`
- `Delete`: `O(log n)`
- `Has`: `O(log n)`
- In-order traversal: `O(n)`
- Space: `O(n)`

## Popular Go Libraries
- `github.com/zyedidia/generic/avl`
- `github.com/emirpasic/gods/trees/avltree`
- `github.com/Workiva/go-datastructures/tree/avl`

## Stdlib (Go 1.25+)
- No AVL tree in stdlib.

## Language Built-ins
- No built-in balanced tree type.

## Implementation Rules
- Read and follow `tree-avl/RULES.md` before writing code.
