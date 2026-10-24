# DUPFIND

- Short description: Drydock uat - DupFind is the simplest Go kit and exercises go build end to end. The source is an original brief for a command-line duplicate file finder that groups regular files by size and SHA-256 digest and reports each group or a summary of reclaimable bytes. The supplied grader, sources/run.sh, builds ./cmd/dupfind with go build, creates a fixed directory tree (hidden files, a symbolic link, empty files, a path with a comma), and runs 15 checks of exit code and output. Standard library only, no database. Run with drydock uat DupFind. Source Repo: none, the brief is original to this kit.
- Status: PLAN CREATED

1 of 4 receipt claims proven. Open `index.html` for the linked proof kit; verify it with `sha256sum -c SHA256SUMS`.

## Receipt

| Claim | Verdict | Recorded outcome | Proof |
|---|---|---|---|
| Lifecycle completed | FAIL | 10 lifecycle commands executed; the run ended at 10-after-plan-workspace-status. | [result.json](result.json) |
| External conformance suite passed | UNPROVEN | No external test command is defined for this project. | — |
| Release score passed | UNPROVEN | No release score was recorded for this run. | — |
| Integrity verification passed | PASS | 200 files digested; verify with sha256sum -c SHA256SUMS. | [SHA256SUMS](SHA256SUMS) |

## Run facts

- Drydock: `0.8.0` (commit `d9553a63c6311a238625f729ccd1c32e8b7c3230`)
- Provider and model: `codex` / `gpt-6-luna`
- Platform: `Linux-6.6.87.2-microsoft-standard-WSL2-x86_64-with-glibc2.35` on Python `3.12.13`
- Run: `20261024.080843`
- Elapsed: 212.2s
- LLM calls: 9
- Tokens: cached 160,512; uncached 161,901; output 11,244
- LLM elapsed: 106.7s
- Build passes: 0; repairs: not recorded
- Conformance: no external suite defined
- Advisory scores: none recorded

## RUN SUMMARY

- Input specification: [`sources/INSTRUCTIONS.md`](sources/INSTRUCTIONS.md)
- Delivered Code: [`build/DupFind/`](build/DupFind)
- Test Results: not recorded for this run.

## Shortfall

- Failure: DupFind: stopped at plan-repair: acceptance criteria still cannot run after one repair pass

## Commands

| # | Command | Exit | Elapsed | Output |
|---|---|---|---|---|
| 01-init | `/mnt/c/Users/barlo/projects/drydock/.venv/bin/python -m drydock init DupFind` | 0 | 3.5s | [stdout](evidence/commands/01-init.stdout.log) · [stderr](evidence/commands/01-init.stderr.log) |
| 02-import-sources | `/mnt/c/Users/barlo/projects/drydock/.venv/bin/python -m drydock import DupFind sources --format markdown` | 0 | 5.9s | [stdout](evidence/commands/02-import-sources.stdout.log) · [stderr](evidence/commands/02-import-sources.stderr.log) |
| 03-analyze | `/mnt/c/Users/barlo/projects/drydock/.venv/bin/python -m drydock analyze DupFind` | 0 | 22.4s | [stdout](evidence/commands/03-analyze.stdout.log) · [stderr](evidence/commands/03-analyze.stderr.log) · [llm](evidence/commands/03-analyze.llm.log) |
| 04-plan | `/mnt/c/Users/barlo/projects/drydock/.venv/bin/python -m drydock plan DupFind --override` | 0 | 116.5s | [stdout](evidence/commands/04-plan.stdout.log) · [stderr](evidence/commands/04-plan.stderr.log) · [llm](evidence/commands/04-plan.llm.log) |
| 05-plan-verify | `/mnt/c/Users/barlo/projects/drydock/.venv/bin/python -m drydock plan verify DupFind` | 1 | 13.0s | [stdout](evidence/commands/05-plan-verify.stdout.log) · [stderr](evidence/commands/05-plan-verify.stderr.log) |
| 06-plan-repair | `/mnt/c/Users/barlo/projects/drydock/.venv/bin/python -m drydock plan repair DupFind` | 1 | 29.1s | [stdout](evidence/commands/06-plan-repair.stdout.log) · [stderr](evidence/commands/06-plan-repair.stderr.log) · [llm](evidence/commands/06-plan-repair.llm.log) |
| 07-plan-verify-after-repair | `/mnt/c/Users/barlo/projects/drydock/.venv/bin/python -m drydock plan verify DupFind` | 1 | 12.3s | [stdout](evidence/commands/07-plan-verify-after-repair.stdout.log) · [stderr](evidence/commands/07-plan-verify-after-repair.stderr.log) |
| 08-after-plan-build-status | `/mnt/c/Users/barlo/projects/drydock/.venv/bin/python -m drydock build status DupFind` | 0 | 2.9s | [stdout](evidence/commands/08-after-plan-build-status.stdout.log) · [stderr](evidence/commands/08-after-plan-build-status.stderr.log) |
| 09-after-plan-target-status | `/mnt/c/Users/barlo/projects/drydock/.venv/bin/python -m drydock status DupFind` | 0 | 2.1s | [stdout](evidence/commands/09-after-plan-target-status.stdout.log) · [stderr](evidence/commands/09-after-plan-target-status.stderr.log) |
| 10-after-plan-workspace-status | `/mnt/c/Users/barlo/projects/drydock/.venv/bin/python -m drydock status` | 0 | 1.3s | [stdout](evidence/commands/10-after-plan-workspace-status.stdout.log) · [stderr](evidence/commands/10-after-plan-workspace-status.stderr.log) |

## Evidence

- [`evidence/commands/`](evidence/commands) — stdout and stderr of every command
- [`evidence/prompts/`](evidence/prompts) — the assembled prompt for every LLM call
- [`evidence/provider_raw/`](evidence/provider_raw) — unmodified provider transcripts
- [`evidence/llm_logs/`](evidence/llm_logs) — call banners and token accounting
- [`result.json`](result.json) — the machine-readable record of this run
