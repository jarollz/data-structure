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
  FORCE=1             Ignore fresh SUCCESS report and rerun folder.
  STOP_ON_FAILURE=1   Stop after first failed folder in all mode.
  MAX_ATTEMPTS=5      Maximum implementation attempts per folder.
"""


def main(argv: list[str]) -> int:
    if len(argv) != 2:
        print(usage(), file=sys.stderr)
        return 1

    target = argv[1]
    cfg = Config.from_env()
    try:
        return run(cfg, target)
    except RuntimeError as exc:
        die(str(exc))
        return 1


if __name__ == "__main__":
    raise SystemExit(main(sys.argv))
