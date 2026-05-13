package engine

import (
    "bufio"
    "context"
    "encoding/json"
    "io"
    "strings"

    "ratelimiter/internal/config"
    "ratelimiter/internal/limiter"
    "ratelimiter/internal/utils"
)

type Request struct {
    User     string  `json:"user"`
    Endpoint string  `json:"endpoint"`
    Ts       float64 `json:"ts"`
}

type RateLimiterApp struct {
    cfg    *config.AppConfig
    engine *Engine
}

func NewRateLimiterApp(cfg *config.AppConfig) *RateLimiterApp {
    return &RateLimiterApp{
        cfg:    cfg,
        engine: NewEngine(copyLimiterCfg(cfg.Limiters))
    }
}

func copyLimiterCfg(src map[string]config.LimiterConfig) map[string]config.LimiterConfig {
    dst := make(map[string]config.LimiterConfig)
    for k, v := range src {
        dst[k] = v
    }
    return dst
}

func parseRequest(line string) (*Request, error) {
    var m map[string]any
    if err := json.Unmarshal([]byte(line), &m); err != nil {
        return nil, &utils.ValidationError{Msg: "Invalid request JSON: " + err.Error()}
    }

    userAny, err := utils.RequireField(m, "user")
    if err != nil {
        return nil, err
    }
    endpointAny, err := utils.RequireField(m, "endpoint")
    if err != nil {
        return nil, err
    }
    tsAny, err := utils.RequireField(m, "ts")
    if err != nil {
        return nil, err
    }

    user, ok := userAny.(string)
    if !ok {
        return nil, &utils.ValidationError{Msg: "user must be a string"}
    }
    endpoint, ok := endpointAny.(string)
    if !ok {
        return nil, &utils.ValidationError{Msg: "endpoint must be a string"}
    }

    var ts float64
    switch v := tsAny.(type) {
    case float64:
        ts = v
    case int:
        ts = float64(v)
    default:
        return nil, &utils.ValidationError{Msg: "ts must be a number"}
    }

    return &Request{
        User:     user,
        Endpoint: endpoint,
        Ts:       ts,
    }, nil
}

func (a *RateLimiterApp) ProcessStream(ctx context.Context, r io.Reader) int {
    reader := bufio.NewReader(r)
for {
    line, err := reader.ReadString('\n')
    if err == io.EOF {
        break
    }
    if err != nil {
        break
    }
    line = strings.TrimSpace(line)
    if line == "" {
        continue
    }
    a.handleLine(line)
}
}

func (a *RateLimiterApp) handleLine(line string) {
    req, err := parseRequest(line)
    if err != nil {
        utils.LogSimple("request_error", "error", err.Error())
        return
    }

    allowed, count, reset := a.engine.Allow(req.Endpoint, req.User, req.Ts)

    utils.LogSimple(
        "decision",
        "user", req.User,
        "endpoint", req.Endpoint,
        "ts", req.Ts,
        "allowed", allowed,
        "count", count,
        "reset", reset,
    )
}
