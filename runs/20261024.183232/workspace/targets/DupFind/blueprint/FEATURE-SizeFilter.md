# FEATURE: Size Filter

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Drop zero-byte files and files smaller than the configured minimum size. |
| Depends On  | FEATURE-Walker.md |
| Provides    | scan.FilterSize |
| Consumes    | scan.Walk |

## Purpose

Remove files that can never be duplicates of interest before any contents are read.

## Behavior

- Exclude every file whose size is 0.
- Exclude every file whose size is smaller than the configured minimum; the default minimum is 1 byte.
- Operate as a pure function over the walker output.

## Programmatic Acceptance

```python
import tempfile
from pathlib import Path

from drydock_go import cli_runner


def test_zero_byte_file_is_excluded_from_stdout():
    with tempfile.TemporaryDirectory() as tmp:
        root = Path(tmp)
        (root / "empty.bin").write_bytes(b"")
        code, out, err = cli_runner(str(root))
        assert code == 0
        assert "empty.bin" not in out


def test_non_empty_duplicates_survive_the_size_filter():
    with tempfile.TemporaryDirectory() as tmp:
        root = Path(tmp)
        (root / "empty.bin").write_bytes(b"")
        (root / "a.bin").write_bytes(b"same")
        (root / "b.bin").write_bytes(b"same")
        code, out, err = cli_runner(str(root))
        assert code == 0
        assert "a.bin" in out
        assert "b.bin" in out
        assert "empty.bin" not in out
```

## User Acceptance

- None.

## Guardrails

- The filter never reads file contents.
