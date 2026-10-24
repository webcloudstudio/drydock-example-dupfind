# Blueprint Analysis: DupFind

## Commander Expectations

- assert the product finds and reports duplicate regular files and the space reclaimable by keeping one copy.
- assert the product never deletes, moves, or modifies files.

## Crew

| Crew | Charge |
|---|---|
| Commander | Defines intent and decides what done means. |
| Team Lead | Confirms epic completeness and stakeholder expectations. |
| Planning Crew | Authors atomic specifications and the ordered Manifest. |
| Shipyard Crew | Builds the tickets without synchronous Commander access. |

## Story List

### Feature: Duplicate Discovery

| ID | Story | High-level AC |
|---|---|---|
| DISCOVERY-001 | Scan a directory for duplicate files | Recursively inspect regular files, including hidden paths, skip symbolic links and unreadable files with a warning, and group eligible files by size and SHA-256 digest. |
| DISCOVERY-002 | List duplicate groups as CSV | Emit the required CSV header and rows, ordered by reclaimable bytes descending, digest ascending, and paths by byte order; quote CSV fields as required. |
| DISCOVERY-003 | Filter files by minimum size | Accept a positive `--min-size` value and exclude smaller files, defaulting to 1 byte. |

### Feature: Summary and Command Operation

| ID | Story | High-level AC |
|---|---|---|
| COMMAND-001 | Summarize duplicate results | With `--summary`, emit group count, duplicate file count, and reclaimable bytes for the matching listing. |
| COMMAND-002 | Handle command errors and usage | Return the specified exit codes and a diagnostic on standard error for operational and usage failures; leave standard output empty on failure. |

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
| `sources/INSTRUCTIONS.md` | instruction-to-test | `sources/run.sh` | Instructions identify the staged grader and specify its command, fixed tree, and expected outputs; runner builds `./cmd/dupfind`. | Build the CLI behavior to the source requirements; preserve and run the supplied harness as the terminal conformance check. |
| `sources/run.sh` | test-kit-to-implementation | `./cmd/dupfind` | Harness builds and invokes the executable for listing, summary, filtering, and error cases. | The executable must expose all required command behavior through its CLI. |

## Source Roles

| Path | Role | Plan disposition | Build disposition |
|---|---|---|---|
| `sources/INSTRUCTIONS.md` | normative specification and conformance test suite | context | prompt-only |
| `sources/run.sh` | conformance harness | context | stage |

## Planning Instructions

### Delivery Shape

A Go command-line tool accepts options followed by one directory, recursively discovers eligible files, groups duplicates, and emits either a CSV listing or CSV summary. The flow validates arguments and directory, scans candidates, hashes same-size files, forms and orders duplicate groups, then writes the selected output.

### Story Realization Map

| Story ID | Blueprint scope | Evidence | Related files | Requires |
|---|---|---|---|---|
| DISCOVERY-001 | Recursive file discovery and duplicate grouping | `sources/INSTRUCTIONS.md` Rules 1–3, 6 | `./cmd/dupfind` | capability |
| DISCOVERY-002 | Listing output format and ordering | `sources/INSTRUCTIONS.md` Rules 4–5 and Output | `./cmd/dupfind` | capability |
| DISCOVERY-003 | Minimum-size option | `sources/INSTRUCTIONS.md` Command and Options | `./cmd/dupfind` | capability |
| COMMAND-001 | Summary output | `sources/INSTRUCTIONS.md` Command and Output | `./cmd/dupfind` | capability |
| COMMAND-002 | Argument validation and operational errors | `sources/INSTRUCTIONS.md` Exit codes and failing command behavior | `./cmd/dupfind` | capability |

### Test and Acceptance Strategy

Keep focused tests bounded to each story: temporary directory cases for scanning/grouping, CSV and ordering, minimum-size filtering, summary totals, and argument and directory failures. The supplied `sources/run.sh` is the terminal verification gate and must pass in full; no story-level test substitutes for that conformance run.

### Sequencing and Dependencies

The Go module and executable entry point underpin all stories. Argument parsing and directory validation precede scanning; scanning and grouping precede listing or summary formatting. The final conformance verification follows implementation of all stories. No migration or external service dependency is described.

### Source Conflicts and Gaps

The source consistently specifies a Go CLI and its behavior. It does not define an interactive help screen or a dedicated help flag; the core command contract otherwise gives the required options, arguments, and errors.

## Analysis Notes
generated: 2026-10-24
blueprint: /mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261024.080843/workspace/targets/DupFind/blueprint

Quality: Ready
  blockers: 0
  questions: 0
  features: 2
  stories: 5
  stack: Go 1.22+, standard library
  display_name: not proposed
  short_description: not proposed

None.
