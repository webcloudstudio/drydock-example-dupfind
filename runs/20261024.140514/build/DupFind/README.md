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

The application root is a Go module. Building from that root with `go build ./cmd/dupfind` produces the dupfind command. The main package lives in `cmd/dupfind` and uses only the standard library.

The default command output is a CSV listing of every file in every duplicate group. It uses the ordered groups and relative paths provided by `FEATURE-Ordering.md`.

When `--summary` is supplied, write a single CSV row summarizing the duplicate groups that the listing would produce with the same options.

Provide a POSIX-compatible `bin/test.sh` at the application root. It runs `go vet ./...` and `go test ./...` from that root and exits zero only when both pass.

## Stack

Go

## Verification

```bash
bash bin/test.sh
```

## Next Steps

- Complete the first-run checks above.
- Run `bash bin/test.sh` before making release or deployment changes.
