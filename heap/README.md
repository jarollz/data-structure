# heap

## Overview
A heap is a complete binary tree usually stored in an array. It is commonly used for priority queues.

Min-heap: smallest element at root. Max-heap: largest element at root.

## When to use
- You repeatedly need highest or lowest priority item.
- You need efficient `Push` and `PopTop` operations.

## When not to use
- You need fast lookup of arbitrary values.
- You need full sorted order at all times.

## Pros and cons
- Pros: simple array layout, strong priority-queue performance.
- Cons: not good for membership checks or ordered traversal.

## Complexity
- `Push`: `O(log n)`
- `PopTop`: `O(log n)`
- `PeekTop`: `O(1)`
- Build heap: `O(n)`
- Space: `O(n)`

## Popular Go Libraries
- `container/heap` usage appears in many production systems.
- `github.com/emirpasic/gods/trees/binaryheap` - feature-rich heap.
- `github.com/zyedidia/generic/heap` - generic heap implementation.

## Stdlib (Go 1.25+)
- `container/heap` provides heap operations.

## Language Built-ins
- No built-in heap type.

## Implementation Rules
- Read and follow `heap/RULES.md` before writing code.
