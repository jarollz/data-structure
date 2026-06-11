#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)

# shellcheck source=.scripts/gen-impl/helpers/common.sh
source "$SCRIPT_DIR/helpers/common.sh"
# shellcheck source=.scripts/gen-impl/helpers/folders.sh
source "$SCRIPT_DIR/helpers/folders.sh"
# shellcheck source=.scripts/gen-impl/helpers/spawner.sh
source "$SCRIPT_DIR/helpers/spawner.sh"
# shellcheck source=.scripts/gen-impl/helpers/report.sh
source "$SCRIPT_DIR/helpers/report.sh"
# shellcheck source=.scripts/gen-impl/helpers/verify.sh
source "$SCRIPT_DIR/helpers/verify.sh"

MAX_ATTEMPTS=${MAX_ATTEMPTS:-5}
REPORT_ATTEMPTS=${REPORT_ATTEMPTS:-3}
STOP_ON_FAILURE=${STOP_ON_FAILURE:-0}
FORCE=${FORCE:-0}

usage() {
  cat <<'EOF'
Usage: ./.scripts/gen-impl/gen.sh <folder|all>

Examples:
  ./.scripts/gen-impl/gen.sh list-array
  ./.scripts/gen-impl/gen.sh all

Environment:
  FORCE=1             Ignore fresh SUCCESS report and rerun folder.
  STOP_ON_FAILURE=1   Stop after first failed folder in all mode.
  MAX_ATTEMPTS=5      Maximum implementation attempts per folder.
EOF
}

append_failure_cause() {
  local file="$1"
  local cause="$2"
  local evidence="$3"
  local suggestion="$4"

  printf '%s\t%s\t%s\n' "$cause" "$evidence" "$suggestion" >>"$file"
}

write_attempt_summary() {
  local output="$1"
  local folder="$2"
  local attempt="$3"
  local command_exit="$4"
  local attempt_dir="$5"
  local doc_status="$6"
  local test_status="$7"
  local bench_status="$8"

  cat >"$output" <<EOF
# Attempt Summary

- Folder: $folder
- Attempt: $attempt of $MAX_ATTEMPTS
- AI command exit status: $command_exit
- Protected file violations: $(count_lines "$attempt_dir/protected_violations.tsv")
- Out-of-scope violations: $(count_lines "$attempt_dir/scope_violations.tsv")
- Doc comment audit: $doc_status
- Unit tests: $test_status
- Benchmarks: $bench_status

## Evidence Files

- AI output log: `$attempt_dir/ai_output.log`
- Protected file violations: `$attempt_dir/protected_violations.tsv`
- Out-of-scope violations: `$attempt_dir/scope_violations.tsv`
- Doc comment audit log: `$attempt_dir/doc_comment_audit.log`
- Unit test log: `$attempt_dir/tests.log`
- Benchmark log: `$attempt_dir/benchmarks.log`

## Retry Guidance

- Read the failing evidence files above before changing code again.
- Re-read `AGENTS.md`, `STRUCTURE-OVERVIEW.md`, `$folder/SPECS.md`, `$folder/api_contract.go`, and all tests.
- Keep doc comments aligned with the implementation and specs.
EOF
}

write_report_input() {
  local output="$1"
  local status="$2"
  local folder="$3"
  local attempts_used="$4"
  local change_table_file="$5"
  local failure_causes_file="$6"
  local folder_run_dir="$7"

  cat >"$output" <<EOF
# Report Facts

Operation Status: $status
Folder: $folder
Attempts Used: $attempts_used

## Files Changed
EOF
  cat "$change_table_file" >>"$output"
  printf '\n' >>"$output"

  if [ "$status" = 'FAILURE' ]; then
    cat >>"$output" <<'EOF'
## Failure Causes

| Cause | Evidence | Suggestion |
|---|---|---|
EOF
    if [ -s "$failure_causes_file" ]; then
      while IFS=$'\t' read -r cause evidence suggestion; do
        printf '| %s | %s | %s |\n' "$cause" "$evidence" "$suggestion" >>"$output"
      done <"$failure_causes_file"
    else
      printf '| Unknown failure | `%s/summary.md` | Review the latest attempt logs and rerun after fixing the root cause. |\n' "$folder_run_dir" >>"$output"
    fi
    printf '\n' >>"$output"
  fi

  cat >>"$output" <<EOF
## Run Artifacts

- Folder run directory: `$folder_run_dir`
- Folder summary: `$folder_run_dir/summary.md`
EOF
}

