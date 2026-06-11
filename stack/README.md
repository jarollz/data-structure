# stack

## Overview
A stack is a LIFO data structure: the last value pushed is the first value popped.

In this repository, the implementation must use array-backed storage only. Do not use slices or maps in the implementation.

## Project contract
- `New(capacity int)` creates an empty stack. A capacity less than or equal to `0` is normalized to `16`.
- The stack grows and shrinks internally. Callers do not manage resizing directly.
- `Push(v)` always returns `true`.
- `Pop()` and `PeekTop()` return `(zero, false)` when the stack is empty.
- The top element is the element at index `Len()-1`.
- `Clone()` returns independent stack copy with same `Len()`, `Cap()`, and top-to-bottom order. Elements are copied with normal Go assignment.
- `CloneWith(cloneValue)` returns independent stack copy with same `Len()`, `Cap()`, and top-to-bottom order. A nil hook uses normal Go assignment.
- `Values()` yields values from top to bottom.
- Mutation during iteration is not safe.

## When to use
- You need reverse-order processing.
- You are modeling DFS, undo, or call-stack style work.

## When not to use
- You need FIFO behavior.
- You need efficient access to middle elements.

## Complexity
- `Push(v)`: amortized `O(1)`
- `Pop()`: amortized `O(1)`
- `PeekTop()`: `O(1)`
- `Clone()`: `O(n)`
- `CloneWith(cloneValue)`: `O(n)`
- Space: `O(n)`

## Implementation notes
- Keep live elements in backing-storage indexes `[0, Len())`.
- Grow by allocating new storage and copying values manually.
- Shrink with hysteresis after enough pops.

## Implementation Rules
- Read and follow `stack/SPECS.md` before writing code.
