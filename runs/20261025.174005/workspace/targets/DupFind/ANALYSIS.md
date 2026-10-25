# Blueprint Analysis: drydock DupFind Uat Case

## Commander Expectations

- assert the product is a command-line duplicate file finder that groups regular files by size and SHA-256 digest
- assert the product never deletes, moves, or modifies any file it examines
- assert the product reports either a full duplicate listing or a summary of reclaimable bytes
- assert the product is implemented in Go 1.22 or newer using only the standard library
- assert the supplied grader `sources/run.sh` is the sole acceptance verdict

## Crew

| Crew | Charge |
|---|---|
| Commander | Defines intent and decides what done means. |
| Team Lead | Confirms epic completeness and stakeholder expectations. |
| Planning Crew | Authors atomic specifications and the ordered Manifest. |
| Shipyard Crew | Builds the tickets without synchronous Commander access. |

## Story List

### Feature: Module and Build

| ID | Story | High-level AC |
|---|---|---|
| MODBUILD-001 | Establish Go module at application root with main package at ./cmd/dupfind | `go build ./cmd/dupfind` succeeds using only the standard library |
| MODBUILD-002 | Provide POSIX bin/test.sh running go vet and go test | `bin/test.sh` exits zero only when `go vet ./...` and `go test ./...` both pass |

### Feature: Directory Walk and Eligibility

| ID | Story | High-level AC |
|---|---|---|
| WALK-001 | Walk the directory recursively including hidden files and directories | Every regular file under the directory, including names beginning `.`, is considered |
| WALK-002 | Exclude symbolic links from consideration | Symbolic links are neither followed nor reported |
| WALK-003 | Exclude files smaller than the minimum size | Files with size below `--min-size` (default 1) are ignored; zero-byte files are never reported |

### Feature: Duplicate Detection

| ID | Story | High-level AC |
|---|---|---|
| DUP-001 | Group candidate files by exact size | Only files sharing a size are read for hashing |
| DUP-002 | Group same-size files by SHA-256 digest | Files with identical size and digest form one group of two or more files |
| DUP-003 | Compute reclaimable bytes per group | Reclaimable bytes equal `size × (files − 1)` |
| DUP-004 | Order groups and assign group numbers | Groups ordered by reclaimable bytes descending then digest ascending; numbered from 1 |
| DUP-005 | Order paths within a group | Paths relative to directory, `/` separator, no leading `./`, byte-ascending order |
| DUP-006 | Skip unreadable files with a stderr warning | Unreadable file produces one stderr warning line and does not change exit code |

### Feature: Command Line Interface

| ID | Story | High-level AC |
|---|---|---|
| CLI-001 | Parse `--summary`, `--min-size <bytes>`, and one directory argument | Options precede the directory; exactly one directory is accepted |
| CLI-002 | Reject malformed arguments with exit code 2 | Missing or extra arguments, unknown option, or non-positive-integer `--min-size` exit 2 with empty stdout |
| CLI-003 | Reject missing or non-directory target with exit code 1 | Nonexistent path or non-directory exits 1 with empty stdout |

### Feature: Listing Output

| ID | Story | High-level AC |
|---|---|---|
| LIST-001 | Write CSV listing with header `group,size,sha256,path` | One row per file in every duplicate group; header only when no duplicates |
| LIST-002 | Quote CSV fields per RFC 4180 | Fields containing comma, double quote, or newline are quoted via `encoding/csv` |

### Feature: Summary Output

| ID | Story | High-level AC |
|---|---|---|
| SUM-001 | Write summary row with header `groups,duplicate_files,reclaimable_bytes` | Single row totals groups, duplicate files, and reclaimable bytes; `0,0,0` when none |

## Story Guidance

Named stories planning must produce. Authoritative record: `STORY_GUIDANCE.json`.

None.

## Surfaced Acceptance Criteria

| ID | Story ID | Criterion |
|---|---|---|
| AC-001 | WALK-001 | Hidden files and directories are traversed and their regular files considered |
| AC-002 | WALK-002 | A symbolic link pointing at a regular file is not reported and its target is not followed |
| AC-003 | DUP-001 | Only files sharing a size are opened for hashing |
| AC-004 | DUP-002 | Digest is lowercase hexadecimal SHA-256 of full file contents |
| AC-005 | LIST-001 | Listing rows match header order `group,size,sha256,path` |
| AC-006 | SUM-001 | Summary `duplicate_files` equals sum of `files − 1` across groups |
| AC-007 | CLI-002 | Usage errors write nothing to standard output and one diagnostic line to standard error |
| AC-008 | CLI-003 | Operational errors write nothing to standard output and one diagnostic line to standard error |
| AC-009 | DUP-006 | Unreadable file warning is written to standard error only |
| AC-010 | MODBUILD-002 | Test suite covers each rule, each output form, and each exit code against real temporary directories |

