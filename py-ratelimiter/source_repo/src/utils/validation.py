from typing import Any, Dict


class ValidationError(Exception):
    pass


def require_field(data: Dict[str, Any], field: str) -> Any:
    if field not in data:
        raise ValidationError(f"Missing required field: {field}")
    return data[field]


def require_positive_int(data: Dict[str, Any], field: str) -> int:
    value = require_field(data, field)
    if not isinstance(value, int) or value <= 0:
        raise ValidationError(f"{field} must be a positive integer")
    return value
