# FEATURE: Duplicate Discovery

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Finds duplicate regular files by size and SHA-256 digest. |
| Depends On  | ARCHITECTURE.md |
| Provides    | duplicate file groups |

## Purpose

Discover duplicate files beneath the requested directory. Recursively include hidden files and directories. Consider regular files only; do not follow or report symbolic links. Ignore zero-byte files.

## Discovery behavior

Walk the directory tree and retain eligible files that meet the configured minimum size. Group candidates by size before opening them. Read contents only when at least two candidates share a size, then group by SHA-256 digest. Only groups with at least two files sharing both size and digest are duplicates.

A file that cannot be read is skipped, produces one warning line on standard error, and does not change the command's exit status. Discovery must not modify any examined file.

Paths in discovered results are relative to the requested directory. The listing capability formats those paths with `/` separators and sorts them by byte value.

## Programmatic Acceptance

=== AC discovery-groups ===
Intent: Duplicate discovery returns the paths supplied in the fixture for identical nonempty files.

import csv
import io
import tempfile
from pathlib import Path
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as directory:
    root = Path(directory)
    left = root / "left"
    right = root / "nested" / "right"
    right.parent.mkdir()
    payload = b"same contents"
    left.write_bytes(payload)
    right.write_bytes(payload)
    code, out, err = cli_runner(str(root), package="./cmd/dupfind")
    rows = list(csv.reader(io.StringIO(out)))
    expected = {str(left.relative_to(root)), str(right.relative_to(root))}
    assert code == 0, (code, out, err)
    assert {row[3] for row in rows[1:]} == expected
=== END AC discovery-groups ===

=== AC discovery-filters-links-and-empty ===
Intent: Discovery ignores empty files and symbolic links while including hidden regular files.

import csv
import io
import tempfile
from pathlib import Path
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as directory:
    root = Path(directory)
    hidden_dir = root / ".hidden"
    hidden_dir.mkdir()
    hidden = hidden_dir / "copy"
    hidden.write_bytes(b"hidden data")
    ordinary = root / "ordinary"
    ordinary.write_bytes(hidden.read_bytes())
    empty_a = root / "empty-a"
    empty_b = root / "empty-b"
    empty_a.write_bytes(b"")
    empty_b.write_bytes(b"")
    link = root / "linked-copy"
    link.symlink_to(hidden)
    code, out, err = cli_runner(str(root), package="./cmd/dupfind")
    rows = list(csv.reader(io.StringIO(out)))
    paths = {row[3] for row in rows[1:]}
    expected = {str(hidden.relative_to(root)), str(ordinary.relative_to(root))}
    assert code == 0, (code, out, err)
    assert paths == expected
=== END AC discovery-filters-links-and-empty ===

=== AC discovery-preserves-files ===
Intent: Scanning leaves examined file contents unchanged.

import tempfile
from pathlib import Path
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as directory:
    root = Path(directory)
    first = root / "first"
    second = root / "second"
    first.write_bytes(b"preserve me")
    second.write_bytes(first.read_bytes())
    before = {path.name: path.read_bytes() for path in (first, second)}
    code, out, err = cli_runner(str(root), package="./cmd/dupfind")
    after = {path.name: path.read_bytes() for path in (first, second)}
    assert code == 0, (code, out, err)
    assert after == before
=== END AC discovery-preserves-files ===

## User Acceptance

- None.

## Guardrails

- Read file contents only after multiple eligible files share a size.
- Skip unreadable files with one standard-error warning and no exit-status change.
- Never follow or report symbolic links.
- Never modify files beneath the examined directory.
