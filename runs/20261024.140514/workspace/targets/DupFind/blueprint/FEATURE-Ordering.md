# FEATURE: Ordering

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Orders duplicate groups and their relative paths deterministically. |
| Depends On  | ARCHITECTURE.md, FEATURE-Scanning.md |
| Provides    | ordered duplicate groups and relative paths |
| Consumes    | duplicate groups from a directory tree |

## Ordering contract

Sort groups by reclaimable bytes descending, then by SHA-256 digest ascending. Reclaimable bytes equal file size times the number of duplicate copies beyond the first. Number groups starting at 1 in sorted order. Sort paths inside each group by byte value ascending. Emit each path relative to the requested directory using `/` separators and without a leading `./`.

## Programmatic Acceptance

=== AC ordering-conformance ===
Intent: The supplied conformance grader passes the duplicate ordering cases.

Suite: scoped

import subprocess

result = subprocess.run(
    ["sh", "sources/run.sh"],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC ordering-conformance ===

## User Acceptance

- None.

## Guardrails

- Ordering is deterministic and independent of filesystem traversal order.
- Do not include files outside the requested directory.
