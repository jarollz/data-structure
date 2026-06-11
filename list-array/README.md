# list-array

## Overview
An array list stores elements in contiguous storage and supports direct index access.

In this repository, the implementation must manage storage manually with arrays only. Do not use Go slices or maps in the implementation.

## Project contract
- `New(capacity int)` creates an empty list. A capacity less than or equal to `0` is normalized to `16`.
- The list grows and shrinks internally. Callers do not manage resizing directly.
- `Append(v)` always succeeds and returns `true`.
- `Insert(i, v)` inserts before index `i`. Valid indexes are `0..Len()`.
- `Get`, `Set`, and `Delete` handle out-of-range indexes safely.
- `Clone()` returns independent list copy with same `Len()`, `Cap()`, and element order. Elements are copied with normal Go assignment.
- `CloneWith(cloneValue)` returns independent list copy with same `Len()`, `Cap()`, and element order. A nil hook uses normal Go assignment.
- `Values()` yields elements in index order from `0` to `Len()-1`.
- Mutation during iteration is not safe.

## When to use
- You need fast random access by index.
- Most writes are appends or tail updates.

## When not to use
- You perform many inserts or deletes in the middle.
- You need stable element positions while the structure resizes.

## Complexity
- `Get(i)`: `O(1)`
- `Set(i, v)`: `O(1)`
- `Append(v)`: amortized `O(1)`
- `Insert(i, v)`: `O(n)`
- `Delete(i)`: `O(n)`
- `Clone()`: `O(n)`
- `CloneWith(cloneValue)`: `O(n)`
- Space: `O(n)`

## Implementation notes
- Keep live elements in indexes `[0, Len())` with no gaps.
- Resize by allocating new storage and copying elements manually.
- Preserve relative order after insert and delete.

## Implementation Rules
- Read and follow `list-array/SPECS.md` before writing code.
