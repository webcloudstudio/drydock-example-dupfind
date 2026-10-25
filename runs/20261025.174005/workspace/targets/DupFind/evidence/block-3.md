# Evidence: Block 3 · Feature (block-3)

- block type: block
- date: 2026-10-25
- resulting state: closed/verified
- story points (combined assembled cost): 6074
- execution id: 20261025.175614.224Z-5c16fcc2

## Stories built
- Assemble the finished tool and confirm the supplied grader passes end to end. (verification) [story]

## Acceptance tooling authorization
- VERIFICATION.md#verification-test-gate-passes: executable=sh; scope=test; authorization=existing Target environment
- VERIFICATION.md#verification-test-gate-passes: executable=go; scope=test; authorization=existing Target environment
- VERIFICATION.md#verification-no-grader-coupling: executable=go; scope=test; authorization=existing Target environment
- VERIFICATION.md#verification-grader-passes: executable=sh; scope=test; authorization=existing Target environment
- VERIFICATION.md#verification-grader-passes: executable=go; scope=test; authorization=existing Target environment

## Stacked context
- stack: go_compact.md (SP 2105)
- compass: COMPASS.md (SP 982)
- implements: VERIFICATION.md (SP 1227)
- stack: common_compact.md (SP 1653)

## Build directory changes
- data/logs/test_2026-10-25_135622.log

## Pre-build acceptance observation
- GREEN (prepassed): verification-test-gate-passes (VERIFICATION.md)
  intent: The project's own gate, bin/test.sh, exits zero on the assembled application.
  return code: 0
  stdout:
    timestamp: 2026-10-25_135604
    arguments: 
    working directory: /mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261025.174005/build/DupFind
    ?   	dupfind/cmd/dupfind	[no test files]
    ok  	dupfind/internal/cli	(cached)
    ok  	dupfind/internal/dupes	(cached)
    ok  	dupfind/internal/report	(cached)
    ok  	dupfind/internal/walk	(cached)
    completed: exit 0
- GREEN (prepassed): verification-no-grader-coupling (VERIFICATION.md)
  intent: Application code, tests excluded, never refers to the supplied grader script.
  return code: 0
  stdout:
    []
- GREEN (prepassed): verification-grader-passes (VERIFICATION.md)
  intent: The supplied grader sources/run.sh builds ./cmd/dupfind and exits zero with every check passing.
  return code: 0
  stdout:
    PASS A1: dupfind tree
    PASS A2: dupfind /tmp/tmp.Qy3YkyFpCI/tree
    PASS B1: dupfind --summary tree
    PASS C1: dupfind --min-size 20 tree
    PASS C2: dupfind --summary --min-size 20 tree
    PASS C3: dupfind --min-size 200 tree
    PASS D1: dupfind tree/unique
    PASS D2: dupfind --summary tree/empty
    PASS E1: dupfind tree/no-such-directory
    PASS E2: dupfind tree/env.bak
    PASS U1: dupfind 
    PASS U2: dupfind tree tree/a
    PASS U3: dupfind --min-size abc tree
    PASS U4: dupfind --min-size 0 tree
    PASS U5: dupfind --frobnicate tree
    DupFind: 15 passed, 0 failed (of 15)

## Post-build programmatic acceptance
- PASS: verification-test-gate-passes (VERIFICATION.md)
  intent: The project's own gate, bin/test.sh, exits zero on the assembled application.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    timestamp: 2026-10-25_135642
    arguments: 
    working directory: /mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261025.174005/build/DupFind
    ?   	dupfind/cmd/dupfind	[no test files]
    ok  	dupfind/internal/cli	(cached)
    ok  	dupfind/internal/dupes	(cached)
    ok  	dupfind/internal/report	(cached)
    ok  	dupfind/internal/walk	(cached)
    completed: exit 0
- PASS: verification-no-grader-coupling (VERIFICATION.md)
  intent: Application code, tests excluded, never refers to the supplied grader script.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    []
- PASS: verification-grader-passes (VERIFICATION.md)
  intent: The supplied grader sources/run.sh builds ./cmd/dupfind and exits zero with every check passing.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    PASS A1: dupfind tree
    PASS A2: dupfind /tmp/tmp.TQE9axe6E5/tree
    PASS B1: dupfind --summary tree
    PASS C1: dupfind --min-size 20 tree
    PASS C2: dupfind --summary --min-size 20 tree
    PASS C3: dupfind --min-size 200 tree
    PASS D1: dupfind tree/unique
    PASS D2: dupfind --summary tree/empty
    PASS E1: dupfind tree/no-such-directory
    PASS E2: dupfind tree/env.bak
    PASS U1: dupfind 
    PASS U2: dupfind tree tree/a
    PASS U3: dupfind --min-size abc tree
    PASS U4: dupfind --min-size 0 tree
    PASS U5: dupfind --frobnicate tree
    DupFind: 15 passed, 0 failed (of 15)

## Build summary
I'll start by inspecting the build directory and the grader.The grader is in place. Next I'll read the existing implementation and run both gates to see where things stand.Both gates pass already. The verification AC `verification-no-grader-coupling` is still unchecked, so I'll run that scan, plus an uncached test run to rule out stale results.The finished tool passes both gates. This step adds no code, so I changed no files.

- `sh bin/test.sh` exits 0: `go vet` is clean and all four packages pass.
- `sh sources/run.sh` exits 0: 15 passed, 0 failed.
- `verification-no-grader-coupling` passes: no non-test Go file under `cmd/` or `internal/` mentions the grader.
- An uncached `go vet` and `go test -race -count=1 ./...` also passes.

`git status` is clean after these runs; the test gate's logs under `data/logs/` are git-ignored.

RESULT: SUCCESS

FILES CHANGED:
- None

SUMMARY:
Verification-only terminal story; no application changes were needed. Checks run:
- `sh bin/test.sh`: exit 0. `go vet ./...` clean; `go test ./...` passes for `internal/cli`, `internal/dupes`, `internal/report`, `internal/walk`.
- `sh sources/run.sh`: exit 0, "DupFind: 15 passed, 0 failed (of 15)".
- `verification-no-grader-coupling` (Python scan of non-test `.go` files under `cmd/` and `internal/`): no offenders.
- `go vet ./...` and `go test -race -count=1 ./...`: exit 0.

The grader was run as-is and not modified, regenerated, or read by application code. Test pass/fail counts for this step: 15/15 grader checks; the four package suites all passed. The `verification-test-gate-passes`, `verification-no-grader-coupling`, and `verification-grader-passes` assertions each passed on this run.

BLOCKERS:
- None
