# FEATURE: Duplicate Listing

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Renders duplicate file groups as ordered CSV rows. |
| Depends On  | FEATURE-Discovery.md |
| Provides    | CSV duplicate listing |
| Consumes    | duplicate file groups |

## Purpose

Render the duplicate groups discovered by `FEATURE-Discovery.md` as a CSV listing.

## Behavior

- Write the header `group,size,sha256,path`, followed by one row for each file in every duplicate group.
- Use lowercase hexadecimal SHA-256 digests.
- Write paths relative to the scanned directory, using `/` separators and no leading `./`.
- Order groups by reclaimable bytes descending, then digest ascending.
- Number groups from 1 in that order.
- Order paths within each group by byte value ascending.
- Use `encoding/csv` so fields containing commas, double quotes, or newlines follow RFC 4180 quoting. Use `\n` line endings.
- With no duplicate groups, write only the header.

## Programmatic Acceptance

=== AC listing-conformance ===
Intent: The supplied grader's listing, ordering, path, CSV, and empty-result cases pass.
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
=== END AC listing-conformance ===

## User Acceptance

- None.

## Guardrails

- Never modify files under the directory being examined.
- Do not read or depend on `sources/run.sh` from application code.
