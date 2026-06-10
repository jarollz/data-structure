# stack

## Overview
A stack is a LIFO (last-in, first-out) structure.

Typical operations are push, pop, and peek top.

## When to use
- You need reverse-order processing.
- You model call stacks, undo history, or DFS traversal.

## When not to use
- You need FIFO processing (use queue).
- You need access by arbitrary key or index.

## Pros and cons
- Pros: very simple API, efficient top operations.
- Cons: no direct access to middle elements.

## Complexity
- `Push`: `O(1)`
- `Pop`: `O(1)`
- `PeekTop`: `O(1)`
- Space: `O(n)`

## Popular Go Libraries
- `github.com/zyedidia/generic/stack`
- `github.com/emirpasic/gods/stacks/arraystack` and `.../linkedliststack`
- `github.com/golang-collections/collections/stack` (older library)

## Stdlib (Go 1.25+)
- No dedicated stack package.

## Language Built-ins
- No built-in stack type.

## Implementation Rules
- Read and follow `stack/RULES.md` before writing code.
