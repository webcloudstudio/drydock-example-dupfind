# ARCHITECTURE: DupFind

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Defines the Go command structure and its module boundaries. |
| Provides    | ./cmd/dupfind, bin/test.sh |

## Overview

DupFind is a Go command-line application with one executable at `./cmd/dupfind`. The root Go module uses Go 1.22 or newer and the standard library only.

## Modules

| Module | Responsibility |
|---|---|
| `cmd/dupfind` | Process entry point; calls the CLI runner and exits with its returned status. |
| `internal/cli` | Parses arguments, validates the requested directory, coordinates discovery and output, and maps failures to exit codes. |
| `internal/duplicate` | Walks the requested directory, filters regular files, groups candidates by size and SHA-256 digest, and orders duplicate groups. |
| `internal/output` | Writes listing and summary CSV through `encoding/csv`. |
| `bin/test.sh` | Runs `go vet ./...` followed by `go test ./...`, stopping on failure. |

The command reads examined files only. It does not modify, delete, or move files under the requested directory.

## Module ownership

| Boundary | Owner | Allowed access |
|---|---|---|
| Persistence | None | The application has no persistent state. |
| Configuration | `internal/cli` | Command-line options only; no configuration store. |
| File discovery and reads | `internal/duplicate` | Read-only directory traversal and file reads for hashing. |
| External services | None | The application uses no external service. |

## Technology Stack

| Technology | Application |
|---|---|
| Go | Go 1.22 or newer; implementation and executable. |
| Go standard library | File traversal, SHA-256, CSV, flags, and tests. |
| POSIX shell | `bin/test.sh` quality-check launcher. |

## Programmatic Acceptance

=== AC foundation-module ===
Intent: The command package builds from the application root with the declared Go toolchain.
import subprocess

result = subprocess.run(
    ["go", "build", "./cmd/dupfind"],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC foundation-module ===

=== AC quality-runner ===
Intent: The project quality script runs both Go quality checks successfully.
import subprocess

result = subprocess.run(
    ["sh", "bin/test.sh"],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC quality-runner ===

## User Acceptance

- None.

## Guardrails

- The source distribution grader is supplied and must remain unchanged.
- The application must not read the grader or embed expected answers.
- The tool must never modify files in the directory it examines.