write_folder_summary() {
  local output="$1"
  local folder="$2"
  local status="$3"
  local attempts_used="$4"
  local report_path="$5"
  local note="$6"

  cat >"$output" <<EOF
# Folder Summary

- Folder: $folder
- Final status: $status
- Attempts used: $attempts_used
- Report path: `$report_path`
- Note: $note
EOF
}

write_run_summary_markdown() {
  local summary_tsv="$1"
  local output="$2"

  cat >"$output" <<'EOF'
# gen-impl Summary

| Folder | Status | Attempts | Report |
|---|---|---:|---|
EOF

  while IFS=$'\t' read -r folder status attempts report_path; do
    [ -n "$folder" ] || continue
    printf '| `%s` | %s | %s | `%s` |\n' "$folder" "$status" "$attempts" "$report_path" >>"$output"
  done <"$summary_tsv"
}

print_run_summary() {
  local summary_tsv="$1"

  printf '| Folder | Status | Attempts | Report |\n'
  printf '|---|---|---:|---|\n'
  while IFS=$'\t' read -r folder status attempts report_path; do
    [ -n "$folder" ] || continue
    printf '| `%s` | %s | %s | `%s` |\n' "$folder" "$status" "$attempts" "$report_path"
  done <"$summary_tsv"
}

generate_report() {
  local folder="$1"
  local folder_run_dir="$2"
  local spawner_command="$3"
  local expected_status="$4"
  local attempts_used="$5"
  local report_input_file="$6"
  local report_path="$7"
  local prompt_template="$GEN_IMPL_ROOT/helpers/prompt_report.md.tmpl"
  local report_attempt=1

  while [ "$report_attempt" -le "$REPORT_ATTEMPTS" ]; do
    local report_dir="$folder_run_dir/report_phase_$report_attempt"
    local report_prompt="$report_dir/prompt.md"
    local report_output_log="$report_dir/ai_output.log"
    local report_snapshot="$report_dir/snapshot"
    local validation_log="$report_dir/validation.log"

    mkdir -p "$report_dir"
    render_template "$prompt_template" "$report_prompt" \
      '@@FOLDER@@' "$folder" \
      '@@REPO_ROOT@@' "$GEN_IMPL_REPO_ROOT" \
      '@@REPORT_INPUT_PATH@@' "$report_input_file" \
      '@@OUTPUT_REPORT_PATH@@' "$report_path"

    create_phase_snapshot "$report_snapshot"
    if run_spawner_command "$spawner_command" "$(read_file_contents "$report_prompt")" "$report_output_log" "${SPAWNER_TIMEOUT_SECONDS:-1800}"; then
      :
    fi

    if ! enforce_phase_scope "$folder" "$report_snapshot" 1 "$report_dir"; then
      warn "report phase $report_attempt for $folder edited protected or out-of-scope files; restored snapshot content"
    fi

    if validate_generated_report "$report_path" "$expected_status" "$folder" "$attempts_used" >"$validation_log" 2>&1; then
      return 0
    fi

    report_attempt=$((report_attempt + 1))
  done

  return 1
}

