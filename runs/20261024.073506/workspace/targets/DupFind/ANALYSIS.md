# Blueprint Analysis: DupFind

## Commander Expectations

- assert the product finds and reports duplicate files under a directory.
- assert the product reports reclaimable space without modifying files.
- assert the supplied grader passes.

## Crew

| Crew | Charge |
|---|---|
| Commander | Defines intent and decides what done means. |
| Team Lead | Confirms epic completeness and stakeholder expectations. |
| Planning Crew | Authors atomic specifications and the ordered Manifest. |
| Shipyard Crew | Builds the tickets without synchronous Commander access. |

## Story List

### Feature: Command Interface

| ID | Story | High-level AC |
|---|---|---|
| COMMAND-001 | Parse command options and validate arguments | Accept the documented options before exactly one directory; reject malformed or unknown arguments with the documented usage exit code. |
| COMMAND-002 | Report directory and usage errors | Missing or non-directory paths produce the documented operational exit code; failures write no standard output and one diagnostic line to standard error. |

### Feature: Duplicate Discovery and Reporting

| ID | Story | High-level AC |
|---|---|---|
| DISCOVERY-001 | Find duplicate regular files recursively | Group non-empty eligible regular files by size and SHA-256 digest; include hidden paths, skip unreadable files with a warning, and neither follow nor report symbolic links. |
| DISCOVERY-002 | List duplicate groups as CSV | Emit the documented header and rows, with groups ordered by reclaimable bytes descending then digest ascending, and paths ordered by byte value ascending. |
| DISCOVERY-003 | Summarize duplicate groups | With `--summary`, emit group count, duplicate file count, and reclaimable bytes for the selected files. |

### Feature: Build and Verification

| ID | Story | High-level AC |
|---|---|---|
| VERIFY-001 | Provide the project test command | Provide POSIX-compatible `bin/test.sh` that runs `go vet ./...` and `go test ./...`, returning success only when both pass. |

## Story Guidance

Named stories planning must produce. Authoritative record: `STORY_GUIDANCE.json`.

None.

## Surfaced Acceptance Criteria

None.

## Source Inventory

| Path | Content kind | Disposition | Reason |
|---|---|---|---|
| `sources/INSTRUCTIONS.md` | markdown | analyzed | readable UTF-8 |
| `sources/run.sh` | code | analyzed | readable UTF-8 |

## Relationship Model

| Source or group | Relationship type | Related source or group | Evidence | Delivery implication |
|---|---|---|---|---|
| `sources/INSTRUCTIONS.md` | instruction-to-test | `sources/run.sh` | Instructions identify the harness as the scoring entry point; the harness builds `./cmd/dupfind` and checks command output and exit codes. | Implement the documented CLI behavior and keep the supplied harness available for final verification. |
| `sources/INSTRUCTIONS.md` | dependency | `sources/run.sh` | Instructions require the Go module root and executable at `./cmd/dupfind`; harness invokes `go build` there. | Build and command implementation precede harness verification. |
| `sources/INSTRUCTIONS.md` | dependency | `bin/test.sh` | Instructions require the test command to run `go vet ./...` and `go test ./...`. | Provide the test command after the Go module and application exist. |

## Source Roles

| Path | Role | Plan disposition | Build disposition |
|---|---|---|---|
| `sources/INSTRUCTIONS.md` | author intent | prompt-only | none |
| `sources/run.sh` | conformance harness | context | stage |

## Planning Instructions

### Delivery Shape

A Go command-line application accepts options and a directory, recursively examines eligible files, groups duplicates by size and SHA-256 digest, and writes either a CSV listing or CSV summary. The flow is argument validation, directory validation, file discovery and hashing, grouping and ordering, then output. The supplied shell harness builds the command and verifies its observable behavior.

### Story Realization Map

| Story ID | Durable Blueprint scope | Evidence | Related files | Requires |
|---|---|---|---|---|
| COMMAND-001 | CLI option and argument parsing | `sources/INSTRUCTIONS.md` — Command and Exit codes | `cmd/dupfind` | capability |
| COMMAND-002 | Operational and usage error behavior | `sources/INSTRUCTIONS.md` — Exit codes and failing command behavior | `cmd/dupfind` | capability |
| DISCOVERY-001 | Recursive file scanning and duplicate grouping | `sources/INSTRUCTIONS.md` — Rules 1–3 and 6 | `cmd/dupfind` | capability |
| DISCOVERY-002 | Default CSV listing and deterministic ordering | `sources/INSTRUCTIONS.md` — Rules 4–5 and Output | `cmd/dupfind` | capability |
| DISCOVERY-003 | Summary CSV output | `sources/INSTRUCTIONS.md` — Summary | `cmd/dupfind` | capability |
| VERIFY-001 | Project test command and terminal verification | `sources/INSTRUCTIONS.md` — Tests; `sources/run.sh` | `bin/test.sh`, `sources/run.sh` | test harness, acceptance contract |

### Test and Acceptance Strategy

Keep focused tests within each story's behavior. The terminal verification story runs the complete supplied grader with `sh sources/run.sh`; the grader builds the executable and checks the specified command cases. The project test command runs `go vet ./...` and `go test ./...`.

### Sequencing and Dependencies

The module and command entry point must exist before either verification command can run. Argument parsing and directory validation precede scanning. Scanning and grouping precede listing and summary output. The complete supplied grader is terminal verification after implementation.

### Source Conflicts and Gaps

None.

## Analysis Notes
generated: 2026-10-24
blueprint: /mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261024.073506/workspace/targets/DupFind/blueprint

Quality: Ready
  blockers: 0
  questions: 0
  features: 3
  stories: 6
  stack: Go 1.22 or newer; standard library
  display_name: not proposed
  short_description: not proposed

None.
