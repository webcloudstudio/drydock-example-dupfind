# MANIFEST: DupFind
updated:     2026-10-25T17:45:27+00:00
state:       approved
blocks:      3

applied_specs: |
  ARCHITECTURE.md sha256=08109b711fcdc50a0973bc1ea1ec716067e63f13e3ec8640a0d28e527aa4617b commit=41117e37e175023f9bd7b8694784d43a6da9457b applied_by=block-2 applied_at=2026-10-25T17:55:52+00:00 build_sha256=aebff5c19477750419aa3a1dbc2fc1873cbd2154b060d3bf1d7acf29e855eb21
  FEATURE-ArgParsing.md sha256=1287bc6702f199267145a49cbf3bdfec639f942620fb06ea1eda87721ab4ea1d commit=41117e37e175023f9bd7b8694784d43a6da9457b applied_by=block-2 applied_at=2026-10-25T17:55:52+00:00 build_sha256=7e4a3fc649d48e1a2ac8cc92cf22d439bcd90a1f72ebe8c3b9e2202861a15d8f
  FEATURE-DigestGrouping.md sha256=f9f1967c6b58c2d855d9bd565ed963f0bb25d4f5e7486631dd9d4d9da696ea06 commit=41117e37e175023f9bd7b8694784d43a6da9457b applied_by=block-2 applied_at=2026-10-25T17:55:52+00:00 build_sha256=f9033d1789dd1b67fb3fbb8e7c242338affefbf337425bdc6866e8e32b3dafe5
  FEATURE-DirectoryValidation.md sha256=f0b618ddb077eb6eff0967afb82d6e7564872457ceb695d1757720b98fa583db commit=41117e37e175023f9bd7b8694784d43a6da9457b applied_by=block-2 applied_at=2026-10-25T17:55:52+00:00 build_sha256=7222048a89a3b7d4e06f1328a553dacd0d9c4dc511f717e082469063ffcb2332
  FEATURE-GroupOrdering.md sha256=d3e44700191eb2cdd2fb17a7af6a1ef888846eecee995504f4180e55949363d2 commit=41117e37e175023f9bd7b8694784d43a6da9457b applied_by=block-2 applied_at=2026-10-25T17:55:52+00:00 build_sha256=661174a5cab0164211ed3cd00bf1d8e6e351e43fb36f5d5fd420a4ab7f11a764
  FEATURE-Listing.md sha256=d836e942346be7f52b8393f9d733cd18ff6c94b1d504b411bd3e18600f6da79f commit=41117e37e175023f9bd7b8694784d43a6da9457b applied_by=block-2 applied_at=2026-10-25T17:55:52+00:00 build_sha256=4e44e032422c57afaae98a796091831ab3d743911392ffc6ceddf60640b5dd5c
  FEATURE-ModuleBuild.md sha256=5df0b0c22ea6fc3ba5447c1d86dad20891fdf5d8c1e2102b6d5faa6a93e01027 commit=41117e37e175023f9bd7b8694784d43a6da9457b applied_by=block-2 applied_at=2026-10-25T17:55:52+00:00 build_sha256=1dbf51042f9de4db40ac08857d33f5c961244c4071cf2012e88c84ae57a0185b
  FEATURE-SizeGrouping.md sha256=7b05d752ec997083b83a89aee7a6084342bf8ccea3c22317490519814102b397 commit=41117e37e175023f9bd7b8694784d43a6da9457b applied_by=block-2 applied_at=2026-10-25T17:55:52+00:00 build_sha256=c7f746ae3925092b3478e8018eafc7bba4d290d8abdeb08b3b5d1abc0d6fefb1
  FEATURE-Summary.md sha256=48b34e18f55c4f242fdf2c695fd5721dbf9c09dca5369b709fe6150466083f65 commit=41117e37e175023f9bd7b8694784d43a6da9457b applied_by=block-2 applied_at=2026-10-25T17:55:52+00:00 build_sha256=55808c5b807acd0e472857a20b1020a9a6da1f081e69e873e26994a268eed4a8
  FEATURE-TestGate.md sha256=abde1e383926f4299dbffb9294a54c6c9ed330ca64a1760dd9e8d0767aa72a5b commit=41117e37e175023f9bd7b8694784d43a6da9457b applied_by=block-2 applied_at=2026-10-25T17:55:52+00:00 build_sha256=78ca67350897e92d61390e3c46d238cd12df910a50b1ae91561a0360322dc026
  FEATURE-Walk.md sha256=802bb8da2686a42bfb81938d57111dd30e80839d8439832fb8fdb7edb078cda1 commit=41117e37e175023f9bd7b8694784d43a6da9457b applied_by=block-2 applied_at=2026-10-25T17:55:52+00:00 build_sha256=c314359704c381857b9900d8e2a7228d488f7c605158d42a886b2ae86b8f2326
  VERIFICATION.md sha256=daf5a6f1e4c99d434659cb378e1713643337ba132a5104775719d09061d54e68 commit=41117e37e175023f9bd7b8694784d43a6da9457b applied_by=block-3 applied_at=2026-10-25T17:58:00+00:00 build_sha256=ff5ef19f1fc2019ee6fb5c8364e0dc0650d70901a15d4cc8a8bf7c6d6dad4efb
