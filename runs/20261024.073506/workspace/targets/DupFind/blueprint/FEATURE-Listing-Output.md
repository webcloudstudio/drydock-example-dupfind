# FEATURE: Listing Output

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Emits duplicate groups as deterministic CSV. |
| Depends On  | FEATURE-Duplicate-Discovery.md |
| Provides    | CSV duplicate listing |
| Consumes    | duplicate groups |

## Purpose

The default command output lists each file in every duplicate group. The discovery capability supplies groups with file sizes, SHA-256 digests, and paths relative to the selected directory.

## Workflow

1. Receive the discovered duplicate groups after the directory has been validated and scanned.
2. Order groups by reclaimable bytes descending, then by digest ascending.
3. Number groups from 1 in that order.
4. Order paths within each group by byte value ascending.
5. Write the header and one row per file as CSV to standard output, with newline line endings.

The listing header is `group,size,sha256,path`. Each row contains the group number, file size in bytes, lowercase hexadecimal SHA-256 digest, and relative path using `/` separators without a leading `./`. A path containing a comma, double quote, or newline is quoted according to RFC 4180. If there are no duplicate groups, write only the header.

## Programmatic Acceptance

=== AC listing-empty ===
Intent: The default listing succeeds and writes only its header when the selected directory has no duplicate files.

import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)
    unique = root / "unique"
    unique.mkdir()
    (unique / "single.txt").write_bytes(b"only one")
    result = subprocess.run(
        ["./cmd/dupfind", str(unique)], capture_output=True
    )
    print(result.stdout.decode(errors="replace"))
    print(result.stderr.decode(errors="replace"))
    assert result.returncode == 0
=== END AC listing-empty ===

=== AC listing-relative-paths ===
Intent: Listing uses relative slash-separated paths for files supplied in a nested directory.

import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)
    data = b"shared bytes"
    first = root / "scan" / "a" / "copy.bin"
    second = root / "scan" / "b" / "copy.bin"
    first.parent.mkdir(parents=True)
    second.parent.mkdir(parents=True)
    first.write_bytes(data)
    second.write_bytes(data)
    result = subprocess.run(
        ["./cmd/dupfind", str(root / "scan")], capture_output=True
    )
    print(result.stdout.decode(errors="replace"))
    print(result.stderr.decode(errors="replace"))
    assert result.returncode == 0
=== END AC listing-relative-paths ===

=== AC listing-csv-path ===
Intent: Listing succeeds when a duplicate path contains a comma and is encoded through the CSV writer.

import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)
    data = b"same"
    (root / "a,copy").write_bytes(data)
    (root / "b").write_bytes(data)
    result = subprocess.run(
        ["./cmd/dupfind", str(root)], capture_output=True
    )
    print(result.stdout.decode(errors="replace"))
    print(result.stderr.decode(errors="replace"))
    assert result.returncode == 0
=== END AC listing-csv-path ===

## User Acceptance

- None.

## Guardrails

- Never modify files under the directory being scanned.
- Keep deterministic group and path ordering.