process_folder() {
  local folder="$1"
  local run_root="$2"
  local spawner_command="$3"
  local summary_tsv="$4"

  local folder_run_dir="$run_root/$folder"
  local summary_file="$folder_run_dir/summary.md"
  local report_path
  local report_status_current
  local prior_attempts
  local prior_report_copy="$folder_run_dir/prior_IMPLEMENTATION_REPORT.md"
  local prior_suggestions_file="$folder_run_dir/prior_failure_suggestions.md"
  local empty_previous_summary="$folder_run_dir/no_previous_failure_summary.md"
  local folder_start_snapshot="$folder_run_dir/folder_start_snapshot"
  local change_table_file="$folder_run_dir/change_table.md"
  local failure_causes_file="$folder_run_dir/failure_causes.tsv"
  local report_input_file="$folder_run_dir/report_input.md"
  local last_failure_summary="$empty_previous_summary"
  local attempts_used=0
  local status='FAILURE'
  local note='implementation attempts exhausted'
  local need_reset=0

  require_folder_layout "$folder"
  mkdir -p "$folder_run_dir"
  : >"$failure_causes_file"
  printf 'No previous attempt summary exists for this run.\n' >"$empty_previous_summary"

  report_path=$(folder_report_path "$folder")
  report_status_current=$(report_status_from_file "$report_path")
  prior_attempts=$(report_attempts_from_file "$report_path")

  if [ -f "$report_path" ]; then
    cp "$report_path" "$prior_report_copy"
  else
    printf 'No prior implementation report found.\n' >"$prior_report_copy"
  fi
  extract_failure_suggestions "$report_path" "$prior_suggestions_file"

  if [ "$FORCE" = '0' ] && [ "$report_status_current" = 'SUCCESS' ] && report_is_fresh "$folder" "$report_path"; then
    status='SKIPPED'
    attempts_used="$prior_attempts"
    note='fresh SUCCESS report already exists'
    write_folder_summary "$summary_file" "$folder" "$status" "$attempts_used" "$report_path" "$note"
    printf '%s\t%s\t%s\t%s\n' "$folder" "$status" "$attempts_used" "$report_path" >>"$summary_tsv"
    return 0
  fi

  create_phase_snapshot "$folder_start_snapshot"

  case "$report_status_current" in
    FAILURE|INVALID|MISSING)
      if folder_has_non_test_impl_files "$folder"; then
        need_reset=1
      fi
      ;;
  esac

  if [ "$need_reset" -eq 1 ]; then
    local reset_dir="$folder_run_dir/reset_phase"
    local reset_prompt="$reset_dir/prompt.md"
    local reset_output="$reset_dir/ai_output.log"
    local reset_snapshot="$reset_dir/snapshot"

    mkdir -p "$reset_dir"
    render_template "$GEN_IMPL_ROOT/helpers/prompt_reset.md.tmpl" "$reset_prompt" \
      '@@FOLDER@@' "$folder" \
      '@@REPO_ROOT@@' "$GEN_IMPL_REPO_ROOT" \
      '@@PRIOR_REPORT_PATH@@' "$prior_report_copy" \
      '@@PRIOR_SUGGESTIONS_PATH@@' "$prior_suggestions_file"

    create_phase_snapshot "$reset_snapshot"
    if ! run_spawner_command "$spawner_command" "$(read_file_contents "$reset_prompt")" "$reset_output" "${SPAWNER_TIMEOUT_SECONDS:-1800}"; then
      append_failure_cause "$failure_causes_file" 'Reset phase failed' "`$reset_output`" 'Check the reset prompt and spawner command, then rerun the script.'
    fi
    if ! enforce_phase_scope "$folder" "$reset_snapshot" 0 "$reset_dir"; then
      append_failure_cause "$failure_causes_file" 'Reset phase touched protected or out-of-scope files' "`$reset_dir/protected_violations.tsv` and `$reset_dir/scope_violations.tsv`" 'Keep reset changes inside target implementation files only.'
    fi
  fi

  local attempt=1
  while [ "$attempt" -le "$MAX_ATTEMPTS" ]; do
    local attempt_dir="$folder_run_dir/attempt_$attempt"
    local prompt_file="$attempt_dir/prompt.md"
    local output_log="$attempt_dir/ai_output.log"
    local phase_snapshot="$attempt_dir/snapshot"
    local doc_log="$attempt_dir/doc_comment_audit.log"
    local test_log="$attempt_dir/tests.log"
    local bench_log="$attempt_dir/benchmarks.log"
    local command_exit=0
    local phase_ok=1
    local doc_status='SKIPPED'
    local test_status='SKIPPED'
    local bench_status='SKIPPED'

    mkdir -p "$attempt_dir"
    render_template "$GEN_IMPL_ROOT/helpers/prompt_impl.md.tmpl" "$prompt_file" \
      '@@FOLDER@@' "$folder" \
      '@@REPO_ROOT@@' "$GEN_IMPL_REPO_ROOT" \
      '@@ATTEMPT@@' "$attempt" \
      '@@MAX_ATTEMPTS@@' "$MAX_ATTEMPTS" \
      '@@PRIOR_REPORT_PATH@@' "$prior_report_copy" \
      '@@PRIOR_SUGGESTIONS_PATH@@' "$prior_suggestions_file" \
      '@@FAILURE_SUMMARY_PATH@@' "$last_failure_summary" \
      '@@RUN_DIR@@' "$folder_run_dir"

    create_phase_snapshot "$phase_snapshot"
    if run_spawner_command "$spawner_command" "$(read_file_contents "$prompt_file")" "$output_log" "${SPAWNER_TIMEOUT_SECONDS:-1800}"; then
      command_exit=0
    else
      command_exit=1
      phase_ok=0
      append_failure_cause "$failure_causes_file" 'AI command exited non-zero' "`$output_log`" 'Fix the spawner command or prompt so the AI run finishes successfully.'
    fi

    if ! enforce_phase_scope "$folder" "$phase_snapshot" 0 "$attempt_dir"; then
      phase_ok=0
      if [ -s "$attempt_dir/protected_violations.tsv" ]; then
        append_failure_cause "$failure_causes_file" 'Protected files were modified' "`$attempt_dir/protected_violations.tsv`" 'Do not edit tests, specs, docs, api_contract.go, go.mod, or go.sum.'
      fi
      if [ -s "$attempt_dir/scope_violations.tsv" ]; then
        append_failure_cause "$failure_causes_file" 'Out-of-scope files were modified' "`$attempt_dir/scope_violations.tsv`" 'Keep all edits inside the target folder implementation files only.'
      fi
    fi

    if [ "$phase_ok" -eq 1 ]; then
      if run_doc_comment_audit "$folder" "$doc_log"; then
        doc_status='PASS'
      else
        phase_ok=0
        doc_status='FAIL'
        append_failure_cause "$failure_causes_file" 'Doc comment audit failed' "`$doc_log`" 'Align exported function comments with SPECS.md, API header rules, and actual implementation behavior.'
      fi
    fi

    if [ "$phase_ok" -eq 1 ]; then
      if run_make_with_timeout "${TEST_TIMEOUT_SECONDS:-900}" "$test_log" test-folder FOLDER="$folder"; then
        test_status='PASS'
      else
        phase_ok=0
        test_status='FAIL'
        append_failure_cause "$failure_causes_file" 'Unit tests failed' "`$test_log`" 'Read the failing tests carefully and align implementation behavior with the fixed test contract.'
      fi
    fi

    if [ "$phase_ok" -eq 1 ]; then
      if run_make_with_timeout "${BENCH_TIMEOUT_SECONDS:-1800}" "$bench_log" bench-folder FOLDER="$folder"; then
        bench_status='PASS'
      else
        phase_ok=0
        bench_status='FAIL'
        append_failure_cause "$failure_causes_file" 'Benchmark checks failed' "`$bench_log`" 'Review benchmark failures and adjust implementation until benchmark validation passes.'
      fi
    fi

    write_attempt_summary "$attempt_dir/summary.md" "$folder" "$attempt" "$command_exit" "$attempt_dir" "$doc_status" "$test_status" "$bench_status"
    last_failure_summary="$attempt_dir/summary.md"
    attempts_used="$attempt"

    if [ "$phase_ok" -eq 1 ]; then
      status='SUCCESS'
      note='implementation and verification succeeded'
      break
    fi

    note='last attempt failed'
    attempt=$((attempt + 1))
  done

  write_change_table_markdown "$folder" "$folder_start_snapshot" "$change_table_file"
  write_report_input "$report_input_file" "$status" "$folder" "$attempts_used" "$change_table_file" "$failure_causes_file" "$folder_run_dir"

  if generate_report "$folder" "$folder_run_dir" "$spawner_command" "$status" "$attempts_used" "$report_input_file" "$report_path"; then
    write_change_table_markdown "$folder" "$folder_start_snapshot" "$change_table_file"
    write_report_input "$report_input_file" "$status" "$folder" "$attempts_used" "$change_table_file" "$failure_causes_file" "$folder_run_dir"
    if ! generate_report "$folder" "$folder_run_dir" "$spawner_command" "$status" "$attempts_used" "$report_input_file" "$report_path"; then
      status='FAILURE'
      note='report finalization failed'
      append_failure_cause "$failure_causes_file" 'Implementation report finalization failed' "`$folder_run_dir/report_phase_*/validation.log`" 'Fix the report prompt or spawner command, then rerun the script.'
      write_report_input "$report_input_file" "$status" "$folder" "$attempts_used" "$change_table_file" "$failure_causes_file" "$folder_run_dir"
      if ! generate_report "$folder" "$folder_run_dir" "$spawner_command" "$status" "$attempts_used" "$report_input_file" "$report_path"; then
        note='report generation failed'
      fi
    fi
  else
    status='FAILURE'
    note='report generation failed'
    append_failure_cause "$failure_causes_file" 'Implementation report generation failed' "`$folder_run_dir/report_phase_*/validation.log`" 'Fix the report prompt or spawner command, then rerun the script.'
    write_report_input "$report_input_file" "$status" "$folder" "$attempts_used" "$change_table_file" "$failure_causes_file" "$folder_run_dir"
    generate_report "$folder" "$folder_run_dir" "$spawner_command" "$status" "$attempts_used" "$report_input_file" "$report_path" || true
  fi

  write_folder_summary "$summary_file" "$folder" "$status" "$attempts_used" "$report_path" "$note"
  printf '%s\t%s\t%s\t%s\n' "$folder" "$status" "$attempts_used" "$report_path" >>"$summary_tsv"

  [ "$status" != 'FAILURE' ]
}

