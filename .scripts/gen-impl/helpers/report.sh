#!/usr/bin/env bash

folder_report_path() {
  local folder="$1"
  printf '%s/%s/IMPLEMENTATION_REPORT.md\n' "$GEN_IMPL_REPO_ROOT" "$folder"
}

report_field_from_file() {
  local report_file="$1"
  local field_name="$2"

  if [ ! -f "$report_file" ]; then
    return 1
  fi

  awk -F': ' -v field_name="$field_name" '$1 == field_name { print $2; exit }' "$report_file"
}

report_status_from_file() {
  local report_file="$1"
  local status

  if [ ! -f "$report_file" ]; then
    printf 'MISSING\n'
    return 0
  fi

  status=$(report_field_from_file "$report_file" 'Operation Status' || true)
  case "$status" in
    SUCCESS|FAILURE)
      printf '%s\n' "$status"
      ;;
    *)
      printf 'INVALID\n'
      ;;
  esac
}

report_attempts_from_file() {
  local report_file="$1"
  local attempts

  if [ ! -f "$report_file" ]; then
    printf '0\n'
    return 0
  fi

  attempts=$(report_field_from_file "$report_file" 'Attempts Used' || true)
  case "$attempts" in
    ''|*[!0-9]*)
      printf '0\n'
      ;;
    *)
      printf '%s\n' "$attempts"
      ;;
  esac
}

list_non_test_impl_files() {
  local folder="$1"
  local folder_dir="$GEN_IMPL_REPO_ROOT/$folder"
  local file

  for file in "$folder_dir"/*.go; do
    if [ ! -f "$file" ]; then
      continue
    fi
    case "${file##*/}" in
      api_contract.go|*_test.go|*_bench_test.go|helpers_test.go|bench_policy_test.go)
        continue
        ;;
    esac
    printf '%s\n' "$file"
  done
}

folder_has_non_test_impl_files() {
  local folder="$1"
  local file

  while IFS= read -r file; do
    [ -n "$file" ] && return 0
  done < <(list_non_test_impl_files "$folder")
  return 1
}

report_is_fresh() {
  local folder="$1"
  local report_file="$2"
  local report_mtime
  local file
  local file_mtime_value

  if [ ! -f "$report_file" ]; then
    return 1
  fi

  report_mtime=$(file_mtime "$report_file")
  while IFS= read -r file; do
    [ -n "$file" ] || continue
    file_mtime_value=$(file_mtime "$file")
    if [ "$file_mtime_value" -gt "$report_mtime" ]; then
      return 1
    fi
  done < <(list_non_test_impl_files "$folder")

  return 0
}

extract_failure_suggestions() {
  local report_file="$1"
  local output_file="$2"

  if [ ! -f "$report_file" ]; then
    printf 'No prior implementation report exists.\n' >"$output_file"
    return 0
  fi

  awk '
    /^## Failure Causes$/ { printing = 1 }
    printing && /^## / && $0 != "## Failure Causes" { exit }
    printing { print }
  ' "$report_file" >"$output_file"

  if [ ! -s "$output_file" ]; then
    printf 'No prior failure suggestions were found in %s.\n' "$report_file" >"$output_file"
  fi
}

validate_generated_report() {
  local report_file="$1"
  local expected_status="$2"
  local expected_folder="$3"
  local expected_attempts="$4"

  if [ ! -f "$report_file" ]; then
    printf 'missing report file: %s\n' "$report_file" >&2
    return 1
  fi
  if [ "$(report_status_from_file "$report_file")" != "$expected_status" ]; then
    printf 'unexpected report status in %s\n' "$report_file" >&2
    return 1
  fi
  if [ "$(report_field_from_file "$report_file" 'Folder' || true)" != "$expected_folder" ]; then
    printf 'unexpected folder field in %s\n' "$report_file" >&2
    return 1
  fi
  if [ "$(report_attempts_from_file "$report_file")" != "$expected_attempts" ]; then
    printf 'unexpected attempts field in %s\n' "$report_file" >&2
    return 1
  fi
  if ! grep -q '^| File | Change |$' "$report_file"; then
    printf 'missing file change table in %s\n' "$report_file" >&2
    return 1
  fi
  if [ "$expected_status" = 'FAILURE' ] && ! grep -q '^## Failure Causes$' "$report_file"; then
    printf 'missing failure cause table in %s\n' "$report_file" >&2
    return 1
  fi
  return 0
}
