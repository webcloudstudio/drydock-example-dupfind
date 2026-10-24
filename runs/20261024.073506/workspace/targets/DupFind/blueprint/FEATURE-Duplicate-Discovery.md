# FEATURE: Duplicate Discovery

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Defines recursive, read-only discovery and grouping of duplicate regular files. |
| Depends On  | FEATURE-Command-Interface.md |
| Provides    | duplicate groups |
| Consumes    | dupfind command interface |

## Purpose

Discover duplicate file groups beneath the selected directory. Two files are duplicates only when their byte contents match; size grouping followed by SHA-256 digest grouping implements this rule while avoiding reads for files with no same-size candidate.

## Discovery Rules

- Walk recursively, including hidden files and directories.
- Consider regular files only. Do not follow or report symbolic links.
- Ignore zero-byte files and files below the selected minimum size.
- Group candidates by size before opening file contents. Read contents only for sizes shared by multiple candidates.
- Hash candidate contents with SHA-256 and group files by digest within each size.
- A duplicate group contains at least two files with the same size and digest.
- If a file cannot be read, skip it and emit one warning line on stderr; this does not change the successful exit status.
- Discovery is read-only: it never deletes, moves, or modifies examined files.

## Programmatic Acceptance

=== AC recursively-finds-identical-files ===
Intent: Identical non-empty regular files in nested directories are recognized as duplicates.

from drydock_go import cli_runner
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as tmp:
    root = Path(tmp)
    payload = b"identical bytes"
    left = root / "left"
    right = root / ".hidden" / "right"
    right.parent.mkdir()
    left.write_bytes(payload)
    right.write_bytes(payload)
    code, out, err = cli_runner(str(root))
    print(out)
    print(err)
    assert code == 0
    assert len(out.splitlines()) == 3
=== END AC recursively-finds-identical-files ===

=== AC ignores-links-and-empty-files ===
Intent: Empty files and symbolic links do not contribute duplicate entries.

from drydock_go import cli_runner
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as tmp:
    root = Path(tmp)
    payload = b"target"
    target = root / "target"
    target.write_bytes(payload)
    (root / "link-a").symlink_to(target)
    (root / "link-b").symlink_to(target)
    empty_a = root / "empty-a"
    empty_b = root / "empty-b"
    empty_a.touch()
    empty_b.touch()
    code, out, err = cli_runner(str(root))
    print(out)
    print(err)
    assert code == 0
    assert len(out.splitlines()) == 1
=== END AC ignores-links-and-empty-files===

=== AC unreadable-file-is-skipped ===
Intent: An unreadable candidate is skipped without changing the command's successful status.

from drydock_go import cli_runner
import tempfile
import os
from pathlib import Path

with tempfile.TemporaryDirectory() as tmp:
    root = Path(tmp)
    payload = b"shared"
    (root / "readable").write_bytes(payload)
    unreadable = root / "unreadable"
    unreadable.write_bytes(payload)
    unreadable.chmod(0)
    code, out, err = cli_runner(str(root))
    print(out)
    print(err)
    assert code == 0
    unreadable.chmod(0o600)
=== END AC unreadable-file-is-skipped===

=== AC read-only-preserves-files ===
Intent: Discovery leaves examined file contents unchanged.

from drydock_go import cli_runner
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as tmp:
    root = Path(tmp)
    payload = b"preserve"
    path = root / "one"
    path.write_bytes(payload)
    before = path.read_bytes()
    code, out, err = cli_runner(str(root))
    print(out)
    print(err)
    assert code == 0
    assert path.read_bytes() == before
=== END AC read-only-preserves-files===

## User Acceptance

- None.

## Guardrails

- Do not hash or read a file unless at least one other eligible file has the same size.
- Never follow or report symbolic links.
- Never modify, move, or delete files under the examined directory.
