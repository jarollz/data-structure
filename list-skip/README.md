# list-skip

## Overview
A skip list is an ordered structure built from multiple forward-link levels. Higher levels skip over many values, which gives expected logarithmic search, insert, and delete.

In this repository, it is a set-like structure: values are unique. The implementation must use array-backed index pools only.

## Project contract
- `New(maxLevel, cmp)` creates an empty skip list. A `maxLevel` less than `1` is normalized to `1`.
- The comparator defines sorted order.
- Duplicate values are not stored. `Insert(v)` returns `false` when the value already exists.
- `Delete(v)` returns `false` when the value is missing.
- `Clone()` returns independent skip-list copy with same sorted order, comparator, `maxLevel`, `currentLevel`, and deterministic RNG state. Elements are copied with normal Go assignment.
- `CloneWith(cloneValue)` returns independent skip-list copy with same sorted order, comparator, `maxLevel`, `currentLevel`, and deterministic RNG state. A nil hook uses normal Go assignment.
- `Values()` yields values in sorted order by traversing level `0`.
- Level generation is deterministic for tests. The same operation sequence must produce the same levels.
- Mutation during iteration is not safe.

## When to use
- You need ordered set behavior with expected `O(log n)` updates.
- You want simpler balancing rules than tree rotations.

## When not to use
- You need strict worst-case balancing guarantees.
- You need minimal per-element overhead.

## Complexity
- `Has(v)`: expected `O(log n)`
- `Insert(v)`: expected `O(log n)`
- `Delete(v)`: expected `O(log n)`
- `Clone()`: `O(n)`
- `CloneWith(cloneValue)`: `O(n)`
- Space: `O(n)` on average

## Implementation notes
- Use a head sentinel node and `-1` as the nil sentinel index.
- Track `currentLevel`, length, free-list head, and deterministic RNG state.
- Reduce `currentLevel` after deletions when upper levels become empty.

## Implementation Rules
- Read and follow `list-skip/SPECS.md` before writing code.
