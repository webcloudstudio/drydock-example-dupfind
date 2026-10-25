# FEATURE: Command Line

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Parses `[--summary] [--min-size <bytes>] <directory>` with options before the directory and maps usage errors to exit code 2. |
| Depends On  | ARCHITECTURE.md |
| Provides    | cli.Run |

## Purpose

Own the command-line contract: argument parsing, the single required directory, the `--summary` and
`--min-size` options, and usage-error reporting. Directory existence checks belong to
FEATURE-TargetValidation.md. Scanning and rendering belong to their own features.

## Behavior

- `cli.Run(args []string, stdout, stderr io.Writer) int` parses arguments with a local `flag.FlagSet` configured with `flag.ContinueOnError`.
- Options are accepted before the directory: `[--summary] [--min-size <bytes>] <directory>`.
- Exactly one directory argument is required.
- An unknown option, missing or extra arguments, or a `--min-size` value that is not a positive integer returns exit code 2, writes exactly one line to stderr, and writes nothing to stdout.
- A valid invocation passes the parsed `summary` flag, `minSize` value (default 1), and directory to the injected execute function and returns that function's code.
- `Run` never calls `os.Exit`; only `cmd/dupfind/main.go` does.

## Interface

- `cli.Run(args []string, stdout, stderr io.Writer, execute func(opts cli.Options) int) int`

## Test Contract

The `internal/cli` package carries one Go test per criterion below, named exactly as stated. Each test
drives `Run` with `bytes.Buffer` writers.

## Programmatic Acceptance

=== AC cli-options-before-directory ===
Intent: Options given before the directory are parsed and the directory is passed to the execute function.

from drydock_go import run_test

result = run_test("./internal/cli", "TestRunOptionsBeforeDirectory")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC cli-options-before-directory ===

=== AC cli-default-min-size ===
Intent: When --min-size is omitted the execute function receives a minimum size of 1.

from drydock_go import run_test

result = run_test("./internal/cli", "TestRunDefaultMinSizeIsOne")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC cli-default-min-size ===

=== AC cli-requires-one-directory ===
Intent: Zero or more than one directory argument returns exit code 2.

from drydock_go import run_test

result = run_test("./internal/cli", "TestRunRequiresExactlyOneDirectory")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC cli-requires-one-directory ===

=== AC cli-unknown-option ===
Intent: An unknown option returns exit code 2 without invoking the execute function.

from drydock_go import run_test

result = run_test("./internal/cli", "TestRunRejectsUnknownOption")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC cli-unknown-option ===

=== AC cli-min-size-positive ===
Intent: A --min-size value that is not a positive integer returns exit code 2 without invoking the execute function.

from drydock_go import run_test

result = run_test("./internal/cli", "TestRunRejectsNonPositiveMinSize")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC cli-min-size-positive ===

=== AC cli-usage-error-output ===
Intent: A usage error writes exactly one stderr line and nothing to stdout.

from drydock_go import run_test

result = run_test("./internal/cli", "TestRunUsageErrorOneStderrLineNoStdout")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC cli-usage-error-output ===

## User Acceptance

- None.

## Guardrails

- `Run` performs no file system access; directory existence is checked only by FEATURE-TargetValidation.md.
- Usage diagnostics go only to stderr and are never written to stdout.
- The parsing and options contract do not change the exit code of a valid invocation; that value is returned by the execute function.
