# tree-red-black

## Overview
A red-black tree is a self-balancing binary search tree maintained by color rules.

This folder is set-like: it stores unique values of type `T`, not key-value pairs.

## Project contract
- `New(cmp)` creates an empty tree.
- The comparator defines sorted order.
- Duplicate values are not stored. `Insert(v)` returns `false` when the value already exists.
- `Delete(v)` returns `false` when the value is missing.
- `Min()` and `Max()` return `(zero, false)` when the tree is empty.
- `InOrder()` yields values in ascending order.
- Mutation during iteration is not safe.

## When to use
- You need ordered set behavior with reliable `O(log n)` operations.
- You expect frequent inserts and deletes.

## When not to use
- You need key-value map behavior.
- You only need the fastest average exact lookup without order.

## Complexity
- `Insert(v)`: `O(log n)`
- `Delete(v)`: `O(log n)`
- `Has(v)`: `O(log n)`
- `Min()` and `Max()`: `O(log n)`
- `InOrder()`: `O(n)`
- Space: `O(n)`

## Implementation notes
- Use an index-based node pool backed by arrays.
- Use `-1` as the nil sentinel index and track the root explicitly.
- Maintain the root-black, no-red-red, and equal black-height invariants.

## Implementation Rules
- Read and follow `tree-red-black/RULES.md` before writing code.
