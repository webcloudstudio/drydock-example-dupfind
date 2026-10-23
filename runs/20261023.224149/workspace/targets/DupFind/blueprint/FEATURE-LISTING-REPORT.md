# FEATURE: Listing Report

| Field       | Value |
|-------------|-------|
| Version     | 20261023 V1 |
| Description | Emit deterministic CSV rows for duplicate file groups. |
| Depends On  | FEATURE-FILE-HASHING.md |
| Provides    | duplicate listing CSV |
| Consumes    | duplicate candidate grouping |

## Purpose

The listing mode writes duplicate groups as CSV on standard output. It emits no data rows when no duplicate groups exist.

## Behavior

- Use `encoding/csv` with LF line endings.
- Emit `group,size,sha256,path` followed by one row per duplicate file.
- Order groups by reclaimable bytes descending, then SHA-256 digest ascending.
- Number groups from 1.
- Write relative paths with `/` separators and sort them by byte value.
- Preserve RFC 4180 escaping for commas, quotes, and newlines.
- Never write beneath the scanned directory.

## Programmatic Acceptance

=== AC listing-duplicate-rows ===
Intent: Duplicate files are emitted as parsed CSV rows with derived size and digest values.

import csv
import hashlib
import os
import tempfile
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as root:
    payload = b"same listing payload\n"
    paths = [
        os.path.join(root, "b", "copy.txt"),
        os.path.join(root, "a", "original.txt"),
    ]
    for path in paths:
        os.makedirs(os.path.dirname(path), exist_ok=True)
        with open(path, "wb") as handle:
            handle.write(payload)

    code, output, error = cli_runner(root)
    print(output)
    print(error)
    assert code == 0
    rows = list(csv.reader(output.splitlines()))
    header = ["group", "size", "sha256", "path"]
    assert rows[0] == header
    data = rows[1:]
    expected_size = len(payload)
    expected_digest = hashlib.sha256(payload).hexdigest()
    expected_paths = sorted(["a/original.txt", "b/copy.txt"])
    assert len(data) == len(expected_paths)
    assert [row[3] for row in data] == expected_paths
    assert all(row[0] == "1" for row in data)
    assert all(int(row[1]) == expected_size for row in data)
    assert all(row[2] == expected_digest for row in data)
=== END AC listing-duplicate-rows ===

=== AC listing-csv-escaping ===
Intent: Paths containing commas remain one correctly parsed CSV field.

import csv
import os
import tempfile
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as root:
    payload = b"comma path payload"
    names = ["one,copy.txt", "two.txt"]
    for name in names:
        path = os.path.join(root, name)
        with open(path, "wb") as handle:
            handle.write(payload)

    code, output, error = cli_runner(root)
    print(output)
    print(error)
    assert code == 0
    rows = list(csv.reader(output.splitlines()))
    assert len(rows) == 1 + len(names)
    paths = [row[3] for row in rows[1:]]
    assert paths == sorted(names)
    assert "one,copy.txt" in paths
=== END AC listing-csv-escaping ===

=== AC listing-no-duplicates ===
Intent: A directory without duplicate groups produces only the listing header.

import csv
import os
import tempfile
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as root:
    values = [b"first", b"second"]
    for index, value in enumerate(values):
        with open(os.path.join(root, f"unique-{index}.txt"), "wb") as handle:
            handle.write(value)

    code, output, error = cli_runner(root)
    print(output)
    print(error)
    assert code == 0
    rows = list(csv.reader(output.splitlines()))
    assert rows == [["group", "size", "sha256", "path"]]
=== END AC listing-no-duplicates ===

## User Acceptance

- None.

## Guardrails

- The implementation must not delete, move, or modify examined files.
- The supplied conformance harness remains authoritative and is not copied into application code.
