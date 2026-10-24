# Build Instructions: A Duplicate File Finder

## Objective

Build a command-line tool that finds duplicate files under a directory. Two files are duplicates
when their contents are byte-for-byte identical. The tool reports each group of duplicates and the
disk space that deleting all but one copy of each would recover. It never deletes, moves, or
modifies a file.

## Stack

Implement in Go 1.22 or newer using only the standard library. The application root is a Go
module, and its executable is the main package at `./cmd/dupfind`. Tests run under `go test`.

## Run Harness

`sources/run.sh` is the scoring entry point. It is supplied, not authored, and is staged verbatim
into the application root. `drydock uat` runs `sh sources/run.sh` from the completed application
root and takes its exit code and output as the score.

The harness runs `go build -o <temporary path> ./cmd/dupfind`, creates a fixed directory tree in a
temporary directory, and then runs 15 checks. Each check compares one command's exit code and
standard output with a known answer. Standard error is never compared. The harness prints `PASS`
or `FAIL` per check and ends with:

```text
DupFind: NN passed, NN failed (of 15)
```

`sources/run.sh` is the exam. Do not modify, rewrite, regenerate, or substitute it. Application
code never reads it and never embeds an answer from it.

## Command

```text
dupfind [--summary] [--min-size <bytes>] <directory>
```

Options precede the directory. Exactly one directory is required.

| Option | Meaning |
|---|---|
| `--summary` | Write the summary instead of the listing. |
| `--min-size <bytes>` | Ignore files smaller than `<bytes>`. A positive integer; default `1`. |

## Rules

1. **Scope.** The tool walks `<directory>` recursively, including hidden files and directories
   (names beginning `.`). It considers regular files only. It neither follows nor reports
   symbolic links.
2. **Empty files are never duplicates.** A file of zero bytes is always ignored, which the
   default `--min-size 1` guarantees.
3. **Duplicate test.** Files are grouped by size first. Only files that share a size are read.
   Within a size, files are grouped by the SHA-256 digest of their contents. A group of two or
   more files with the same size and digest is a duplicate group.
4. **Paths** are written relative to `<directory>`, with `/` as the separator and no leading
   `./`.
5. **Order.** Groups are ordered by reclaimable bytes descending, then by digest ascending.
   A group's reclaimable bytes are `size × (files − 1)`. Groups are numbered from 1 in that order.
   Within a group, paths are ordered by byte value ascending.
6. **Unreadable files.** A file that cannot be read is skipped with one warning line on standard
   error. It does not change the exit code.

## Output

Output is CSV on standard output: a header row, then data rows, with `\n` line endings. A field
that contains a comma, a double quote, or a newline is quoted as RFC 4180 specifies. Go's
`encoding/csv` writer produces this format.

**Listing** (default): header `group,size,sha256,path`, then one row per file in every duplicate
group. `sha256` is the lowercase hexadecimal digest. With no duplicates, only the header is
written.

```text
group,size,sha256,path
1,17,8a3c67892f82af22b58377b8e85bb41899053788d3aa3b81d2e6b4580e917c29,a/report.txt
1,17,8a3c67892f82af22b58377b8e85bb41899053788d3aa3b81d2e6b4580e917c29,b/report-copy.txt
```

**Summary** (`--summary`): header `groups,duplicate_files,reclaimable_bytes`, then one row.
`duplicate_files` is the sum of `files − 1` over all groups. `reclaimable_bytes` is the sum of
every group's reclaimable bytes. With no duplicates the row is `0,0,0`.

A failing command writes nothing to standard output and one diagnostic line to standard error.

## Exit codes

| Code | Meaning |
|---|---|
| `0` | Success, whether or not duplicates were found. |
| `1` | Operational failure: `<directory>` does not exist or is not a directory. |
| `2` | Usage error: missing or extra arguments, unknown option, `--min-size` not a positive integer. |

## Tests

The completed application provides a POSIX-compatible `bin/test.sh` that runs `go vet ./...` and
`go test ./...` from the application root and exits zero only when both pass. The suite covers
each rule, each output form, and each exit code against real temporary directories.
