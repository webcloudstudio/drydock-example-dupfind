# FEATURE: Test Gate

| Field       | Value |
|-------------|-------|
| Version     | 20261025 V1 |
| Description | Provide the POSIX bin/test.sh launcher and a Go test suite that covers every rule, output form, and exit code of dupfind against real temporary directories. |
| Depends On  | FEATURE-ModuleBuild.md, FEATURE-DirectoryValidation.md, FEATURE-Listing.md, FEATURE-Summary.md |
| Provides    | bin/test.sh |
| Consumes    | cli.Run |

## Purpose

The application needs one command that decides whether the code is healthy. `bin/test.sh` runs
`go vet ./...` and then `go test ./...` from the application root and exits zero only when both
pass. The Go test suite behind it exercises the finished tool, so a regression in any earlier story
turns this gate red.

## Trigger

A developer, the build engine, or the terminal verification story runs `sh bin/test.sh`.

## Sequence

1. `bin/test.sh` is POSIX `sh`: no bashisms, no arrays, no `[[ ]]`.
2. It derives the application root from the script's own location (`dirname "$0"`/`..`), changes
   into it, so the result is the same from any working directory.
3. It runs `go vet ./...`. A non-zero status stops the script with that non-zero status.
4. It runs `go test ./...`. Its status becomes the script's exit status.
5. The script exits `0` only when both commands succeeded.

## Test Suite Coverage

The Go test suite lives in `*_test.go` files beside the packages under `internal/`. Tests call
`cli.Run(args, &stdout, &stderr)` with `bytes.Buffer` writers over real directories made with
`t.TempDir()`, and parse standard output with `encoding/csv`. Expectations are written against
observed real output. Each rule has at least one test named for the behavior:

| Area | Behavior covered |
|---|---|
| Scope | Hidden files and hidden directories are walked; nested directories are walked |
| Scope | A symbolic link to a file is neither followed nor reported; a link to a directory is not traversed |
| Empty files | Zero-byte files are never reported, including with an explicit small `--min-size` |
| Minimum size | Files below `--min-size` are ignored; the default is 1 |
| Duplicate test | Same size but different content is not a duplicate; same size and content is |
| Duplicate test | Files with a unique size are never opened for hashing |
| Paths | Relative to the directory, `/` separator, no leading `./` |
| Order | Groups by reclaimable bytes descending then digest ascending, numbered from 1; paths by byte value |
| Unreadable | An unreadable file is skipped, one warning on stderr, exit code unchanged (skipped when running as root) |
| Listing | Header `group,size,sha256,path`; header only when there are no duplicates |
| CSV quoting | A path containing a comma is quoted and round-trips through `encoding/csv` |
| Summary | Header `groups,duplicate_files,reclaimable_bytes`; `0,0,0` when there are no duplicates; totals equal those of the listing under the same options |
| Exit codes | `0` for success with or without duplicates; `1` for a missing directory or a regular file; `2` for no arguments, extra arguments, unknown option, `--min-size` of `abc` and of `0` |
| Failure output | Exit codes 1 and 2 write nothing to standard output |

## Reads

- The application source tree, through `go vet` and `go test`.

## Writes

- Nothing outside temporary directories created by the tests themselves. The suite never writes
  under any directory the tool examines other than the temporary fixture it created.

## Programmatic Acceptance

=== AC test-script-exits-zero ===
Intent: bin/test.sh runs go vet and go test on the finished application and exits zero.
Requires: executable=sh; scope=test
Requires: executable=go; scope=test

import subprocess
import sys

result = subprocess.run(["sh", "bin/test.sh"], capture_output=True, text=True)
print(result.stdout)
print(result.stderr, file=sys.stderr)
assert result.returncode == 0, result.returncode
=== END AC test-script-exits-zero ===

=== AC test-script-posix-syntax ===
Intent: bin/test.sh exists and parses as POSIX sh without executing.
Requires: executable=sh; scope=test

import os
import subprocess
import sys

script = os.path.join("bin", "test.sh")
assert os.path.isfile(script)
result = subprocess.run(["sh", "-n", script], capture_output=True, text=True)
print(result.stdout)
print(result.stderr, file=sys.stderr)
assert result.returncode == 0, result.returncode
=== END AC test-script-posix-syntax ===

