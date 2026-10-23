# FEATURE: Summary Report

| Field       | Value |
|-------------|-------|
| Version     | 20261023 V1 |
| Description | Emit duplicate-group counts and reclaimable-space totals. |
| Depends On  | FEATURE-FILE-HASHING.md |
| Provides    | duplicate summary CSV |
| Consumes    | duplicate candidate grouping |

## Purpose

Summary mode reports the number of duplicate groups, duplicate files beyond one retained copy per group, and reclaimable bytes.

## Behavior

- Emit `groups,duplicate_files,reclaimable_bytes`.
- Emit exactly one data row, including when no duplicates exist.
- Compute totals from the same duplicate groups used by listing mode.
- Apply the positive `--min-size` filter before forming totals.
- Write CSV with LF line endings.

## Programmatic Acceptance

=== AC summary-derived-totals ===
Intent: Summary values equal totals derived from duplicate files created by the criterion.

import csv
import hashlib
import os
import tempfile
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as root:
    groups = [
        (b"alpha", ["a.txt", "b.txt", "c.txt"]),
        (b"longer payload", ["d.txt", "e.txt"]),
    ]
    expected_groups = len(groups)
    expected_duplicate_files = 0
    expected_reclaimable = 0
    for payload, names in groups:
        for name in names:
            with open(os.path.join(root, name), "wb") as handle:
                handle.write(payload)
        expected_duplicate_files += len(names) - 1
        expected_reclaimable += len(payload) * (len(names) - 1)

    code, output, error = cli_runner("--summary", root)
    print(output)
    print(error)
    assert code == 0
    rows = list(csv.reader(output.splitlines()))
    assert rows[0] == ["groups", "duplicate_files", "reclaimable_bytes"]
    assert len(rows) == 2
    actual = [int(value) for value in rows[1]]
    expected = [expected_groups, expected_duplicate_files, expected_reclaimable]
    assert actual == expected
=== END AC summary-derived-totals ===

=== AC summary-min-size ===
Intent: The minimum-size option excludes smaller duplicate groups from summary totals.

import csv
import os
import tempfile
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as root:
    small = b"tiny"
    large = b"large duplicate payload"
    for index in range(2):
        with open(os.path.join(root, f"small-{index}.txt"), "wb") as handle:
            handle.write(small)
        with open(os.path.join(root, f"large-{index}.txt"), "wb") as handle:
            handle.write(large)

    code, output, error = cli_runner("--summary", "--min-size", str(len(large)), root)
    print(output)
    print(error)
    assert code == 0
    rows = list(csv.reader(output.splitlines()))
    assert len(rows) == 2
    actual = [int(value) for value in rows[1]]
    expected = [1, 1, len(large)]
    assert actual == expected
=== END AC summary-min-size ===

=== AC summary-zero-totals ===
Intent: An empty or duplicate-free directory reports zero summary totals.

import csv
import os
import tempfile
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as root:
    with open(os.path.join(root, "empty"), "wb") as handle:
        handle.write(b"")

    code, output, error = cli_runner("--summary", root)
    print(output)
    print(error)
    assert code == 0
    rows = list(csv.reader(output.splitlines()))
    assert rows == [
        ["groups", "duplicate_files", "reclaimable_bytes"],
        ["0", "0", "0"],
    ]
=== END AC summary-zero-totals ===

## User Acceptance

- None.

## Guardrails

- Summary mode must not mutate the scanned directory.
- Summary totals must be derived from discovered duplicate groups, never hardcoded.
