# MANIFEST: DupFind
updated:     2026-10-24T14:07:18+00:00
state:       approved
blocks:      2

applied_specs: |
  ARCHITECTURE.md sha256=27895ac5b7eedfdb61f0173fdf005a797cd8f496bdce6af4d28dff4066c667a2 commit=bf0482e8365898a69bf4cc67c934ea66c40f567f applied_by=block-2 applied_at=2026-10-24T14:13:10+00:00 build_sha256=0740e0ade7d580855c41949b4c3d6a5550f1c10cd584509d6a147a496c4d39de
  FEATURE-Build.md sha256=ca517c3e3d49f8648fb94d53d7bd1be8f6d8bf1cbd7207e1d426192fcab46cd5 commit=bf0482e8365898a69bf4cc67c934ea66c40f567f applied_by=block-2 applied_at=2026-10-24T14:13:10+00:00 build_sha256=980b7c6aaf0190da0eb5f2fe6c735080f32961143330f7aa522a476d93897d30
  FEATURE-Command.md sha256=6913d12ffa50e16c76e990e059520a9a1de8f7d4826e6eb35999591026f56fca commit=bf0482e8365898a69bf4cc67c934ea66c40f567f applied_by=block-2 applied_at=2026-10-24T14:13:10+00:00 build_sha256=a20edcda8c152ea4ab0db5bf817839739311b98aaa996d373ea9dd3331ef02bd
  FEATURE-Listing.md sha256=00ec23f399fbc248203586bc42ec9f8a714e9fcda6e3b221b6680157b8850202 commit=bf0482e8365898a69bf4cc67c934ea66c40f567f applied_by=block-2 applied_at=2026-10-24T14:13:10+00:00 build_sha256=6c420a3c71e992be1ff838e63c1815645b4bafe38fa80a1ee7a16c812ff24462
  FEATURE-Ordering.md sha256=a8a77eb8fc4793037fdee75f4077785e07b2cb705d5e782100a9881dbf6b49c1 commit=bf0482e8365898a69bf4cc67c934ea66c40f567f applied_by=block-2 applied_at=2026-10-24T14:13:10+00:00 build_sha256=905f9b88478e8ed2b8d3a18065544762163e49c449f0ee7170bf55dc83ebb475
  FEATURE-Scanning.md sha256=856647526b5dba10af0b423ebf5efe1f2865988f0c05334004d7a9ff28b049ec commit=bf0482e8365898a69bf4cc67c934ea66c40f567f applied_by=block-2 applied_at=2026-10-24T14:13:10+00:00 build_sha256=f96db7fa6b9be70a78c0b90ee569c7629a829e371945b4325c73037738d1cae6
  FEATURE-Summary.md sha256=2d6e27a9298bac8137fb4bc4c3ae80061b4a24f70100ff79f9f865c6e8e3cb1e commit=bf0482e8365898a69bf4cc67c934ea66c40f567f applied_by=block-2 applied_at=2026-10-24T14:13:10+00:00 build_sha256=94d287e568941714d35f2936f7b2aae9325a7e6bf2f07c93bb566e90f96423d9
  FEATURE-Tests.md sha256=18fd086294424f4420cffb68977075ae6f002a63e4cb2d0e71e510c1fedd3f87 commit=bf0482e8365898a69bf4cc67c934ea66c40f567f applied_by=block-2 applied_at=2026-10-24T14:13:10+00:00 build_sha256=76df475c7a93ea094d0572bb55f75563eb77da41a4f94f9a7a68ed001f8372c9
  VERIFICATION.md sha256=5c2c08c97adaef587f9a218a3777357113395c5dbe9932b9d9dc9aac535010da commit=bf0482e8365898a69bf4cc67c934ea66c40f567f applied_by=block-2 applied_at=2026-10-24T14:13:10+00:00 build_sha256=6998845ff8b3be14904f7b6ec1efcff8b837d2c1d667600b5015ebe0544923eb
  run.sh sha256=dc265eae016629662090b0465362c0642ea7deb57afd77c5725652d0969f880d commit=bf0482e8365898a69bf4cc67c934ea66c40f567f applied_by=block-2 applied_at=2026-10-24T14:13:10+00:00 build_sha256=dc265eae016629662090b0465362c0642ea7deb57afd77c5725652d0969f880d
## story 1: Establish the Go module and application boundaries.
id:           architecture
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-1.md
summary:      Establish the Go module and application boundaries.
type:         foundational
kind:         capability
phase:        1
block:        1
implements:   ARCHITECTURE.md
instructions: |
  Establish the Go module and package boundaries required by the command-line application. Keep the implementation in Go 1.22 or newer and use only the standard library. Define clear ownership for command startup, CLI behavior, scanning, and output formatting. Do not read or modify files under the directory being examined.
stack:        go.md
stack_mode:   builder
size:         6536
provides:     Go module and application package layout
acceptance:   yes
state: closed/implemented

