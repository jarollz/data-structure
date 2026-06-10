# list-array

## Overview
An array list stores elements in contiguous memory and supports index-based access.

In this project, you still follow array-only constraints and avoid using slices in the implementation.

## When to use
- You need fast random access by index.
- Most updates are append or tail removals.

## When not to use
- You do many inserts/deletes in the middle.
- Capacity limits are very small and resizing cost is unacceptable.

## Pros and cons
- Pros: fast indexing, cache-friendly, simple model.
- Cons: middle insert/delete shifts many items, resize copies data.

## Complexity
- `Get(i)`: `O(1)`
- `Set(i)`: `O(1)`
- `Append`: amortized `O(1)` with growth strategy
- `Insert/Delete` middle: `O(n)`
- Space: `O(n)`

## Popular Go Libraries
- `github.com/emirpasic/gods/lists/arraylist` - array-list style API.
- `github.com/zyedidia/generic/rope` - sequence alternative for heavy middle edits.

## Stdlib (Go 1.25+)
- No dedicated ArrayList type in stdlib.

## Language Built-ins
- Arrays and slices are the common sequence primitives in Go.
- This project forbids using slices in your implementation.

## Implementation Rules
- Read and follow `list-array/RULES.md` before writing code.
