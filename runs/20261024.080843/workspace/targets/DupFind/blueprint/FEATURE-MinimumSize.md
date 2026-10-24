# FEATURE: Minimum Size Filter

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Excludes files below the configured minimum size from duplicate results. |
| Depends On  | FEATURE-Discovery.md |
| Provides    | `--min-size` |
| Consumes    | duplicate file groups |

## Purpose

Let callers set the minimum file size considered by duplicate discovery.

## Behavior

- Accept `--min-size <bytes>` before the directory argument.
- Require a positive integer. The default is `1`.
- Exclude files smaller than the configured size before forming duplicate results.
- Apply the same filtered duplicate groups to listing and summary output.

## Programmatic Acceptance

=== AC minimum-size-conformance ===
Intent: The supplied grader's minimum-size listing and summary cases pass.
Suite: scoped

import os
import subprocess

result = subprocess.run(
    ["sh", "sources/run.sh"],
    capture_output=True,
    text=True,
    env=os.environ.copy(),
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC minimum-size-conformance ===

## User Acceptance

- None.

## Guardrails

- Do not modify files under the directory being examined.
