# ARCHITECTURE: DupFind

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Defines the Go module, executable boundary, and test runner for DupFind. |
| Provides    | Go module, ./cmd/dupfind, bin/test.sh |

## Purpose

DupFind is a read-only command-line tool that recursively finds duplicate regular files beneath one directory. It groups candidates by size before reading contents, then compares SHA-256 digests. The executable reports duplicate groups or a summary as CSV.

## Technology Stack

- Go 1.22 or newer, standard library only: application module, executable, behavior, and tests.
- Go test: automated tests.
- POSIX sh: `bin/test.sh` and the supplied grader.

## Module boundaries

| Module | Owns | Allowed low-level access |
|---|---|---|
| `cmd/dupfind` | Process startup and `os.Exit` | Calls `internal/cli`; no product logic |
| `internal/cli` | Argument parsing, directory validation, command orchestration, diagnostics, and output selection | Calls discovery and output packages |
| `internal/discovery` | Walking, candidate filtering, reading same-size files, digesting, and duplicate grouping | Filesystem reads only |
| `internal/output` | CSV listing and summary formatting | `encoding/csv`; receives grouped results |
| `bin/test.sh` | Project verification entry point | Runs `go vet ./...` then `go test ./...` |

The application does not use a database, configuration store, file store, or external service. It never writes to the examined directory.

## Process contract

The CLI behavior is exposed as `Run(args []string, stdout, stderr io.Writer) int`. `cmd/dupfind` passes its arguments and standard streams to this function and is the only package that calls `os.Exit`.

Options precede the directory argument. Exactly one directory is required. Operational failures return 1; usage failures return 2; successful runs return 0. Failures leave standard output empty and write one diagnostic line to standard error.

## Test and build contract

The root is a Go module declaring Go 1.22 or newer. Build the executable with `go build ./cmd/dupfind`. `bin/test.sh` runs `go vet ./...` and `go test ./...` from the module root and succeeds only if both commands succeed.

## Programmatic Acceptance

=== AC module-build ===
Intent: The module builds the declared executable package.

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

=== AC test-runner ===
Intent: The project test runner succeeds when vet and all Go tests pass.

import subprocess

result = subprocess.run(
    ["sh", "bin/test.sh"],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC test-runner ===

## User Acceptance

- None.

## Guardrails

- Use only Go's standard library.
- Keep product behavior outside `cmd/dupfind`; only main handles process startup and `os.Exit`.
- Never modify, move, delete, or write files under the examined directory.
- Do not modify, rewrite, regenerate, substitute, or read `sources/run.sh` from application code.
