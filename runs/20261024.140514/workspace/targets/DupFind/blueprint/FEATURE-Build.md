# FEATURE: Build

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Provides the dupfind executable at the documented Go package path. |
| Depends On  | ARCHITECTURE.md |
| Provides    | ./cmd/dupfind |
| Consumes    | Go module and application package layout |

## Purpose

The application root is a Go module. Building from that root with `go build ./cmd/dupfind` produces the dupfind command. The main package lives in `cmd/dupfind` and uses only the standard library.

## Build contract

The executable delegates command behavior to the internal CLI package. Process startup passes command-line arguments and standard streams to `Run`; only the main package exits the process using the returned status code. Build output is produced by the Go toolchain, not by application code.

## Programmatic Acceptance

=== AC build-command ===
Intent: Building the declared executable package succeeds.

import subprocess

result = subprocess.run(
    ["go", "build", "./cmd/dupfind"],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC build-command ===

=== AC build-test-entry ===
Intent: The Go test runner recognizes the application packages.

import subprocess

result = subprocess.run(
    ["go", "test", "./..."],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC build-test-entry ===

## User Acceptance

- None.

## Guardrails

- Use Go 1.22 or newer and the standard library only.
- Do not add command behavior to process startup; keep it in the CLI package.
