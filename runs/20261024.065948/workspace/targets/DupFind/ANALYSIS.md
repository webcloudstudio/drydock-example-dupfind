# Blueprint Analysis: DupFind

## Commander Expectations

- assert the product finds and reports duplicate files under a directory, or summarizes the space reclaimable by removing duplicate copies.
- assert the product never modifies files it examines.

## Crew

| Crew | Charge |
|---|---|
| Commander | Defines intent and decides what done means. |
| Team Lead | Confirms epic completeness and stakeholder expectations. |
| Planning Crew | Authors atomic specifications and the ordered Manifest. |
| Shipyard Crew | Builds the tickets without synchronous Commander access. |

## Story List

### Feature: Duplicate Discovery and Listing

| ID | Story | High-level AC |
|---|---|---|
| DUPLICATE-DISCOVERY-001 | Find duplicate regular files recursively | Groups nonempty regular files by size and SHA-256 digest, including files under hidden directories; skips unreadable files with one warning per file. |
| DUPLICATE-DISCOVERY-002 | List duplicate groups as CSV | Outputs the required CSV header and rows, with groups ordered by reclaimable bytes descending then digest ascending, and paths relative to the directory in byte order. |
| DUPLICATE-DISCOVERY-003 | Filter files by minimum size | Accepts a positive integer minimum size, defaults to 1 byte, and lists only matching duplicate groups. |

### Feature: Summary and Command Errors

| ID | Story | High-level AC |
|---|---|---|
| SUMMARY-001 | Summarize duplicate groups | With `--summary`, outputs the group count, duplicate file count, and reclaimable bytes for the groups selected by the same options. |
| COMMAND-ERRORS-001 | Handle usage and directory errors | Malformed arguments exit 2; a missing or non-directory path exits 1. Failures write no standard output and produce a diagnostic on standard error. |

### Feature: Test and Acceptance

| ID | Story | High-level AC |
|---|---|---|
| TEST-ACCEPTANCE-001 | Run the Go quality checks | Provides a POSIX-compatible `bin/test.sh` that runs `go vet ./...` and `go test ./...`, exiting zero only when both pass. |
| TEST-ACCEPTANCE-002 | Pass the supplied grader | The completed application passes every supplied grader check. Assertion heading: `Suite: full`. |

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
| `sources/INSTRUCTIONS.md` | instruction-to-test | `sources/run.sh` | The instructions define the command contract; the harness builds `./cmd/dupfind` and checks exit codes and standard output. | Implement the source-defined behavior and use the harness as the final conformance gate. |
| `sources/INSTRUCTIONS.md` | dependency | `./cmd/dupfind` | The command depends on recursive file discovery, size and digest grouping, CSV listing, summary, and error handling. | Plan the discovery and output behavior before the terminal grader verification story. |
| `sources/INSTRUCTIONS.md` | dependency | `bin/test.sh` | The instructions require `go vet ./...` and `go test ./...` from the application root. | Deliver the quality-check script with the application; its test suite covers the command behavior. |

## Source Roles

| Path | Role | Plan disposition | Build disposition |
|---|---|---|---|
| `sources/INSTRUCTIONS.md` | normative specification and conformance test suite | context | prompt-only |
| `sources/run.sh` | conformance harness | context | stage |

Build disposition governs what exists on disk when the build agent runs and when acceptance executes:

- `stage` places the file in the build directory at `sources/{path relative to sources/}`. Every assertion, `ac` check, and Sea Trial `Command:` references it by that build-relative path. A test suite or harness the project must execute is always `stage` — a file present only in the prompt can be read but not run. **Everything a staged file needs at run time is also `stage`**: a harness's imported modules, its normalizer, and its fixtures. Read each staged file's imports and open() calls and stage what they name; a harness missing one dependency cannot run at all.
- `prompt-only` supplies the file as prompt context and places nothing on disk.
- `none` neither stages nor supplies it.

Markdown is never staged; use `prompt-only` for it.

## Planning Instructions

### Delivery Shape

Deliver a Go command-line application that recursively examines one directory and writes either a CSV duplicate listing or a CSV summary. Its input is the requested directory and command options; its output is standard output, with diagnostics on standard error. The execution flow is argument validation, directory validation, recursive discovery of eligible files, size-based filtering and SHA-256 grouping, then ordered CSV output. The application also provides a script for Go quality checks, and the supplied grader is the final conformance gate.

### Story Realization Map

| Story ID | Durable Blueprint scope | Evidence | Related files | Capability, integration, migration, test harness, or acceptance contract |
|---|---|---|---|---|
| DUPLICATE-DISCOVERY-001 | Recursive discovery and duplicate grouping | `sources/INSTRUCTIONS.md` Rules 1–3, 6 | `./cmd/dupfind` | Capability |
| DUPLICATE-DISCOVERY-002 | Listing format, ordering, and relative paths | `sources/INSTRUCTIONS.md` Rules 4–5 and Output | `./cmd/dupfind` | Capability |
| DUPLICATE-DISCOVERY-003 | Minimum-size option | `sources/INSTRUCTIONS.md` Command and Rules 2–3 | `./cmd/dupfind` | Capability |
| SUMMARY-001 | Summary output | `sources/INSTRUCTIONS.md` Command and Output | `./cmd/dupfind` | Capability |
| COMMAND-ERRORS-001 | Usage and operational errors | `sources/INSTRUCTIONS.md` Exit codes and Output | `./cmd/dupfind` | Capability |
| TEST-ACCEPTANCE-001 | Quality-check script and automated tests | `sources/INSTRUCTIONS.md` Tests | `bin/test.sh`, Go tests | Test harness |
| TEST-ACCEPTANCE-002 | Complete supplied grader verification | `sources/run.sh`, `sources/INSTRUCTIONS.md` Run Harness | `sources/run.sh` | Acceptance contract |

### Test and Acceptance Strategy

Keep focused Go tests bounded to each behavior story: discovery and grouping, CSV listing and ordering, minimum-size filtering, summary totals, and argument and directory errors. The quality-check script runs `go vet ./...` and `go test ./...`. The terminal acceptance story runs the supplied grader unfiltered and declares `Suite: full`; it depends on every implementation story. Sea Trial verification runs `sh sources/run.sh` from the application root.

### Sequencing and Dependencies

Implement and build the command before running Go tests or the grader. Discovery and grouping support listing and summary; option validation and directory validation establish the command boundary. The quality-check script follows the testable implementation. The full grader verification story follows every implementation story and is terminal.

### Source Conflicts and Gaps

None.

## Analysis Notes
generated: 2026-10-24
blueprint: /mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261024.065948/workspace/targets/DupFind/blueprint

Quality: Ready
  blockers: 0
  questions: 0
  features: 3
  stories: 7
  stack: Go 1.22 or newer, standard library
  display_name: not proposed
  short_description: not proposed

None.
