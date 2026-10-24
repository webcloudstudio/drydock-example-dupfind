# FEATURE: Go Test Harness

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Defines the project quality test entry point. |
| Depends On  | ARCHITECTURE.md, FEATURE-Build.md |
| Provides    | bin/test.sh |
| Consumes    | go vet ./..., go test ./... |

## Purpose

Provide a POSIX-compatible `bin/test.sh` at the application root. It runs `go vet ./...` and `go test ./...` from that root and exits zero only when both pass.

## Programmatic Acceptance

=== AC test-harness ===
Intent: The test entry point runs vet and the Go test suite and succeeds when both succeed.
import os
import subprocess

script = "bin/test.sh"
os.makedirs("bin", exist_ok=True)
result = subprocess.run(
    ["sh", script],
    capture_output=True,
    text=True,
)
assert result.returncode == 0, (result.stdout, result.stderr)
=== END AC test-harness ===

## User Acceptance

- None.

## Guardrails

- The script must not report success when either `go vet ./...` or `go test ./...` fails.
