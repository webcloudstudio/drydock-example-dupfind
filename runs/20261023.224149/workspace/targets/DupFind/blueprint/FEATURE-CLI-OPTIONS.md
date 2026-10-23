# FEATURE: CLI Options

| Field       | Value |
|-------------|-------|
| Version     | 20261023 V1 |
| Description | Defines DupFind option parsing and command argument validation. |
| Depends On  | ARCHITECTURE.md |
| Provides    | dupfind option parsing |
| Consumes    | Go module and cmd/dupfind entry point |

## Purpose

The command accepts `--summary`, `--min-size <bytes>`, and exactly one directory argument. The default minimum size is one byte.

## Command Contract

Valid forms are:

```text
dupfind [--summary] [--min-size <bytes>] <directory>
```

Options precede the directory. `--min-size` must be a positive integer. Unknown options, missing arguments, extra arguments, and invalid minimum sizes are usage errors with exit status 2, no standard output, and a diagnostic on standard error. The conventional Go help path is supported.

## Programmatic Acceptance

=== AC cli-valid-options ===
Intent: Valid option combinations are accepted and produce successful command exits.
import csv
import os
import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)
    first = root / "first.txt"
    second = root / "second.txt"
    first.write_bytes(b"same")
    second.write_bytes(b"same")

    result = subprocess.run(
        ["go", "run", "./cmd/dupfind", "--min-size", "1", str(root)],
        capture_output=True,
        text=True,
    )
    assert result.returncode == 0
    rows = list(csv.reader(result.stdout.splitlines()))
    assert rows
    assert rows[0][0] == "group"

=== END AC cli-valid-options ===
=== AC cli-summary-option ===
Intent: The summary option is accepted and selects summary output.
import csv
import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)
    (root / "one").write_bytes(b"same")
    (root / "two").write_bytes(b"same")

    result = subprocess.run(
        ["go", "run", "./cmd/dupfind", "--summary", str(root)],
        capture_output=True,
        text=True,
    )
    assert result.returncode == 0
    rows = list(csv.reader(result.stdout.splitlines()))
    assert rows[0][0] == "groups"

=== END AC cli-summary-option ===
=== AC cli-usage-errors ===
Intent: Malformed arguments return usage status 2 without standard output.
import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)

    cases = [
        ["go", "run", "./cmd/dupfind"],
        ["go", "run", "./cmd/dupfind", str(root), str(root)],
        ["go", "run", "./cmd/dupfind", "--min-size", "abc", str(root)],
        ["go", "run", "./cmd/dupfind", "--min-size", "0", str(root)],
        ["go", "run", "./cmd/dupfind", "--frobnicate", str(root)],
    ]

    for command in cases:
        result = subprocess.run(command, capture_output=True, text=True)
        assert result.returncode == 2
        assert result.stdout == ""
=== END AC cli-usage-errors ===

## User Acceptance

- None.

## Guardrails

- Usage errors must not write CSV to standard output.
- Only positive integer minimum sizes are accepted.
- The implementation must not read or depend on `sources/run.sh`.
