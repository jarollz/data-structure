#!/usr/bin/env sh
set -eu

usage() {
  echo "Usage: ./.scripts/release-tag.sh <folder> <vX.Y.Z>"
  echo "Example: ./.scripts/release-tag.sh stack v1.2.3"
}

if [ "$#" -ne 2 ]; then
  usage
  exit 1
fi

folder="$1"
version="$2"
tag="$folder/$version"

if [ ! -d "$folder" ]; then
  echo "Error: folder does not exist: $folder"
  exit 1
fi

if [ ! -f "$folder/go.mod" ]; then
  echo "Error: missing go.mod in folder: $folder"
  exit 1
fi

if ! printf '%s\n' "$version" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+$'; then
  echo "Error: VERSION must match vX.Y.Z, got: $version"
  exit 1
fi

module_path="$(awk '/^module /{print $2; exit}' "$folder/go.mod")"
if [ -z "$module_path" ]; then
  echo "Error: cannot read module path from $folder/go.mod"
  exit 1
fi

tag_major="$(printf '%s' "$version" | cut -d. -f1 | tr -d 'v')"
module_major="$(printf '%s\n' "$module_path" | sed -n 's#.*/v\([0-9][0-9]*\)$#\1#p')"
if [ -z "$module_major" ]; then
  module_major="1"
fi

if [ "$module_major" -ne "$tag_major" ]; then
  echo "Error: major mismatch. Tag=$version (major $tag_major), module=$module_path (major $module_major)"
  exit 1
fi

if git rev-parse -q --verify "refs/tags/$tag" >/dev/null 2>&1; then
  echo "Error: local tag already exists: $tag"
  exit 1
fi

if git ls-remote --exit-code --tags origin "refs/tags/$tag" >/dev/null 2>&1; then
  echo "Error: remote tag already exists on origin: $tag"
  exit 1
fi

if [ "${DRY_RUN:-0}" = "1" ]; then
  echo "DRY_RUN=1"
  echo "Would run: git tag $tag"
  echo "Would run: git push origin refs/tags/$tag"
  exit 0
fi

git tag "$tag"
git push origin "refs/tags/$tag"
echo "Created and pushed tag: $tag"
