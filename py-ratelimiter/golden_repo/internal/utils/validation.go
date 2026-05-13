package utils

import "fmt"

type ValidationError struct {
    Msg string
}

func (e *ValidationError) Error() string {
    return e.Msg
}

func RequireField(m map[string]any, field string) (any, error) {
    v, ok := m[field]
    if !ok {
        return nil, &ValidationError{Msg: fmt.Sprintf("Missing required field: %s", field)}
    }
    return v, nil
}

func RequirePositiveInt(m map[string]any, field string) (int, error) {
    v, err := RequireField(m, field)
    if err != nil {
        return 0, err
    }
    f, ok := v.(float64)
    if !ok || f <= 0 || float64(int(f)) != f {
        return 0, &ValidationError{Msg: fmt.Sprintf("%s must be a positive integer", field)}
    }
    return int(f), nil
}
