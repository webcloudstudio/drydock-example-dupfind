# Architecture: DupFind

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Single-binary Go command-line tool that lists byte-identical duplicate files under a directory, built from ./cmd/dupfind on the standard library only. |
| Depends On  | COMPASS.md |
| Provides    | cmd/dupfind executable, internal/scan, internal/output, internal/cli |
| Consumes    | — |

## Purpose

DupFind walks a directory, groups regular files by size and then by SHA-256 digest, and reports
either a CSV listing of every duplicate group or a one-row CSV summary of reclaimable space.
Diagnostics go to standard error. The tool never modifies the files it examines.

## Technology Stack

- Go 1.22 or newer, standard library only, declared in `go.mod` at the application root.
- `go test` with `go vet` behind `bin/test.sh`.
- POSIX `sh` for `bin/test.sh`.

## Modules

| Module | Location | Responsibility |
|---|---|---|
| Entry point | `cmd/dupfind/main.go` | Calls `cli.Run` with the execute function; the only caller of `os.Exit`. |
| CLI | `internal/cli` | Parses `[--summary] [--min-size <bytes>] <directory>`, validates the target, and returns an exit code. |
| Scan | `internal/scan` | Recursive walk, size filter, size-then-SHA-256 grouping, and unreadable-file warnings. |
| Output | `internal/output` | Ordering of groups and paths, listing CSV, and summary CSV via `encoding/csv`. |

## Module Boundaries

- `cli` depends on no other application package except through the injected execute function.
- `scan` performs read-only access to files; it never writes, moves, renames, or deletes.
- `output` consumes grouping results from `scan` and writes only to the supplied writer.
- Only `cmd/dupfind/main.go` calls `os.Exit`.

## Interfaces

| Symbol | Package | Signature |
|---|---|---|
| `Run` | `internal/cli` | `func Run(args []string, stdout, stderr io.Writer) int` |
| `Walk` | `internal/scan` | Yields regular files under a root; hidden entries included; symbolic links neither followed nor reported. |
| `FilterSize` | `internal/scan` | Excludes zero-byte files and files below the configured minimum size. |
| `Group` | `internal/scan` | Groups size-matched files by SHA-256 digest; reads contents only for size-matched files. |
| `WriteListing` | `internal/output` | Writes the `group,size,sha256,path` CSV. |
| `WriteSummary` | `internal/output` | Writes the `groups,duplicate_files,reclaimable_bytes` CSV. |

## Technical Decisions

- Standard library only; `flag.FlagSet` with `flag.ContinueOnError` for argument parsing.
- Hashing uses `crypto/sha256`; files are opened read-only.
- Output uses `encoding/csv` with `\n` line endings.
- Exit codes: `0` success, `1` operational failure (missing or non-directory target), `2` usage error.
- Diagnostics are written to standard error as single lines and are never compared by the grader.

## Module Ownership

| Boundary | Owner | Low-level access |
|---|---|---|
| Filesystem reads | `internal/scan` | `os`, `io`, `crypto/sha256` |
| CSV rendering | `internal/output` | `encoding/csv` |
| Process exit | `cmd/dupfind` | `os.Exit` |

## Programmatic Acceptance

=== AC architecture-go-module ===
Intent: The application root is a Go module declaring Go 1.22 or newer with no third-party requires.

import os
import re

with open("go.mod", encoding="utf-8") as handle:
    text = handle.read()

module_match = re.search(r"^module\s+\S+", text, re.MULTILINE)
go_match = re.search(r"^go\s+(\d+)\.(\d+)", text, re.MULTILINE)
assert module_match is not None, text
assert go_match is not None, text
assert int(go_match.group(1)) == 1 and int(go_match.group(2)) >= 22, text
assert "require" not in text, text
=== END AC architecture-go-module ===

=== AC architecture-cmd-layout ===
Intent: The main package directory and internal packages exist at the declared locations.

import os

assert os.path.isdir("cmd/dupfind"), "cmd/dupfind missing"
assert os.path.isdir("internal"), "internal missing"
assert os.path.isdir("bin"), "bin missing"
=== END AC architecture-cmd-layout ===

## User Acceptance

- None.

## Guardrails

- The tool never deletes, moves, renames, or writes any file under the examined directory.
- Use the Go standard library only; no third-party modules.
- Only `cmd/dupfind/main.go` calls `os.Exit`.
