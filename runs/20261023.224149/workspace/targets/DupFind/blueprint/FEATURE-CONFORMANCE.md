# FEATURE: Complete Conformance

| Field       | Value |
|-------------|-------|
| Version     | 20261023 V1 |
| Description | Verify the completed DupFind implementation against the supplied conformance suite. |
| Depends On  | FEATURE-CLI-OPTIONS.md, FEATURE-DIRECTORY-VALIDATION.md, FEATURE-FILE-DISCOVERY.md, FEATURE-FILE-HASHING.md, FEATURE-LISTING-REPORT.md, FEATURE-SUMMARY-REPORT.md, FEATURE-TEST-GATE.md |
| Provides    | complete supplied conformance verification |
| Consumes    | bin/test.sh, dupfind executable |

## Purpose

The supplied `sources/run.sh` is the terminal acceptance authority for the completed command-line application.

## Verification

- Run the staged harness unfiltered from the application root.
- Preserve `sources/run.sh` byte-for-byte.
- Require the harness process to exit zero.
- Do not interpret or replace individual harness checks.

## Programmatic Acceptance

=== AC complete-conformance ===
Intent: The supplied terminal harness passes for the completed DupFind application.
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
=== END AC complete-conformance ===

## User Acceptance

- None.

## Guardrails

- Do not modify, rewrite, regenerate, filter, or substitute `sources/run.sh`.
- Do not execute the full conformance suite before this terminal story.
