import json
import sys
from typing import Any, Dict


def log(event: str, **fields: Any) -> None:
    payload: Dict[str, Any] = {"event": event}
    payload.update(fields)
    sys.stdout.write(json.dumps(payload) + "\n")
    sys.stdout.flush()
