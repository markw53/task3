package config

import (
    "encoding/json"
    "fmt"

    "ratelimiter/internal/utils"
)

type Mode string

const (
    ModeFixed   Mode = "fixed"
    ModeSliding Mode = "sliding"
)

type LimiterConfig struct {
    Name          string
    Mode          Mode
    WindowSeconds int
    MaxRequests   int
}

type AppConfig struct {
    Limiters map[string]LimiterConfig
}

func ParseConfig(raw string) (*AppConfig, error) {
    var data map[string]any
    if err := json.Unmarshal([]byte(raw), &data); err != nil {
        return nil, &utils.ValidationError{Msg: fmt.Sprintf("Invalid JSON config: %v", err)}
    }
    limitersAny, err := utils.RequireField(data, "limiters")
    if err != nil {
        return nil, err
    }
    limitersSlice, ok := limitersAny.([]any)
    if !ok {
        return nil, &utils.ValidationError{Msg: "limiters must be a list"}
    }
    limiters := make(map[string]LimiterConfig)
    for _, itemAny := range limitersSlice {
        item, ok := itemAny.(map[string]any)
        if !ok {
            return nil, &utils.ValidationError{Msg: "each limiter must be an object"}
        }
        nameAny, err := utils.RequireField(item, "name")
        if err != nil {
            return nil, err
        }
        name, ok := nameAny.(string)
        if !ok {
            return nil, &utils.ValidationError{Msg: "name must be a string"}
        }
        modeAny, err := utils.RequireField(item, "mode")
        if err != nil {
            return nil, err
        }
        modeStr, ok := modeAny.(string)
        if !ok {
            return nil, &utils.ValidationError{Msg: "mode must be a string"}
        }
        var mode Mode
        switch modeStr {
        case string(ModeFixed):
            mode = ModeFixed
        case string(ModeSliding):
            mode = ModeSliding
        default:
            return nil, &utils.ValidationError{Msg: fmt.Sprintf("invalid mode: %s", modeStr)}
        }
        windowSeconds, err := utils.RequirePositiveInt(item, "window_seconds")
        if err != nil {
            return nil, err
        }
        maxRequests, err := utils.RequirePositiveInt(item, "max_requests")
        if err != nil {
            return nil, err
        }
        limiters[name] = LimiterConfig{
            Name:          name,
            Mode:          mode,
            WindowSeconds: windowSeconds,
            MaxRequests:   maxRequests,
        }
    }
    if len(limiters) == 0 {
        return nil, &utils.ValidationError{Msg: "at least one limiter required"}
    }
    return &AppConfig{Limiters: limiters}, nil
}
