# FEATURE: File Hashing

| Field       | Value |
|-------------|-------|
| Version     | 20261023 V1 |
| Description | Defines candidate filtering and exact duplicate grouping by size and SHA-256 digest. |
| Depends On  | FEATURE-FILE-DISCOVERY.md |
| Provides    | duplicate candidate grouping |
| Consumes    | recursive eligible-file discovery |

## Purpose

DupFind applies the positive minimum-size filter, groups candidates by file size before reading contents, and hashes only files sharing a size. Files with equal size but different contents remain separate. Unreadable files are skipped with one warning and do not change the successful exit status.

## Grouping Rules

For each eligible file at or above the minimum size, files are grouped by size. Contents are read only within a size group containing multiple candidates. SHA-256 digests then identify exact duplicate groups. A group is reportable only when it contains at least two files.

## Programmatic Acceptance

=== AC hashing-min-size ===
Intent: The minimum-size option excludes smaller files while retaining qualifying duplicates.
import csv
import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)
    (root / "small-a").write_bytes(b"x")
    (root / "small-b").write_bytes(b"x")
    (root / "large-a").write_bytes(b"large")
    (root / "large-b").write_bytes(b"large")

    result = subprocess.run(
        ["go", "run", "./cmd/dupfind", "--min-size", "5", str(root)],
        capture_output=True,
        text=True,
    )
    assert result.returncode == 0
    rows = list(csv.DictReader(result.stdout.splitlines()))
    paths = {row["path"] for row in rows}
    assert paths == {"large-a", "large-b"}

=== END AC hashing-min-size ===
=== AC hashing-content-identity ===
Intent: Equal-sized files with different contents are not grouped as duplicates.
import csv
import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)
    (root / "first").write_bytes(b"abcd")
    (root / "second").write_bytes(b"abce")

    result = subprocess.run(
        ["go", "run", "./cmd/dupfind", str(root)],
        capture_output=True,
        text=True,
    )
    assert result.returncode == 0
    rows = list(csv.DictReader(result.stdout.splitlines()))
    assert rows == []

=== END AC hashing-content-identity ===
=== AC hashing-digest ===
Intent: Exact duplicate groups report the lowercase SHA-256 digest computed from file contents.
import csv
import hashlib
import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)
    content = b"same bytes"
    (root / "first").write_bytes(content)
    (root / "second").write_bytes(content)
    expected_digest = hashlib.sha256(content).hexdigest()

    result = subprocess.run(
        ["go", "run", "./cmd/dupfind", str(root)],
        capture_output=True,
        text=True,
    )
    assert result.returncode == 0
    rows = list(csv.DictReader(result.stdout.splitlines()))
    assert len(rows) == 2
    assert {row["sha256"] for row in rows} == {expected_digest}
=== END AC hashing-digest ===

=== AC hashing-unreadable ===
Intent: An unreadable candidate is skipped without turning a scan into an operational failure.
import os
import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)
    readable_a = root / "readable-a"
    readable_b = root / "readable-b"
    unreadable = root / "unreadable"
    readable_a.write_bytes(b"same")
    readable_b.write_bytes(b"same")
    unreadable.write_bytes(b"same")
    unreadable.chmod(0)

    try:
        result = subprocess.run(
            ["go", "run", "./cmd/dupfind", str(root)],
            capture_output=True,
            text=True,
        )
        assert result.returncode == 0
    finally:
        unreadable.chmod(0o600)
=== END AC hashing-unreadable ===

## User Acceptance

- None.

## Guardrails

- Apply the positive minimum-size filter before hashing.
- Read contents only for files sharing a size.
- Use SHA-256 and lowercase hexadecimal output.
- Skip unreadable files with one stderr warning and preserve successful exit status.
