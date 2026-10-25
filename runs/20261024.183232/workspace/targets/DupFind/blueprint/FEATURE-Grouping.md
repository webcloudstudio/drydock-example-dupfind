# FEATURE: Duplicate Grouping

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Group regular files by size, then by SHA-256 digest, reading only size-matched files. |
| Depends On  | FEATURE-SizeFilter.md |
| Provides    | scan.Group |
| Consumes    | scan.FilterSize |

## Purpose

Identify byte-identical files. A duplicate group is two or more files sharing both size and
SHA-256 digest.

## Behavior

- Bucket files by size. Read a file's contents only when another file shares its size.
- Within a size bucket, group by lowercase hexadecimal SHA-256 digest.
- A group of two or more files with equal size and digest is a duplicate group.
- Open files read-only.

## Programmatic Acceptance

=== AC grouping-duplicate-groups ===
Intent: Files with identical bytes form one group; files with differing bytes do not.

import os
import subprocess
import tempfile

root = tempfile.mkdtemp()
with open(os.path.join(root, "a.txt"), "w", encoding="utf-8") as handle:
    handle.write("same\n")
with open(os.path.join(root, "b.txt"), "w", encoding="utf-8") as handle:
    handle.write("same\n")
with open(os.path.join(root, "c.txt"), "w", encoding="utf-8") as handle:
    handle.write("diff\n")

result = subprocess.run(
    ["go", "test", "-count=1", "-run", "TestGroupDuplicates", "./internal/scan"],
    capture_output=True, text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC grouping-duplicate-groups ===

## User Acceptance

- None.

## Guardrails

- Files are opened read-only and never modified, moved, or deleted.
- Contents are read only for files whose size matches another file.
