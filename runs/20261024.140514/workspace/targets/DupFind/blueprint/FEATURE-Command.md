# FEATURE: Command

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Parses dupfind options and validates its directory argument. |
| Depends On  | ARCHITECTURE.md, FEATURE-Build.md |
| Provides    | dupfind command-line options and usage |
| Consumes    | ./cmd/dupfind |

## Command contract

The command syntax is `dupfind [--summary] [--min-size <bytes>] <directory>`. Options precede exactly one directory argument. `--summary` selects summary CSV. `--min-size` accepts a positive integer and defaults to 1. Help text describes the command and its options.

Unknown options, missing or extra directory arguments, and invalid minimum sizes are usage errors. They return 2, write no standard output, and write one diagnostic line to standard error. A missing path or a path that is not a directory is an operational error: return 1, leave standard output empty, and write one diagnostic line to standard error. Other operational errors use the same conventional operational-failure behavior. File read errors during scanning are handled as warnings and do not change the exit code.

## Programmatic Acceptance

=== AC command-help-usage ===
Intent: Help includes a description of the command and its supported options.

import subprocess

result = subprocess.run(
    ["./dupfind", "--help"],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC command-help-usage ===

=== AC command-invalid-option ===
Intent: Unknown options return the usage-error exit status.

import subprocess

result = subprocess.run(
    ["./dupfind", "--frobnicate", "."],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 2
=== END AC command-invalid-option ===

=== AC command-invalid-size ===
Intent: A zero minimum size returns the usage-error exit status.

import subprocess

result = subprocess.run(
    ["./dupfind", "--min-size", "0", "."],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 2
=== END AC command-invalid-size ===

## User Acceptance

- None.

## Guardrails

- Malformed arguments and operational failures write nothing to standard output.
- Diagnostics go to standard error.
