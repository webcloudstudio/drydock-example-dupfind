# MANIFEST: DupFind
updated: 2026-10-24T08:11:52+00:00
state: draft
blocks:      2

## story 1: Establish the Go application module, executable scaffold, and test runner.
id:           go-foundation
summary:      Establish the Go application module, executable scaffold, and test runner.
type:         foundational
kind:         capability
phase:        1
block:        1
implements:   ARCHITECTURE.md
instructions: |
  Establish the Go 1.22+ standard-library-only module and ./cmd/dupfind entry point. Add a POSIX-compatible bin/test.sh that runs go vet ./... and go test ./... from the application root and exits successfully only when both pass. Keep application behavior in internal packages and keep only startup and os.Exit handling in main.
stack:        go.md, common.md
stack_mode:   builder
size:         8294
provides:     Go module, ./cmd/dupfind, bin/test.sh
acceptance:   yes
state:        pending

## story 2: Discover and group duplicate regular files by size and SHA-256.
id: duplicate-discovery
summary: Discover and group duplicate regular files by size and SHA-256.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Discovery.md
covers: DISCOVERY-001
accepts: st-004, st-005, st-008
context: ARCHITECTURE_compact.md
stack: go.md, common.md
stack_mode: consumer
size: 5018
provides: duplicate file groups
instructions: Recursively inspect the requested directory, including hidden paths. Consider regular files only, do not follow or report symbolic links, ignore zero-byte files, and skip unreadable files with one warning on standard error without changing the exit status. Read contents only for files sharing a size, then group by size and SHA-256 digest. Do not modify files under the examined directory.
acceptance: yes
depends: go-foundation
state: pending

## story 3: Filter scan candidates using the minimum file size option.
id: minimum-size
summary: Filter scan candidates using the minimum file size option.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-MinimumSize.md
covers: DISCOVERY-003
context: ARCHITECTURE_compact.md
stack: go.md, common.md
stack_mode: consumer
size: 4270
provides: --min-size
consumes: duplicate file groups
instructions: Accept --min-size <bytes> before the directory argument. Require a positive integer, default to 1 byte, and exclude smaller files from duplicate discovery. Apply the same filtered groups to listing and summary output.
acceptance: yes
depends: duplicate-discovery
state: pending

## story 4: Render ordered duplicate groups as CSV.
id: duplicate-listing
summary: Render ordered duplicate groups as CSV.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Listing.md
covers: DISCOVERY-002
context: ARCHITECTURE_compact.md
stack: go.md, common.md
stack_mode: consumer
size: 4405
provides: CSV duplicate listing
consumes: duplicate file groups
instructions: Write the listing as CSV using encoding/csv with the required header and one row per file in each duplicate group. Use paths relative to the scanned directory with forward-slash separators and no leading ./; order paths by byte value. Order groups by reclaimable bytes descending, then digest ascending, and number them from 1. Emit only the header when there are no duplicates.
acceptance: yes
depends: duplicate-discovery
state: pending

## story 5: Report duplicate counts and reclaimable bytes as CSV.
id: duplicate-summary
summary: Report duplicate counts and reclaimable bytes as CSV.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Summary.md
covers: COMMAND-001
accepts: st-006
context: ARCHITECTURE_compact.md
stack: go.md, common.md
stack_mode: consumer
size: 4343
provides: --summary CSV
consumes: duplicate file groups, --min-size
instructions: When --summary is supplied, write the required CSV header and exactly one data row. Compute groups, duplicate_files as the sum of files minus one per group, and reclaimable_bytes as the sum of size times files minus one. With no duplicate groups, report zero for every total.
acceptance: yes
depends: duplicate-listing, minimum-size
state: pending

## story 6: Validate CLI arguments and directory and complete release conformance.
id: command-operation
summary: Validate CLI arguments and directory and complete release conformance.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Command.md
covers: COMMAND-002
accepts: st-001, st-007
context: run.sh, ARCHITECTURE_compact.md
stack: go.md, common.md
stack_mode: consumer
size: 6506
instructions: Validate that exactly one directory is supplied and options precede it. Return exit 2 for missing or extra arguments, unknown options, or a non-positive or invalid --min-size value. Return exit 1 when the directory does not exist or is not a directory. On failure, leave standard output empty and write one diagnostic line to standard error. Successful commands return 0. As the final implementation step, run the complete supplied grader via sh sources/run.sh and make it pass.
acceptance: yes
depends: duplicate-listing, minimum-size, duplicate-summary
state: pending
