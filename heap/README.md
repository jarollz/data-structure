# heap

## Overview
A binary heap is a complete binary tree stored in array form. It is the standard base structure for a priority queue.

The comparator decides whether the heap behaves as a min-heap or a max-heap.

## Project contract
- `New(capacity, cmp)` creates an empty heap. A capacity less than or equal to `0` is normalized to `16`.
- `Push(v)` inserts a value, restores the heap property, and always returns `true`.
- `PopTop()` and `PeekTop()` return `(zero, false)` when the heap is empty.
- `Values()` yields internal array order, not sorted order.
- Mutation during iteration is not safe.

## When to use
- You repeatedly need the highest-priority or lowest-priority value.
- You need efficient `Push` and `PopTop` operations.

## When not to use
- You need fast lookup of arbitrary stored values.
- You need sorted traversal without repeated pops.

## Complexity
- `Push(v)`: `O(log n)`
- `PopTop()`: `O(log n)`
- `PeekTop()`: `O(1)`
- Space: `O(n)`

## Implementation notes
- The parent index is `(i-1)/2`.
- The child indexes are `2*i+1` and `2*i+2`.
- Resize backing storage without breaking heap order.

## Implementation Rules
- Read and follow `heap/RULES.md` before writing code.
