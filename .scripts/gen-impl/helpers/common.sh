#!/usr/bin/env bash

GEN_IMPL_HELPERS_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
GEN_IMPL_ROOT=$(cd "$GEN_IMPL_HELPERS_DIR/.." && pwd)
GEN_IMPL_REPO_ROOT=$(cd "$GEN_IMPL_ROOT/../.." && pwd)

log() {
  printf '[gen-impl] %s\n' "$*"
}

warn() {
  printf '[gen-impl] Warning: %s\n' "$*" >&2
}

die() {
  printf '[gen-impl] Error: %s\n' "$*" >&2
  exit 1
}

timestamp_utc() {
  date '+%Y%m%d_%H%M%S'
}

ensure_runtime_root() {
  mkdir -p "$GEN_IMPL_REPO_ROOT/tmp/gen-impl/runs"
}

file_mtime() {
  stat -f '%m' "$1"
}

shell_quote() {
  printf '%q' "$1"
}

run_with_timeout() {
  local timeout_seconds="$1"
  shift
  perl -e 'alarm shift; exec @ARGV or die $!;' "$timeout_seconds" "$@"
}

strip_ansi_to_file() {
  local input_file="$1"
  local output_file="$2"
  perl -pe 's/\e\[[0-9;?]*[ -\\/]*[@-~]//g; s/\r/\n/g' "$input_file" >"$output_file"
}

last_non_empty_line() {
  local input_file="$1"
  awk 'NF { line = $0 } END { print line }' "$input_file"
}

render_template() {
  local template_file="$1"
  local output_file="$2"
  shift 2

  cp "$template_file" "$output_file"
  while [ "$#" -gt 0 ]; do
    local placeholder="$1"
    local value="$2"
    shift 2
    GEN_IMPL_PLACEHOLDER="$placeholder" GEN_IMPL_RENDER_VALUE="$value" \
      perl -0pi -e 's/\Q$ENV{GEN_IMPL_PLACEHOLDER}\E/$ENV{GEN_IMPL_RENDER_VALUE}/g' "$output_file"
  done
}

read_file_contents() {
  local input_file="$1"
  perl -0pe '' "$input_file"
}

count_lines() {
  local input_file="$1"
  if [ ! -f "$input_file" ] || [ ! -s "$input_file" ]; then
    printf '0\n'
    return 0
  fi
  awk 'END { print NR }' "$input_file"
}
