# FEATURE: Verification

| Field       | Value |
|-------------|-------|
| Version     | 20261025 V1 |
| Description | Terminal release verification that the finished dupfind passes its own test gate and the supplied grader end to end. |
| Depends On  | FEATURE-TestGate.md, FEATURE-ModuleBuild.md |
| Provides    | None |
| Consumes    | bin/test.sh, cmd/dupfind |

## Purpose

This assembly story proves that the stories fit together. The supplied grader `sources/run.sh`
builds `./cmd/dupfind`, creates a fixed directory tree, and makes 15 checks of exit code and
standard output. The grader is the exam: it is never modified, regenerated, substituted, or read by
application code. A defect it exposes is fixed in the package that owns the behavior.

## Assembly

All implementation stories are complete before this story runs: argument parsing, directory
validation, walk, size grouping, digest grouping, group ordering, listing, summary, module build,
and the test gate. This story adds no code of its own.

## Reads

- `sources/run.sh`, executed as a script only, from the application root.
- `bin/test.sh`.

## Writes

- Nothing. The grader works in its own temporary directory.

## Programmatic Acceptance

=== AC verification-test-gate-passes ===
Intent: The project's own gate, bin/test.sh, exits zero on the assembled application.
Requires: executable=sh; scope=test
Requires: executable=go; scope=test

import subprocess
import sys

result = subprocess.run(["sh", "bin/test.sh"], capture_output=True, text=True)
print(result.stdout)
print(result.stderr, file=sys.stderr)
assert result.returncode == 0, result.returncode
=== END AC verification-test-gate-passes ===

=== AC verification-no-grader-coupling ===
Intent: Application code, tests excluded, never refers to the supplied grader script.
Requires: executable=go; scope=test

import os

needle = "run" + ".sh"
offenders = []
for base in ("cmd", "internal"):
    for folder, _dirs, files in os.walk(base):
        for name in files:
            if name.endswith(".go") and not name.endswith("_test.go"):
                path = os.path.join(folder, name)
                with open(path, encoding="utf-8") as handle:
                    if needle in handle.read():
                        offenders.append(path)
print(offenders)
assert offenders == [], offenders
=== END AC verification-no-grader-coupling ===

=== AC verification-grader-passes ===
Intent: The supplied grader sources/run.sh builds ./cmd/dupfind and exits zero with every check passing.
Suite: full
Requires: executable=sh; scope=test
Requires: executable=go; scope=test

import subprocess
import sys

result = subprocess.run(["sh", "sources/run.sh"], capture_output=True, text=True)
print(result.stdout)
print(result.stderr, file=sys.stderr)
assert result.returncode == 0, result.returncode
=== END AC verification-grader-passes ===

## User Acceptance

- None.

## Guardrails

- Do not modify, rewrite, regenerate, or substitute `sources/run.sh`.
- Application code never reads `sources/run.sh` and never hardcodes a digest, a path, or a total; every value is computed from the directory examined.
- Never delete, move, or write a file under the directory examined.
