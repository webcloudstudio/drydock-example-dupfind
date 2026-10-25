# FEATURE: Pipeline

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Wire the command-line entry point to walk, size filter, group, order, and render, producing the built `dupfind` executable. |
| Depends On  | ARCHITECTURE.md, FEATURE-Walker.md, FEATURE-SizeFilter.md, FEATURE-Grouping.md, FEATURE-UnreadableFiles.md, FEATURE-ListingCSV.md, FEATURE-ListingOrder.md, FEATURE-ListingEmpty.md, FEATURE-Summary.md, FEATURE-CommandLine.md, FEATURE-TargetValidation.md |
| Provides    | cmd/dupfind executable |
| Consumes    | cli.Run, cli.ValidateTarget, scan.Walk, scan.FilterSize, scan.Group, output.WriteListing, output.WriteSummary |

## Purpose

This is the assembly feature. It owns `cmd/dupfind/main.go`, which passes an execute function to
`cli.Run`. The execute function runs the full pipeline in sequence:

1. walk the target directory (`scan.Walk`)
2. drop zero-byte and too-small files (`scan.FilterSize`)
3. group by size then SHA-256 (`scan.Group`), warning on unreadable files
4. order groups and paths (`output.Order`)
5. render either the listing (`output.WriteListing`) or the summary (`output.WriteSummary`),
   chosen by the `--summary` flag

Only `main` calls `os.Exit`, with the code returned by `cli.Run`.

## Behavior

- The built executable is `dupfind`, produced by `go build ./cmd/dupfind`.
- Listing and summary modes both derive from the same grouping result, so their totals agree.
- Diagnostics go to standard error; standard output carries only the CSV.

## Programmatic Acceptance

=== AC pipeline-listing-end-to-end ===
Intent: The built executable emits the exact listing header and the duplicate rows for a known tree.
Requires: executable=go; scope=test

import os
import tempfile

from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as d:
    with open(os.path.join(d, "a.txt"), "w") as f:
        f.write("same\n")
    with open(os.path.join(d, "b.txt"), "w") as f:
        f.write("same\n")
    with open(os.path.join(d, "c.txt"), "w") as f:
        f.write("different\n")
    code, out, err = cli_runner(d, package="./cmd/dupfind")
    assert code == 0, (code, out, err)
    rows = out.splitlines()
    assert rows[0] == "group,size,sha256,path", out
    paths = sorted(line.split(",")[-1] for line in rows[1:])
    assert paths == ["a.txt", "b.txt"], out
=== END AC pipeline-listing-end-to-end ===

=== AC pipeline-summary-end-to-end ===
Intent: The built executable emits the summary header and one totals row for the same tree.
Requires: executable=go; scope=test

import os
import tempfile

from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as d:
    with open(os.path.join(d, "a.txt"), "w") as f:
        f.write("same\n")
    with open(os.path.join(d, "b.txt"), "w") as f:
        f.write("same\n")
    code, out, err = cli_runner("--summary", d, package="./cmd/dupfind")
    assert code == 0, (code, out, err)
    rows = out.splitlines()
    assert rows[0] == "groups,duplicate_files,reclaimable_bytes", out
    assert rows[1] == "1,1,5", out
=== END AC pipeline-summary-end-to-end ===

=== AC pipeline-usage-exit ===
Intent: The executable exits 2 with no stdout on a usage error.
Requires: executable=go; scope=test

from drydock_go import cli_runner

code, out, err = cli_runner("--min-size", "abc", "tree", package="./cmd/dupfind")
assert code == 2, (code, out, err)
assert out == "", out
=== END AC pipeline-usage-exit ===

## User Acceptance

- None.

## Guardrails

- `main` is the only caller of `os.Exit`.
- The pipeline never writes, moves, or deletes any file under the directory examined.
