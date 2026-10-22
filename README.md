# DupFind

A Drydock UAT kit that exercises the Go toolchain with a duplicate file finder. The brief is
original to this kit. Complexity: 1 of 5, about 4 to 5 stories.

| Path | Role |
|---|---|
| `sources/INSTRUCTIONS.md` | The specification imported into the Target. |
| `bin/run.sh` | The grader. Staged into the build as `sources/run.sh`; `acceptance.full` runs it. |
| `tools/dupfind.py` | The reference that produced the answers. Not imported. |
| `inputs/` | Frozen `SEA_TRIALS.md`, `TECHNOLOGY_STACK.md`, `COMPASS.md`, `METADATA.md`. |

## Grading

`bin/run.sh [APP]` runs `go build ./cmd/dupfind` in `APP` (default: the current directory),
builds a fixed tree in a temporary directory, and makes 15 checks of exit code and standard
output. The tree has three copies of one file, a pair whose paths sort by byte value
(`img1 (copy).jpg` before `img1.jpg`), a path containing a comma, a hidden pair, two 1000-byte files plus a
same-size non-duplicate, two empty files, and a symbolic link.

| Group | Checks |
|---|---|
| A | full listing, by relative and absolute directory |
| B | `--summary` |
| C | `--min-size` 20 and 200, with and without `--summary` |
| D | a directory with no duplicates; empty files only |
| E | missing directory, file given as directory |
| U | usage errors |

Exit `0` only when the build succeeds and all 15 pass. No data files; the grader creates the
tree, so the answers depend only on the bytes it writes.
