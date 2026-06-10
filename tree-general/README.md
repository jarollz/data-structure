# tree-general

## Overview
This folder is for a general n-ary tree for hierarchical data, not a binary search tree.

Each node can have zero or more children. The public API is ID-based rather than pointer-based.

## Project contract
- `New(rootValue)` creates a tree with one root node.
- Root node ID is always `0` while the tree is non-empty.
- `AddChild(parentID, value)` appends a new last child and returns a new stable child ID.
- Node IDs increase monotonically and are never reused.
- `RemoveSubtree(nodeID)` removes the node and all descendants. `RemoveSubtree(0)` empties the tree.
- `Get` and `Parent` report failure for invalid or removed IDs.
- `ChildCount` returns `-1` for invalid or removed IDs.
- `PreOrder()` yields each parent before its children.
- Mutation during iteration is not safe.

## When to use
- You are modeling parent-child hierarchies such as menus, org charts, or document trees.
- You need stable node IDs.

## When not to use
- You need ordered lookup by comparator.
- You need exact-key lookup like a hash map.

## Complexity
- `Get(nodeID)`: `O(1)`
- `Parent(nodeID)`: `O(1)`
- `AddChild(parentID, value)`: depends on the child-link traversal strategy
- `RemoveSubtree(nodeID)`: `O(size of removed subtree)`
- `PreOrder()`: `O(n)`
- Space: `O(n)`

## Implementation notes
- Use array-backed storage with stable IDs.
- Store `parent`, `firstChild`, `nextSibling`, and `prevSibling` indexes.
- Use `-1` as the nil sentinel index.

## Implementation Rules
- Read and follow `tree-general/RULES.md` before writing code.
