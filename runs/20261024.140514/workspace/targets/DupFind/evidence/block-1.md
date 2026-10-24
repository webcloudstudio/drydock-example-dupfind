# Evidence: Block 1 · Foundational (block-1)

- block type: block
- date: 2026-10-24
- resulting state: closed/verified
- story points (combined assembled cost): 6530
- execution id: 20261024.140854.923Z-00f76a5e

## Stories built
- Establish the Go module and application boundaries. (architecture) [story]

## Stacked context
- implements: go.md (SP 4944)
- compass: COMPASS.md (SP 641)
- implements: ARCHITECTURE.md (SP 864)

## Build directory changes
- cmd/dupfind/main.go
- cmd/dupfind/main_test.go
- go.mod
- internal/cli/cli.go
- internal/cli/cli_test.go
- internal/output/output.go
- internal/scan/scan.go

## Pre-build acceptance observation
- RED: architecture-cli-contract (ARCHITECTURE.md)
  intent: The architecture exposes the documented CLI runner contract and the executable package builds.
  return code: 1
  stdout:
    go: cannot find main module, but found .git/config in /mnt/c/Users/barlo/projects/drydock/uat/DupFind
    	to create a module there, run:
    	cd ../../../.. && go mod init
  stderr:
    --- drydock: values at failure ---
      result = CompletedProcess(args=['go', 'build', './cmd/dupfind'], returncode=1, stdout='', stderr='go: cannot find main module, but found .git/config in /mnt/c/Users/barlo/projects/drydock/uat/DupFind\n\tto create a module there, run:\n\tcd ../../../.. && go mod init\n')
    --- drydock: end values ---
    Traceback (most recent call last):
      File "architecture-cli-contract.py", line 10, in <module>
        assert result.returncode == 0
               ^^^^^^^^^^^^^^^^^^^^^^
    AssertionError
- RED: architecture-package-layout (ARCHITECTURE.md)
  intent: The internal application packages and command entry point compile together.
  return code: 1
  stdout:
    FAIL	./... [setup failed]
    FAIL
    
    # ./...
    pattern ./...: directory prefix . does not contain main module or its selected dependencies
  stderr:
    --- drydock: values at failure ---
      result = CompletedProcess(args=['go', 'test', './...'], returncode=1, stdout='FAIL\t./... [setup failed]\nFAIL\n', stderr='# ./...\npattern ./...: directory prefix . does not contain main module or its selected dependencies\n')
    --- drydock: end values ---
    Traceback (most recent call last):
      File "architecture-package-layout.py", line 10, in <module>
        assert result.returncode == 0
               ^^^^^^^^^^^^^^^^^^^^^^
    AssertionError

## Post-build programmatic acceptance
- PASS: architecture-cli-contract (ARCHITECTURE.md)
  intent: The architecture exposes the documented CLI runner contract and the executable package builds.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
- PASS: architecture-package-layout (ARCHITECTURE.md)
  intent: The internal application packages and command entry point compile together.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    ok  	dupfind/cmd/dupfind	3.494s
    ok  	dupfind/internal/cli	(cached)
    ?   	dupfind/internal/output	[no test files]
    ?   	dupfind/internal/scan	[no test files]

## Build summary
RESULT: SUCCESS

FILES CHANGED:
- go.mod
- cmd/dupfind/main.go
- cmd/dupfind/main_test.go
- internal/cli/cli.go
- internal/cli/cli_test.go
- internal/scan/scan.go
- internal/output/output.go

SUMMARY:
Established the Go 1.22 module and `cmd/dupfind`, `internal/cli`, `internal/scan`, and `internal/output` package boundaries. Added Go tests for the executable build, package layout, and initial CLI contract. `go build ./cmd/dupfind` and `go test ./...` passed.

BLOCKERS:
- None.
