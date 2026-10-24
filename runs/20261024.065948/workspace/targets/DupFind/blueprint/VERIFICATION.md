# FEATURE: Supplied Grader Verification

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Runs the supplied grader as the terminal full conformance gate. |
| Depends On  | FEATURE-Duplicate-Finding.md |
| Provides    | Full supplied-grader verification |

## Purpose

The supplied grader is the final acceptance gate for the completed command. It builds `./cmd/dupfind`, runs all 15 checks against its fixed directory tree, and decides success by its exit status. Keep the grader unchanged.

## Programmatic Acceptance

=== AC full-grader ===
Intent: The supplied grader passes all 15 checks against the completed command.
Suite: full

import subprocess

result = subprocess.run(
    ["sh", "sources/run.sh"],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC full-grader ===

## User Acceptance

- None.

## Guardrails

- Do not alter, replace, filter, or reinterpret `sources/run.sh`.
- The grader's exit status is the conformance verdict.
