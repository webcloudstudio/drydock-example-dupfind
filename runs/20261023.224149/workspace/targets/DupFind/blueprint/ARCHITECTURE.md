# ARCHITECTURE: DupFind

| Field       | Value |
|-------------|-------|
| Version     | 20261023 V1 |
| Description | Defines the Go module structure and boundaries for the DupFind command-line application. |
| Provides    | Go module, cmd/dupfind entry point |
| Consumes    | — |

## Purpose

DupFind is a standard-library Go command-line application that recursively scans a user-supplied directory and reports duplicate regular files.

## Module Structure

- `cmd/dupfind`: process entry point; parses process arguments, invokes application behavior, and calls `os.Exit`.
- `internal/cli`: command-line parsing, validation, and orchestration.
- `internal/discovery`: recursive traversal and eligible-file selection.
- `internal/grouping`: size filtering, SHA-256 hashing, and duplicate grouping.
- `internal/reporting`: deterministic CSV listing and summary output.
- `bin/test.sh`: POSIX local verification launcher.

## Boundaries

The application reads filesystem metadata and file contents only beneath the requested scan directory. It never writes, deletes, or moves files beneath that directory. Reporting writes only to standard output, while diagnostics are written to standard error.

The process entry point owns process exit status. Application behavior is exposed through testable functions that accept arguments and I/O streams.

## Technology Stack

- Go 1.22 or newer.
- Go standard library only.
- POSIX `sh` for the local test gate.

## Technical Decisions

- Files are grouped by size before contents are read.
- Candidate contents are hashed with SHA-256.
- CSV output uses `encoding/csv`.
- Diagnostics use concise conventional stderr lines; diagnostic wording is not part of the output contract.

## Programmatic Acceptance

=== AC architecture-module ===
Intent: The application is a Go module with the required executable package.
import subprocess

module_check = subprocess.run(
    ["go", "list", "-m"],
    capture_output=True,
    text=True,
)
assert module_check.returncode == 0

build_check = subprocess.run(
    ["go", "build", "./cmd/dupfind"],
    capture_output=True,
    text=True,
)
assert build_check.returncode == 0
=== END AC architecture-module ===

=== AC architecture-standard-library ===
Intent: The executable package builds successfully with the module's declared dependencies.
import subprocess

check = subprocess.run(
    ["go", "build", "-trimpath", "./cmd/dupfind"],
    capture_output=True,
    text=True,
)
assert check.returncode == 0
=== END AC architecture-standard-library ===

## User Acceptance

- None.

## Guardrails

- The implementation must use only the Go standard library.
- The application must never modify files beneath the examined directory.
- `sources/run.sh` remains supplied and unmodified.
