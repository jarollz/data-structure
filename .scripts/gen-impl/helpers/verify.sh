#!/usr/bin/env bash

create_phase_snapshot() {
  local snapshot_dir="$1"

  mkdir -p "$snapshot_dir/repo"
  (
    cd "$GEN_IMPL_REPO_ROOT"
    tar -cf - --exclude='.git' --exclude='tmp/gen-impl/runs' .
  ) | (
    cd "$snapshot_dir/repo"
    tar -xf -
  )
}

build_manifest() {
  local root_dir="$1"
  local output_file="$2"
  local relative_path
  local hash

  (
    cd "$root_dir"
    find . -type f ! -path './.git/*' ! -path './tmp/gen-impl/runs/*' -print0 |
      LC_ALL=C sort -z |
      while IFS= read -r -d '' path; do
        relative_path=${path#./}
        hash=$(shasum -a 256 "$relative_path" | awk '{ print $1 }')
        printf '%s\t%s\n' "$relative_path" "$hash"
      done
  ) >"$output_file"
}

diff_manifests() {
  local before_manifest="$1"
  local after_manifest="$2"
  local output_file="$3"

  awk -F'\t' '
    NR == FNR { before[$1] = $2; next }
    { after[$1] = $2 }
    END {
      for (path in before) {
        if (!(path in after)) {
          printf "D\t%s\n", path
        } else if (before[path] != after[path]) {
          printf "M\t%s\n", path
        }
      }
      for (path in after) {
        if (!(path in before)) {
          printf "A\t%s\n", path
        }
      }
    }
  ' "$before_manifest" "$after_manifest" | LC_ALL=C sort >"$output_file"
}

path_is_protected() {
  local path="$1"
  local folder="$2"
  local allow_report_write="$3"
  local base_name

  case "$path" in
    "$folder"/*)
      ;;
    *)
      return 1
      ;;
  esac

  base_name=${path##*/}
  case "$base_name" in
    api_contract.go|go.mod|go.sum|SPECS.md|*_test.go|*_bench_test.go|helpers_test.go|bench_policy_test.go)
      return 0
      ;;
    *.md)
      if [ "$base_name" = 'IMPLEMENTATION_REPORT.md' ] && [ "$allow_report_write" -eq 1 ]; then
        return 1
      fi
      return 0
      ;;
  esac

  return 1
}

restore_path_from_snapshot() {
  local snapshot_dir="$1"
  local relative_path="$2"
  local snapshot_file="$snapshot_dir/repo/$relative_path"
  local repo_file="$GEN_IMPL_REPO_ROOT/$relative_path"

  if [ -f "$snapshot_file" ]; then
    mkdir -p "$(dirname "$repo_file")"
    cp "$snapshot_file" "$repo_file"
    return 0
  fi

  rm -f "$repo_file"
}

enforce_phase_scope() {
  local folder="$1"
  local snapshot_dir="$2"
  local allow_report_write="$3"
  local result_dir="$4"
  local before_manifest="$result_dir/before_manifest.tsv"
  local after_manifest="$result_dir/after_manifest.tsv"
  local changes_manifest="$result_dir/changes.tsv"
  local status=0
  local change_type
  local path

  mkdir -p "$result_dir"
  : >"$result_dir/kept_changes.tsv"
  : >"$result_dir/protected_violations.tsv"
  : >"$result_dir/scope_violations.tsv"

  build_manifest "$snapshot_dir/repo" "$before_manifest"
  build_manifest "$GEN_IMPL_REPO_ROOT" "$after_manifest"
  diff_manifests "$before_manifest" "$after_manifest" "$changes_manifest"

  while IFS=$'\t' read -r change_type path; do
    [ -n "$path" ] || continue
    case "$path" in
      tmp/gen-impl/runs/*)
        continue
        ;;
    esac

    case "$path" in
      "$folder"/*)
        if path_is_protected "$path" "$folder" "$allow_report_write"; then
          printf '%s\t%s\n' "$change_type" "$path" >>"$result_dir/protected_violations.tsv"
          restore_path_from_snapshot "$snapshot_dir" "$path"
          status=1
        else
          printf '%s\t%s\n' "$change_type" "$path" >>"$result_dir/kept_changes.tsv"
        fi
        ;;
      *)
        printf '%s\t%s\n' "$change_type" "$path" >>"$result_dir/scope_violations.tsv"
        restore_path_from_snapshot "$snapshot_dir" "$path"
        status=1
        ;;
    esac
  done <"$changes_manifest"

  return "$status"
}

write_change_table_markdown() {
  local folder="$1"
  local start_snapshot="$2"
  local output_file="$3"
  local before_manifest="$output_file.before.tsv"
  local after_manifest="$output_file.after.tsv"
  local changes_manifest="$output_file.changes.tsv"
  local count=0
  local change_type
  local path
  local label

  build_manifest "$start_snapshot/repo" "$before_manifest"
  build_manifest "$GEN_IMPL_REPO_ROOT" "$after_manifest"
  diff_manifests "$before_manifest" "$after_manifest" "$changes_manifest"

  cat >"$output_file" <<'EOF'
| File | Change |
|---|---|
EOF

  while IFS=$'\t' read -r change_type path; do
    case "$path" in
      "$folder"/*)
        ;;
      *)
        continue
        ;;
    esac

    case "$change_type" in
      A) label='added' ;;
      D) label='deleted' ;;
      *) label='modified' ;;
    esac
    printf '| `%s` | %s |\n' "$path" "$label" >>"$output_file"
    count=$((count + 1))
  done <"$changes_manifest"

  if [ "$count" -eq 0 ]; then
    printf '| `(none)` | unchanged |\n' >>"$output_file"
  fi
}

run_doc_comment_audit() {
  local folder="$1"
  local output_log="$2"

  run_with_timeout "${DOC_AUDIT_TIMEOUT_SECONDS:-300}" go run "$GEN_IMPL_ROOT/helpers/audit_doc_comments.go" --folder "$folder" >"$output_log" 2>&1
}

run_make_with_timeout() {
  local timeout_seconds="$1"
  local output_log="$2"
  shift 2
  run_with_timeout "$timeout_seconds" make -C "$GEN_IMPL_REPO_ROOT" "$@" >"$output_log" 2>&1
}