main() {
  local target="$1"
  local run_root
  local summary_tsv
  local summary_md
  local overall_status=0
  local spawner_command
  local folder
  local targets=()

  ensure_runtime_root
  run_root="$GEN_IMPL_REPO_ROOT/tmp/gen-impl/runs/$(timestamp_utc)"
  mkdir -p "$run_root"

  while IFS= read -r folder; do
    [ -n "$folder" ] || continue
    targets+=("$folder")
  done < <(resolve_targets "$target")

  spawner_command=$(prompt_for_spawner_command "$run_root")

  summary_tsv="$run_root/summary.tsv"
  : >"$summary_tsv"

  for folder in "${targets[@]}"; do
    log "Processing $folder"
    if ! process_folder "$folder" "$run_root" "$spawner_command" "$summary_tsv"; then
      overall_status=1
      if [ "$target" = 'all' ] && [ "$STOP_ON_FAILURE" = '1' ]; then
        break
      fi
    fi
  done

  summary_md="$run_root/summary.md"
  write_run_summary_markdown "$summary_tsv" "$summary_md"
  print_run_summary "$summary_tsv"
  log "Run artifacts: $run_root"

  exit "$overall_status"
}

if [ "$#" -ne 1 ]; then
  usage
  exit 1
fi

main "$1"
