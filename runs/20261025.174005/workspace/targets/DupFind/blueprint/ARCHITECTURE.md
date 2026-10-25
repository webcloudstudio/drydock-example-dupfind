# ARCHITECTURE: DupFind

| Field       | Value |
|-------------|-------|
| Version     | 20261025 V1 |
| Description | Module layout, package boundaries, and pipeline of the dupfind command-line duplicate file finder. |
| Depends On  | None |
| Provides    | go.mod, internal/, cmd/dupfind, bin/ |

## Overview

`dupfind` is a single-binary Go command-line tool. It finds regular files under a directory whose
contents are byte-for-byte identical, reports them as a CSV listing or a one-row CSV summary, and
never deletes, moves, or modifies any file.

Pipeline: parse arguments, validate the directory, walk, group by size, hash same-size groups with
SHA-256, order and number groups, emit CSV.

## Technology Stack

Derived from `TECHNOLOGY_STACK.md`; that file is the decision of record.

| Technology | Where it applies |
|---|---|
| Go 1.22 or newer, standard library only | All application code; `go.mod` declares no third-party requirement |
| go test and go vet | Test runner behind `bin/test.sh` |
| POSIX sh | `bin/test.sh` |

## Module Layout

| Path | Package | Responsibility |
|---|---|---|
| `go.mod` | module root | Module path and Go toolchain version; no `require` lines |
| `cmd/dupfind` | `main` | Calls `cli.Run(os.Args[1:], os.Stdout, os.Stderr)` and passes the result to `os.Exit`; the only caller of `os.Exit` |
| `internal/cli` | `cli` | `Run(args []string, stdout, stderr io.Writer) int`: argument parsing, directory validation, pipeline orchestration, exit codes |
| `internal/walk` | `walk` | Recursive directory walk yielding eligible regular files |
| `internal/dupes` | `dupes` | Size grouping, SHA-256 digest grouping, reclaimable bytes, ordering and numbering |
| `internal/report` | `report` | CSV listing and summary writers using `encoding/csv` |
| `bin/test.sh` | POSIX sh | Runs `go vet ./...` then `go test ./...` from the application root |

## Boundaries

- Only `main` calls `os.Exit`. `cli.Run` returns the exit code: 0 success, 1 operational failure, 2 usage error.
- Standard output carries only CSV. A failing command writes nothing to standard output and one diagnostic line to standard error.
- Warnings (unreadable file or subdirectory) go to standard error only and never change the exit code.
- File contents are read only by `dupes`, and only for files that share a size with another file.
- Files under the examined directory are opened read-only. Nothing is written, moved, or deleted there.
- Application code never reads `sources/run.sh` and never embeds a digest, path, or total.

## Module Ownership

| Boundary | Owning module | Allowed low-level APIs |
|---|---|---|
| Directory traversal and file metadata | `internal/walk` | `os.Lstat`, `os.ReadDir`, `io/fs` |
| File content reading and hashing | `internal/dupes` | `os.Open`, `crypto/sha256`, `io.Copy` |
| CSV output | `internal/report` | `encoding/csv` |
| Process arguments and exit codes | `internal/cli` and `cmd/dupfind` | `flag`, `os.Args`, `os.Exit` (main only) |
| Persistence, configuration files, external services | None; the tool has no persistent state | None |

## Foundation Scope

The foundation story creates `go.mod`, the package tree above, a `cmd/dupfind/main.go` that
delegates to `cli.Run`, and a minimal compiling placeholder in each internal package, so that
`go build ./cmd/dupfind`, `go vet ./...`, and `go test ./...` succeed before feature stories land.
The placeholder `Run` returns 0. `bin/` exists; `bin/test.sh` itself is delivered by the test-gate story.

## Programmatic Acceptance

=== AC arch-module-no-third-party ===
Intent: The application root is a Go module that declares no third-party requirement.

import pathlib

text = pathlib.Path("go.mod").read_text()
lines = [ln.strip() for ln in text.splitlines()]
module_lines = [ln for ln in lines if ln.startswith("module ")]
go_lines = [ln for ln in lines if ln.startswith("go ")]
assert len(module_lines) == 1, text
assert len(go_lines) == 1, text
assert "require" not in text, text
=== END AC arch-module-no-third-party ===

=== AC arch-package-tree ===
Intent: The package directories named by the architecture exist in the build tree.

import pathlib

for rel in ["cmd/dupfind", "internal/cli", "internal/walk", "internal/dupes", "internal/report", "bin"]:
    assert pathlib.Path(rel).is_dir(), rel
assert pathlib.Path("cmd/dupfind/main.go").is_file()
=== END AC arch-package-tree ===

=== AC arch-vet-and-build ===
Intent: The foundation compiles and passes go vet across every package.

import subprocess
import sys

for cmd in (["go", "vet", "./..."], ["go", "build", "./..."]):
    result = subprocess.run(cmd, capture_output=True, text=True)
    print(result.stdout)
    print(result.stderr, file=sys.stderr)
    assert result.returncode == 0, cmd
=== END AC arch-vet-and-build ===

## User Acceptance

- None.

## Guardrails

- Standard library only; no third-party modules.
- Never delete, move, or write a file under the directory examined.
- Do not modify, rewrite, regenerate, or substitute `sources/run.sh`, and never read it from application code.
- Do not hardcode a digest, a path, or a total; every value is computed from the directory examined.
- Read a file's contents only when another file has the same size.
