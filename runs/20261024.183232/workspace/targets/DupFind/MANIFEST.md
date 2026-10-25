# MANIFEST: DupFind
updated: 2026-10-24T18:42:13+00:00
state: draft
blocks:      3

## story 1: Establish the Go module, the cmd and internal layout, and the bin launcher directory.
id:           foundation
summary:      Establish the Go module, the cmd and internal layout, and the bin launcher directory.
type:         foundational
kind:         capability
phase:        1
block:        1
implements:   ARCHITECTURE.md
instructions: |
  Create go.mod at the application root declaring Go 1.22 and a module path, with no third-party requires.
  Create cmd/dupfind/ as the main package directory, internal/ for application packages, and bin/ for launchers.
  Use the standard library only.
stack:        go.md, common.md
stack_mode:   builder
size:         13093
provides:     go.mod module, cmd/dupfind main package directory, internal/ package directory
acceptance:   yes
state:        pending

## story 2: Walk the target directory recursively and yield regular files only.
id: scan-walk
summary: Walk the target directory recursively and yield regular files only.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Walker.md
covers: SCAN-001
accepts: st-005
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 5563
provides: scan.Walk
instructions: |
  Implement the recursive walker in internal/scan.
  Descend into hidden directories and include hidden files such as .env.
  Yield regular files only. Neither follow nor report symbolic links, even when a link targets a duplicate file.
  Skip directories and special files.
acceptance: yes
depends: foundation
state: pending

## story 3: Drop zero-byte files and files smaller than the minimum size.
id: scan-size
summary: Drop zero-byte files and files smaller than the minimum size.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-SizeFilter.md
covers: SCAN-002
accepts: st-005
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 5403
provides: scan.FilterSize
consumes: scan.Walk
instructions: |
  Implement FilterSize in internal/scan as a pure function over the walker output.
  Exclude every file of size 0. Exclude files smaller than the configured minimum, which defaults to 1 byte.
acceptance: yes
depends: scan-walk
state: pending

## story 4: Group files by size, then by SHA-256 digest, reading only size-matched files.
id: scan-group
summary: Group files by size, then by SHA-256 digest, reading only size-matched files.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Grouping.md
covers: SCAN-003
accepts: st-004, st-008
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 5520
provides: scan.Group
consumes: scan.FilterSize
instructions: |
  Implement Group in internal/scan.
  Bucket files by size. Read a file's contents only when another file shares its size.
  Within a size bucket, group by SHA-256 digest. A group of two or more files with equal size and digest is a duplicate group.
  Open files read-only. Never write, move, rename, or delete any file.
acceptance: yes
depends: scan-size
state: pending

## story 5: Skip unreadable files with one stderr warning line each, without changing the exit code.
id: scan-unreadable
summary: Skip unreadable files with one stderr warning line each, without changing the exit code.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-UnreadableFiles.md
covers: SCAN-004
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 5663
provides: scan.Warn
consumes: scan.Group
instructions: |
  When a file cannot be opened or hashed, exclude it from grouping.
  Write exactly one warning line naming the file to standard error, and write no stdout row for it.
  The warning does not alter the exit code.
acceptance: yes
depends: scan-group
state: pending

## story 6: Emit the --summary CSV totals computed from the same grouping as the listing.
id: summary
summary: Emit the --summary CSV totals computed from the same grouping as the listing.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Summary.md
covers: SUM-001
accepts: st-006
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 6152
provides: output.WriteSummary
consumes: scan.Group
instructions: |
  Implement WriteSummary in internal/output.
  Write the header groups,duplicate_files,reclaimable_bytes and one totals row.
  duplicate_files is the sum of files minus one over all groups. reclaimable_bytes is the sum of every group's reclaimable bytes.
  With no duplicates, the row is 0,0,0.
acceptance: yes
depends: scan-group
state: pending

## story 7: Emit the duplicate listing as CSV with the exact header.
id: list-csv
summary: Emit the duplicate listing as CSV with the exact header.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-ListingCSV.md
covers: LIST-001
accepts: st-004
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 6285
provides: output.WriteListing
consumes: scan.Group
instructions: |
  Implement WriteListing in internal/output using encoding/csv.
  Write the header group,size,sha256,path, then one row per file in every duplicate group.
  Write sha256 as lowercase hexadecimal. Use \n line endings.
