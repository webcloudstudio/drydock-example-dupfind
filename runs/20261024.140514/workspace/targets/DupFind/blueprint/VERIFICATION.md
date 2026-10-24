# FEATURE: Supplied Grader Verification

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Defines the terminal verification against the supplied grader. |
| Depends On  | ARCHITECTURE.md, FEATURE-Build.md, FEATURE-Command.md, FEATURE-Scanning.md, FEATURE-Ordering.md, FEATURE-Listing.md, FEATURE-Summary.md, FEATURE-Tests.md |
| Provides    | full supplied grader verification |
| Consumes    | sources/run.sh |

## Purpose

The supplied `sources/run.sh` is the authoritative final conformance grader. Run it from the completed application root; its exit status is the acceptance verdict. The staged grader compares command exit codes and standard output against its fixed directory tree. Standard error is diagnostic only.

Do not modify, rewrite, regenerate, substitute, or read the grader from application code.

## Programmatic Acceptance

=== AC full-grader ===
Intent: The complete supplied grader passes all 15 checks.
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

- The grader is staged context and must remain verbatim.
- Verification runs only after all implementation stories are complete.
