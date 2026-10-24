# Evidence: Block 1 · Foundational (block-1)

- block type: block
- date: 2026-10-24
- resulting state: closed/verified
- story points (combined assembled cost): 8212
- execution id: 20261024.073854.707Z-087e7a43

## Stories built
- Establish the Go module and command scaffold. (foundation) [story]

## Reusable compacts
- ARCHITECTURE_compact.md

## Stacked context
- implements: go.md (SP 4944)
- compass: COMPASS.md (SP 641)
- implements: ARCHITECTURE.md (SP 781)
- stack: common.md (SP 1812)

## Build directory changes
- cmd/dupfind/main.go
- go.mod
- internal/cli/cli.go
- internal/cli/scaffold_test.go
- internal/discovery/discovery.go
- internal/output/output.go

## Pre-build acceptance observation
- RED: module-build (ARCHITECTURE.md)
  intent: The module exposes the declared command package and builds it successfully.
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
      File "module-build.py", line 10, in <module>
        assert result.returncode == 0
               ^^^^^^^^^^^^^^^^^^^^^^
    AssertionError
- RED: internal-packages (ARCHITECTURE.md)
  intent: The application provides the declared internal package structure.
  return code: 1
  stderr:
    Traceback (most recent call last):
      File "internal-packages.py", line 4, in <module>
        assert importlib.util.find_spec(f"internal.{package}") is not None
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
      File "<frozen importlib.util>", line 91, in find_spec
    ModuleNotFoundError: No module named 'internal'

## Post-build programmatic acceptance
- PASS: module-build (ARCHITECTURE.md)
  intent: The module exposes the declared command package and builds it successfully.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
- PASS: internal-packages (ARCHITECTURE.md)
  intent: The application provides the declared internal package structure.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0

## Build summary
<reusable-compact filename="ARCHITECTURE.md">
- Build the executable from `./cmd/dupfind`; `main` delegates to `internal/cli`.
- Keep discovery and grouping in `internal/discovery`; keep ordering and CSV output in `internal/output`.
- Use the Go standard library. The application is read-only, and reads file contents only when another eligible file has the same size.
- Do not access the grader from application code.
</reusable-compact>

RESULT: SUCCESS

FILES CHANGED:
- go.mod
- cmd/dupfind/main.go
- internal/cli/cli.go
- internal/cli/scaffold_test.go
- internal/discovery/discovery.go
- internal/output/output.go

SUMMARY:
Established the Go 1.22 module and command scaffold, with the declared internal packages. `go test ./...`, `go build ./cmd/dupfind`, and the package-structure assertion passed.

BLOCKERS:
- None