acceptance: yes
depends: scan-group
state: pending

## story 8: Order groups and paths, number groups from 1, and render relative paths.
id: list-order
summary: Order groups and paths, number groups from 1, and render relative paths.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-ListingOrder.md
covers: LIST-002
accepts: st-004
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 6168
provides: output.Order
consumes: output.WriteListing
instructions: |
  Order groups by reclaimable bytes descending, where reclaimable bytes equal size times (files minus one), then by digest ascending.
  Number groups from 1 in that order.
  Within a group, order paths by byte value ascending.
  Render paths relative to the target directory with / separators and no leading ./.
acceptance: yes
depends: list-csv
state: pending

## story 9: Emit only the header when no duplicate groups exist.
id: list-empty
summary: Emit only the header when no duplicate groups exist.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-ListingEmpty.md
covers: LIST-003
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 5651
provides: output.WriteListing empty case
consumes: output.Order
instructions: When there are no duplicate groups, output is exactly group,size,sha256,path followed by a newline and nothing else.
acceptance: yes
depends: list-order
state: pending

## story 10: Parse the command line with options before the single directory argument.
id: cli-args
summary: Parse the command line with options before the single directory argument.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-CommandLine.md
covers: CLI-001
accepts: st-007
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 6578
provides: cli.Run
instructions: |
  Implement Run(args []string, stdout, stderr io.Writer) int in internal/cli using a local flag.FlagSet with flag.ContinueOnError.
  Accept [--summary] [--min-size <bytes>] <directory>. Exactly one directory is required.
  Unknown options, missing or extra arguments, and a --min-size that is not a positive integer exit 2 with exactly one stderr line and no stdout.
  Run receives an injected execute function that the pipeline story supplies. Only main calls os.Exit.
acceptance: yes
depends: foundation
state: pending

## story 11: Reject a missing or non-directory target with exit 1.
id: cli-targets
summary: Reject a missing or non-directory target with exit 1.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-TargetValidation.md
covers: CLI-002
accepts: st-007
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 6118
provides: cli.ValidateTarget
consumes: cli.Run
instructions: |
  After argument parsing succeeds, check the directory argument.
  A nonexistent path or a regular file exits 1 with exactly one stderr line and empty stdout.
  Validation runs after usage validation, so usage errors keep exit 2.
acceptance: yes
depends: cli-args
state: pending

## story 12: Wire the main package to parsing, scanning, ordering, and rendering as the built executable.
id: pipeline
summary: Wire the main package to parsing, scanning, ordering, and rendering as the built executable.
type: feature
kind: integration
phase: 1
block: 3
implements: FEATURE-Pipeline.md
accepts: st-003, st-004, st-006
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 6286
provides: cmd/dupfind executable
consumes: cli.Run, cli.ValidateTarget, scan.Walk, scan.FilterSize, scan.Group, output.WriteListing, output.WriteSummary
instructions: |
  Implement cmd/dupfind/main.go, which calls cli.Run with an execute function.
  The execute function runs walk, size filter, group, order, and render in sequence, choosing listing or summary output.
  Only main calls os.Exit.
  This is the assembly story. Its acceptance verifies the seams by building the executable and running it through the command-line runner.
acceptance: yes
depends: scan-unreadable, list-empty, summary, cli-targets
state: pending

## story 13: Ship the POSIX bin/test.sh and act as the terminal verification gate.
id: test-harness
summary: Ship the POSIX bin/test.sh and act as the terminal verification gate.
type: feature
kind: test harness
phase: 1
block: 3
implements: FEATURE-TestHarness.md
covers: CLI-003
accepts: st-001, st-002
context: ARCHITECTURE_compact.md
stack: go.md, common.md
stack_mode: consumer
size: 7451
instructions: |
  Create bin/test.sh as POSIX sh. Run go vet ./... and go test ./... from the application root.
  Exit zero only when both pass.
  This is the terminal story. Its acceptance runs the supplied grader sh sources/run.sh as the full suite and requires exit zero.
acceptance: yes
depends: scan-walk, scan-size, scan-group, scan-unreadable, list-csv, list-order, list-empty, summary, cli-args, cli-targets, pipeline
state: pending
