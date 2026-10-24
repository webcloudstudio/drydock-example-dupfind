# FEATURE: Duplicate Summary

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Reports duplicate-group totals as a CSV summary. |
| Depends On  | FEATURE-Discovery.md, FEATURE-Listing.md, FEATURE-MinimumSize.md |
| Provides    | `--summary` CSV |
| Consumes    | duplicate file groups, `--min-size` |

## Purpose

Provide a compact report of duplicate groups, redundant file counts, and reclaimable bytes.

## Behavior

- With `--summary`, write the header `groups,duplicate_files,reclaimable_bytes` and exactly one data row.
- Compute `groups` as the number of duplicate groups.
- Compute `duplicate_files` as the sum of `files - 1` over each group.
- Compute `reclaimable_bytes` as the sum of `size × (files - 1)` over each group.
- When there are no duplicate groups, report zero for all totals.
- Use the same groups selected by the active options, including `--min-size`.

## Programmatic Acceptance

=== AC summary-conformance ===
Intent: The supplied grader's summary and minimum-size summary cases pass.
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
=== END AC summary-conformance ===

## User Acceptance

- None.

## Guardrails

- Never modify files under the directory being examined.
