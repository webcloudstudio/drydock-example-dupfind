# FEATURE: Digest Grouping

| Field       | Value |
|-------------|-------|
| Version     | 20261025 V1 |
| Description | Hash same-size files with SHA-256 and form duplicate groups, skipping unreadable files with a stderr warning. |
| Depends On  | FEATURE-SizeGrouping.md, ARCHITECTURE.md |
| Provides    | dupes.GroupByDigest |
| Consumes    | dupes.GroupBySize |

## Purpose

Rule 3 groups same-size files by the SHA-256 digest of their contents; a group of two or more files
with the same size and digest is a duplicate group. Rule 6 skips unreadable files with one warning.

## Interface

`dupes.GroupByDigest` takes the size buckets from `dupes.GroupBySize`. For each bucket it opens each
file read-only and streams it through `crypto/sha256`. The digest is the lowercase hexadecimal
encoding of the full contents. Files with identical size and digest form one group; groups with a
single member are dropped.

A file that cannot be opened or read yields exactly one warning line on standard error, is excluded
from grouping, and does not change the exit code. Files are opened read-only; nothing is written,
moved, or deleted under the examined directory.

## Reads and Writes

- Reads: the contents of files in shared-size buckets.
- Writes: warning lines to standard error only.

## Required Project Tests

`internal/dupes` contains `TestGroupByDigest` (identical contents group, same-size different
contents do not) and `TestGroupByDigestSkipsUnreadable` (an unreadable file is skipped with one
warning and does not fail the call).

## Programmatic Acceptance

=== AC digest-columns-match-contents ===
Intent: Reported size and sha256 are the byte length and lowercase SHA-256 of the fixture contents.

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

content = b"quarterly report\n"
with tempfile.TemporaryDirectory() as root:
    write(root, "a/one.txt", content)
    write(root, "b/two.txt", content)
    write(root, "c/three.txt", content)
    code, out, err = cli_runner(root)
    print(out)
    print(err)
    assert code == 0
    rows = list(csv.reader(out.splitlines()))[1:]
    assert len(rows) == 3
    for r in rows:
        assert int(r[1]) == len(content)
        assert r[2] == hashlib.sha256(content).hexdigest()
    assert sorted(r[3] for r in rows) == ["a/one.txt", "b/two.txt", "c/three.txt"]
=== END AC digest-columns-match-contents ===

=== AC digest-splits-same-size-groups ===
Intent: Same-size files with different digests form separate groups and a lone digest is dropped.

import csv
import os
import tempfile
from drydock_go import cli_runner

def write(root, rel, data):
    p = os.path.join(root, rel)
    os.makedirs(os.path.dirname(p), exist_ok=True)
    with open(p, "wb") as f:
        f.write(data)

first = b"abcd"
second = b"abce"
third = b"abcf"
with tempfile.TemporaryDirectory() as root:
    write(root, "x1", first)
    write(root, "x2", first)
    write(root, "y1", second)
    write(root, "y2", second)
    write(root, "z1", third)
    code, out, err = cli_runner(root)
    print(out)
    print(err)
    assert code == 0
    rows = list(csv.reader(out.splitlines()))[1:]
    by_group = {}
    for r in rows:
        by_group.setdefault(r[0], []).append(r[3])
    assert len(by_group) == 2
    assert sorted(sorted(v) for v in by_group.values()) == [["x1", "x2"], ["y1", "y2"]]
=== END AC digest-splits-same-size-groups ===

=== AC digest-unreadable-skipped ===
Intent: An unreadable duplicate candidate is skipped without changing the exit code (skipped when permissions cannot be enforced).

import csv
import os
import tempfile
from drydock_go import cli_runner

def write(root, rel, data):
    p = os.path.join(root, rel)
    os.makedirs(os.path.dirname(p), exist_ok=True)
    with open(p, "wb") as f:
        f.write(data)

content = b"same bytes\n"
with tempfile.TemporaryDirectory() as root:
    write(root, "ok1.txt", content)
    write(root, "ok2.txt", content)
    write(root, "locked.txt", content)
    os.chmod(os.path.join(root, "locked.txt"), 0)
    try:
        enforced = False
        try:
            open(os.path.join(root, "locked.txt"), "rb").close()
        except PermissionError:
            enforced = True
        code, out, err = cli_runner(root)
        print(out)
        print(err)
        assert code == 0
        rows = list(csv.reader(out.splitlines()))[1:]
        paths = sorted(r[3] for r in rows)
        if enforced:
            assert paths == ["ok1.txt", "ok2.txt"]
        else:
            assert paths == ["locked.txt", "ok1.txt", "ok2.txt"]
    finally:
        os.chmod(os.path.join(root, "locked.txt"), 0o600)
=== END AC digest-unreadable-skipped ===

=== AC digest-tree-unmodified ===
Intent: Running the tool leaves every file's bytes and the directory listing unchanged.

import os
import tempfile
from drydock_go import cli_runner

def write(root, rel, data):
    p = os.path.join(root, rel)
    os.makedirs(os.path.dirname(p), exist_ok=True)
    with open(p, "wb") as f:
        f.write(data)

def snapshot(root):
    snap = {}
    for dirpath, dirnames, filenames in os.walk(root):
        for n in filenames:
            p = os.path.join(dirpath, n)
            with open(p, "rb") as f:
                snap[os.path.relpath(p, root)] = f.read()
        for n in dirnames:
            snap[os.path.relpath(os.path.join(dirpath, n), root) + "/"] = None
    return snap

content = b"keep me\n"
with tempfile.TemporaryDirectory() as root:
    write(root, "a/f1", content)
    write(root, "b/f2", content)
    write(root, "c/other", b"different!")
    before = snapshot(root)
    code, out, err = cli_runner(root)
    print(out)
    print(err)
    assert code == 0
    code, out, err = cli_runner("--summary", root)
    assert code == 0
    assert snapshot(root) == before
=== END AC digest-tree-unmodified ===

=== AC digest-unit-tests ===
Intent: The project's own digest grouping tests pass and actually execute.

from drydock_go import run_test

result = run_test("./internal/dupes", "TestGroupByDigest")
assert result.passed, (result.status, result.stdout, result.stderr)
result = run_test("./internal/dupes", "TestGroupByDigestSkipsUnreadable")
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC digest-unit-tests ===

## User Acceptance

- None.

## Guardrails

- Files are opened read-only; nothing under the examined directory is written, moved, or deleted.
- Warnings go to standard error only and never alter the exit code.
- Only the standard library is used; no digest or path is hardcoded.