## story 2: Provide the dupfind executable and build contract.
id: build-executable
finding: ADVISORY: implemented, not verified — every binding criterion this story passed was already green before the story built, so none measured the work. Repair the criteria in the Blueprint to exercise this story and rebuild to obtain a verdict.
evidence: evidence/block-2.md
summary: Provide the dupfind executable and build contract.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Build.md
covers: DEL-001
accepts: st-003
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 3341
provides: ./cmd/dupfind
instructions: Provide the main package at ./cmd/dupfind and ensure the application root builds it with go build ./cmd/dupfind using only the standard library.
acceptance: yes
depends: architecture
state: closed/implemented

## story 3: Parse command options and validate the directory argument.
id: command-options
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Parse command options and validate the directory argument.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Command.md
covers: CLI-001
accepts: st-007
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 3554
provides: dupfind command-line options and usage
consumes: ./cmd/dupfind
instructions: Implement the documented command syntax, --summary, --min-size, help text, and validation of options and the single directory argument. Malformed arguments return exit code 2 and write no standard output. Missing or non-directory paths return exit code 1 and write no standard output. Operational diagnostics go to standard error.
acceptance: yes
depends: build-executable
state: closed/implemented

## story 4: Find exact duplicate regular files in a directory tree.
id: scan-and-group
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Find exact duplicate regular files in a directory tree.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Scanning.md
covers: CLI-002
accepts: st-004, st-005, st-008
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 3615
provides: duplicate groups from a directory tree
consumes: dupfind command-line options and usage
instructions: Recursively examine regular files, including hidden paths, without following or reporting symbolic links. Ignore zero-byte files and files smaller than the selected minimum size. Group by size before reading contents, then by SHA-256 digest. Skip unreadable files with one warning on standard error and do not change the exit code. Only read examined files; never delete, move, or write anything under the examined directory.
acceptance: yes
depends: command-options
state: closed/implemented

## story 5: Order duplicate groups and paths deterministically.
id: order-results
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Order duplicate groups and paths deterministically.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Ordering.md
covers: CLI-003
accepts: st-004
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 3278
provides: ordered duplicate groups and relative paths
consumes: duplicate groups from a directory tree
instructions: Order groups by reclaimable bytes descending and then digest ascending. Number groups from 1 in that order. Order paths within each group by byte value and emit paths relative to the requested directory using / separators without a leading ./.
acceptance: yes
depends: scan-and-group
state: closed/implemented

## story 6: Write the duplicate listing as CSV.
id: listing-csv
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Write the duplicate listing as CSV.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Listing.md
covers: OUT-001
accepts: st-004
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 3569
provides: default CSV listing
consumes: ordered duplicate groups and relative paths
instructions: Write the default listing to standard output with encoding/csv. Emit the specified header and one row per file in each duplicate group, with lowercase hexadecimal SHA-256 digests. With no duplicates, emit only the header. Use newline line endings and RFC 4180 field quoting.
acceptance: yes
depends: order-results
state: closed/implemented

## story 7: Write duplicate group and reclaimable-space totals.
id: summary-csv
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Write duplicate group and reclaimable-space totals.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Summary.md
covers: OUT-002
accepts: st-006
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 3594
provides: --summary CSV output
consumes: ordered duplicate groups and relative paths
instructions: When --summary is selected, emit the specified CSV header and one data row containing the number of duplicate groups, the sum of files minus one per group, and total reclaimable bytes. With no duplicates, emit zero for all totals. Derive the summary from the same duplicate groups as the listing.
acceptance: yes
depends: order-results
state: closed/implemented

## story 8: Provide the project quality test entry point.
id: test-harness
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Provide the project quality test entry point.
type: service
kind: test harness
phase: 1
block: 2
implements: FEATURE-Tests.md
covers: DEL-002
accepts: st-002
context: ARCHITECTURE_compact.md
stack: go.md, common.md
stack_mode: builder
size: 4993
provides: bin/test.sh
instructions: Provide a POSIX-compatible bin/test.sh that runs go vet ./... and go test ./... from the application root. Exit zero only when both commands pass.
acceptance: yes
depends: build-executable
state: closed/implemented

## story 9: Pass the supplied duplicate-finder conformance grader.
id: conformance-grader
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Pass the supplied duplicate-finder conformance grader.
type: service
kind: test harness
phase: 1
block: 2
implements: VERIFICATION.md
covers: DEL-003
accepts: st-001
context: run.sh
stack: go.md, common.md
stack_mode: consumer
size: 5364
instructions: Run the supplied sources/run.sh from the completed application root and require its exit status to be zero. Do not modify, rewrite, regenerate, or substitute the supplied harness.
acceptance: yes
depends: architecture, build-executable, command-options, scan-and-group, order-results, listing-csv, summary-csv, test-harness
state: closed/implemented
