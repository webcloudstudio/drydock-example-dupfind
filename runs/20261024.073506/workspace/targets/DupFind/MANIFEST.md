# MANIFEST: DupFind
updated:     2026-10-24T07:37:18+00:00
state:       approved
blocks:      2

applied_specs: |
  ARCHITECTURE.md sha256=c86333e5266241be5ba8662f0c3c3cdd9767d9f6a6ff07f171b21bacd125984f commit=7d29b09f9ad5862957dfd874aadba104675217d8 applied_by=block-1 applied_at=2026-10-24T07:39:40+00:00 build_sha256=c86333e5266241be5ba8662f0c3c3cdd9767d9f6a6ff07f171b21bacd125984f

## story 1: Establish the Go module and command scaffold.
id:           foundation
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-1.md
summary:      Establish the Go module and command scaffold.
type:         foundational
kind:         capability
phase:        1
block:        1
implements:   ARCHITECTURE.md
instructions: |
  Create the Go 1.22 module and the ./cmd/dupfind main package. Keep application code within the module and use only the standard library.
stack:        go.md, common.md
stack_mode:   builder
size:         8217
provides:     ./cmd/dupfind
acceptance:   yes
state: closed/implemented

## story 2: Parse command options and validate directory arguments.
id: command_interface
finding: step block-2 failed: summary-no-duplicates, summary-filtered, listing-empty, listing-relative-paths, listing-csv-path
evidence: evidence/block-2.md
summary: Parse command options and validate directory arguments.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Command-Interface.md
covers: COMMAND-001, COMMAND-002
accepts: st-007
context: ARCHITECTURE_compact.md
stack: go.md, common.md
stack_mode: consumer
size: 4714
provides: dupfind command interface
consumes: ./cmd/dupfind
instructions: Implement the CLI options, argument validation, directory validation, and documented failure exit codes and output channels.
acceptance: yes
depends: foundation
state: closed/failed

## story 3: Find duplicate regular files recursively.
id: duplicate_discovery
finding: step block-2 failed: summary-no-duplicates, summary-filtered, listing-empty, listing-relative-paths, listing-csv-path
evidence: evidence/block-2.md
summary: Find duplicate regular files recursively.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Duplicate-Discovery.md
covers: DISCOVERY-001
accepts: st-004, st-005, st-008
context: ARCHITECTURE_compact.md
stack: go.md, common.md
stack_mode: consumer
size: 4933
provides: duplicate groups
consumes: dupfind command interface
instructions: Implement recursive file discovery, size-first hashing, digest grouping, and safe read-only handling for files and symbolic links.
acceptance: yes
depends: command_interface
state: closed/failed

## story 4: Emit duplicate summary totals as CSV.
id: summary_output
finding: programmatic acceptance failed: summary-no-duplicates, summary-filtered: PermissionError: [Errno 13] Permission denied: './cmd/dupfind'
evidence: evidence/block-2.md
summary: Emit duplicate summary totals as CSV.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Summary-Output.md
covers: DISCOVERY-003
accepts: st-006
context: ARCHITECTURE_compact.md
stack: go.md, common.md
stack_mode: consumer
size: 4551
provides: CSV duplicate summary
consumes: duplicate groups
instructions: Implement summary output for duplicate group count, duplicate file count, and reclaimable bytes, respecting the selected minimum size.
acceptance: yes
depends: duplicate_discovery
state: closed/failed

## story 5: Emit ordered duplicate groups as CSV.
id: listing_output
finding: programmatic acceptance failed: listing-empty, listing-relative-paths, listing-csv-path: PermissionError: [Errno 13] Permission denied: './cmd/dupfind'
evidence: evidence/block-2.md
summary: Emit ordered duplicate groups as CSV.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Listing-Output.md
covers: DISCOVERY-002
accepts: st-004
context: ARCHITECTURE_compact.md
stack: go.md, common.md
stack_mode: consumer
size: 4766
provides: CSV duplicate listing
consumes: duplicate groups
instructions: Implement deterministic listing output with the required CSV fields, group order, numbering, and path order.
acceptance: yes
depends: duplicate_discovery
state: closed/failed

## story 6: Provide and run the project verification suite.
id: verification
finding: step block-2 failed: summary-no-duplicates, summary-filtered, listing-empty, listing-relative-paths, listing-csv-path
evidence: evidence/block-2.md
summary: Provide and run the project verification suite.
type: service
kind: test harness
phase: 1
block: 2
implements: FEATURE-Verification.md
covers: VERIFY-001
accepts: st-001, st-002, st-003
context: run.sh, ARCHITECTURE_compact.md
stack: go.md, common.md
stack_mode: consumer
size: 6292
instructions: Provide bin/test.sh to run go vet ./... and go test ./... successfully. As the terminal story, run the supplied full grader and require its exit status to be zero.
acceptance: yes
depends: command_interface, duplicate_discovery, listing_output, summary_output
state: closed/failed
