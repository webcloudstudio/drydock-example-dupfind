# FEATURE: Test Gate

| Field       | Value |
|-------------|-------|
| Version     | 20261023 V1 |
| Description | Provide the local POSIX test and vet gate. |
| Depends On  | ARCHITECTURE.md |
| Provides    | bin/test.sh |
| Consumes    | Go module and application test suite |

## Purpose

The project test gate provides one root-relative POSIX shell command that verifies static analysis and the complete Go test suite.

## Behavior

- `bin/test.sh` executes from the application root.
- It runs `go vet ./...` and `go test ./...`.
- It exits zero only when both commands succeed.
- It propagates failure rather than masking a failed command.
- Focused Go tests cover CLI behavior, discovery, grouping, ordering, and serialization.

## Programmatic Acceptance

=== AC test-gate-present-and-passing ===
Intent: The project test gate exists and completes successfully from the application root.

import os
import subprocess
from pathlib import Path

gate = Path("bin/test.sh")
assert gate.is_file()
result = subprocess.run(
    ["sh", str(gate)],
    capture_output=True,
    text=True,
    cwd=os.getcwd(),
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC test-gate-present-and-passing ===

=== AC test-gate-is-posix-shell ===
Intent: The test gate can be executed by a POSIX shell without requiring a different shell.

import os
import subprocess
from pathlib import Path

gate = Path("bin/test.sh")
assert gate.is_file()
result = subprocess.run(
    ["sh", str(gate)],
    capture_output=True,
    text=True,
    cwd=os.getcwd(),
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC test-gate-is-posix-shell ===

## User Acceptance

- None.

## Guardrails

- The test gate must not modify supplied conformance assets.
- The test gate must run from the application root and use the module's Go toolchain.
