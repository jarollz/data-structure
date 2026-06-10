# queue

## Overview
A queue is a FIFO (first-in, first-out) structure.

Typical operations are enqueue at the back and dequeue from the front.

## When to use
- You need processing in arrival order.
- You model buffers, jobs, or breadth-first traversal.

## When not to use
- You need LIFO behavior (use stack).
- You need priority ordering (use heap/priority queue).

## Pros and cons
- Pros: natural FIFO semantics, simple API.
- Cons: implementation details matter for performance (ring buffer vs linked list).

## Complexity
- Enqueue: `O(1)`
- Dequeue: `O(1)`
- PeekFront: `O(1)`
- Space: `O(n)`

## Popular Go Libraries
- `github.com/eapache/queue` - fast ring-buffer queue.
- `github.com/gammazero/deque` - deque that supports queue usage.
- `github.com/emirpasic/gods/queues/arrayqueue` and `.../linkedlistqueue`.

## Stdlib (Go 1.25+)
- No dedicated queue package.
- Common building blocks: `container/list`, `container/ring`.

## Language Built-ins
- No built-in queue type.

## Implementation Rules
- Read and follow `queue/RULES.md` before writing code.
