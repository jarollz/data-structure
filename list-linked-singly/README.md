# list-linked-singly

## Overview
A singly linked list stores each element in a node that points to the next node.

In this repository, the implementation must use an index-based node pool backed by arrays. Do not use slices, maps, or pointer-linked nodes.

## Project contract
- `PushFront(v)` adds a value at the head.
- `PopFront()` returns `(zero, false)` when the list is empty.
- `Append(v)` appends at the tail in `O(1)` time.
- `DeleteFirst(match)` removes the first value that matches and returns `false` when no match exists.
- `Clone()` returns independent list copy with same `Len()` and head-to-tail order. Elements are copied with normal Go assignment.
- `CloneWith(cloneValue)` returns independent list copy with same `Len()` and head-to-tail order. A nil hook uses normal Go assignment.
- `Values()` yields values from head to tail.
- Mutation during iteration is not safe.

## When to use
- You need frequent inserts and removals at the head.
- You only need forward traversal.

## When not to use
- You need fast random access by index.
- You need reverse traversal.

## Complexity
- `PushFront(v)`: `O(1)`
- `PopFront()`: `O(1)`
- `Append(v)`: `O(1)`
- `DeleteFirst(match)`: `O(n)`
- `Clone()`: `O(n)`
- `CloneWith(cloneValue)`: `O(n)`
- Space: `O(n)`

## Implementation notes
- Track the head index, tail index, free-list head, and length.
- Tail tracking is required by contract, not optional optimization.
- Use `-1` as the nil sentinel index.
- When the node pool is full, allocate larger arrays and copy node fields by index.

## Implementation Rules
- Read and follow `list-linked-singly/SPECS.md` before writing code.