## story 1: Stand up the Go module, package tree, and test launcher that every other story builds into.
id:           foundation
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-1.md
summary:      Stand up the Go module, package tree, and test launcher that every other story builds into.
type:         foundational
kind:         capability
phase:        1
block:        1
implements:   ARCHITECTURE.md
instructions: |
  Initialize the root go.mod (Go 1.22 or newer, standard library only, no third-party modules).
  Create the package tree: cmd/dupfind (package main), internal/cli, internal/walk,
  internal/dupes, internal/report, and bin/.

  Document in ARCHITECTURE the pipeline: parse arguments, validate directory, walk, group by
  size, hash same-size groups, order and number groups, emit CSV. Expose
  Run(args []string, stdout, stderr io.Writer) int in internal/cli; only main calls os.Exit.

  Add a minimal compiling placeholder in each internal package so go vet ./... and go test ./...
  succeed before feature stories land. Do not read sources/run.sh from application code.
stack:        common.md, go.md
stack_mode:   builder
size:         13397
provides:     go.mod, internal/, cmd/dupfind, bin/
acceptance:   yes
state: closed/implemented

## story 2: Parse --summary, --min-size, and exactly one directory argument; reject malformed usage with exit 2.
id: arg-parsing
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Parse --summary, --min-size, and exactly one directory argument; reject malformed usage with exit 2.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-ArgParsing.md
covers: CLI-001, CLI-002
accepts: st-007
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 6658
provides: cli.Run, cli.Config
instructions: |
  Implement Run in internal/cli using a local flag.FlagSet with flag.ContinueOnError and its
  output discarded. Options precede the directory; exactly one directory argument is required.
  
  Missing or extra arguments, an unknown option, a help flag, or a --min-size value that is not a
  positive decimal integer (no sign, no zero, no non-decimal forms) return 2 with empty stdout
  and one diagnostic line on stderr. The default min-size is 1.
  
  Produce a parsed configuration (summary flag, min size, directory) for later stories.
acceptance: yes
depends: foundation
state: closed/implemented

## story 3: Provide the cmd/dupfind main package that builds with go build and delegates to cli.Run.
id: module-build
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Provide the cmd/dupfind main package that builds with go build and delegates to cli.Run.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-ModuleBuild.md
covers: MODBUILD-001
accepts: st-003
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 6119
provides: cmd/dupfind
consumes: cli.Run
instructions: |
  Write cmd/dupfind/main.go so main calls cli.Run(os.Args[1:], os.Stdout, os.Stderr) and passes
  the result to os.Exit. go build ./cmd/dupfind must succeed using only the standard library.
acceptance: yes
depends: arg-parsing
state: closed/implemented

## story 4: Exit 1 with empty stdout when the directory is missing or is not a directory.
id: directory-validation
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Exit 1 with empty stdout when the directory is missing or is not a directory.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-DirectoryValidation.md
covers: CLI-003
accepts: st-007
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 6371
provides: cli.validateDirectory
consumes: cli.Config
instructions: |
  After successful argument parsing, stat the directory. A nonexistent path or a path that is not
  a directory returns exit 1 with nothing on stdout and one diagnostic line on stderr.
  Validation happens before any walking or output.
acceptance: yes
depends: arg-parsing
state: closed/implemented

## story 5: Walk the directory recursively and yield eligible regular files at or above the minimum size.
id: walk
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Walk the directory recursively and yield eligible regular files at or above the minimum size.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Walk.md
covers: WALK-001, WALK-002, WALK-003
accepts: st-005
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 7061
provides: walk.Files
consumes: cli.Config
instructions: |
  Implement walk.Files in internal/walk, returning each eligible file's relative path (slash
  separated, no leading ./) and size.
  
  Traverse hidden files and directories. Consider regular files only: symbolic links are neither
  followed nor reported, and directories reached only through a symbolic link are not traversed.
  Ignore files smaller than min-size; zero-byte files are never returned. Use lstat semantics
  (os.Lstat or DirEntry type) and never open file contents here.
  
  A subdirectory that cannot be read yields one stderr warning and the walk continues.
acceptance: yes
depends: directory-validation
state: closed/implemented

## story 6: Group candidate files by exact size and keep only sizes shared by two or more files.
id: size-grouping
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Group candidate files by exact size and keep only sizes shared by two or more files.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-SizeGrouping.md
covers: DUP-001
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 6419
provides: dupes.GroupBySize
consumes: walk.Files
instructions: |
  Implement dupes.GroupBySize, returning only size buckets with at least two files. A file whose
  size is unique is never opened for hashing. The size is taken from the walk result, not from
  reading contents.
acceptance: yes
depends: walk
state: closed/implemented

