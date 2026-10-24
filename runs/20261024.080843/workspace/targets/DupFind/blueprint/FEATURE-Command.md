# FEATURE: Command Operation

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Validates command arguments and directories and completes grader conformance. |
| Depends On  | FEATURE-Listing.md, FEATURE-MinimumSize.md, FEATURE-Summary.md |
| Provides    | `./cmd/dupfind` |
| Consumes    | CSV duplicate listing, `--min-size`, `--summary` CSV |

## Purpose

Define command invocation, argument validation, directory validation, and the final supplied-grader check.

## Command

`dupfind [--summary] [--min-size <bytes>] <directory>`

Options precede the directory, and exactly one directory is required.

## Behavior

- Return `0` on success, whether or not duplicates are found.
- Return `1` when the directory does not exist or is not a directory.
- Return `2` for missing or extra arguments, unknown options, or a `--min-size` value that is not a positive integer.
- On failure, leave standard output empty and write one diagnostic line to standard error.
- The completed implementation must pass the supplied grader, run from the application root as `sh sources/run.sh`.

## Programmatic Acceptance

=== AC command-errors ===
Intent: Missing and extra arguments, unknown options, invalid minimum sizes, and invalid directories return their specified exit codes without producing output.
Cite: sources/INSTRUCTIONS.md:51-66

from drydock_go import cli_runner

cases = [
    ([], 2),
    (["tree", "tree/a"], 2),
    (["--frobnicate", "tree"], 2),
    (["--min-size", "abc", "tree"], 2),
    (["--min-size", "0", "tree"], 2),
    (["tree/no-such-directory"], 1),
    (["tree/env.bak"], 1),
]
for args, expected in cases:
    code, out, err = cli_runner(*args)
    print(out)
    print(err)
    assert code == expected
    assert out == ""
=== END AC command-errors ===

=== AC command-success ===
Intent: A valid invocation returns success and emits parseable CSV.
Cite: sources/INSTRUCTIONS.md:36-50

import csv
import io
from drydock_go import cli_runner

code, out, err = cli_runner("sources")
print(out)
print(err)
assert code == 0
rows = list(csv.reader(io.StringIO(out)))
assert rows
=== END AC command-success ===

## User Acceptance

- None.

## Guardrails

- The application must never delete, move, or modify files under the directory examined.
- Application code must not read the supplied grader or hardcode its answers.
