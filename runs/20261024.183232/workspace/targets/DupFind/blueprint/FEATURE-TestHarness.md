# FEATURE: Test Harness

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Ship a POSIX `bin/test.sh` that runs `go vet` and `go test`, and serve as the terminal verification gate that runs the supplied grader. |
| Depends On  | ARCHITECTURE.md, FEATURE-Walker.md, FEATURE-SizeFilter.md, FEATURE-Grouping.md, FEATURE-UnreadableFiles.md, FEATURE-ListingCSV.md, FEATURE-ListingOrder.md, FEATURE-ListingEmpty.md, FEATURE-Summary.md, FEATURE-CommandLine.md, FEATURE-TargetValidation.md, FEATURE-Pipeline.md |
| Provides    | bin/test.sh |
| Consumes    | cli.Run, cmd/dupfind executable |

## Purpose

Provide `bin/test.sh`, a POSIX-compatible script that runs `go vet ./...` and `go test ./...` from
the application root and exits zero only when both pass. This story is the terminal verification
gate: it depends on every implementation story, and its acceptance runs the supplied grader
`sh sources/run.sh` as the full suite, requiring exit zero.

## Behavior

- `bin/test.sh` is POSIX `sh`, not bash.
- It runs `go vet ./...` then `go test ./...` from the application root.
- It exits zero only when both succeed; any failure exits non-zero.
- The grader `sources/run.sh` is staged verbatim and is never modified, rewritten, or substituted.

## Programmatic Acceptance

=== AC harness-test-script-passes ===
Intent: bin/test.sh exits zero when go vet and go test both pass for the completed application.
Requires: executable=sh; scope=test

import subprocess

result = subprocess.run(["sh", "bin/test.sh"], capture_output=True, text=True)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0, result.returncode
=== END AC harness-test-script-passes ===

=== AC harness-grader-full ===
Intent: The supplied grader sources/run.sh exits zero against the completed application.
Suite: full
Requires: executable=sh; scope=test

import subprocess

result = subprocess.run(["sh", "sources/run.sh"], capture_output=True, text=True)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0, result.returncode
=== END AC harness-grader-full ===

## User Acceptance

- None.

## Guardrails

- `sources/run.sh` is never modified, regenerated, or substituted, and application code never reads it.
- `bin/test.sh` never runs the grader itself; the grader is verified only by the terminal acceptance.
- Application code never hardcodes a digest, path, or total derived from the grader's fixed tree.
