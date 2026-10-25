# FEATURE: Size Grouping

| Field       | Value |
|-------------|-------|
| Version     | 20261025 V1 |
| Description | Group candidate files by exact size and keep only sizes shared by two or more files, so unique-size files are never read. |
| Depends On  | FEATURE-Walk.md, ARCHITECTURE.md |
| Provides    | dupes.GroupBySize |
| Consumes    | walk.Files |

## Purpose

Rule 3 of the build instructions says files are grouped by size first and only files that share a
size are read. This feature implements that first stage in `internal/dupes`.

## Interface

`dupes.GroupBySize` takes the eligible files produced by `walk.Files` (relative path and size) and
returns size buckets. Only buckets with two or more files are returned. The size is taken from the
walk result; file contents are never opened in this stage. A file whose size is unique appears in
no bucket and is therefore never opened for hashing.

## Reads and Writes

- Reads: the in-memory walk result only.
- Writes: nothing.

## Required Project Tests

The project test suite in `internal/dupes` contains `TestGroupBySizeOmitsUniqueSizes` (a unique
size produces no bucket) and `TestGroupBySizeKeepsSharedSizes` (two or more files with one size
produce one bucket holding all of them).

## Programmatic Acceptance

=== AC size-unique-and-shared ===
Intent: Files with a unique size are never reported and a shared size with identical contents is reported.

import csv
import os
import tempfile
from drydock_go import cli_runner

def write(root, rel, data):
    p = os.path.join(root, rel)
    os.makedirs(os.path.dirname(p), exist_ok=True)
    with open(p, "wb") as f:
        f.write(data)

shared = b"dddddd"
unique = b"cc"
with tempfile.TemporaryDirectory() as root:
    write(root, "d1.bin", shared)
    write(root, "sub/d2.bin", shared)
    write(root, "c.bin", unique)
    code, out, err = cli_runner(root)
    print(out)
    print(err)
    assert code == 0
    rows = list(csv.reader(out.splitlines()))[1:]
    paths = sorted(r[3] for r in rows)
    assert paths == ["d1.bin", "sub/d2.bin"]
=== END AC size-unique-and-shared ===

=== AC size-same-size-different-content ===
Intent: Files that share a size but differ in content are not reported as duplicates.

import csv
import os
import tempfile
from drydock_go import cli_runner

def write(root, rel, data):
    p = os.path.join(root, rel)
    os.makedirs(os.path.dirname(p), exist_ok=True)
    with open(p, "wb") as f:
        f.write(data)

with tempfile.TemporaryDirectory() as root:
    write(root, "a.bin", b"aaaa")
    write(root, "b.bin", b"bbbb")
    code, out, err = cli_runner(root)
    print(out)
    print(err)
    assert code == 0
    rows = list(csv.reader(out.splitlines()))
    assert len(rows) == 1
=== END AC size-same-size-different-content ===

=== AC size-unit-tests ===
Intent: The project's own tests for the size bucketing pass and actually execute.

from drydock_go import run_test

result = run_test("./internal/dupes", "TestGroupBySizeOmitsUniqueSizes")
assert result.passed, (result.status, result.stdout, result.stderr)
result = run_test("./internal/dupes", "TestGroupBySizeKeepsSharedSizes")
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC size-unit-tests ===

## User Acceptance

- None.

## Guardrails

- Contents of a file whose size is unique are never opened.
- No file under the examined directory is modified, moved, or deleted.
- Only the standard library is used.
