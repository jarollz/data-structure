# queue

## Overview
A queue is a FIFO data structure: the first value enqueued is the first value dequeued.

In this repository, the implementation must use an array-backed circular buffer. Do not use slices or maps in the implementation.

## Project contract
- `New(capacity int)` creates an empty queue. A capacity less than or equal to `0` is normalized to `16`.
- The queue grows and shrinks internally. Callers do not manage resizing directly.
- `Enqueue(v)` adds a value at the back and always returns `true`.
- `Dequeue()` and `PeekFront()` return `(zero, false)` when the queue is empty.
- `Clone()` returns independent queue copy with same `Len()`, `Cap()`, and front-to-back order. Elements are copied with normal Go assignment.
- `CloneWith(cloneValue)` returns independent queue copy with same `Len()`, `Cap()`, and front-to-back order. A nil hook uses normal Go assignment.
- `Values()` yields values from front to back.
- Mutation during iteration is not safe.

## When to use
- You need processing in arrival order.
- You are modeling buffers, jobs, or breadth-first traversal.

## When not to use
- You need LIFO behavior.
- You need priority ordering.

## Complexity
- `Enqueue(v)`: amortized `O(1)`
- `Dequeue()`: amortized `O(1)`
- `PeekFront()`: `O(1)`
- `Clone()`: `O(n)`
- `CloneWith(cloneValue)`: `O(n)`
- Space: `O(n)`

## Implementation notes
- Use circular-buffer indexing.
- Track the head index and size. The tail may be stored or derived.
- On resize, repack values in logical front-to-back order.

## Implementation Rules
- Read and follow `queue/SPECS.md` before writing code.
