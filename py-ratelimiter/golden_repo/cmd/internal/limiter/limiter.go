package limiter

import (
    "time"

    "ratelimiter/internal/config"
)

type FixedWindowLimiter struct {
    Cfg         config.LimiterConfig
    WindowStart float64
    Count       int
}

func NewFixedWindowLimiter(cfg config.LimiterConfig) *FixedWindowLimiter {
    return &FixedWindowLimiter{Cfg: cfg}
}

func (l *FixedWindowLimiter) Allow(ts float64) (bool, int, float64) {
    if l.WindowStart == 0 {
        l.WindowStart = ts
    }
    windowEnd := l.WindowStart + float64(l.Cfg.WindowSeconds)
    if ts >= windowEnd {
        l.WindowStart = ts
        l.Count = 0
        windowEnd = l.WindowStart + float64(l.Cfg.WindowSeconds)
    }
    if l.Count < l.Cfg.MaxRequests {
        l.Count++
        return true, l.Count, windowEnd
    }
    return false, l.Count, windowEnd
}

type SlidingWindowLimiter struct {
    Cfg    config.LimiterConfig
    Events []float64
}

func NewSlidingWindowLimiter(cfg config.LimiterConfig) *SlidingWindowLimiter {
    return &SlidingWindowLimiter{Cfg: cfg, Events: []float64{}}
}

func (l *SlidingWindowLimiter) Allow(ts float64) (bool, int, float64) {
    windowStart := ts - float64(l.Cfg.WindowSeconds)
    // evict old
    i := 0
    for _, v := range l.Events {
        if v > windowStart {
            l.Events[i] = v
            i++
        }
    }
    l.Events = l.Events[:i]
    if len(l.Events) < l.Cfg.MaxRequests {
        l.Events = append(l.Events, ts)
        reset := l.Events[0] + float64(l.Cfg.WindowSeconds)
        return true, len(l.Events), reset
    }
    reset := l.Events[0] + float64(l.Cfg.WindowSeconds)
    return false, len(l.Events), reset
}

type Engine struct {
    limiters map[string]any
}

func NewEngine(cfgs map[string]config.LimiterConfig) *Engine {
    m := make(map[string]any)
    for name, cfg := range cfgs {
        if cfg.Mode == config.ModeFixed {
            m[name] = NewFixedWindowLimiter(cfg)
        } else {
            m[name] = NewSlidingWindowLimiter(cfg)
        }
    }
    return &Engine{limiters: m}
}

func (e *Engine) Allow(name string, ts float64) (bool, int, float64) {
    l, ok := e.limiters[name]
    if !ok {
        return true, 0, ts
    }
    switch v := l.(type) {
    case *FixedWindowLimiter:
        return v.Allow(ts)
    case *SlidingWindowLimiter:
        return v.Allow(ts)
    default:
        return true, 0, ts
    }
}

// helper for tests if needed
func NowSeconds() float64 {
    return float64(time.Now().UnixNano()) / 1e9
}
