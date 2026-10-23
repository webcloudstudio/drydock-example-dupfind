# COMPASS: DupFind

## Compass

Build a command-line duplicate file finder, the Go executable `./cmd/dupfind`, as described in
`sources/INSTRUCTIONS.md`. It walks a directory, groups regular files by size and then by SHA-256
digest, and reports each duplicate group or a one-row summary of the space that removing the
copies would recover.

Correctness is measured by `sources/run.sh`: it builds the executable with `go build` and makes
15 checks of exit code and standard output against a fixed directory tree. The goal is every check
passing.

## Constraints

- Implement in Go 1.22 or newer with the standard library only. No third-party modules.
- The application root is a Go module; the executable's main package is `./cmd/dupfind`.
- Read a file's contents only when another file has the same size.
- Write CSV with `encoding/csv`; diagnostics go to standard error and are never compared.

## Guardrails

- Do not modify, rewrite, regenerate, or substitute `sources/run.sh`.
- Do not read `sources/run.sh` from application code, and do not hardcode a digest, a path, or a
  total. Every value is computed from the directory examined.
- Never delete, move, or write a file under the directory examined.

## Verification Protocol

| Story kind | Command | Asserts |
|---|---|---|
| Every behavioral story | `sh bin/test.sh` | exit `0` |
| Terminal story (once, last) | `sh sources/run.sh` | exit `0` |

<!-- drydock:build-write-guardrail:start -->
## Build Write Guardrail

- Authorized build directory: `/mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261023.224149/build/DupFind`
- Authorized Target directory: `/mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261023.224149/workspace/targets/DupFind`
- Build agents have permission to create, modify, and remove files required by the active build block inside these authorized directories.
- No path outside these authorized directories may be modified.
- Protected Drydock artifacts:
  - `/mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261023.224149/workspace/targets/DupFind/blueprint/`
  - `/mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261023.224149/workspace/targets/DupFind/MANIFEST.md`
  - `/mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261023.224149/workspace/targets/DupFind/COMPASS.md`
  - `/mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261023.224149/workspace/targets/DupFind/QuarterDeck/`
  - `/mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261023.224149/workspace/targets/DupFind/evidence/`
<!-- drydock:build-write-guardrail:end -->
