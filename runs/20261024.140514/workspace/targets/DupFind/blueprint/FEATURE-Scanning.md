# FEATURE: Scanning

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Finds exact duplicate regular files by scanning a directory tree. |
| Depends On  | ARCHITECTURE.md, FEATURE-Command.md |
| Provides    | duplicate groups from a directory tree |
| Consumes    | dupfind command-line options and usage |

## Scope and filtering

Walk the requested directory recursively, including hidden files and directories. Examine regular files only. Do not follow or report symbolic links. Ignore files smaller than the selected minimum size; zero-byte files are therefore ignored under the default minimum of 1.

## Candidate grouping and reads

Partition eligible files by size. Read file contents only for sizes represented by at least two eligible files. Compute SHA-256 for each readable candidate and group files with the same size and digest. Report groups only when at least two files share both. If a candidate cannot be read, skip it, write one warning line to standard error, and keep the command's success exit status.

## Read-only behavior

The scan only reads files. It never deletes, moves, or writes files in the examined directory. It does not use fixture-specific paths or expected results.

## Programmatic Acceptance

=== AC scan-groups-exact-files ===
Intent: The supplied conformance grader passes its duplicate scanning cases.

Suite: scoped

import subprocess

result = subprocess.run(
    ["sh", "sources/run.sh"],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC scan-groups-exact-files ===

=== AC scan-symlinks-and-empty ===
Intent: The supplied conformance grader passes its symbolic-link and empty-file cases.

Suite: scoped

import subprocess

result = subprocess.run(
    ["sh", "sources/run.sh"],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC scan-symlinks-and-empty ===

=== AC scan-read-only ===
Intent: The supplied conformance grader passes while examined files remain untouched.

Suite: scoped

import subprocess

result = subprocess.run(
    ["sh", "sources/run.sh"],
    capture_output=True,
    text=True,
)
print(result.stdout)
print(result.stderr)
assert result.returncode == 0
=== END AC scan-read-only ===

## User Acceptance

- None.

## Guardrails

- Never follow or report symbolic links.
- Never write to the examined directory.
