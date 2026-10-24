# FEATURE: Summary Output

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Reports duplicate groups and recoverable bytes as one CSV summary row. |
| Depends On  | FEATURE-Duplicate-Discovery.md |
| Provides    | CSV duplicate summary |
| Consumes    | duplicate groups |

## Purpose

The `--summary` option replaces the default listing with one summary row. It uses the same discovered groups and selected minimum-size threshold as the listing.

## Workflow

1. Receive the duplicate groups discovered for the selected directory.
2. Compute the number of groups.
3. Sum duplicate files as `files − 1` per group.
4. Sum reclaimable bytes as `size × (files − 1)` per group.
5. Write the CSV header and one data row to standard output.

The header is `groups,duplicate_files,reclaimable_bytes`. With no duplicate groups, the row contains zero for all three totals.

## Programmatic Acceptance

=== AC summary-no-duplicates ===
Intent: Summary succeeds with zero totals when no non-empty duplicate groups exist.

import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)
    empty = root / "empty"
    empty.mkdir()
    (empty / "one").write_bytes(b"")
    (empty / "two").write_bytes(b"")
    result = subprocess.run(
        ["./cmd/dupfind", "--summary", str(empty)], capture_output=True
    )
    print(result.stdout.decode(errors="replace"))
    print(result.stderr.decode(errors="replace"))
    assert result.returncode == 0
=== END AC summary-no-duplicates ===

=== AC summary-filtered ===
Intent: Summary succeeds when the minimum-size option filters smaller duplicate files.

import subprocess
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as temp:
    root = Path(temp)
    small = b"x"
    large = b"large payload"
    (root / "small-a").write_bytes(small)
    (root / "small-b").write_bytes(small)
    (root / "large-a").write_bytes(large)
    (root / "large-b").write_bytes(large)
    result = subprocess.run(
        ["./cmd/dupfind", "--summary", "--min-size", str(len(large)), str(root)],
        capture_output=True,
    )
    print(result.stdout.decode(errors="replace"))
    print(result.stderr.decode(errors="replace"))
    assert result.returncode == 0
=== END AC summary-filtered ===

## User Acceptance

- None.

## Guardrails

- Summary mode is read-only and does not delete or modify scanned files.
- Derive totals from the selected duplicate groups.
