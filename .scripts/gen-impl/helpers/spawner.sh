#!/usr/bin/env bash

validate_prompt_placeholder() {
  local command_template="$1"
  local in_single=0
  local in_double=0
  local escaped=0
  local count=0
  local i=0
  local length=${#command_template}
  local char
  local segment

  while [ "$i" -lt "$length" ]; do
    char=${command_template:$i:1}

    if [ "$escaped" -eq 1 ]; then
      escaped=0
      i=$((i + 1))
      continue
    fi

    if [ "$in_single" -eq 0 ] && [ "$char" = '\\' ]; then
      escaped=1
      i=$((i + 1))
      continue
    fi

    if [ "$char" = "'" ] && [ "$in_double" -eq 0 ]; then
      if [ "$in_single" -eq 0 ]; then
        in_single=1
      else
        in_single=0
      fi
      i=$((i + 1))
      continue
    fi

    if [ "$char" = '"' ] && [ "$in_single" -eq 0 ]; then
      if [ "$in_double" -eq 0 ]; then
        in_double=1
      else
        in_double=0
      fi
      i=$((i + 1))
      continue
    fi

    segment=${command_template:$i:8}
    if [ "$segment" = '[prompt]' ]; then
      if [ "$in_single" -eq 1 ] || [ "$in_double" -eq 1 ]; then
        printf 'quoted\n'
        return 1
      fi
      count=$((count + 1))
      i=$((i + 8))
      continue
    fi

    i=$((i + 1))
  done

  if [ "$count" -eq 0 ]; then
    printf 'missing\n'
    return 1
  fi
  if [ "$count" -gt 1 ]; then
    printf 'multiple\n'
    return 1
  fi
  return 0
}

validate_spawner_command_syntax() {
  local command_template="$1"
  local result

  if [ -z "$command_template" ]; then
    printf 'empty command\n'
    return 1
  fi

  if ! result=$(validate_prompt_placeholder "$command_template"); then
    case "$result" in
      quoted)
        printf '[prompt] must be unquoted\n'
        ;;
      missing)
        printf 'missing [prompt] placeholder\n'
        ;;
      multiple)
        printf '[prompt] must appear exactly once\n'
        ;;
      *)
        printf 'invalid [prompt] placeholder usage\n'
        ;;
    esac
    return 1
  fi

  return 0
}

render_spawner_command() {
  local command_template="$1"
  local prompt_text="$2"
  local quoted_prompt

  quoted_prompt=$(shell_quote "$prompt_text")
  printf '%s\n' "${command_template//\[prompt\]/$quoted_prompt}"
}

run_spawner_command() {
  local command_template="$1"
  local prompt_text="$2"
  local output_file="$3"
  local timeout_seconds="$4"
  local rendered_command

  rendered_command=$(render_spawner_command "$command_template" "$prompt_text")
  run_with_timeout "$timeout_seconds" bash -lc 'cd "$1" && eval "$2"' bash "$GEN_IMPL_REPO_ROOT" "$rendered_command" >"$output_file" 2>&1
}

probe_spawner_command() {
  local command_template="$1"
  local probe_dir="$2"
  local raw_output="$probe_dir/probe_raw.log"
  local clean_output="$probe_dir/probe_clean.log"
  local last_line
  local probe_prompt='Reply with exactly this single word and nothing else: Hello'

  mkdir -p "$probe_dir"
  if ! run_spawner_command "$command_template" "$probe_prompt" "$raw_output" "${PROBE_TIMEOUT_SECONDS:-120}"; then
    return 1
  fi

  strip_ansi_to_file "$raw_output" "$clean_output"
  last_line=$(last_non_empty_line "$clean_output")
  [ "$last_line" = 'Hello' ]
}

prompt_for_spawner_command() {
  local run_root="$1"
  local attempt=1
  local command_template
  local reason

  while true; do
    printf 'type in the AI agent spawner command (must have "[prompt]" in it): '
    IFS= read -r command_template || die 'failed to read spawner command from stdin'

    if ! reason=$(validate_spawner_command_syntax "$command_template"); then
      warn "$reason"
      attempt=$((attempt + 1))
      continue
    fi

    local probe_dir="$run_root/spawner_probe_$attempt"
    if probe_spawner_command "$command_template" "$probe_dir"; then
      printf '%s\n' "$command_template"
      return 0
    fi

    warn "probe command failed or final output was not Hello. Review $probe_dir/probe_clean.log and try again."
    attempt=$((attempt + 1))
  done
}
