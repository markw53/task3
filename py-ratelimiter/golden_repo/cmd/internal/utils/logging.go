package utils

import (
    "encoding/json"
    "os"
)

func Log(event string, fields map[string]any) {
    payload := map[string]any{"event": event}
    for k, v := range fields {
        payload[k] = v
    }
    enc := json.NewEncoder(os.Stdout)
    _ = enc.Encode(payload)
}

func LogSimple(event string, kv ...any) {
    fields := map[string]any{}
    for i := 0; i+1 < len(kv); i += 2 {
        key, ok := kv[i].(string)
        if !ok {
            continue
        }
        fields[key] = kv[i+1]
    }
    Log(event, fields)
}
