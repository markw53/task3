from __future__ import annotations

import collections
import time
from dataclasses import dataclass, field
from typing import Deque, Dict, Tuple

from .config import LimiterConfig, Mode


@dataclass
class FixedWindowLimiter:
    cfg: LimiterConfig
    window_start: float = field(default_factory=lambda: 0.0)
    count: int = 0

    def allow(self, ts: float) -> Tuple[bool, int, float]:
        if self.window_start == 0.0:
            self.window_start = ts
        window_end = self.window_start + self.cfg.window_seconds
        if ts >= window_end:
            self.window_start = ts
            self.count = 0
            window_end = self.window_start + self.cfg.window_seconds
        if self.count < self.cfg.max_requests:
            self.count += 1
            return True, self.count, window_end
        return False, self.count, window_end


@dataclass
class SlidingWindowLimiter:
    cfg: LimiterConfig
    events: Deque[float] = field(default_factory=collections.deque)

    def allow(self, ts: float) -> Tuple[bool, int, float]:
        window_start = ts - self.cfg.window_seconds
        while self.events and self.events[0] <= window_start:
            self.events.popleft()
        if len(self.events) < self.cfg.max_requests:
            self.events.append(ts)
            reset = self.events[0] + self.cfg.window_seconds
            return True, len(self.events), reset
        reset = self.events[0] + self.cfg.window_seconds
        return False, len(self.events), reset


class RateLimiterEngine:
    def __init__(self, cfgs: Dict[str, LimiterConfig]) -> None:
        self._limiters: Dict[str, object] = {}
        for name, cfg in cfgs.items():
            if cfg.mode == Mode.FIXED:
                self._limiters[name] = FixedWindowLimiter(cfg)
            else:
                self._limiters[name] = SlidingWindowLimiter(cfg)

    def allow(self, limiter_name: str, ts: float) -> Tuple[bool, int, float]:
        limiter = self._limiters.get(limiter_name)
        if limiter is None:
            # unknown limiter: treat as allowed with dummy window
            return True, 0, ts
        return limiter.allow(ts)
