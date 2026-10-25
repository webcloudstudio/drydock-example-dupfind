# Blueprint Analysis: drydock DupFind Uat Case

## Commander Expectations

- assert the product is a command-line tool that finds byte-identical duplicate regular files under a directory
- assert the tool never deletes, moves, or modifies any file it examines
- assert the built executable is `dupfind`, produced from `./cmd/dupfind` with `go build`, using only the Go standard library
- assert the supplied grader `sources/run.sh` exits zero against the completed application

## Crew

| Crew | Charge |
|---|---|
| Commander | Defines intent and decides what done means. |
| Team Lead | Confirms epic completeness and stakeholder expectations. |
| Planning Crew | Authors atomic specifications and the ordered Manifest. |
| Shipyard Crew | Builds the tickets without synchronous Commander access. |

## Story List

features: 4
stories: 11

### Feature: Scan and Group

| ID | Story | High-level AC |
|---|---|---|
| SCAN-001 | Walk the directory recursively, including hidden entries, considering regular files only | Walk includes `.hidden/` and `.env`; symlinks and directories are not reported as files |
| SCAN-002 | Ignore zero-byte files and files smaller than `--min-size` | Files of size 0 never appear; `--min-size N` excludes files smaller than N bytes |
| SCAN-003 | Group files by size, then by SHA-256 digest of contents, reading only size-matched files | Only groups of two or more files sharing size and digest form duplicate groups |
| SCAN-004 | Skip unreadable files with one warning line on stderr without changing exit code | An unreadable file produces one stderr line and the run still exits per the normal rules |

### Feature: Listing Output

| ID | Story | High-level AC |
|---|---|---|
| LIST-001 | Emit the duplicate listing as CSV with header `group,size,sha256,path` | Header row is exact; one row per file in every group; RFC 4180 quoting applied |
| LIST-002 | Order groups by reclaimable bytes descending then digest ascending, number groups from 1, and order paths within a group by byte value | Group numbers follow order; paths are relative to the directory, `/`-separated, no leading `./` |
| LIST-003 | Emit only the header when no duplicates exist | Output is exactly `group,size,sha256,path` followed by `\n` |

### Feature: Summary Output

| ID | Story | High-level AC |
|---|---|---|
| SUM-001 | Emit `--summary` CSV with header `groups,duplicate_files,reclaimable_bytes` and one totals row | Totals equal sums over the same listing; no duplicates yields `0,0,0` |

### Feature: Command Line Interface

| ID | Story | High-level AC |
|---|---|---|
| CLI-001 | Parse `[--summary] [--min-size <bytes>] <directory>` with options before the directory | Exactly one directory required; unknown options, missing or extra arguments, and non-positive-integer `--min-size` exit 2 with no stdout |
| CLI-002 | Report a missing or non-directory target with exit 1 and no stdout | Nonexistent path or a regular file as target exits 1 with one stderr line and empty stdout |
| CLI-003 | Ship a POSIX-compatible `bin/test.sh` running `go vet ./...` and `go test ./...` | Script exits zero only when both pass |

## Story Guidance

Named stories planning must produce. Authoritative record: `STORY_GUIDANCE.json`.

None.

## Surfaced Acceptance Criteria

| ID | Story ID | Criterion |
|---|---|---|
| AC-001 | SCAN-001 | The walker descends into hidden directories and includes hidden files such as `.env` |
| AC-002 | SCAN-001 | Symbolic links are neither followed nor reported, even when they point at a duplicate file |
| AC-003 | SCAN-004 | An unreadable file is named in exactly one stderr warning line and produces no stdout row |
| AC-004 | LIST-001 | Fields containing commas, double quotes, or newlines are quoted per RFC 4180 |
| AC-005 | LIST-002 | Paths in listing use `/` separators with no leading `./`, sorted by byte value within each group |
| AC-006 | CLI-001 | Usage errors write exactly one diagnostic line to stderr and nothing to stdout |
| AC-007 | CLI-002 | Operational failures write exactly one diagnostic line to stderr and nothing to stdout |
| AC-008 | SUM-001 | Summary totals are computed from the same grouping the listing would produce for identical options |

## Source Inventory

| Path | Content kind | Disposition | Reason |
|---|---|---|---|
| `sources/INSTRUCTIONS.md` | markdown | analyzed | readable UTF-8 |
| `sources/run.sh` | code | analyzed | readable UTF-8 |

