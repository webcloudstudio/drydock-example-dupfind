# FEATURE: Listing

| Field       | Value |
|-------------|-------|
| Version     | 20261025 V1 |
| Description | Write the CSV duplicate listing with header group,size,sha256,path using encoding/csv and wire the default mode of cli.Run. |
| Depends On  | FEATURE-GroupOrdering.md, FEATURE-ModuleBuild.md, ARCHITECTURE.md |
| Provides    | report.WriteListing |
| Consumes    | dupes.Order |

## Purpose

The default output of `dupfind <directory>` is a CSV listing of every file in every duplicate group.

## Interface

`report.WriteListing` writes to standard output through `encoding/csv` with `\n` line endings
(`UseCRLF` false). The header row is `group,size,sha256,path`. One row follows per file of every
ordered group: group number, size in bytes, lowercase hexadecimal SHA-256, relative path. With no
duplicates only the header is written. Fields containing a comma, a double quote, or a newline are
quoted as RFC 4180 specifies, by the csv writer.

The default (non-summary) mode of `cli.Run` runs the walk, size grouping, digest grouping, and
ordering with the configured `--min-size`, then this writer, and returns 0. `--min-size` excludes
files smaller than the given number of bytes.

## Reads and Writes

- Reads: ordered groups.
- Writes: CSV to standard output only.

## Programmatic Acceptance

=== AC listing-full-output ===
Intent: The entire standard output equals the CSV the standard csv writer produces from the fixture-derived rows.

import csv
import hashlib
import io
import os
import tempfile
from drydock_go import cli_runner

def write(root, rel, data):
    p = os.path.join(root, rel)
    os.makedirs(os.path.dirname(p), exist_ok=True)
    with open(p, "wb") as f:
        f.write(data)

fixtures = {
    b"alpha alpha\n": ["a/one.txt", "b/two.txt", "c/deep/three.txt"],
    b"beta\n": [".hidden/.env", "env.bak"],
}
header = ["group", "size", "sha256", "path"]
groups = []
for content, paths in fixtures.items():
    groups.append((len(content) * (len(paths) - 1), hashlib.sha256(content).hexdigest(), len(content), sorted(paths)))
groups.sort(key=lambda g: (-g[0], g[1]))
buf = io.StringIO()
w = csv.writer(buf, lineterminator="\n")
w.writerow(header)
for n, g in enumerate(groups, start=1):
    for p in g[3]:
        w.writerow([n, g[2], g[1], p])
expected = buf.getvalue()

with tempfile.TemporaryDirectory() as root:
    for content, paths in fixtures.items():
        for p in paths:
            write(root, p, content)
    code, out, err = cli_runner(root)
    print(out)
    print(err)
    assert code == 0
    assert out == expected
    assert "\r" not in out
=== END AC listing-full-output ===

=== AC listing-csv-quoting ===
Intent: Paths containing a comma or a double quote survive a CSV round trip unchanged.

import csv
import os
import tempfile
from drydock_go import cli_runner

def write(root, rel, data):
    p = os.path.join(root, rel)
    os.makedirs(os.path.dirname(p), exist_ok=True)
    with open(p, "wb") as f:
        f.write(data)

content = b"quote me\n"
paths = ["song,live.mp3", 'say "hi".txt', "plain.txt"]
with tempfile.TemporaryDirectory() as root:
    for p in paths:
        write(root, p, content)
    code, out, err = cli_runner(root)
    print(out)
    print(err)
    assert code == 0
    rows = list(csv.reader(out.splitlines()))
    assert all(len(r) == 4 for r in rows)
    assert sorted(r[3] for r in rows[1:]) == sorted(paths)
=== END AC listing-csv-quoting ===

=== AC listing-header-only-when-no-duplicates ===
Intent: A directory without duplicates (or with only zero-byte files) yields the header row only, with exit 0.

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
    write(root, "solo.txt", b"only one of me\n")
    write(root, "empty/one", b"")
    write(root, "empty/two", b"")
    code, out, err = cli_runner(root)
    print(out)
    print(err)
    assert code == 0
    assert out.splitlines() == ["group,size,sha256,path"]
    assert out.endswith("\n")
=== END AC listing-header-only-when-no-duplicates ===

=== AC listing-min-size ===
Intent: Files smaller than --min-size are excluded from the listing while larger duplicates remain.

import csv
import os
import tempfile
from drydock_go import cli_runner

def write(root, rel, data):
    p = os.path.join(root, rel)
    os.makedirs(os.path.dirname(p), exist_ok=True)
    with open(p, "wb") as f:
        f.write(data)

small = b"tiny"
large = b"L" * 50
threshold = len(small) + 1
with tempfile.TemporaryDirectory() as root:
    write(root, "s1", small)
    write(root, "s2", small)
    write(root, "l1", large)
    write(root, "l2", large)
    code, out, err = cli_runner("--min-size", str(threshold), root)
    print(out)
    print(err)
    assert code == 0
    rows = list(csv.reader(out.splitlines()))[1:]
    assert sorted(r[3] for r in rows) == ["l1", "l2"]
    code, out, err = cli_runner(root)
    assert code == 0
    rows = list(csv.reader(out.splitlines()))[1:]
    assert len(rows) == 4
=== END AC listing-min-size ===

=== AC listing-symlinks-and-hidden ===
Intent: Hidden files are listed while symbolic links are neither followed nor reported.

import csv
import os
import tempfile
from drydock_go import cli_runner

def write(root, rel, data):
    p = os.path.join(root, rel)
    os.makedirs(os.path.dirname(p), exist_ok=True)
    with open(p, "wb") as f:
        f.write(data)

content = b"KEY=1\n"
with tempfile.TemporaryDirectory() as root:
    write(root, ".hidden/.env", content)
    write(root, "env.bak", content)
    os.symlink(os.path.join(root, "env.bak"), os.path.join(root, "link.bak"))
    os.symlink(os.path.join(root, ".hidden"), os.path.join(root, "linkdir"))
    code, out, err = cli_runner(root)
    print(out)
    print(err)
    assert code == 0
    rows = list(csv.reader(out.splitlines()))[1:]
    assert sorted(r[3] for r in rows) == [".hidden/.env", "env.bak"]
=== END AC listing-symlinks-and-hidden ===

## User Acceptance

- None.

## Guardrails

- Output is produced with `encoding/csv`; fields are never concatenated by hand.
- A failing command writes nothing to standard output.
- No file under the examined directory is modified, moved, or deleted.
