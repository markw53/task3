import json
import subprocess
import os
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable


@dataclass
class RunResult:
    exit_code: int
    stdout: str
    stderr: str


def run_binary(config_text: str, lines: Iterable[str]) -> RunResult:
    cfg_path = Path("config.json")
    cfg_path.write_text(config_text, encoding="utf-8")

    bin_path = Path(os.environ.get("RATELIMITER_BIN","/target/ratelimiter"))
    proc = subprocess.Popen(
        [str(bin_path), "--config", str(cfg_path)],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )

    stdin_data = "".join(line + "\n" for line in lines)
    out, err = proc.communicate(stdin_data)
    return RunResult(proc.returncode, out, err)


def parse_events(stdout: str):
    for line in stdout.splitlines():
        line = line.strip()
        if not line:
            continue
        try:
            yield json.loads(line)
        except json.JSONDecodeError:
            continue
