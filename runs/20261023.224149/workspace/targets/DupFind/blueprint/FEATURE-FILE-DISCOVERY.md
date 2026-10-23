# FEATURE: File Discovery

| Field       | Value |
|-------------|-------|
| Version     | 20261023 V1 |
| Description | Defines recursive selection of eligible regular files beneath the scan directory. |
| Depends On  | FEATURE-DIRECTORY-VALIDATION.md |
| Provides    | recursive eligible-file discovery |
| Consumes    | scan-directory validation |

## Purpose

DupFind recursively walks the validated directory, including hidden entries and directories. It considers regular files only and ignores symbolic links and zero-byte files.

## Discovery Rules

- Hidden files and directories are included.
- Regular files are eligible.
- Symbolic links are neither followed nor reported.
- Empty files are ignored.
- The examined tree is read-only from DupFind's perspective.

## Programmatic Acceptance

=== AC discovery-hidden-and-recursive ===
Intent: Recursive discovery includes hidden files and nested regular files.
import csv
import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)
    hidden = root / ".hidden"
    nested = root / "nested" / "deep"
    hidden.mkdir()
    nested.mkdir(parents=True)
    (hidden / "secret").write_bytes(b"duplicate")
    (nested / "copy").write_bytes(b"duplicate")

    result = subprocess.run(
        ["go", "run", "./cmd/dupfind", str(root)],
        capture_output=True,
        text=True,
    )
    assert result.returncode == 0
    rows = list(csv.DictReader(result.stdout.splitlines()))
    paths = {row["path"] for row in rows}
    assert ".hidden/secret" in paths
    assert "nested/deep/copy" in paths

=== END AC discovery-hidden-and-recursive ===
=== AC discovery-symbolic-link-and-empty ===
Intent: Symbolic links and empty files are excluded from duplicate listings.
import csv
import os
import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)
    source = root / "source"
    source.write_bytes(b"duplicate")
    os.symlink("source", root / "link")
    (root / "empty-a").write_bytes(b"")
    (root / "empty-b").write_bytes(b"")

    result = subprocess.run(
        ["go", "run", "./cmd/dupfind", str(root)],
        capture_output=True,
        text=True,
    )
    assert result.returncode == 0
    rows = list(csv.DictReader(result.stdout.splitlines()))
    paths = {row["path"] for row in rows}
    assert "link" not in paths
    assert "empty-a" not in paths
    assert "empty-b" not in paths

=== END AC discovery-symbolic-link-and-empty ===
=== AC discovery-read-only ===
Intent: Scanning does not modify files in the examined directory.
import hashlib
import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)
    first = root / "first"
    second = root / "second"
    first.write_bytes(b"same")
    second.write_bytes(b"same")
    before = {
        path.relative_to(root): hashlib.sha256(path.read_bytes()).digest()
        for path in root.rglob("*")
        if path.is_file() and not path.is_symlink()
    }

    result = subprocess.run(
        ["go", "run", "./cmd/dupfind", str(root)],
        capture_output=True,
        text=True,
    )
    assert result.returncode == 0
    after = {
        path.relative_to(root): hashlib.sha256(path.read_bytes()).digest()
        for path in root.rglob("*")
        if path.is_file() and not path.is_symlink()
    }
    assert result.returncode == 0
    assert after == before
=== END AC discovery-read-only ===

## User Acceptance

- None.

## Guardrails

- Never follow or report symbolic links.
- Never report empty files.
- Never write, delete, move, or modify anything beneath the scan directory.
