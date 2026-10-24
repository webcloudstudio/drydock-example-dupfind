# ARCHITECTURE: DupFind

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Defines the Go module and package boundaries for the dupfind command-line application. |
| Depends On  |  |
| Provides    | Go module and application package layout |
| Consumes    |  |

## Purpose

DupFind is a Go command-line application that recursively scans a requested directory and reports regular files with byte-for-byte identical contents. The executable is `./cmd/dupfind`; implementation uses Go 1.22 or newer and the standard library only.

## Modules and ownership

| Module | Owns | Allowed low-level access |
|---|---|---|
| `cmd/dupfind` | Process startup and exit status | `os.Args`, `os.Exit`, standard streams |
| `internal/cli` | Argument parsing, validation, command flow, diagnostics, and output selection | Standard library; delegates scanning and formatting |
| `internal/scan` | Directory traversal, candidate selection, content hashing, and duplicate groups | Read-only filesystem access and SHA-256 |
| `internal/output` | Deterministic ordering and CSV listing or summary formatting | `encoding/csv` and standard library |

Only `cmd/dupfind` calls `os.Exit`. CLI behavior is exposed through `Run(args []string, stdout, stderr io.Writer) int` so it can be exercised without terminating the test process. The CLI uses a local `flag.FlagSet` with `flag.ContinueOnError`.

## Data flow

The CLI validates options and the directory argument, confirms the path names a directory, and calls the scanner with the minimum size. The scanner walks the tree without following symbolic links, selects regular files at or above the minimum size, and reads file contents only when at least two selected files share a size. It groups candidates by SHA-256 digest and returns duplicate groups. The output package orders groups and paths, then writes listing or summary CSV.

Scanner and output interfaces use in-memory Go types; no persistent state or external service is used. Diagnostics are written to standard error. A failing command leaves standard output empty.

## Technology Stack

- Go 1.22 or newer, standard library only: module, executable, CLI, scanning, formatting, and tests.
- POSIX shell: project test entry point and supplied grader.

## Guardrails

- Never read, write, delete, or move files under the directory being examined; examined files are read only.
- Do not read or depend on `sources/run.sh` from application code.
- Do not hardcode file paths, digests, or totals.
- Do not add third-party modules.
- Read file contents only when another eligible regular file has the same size.

## Programmatic Acceptance

=== AC architecture-cli-contract ===
Intent: The architecture exposes the documented CLI runner contract and the executable package builds.

import subprocess

result = subprocess.run(
    ["go", "build", "./cmd/dupfind"],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC architecture-cli-contract ===

=== AC architecture-package-layout ===
Intent: The internal application packages and command entry point compile together.

import subprocess

result = subprocess.run(
    ["go", "test", "./..."],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC architecture-package-layout ===

## User Acceptance

- None.

## Guardrails

- None.
