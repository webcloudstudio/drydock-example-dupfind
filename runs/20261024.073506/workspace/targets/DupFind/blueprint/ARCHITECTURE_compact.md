<!-- Compacted from ARCHITECTURE.md sha256=c86333e5266241be5ba8662f0c3c3cdd9767d9f6a6ff07f171b21bacd125984f on 2026-10-24 by drydock build agent -->

- Build the executable from `./cmd/dupfind`; `main` delegates to `internal/cli`.
- Keep discovery and grouping in `internal/discovery`; keep ordering and CSV output in `internal/output`.
- Use the Go standard library. The application is read-only, and reads file contents only when another eligible file has the same size.
- Do not access the grader from application code.
