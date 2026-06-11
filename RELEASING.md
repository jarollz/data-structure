# Releasing and Tagging

## Scope

This repository is a monorepo with many Go modules (one module per structure folder).
Each module has its own semantic version stream.

## Tag format

- Submodule release tag MUST be `<folder>/vX.Y.Z`.
- Example: `stack/v1.2.3`, `list-array/v0.4.0`.
- Root tag `vX.Y.Z` is only valid if a root `go.mod` exists.
- Current state: no root `go.mod`, so use submodule tags only.

## Semantic version rules

- Bump `PATCH` for backward-compatible bug fixes.
- Bump `MINOR` for backward-compatible feature additions.
- Bump `MAJOR` for breaking API/behavior changes.
- Versioning is per module folder, not global for repo.

## Major version path rule (`v2+`)

- If tag major is `v2` or higher, module path in `<folder>/go.mod` MUST end with `/vN`.
- Example:
  - tag: `stack/v2.0.0`
  - module path: `github.com/jarollz/data-structure/stack/v2`
- Tag major and module path major MUST match.

## Tag immutability policy

- Published release tags are immutable.
- Do not move, rewrite, or force-push an existing release tag.
- If a release is wrong, publish a new version (for example `v1.1.2`) instead of retagging `v1.1.1`.

## Release checklist

1. Run tests for target module folder (for example `make test-folder FOLDER=stack`).
2. Create and push tag with guard:
   - `make tag FOLDER=stack VERSION=v1.2.3`
   - or shorthand: `make tag stack v1.2.3`
3. Verify tag is visible on remote:
   - `git ls-remote --tags origin "refs/tags/stack/v1.2.3"`
