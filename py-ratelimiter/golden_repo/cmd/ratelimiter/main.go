package main

import (
    "context"
    "fmt"
    "os"

    "ratelimiter/internal/config"
    "ratelimiter/internal/engine"
    "ratelimiter/internal/utils"
)

func main() {
    if len(os.Args) < 3 || os.Args[1] != "--config" {
        fmt.Fprintln(os.Stderr, "usage: ratelimiter --config <file>")
        os.Exit(1)
    }
    configPath := os.Args[2]
    raw, err := os.ReadFile(configPath)
    if err != nil {
        utils.LogSimple("config_error", "error", err.Error())
        os.Exit(1)
    }
    appCfg, err := config.ParseConfig(string(raw))
    if err != nil {
        utils.LogSimple("config_error", "error", err.Error())
        os.Exit(1)
    }
    app := engine.NewRateLimiterApp(appCfg)
    code := app.ProcessStream(context.Background(), os.Stdin)
    os.Exit(code)
}
