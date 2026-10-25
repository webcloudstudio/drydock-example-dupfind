# FEATURE: Target Validation

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Reject a missing or non-directory target argument with exit status 1 and no standard output. |
| Depends On  | ARCHITECTURE.md, FEATURE-CommandLine.md |
| Provides    | cli.ValidateTarget |
| Consumes    | cli.Run |

## Purpose

After argument parsing succeeds, the directory argument is checked. A path that does not exist, or
that exists but is a regular file rather than a directory, is an operational failure: the process
writes exactly one diagnostic line to standard error, writes nothing to standard output, and exits
with status 1.

Validation runs after usage validation, so a usage error (exit 2) takes precedence over an invalid
target.

## Behavior

- Input: the single directory argument returned by `cli.Run` argument parsing.
- A nonexistent path: exit 1, one stderr line, empty stdout.
- A path that exists and is a regular file: exit 1, one stderr line, empty stdout.
- A path that is a directory: validation passes and control returns to the pipeline.
- Symbolic links that resolve to a directory are accepted as the target; the walker governs what is
  examined beneath it.

## Programmatic Acceptance

=== AC target-missing-path ===
Intent: A nonexistent target exits 1, writes exactly one stderr line, and writes nothing to stdout.
Requires: executable=go; scope=test

from drydock_go import cli_runner

code, out, err = cli_runner("tree/no-such-directory")
assert code == 1, (code, out, err)
assert out == "", out
assert len([line for line in err.splitlines() if line.strip()]) == 1, err
=== END AC target-missing-path ===

=== AC target-regular-file ===
Intent: A regular file given as the target exits 1, writes one stderr line, and writes nothing to stdout.
Requires: executable=go; scope=test

import os
import tempfile

from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as d:
    path = os.path.join(d, "plain.txt")
    with open(path, "w") as f:
        f.write("not a directory\n")
    code, out, err = cli_runner(path)
    assert code == 1, (code, out, err)
    assert out == "", out
    assert len([line for line in err.splitlines() if line.strip()]) == 1, err
=== END AC target-regular-file ===

=== AC target-directory-accepted ===
Intent: An existing directory passes validation and produces a listing with exit 0.
Requires: executable=go; scope=test

import os
import tempfile

from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as d:
    code, out, err = cli_runner(d)
    assert code == 0, (code, out, err)
    assert out.splitlines()[0] == "group,size,sha256,path", out
=== END AC target-directory-accepted ===

=== AC target-usage-precedes ===
Intent: A usage error takes precedence over target validation and exits 2, not 1.
Requires: executable=go; scope=test

from drydock_go import cli_runner

code, out, err = cli_runner("--frobnicate", "tree/no-such-directory")
assert code == 2, (code, out, err)
assert out == "", out
=== END AC target-usage-precedes ===

## User Acceptance

- None.

## Guardrails

- The validator only reads the target's metadata. It never creates, modifies, or removes anything
  under the directory examined.
- Diagnostics go to standard error only; standard output stays empty on every failure.
