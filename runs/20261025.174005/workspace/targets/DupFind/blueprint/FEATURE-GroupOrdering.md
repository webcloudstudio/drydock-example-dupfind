# FEATURE: Group Ordering

| Field       | Value |
|-------------|-------|
| Version     | 20261025 V1 |
| Description | Compute reclaimable bytes per group, order and number the groups, and order paths within each group. |
| Depends On  | FEATURE-DigestGrouping.md, ARCHITECTURE.md |
| Provides    | dupes.Order |
| Consumes    | dupes.GroupByDigest |

## Purpose

Rules 4 and 5 fix how paths are written and how groups and paths are ordered so output is
deterministic.

## Interface

`dupes.Order` takes the duplicate groups and returns them ordered and numbered.

- Reclaimable bytes of a group are `size x (files - 1)`.
- Groups are ordered by reclaimable bytes descending, then by digest ascending, and numbered from 1
  in that order.
- Within a group, paths are ordered by byte value ascending (plain string comparison).
- Paths are relative to the examined directory, use `/` as the separator, and carry no leading `./`.

## Reads and Writes

- Reads: the in-memory groups.
- Writes: nothing.

## Required Project Tests

`internal/dupes` contains `TestOrder` covering reclaimable-first ordering, digest tie-breaking,
and byte-order paths.

## Programmatic Acceptance

=== AC order-by-reclaimable-then-digest ===
Intent: Groups appear by reclaimable bytes descending then digest ascending and are numbered from 1.

import csv
import hashlib
import os
import tempfile
from drydock_go import cli_runner

def write(root, rel, data):
    p = os.path.join(root, rel)
    os.makedirs(os.path.dirname(p), exist_ok=True)
    with open(p, "wb") as f:
        f.write(data)

fixtures = {
    b"0123456789": ["big/a", "big/b"],
    b"xyz": ["s/1", "s/2", "s/3", "s/4", "s/5"],
    b"tie-aa": ["t1/a", "t1/b"],
    b"tie-bb": ["t2/a", "t2/b"],
    b"tie-cc": ["t3/a", "t3/b"],
}
groups = []
for content, paths in fixtures.items():
    groups.append((len(content) * (len(paths) - 1), hashlib.sha256(content).hexdigest(), len(content), sorted(paths)))
expected = sorted(groups, key=lambda g: (-g[0], g[1]))

with tempfile.TemporaryDirectory() as root:
    for content, paths in fixtures.items():
        for p in paths:
            write(root, p, content)
    code, out, err = cli_runner(root)
    print(out)
    print(err)
    assert code == 0
    rows = list(csv.reader(out.splitlines()))[1:]
    seen = []
    for r in rows:
        if not seen or seen[-1][0] != r[0]:
            seen.append((r[0], r[2], int(r[1]), []))
        seen[-1][3].append(r[3])
    assert [int(s[0]) for s in seen] == list(range(1, len(expected) + 1))
    assert [s[1] for s in seen] == [g[1] for g in expected]
    assert [s[3] for s in seen] == [g[3] for g in expected]
=== END AC order-by-reclaimable-then-digest ===

=== AC order-paths-byte-ascending ===
Intent: Paths within a group are in byte order and written relative with slash separators and no leading dot-slash.

import csv
import os
import tempfile
from drydock_go import cli_runner

def write(root, rel, data):
    p = os.path.join(root, rel)
    os.makedirs(os.path.dirname(p), exist_ok=True)
    with open(p, "wb") as f:
        f.write(data)

content = b"order me\n"
paths = ["b.txt", "B.txt", "a b.txt", "a/z.txt", ".hid/x.txt", "a.txt"]
with tempfile.TemporaryDirectory() as root:
    for p in paths:
        write(root, p, content)
    code, out, err = cli_runner(root)
    print(out)
    print(err)
    assert code == 0
    got = [r[3] for r in list(csv.reader(out.splitlines()))[1:]]
    assert got == sorted(paths, key=lambda s: s.encode("utf-8"))
    for g in got:
        assert not g.startswith("./")
        assert not g.startswith("/")
        assert "\\" not in g
=== END AC order-paths-byte-ascending ===

=== AC order-reclaimable-beats-size ===
Intent: A smaller file with more copies precedes a larger file with fewer copies when it reclaims more bytes.

import csv
import hashlib
import os
import tempfile
from drydock_go import cli_runner

def write(root, rel, data):
    p = os.path.join(root, rel)
    os.makedirs(os.path.dirname(p), exist_ok=True)
    with open(p, "wb") as f:
        f.write(data)

large = b"L" * 10
small = b"s" * 3
large_paths = ["l1", "l2"]
small_paths = ["s1", "s2", "s3", "s4", "s5"]
assert len(small) * (len(small_paths) - 1) > len(large) * (len(large_paths) - 1)
with tempfile.TemporaryDirectory() as root:
    for p in large_paths:
        write(root, p, large)
    for p in small_paths:
        write(root, p, small)
    code, out, err = cli_runner(root)
    print(out)
    print(err)
    assert code == 0
    rows = list(csv.reader(out.splitlines()))[1:]
    assert rows[0][2] == hashlib.sha256(small).hexdigest()
    assert rows[0][0] == rows[len(small_paths) - 1][0]
    assert rows[len(small_paths)][2] == hashlib.sha256(large).hexdigest()
=== END AC order-reclaimable-beats-size ===

=== AC order-unit-tests ===
Intent: The project's own ordering test passes and actually executes.

from drydock_go import run_test

result = run_test("./internal/dupes", "TestOrder")
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC order-unit-tests ===

## User Acceptance

- None.

## Guardrails

- Ordering is deterministic and independent of directory enumeration order and locale.
- Paths never contain a leading `./` or a backslash separator.
