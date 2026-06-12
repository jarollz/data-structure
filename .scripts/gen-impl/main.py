#!/usr/bin/env python3

from __future__ import annotations

import sys

from core.config import Config
from core.io_safe import die
from core.runner import run


def usage() -> str:
    return """Usage: ./.scripts/gen-impl/gen.sh <folder|all>

Examples:
  ./.scripts/gen-impl/gen.sh list-array
  ./.scripts/gen-impl/gen.sh all

Environment:
  FORCE=0             Default. Skip fresh SUCCESS; otherwise rerun clean-first.
  FORCE=1             Always rerun and clean-first reset.
  FORCE=2             Rerun in-place when prior Performance Status is FAIL.
  STOP_ON_FAILURE=1   Stop after first failed folder in all mode.
  MAX_ATTEMPTS=5      Maximum implementation attempts per folder.
  SPAWNER_TIMEOUT_SECONDS
                      Hard timeout for each AI spawner run.
  SPAWNER_IDLE_TIMEOUT_SECONDS
                      Fail AI spawner run when output stays idle.
  AI_SPAWNER_COMMAND  Full spawner command template with [prompt].
"""


def main(argv: list[str]) -> int:
    if len(argv) != 2:
        print(usage(), file=sys.stderr)
        return 1

    target = argv[1]
    try:
        cfg = Config.from_env()
        return run(cfg, target)
    except RuntimeError as exc:
        die(str(exc))
        return 1


if __name__ == "__main__":
    raise SystemExit(main(sys.argv))