## Source Inventory

| Path | Content kind | Disposition | Reason |
|---|---|---|---|
| `sources/INSTRUCTIONS.md` | markdown | analyzed | readable UTF-8 |
| `sources/run.sh` | code | analyzed | readable UTF-8 |

## Relationship Model

| Source or group | Relationship type | Related source or group | Evidence | Delivery implication |
|---|---|---|---|---|
| sources/run.sh | instruction-to-test | sources/INSTRUCTIONS.md | run.sh cites the 15 checks, the `./cmd/dupfind` build target, and the exit-code table defined in INSTRUCTIONS.md | Implementation must satisfy INSTRUCTIONS.md; run.sh is the external acceptance verdict and is staged, not authored |
| sources/INSTRUCTIONS.md | dependency | Command, Output, Exit codes sections | Rules 1–6 define the duplicate test that drives all listing and summary output | Walk, hash, group, and output stories depend on the Rules section |

## Source Roles

| Path | Role | Plan disposition | Build disposition |
|---|---|---|---|
| sources/INSTRUCTIONS.md | author intent | compass | prompt-only |
| sources/run.sh | conformance harness | context | stage |

## Planning Instructions

### Delivery Shape

The system is a single-binary Go CLI. Input: one directory path plus optional `--summary` and `--min-size` flags. Output: CSV on standard output (listing or summary) with exit codes 0, 1, or 2. Flow: parse arguments → validate directory → walk recursively → filter regular non-empty files at or above minimum size → group by size → hash same-size groups with SHA-256 → form duplicate groups → order and number → emit CSV.

### Story Realization Map

- MODBUILD-001: `sources/INSTRUCTIONS.md` Stack section; `sources/run.sh` build step. Requires a Go module and `cmd/dupfind` main package.
- MODBUILD-002: `sources/INSTRUCTIONS.md` Tests section. Requires `bin/test.sh`.
- WALK-001 through WALK-003: `sources/INSTRUCTIONS.md` Scope and Rule 2. Requires a directory walker.
- DUP-001 through DUP-006: `sources/INSTRUCTIONS.md` Rules 3–6. Requires hashing and grouping logic.
- CLI-001 through CLI-003: `sources/INSTRUCTIONS.md` Command and Exit codes. Requires argument parsing.
- LIST-001, LIST-002: `sources/INSTRUCTIONS.md` Output (listing). Requires `encoding/csv`.
- SUM-001: `sources/INSTRUCTIONS.md` Output (summary).

### Test and Acceptance Strategy

Story-scoped tests run per feature area with bounded scope. The terminal verification story depends on every implementation story and runs the full `bin/test.sh` suite (`Suite: full`). Sea Trial `st-001` executes `sh sources/run.sh` as the sole acceptance verdict.

### Sequencing and Dependencies

Build (MODBUILD-001) precedes all feature stories. CLI-001 precedes CLI-002 and CLI-003. Walk stories precede Duplicate Detection. DUP-001 precedes DUP-002. DUP-002 precedes DUP-003 and DUP-004. Listing and Summary output stories depend on DUP-004 and DUP-005. MODBUILD-002 is the terminal gate.

### Source Conflicts and Gaps

- None. The single source defines all behavior; no cross-source conflicts exist. Unresolved edge cases (symlink to a directory, file read errors mid-hash) are resolved by the narrowest reading: directories are not regular files and are not walked into via symlinks; read errors follow Rule 6.

## Analysis Notes
generated: 2026-10-25
blueprint: /mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261025.174005/workspace/targets/DupFind/blueprint

Quality: Ready
  blockers: 0
  questions: 0
  features: 6
  stories: 16
  stack: Go 1.22+ (standard library only)
  display_name: not proposed
  short_description: not proposed

Selected behaviors (unattended policy):
- Symbolic links are excluded at the walk level; directories reached only through symbolic links are not traversed.
- Files whose size is zero are excluded by the default `--min-size 1` and by the explicit zero-byte rule.
- Unreadable files are skipped with one stderr line; they are excluded from grouping.
- `--min-size` is validated as a positive decimal integer; leading `+`, signs, and non-decimal forms are usage errors.
- Path separators are normalized to `/` on all platforms.
