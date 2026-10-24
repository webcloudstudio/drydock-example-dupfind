# ARCHITECTURE: DupFind

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Defines the module structure and boundaries for the DupFind command-line application. |
| Provides    | ./cmd/dupfind |

## Overview

DupFind is a read-only Go command-line application. The `./cmd/dupfind` package starts the process and delegates command behavior to internal application packages. It uses the Go standard library and the module declared in `go.mod`.

## Module Boundaries

| Module | Responsibility | Owns |
|---|---|---|
| `cmd/dupfind` | Process startup; invokes the CLI and exits with its returned status. | `main` |
| `internal/cli` | Parses options and arguments, validates the directory, coordinates discovery and output, and routes diagnostics to stderr. | CLI behavior |
| `internal/discovery` | Walks directories, identifies eligible regular files, hashes same-size candidates, and forms duplicate groups. | Discovery and grouping |
| `internal/output` | Orders duplicate groups and encodes listing or summary CSV. | CSV serialization |

The CLI receives arguments and the stdout/stderr writers from `main`. Discovery returns data to the CLI; output consumes those groups and writes CSV. Application packages do not read the supplied grader or write to examined directories.

## Module Ownership

| Boundary | Owner | Allowed low-level access |
|---|---|---|
| Persistence | None | DupFind has no persistent state or database. |
| Configuration | `internal/cli` | Command-line options only; no external configuration store. |
| File store | `internal/discovery` | Read directory entries and candidate file contents through the Go standard library. Never write, move, or delete examined files. |
| External services | None | No external services. |

## Technology Stack

| Technology | Use |
|---|---|
| Go | Go 1.22 or newer; application and tests use the standard library. |
| go test | Project test runner; `bin/test.sh` also runs `go vet ./...`. |
| Shell | POSIX `sh` for `bin/test.sh` and the supplied grader. |

## Guardrails

- Build the executable as `./cmd/dupfind` from the Go module root.
- Use only the Go standard library; do not add third-party modules.
- Keep the application read-only with respect to every examined directory.
- Do not read or depend on `sources/run.sh` from application code.
- Read file contents only when another eligible file has the same size.

## Programmatic Acceptance

=== AC module-build ===
Intent: The module exposes the declared command package and builds it successfully.

import subprocess

result = subprocess.run(
    ["go", "build", "./cmd/dupfind"],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC module-build ===

=== AC internal-packages ===
Intent: The application provides the declared internal package structure.

import importlib.util

for package in ("cli", "discovery", "output"):
    assert importlib.util.find_spec(f"internal.{package}") is not None
=== END AC internal-packages ===

## User Acceptance

- None.

## Guardrails

- None.
