# MANIFEST: DupFind
updated: 2026-10-23T22:44:22+00:00
state: draft
blocks:      3

## story 1: Establish the Go module structure and application boundaries for DupFind.
id:           architecture
summary:      Establish the Go module structure and application boundaries for DupFind.
type:         foundational
kind:         capability
phase:        1
block:        1
implements:   ARCHITECTURE.md
scope:        target
instructions: |
  Create the Go module and the standard-library application structure required for DupFind.
  Establish ./cmd/dupfind as the executable entry point and keep filesystem scanning, grouping,
  reporting, and CLI orchestration in maintainable packages.
stack:        go.md
stack_mode:   builder
size:         6345
provides:     Go module, cmd/dupfind entry point
acceptance:   yes
state:        pending

## story 2: Parse DupFind options and enforce the command argument contract.
id: cli-options
summary: Parse DupFind options and enforce the command argument contract.
type: service
kind: capability
phase: 2
block: 2
implements: FEATURE-CLI-OPTIONS.md
covers: CLI-001
accepts: st-007
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 3578
provides: dupfind option parsing
consumes: Go module and cmd/dupfind entry point
instructions: |
  Implement the documented command-line contract: --summary, --min-size with a positive integer,
  and exactly one directory argument. Malformed input and unknown options must produce no stdout
  and exit 2. Use the conventional Go help path and keep only the process entry point responsible
  for os.Exit.
acceptance: yes
depends: architecture
state: pending
scope: target

## story 3: Validate that the requested scan path exists and is a directory.
id: directory-validation
summary: Validate that the requested scan path exists and is a directory.
type: service
kind: capability
phase: 2
block: 2
implements: FEATURE-DIRECTORY-VALIDATION.md
covers: CLI-002
accepts: st-007
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 3390
provides: scan-directory validation
consumes: dupfind option parsing
instructions: |
  Validate the single scan argument before traversal. Missing paths and non-directory paths must
  produce no stdout, emit a diagnostic on stderr, and exit 1. Successful validation must pass the
  directory to the discovery pipeline without modifying it.
acceptance: yes
depends: cli-options
state: pending
scope: target

## story 4: Walk the scan tree and select eligible regular files.
id: file-discovery
summary: Walk the scan tree and select eligible regular files.
type: service
kind: capability
phase: 2
block: 2
implements: FEATURE-FILE-DISCOVERY.md
covers: DISCOVERY-001
accepts: st-005, st-008
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 3686
provides: recursive eligible-file discovery
consumes: scan-directory validation
instructions: |
  Walk the supplied directory recursively, including hidden entries. Include regular files,
  exclude symbolic links, and do not follow symbolic links. Ignore empty files and never delete,
  move, or write anything under the examined directory.
acceptance: yes
depends: directory-validation
state: pending
scope: target

## story 5: Filter candidates and group files by size and SHA-256 digest.
id: file-hashing
summary: Filter candidates and group files by size and SHA-256 digest.
type: service
kind: capability
phase: 2
block: 2
implements: FEATURE-FILE-HASHING.md
covers: DISCOVERY-002
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 3838
provides: duplicate candidate grouping
consumes: recursive eligible-file discovery
instructions: |
  Apply the positive minimum-size filter, group candidates by size before reading contents, and
  compute SHA-256 only for files sharing a size. Skip unreadable files with one stderr warning and
  continue successfully. Produce exact duplicate groups based on size and digest.
acceptance: yes
depends: file-discovery
state: pending
scope: target

## story 6: Emit the duplicate-group summary and reclaimable-space totals.
id: summary-report
summary: Emit the duplicate-group summary and reclaimable-space totals.
type: service
kind: capability
phase: 2
block: 2
implements: FEATURE-SUMMARY-REPORT.md
covers: REPORT-002
accepts: st-006
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 3699
provides: duplicate summary CSV
consumes: duplicate candidate grouping
instructions: |
  Implement --summary using the same duplicate groups as listing mode. Emit the summary header
  and exactly one row containing group count, duplicate-file count, and reclaimable bytes,
  including zero totals when no duplicates exist.
acceptance: yes
depends: file-hashing
state: pending
scope: target

## story 7: Emit deterministic CSV rows for duplicate groups.
id: listing-report
summary: Emit deterministic CSV rows for duplicate groups.
type: service
kind: capability
phase: 2
block: 2
implements: FEATURE-LISTING-REPORT.md
covers: REPORT-001
accepts: st-004
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 3739
provides: duplicate listing CSV
consumes: duplicate candidate grouping
instructions: |
  Emit the listing header and one CSV row per file in each duplicate group using encoding/csv.
  Order groups by reclaimable bytes descending and digest ascending, number them from one, and
  order relative slash-separated paths by byte value. Preserve RFC 4180 escaping and LF endings.
acceptance: yes
depends: file-hashing
state: pending
scope: target

## story 8: Provide the local POSIX test and vet gate.
id: test-gate
summary: Provide the local POSIX test and vet gate.
type: service
kind: test harness
phase: 2
block: 2
implements: FEATURE-TEST-GATE.md
covers: VERIFY-001
accepts: st-002
context: ARCHITECTURE_compact.md
stack: common.md, go.md
stack_mode: builder
size: 5071
provides: bin/test.sh
consumes: Go module and application test suite
instructions: |
  Add a POSIX-compatible bin/test.sh that runs go vet ./... and go test ./... from the application
  root, preserving failure status and exiting zero only when both commands pass. Add focused Go
  tests for the implemented CLI, discovery, grouping, ordering, and serialization behavior.
acceptance: yes
depends: architecture
state: pending
scope: target

## story 9: Run the complete supplied DupFind conformance suite.
id: conformance
summary: Run the complete supplied DupFind conformance suite.
type: service
kind: test harness
phase: 3
block: 3
implements: FEATURE-CONFORMANCE.md
covers: VERIFY-002
accepts: st-001, st-003
context: run.sh, ARCHITECTURE_compact.md
stack: common.md, go.md
stack_mode: consumer
size: 5814
provides: complete supplied conformance verification
consumes: bin/test.sh, dupfind executable
instructions: |
  Preserve the supplied sources/run.sh as the authoritative terminal conformance harness.
  Build and verify the completed ./cmd/dupfind implementation against the full unfiltered suite.
  Do not modify or replace the staged harness, and do not run the full suite before this terminal
  story.
acceptance: yes
depends: cli-options, directory-validation, file-discovery, file-hashing, listing-report, summary-report, test-gate
state: pending
scope: target
