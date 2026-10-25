# FEATURE: Module Build

| Field       | Value |
|-------------|-------|
| Version     | 20261025 V1 |
| Description | The cmd/dupfind main package builds with go build using only the standard library and delegates to cli.Run. |
| Depends On  | ARCHITECTURE.md, FEATURE-ArgParsing.md |
| Provides    | cmd/dupfind |
| Consumes    | cli.Run |

## Purpose

The application root is a Go module whose executable is the main package at `./cmd/dupfind`.
`go build ./cmd/dupfind` must succeed using only the standard library.

## Behavior

- `cmd/dupfind/main.go` is `package main`. `main` calls `cli.Run(os.Args[1:], os.Stdout, os.Stderr)` and passes the returned code to `os.Exit`.
- `main` contains no other logic and is the only caller of `os.Exit`.
- Every package in the build's dependency graph is either in the standard library or inside this module.

## Reads

Process arguments, standard output, standard error.

## Writes

None.

## Programmatic Acceptance

=== AC modbuild-go-build ===
Intent: go build ./cmd/dupfind succeeds and produces an executable.

import os
import subprocess
import sys
import tempfile

with tempfile.TemporaryDirectory() as d:
    out_path = os.path.join(d, "dupfind")
    result = subprocess.run(["go", "build", "-o", out_path, "./cmd/dupfind"], capture_output=True, text=True)
    print(result.stdout)
    print(result.stderr, file=sys.stderr)
    assert result.returncode == 0
    assert os.path.isfile(out_path)
=== END AC modbuild-go-build ===

=== AC modbuild-standard-library-only ===
Intent: The executable's dependency graph contains only standard library and in-module packages.

import pathlib
import subprocess
import sys

mod = ""
for line in pathlib.Path("go.mod").read_text().splitlines():
    if line.startswith("module "):
        mod = line.split()[1]
assert mod
result = subprocess.run(
    ["go", "list", "-deps", "-f", "{{if not .Standard}}{{.ImportPath}}{{end}}", "./cmd/dupfind"],
    capture_output=True, text=True,
)
print(result.stdout)
print(result.stderr, file=sys.stderr)
assert result.returncode == 0
for pkg in result.stdout.split():
    assert pkg == mod or pkg.startswith(mod + "/"), pkg
=== END AC modbuild-standard-library-only ===

=== AC modbuild-executable-runs ===
Intent: The built executable runs and rejects an empty command line with exit 2.

from drydock_go import cli_runner

code, out, err = cli_runner()
print(code, repr(out), repr(err))
assert code == 2
assert out == ""
=== END AC modbuild-executable-runs ===

## User Acceptance

- None.

## Guardrails

- Standard library only; no third-party modules.
- `main` never reads `sources/run.sh` and embeds no digest, path, or total.
- Only `main` calls `os.Exit`.
