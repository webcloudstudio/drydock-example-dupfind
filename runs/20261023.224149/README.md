# DUPFIND

- Short description: Drydock uat - DupFind is the simplest Go kit and exercises go build end to end. The source is an original brief for a command-line duplicate file finder that groups regular files by size and SHA-256 digest and reports each group or a summary of reclaimable bytes. The supplied grader, sources/run.sh, builds ./cmd/dupfind with go build, creates a fixed directory tree (hidden files, a symbolic link, empty files, a path with a comma), and runs 15 checks of exit code and output. Standard library only, no database. Run with drydock uat DupFind. Source Repo: none, the brief is original to this kit.
- Status: PLAN CREATED

1 of 4 receipt claims proven. Open `index.html` for the linked proof kit; verify it with `sha256sum -c SHA256SUMS`.

## Receipt

| Claim | Verdict | Recorded outcome | Proof |
|---|---|---|---|
| Lifecycle completed | FAIL | 4 lifecycle commands executed; the run ended at 04-plan. | [result.json](result.json) |
| External conformance suite passed | UNPROVEN | No external test command is defined for this project. | — |
| Release score passed | UNPROVEN | No release score was recorded for this run. | — |
| Integrity verification passed | PASS | 145 files digested; verify with sha256sum -c SHA256SUMS. | [SHA256SUMS](SHA256SUMS) |

## Run facts

- Drydock: `0.2.0` (commit `fd5d16d2f65171f4668e018c4c8955fa07e68212`)
- Provider and model: `codex` / `gpt-5.6-luna`
- Platform: `Linux-6.6.87.2-microsoft-standard-WSL2-x86_64-with-glibc2.35` on Python `3.12.13`
- Run: `20261023.224149`
- Elapsed: 154.2s
- LLM calls: 6
- Tokens: cached 100,480; uncached 111,961; output 15,427
- LLM elapsed: 115.1s
- Build passes: 0; repairs: not recorded
- Conformance: no external suite defined
- Advisory scores: none recorded

## RUN SUMMARY

- Input specification: [`sources/INSTRUCTIONS.md`](sources/INSTRUCTIONS.md)
- Delivered Code: [`build/DupFind/`](build/DupFind)
- Test Results: not recorded for this run.

## Shortfall

- Failure: DupFind: plan exited 1

## Commands

| # | Command | Exit | Elapsed | Output |
|---|---|---|---|---|
| 01-init | `/mnt/c/Users/barlo/projects/drydock/.venv/bin/python -m drydock init DupFind` | 0 | 3.4s | [stdout](evidence/commands/01-init.stdout.log) · [stderr](evidence/commands/01-init.stderr.log) |
| 02-import-sources | `/mnt/c/Users/barlo/projects/drydock/.venv/bin/python -m drydock import DupFind sources --format markdown` | 0 | 5.7s | [stdout](evidence/commands/02-import-sources.stdout.log) · [stderr](evidence/commands/02-import-sources.stderr.log) |
| 03-analyze | `/mnt/c/Users/barlo/projects/drydock/.venv/bin/python -m drydock analyze DupFind` | 0 | 25.0s | [stdout](evidence/commands/03-analyze.stdout.log) · [stderr](evidence/commands/03-analyze.stderr.log) · [llm](evidence/commands/03-analyze.llm.log) |
| 04-plan | `/mnt/c/Users/barlo/projects/drydock/.venv/bin/python -m drydock plan DupFind --override` | 1 | 118.3s | [stdout](evidence/commands/04-plan.stdout.log) · [stderr](evidence/commands/04-plan.stderr.log) · [llm](evidence/commands/04-plan.llm.log) |

## Evidence

- [`evidence/commands/`](evidence/commands) — stdout and stderr of every command
- [`evidence/prompts/`](evidence/prompts) — the assembled prompt for every LLM call
- [`evidence/provider_raw/`](evidence/provider_raw) — unmodified provider transcripts
- [`evidence/llm_logs/`](evidence/llm_logs) — call banners and token accounting
- [`result.json`](result.json) — the machine-readable record of this run
