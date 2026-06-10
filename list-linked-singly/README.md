# list-linked-singly

## Overview
A singly linked list stores elements in nodes where each node points to the next node.

It supports efficient inserts/removes near the head, but random indexing is slow.

## When to use
- You need frequent head inserts/removes.
- You do not need fast random access.

## When not to use
- You need fast index lookup.
- You need reverse traversal without extra work.

## Pros and cons
- Pros: simple node model, no large contiguous memory required.
- Cons: `O(n)` indexing/search, pointer chasing hurts cache locality.

## Complexity
- `PushFront`: `O(1)`
- `PopFront`: `O(1)`
- `Append`: `O(1)` with tail pointer, else `O(n)`
- `Find`: `O(n)`
- Space: `O(n)`

## Popular Go Libraries
- `github.com/emirpasic/gods/lists/singlylinkedlist`
- `github.com/golang-collections/collections` (older collection implementations)

## Stdlib (Go 1.25+)
- No dedicated singly linked list in stdlib.

## Language Built-ins
- No built-in linked-list type.

## Implementation Rules
- Read and follow `list-linked-singly/RULES.md` before writing code.
