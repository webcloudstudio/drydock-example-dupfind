# FEATURE: Directory Validation

| Field       | Value |
|-------------|-------|
| Version     | 20261023 V1 |
| Description | Defines validation of the requested scan path before filesystem traversal. |
| Depends On  | FEATURE-CLI-OPTIONS.md |
| Provides    | scan-directory validation |
| Consumes    | dupfind option parsing |

## Purpose

After successful option parsing, DupFind validates that the supplied path exists and is a directory.

## Behavior

A missing path or a path naming a regular file is an operational failure. The command returns exit status 1, writes no standard output, and emits a diagnostic on standard error. A valid directory proceeds to discovery.

## Programmatic Acceptance

=== AC directory-missing ===
Intent: A missing scan path returns operational failure without CSV output.
import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    missing = Path(temp) / "missing"
    result = subprocess.run(
        ["go", "run", "./cmd/dupfind", str(missing)],
        capture_output=True,
        text=True,
    )
    assert result.returncode == 1
    assert result.stdout == ""
    assert result.stderr != ""

=== END AC directory-missing ===
=== AC directory-file ===
Intent: A regular file supplied as the scan path returns operational failure without CSV output.
import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    path = Path(temp) / "input"
    path.write_bytes(b"content")
    result = subprocess.run(
        ["go", "run", "./cmd/dupfind", str(path)],
        capture_output=True,
        text=True,
    )
    assert result.returncode == 1
    assert result.stdout == ""
    assert result.stderr != ""

=== END AC directory-file ===
=== AC directory-success ===
Intent: An existing directory is accepted and produces a successful result.
import csv
import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)
    result = subprocess.run(
        ["go", "run", "./cmd/dupfind", str(root)],
        capture_output=True,
        text=True,
    )
    assert result.returncode == 0
    rows = list(csv.reader(result.stdout.splitlines()))
    assert rows[0][0] == "group"
=== END AC directory-success ===

## User Acceptance

- None.

## Guardrails

- Directory validation occurs before traversal.
- Operational failures write no standard output.
- Validation must not modify the supplied path.
