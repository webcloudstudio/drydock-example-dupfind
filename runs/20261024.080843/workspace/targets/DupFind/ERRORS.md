# BIG ERRORS — action required

- Command: `plan`
- Phase: plan verify
- State: Error
- Timestamp: 2026-10-24T08:11:09+00:00
- Execution ID: 20261024.081025.381Z-ff20c60b
- Challenge Execution ID: -
- Classification: plan cannot be built as written
- Evidence / logs: /mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261024.080843/workspace/logs

## Diagnostic

plan verify: DupFind
  2 criteria cannot run as written
  ✗  FEATURE-Command.md [command-errors]: Cite: sources/INSTRUCTIONS.md:51-66 contains none of this criterion's expected literals verbatim; cite the lines that state the expected value
  ✗  FEATURE-Command.md [command-success]: Cite: sources/INSTRUCTIONS.md:36-50 contains none of this criterion's expected literals verbatim; cite the lines that state the expected value
  Run: drydock plan repair DupFind

## Recovery

Run: drydock plan repair DupFind
