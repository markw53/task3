from __future__ import annotations

import asyncio
import json
import sys
from dataclasses import dataclass
from typing import Any, Dict

from .config import AppConfig
from .limiter import RateLimiterEngine
from .utils.logging import log
from .utils.validation import ValidationError, require_field


@dataclass
class Request:
    user: str
    endpoint: str
    ts: float


def parse_request(line: str) -> Request:
    try:
        data = json.loads(line)
    except json.JSONDecodeError as exc:
        raise ValidationError(f"Invalid request JSON: {exc}") from exc
    user = require_field(data, "user")
    endpoint = require_field(data, "endpoint")
    ts = require_field(data, "ts")
    if not isinstance(user, str):
        raise ValidationError("user must be a string")
    if not isinstance(endpoint, str):
        raise ValidationError("endpoint must be a string")
    if not isinstance(ts, (int, float)):
        raise ValidationError("ts must be a number")
    return Request(user=user, endpoint=endpoint, ts=float(ts))


class RateLimiterApp:
    def __init__(self, cfg: AppConfig) -> None:
        self.cfg = cfg
        self.engine = RateLimiterEngine(cfg.limiters)

    async def process_stream(self, stream) -> int:
        loop = asyncio.get_running_loop()
        tasks = []
        for line in stream:
            line = line.strip()
            if not line:
                continue
            tasks.append(loop.create_task(self._handle_line(line)))
        if tasks:
            await asyncio.gather(*tasks)
        return 0

    async def _handle_line(self, line: str) -> None:
        try:
            req = parse_request(line)
        except ValidationError as exc:
            log("request_error", error=str(exc))
            return
        allowed, count, reset = self.engine.allow(req.endpoint, req.ts)
        log(
            "decision",
            user=req.user,
            endpoint=req.endpoint,
            ts=req.ts,
            allowed=allowed,
            count=count,
            reset=reset,
        )
