# FEATURE: Arg Parsing

| Field       | Value |
|-------------|-------|
| Version     | 20261025 V1 |
| Description | Parse --summary, --min-size, and exactly one directory argument, rejecting malformed usage with exit 2. |
| Depends On  | ARCHITECTURE.md |
| Provides    | cli.Run, cli.Config |

## Purpose

`cli.Run(args []string, stdout, stderr io.Writer) int` in `internal/cli` interprets the command
line `dupfind [--summary] [--min-size <bytes>] <directory>` and produces a parsed configuration for
later stories.

## Behavior

- Parsing uses a local `flag.FlagSet` with `flag.ContinueOnError` and its own output discarded.
- `cli.Config` carries `Summary bool`, `MinSize int64`, and `Dir string`. `MinSize` defaults to 1.
- Options precede the directory. Exactly one directory argument is required.
- Usage errors return exit code 2, write nothing to stdout, and write exactly one diagnostic line to stderr. Usage errors are:
  - no arguments, or an option with no directory;
  - more than one positional argument, or an option placed after the directory (it counts as an extra argument);
  - an unknown option, including `-h` and `--help`;
  - `--min-size` with a missing value, or a value that is not a positive decimal integer.
- A `--min-size` value is accepted only when it consists solely of ASCII decimal digits, is at least 1, and fits in an `int64`. Signs (`+5`, `-5`), zero, hexadecimal, exponent forms, and overflow are usage errors.
- Argument errors take precedence over directory validation.

## Reads

None.

## Writes

Standard error only, on usage errors.

## Programmatic Acceptance

=== AC argparse-usage-errors-exit-2 ===
Intent: Every malformed command line exits 2 with empty standard output.

import os
import tempfile
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as d:
    other = os.path.join(d, "other")
    os.mkdir(other)
    bad = [
        [],
        ["--summary"],
        ["--min-size", "5"],
        [d, other],
        ["--frobnicate", d],
        ["-h", d],
        ["--min-size", "abc", d],
        ["--min-size", "0", d],
        ["--min-size", "-5", d],
        ["--min-size", "+5", d],
        ["--min-size", "0x10", d],
        ["--min-size", "1e3", d],
        ["--min-size", "99999999999999999999999", d],
        ["--min-size", d],
        [d, "--summary"],
    ]
    for args in bad:
        code, out, err = cli_runner(*args)
        print(args, code, repr(out), repr(err))
        assert code == 2, args
        assert out == "", args
=== END AC argparse-usage-errors-exit-2 ===

=== AC argparse-diagnostic-on-stderr ===
Intent: A usage error writes exactly one diagnostic line to standard error.

from drydock_go import cli_runner

code, out, err = cli_runner("--min-size", "0", ".")
print(code, repr(out), repr(err))
assert code == 2
assert out == ""
assert len(err.splitlines()) == 1
=== END AC argparse-diagnostic-on-stderr ===

=== AC argparse-valid-forms-not-usage-errors ===
Intent: Well-formed command lines are not rejected as usage errors.

import tempfile
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as d:
    valid = [
        [d],
        ["--summary", d],
        ["--min-size", "5", d],
        ["--summary", "--min-size", "5", d],
        ["--min-size", "5", "--summary", d],
    ]
    for args in valid:
        code, out, err = cli_runner(*args)
        print(args, code, repr(out), repr(err))
        assert code != 2, args
=== END AC argparse-valid-forms-not-usage-errors ===

## User Acceptance

- None.

## Guardrails

- Standard library only; parsing uses `flag` with `flag.ContinueOnError`.
- Only `main` calls `os.Exit`; `Run` returns the code.
- A failing command writes nothing to standard output.
- Never touch any file under the directory examined.
