# Evidence: Block 1 · Foundational (block-1)

- block type: block
- date: 2026-10-25
- resulting state: closed/verified
- story points (combined assembled cost): 13594
- execution id: 20261025.174754.270Z-8643e753

## Stories built
- Stand up the Go module, package tree, and test launcher that every other story builds into. (foundation) [story]

## Stacked context
- implements: go.md (SP 7578)
- compass: COMPASS.md (SP 982)
- implements: ARCHITECTURE.md (SP 2007)
- stack: common.md (SP 2776)

## Build directory changes
- bin/.gitkeep
- cmd/dupfind/main.go
- go.mod
- internal/cli/cli.go
- internal/cli/cli_test.go
- internal/dupes/dupes.go
- internal/report/report.go
- internal/walk/walk.go

## Pre-build acceptance observation
- RED: arch-module-no-third-party (ARCHITECTURE.md)
  intent: The application root is a Go module that declares no third-party requirement.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      self = PosixPath('go.mod')
      mode = 'r'
      buffering = -1
      encoding = 'locale'
      errors = None
      newline = None
    --- drydock: end values ---
    Traceback (most recent call last):
      File "arch-module-no-third-party.py", line 3, in <module>
        text = pathlib.Path("go.mod").read_text()
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
      File "/home/barlo/.local/share/uv/python/cpython-3.12.13-linux-x86_64-gnu/lib/python3.12/pathlib.py", line 1027, in read_text
        with self.open(mode='r', encoding=encoding, errors=errors) as f:
             ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
      File "/home/barlo/.local/share/uv/python/cpython-3.12.13-linux-x86_64-gnu/lib/python3.12/pathlib.py", line 1013, in open
        return io.open(self, mode, buffering, encoding, errors, newline)
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
    FileNotFoundError: [Errno 2] No such file or directory: 'go.mod'
- RED: arch-package-tree (ARCHITECTURE.md)
  intent: The package directories named by the architecture exist in the build tree.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      rel = 'cmd/dupfind'
    --- drydock: end values ---
    Traceback (most recent call last):
      File "arch-package-tree.py", line 4, in <module>
        assert pathlib.Path(rel).is_dir(), rel
               ^^^^^^^^^^^^^^^^^^^^^^^^^^
    AssertionError: cmd/dupfind
- RED: arch-vet-and-build (ARCHITECTURE.md)
  intent: The foundation compiles and passes go vet across every package.
  return code: 1
  stderr:
    pattern ./...: directory prefix . does not contain main module or its selected dependencies
    
    --- drydock: values at failure ---
      cmd = ['go', 'vet', './...']
      result = CompletedProcess(args=['go', 'vet', './...'], returncode=1, stdout='', stderr='pattern ./...: directory prefix . does not contain main module or its selected dependencies\n')
    --- drydock: end values ---
    Traceback (most recent call last):
      File "arch-vet-and-build.py", line 8, in <module>
        assert result.returncode == 0, cmd
               ^^^^^^^^^^^^^^^^^^^^^^
    AssertionError: ['go', 'vet', './...']

## Post-build programmatic acceptance
- PASS: arch-module-no-third-party (ARCHITECTURE.md)
  intent: The application root is a Go module that declares no third-party requirement.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
- PASS: arch-package-tree (ARCHITECTURE.md)
  intent: The package directories named by the architecture exist in the build tree.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
- PASS: arch-vet-and-build (ARCHITECTURE.md)
  intent: The foundation compiles and passes go vet across every package.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0

## Build summary
The build directory holds only `sources/run.sh` and an empty git repo. The foundation spec calls for a root module, the package tree, and placeholders that compile, so I'll write those now. `run.sh` is only for the terminal story, so I won't run it.Files are written. Next I'll run the gates the foundation acceptance criteria name: gofmt, vet, build, the tests, and the two AC scripts.RESULT: SUCCESS

FILES CHANGED:
- go.mod
- cmd/dupfind/main.go
- internal/cli/cli.go
- internal/cli/cli_test.go
- internal/walk/walk.go
- internal/dupes/dupes.go
- internal/report/report.go
- bin/.gitkeep

SUMMARY:
The foundation is in place. The root module is `dupfind` with `go 1.22` and no `require` lines. `cmd/dupfind/main.go` calls `cli.Run` and is the only place that calls `os.Exit`. `internal/cli` exposes `Run(args, stdout, stderr) int`, which returns 0 for now. `internal/walk`, `internal/dupes`, and `internal/report` are package stubs that compile. `bin/` exists.

Checks run, all passing:
- `gofmt -l .` printed nothing.
- `go vet ./...` and `go build ./...` passed, and `go build ./cmd/dupfind` passed (Go 1.24.6 toolchain).
- `go test ./...` passed. The one test, `TestRunPlaceholderSucceedsSilently`, checks that `Run` returns 0 and writes nothing to stdout.
- AC `arch-module-no-third-party` passed.
- AC `arch-package-tree` passed.
- AC `arch-vet-and-build` passed (the vet and build commands above).

Not run: `sources/run.sh`. It is the terminal-story grader and will fail until the feature stories land.

Choices I made:
- Module path is `dupfind`. The spec did not name one.
- `bin/test.sh` is not written. ARCHITECTURE says the test-gate story delivers it, and that file is authoritative for this step. The story title mentions a "test launcher", so that wording conflicts with ARCHITECTURE. `sh bin/test.sh` will not exist until the test-gate story runs.
- `bin/.gitkeep` keeps the empty `bin/` directory in git.
- No commit was made, per the build rules.

BLOCKERS:
- None.
