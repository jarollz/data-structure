from __future__ import annotations

import re
import shutil
import sys
from dataclasses import dataclass


ANSI_PATTERN = re.compile(r"\x1b\[[0-?]*[ -~]*[@-~]")
CONTROL_PATTERN = re.compile(r"[\x00-\x08\x0b\x0c\x0e-\x1f\x7f]")


def _sanitize_tail(text: str) -> str:
    clean = ANSI_PATTERN.sub("", text)
    clean = clean.replace("\r", "\n")
    clean = " ".join(clean.split())
    clean = CONTROL_PATTERN.sub("", clean)
    return clean.strip()


def _truncate_to_width(text: str, width: int) -> str:
    if width <= 0:
        return ""
    if len(text) <= width:
        return text
    if width <= 3:
        return text[:width]
    return text[: width - 3] + "..."


@dataclass
class LiveProgress:
    enabled: bool = False
    _active: bool = False
    _phase: str = "waiting"
    _label: str = "ai"
    _spinner_index: int = 0

    def __post_init__(self) -> None:
        self.enabled = sys.stdout.isatty()

    def begin_phase(self, phase: str, label: str) -> None:
        self._phase = phase
        self._label = label
        self._spinner_index = 0

    def _spinner(self) -> str:
        frames = (".", "..", "...")
        frame = frames[self._spinner_index % len(frames)]
        self._spinner_index += 1
        return frame

    def _render_lines(self, elapsed_seconds: float, tail: str | None) -> tuple[str, str]:
        columns = max(shutil.get_terminal_size(fallback=(120, 24)).columns, 50)
        phase_line = f"{self._phase} | {self._spinner()} | {elapsed_seconds:.2f}s"

        if tail is None or not tail.strip():
            if self._label == "ai":
                tail_text = "(waiting output...)"
            else:
                tail_text = "(waiting output...)"
        else:
            tail_text = _sanitize_tail(tail)
            if not tail_text:
                tail_text = "(waiting output...)"

        tail_line = f"{self._label}: {tail_text}"
        return _truncate_to_width(phase_line, columns), _truncate_to_width(tail_line, columns)

    def _clear_dynamic_lines(self) -> None:
        sys.stdout.write("\r\x1b[2K\x1b[1A\r\x1b[2K")

    def update(self, elapsed_seconds: float, tail: str | None) -> None:
        if not self.enabled:
            return

        line1, line2 = self._render_lines(elapsed_seconds, tail)
        if self._active:
            self._clear_dynamic_lines()
        sys.stdout.write(f"{line1}\n{line2}")
        sys.stdout.flush()
        self._active = True

    def clear(self) -> None:
        if not self.enabled or not self._active:
            return
        self._clear_dynamic_lines()
        sys.stdout.flush()
        self._active = False
