# list-linked-doubly

## Overview
A doubly linked list stores each node with both previous and next links.

In this repository, the implementation must use an index-based node pool backed by arrays. Do not use slices, maps, or pointer-linked nodes.

## Project contract
- `PushFront(v)` and `PushBack(v)` both run in `O(1)` time.
- `PopFront()` and `PopBack()` return `(zero, false)` when the list is empty.
- `Clone()` returns independent list copy with same `Len()` and head-to-tail order. Elements are copied with normal Go assignment.
- `CloneWith(cloneValue)` returns independent list copy with same `Len()` and head-to-tail order. A nil hook uses normal Go assignment.
- `Values()` yields values from head to tail.
- The empty state uses the nil sentinel for both head and tail.
- Mutation during iteration is not safe.

## When to use
- You need efficient updates at both ends.
- You need links in both directions.

## When not to use
- You need fast random access by index.
- You want the most compact memory layout.

## Complexity
- `PushFront(v)`: `O(1)`
- `PushBack(v)`: `O(1)`
- `PopFront()`: `O(1)`
- `PopBack()`: `O(1)`
- `Clone()`: `O(n)`
- `CloneWith(cloneValue)`: `O(n)`
- Space: `O(n)`

## Implementation notes
- Track the head index, tail index, free-list head, and length.
- Use `-1` as the nil sentinel index.
- Keep `prev` and `next` links consistent in both directions.

## Implementation Rules
- Read and follow `list-linked-doubly/SPECS.md` before writing code.
