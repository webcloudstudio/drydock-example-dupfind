# FEATURE: Summary

| Field       | Value |
|-------------|-------|
| Version     | 20261025 V1 |
| Description | Write the one-row CSV summary of groups, duplicate files, and reclaimable bytes when --summary is given. |
| Depends On  | FEATURE-Listing.md, ARCHITECTURE.md |
| Provides    | report.WriteSummary |
| Consumes    | dupes.Order |

## Purpose

With `--summary` the tool reports the totals of the listing that the same options would produce,
instead of the listing itself.

## Interface

`report.WriteSummary` writes, through `encoding/csv` with `\n` line endings, the header
`groups,duplicate_files,reclaimable_bytes` and one row:

- `groups`: the number of duplicate groups.
- `duplicate_files`: the sum of `files - 1` over all groups.
- `reclaimable_bytes`: the sum of every group's `size x (files - 1)`.

With no duplicates the row is `0,0,0`. The `--summary` mode of `cli.Run` uses the same walk,
grouping, and ordering as the listing, with the same `--min-size`, and returns 0.

## Reads and Writes

- Reads: ordered groups.
- Writes: CSV to standard output only.

## Programmatic Acceptance

=== AC summary-totals-match-fixture ===
Intent: The summary row equals totals computed from the fixture's known duplicate structure.

import csv
import os
import tempfile
from drydock_go import cli_runner

def write(root, rel, data):
    p = os.path.join(root, rel)
    os.makedirs(os.path.dirname(p), exist_ok=True)
    with open(p, "wb") as f:
        f.write(data)

fixtures = {
    b"quarterly report\n": ["a/r.txt", "b/r2.txt", "c/r3.txt"],
    b"la la la\n": ["m/s.mp3", "m/s2.mp3"],
    b"unique content here\n": ["u/solo"],
}
dup_groups = [(c, p) for c, p in fixtures.items() if len(p) > 1]
want_groups = len(dup_groups)
want_files = sum(len(p) - 1 for c, p in dup_groups)
want_bytes = sum(len(c) * (len(p) - 1) for c, p in dup_groups)
with tempfile.TemporaryDirectory() as root:
    for content, paths in fixtures.items():
        for p in paths:
            write(root, p, content)
    code, out, err = cli_runner("--summary", root)
    print(out)
    print(err)
    assert code == 0
    rows = list(csv.reader(out.splitlines()))
    assert len(rows) == 2
    assert rows[0] == ["groups", "duplicate_files", "reclaimable_bytes"]
    assert [int(x) for x in rows[1]] == [want_groups, want_files, want_bytes]
=== END AC summary-totals-match-fixture ===

=== AC summary-no-duplicates ===
Intent: With no duplicates the summary row holds three zeros and exit status is 0.

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
    write(root, "empty/one", b"")
    write(root, "empty/two", b"")
    code, out, err = cli_runner("--summary", root)
    print(out)
    print(err)
    assert code == 0
    rows = list(csv.reader(out.splitlines()))
    assert len(rows) == 2
    assert [int(x) for x in rows[1]] == [0, 0, 0]
=== END AC summary-no-duplicates ===

=== AC summary-agrees-with-listing ===
Intent: The summary totals equal the totals derived from the listing produced by the same options, including --min-size.

import csv
import os
import tempfile
from drydock_go import cli_runner

def write(root, rel, data):
    p = os.path.join(root, rel)
    os.makedirs(os.path.dirname(p), exist_ok=True)
    with open(p, "wb") as f:
        f.write(data)

fixtures = {
    b"S" * 4: ["s1", "s2"],
    b"M" * 30: ["m1", "m2", "m3"],
    b"B" * 200: ["b1", "b2"],
}
with tempfile.TemporaryDirectory() as root:
    for content, paths in fixtures.items():
        for p in paths:
            write(root, p, content)
    for opts in ([], ["--min-size", "20"], ["--min-size", "100"], ["--min-size", "1000"]):
        code, listing, err = cli_runner(*opts, root)
        assert code == 0
        rows = list(csv.reader(listing.splitlines()))[1:]
        per_group = {}
        for r in rows:
            per_group.setdefault(r[0], []).append(int(r[1]))
        want = [
            len(per_group),
            sum(len(v) - 1 for v in per_group.values()),
            sum(v[0] * (len(v) - 1) for v in per_group.values()),
        ]
        code, summ, err = cli_runner("--summary", *opts, root)
        print(summ)
        print(err)
        assert code == 0
        srows = list(csv.reader(summ.splitlines()))
        assert len(srows) == 2
        assert [int(x) for x in srows[1]] == want
=== END AC summary-agrees-with-listing ===

=== AC summary-unit-tests ===
Intent: The project's own summary tests pass and actually execute.

from drydock_go import run_test

result = run_test("./internal/report", "TestWriteSummary")
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC summary-unit-tests ===

## User Acceptance

- None.

## Guardrails

- The summary never prints the listing and the listing never prints the summary.
- Output is produced with `encoding/csv`.
- No file under the examined directory is modified, moved, or deleted.