## story 7: Hash same-size files with SHA-256 and form duplicate groups, skipping unreadable files with a warning.
id: digest-grouping
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Hash same-size files with SHA-256 and form duplicate groups, skipping unreadable files with a warning.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-DigestGrouping.md
covers: DUP-002, DUP-006
accepts: st-008
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 7737
provides: dupes.GroupByDigest
consumes: dupes.GroupBySize
instructions: |
  For each same-size bucket, open each file read-only and stream it through crypto/sha256. The
  digest is lowercase hexadecimal of the full contents. Files with identical size and digest
  form one group of two or more files.
  
  A file that cannot be opened or read yields exactly one warning line on stderr, is excluded from
  grouping, and does not change the exit code. Never delete, move, or write any file under the
  examined directory.
acceptance: yes
depends: size-grouping
state: closed/implemented

## story 8: Compute reclaimable bytes, order groups, number them from 1, and order paths within each group.
id: group-ordering
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Compute reclaimable bytes, order groups, number them from 1, and order paths within each group.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-GroupOrdering.md
covers: DUP-003, DUP-004, DUP-005
accepts: st-004
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 7136
provides: dupes.Order
consumes: dupes.GroupByDigest
instructions: |
  Reclaimable bytes per group equal size x (files - 1). Order groups by reclaimable bytes
  descending, then digest ascending, and number them from 1. Within a group order paths by byte
  value ascending (plain string comparison, not locale aware). Paths are relative to the
  directory with / separators and no leading ./.
acceptance: yes
depends: digest-grouping
state: closed/implemented

## story 9: Write the CSV duplicate listing with header group,size,sha256,path using encoding/csv.
id: listing
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Write the CSV duplicate listing with header group,size,sha256,path using encoding/csv.
type: service
kind: capability
phase: 1
block: 2
implements: FEATURE-Listing.md
covers: LIST-001, LIST-002
context: ARCHITECTURE_compact.md
stack: go.md
stack_mode: consumer
size: 7587
provides: report.WriteListing
consumes: dupes.Order
instructions: |
  Implement report.WriteListing with encoding/csv: header group,size,sha256,path, then one row
  per file in every ordered group. With no duplicates write only the header. Line endings are \n
  (Writer.UseCRLF false). Fields containing a comma, double quote, or newline are quoted by the
  csv writer.
  
  Wire the default (non-summary) mode of cli.Run through walk, grouping, ordering, and this
  writer, exit 0.
acceptance: yes
depends: group-ordering, module-build
state: closed/implemented

## story 10: Write the one-row summary of groups, duplicate files, and reclaimable bytes with --summary.
id: summary
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Write the one-row summary of groups, duplicate files, and reclaimable bytes with --summary.
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
size: 7107
provides: report.WriteSummary
consumes: dupes.Order
instructions: |
  Implement report.WriteSummary with encoding/csv: header groups,duplicate_files,reclaimable_bytes
  then one row. duplicate_files is the sum of (files - 1) over groups; reclaimable_bytes is the
  sum of every group's reclaimable bytes. With no duplicates the row is 0,0,0.
  
  Wire --summary mode of cli.Run to this writer using the same walk and grouping as the listing
  with the same --min-size, exit 0.
acceptance: yes
depends: listing
state: closed/implemented

## story 11: Provide POSIX bin/test.sh and a Go test suite covering every rule, output form, and exit code.
id: test-gate
finding: ADVISORY: implemented, not verified — no governed acceptance command covers this story. Declare one in ACCEPTANCE.json to gate it.
evidence: evidence/block-2.md
summary: Provide POSIX bin/test.sh and a Go test suite covering every rule, output form, and exit code.
type: service
kind: test harness
phase: 1
block: 2
implements: FEATURE-TestGate.md
covers: MODBUILD-002
accepts: st-002
context: ARCHITECTURE_compact.md
stack: common.md, go.md
stack_mode: consumer
size: 10013
provides: bin/test.sh
consumes: cli.Run
instructions: |
  Write bin/test.sh as POSIX sh: cd to the application root, run go vet ./... then go test
  ./..., and exit non-zero if either fails.
  
  Ensure the Go test suite, using real temporary directories and Run with bytes.Buffer, covers
  each rule (hidden files, symlinks, zero-byte, min-size, size-then-digest, ordering, path
  form, unreadable skip), each output form (listing, summary, empty cases), and each exit code
  (0, 1, 2). Expectations in tests are written against observed real output.
acceptance: yes
depends: module-build, directory-validation, summary, listing
state: closed/implemented

## story 12: Assemble the finished tool and confirm the supplied grader passes end to end.
id:           verification
finding: ADVISORY: implemented, not verified — every binding criterion this story passed was already green before the story built, so none measured the work. Repair the criteria in the Blueprint to exercise this story and rebuild to obtain a verdict.
evidence: evidence/block-3.md
summary:      Assemble the finished tool and confirm the supplied grader passes end to end.
type:         feature
kind:         test harness
phase:        1
block:        3
implements:   VERIFICATION.md
accepts:      st-001
instructions: |
  Terminal assembly story. Confirm that sh bin/test.sh exits 0 and that sh sources/run.sh exits 0
  with every check passing. Do not modify, regenerate, or substitute sources/run.sh, and do not
  read it from application code. Fix any defect in the owning package, not in the grader.
stack:        common.md
stack_mode:   consumer
size:         5989
consumes:     bin/test.sh, cmd/dupfind
depends:      test-gate
acceptance:   yes
state: closed/implemented
