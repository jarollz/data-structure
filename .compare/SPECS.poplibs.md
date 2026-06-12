# SPECS.poplibs.md

## Goal
Generate `.compare/poplibs.json` mapping each in-scope data-structure folder in this repo to up to 3 popular external Go libraries with overlapping APIs.

Output must be reproducible from this spec, while allowing ranking/content changes over time as popularity and upstream APIs change.

## Hard rules
- Do not hardcode structure folder names.
- Do not hardcode candidate library registry.
- Discover structure folders and candidate libraries dynamically each generation run.
- Keep process deterministic given same metric/API snapshot.

## Output file contract

### Path
- `.compare/poplibs.json`

### Encoding and format
- UTF-8 JSON
- 2-space indentation
- Top-level JSON object

### JSON structure (required)
- Type: `object`
- Key: structure folder name (string)
- Value: array of library objects (`[]LibraryMatch`)
- Maximum list length per folder: `3`
- Lists sorted by weighted popularity score descending

`LibraryMatch` schema:
- `import_name` (string): concrete Go import path for matched package
- `git_url` (string): canonical repository URL (HTTPS)
- `common_apis` (`[]CommonAPI`): matched API pairs

`CommonAPI` schema:
- `internal` (string): exact API identifier from this repo folder `api_contract.go`
- `external` (string): exact exported API identifier from external package
- `confidence` (number): closed interval `[0.00, 1.00]`, rounded to 2 decimals

Only include library object when `common_apis` is non-empty after threshold filtering.

Example shape only:

```json
{
  "queue": [
    {
      "import_name": "github.com/example/lib/queue",
      "git_url": "https://github.com/example/lib",
      "common_apis": [
        {
          "internal": "Enqueue",
          "external": "Push",
          "confidence": 0.83
        }
      ]
    }
  ]
}
```

## Step 1: Discover structure folders dynamically

Source of truth: repo root `go.work`.

Folder in scope iff all true:
1. Entry appears in `go.work` `use (...)` block.
2. Entry resolves to existing directory in repo.
3. Directory contains `SPECS.md`.

Normalization:
- Convert `./name` to `name` for output key.
- Preserve deterministic key order:
  - Primary: order in `go.work`
  - Secondary (tie/stability fallback): lexicographic

## Step 2: Build structure intent per folder
For each discovered folder, inspect:
- `<folder>/SPECS.md`
- `<folder>/api_contract.go`
- `<folder>/README.md` (if exists)

Infer intent from folder tokens and contracts (examples: list, queue, stack, heap, trie, map, avl, red-black, tree).

## Step 3: Discover candidates dynamically (no static registry)
For each folder intent, run multiple GitHub repository searches (language Go), such as:
- `golang <intent> data structure`
- `go <intent> library`
- `language:Go <intent>`

Build candidate pool from top search results (recommended 20-50 unique repos per folder).

Candidate eligibility:
- Public reachable repo
- Go library-like repository
- Evidence in description/README/code that target structure is implemented

Exclude obvious non-library results (pure tutorials, interview collections, problemset repos) unless concrete reusable package clearly exists.

## Step 4: Collect popularity metrics
For each candidate repo, collect:
1. `stars` (GitHub stars)
2. `recency` (months since last commit/push)
3. `ecosystem` (prefer pkg.go.dev importers; fallback GitHub forks)

Recommended fetch order:
- GitHub API/search pages
- Shields badges as fallback

## Step 5: Weighted popularity score
Within each folder candidate set:
- `stars_norm = stars / max(stars)` (0 if max is 0)
- `recency_norm = max(0, 1 - months_old/60)`
- `ecosystem_norm = ecosystem / max(ecosystem)` (0 if max is 0)

Score:
- `popularity_score = 0.70*stars_norm + 0.20*recency_norm + 0.10*ecosystem_norm`

Tie-break order:
1. Higher raw stars
2. More recent commit/push (smaller months_old)
3. Higher ecosystem raw value
4. Lexicographically smaller `import_name`

## Step 6: API overlap extraction
For each selected candidate package:
1. Detect concrete import path via `go.mod` module path + package directory exposing target APIs.
2. Extract exported APIs from code/docs/examples.
3. Compare against internal APIs in `<folder>/api_contract.go`.

Allowed mapping cardinality:
- One internal API may map to multiple external APIs.

Do not include guessed mappings without evidence.

## Step 7: Confidence scoring for `common_apis`
For each `(internal, external)` pair, compute three subscores.

### 1) semantic_score
- `1.00`: exact normalized name match
- `0.85`: synonym dictionary match (for example `Push~Enqueue`, `Pop~Dequeue`, `Delete~Remove`, `Has~Contains`)
- `0.70`: token-overlap heuristic match (CamelCase token overlap >= 0.5)
- `0.00`: otherwise

### 2) signature_score
- `1.00`: same operation class and arity-compatible signature
- `0.60`: same operation class, weak arity compatibility
- `0.00`: otherwise

Operation classes:
- mutate insert/update
- mutate remove
- read lookup
- metadata (`Len`, `Cap`, `Size`, etc.)
- iteration (`Values`, `All`, iterators)

### 3) evidence_score
- `1.00`: API on primary structure type and shown in docs/examples
- `0.70`: API on primary structure type in code
- `0.40`: helper/wrapper-level only
- `0.00`: otherwise

Final confidence:
- `confidence = round(0.60*semantic_score + 0.30*signature_score + 0.10*evidence_score, 2)`

Inclusion threshold (strict):
- Keep pair only if `confidence >= 0.70`.

Sort `common_apis` deterministically by:
1. `internal` asc
2. `external` asc

## Step 8: Keep/drop and top-3 selection
Per folder:
1. Drop any library with empty `common_apis` after threshold filtering.
2. Rank remaining libraries by `popularity_score` desc using tie-break rules.
3. Keep top 3.
4. If fewer than 3 remain, keep available set (including empty list).

## Step 9: Write output
Write `.compare/poplibs.json` with all discovered folders as keys.

Each key value is array of `LibraryMatch` objects sorted by ranking.

No extra metadata fields.

## Failure handling
- Retry each failed metric/API fetch up to 2 times.
- If one metric remains unavailable, use neutral normalized fallback `0.50` for that metric only.
- If candidate discovery is rate-limited, retry with backoff and alternate sources.

## Validation checklist
- `.compare/poplibs.json` exists and is valid JSON.
- Keys exactly match dynamically discovered in-scope folders.
- Each folder list length `<= 3`.
- Each library object has required fields.
- `common_apis` non-empty for every kept library.
- Every `confidence` in `[0.00, 1.00]` with 2 decimals.
- `common_apis` sorted deterministically.
