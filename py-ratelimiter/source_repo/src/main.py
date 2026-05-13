import argparse
import asyncio
import sys
from typing import Optional

from .config import parse_config
from .engine import RateLimiterApp
from .utils.logging import log
from .utils.validation import ValidationError


def main(argv: Optional[list[str]] = None) -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--config", required=True, help="JSON config file")
    args = parser.parse_args(argv)

    try:
        with open(args.config, "r", encoding="utf-8") as f:
            raw = f.read()
    except OSError as exc:
        log("config_error", error=str(exc))
        return 1

    try:
        app_cfg = parse_config(raw)
    except ValidationError as exc:
        log("config_error", error=str(exc))
        return 1

    app = RateLimiterApp(app_cfg)
    try:
        return asyncio.run(app.process_stream(sys.stdin))
    except Exception as exc:  # noqa: BLE001
        log("runtime_error", error=str(exc))
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
