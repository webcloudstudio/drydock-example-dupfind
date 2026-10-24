# FEATURE: Duplicate Listing

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Defines the duplicate file CSV listing. |
| Depends On  | FEATURE-Ordering.md |
| Provides    | default CSV listing |
| Consumes    | ordered duplicate groups and relative paths |

## Purpose

The default command output is a CSV listing of every file in every duplicate group. It uses the ordered groups and relative paths provided by `FEATURE-Ordering.md`.

## Output

The CSV header is `group,size,sha256,path`. Each subsequent row contains the group number, file size in bytes, lowercase hexadecimal SHA-256 digest, and relative path. Emit one row per file in every duplicate group. When no duplicates exist, emit the header only.

Use Go's `encoding/csv` writer, with `\n` line endings and RFC 4180 field quoting.

## Programmatic Acceptance

=== AC listing-output ===
Intent: The command emits the specified CSV header and one row for each duplicate file, with paths parsed as CSV fields.

import csv
import hashlib
import io
import os
import pathlib
import subprocess

root = pathlib.Path("listing-fixture")
(root / "a").mkdir(parents=True, exist_ok=True)
payload = b"listing input\n"
(root / "a" / "first.txt").write_bytes(payload)
(root / "second,copy.txt").write_bytes(payload)
result = subprocess.run(
    ["./dupfind", str(root)],
    capture_output=True,
    text=True,
)
assert result.returncode == 0, (result.stdout, result.stderr)
rows = list(csv.reader(io.StringIO(result.stdout)))
digest = hashlib.sha256(payload).hexdigest()
assert rows[0] == ["group", "size", "sha256", "path"]
assert rows[1:] == [
    ["1", str(len(payload)), digest, "a/first.txt"],
    ["1", str(len(payload)), digest, "second,copy.txt"],
]
=== END AC listing-output ===

=== AC listing-empty ===
Intent: A directory with no duplicates produces a header-only listing.
import csv
import io
import pathlib
import subprocess

root = pathlib.Path("listing-empty-fixture")
root.mkdir(parents=True, exist_ok=True)
(root / "single.txt").write_text("only one\n")
result = subprocess.run(
    ["./dupfind", str(root)],
    capture_output=True,
    text=True,
)
assert result.returncode == 0, (result.stdout, result.stderr)
rows = list(csv.reader(io.StringIO(result.stdout)))
assert rows[1:] == []
=== END AC listing-empty ===

## User Acceptance

- None.

## Guardrails

- Listing never modifies, moves, or deletes examined files.