## Relationship Model

| Source or group | Relationship type | Related source or group | Evidence | Delivery implication |
|---|---|---|---|---|
| `sources/run.sh` | instruction-to-test | `sources/INSTRUCTIONS.md` | run.sh checks A1–U5 encode the Rules, Output, and Exit codes sections | Acceptance is the supplied grader; all stories must satisfy its 15 checks |
| `sources/INSTRUCTIONS.md` Rules 1–6 | dependency | Scan and Group stories | Scope, empty files, duplicate test, paths, order, unreadable files | Scan and Group must land before Listing and Summary |
| `sources/INSTRUCTIONS.md` Output | dependency | Listing Output and Summary Output | CSV header and row shapes | Both output stories consume the same grouping result |
| `sources/INSTRUCTIONS.md` Exit codes | dependency | Command Line Interface | Codes 0, 1, 2 mapped to failure modes | CLI validation precedes directory checks |
| `sources/INSTRUCTIONS.md` Tests | instruction-to-test | `bin/test.sh` | Requires `go vet ./...` and `go test ./...` | CLI-003 gates on all implementation stories |

## Source Roles

| Path | Role | Plan disposition | Build disposition |
|---|---|---|---|
| `sources/INSTRUCTIONS.md` | normative specification | compass | prompt-only |
| `sources/run.sh` | conformance harness | context | stage |

## Planning Instructions

### Delivery Shape

The system is a single-binary Go CLI. Input: one directory path and options. Output: CSV on stdout (listing or summary), diagnostics on stderr, exit codes 0, 1, 2. Execution flow: parse arguments, validate the directory, walk recursively, filter by size, group by size then SHA-256, order, render CSV.

### Story Realization Map

| Story ID | Durable scope | Evidence | Related files | Requirement type |
|---|---|---|---|---|
| SCAN-001 | Recursive walker over regular files | INSTRUCTIONS.md Rule 1; run.sh `.hidden`, `env.bak`, `report-link.txt` | `cmd/dupfind/`, walker package | capability |
| SCAN-002 | Size filter | Rule 2; run.sh `empty/`, C1–C3 | walker/filter | capability |
| SCAN-003 | Size then SHA-256 grouping | Rule 3; run.sh A1, A2 | grouping package | capability |
| SCAN-004 | Unreadable-file warning | Rule 6 | grouping/reader | capability |
| LIST-001 | CSV listing writer | Output listing; run.sh A1 | output package | capability |
| LIST-002 | Ordering and path rendering | Rules 4–5; run.sh A1 `"music/song,live.mp3"` | ordering, output | capability |
| LIST-003 | Empty listing | Output; run.sh D1 | output | capability |
| SUM-001 | Summary totals | Output summary; run.sh B1, C2, D2 | output | capability |
| CLI-001 | Argument parsing and usage errors | Command section; Exit codes 2; run.sh U1–U5 | `cmd/dupfind/main.go` | capability |
| CLI-002 | Target validation | Exit codes 1; run.sh E1, E2 | `cmd/dupfind/main.go` | capability |
| CLI-003 | Test script | Tests section | `bin/test.sh` | test harness |

### Test and Acceptance Strategy

Focused story tests run in `go test` against real temporary directories, bounded to each story's slice. The terminal verification story CLI-003 depends on every implementation story and runs the full suite through `bin/test.sh`. The supplied `sources/run.sh` is the final acceptance gate and is referenced only through `st-001`.

### Sequencing and Dependencies

1. SCAN-001 → SCAN-002 → SCAN-003 → SCAN-004
2. SCAN-003 → LIST-001 → LIST-002 → LIST-003
3. SCAN-003 → SUM-001
4. CLI-001 and CLI-002 may proceed in parallel with SCAN work
5. CLI-003 depends on all of the above

### Source Conflicts and Gaps

- None. The sources are internally consistent; `run.sh` agrees with `INSTRUCTIONS.md` on the enumerated outputs and exit codes.

## Analysis Notes
generated: 2026-10-24
blueprint: /mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261024.183232/workspace/targets/DupFind/blueprint

Quality: Ready
  blockers: 0
  questions: 0
  features: 4
  stories: 11
  stack: Go 1.22 with standard library only
  display_name: not proposed
  short_description: not proposed

None.
