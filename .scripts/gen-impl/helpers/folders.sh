#!/usr/bin/env bash

GEN_IMPL_FOLDERS=(
  list-array
  list-linked-singly
  list-linked-doubly
  list-skip
  queue
  stack
  heap
  tree-general
  tree-avl
  tree-red-black
  map-hash
  map-trie
  map-tree-avl
  map-tree-red-black
)

list_supported_folders() {
  printf '%s\n' "${GEN_IMPL_FOLDERS[@]}"
}

is_supported_folder() {
  local folder="$1"
  local item

  for item in "${GEN_IMPL_FOLDERS[@]}"; do
    if [ "$item" = "$folder" ]; then
      return 0
    fi
  done
  return 1
}

resolve_targets() {
  local target="$1"

  if [ "$target" = "all" ]; then
    list_supported_folders
    return 0
  fi

  if ! is_supported_folder "$target"; then
    die "unsupported folder '$target'"
  fi

  printf '%s\n' "$target"
}

require_folder_layout() {
  local folder="$1"

  if [ ! -d "$GEN_IMPL_REPO_ROOT/$folder" ]; then
    die "folder does not exist: $folder"
  fi
  if [ ! -f "$GEN_IMPL_REPO_ROOT/$folder/go.mod" ]; then
    die "missing go.mod in folder: $folder"
  fi
}
