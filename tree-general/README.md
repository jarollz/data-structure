# tree-general

## Overview
This folder is for a general tree (n-ary tree), not a binary search tree.

Each node can have zero or more children. This structure is useful for hierarchical data.

## When to use
- You model parent-child hierarchy (file tree, org chart, menu tree).
- You need traversal across hierarchy levels.

## When not to use
- You need fast ordered key lookup by comparator.
- You need map-style key-value lookup guarantees.

## Pros and cons
- Pros: natural representation for hierarchical data, flexible child count.
- Cons: generic search can be `O(n)` without extra indexing.

## Complexity
- Add/remove child: usually `O(1)` to `O(n)` depending on representation.
- Search by value: `O(n)`
- Traversal: `O(n)`
- Space: `O(n)`

## Popular Go Libraries
- `github.com/blazingorb/ntreego` - small n-ary tree library.
- Many Go codebases implement custom n-ary trees per domain model.

## Stdlib (Go 1.25+)
- No general n-ary tree package in stdlib.

## Language Built-ins
- No built-in tree type.

## Implementation Rules
- Read and follow `tree-general/RULES.md` before writing code.
