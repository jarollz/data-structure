from __future__ import annotations

import re
import shlex
import subprocess
import threading
import time
from collections.abc import Callable
from pathlib import Path

from .io_safe import warn
from .process_control import terminate_timed_out_process


ANSI_PATTERN = re.compile(r"\x1b\[[0-?]*[ -~]*[@-~]")


def validate_prompt_placeholder(command_template: str) -> tuple[bool, str]:
    in_single = False
    in_double = False
    escaped = False
    count = 0
    i = 0
    length = len(command_template)

    while i < length:
        char = command_template[i]

        if escaped:
            escaped = False
            i += 1
            continue

        if not in_single and char == "\\":
            escaped = True
            i += 1
            continue

        if char == "'" and not in_double:
            in_single = not in_single
            i += 1
            continue

        if char == '"' and not in_single:
            in_double = not in_double
            i += 1
            continue

        if command_template[i : i + 8] == "[prompt]":
            if in_single or in_double:
                return False, "[prompt] must be unquoted"
            count += 1
            i += 8
            continue

        i += 1

    if count == 0:
        return False, "missing [prompt] placeholder"
    if count > 1:
        return False, "[prompt] must appear exactly once"
    return True, ""


def validate_spawner_command_syntax(command_template: str) -> tuple[bool, str]:
    if not command_template:
        return False, "empty command"
    return validate_prompt_placeholder(command_template)


def render_spawner_command(command_template: str, prompt_text: str) -> str:
    return command_template.replace("[prompt]", shlex.quote(prompt_text))


def run_spawner_command(
    command_template: str,
    prompt_text: str,
    output_file: Path,
    timeout_seconds: int,
    repo_root: Path,
    progress_callback: Callable[[float, str | None], None] | None = None,
) -> bool:
    rendered = render_spawner_command(command_template, prompt_text)
    output_file.parent.mkdir(parents=True, exist_ok=True)
    with output_file.open("w", encoding="utf-8") as fh:
        proc = subprocess.Popen(
            ["bash", "-lc", rendered],
            cwd=repo_root,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            bufsize=1,
            start_new_session=True,
        )

        tail: dict[str, str | None] = {"value": None}

        def _reader() -> None:
            if proc.stdout is None:
                return
            try:
                for chunk in proc.stdout:
                    fh.write(chunk)
                    fh.flush()
                    for piece in chunk.replace("\r", "\n").split("\n"):
                        if piece.strip():
                            tail["value"] = piece.strip()
            finally:
                if proc.stdout is not None:
                    proc.stdout.close()

        reader = threading.Thread(target=_reader, daemon=True)
        reader.start()

        start = time.monotonic()
        while True:
            elapsed = time.monotonic() - start
            if progress_callback is not None:
                progress_callback(elapsed, tail["value"])

            code = proc.poll()
            if code is not None:
                reader.join()
                return code == 0

            if elapsed >= timeout_seconds:
                terminate_timed_out_process(proc, reader, fh, "[gen-impl] spawner timeout")
                return False

            time.sleep(0.2)


def _last_non_empty_line(path: Path) -> str:
    last = ""
    for line in path.read_text(encoding="utf-8", errors="replace").splitlines():
        if line.strip():
            last = line.strip()
    return last


def strip_ansi_to_file(input_file: Path, output_file: Path) -> None:
    text = input_file.read_text(encoding="utf-8", errors="replace")
    text = ANSI_PATTERN.sub("", text).replace("\r", "\n")
    output_file.write_text(text, encoding="utf-8")


def probe_spawner_command(
    command_template: str,
    probe_dir: Path,
    timeout_seconds: int,
    repo_root: Path,
) -> bool:
    raw_output = probe_dir / "probe_raw.log"
    clean_output = probe_dir / "probe_clean.log"
    prompt = "Reply with exactly this single word and nothing else: Hello"
    probe_dir.mkdir(parents=True, exist_ok=True)

    if not run_spawner_command(command_template, prompt, raw_output, timeout_seconds, repo_root):
        return False
    strip_ansi_to_file(raw_output, clean_output)
    return _last_non_empty_line(clean_output) == "Hello"


def prompt_for_spawner_command(
    run_root: Path,
    repo_root: Path,
    probe_timeout_seconds: int,
    preset_command: str = "",
) -> str:
    if preset_command:
        ok, reason = validate_spawner_command_syntax(preset_command)
        if not ok:
            raise RuntimeError(f"AI_SPAWNER_COMMAND invalid: {reason}")

        probe_dir = run_root / "spawner_probe_env"
        if probe_spawner_command(preset_command, probe_dir, probe_timeout_seconds, repo_root):
            return preset_command
        raise RuntimeError(
            "AI_SPAWNER_COMMAND probe command failed or final output was not Hello. "
            f"Review {probe_dir / 'probe_clean.log'}"
        )

    attempt = 1
    while True:
        try:
            print("AI agent spawner command (must have [prompt] in it):")
            command_template = input().strip()
        except EOFError as exc:
            raise RuntimeError("failed to read spawner command from stdin") from exc

        ok, reason = validate_spawner_command_syntax(command_template)
        if not ok:
            warn(reason)
            attempt += 1
            continue

        probe_dir = run_root / f"spawner_probe_{attempt}"
        if probe_spawner_command(command_template, probe_dir, probe_timeout_seconds, repo_root):
            return command_template

        warn(f"probe command failed or final output was not Hello. Review {probe_dir / 'probe_clean.log'} and try again.")
        attempt += 1
