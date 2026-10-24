# Scorecard: DupFind

## Specification

> **NOT YET RUN**
> drydock score spec DupFind

## Acceptance

> **PASSED**
> drydock score ac DupFind · 2026-10-24T14:14:15.065+00:00

PASS 13   FAIL 0   UNVERIFIED 0   PREPASSED 4   (17 AC total)

## Build

> **RECORDED**
> drydock score build DupFind

Blocks: 2   Repaired: 0   Failed: 0
Tokens: 645,490 (cached 575,360, output 8,102)   Cache hit rate: 90.3%

## Release

> **RELEASED**
> drydock score release DupFind · 2026-10-24T14:15:00.108+00:00

- Verdict: PASSED

DupFind: PASSED — 8 of 8

  st-001  MET       ["sh", "sources/run.sh"] exited 0; governed gate: sh sources/run.sh exited 0; 15 passed, 0 failed
  st-002  MET       bin/test.sh:1-6; sh bin/test.sh exited 0
  st-003  MET       go.mod declares module dupfind; go build -o dupfind ./cmd/dupfind succeeded; cmd/dupfind/main.go and its imported application packages use standard-library dependencies
  st-004  MET       probe: listing contained both identical files with their correct digest and relative paths; internal/output/output.go: Ordered sorts paths and groups by the required keys
  st-005  MET       probe: listing included only a and sub/b; empty files and links were absent; internal/scan/scan.go: skips ModeSymlink entries and requires regular files
  st-006  MET       probe: summary output was 1 group, 1 duplicate file, 5 reclaimable bytes; internal/output/output.go: Summary computes totals from ordered groups
  st-007  MET       probe: missing directory exited 1 with empty stdout; internal/cli/cli.go: malformed parsing, nonpositive size, and wrong argument count return 2 after writing errors to stderr
  st-008  MET       internal/scan/scan.go: uses os.Open and io.Copy to read examined files; internal/cli/cli.go, internal/scan/scan.go, and internal/output/output.go: no filesystem mutation calls
  reported: Build directory has uncommitted changes

## Project acceptance

| ID | Type | Criterion | Verdict | Observed |
|---|---|---|---|---|
| st-001 | technical | The completed dupfind shall make sh sources/run.sh exit zero; that script's exit status is the sole acceptance verdict. | MET | ["sh", "sources/run.sh"] exited 0; governed gate: sh sources/run.sh exited 0; 15 passed, 0 failed |
| st-002 | technical | The application shall provide a POSIX-compatible bin/test.sh that exits zero when go vet and every automated test pass. | MET | bin/test.sh:1-6; sh bin/test.sh exited 0 |
| st-003 | technical | The application root shall be a Go module from which go build ./cmd/dupfind produces the dupfind executable using only the standard library. | MET | go.mod declares module dupfind; go build -o dupfind ./cmd/dupfind succeeded; cmd/dupfind/main.go and its imported application packages use standard-library dependencies |
| st-004 | behavioral | When dupfind lists a directory, the application shall report every group of two or more regular files with identical size and SHA-256 digest, ordered by reclaimable bytes descending then digest ascending, with paths relative to the directory in byte order. | MET | probe: listing contained both identical files with their correct digest and relative paths; internal/output/output.go: Ordered sorts paths and groups by the required keys |
| st-005 | behavioral | The application shall ignore zero-byte files and shall neither follow nor report symbolic links. | MET | probe: listing included only a and sub/b; empty files and links were absent; internal/scan/scan.go: skips ModeSymlink entries and requires regular files |
| st-006 | behavioral | When --summary is given, the application shall write the group count, the duplicate file count, and the reclaimable bytes of the listing that the same options would produce. | MET | probe: summary output was 1 group, 1 duplicate file, 5 reclaimable bytes; internal/output/output.go: Summary computes totals from ordered groups |
| st-007 | behavioral | If the directory is missing or is not a directory, then the application shall exit 1; if the arguments are malformed, then the application shall exit 2; in both cases it shall write nothing to standard output. | MET | probe: missing directory exited 1 with empty stdout; internal/cli/cli.go: malformed parsing, nonpositive size, and wrong argument count return 2 after writing errors to stderr |
| st-008 | technical | The application shall only read the files it examines and shall never delete, move, or write any file under the directory. | MET | internal/scan/scan.go: uses os.Open and io.Copy to read examined files; internal/cli/cli.go, internal/scan/scan.go, and internal/output/output.go: no filesystem mutation calls |

**Advisory warnings:**
- Build directory has uncommitted changes

**Material decisions:**

- Operational failure handling — 
