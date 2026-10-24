# MANIFEST: DupFind
updated: 2026-10-24T07:02:15+00:00
state: draft
blocks:      2

## story 1: Establish the Go module, command entry point, and project test runner.
id:           language-foundation
summary:      Establish the Go module, command entry point, and project test runner.
type:         foundational
kind:         capability
phase:        1
block:        1
implements:   ARCHITECTURE.md
instructions: |
  Establish the Go 1.22+ standard-library-only module and the ./cmd/dupfind entry point. Provide bin/test.sh to run go vet ./... and go test ./... from the application root. Define clear boundaries between CLI argument handling, directory discovery and hashing, and CSV output so the feature can be implemented and tested.
stack:        go.md
stack_mode:   builder
size:         6349
provides:     ./cmd/dupfind, bin/test.sh
acceptance:   yes
state:        pending

## story 2: Discover duplicate regular files and emit ordered CSV listings or summaries.
id: duplicate-discovery
summary: Discover duplicate regular files and emit ordered CSV listings or summaries.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Duplicate-Finding.md
covers: DUPLICATE-DISCOVERY-001, DUPLICATE-DISCOVERY-002, DUPLICATE-DISCOVERY-003, SUMMARY-001, COMMAND-ERRORS-001
accepts: st-004, st-005, st-006, st-007, st-008
context: INSTRUCTIONS.md, ARCHITECTURE_compact.md
stack_mode: builder
size: 5450
provides: dupfind command
consumes: ./cmd/dupfind
instructions: Implement the command contract in sources/INSTRUCTIONS.md. Parse options before the one required directory argument; distinguish usage errors (exit 2) from missing or non-directory paths (exit 1), with failures producing no stdout. Recursively walk hidden directories, skip symbolic links and zero-byte files, group by size before reading contents, hash same-size files with SHA-256, and warn once to stderr for each unreadable file while continuing. Apply the minimum-size filter consistently to listing and summary. Produce RFC-compliant CSV with relative slash-separated paths, bytewise path ordering, and groups ordered by reclaimable bytes descending then digest ascending. The summary reports counts and reclaimable bytes of those selected groups. Never modify examined files. Add focused Go tests for behavior and error paths.
acceptance: yes
depends: language-foundation
state: pending

## story 3: Verify the completed command against the supplied full conformance grader.
id: supplied-grader
summary: Verify the completed command against the supplied full conformance grader.
type: service
kind: test harness
phase: 1
block: 2
implements: VERIFICATION.md
covers: TEST-ACCEPTANCE-001, TEST-ACCEPTANCE-002
accepts: st-001
context: run.sh, INSTRUCTIONS.md
stack: common.md, go.md
stack_mode: builder
size: 7127
instructions: Run the supplied grader unfiltered from the application root after the command implementation is complete. Do not alter or replace the harness. Its exit status is the release conformance verdict; repair the implementation until the complete suite passes.
acceptance: yes
depends: duplicate-discovery
state: pending
