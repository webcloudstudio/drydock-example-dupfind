# FEATURE: Command Interface

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Defines DupFind's command options, argument validation, and operational error behavior. |
| Depends On  | ARCHITECTURE.md |
| Provides    | dupfind command interface |
| Consumes    | ./cmd/dupfind |

## Purpose

The command is `dupfind [--summary] [--min-size <bytes>] <directory>`. Options precede exactly one directory argument. `--summary` selects summary CSV output. `--min-size` takes a positive integer; files smaller than it are ignored, and the default is 1 byte.

## Behavior

- Accept options before the directory and reject unknown options, missing values, non-positive or non-integer minimum sizes, a missing directory argument, and extra arguments with usage exit code 2.
- A nonexistent path or a path that is not a directory exits with code 1.
- On success, return exit code 0 whether or not duplicates exist.
- On failure, write nothing to stdout and one diagnostic line to stderr.
- The process entry point is `./cmd/dupfind`; only `main` calls `os.Exit`. CLI behavior is exposed through `Run(args []string, stdout, stderr io.Writer) int`.

## Guardrails

- Do not print output to stdout for malformed arguments or invalid directory paths.
- Do not write diagnostic text to stdout.
- Do not change files in the directory being examined.

## Programmatic Acceptance

=== AC command-options ===
Intent: The command accepts the summary option and a positive minimum size before its directory argument.

from drydock_go import cli_runner
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as tmp:
    root = Path(tmp)
    a = root / "a"
    b = root / "b"
    a.write_bytes(b"sample")
    b.write_bytes(b"sample")
    code, out, err = cli_runner("--summary", "--min-size", "2", str(root))
    print(out)
    print(err)
    assert code == 0
=== END AC command-options ===

=== AC usage-exit ===
Intent: Malformed and incomplete invocations return the usage exit status.

from drydock_go import cli_runner

cases = [
    [],
    ["one", "two"],
    ["--min-size", "0", "dir"],
    ["--frobnicate", "dir"],
]

for args in cases:
    code, out, err = cli_runner(*args)
    print(out)
    print(err)
    assert code == 2
=== END AC usage-exit ===

=== AC directory-errors ===
Intent: Missing and non-directory paths return the operational failure status.

from drydock_go import cli_runner
import tempfile
from pathlib import Path

with tempfile.TemporaryDirectory() as tmp:
    root = Path(tmp)
    file_path = root / "file"
    file_path.write_bytes(b"x")
    missing = root / "missing"
    for path in (missing, file_path):
        code, out, err = cli_runner(str(path))
        print(out)
        print(err)
        assert code == 1
        assert out == ""
=== END AC directory-errors ===

=== AC minimum-size-validation ===
Intent: A non-integer minimum size is rejected with the usage exit status.

from drydock_go import cli_runner

code, out, err = cli_runner("--min-size", "abc", "dir")
print(out)
print(err)
assert code == 2
=== END AC minimum-size-validation ===

## User Acceptance

- None.

## Guardrails

- None.
