# list-linked-doubly

## Overview
A doubly linked list stores nodes with both `prev` and `next` links.

It supports efficient insert/delete at both ends and around known nodes.

## When to use
- You need frequent inserts/removes at both ends.
- You need bidirectional traversal.

## When not to use
- You need fast random access by index.
- You need compact memory usage.

## Pros and cons
- Pros: efficient local edits, supports forward/backward traversal.
- Cons: extra pointer/index overhead, cache locality worse than arrays.

## Complexity
- `PushFront/PushBack`: `O(1)`
- `PopFront/PopBack`: `O(1)`
- `Insert/Delete` at known node: `O(1)`
- `Find`: `O(n)`
- Space: `O(n)`

## Popular Go Libraries
- `container/list` (stdlib) is widely used.
- `github.com/emirpasic/gods/lists/doublylinkedlist`
- `github.com/zyedidia/generic/list`

## Stdlib (Go 1.25+)
- `container/list` implements a doubly linked list.

## Language Built-ins
- No built-in linked-list type.

## Implementation Rules
- Read and follow `list-linked-doubly/RULES.md` before writing code.
