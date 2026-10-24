# BIG ERRORS — action required

- Command: `plan`
- Phase: plan verify
- State: Error
- Timestamp: 2026-10-24T07:02:16+00:00
- Execution ID: 20261024.070147.767Z-7b65894b
- Challenge Execution ID: -
- Classification: plan cannot be built as written
- Evidence / logs: /mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261024.065948/workspace/logs

## Diagnostic

plan verify: DupFind
  4 criteria cannot run as written
  ✗  FEATURE-Duplicate-Finding.md [duplicate-listing]: runs the application with `go run` at line(s) 5, which reports every non-zero exit as 1; use `from drydock_go import cli_runner` and `code, out, err = cli_runner(<args>)`, which runs the built executable
  ✗  FEATURE-Duplicate-Finding.md [duplicate-summary]: runs the application with `go run` at line(s) 5, which reports every non-zero exit as 1; use `from drydock_go import cli_runner` and `code, out, err = cli_runner(<args>)`, which runs the built executable
  ✗  FEATURE-Duplicate-Finding.md [duplicate-min-size]: runs the application with `go run` at line(s) 4, 12, which reports every non-zero exit as 1; use `from drydock_go import cli_runner` and `code, out, err = cli_runner(<args>)`, which runs the built executable
  ✗  FEATURE-Duplicate-Finding.md [duplicate-errors]: runs the application with `go run` at line(s) 4, 13, 22, 31, which reports every non-zero exit as 1; use `from drydock_go import cli_runner` and `code, out, err = cli_runner(<args>)`, which runs the built executable
  Run: drydock plan repair DupFind

## Recovery

Run: drydock plan repair DupFind
