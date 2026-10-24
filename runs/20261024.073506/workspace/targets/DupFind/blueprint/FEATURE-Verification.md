# FEATURE: Verification

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Provides the project verification script and runs the supplied final grader. |
| Depends On  | FEATURE-Command-Interface.md, FEATURE-Duplicate-Discovery.md, FEATURE-Listing-Output.md, FEATURE-Summary-Output.md |
| Provides    | bin/test.sh, ./cmd/dupfind |
| Consumes    | ./cmd/dupfind, supplied grader |

## Purpose

Provide a POSIX-compatible project test command that runs the vet and automated test suites from the application root, and complete terminal verification with the supplied grader after all implementation stories are complete.

## Verification

`bin/test.sh` runs `go vet ./...` and `go test ./...`. It returns success only if both commands succeed. The supplied `sources/run.sh` is staged verbatim and is the final grader; it must not be edited, replaced, or read by application code.

## Programmatic Acceptance

=== AC project-test-command ===
Intent: The project test command exits successfully after running the Go vet and test commands.

import subprocess

result = subprocess.run(["sh", "bin/test.sh"], capture_output=True)
print(result.stdout.decode(errors="replace"))
print(result.stderr.decode(errors="replace"))
assert result.returncode == 0
=== END AC project-test-command ===

=== AC supplied-grader ===
Intent: The terminal story passes all checks in the supplied duplicate finder grader.
Suite: full

import subprocess

result = subprocess.run(["sh", "sources/run.sh"], capture_output=True)
print(result.stdout.decode(errors="replace"))
print(result.stderr.decode(errors="replace"))
assert result.returncode == 0
=== END AC supplied-grader ===

## User Acceptance

- None.

## Guardrails

- Do not modify, rewrite, regenerate, or substitute `sources/run.sh`.
- Do not read `sources/run.sh` from application code or hardcode answers from it.
- Run the full supplied grader only as terminal verification.
