# FEATURE: Duplicate Summary

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Defines the summary of duplicate groups and reclaimable space. |
| Depends On  | FEATURE-Ordering.md |
| Provides    | --summary CSV output |
| Consumes    | ordered duplicate groups and relative paths |

## Purpose

When `--summary` is supplied, write a single CSV row summarizing the duplicate groups that the listing would produce with the same options.

## Output

The CSV header is `groups,duplicate_files,reclaimable_bytes`, followed by one row containing:

- The number of duplicate groups.
- The sum of `files − 1` across all groups.
- The sum of `size × (files − 1)` across all groups.

With no duplicates, emit `0,0,0` as the data row. Use `encoding/csv` and `\n` line endings.

## Programmatic Acceptance

=== AC summary-totals ===
Intent: The summary reports totals derived from the same duplicate groups as the listing.
import csv
import io
import pathlib
import subprocess

root = pathlib.Path("summary-fixture")
root.mkdir(parents=True, exist_ok=True)
payload_a = b"summary group a\n"
payload_b = b"summary group b\n"
(root / "a1").write_bytes(payload_a)
(root / "a2").write_bytes(payload_a)
(root / "a3").write_bytes(payload_a)
(root / "b1").write_bytes(payload_b)
(root / "b2").write_bytes(payload_b)
result = subprocess.run(
    ["./dupfind", "--summary", str(root)],
    capture_output=True,
    text=True,
)
assert result.returncode == 0, (result.stdout, result.stderr)
rows = list(csv.reader(io.StringIO(result.stdout)))
assert rows[0] == ["groups", "duplicate_files", "reclaimable_bytes"]
assert rows[1:] == [
    [
        "2",
        str((3 - 1) + (2 - 1)),
        str(len(payload_a) * (3 - 1) + len(payload_b) * (2 - 1)),
    ]
]
=== END AC summary-totals ===

=== AC summary-empty ===
Intent: A directory with no duplicates reports zero totals.
import csv
import io
import pathlib
import subprocess

root = pathlib.Path("summary-empty-fixture")
root.mkdir(parents=True, exist_ok=True)
(root / "single.txt").write_bytes(b"unique")
result = subprocess.run(
    ["./dupfind", "--summary", str(root)],
    capture_output=True,
    text=True,
)
assert result.returncode == 0, (result.stdout, result.stderr)
rows = list(csv.reader(io.StringIO(result.stdout)))
assert rows[1:] == [["0", "0", "0"]]
=== END AC summary-empty ===

## User Acceptance

- None.

## Guardrails

- Summary totals are computed from the duplicate groups; do not hardcode counts or byte totals.
