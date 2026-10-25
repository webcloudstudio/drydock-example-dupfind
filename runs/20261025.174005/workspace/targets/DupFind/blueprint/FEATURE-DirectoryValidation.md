# FEATURE: Directory Validation

| Field       | Value |
|-------------|-------|
| Version     | 20261025 V1 |
| Description | Exit 1 with empty standard output when the directory argument is missing or is not a directory. |
| Depends On  | ARCHITECTURE.md, FEATURE-ArgParsing.md |
| Provides    | cli.validateDirectory |
| Consumes    | cli.Config |

## Purpose

After successful argument parsing, `cli.Run` checks that `Config.Dir` names an existing directory
before any walking or output.

## Behavior

- A nonexistent path, or a path that exists but is not a directory (for example a regular file), returns exit code 1.
- A failure writes nothing to stdout and exactly one diagnostic line to stderr.
- Validation runs after argument parsing and before the walk, so a usage error (exit 2) takes precedence over an operational failure (exit 1).
- An existing directory, empty or not, proceeds and the run succeeds with exit 0.
- The check uses `os.Stat`, so a symbolic link given as the directory argument that resolves to a directory is accepted.

## Reads

File metadata of the directory argument.

## Writes

Standard error only, on failure.

## Programmatic Acceptance

=== AC dirval-missing-directory ===
Intent: A nonexistent directory exits 1 with empty standard output and one stderr line.

import os
import tempfile
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as d:
    missing = os.path.join(d, "no-such-directory")
    for args in ([missing], ["--summary", missing], ["--min-size", "3", missing]):
        code, out, err = cli_runner(*args)
        print(args, code, repr(out), repr(err))
        assert code == 1, args
        assert out == "", args
        assert len(err.splitlines()) == 1, args
=== END AC dirval-missing-directory ===

=== AC dirval-not-a-directory ===
Intent: A regular file given as the directory exits 1 with empty standard output.

import os
import tempfile
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as d:
    path = os.path.join(d, "plain.txt")
    with open(path, "w") as fh:
        fh.write("content")
    for args in ([path], ["--summary", path]):
        code, out, err = cli_runner(*args)
        print(args, code, repr(out), repr(err))
        assert code == 1, args
        assert out == "", args
=== END AC dirval-not-a-directory ===

=== AC dirval-usage-precedes-validation ===
Intent: A usage error on a nonexistent directory still exits 2, not 1.

import os
import tempfile
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as d:
    missing = os.path.join(d, "no-such-directory")
    code, out, err = cli_runner("--min-size", "0", missing)
    print(code, repr(out), repr(err))
    assert code == 2
    assert out == ""
=== END AC dirval-usage-precedes-validation ===

=== AC dirval-existing-directory-succeeds ===
Intent: An existing empty directory is accepted and the run exits 0.

import tempfile
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as d:
    code, out, err = cli_runner(d)
    print(code, repr(out), repr(err))
    assert code == 0
=== END AC dirval-existing-directory-succeeds ===

## User Acceptance

- None.

## Guardrails

- A failing command writes nothing to standard output.
- Validation never creates, modifies, or removes anything on disk.
