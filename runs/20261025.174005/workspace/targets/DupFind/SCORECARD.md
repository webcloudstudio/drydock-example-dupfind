# Scorecard: DupFind

## Specification

> **NOT YET RUN**
> drydock score spec DupFind

## Acceptance

> **PASSED**
> drydock score ac DupFind · 2026-10-25T17:59:45.808+00:00

PASS 34   FAIL 0   UNVERIFIED 0   PREPASSED 11   (45 AC total)

## Build

> **RECORDED**
> drydock score build DupFind

Blocks: 3   Repaired: 0   Failed: 0
Tokens: 1,477,424 (cached 1,421,339, output 56,033)   Cache hit rate: 100.0%

## Release

> **RELEASED**
> drydock score release DupFind · 2026-10-25T18:00:45.362+00:00

- Verdict: PASSED

DupFind: PASSED — 8 of 8

  st-001  MET       ["sh", "sources/run.sh"] exited 0; governed gate 'full' PASS, return_code 0, stdout 'DupFind: 15 passed, 0 failed (of 15)'; sources/run.sh:exit status 0 on every check
  st-002  MET       bin/test.sh:1-24 uses #!/bin/sh, set -eu, no bashisms, and runs `go vet ./... && go test ./...`; ran `sh bin/test.sh`: go vet clean, all four test packages ok, 'completed: exit 0'
  st-003  MET       go.mod:1 'module dupfind'; cmd/dupfind/main.go:1-12 package main calling cli.Run; `go list -deps ./cmd/dupfind` shows only stdlib packages and dupfind/internal/* packages; go build -o .drydock-probe/bin/dupfind ./cmd/dupfind succeeded
  st-004  MET       sources/run.sh A1/A2 expected listing matches output: groups ordered 1(1000B) 2(128B) 3(17B) 4(9B) 5(6B), paths sorted within group; internal/dupes/dupes.go Order(): sort.Strings on Paths; ri>rj then Digest<; internal/dupes/dupes.go GroupByDigest(): SHA-256 via digestFile, only groups of 2+ files kept
  st-005  MET       internal/walk/walk.go: `info.Size() == 0` skipped; `!d.Type().IsRegular()` skips symlinks; filepath.WalkDir does not descend symlinked directories; probe: tree with linkdir -> other/ and linkfile.txt -> real/f1.txt printed only real/f1.txt and other/f2.txt, rc=0; grader C/D/E cases with empty files and report-link.txt pass
  st-006  MET       internal/report/report.go WriteSummary: files += len(Paths)-1, reclaimable += Reclaimable(), groups = len(groups); grader B1 '5,6,1177', C2 '2,2,1128', D2 '0,0,0' all PASS
  st-007  MET       probe: `dupfind <missing>` rc=1 stdout_bytes=0, stderr 'stat ...: no such file or directory'; probe: `dupfind <regular file>` rc=1 stdout_bytes=0, stderr 'is not a directory'; grader U1-U5 compare stdout against empty expected output with rc 2, all PASS; internal/cli/cli.go Run(): usage errors return 2 and write only to stderr
  st-008  MET       internal/dupes/dupes.go digestFile: os.Open (read-only) with io.Copy into the hash; internal/walk/walk.go: filepath.WalkDir and d.Info() only, no writes; probe: stat of every path in a tree before and after running dupfind was identical

## Project acceptance

| ID | Type | Criterion | Verdict | Observed |
|---|---|---|---|---|
| st-001 | technical | The completed dupfind shall make sh sources/run.sh exit zero; that script's exit status is the sole acceptance verdict. | MET | ["sh", "sources/run.sh"] exited 0; governed gate 'full' PASS, return_code 0, stdout 'DupFind: 15 passed, 0 failed (of 15)'; sources/run.sh:exit status 0 on every check |
| st-002 | technical | The application shall provide a POSIX-compatible bin/test.sh that exits zero when go vet and every automated test pass. | MET | bin/test.sh:1-24 uses #!/bin/sh, set -eu, no bashisms, and runs `go vet ./... && go test ./...`; ran `sh bin/test.sh`: go vet clean, all four test packages ok, 'completed: exit 0' |
| st-003 | technical | The application root shall be a Go module from which go build ./cmd/dupfind produces the dupfind executable using only the standard library. | MET | go.mod:1 'module dupfind'; cmd/dupfind/main.go:1-12 package main calling cli.Run; `go list -deps ./cmd/dupfind` shows only stdlib packages and dupfind/internal/* packages; go build -o .drydock-probe/bin/dupfind ./cmd/dupfind succeeded |
| st-004 | behavioral | When dupfind lists a directory, the application shall report every group of two or more regular files with identical size and SHA-256 digest, ordered by reclaimable bytes descending then digest ascending, with paths relative to the directory in byte order. | MET | sources/run.sh A1/A2 expected listing matches output: groups ordered 1(1000B) 2(128B) 3(17B) 4(9B) 5(6B), paths sorted within group; internal/dupes/dupes.go Order(): sort.Strings on Paths; ri>rj then Digest<; internal/dupes/dupes.go GroupByDigest(): SHA-256 via digestFile, only groups of 2+ files kept |
| st-005 | behavioral | The application shall ignore zero-byte files and shall neither follow nor report symbolic links. | MET | internal/walk/walk.go: `info.Size() == 0` skipped; `!d.Type().IsRegular()` skips symlinks; filepath.WalkDir does not descend symlinked directories; probe: tree with linkdir -> other/ and linkfile.txt -> real/f1.txt printed only real/f1.txt and other/f2.txt, rc=0; grader C/D/E cases with empty files and report-link.txt pass |
| st-006 | behavioral | When --summary is given, the application shall write the group count, the duplicate file count, and the reclaimable bytes of the listing that the same options would produce. | MET | internal/report/report.go WriteSummary: files += len(Paths)-1, reclaimable += Reclaimable(), groups = len(groups); grader B1 '5,6,1177', C2 '2,2,1128', D2 '0,0,0' all PASS |
| st-007 | behavioral | If the directory is missing or is not a directory, then the application shall exit 1; if the arguments are malformed, then the application shall exit 2; in both cases it shall write nothing to standard output. | MET | probe: `dupfind <missing>` rc=1 stdout_bytes=0, stderr 'stat ...: no such file or directory'; probe: `dupfind <regular file>` rc=1 stdout_bytes=0, stderr 'is not a directory'; grader U1-U5 compare stdout against empty expected output with rc 2, all PASS; internal/cli/cli.go Run(): usage errors return 2 and write only to stderr |
| st-008 | technical | The application shall only read the files it examines and shall never delete, move, or write any file under the directory. | MET | internal/dupes/dupes.go digestFile: os.Open (read-only) with io.Copy into the hash; internal/walk/walk.go: filepath.WalkDir and d.Info() only, no writes; probe: stat of every path in a tree before and after running dupfind was identical |

**Ranked improvements:**
1. st-008 is judged by reading the code; add a grader check that snapshots the tree's paths, sizes and mtimes before and after a run and fails on any change, so it becomes a proof rather than an LLM judgement.
2. st-005 says symlinks are never followed, but cli.validateDirectory resolves a symlinked root argument with filepath.EvalSymlinks. The behaviour is sensible, but the criterion should say whether a symlinked root argument is in scope so the grader can pin it.
3. st-002 names 'every automated test' but the POSIX requirement is only checked by reading the script; a `dash bin/test.sh` run in the grader would make the shell requirement observable.
