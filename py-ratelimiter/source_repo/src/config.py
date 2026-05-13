from __future__ import annotations

import json
from dataclasses import dataclass
from enum import Enum
from typing import Dict, Any, List

from .utils.validation import ValidationError, require_field, require_positive_int


class Mode(str, Enum):
    FIXED = "fixed"
    SLIDING = "sliding"


@dataclass
class LimiterConfig:
    name: str
    mode: Mode
    window_seconds: int
    max_requests: int


@dataclass
class AppConfig:
    limiters: Dict[str, LimiterConfig]


def parse_config(raw: str) -> AppConfig:
    try:
        data = json.loads(raw)
    except json.JSONDecodeError as exc:
        raise ValidationError(f"Invalid JSON config: {exc}") from exc

    limiters_raw = require_field(data, "limiters")
    if not isinstance(limiters_raw, list):
        raise ValidationError("limiters must be a list")

    limiters: Dict[str, LimiterConfig] = {}
    for item in limiters_raw:
        if not isinstance(item, dict):
            raise ValidationError("each limiter must be an object")
        name = require_field(item, "name")
        if not isinstance(name, str):
            raise ValidationError("name must be a string")
        mode_raw = require_field(item, "mode")
        if mode_raw not in (Mode.FIXED.value, Mode.SLIDING.value):
            raise ValidationError(f"invalid mode: {mode_raw}")
        mode = Mode(mode_raw)
        window_seconds = require_positive_int(item, "window_seconds")
        max_requests = require_positive_int(item, "max_requests")
        limiters[name] = LimiterConfig(
            name=name,
            mode=mode,
            window_seconds=window_seconds,
            max_requests=max_requests,
        )

    if not limiters:
        raise ValidationError("at least one limiter required")

    return AppConfig(limiters=limiters)