=== AC test-script-fails-on-failing-test ===
Intent: bin/test.sh exits non-zero when any Go test fails, proven on a temporary copy of the application.
Requires: executable=sh; scope=test
Requires: executable=go; scope=test

import os
import shutil
import subprocess
import sys
import tempfile

root = os.getcwd()
failing_test = 'package cli\n\nimport "testing"\n\nfunc TestForcedFailure(t *testing.T) {\n\tt.Fatal("forced")\n}\n'
with tempfile.TemporaryDirectory() as tmp:
    dst = os.path.join(tmp, "app")
    shutil.copytree(root, dst, ignore=shutil.ignore_patterns(".git"))
    with open(os.path.join(dst, "internal", "cli", "zz_forced_failure_test.go"), "w") as handle:
        handle.write(failing_test)
    result = subprocess.run(["sh", "bin/test.sh"], cwd=dst, capture_output=True, text=True)
    print(result.stdout)
    print(result.stderr, file=sys.stderr)
    assert result.returncode != 0, result.returncode
=== END AC test-script-fails-on-failing-test ===

=== AC test-script-independent-of-cwd ===
Intent: bin/test.sh locates the application root itself, so running it from another directory still exits zero.
Requires: executable=sh; scope=test
Requires: executable=go; scope=test

import os
import subprocess
import sys
import tempfile

script = os.path.abspath(os.path.join("bin", "test.sh"))
with tempfile.TemporaryDirectory() as elsewhere:
    result = subprocess.run(["sh", script], cwd=elsewhere, capture_output=True, text=True)
    print(result.stdout)
    print(result.stderr, file=sys.stderr)
    assert result.returncode == 0, result.returncode
=== END AC test-script-independent-of-cwd ===

=== AC suite-has-tests-for-cli ===
Intent: The Go test suite contains test files for the cli package and they pass when run alone.
Requires: executable=go; scope=test

import glob
import os
import subprocess
import sys

test_files = glob.glob(os.path.join("internal", "cli", "*_test.go"))
assert len(test_files) >= 1, test_files
result = subprocess.run(
    ["go", "test", "-count=1", "./internal/cli"], capture_output=True, text=True
)
print(result.stdout)
print(result.stderr, file=sys.stderr)
assert result.returncode == 0, result.returncode
=== END AC suite-has-tests-for-cli ===

=== AC built-binary-exit-codes ===
Intent: The built executable exits 0 for a directory, 1 for a missing path, and 2 for each malformed argument form the suite covers.
Requires: executable=go; scope=test

import os
import tempfile
from drydock_go import cli_runner

with tempfile.TemporaryDirectory() as tree:
    payload = "same bytes\n"
    for name in ("one.txt", "two.txt"):
        with open(os.path.join(tree, name), "w") as handle:
            handle.write(payload)
    missing = os.path.join(tree, "no-such-directory")
    regular = os.path.join(tree, "one.txt")

    ok_code, ok_out, ok_err = cli_runner(tree)
    assert ok_code == 0, (ok_code, ok_out, ok_err)

    missing_code, missing_out, missing_err = cli_runner(missing)
    assert missing_code == 1, (missing_code, missing_out, missing_err)
    assert missing_out == ""

    file_code, file_out, file_err = cli_runner(regular)
    assert file_code == 1, (file_code, file_out, file_err)
    assert file_out == ""

    for args in ([], [tree, tree], ["--min-size", "abc", tree], ["--min-size", "0", tree], ["--frobnicate", tree]):
        code, out, err = cli_runner(*args)
        print(args, code, err)
        assert code == 2, (args, code, out, err)
        assert out == "", (args, out)
=== END AC built-binary-exit-codes ===

## User Acceptance

- None.

## Guardrails

- `bin/test.sh` is POSIX `sh`; it does not require bash.
- The test suite never reads `sources/run.sh` and never embeds a digest, path, or total taken from it.
- Tests write only inside temporary directories they create; nothing under a directory the tool examines is deleted, moved, or modified by the tool.
- The tool is built with the standard library only; no third-party module is added for testing.
