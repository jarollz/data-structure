# list-skip

## Overview
A skip list is an ordered structure built from multiple linked levels. Higher levels skip over elements to speed up search.

It is often used as a simpler alternative to balanced trees for ordered maps/sets.

## When to use
- You need ordered operations with expected logarithmic performance.
- You want a probabilistic balancing approach.

## When not to use
- You need strict worst-case balancing guarantees.
- You need very compact memory layout.

## Pros and cons
- Pros: conceptually simpler balancing than tree rotations, good expected performance.
- Cons: probabilistic behavior, pointer/index overhead across levels.

## Complexity (expected)
- Search: `O(log n)`
- Insert: `O(log n)`
- Delete: `O(log n)`
- Space: `O(n)` average

## Popular Go Libraries
- `github.com/huandu/skiplist` - popular skip list ordered map implementation.
- `github.com/Workiva/go-datastructures` includes skip list package.

## Stdlib (Go 1.25+)
- No skip list in stdlib.

## Language Built-ins
- No built-in skip-list type.

## Implementation Rules
- Read and follow `list-skip/RULES.md` before writing code.
