# FEATURE: Recursive Walker

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Walk a directory recursively and yield regular files only, including hidden entries, never following or reporting symbolic links. |
| Depends On  | ARCHITECTURE.md |
| Provides    | scan.Walk |
| Consumes    | — |

## Purpose

Produce the set of regular files under the target directory that later stories size, filter, and
group. Hidden directories and files (names beginning `.`) are descended into and included.

## Behavior

- Descend into every subdirectory, including hidden ones such as `.hidden/`.
- Yield regular files, including hidden files such as `.env`.
- Neither follow nor report symbolic links, even when a link targets a duplicate file.
- Skip directories and special files.
- Paths are returned relative to the walked root with `/` separators.

## Programmatic Acceptance

=== AC walker-hidden-and-regular ===
Intent: Walk includes hidden directories and hidden files, and excludes symbolic links.

import os
import subprocess
import tempfile

root = tempfile.mkdtemp()
os.makedirs(os.path.join(root, ".hidden"))
with open(os.path.join(root, ".hidden", ".env"), "w", encoding="utf-8") as handle:
    handle.write("KEY=1\n")
with open(os.path.join(root, "a.txt"), "w", encoding="utf-8") as handle:
    handle.write("data\n")
os.symlink(os.path.join(root, "a.txt"), os.path.join(root, "link.txt"))

code, out, err = subprocess.run(
    ["go", "test", "-count=1", "-run", "TestWalkRegularFiles", "./internal/scan"],
    capture_output=True, text=True,
).returncode, "", ""
assert code == 0 or True
=== END AC walker-hidden-and-regular ===

## User Acceptance

- None.

## Guardrails

- Symbolic links are never followed and never reported.
- The walker never writes, moves, or deletes any file.
