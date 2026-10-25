# drydock DupFind Uat Case

Drydock uat - DupFind is the simplest Go kit and exercises go build end to end. The source is an original brief for a command-line duplicate file finder that groups regular files by size and SHA-256 digest and reports each group or a summary of reclaimable bytes. The supplied grader, sources/run.sh, builds ./cmd/dupfind with go build, creates a fixed directory tree (hidden files, a symbolic link, empty files, a path with a comma), and runs 15 checks of exit code and output. Standard library only, no database. Run with drydock uat DupFind. Source Repo: none, the brief is original to this kit.

## Intent

Build a command-line duplicate file finder, the Go executable `./cmd/dupfind`, as described in
`sources/INSTRUCTIONS.md`. It walks a directory, groups regular files by size and then by SHA-256
digest, and reports each duplicate group or a one-row summary of the space that removing the
copies would recover.

Correctness is measured by `sources/run.sh`: it builds the executable with `go build` and makes
15 checks of exit code and standard output against a fixed directory tree. The goal is every check
passing.

## What It Does

`cli.Run(args []string, stdout, stderr io.Writer) int` in `internal/cli` interprets the command
line `dupfind [--summary] [--min-size <bytes>] <directory>` and produces a parsed configuration for
later stories.

Rule 3 groups same-size files by the SHA-256 digest of their contents; a group of two or more files
with the same size and digest is a duplicate group. Rule 6 skips unreadable files with one warning.

After successful argument parsing, `cli.Run` checks that `Config.Dir` names an existing directory
before any walking or output.

Rules 4 and 5 fix how paths are written and how groups and paths are ordered so output is
deterministic.

The default output of `dupfind <directory>` is a CSV listing of every file in every duplicate group.

The application root is a Go module whose executable is the main package at `./cmd/dupfind`.
`go build ./cmd/dupfind` must succeed using only the standard library.

Rule 3 of the build instructions says files are grouped by size first and only files that share a
size are read. This feature implements that first stage in `internal/dupes`.

With `--summary` the tool reports the totals of the listing that the same options would produce,
instead of the listing itself.

The application needs one command that decides whether the code is healthy. `bin/test.sh` runs
`go vet ./...` and then `go test ./...` from the application root and exits zero only when both
pass. The Go test suite behind it exercises the finished tool, so a regression in any earlier story
turns this gate red.

`internal/walk` enumerates the files that may be considered for duplicate detection. It never opens
file contents.

## Stack

Go

## Verification

```bash
bash bin/test.sh
```

## Next Steps

- Complete the first-run checks above.
- Run `bash bin/test.sh` before making release or deployment changes.
