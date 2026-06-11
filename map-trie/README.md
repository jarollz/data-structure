# map-trie

## Overview
A trie map stores string-keyed values by shared prefixes.

Compared with a hash map, it can answer prefix questions directly because keys that share a prefix also share trie path segments.

## Project contract
- `New()` creates an empty trie map.
- Keys are `string` and traversal is byte-wise.
- Empty string key `""` is allowed.
- `Put(key, value)` inserts a new key or overwrites an existing key without changing `Len()` on overwrite.
- `Delete(key)` returns `false` when the key is missing.
- `HasPrefix(prefix)` reports whether any stored key starts with `prefix`.
- `Clone()` returns independent trie copy with same keys, values, and iteration order. Values are copied with normal Go assignment.
- `CloneWith(cloneValue)` returns independent trie copy with same keys and iteration order. Nil hook uses normal Go assignment.
- `All()` yields key-value pairs in ascending byte-lex key order.
- `WithPrefix(prefix)` yields matching key-value pairs in ascending byte-lex key order.
- Mutation during iteration is not safe.

## When to use
- You need exact lookup plus prefix existence or prefix iteration.
- Many keys share common prefixes.

## When not to use
- You only need fastest average exact-key lookup.
- Your keys do not benefit from shared-prefix structure.

## Complexity
- `Get(key)`: path traversal driven by `len(key)`, plus child lookup work at each step
- `Put(key, value)`: path traversal driven by `len(key)`, plus child lookup work at each step
- `Delete(key)`: path traversal driven by `len(key)`, plus child lookup work at each step
- `HasPrefix(prefix)`: path traversal driven by `len(prefix)`, plus child lookup work at each step
- `All()`: `O(total live key bytes + entry count)`
- `WithPrefix(prefix)`: `O(len(prefix) + output subtree work)`
- `Clone()`: `O(total live key bytes + entry count)`
- `CloneWith(cloneValue)`: `O(total live key bytes + entry count)`
- Space: `O(total live key bytes + entry count)`

## Implementation notes
- Use array-backed trie storage only; do not use Go slices or maps in implementation code.
- Representation should distinguish live terminal keys from internal prefix-only nodes.
- Keep child traversal deterministic so iteration order stays ascending byte-lexicographic.
- Delete should prune dead suffix state while preserving shared prefixes still needed by other keys.

## Implementation Rules
- Read and follow `map-trie/SPECS.md` before writing code.
