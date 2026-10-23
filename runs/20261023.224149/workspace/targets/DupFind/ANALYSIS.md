# Blueprint Analysis: DupFind

## Commander Expectations

- assert the duplicate-file finder builds successfully and passes the supplied acceptance harness.
- assert duplicate files are grouped accurately, reported deterministically, and never modified.

## Crew

| Crew | Charge |
|---|---|
| Commander | Defines intent and decides what done means. |
| Team Lead | Confirms epic completeness and stakeholder expectations. |
| Planning Crew | Authors atomic specifications and the ordered Manifest. |
| Shipyard Crew | Builds the tickets without synchronous Commander access. |

## Story List

### Feature: Command-Line Contract

| ID | Story | High-level AC |
|---|---|---|
| CLI-001 | Parse duplicate-finder options and arguments | Accepts the documented options and exactly one directory; malformed input exits 2 with no stdout. |
| CLI-002 | Validate the scan directory | Missing paths and non-directories exit 1 with no stdout and a diagnostic on stderr. |

### Feature: Duplicate Discovery

| ID | Story | High-level AC |
|---|---|---|
| DISCOVERY-001 | Walk eligible files recursively | Recursively includes hidden entries, considers regular files, and neither follows nor reports symbolic links. |
| DISCOVERY-002 | Filter and hash candidate files | Applies the positive minimum-size filter, ignores empty files, groups candidates by size and SHA-256 digest, and skips unreadable files with one stderr warning. |

### Feature: Duplicate Reporting

| ID | Story | High-level AC |
|---|---|---|
| REPORT-001 | Emit ordered duplicate listings | Writes the required CSV header and rows, with RFC 4180 escaping, relative slash-separated paths, deterministic group ordering, and byte-ordered paths. |
| REPORT-002 | Emit duplicate summaries | Writes the required summary header and totals for groups, duplicate files, and reclaimable bytes, including zero totals. |

### Feature: Verification and Delivery

| ID | Story | High-level AC |
|---|---|---|
| VERIFY-001 | Provide the project test gate | `bin/test.sh` runs `go vet ./...` and `go test ./...`, succeeding only when both pass. |
| VERIFY-002 | Pass the complete supplied conformance suite | The staged grader runs unfiltered and every supplied case passes with zero failures. Suite: full |

## Story Guidance

Named stories planning must produce. Authoritative record: `STORY_GUIDANCE.json`.

None.

## Surfaced Acceptance Criteria

| ID | Story ID | Criterion |
|---|---|---|
| AC-001 | CLI-001 | The command shall provide help text for its command and documented options when invoked through the conventional Go CLI help path. |
| AC-002 | REPORT-001 | The listing shall use `\n` line endings and quote any CSV field containing a comma, double quote, or newline according to RFC 4180. |
| AC-003 | VERIFY-001 | The test script shall be POSIX-compatible and execute from the application root. |

## Source Inventory

| Path | Content kind | Disposition | Reason |
|---|---|---|---|
| `sources/INSTRUCTIONS.md` | markdown | analyzed | readable UTF-8 |
| `sources/run.sh` | code | analyzed | readable UTF-8 |

## Relationship Model

| Source or group | Relationship type | Related source or group | Evidence | Delivery implication |
|---|---|---|---|---|
| `sources/INSTRUCTIONS.md` | instruction-to-test | `sources/run.sh` | The instructions define the command, rules, outputs, and exit codes exercised by the harness. | Implement the normative CLI and behavioral contract before conformance verification. |
| `sources/run.sh` | conformance harness | `cmd/dupfind`, `bin/test.sh` | The harness builds `./cmd/dupfind` and invokes the resulting executable across listing, summary, filtering, and error cases. | Stage and preserve the harness; terminal verification depends on all implementation stories. |
| `sources/run.sh` | fixture-to-verification | `cmd/dupfind` | The harness creates hidden files, empty files, a symlink, comma-containing paths, and duplicate contents. | Discovery and reporting must be tested against these filesystem edge cases. |
| `sources/INSTRUCTIONS.md` | dependency | Go module and standard library | The product requires Go 1.22 or newer, a module root, and standard library only. | Build configuration precedes implementation and test stories. |

## Source Roles

| Path | Role | Plan disposition | Build disposition |
|---|---|---|---|
| `sources/INSTRUCTIONS.md` | normative specification and conformance test suite | context | prompt-only |
| `sources/run.sh` | conformance harness | context | stage |

## Planning Instructions

### Delivery Shape

A Go command-line application scans one supplied directory recursively. Filesystem entries are the input; duplicate groups, CSV listings, or CSV summaries are the outputs. The execution flow is: parse and validate arguments, validate the directory, walk eligible files, filter and hash candidates, form ordered duplicate groups, then emit either listing or summary output. Failures produce diagnostics on stderr and no stdout.

### Story Realization Map

| Story ID | Blueprint scope | Evidence | Related files | Delivery type |
|---|---|---|---|---|
| CLI-001 | Argument parsing and usage errors | `sources/INSTRUCTIONS.md` command and exit-code tables | `cmd/dupfind` | capability |
| CLI-002 | Directory validation and operational errors | `sources/INSTRUCTIONS.md` exit-code and failure rules | `cmd/dupfind` | capability |
| DISCOVERY-001 | Recursive filesystem traversal | `sources/INSTRUCTIONS.md` scope rules | `cmd/dupfind` | capability |
| DISCOVERY-002 | Size filtering, hashing, unreadable-file handling | `sources/INSTRUCTIONS.md` duplicate and unreadable-file rules | `cmd/dupfind` | capability |
| REPORT-001 | Listing serialization and ordering | `sources/INSTRUCTIONS.md` listing and ordering rules; `sources/run.sh` expected CSV | `cmd/dupfind` | capability |
| REPORT-002 | Summary serialization and totals | `sources/INSTRUCTIONS.md` summary rules; `sources/run.sh` summary cases | `cmd/dupfind` | capability |
| VERIFY-001 | Local vet/test command | `sources/INSTRUCTIONS.md` test requirement | `bin/test.sh` | test harness |
| VERIFY-002 | Full supplied conformance execution | `sources/run.sh` complete 15-case grader | `sources/run.sh` | acceptance contract |

### Test and Acceptance Strategy

Implementation stories use focused Go tests for their own parsing, traversal, grouping, ordering, serialization, and exit behavior. `VERIFY-001` runs the project-local vet and test gate. `VERIFY-002` is the sole terminal story for the complete supplied suite and declares `Suite: full`; it depends on every implementation story. The final Sea Trial invokes the staged harness unfiltered.

### Sequencing and Dependencies

Create the Go module and executable entry point before implementation stories. Implement argument and directory validation before scanning. Implement traversal and candidate hashing before reporting. Implement listing and summary output after duplicate groups are available. Add `bin/test.sh` after the test suite exists. Run the complete conformance harness only after all implementation and local verification stories are complete.

### Source Conflicts and Gaps

No unresolved source conflicts or blocking gaps. The sources consistently define a standard-library Go CLI with no persistence, external services, authentication, or deployment target.

## Analysis Notes
generated: 2026-10-23
blueprint: /mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261023.224149/workspace/targets/DupFind/blueprint

Quality: Ready
  blockers: 0
  questions: 0
  features: 4
  stories: 8
  stack: Go 1.22+ standard library; POSIX sh verification
  display_name: not proposed
  short_description: not proposed

Project type: cli. The supplied harness is staged as the authoritative conformance asset. No Compass rewrite is required because the existing Compass is authoritative and already formatted.
