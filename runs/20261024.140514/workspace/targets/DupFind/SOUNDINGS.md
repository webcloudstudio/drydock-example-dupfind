# Soundings

Per-assertion acceptance board, one row per Blueprint Programmatic Acceptance check.
`drydock plan` projects every assertion as `— UNVERIFIED`; `drydock score ac` sets the
Status column from deterministic proof results. A rerun of `drydock plan` resets Status to
`— UNVERIFIED`; rescore to refresh.

| Status | Blueprint | AC Id | Text | Binding | Baseline | Previously Passed | Last Block | Evidence | Verified At |
|---|---|---|---|---|---|---|---|---|---|
| ✓ PASS | ARCHITECTURE.md | architecture-cli-contract | The architecture exposes the documented CLI runner contract and the executable package builds. | yes | VALID baseline-red | yes | block-2 |  | 2026-10-24T14:13:38+00:00 |
| ✓ PASS | ARCHITECTURE.md | architecture-package-layout | The internal application packages and command entry point compile together. | yes | VALID baseline-red | yes | block-2 |  | 2026-10-24T14:13:38+00:00 |
| ~ PREPASSED | FEATURE-Build.md | build-command | Building the declared executable package succeeds. | yes | VALID existing-behavior | yes | block-2 | green at this block's baseline too, before its code existed — confirm the criterion exercises the story's work | 2026-10-24T14:13:38+00:00 |
| ~ PREPASSED | FEATURE-Build.md | build-test-entry | The Go test runner recognizes the application packages. | yes | VALID existing-behavior | yes | block-2 | green at this block's baseline too, before its code existed — confirm the criterion exercises the story's work | 2026-10-24T14:13:38+00:00 |
| ~ PREPASSED | FEATURE-Command.md | command-help-usage | Help includes a description of the command and its supported options. | yes | VALID existing-behavior | yes | block-2 | green at this block's baseline too, before its code existed — confirm the criterion exercises the story's work | 2026-10-24T14:13:38+00:00 |
| ✓ PASS | FEATURE-Command.md | command-invalid-option | Unknown options return the usage-error exit status. | yes | VALID baseline-red | yes | block-2 |  | 2026-10-24T14:13:38+00:00 |
| ✓ PASS | FEATURE-Command.md | command-invalid-size | A zero minimum size returns the usage-error exit status. | yes | VALID baseline-red | yes | block-2 |  | 2026-10-24T14:13:38+00:00 |
| ✓ PASS | FEATURE-Scanning.md | scan-groups-exact-files | The supplied conformance grader passes its duplicate scanning cases. | yes | VALID baseline-red | yes | block-2 |  | 2026-10-24T14:13:38+00:00 |
| ✓ PASS | FEATURE-Scanning.md | scan-symlinks-and-empty | The supplied conformance grader passes its symbolic-link and empty-file cases. | yes | VALID baseline-red | yes | block-2 |  | 2026-10-24T14:13:38+00:00 |
| ✓ PASS | FEATURE-Scanning.md | scan-read-only | The supplied conformance grader passes while examined files remain untouched. | yes | VALID baseline-red | yes | block-2 |  | 2026-10-24T14:13:38+00:00 |
| ✓ PASS | FEATURE-Ordering.md | ordering-conformance | The supplied conformance grader passes the duplicate ordering cases. | yes | VALID baseline-red | yes | block-2 |  | 2026-10-24T14:13:38+00:00 |
| ✓ PASS | FEATURE-Listing.md | listing-output | The command emits the specified CSV header and one row for each duplicate file, with paths parsed as CSV fields. | yes | VALID baseline-red | yes | block-2 |  | 2026-10-24T14:13:38+00:00 |
| ~ PREPASSED | FEATURE-Listing.md | listing-empty | A directory with no duplicates produces a header-only listing. | yes | VALID existing-behavior | yes | block-2 | green at this block's baseline too, before its code existed — confirm the criterion exercises the story's work | 2026-10-24T14:13:38+00:00 |
| ✓ PASS | FEATURE-Summary.md | summary-totals | The summary reports totals derived from the same duplicate groups as the listing. | yes | VALID baseline-red | yes | block-2 |  | 2026-10-24T14:13:38+00:00 |
| ✓ PASS | FEATURE-Summary.md | summary-empty | A directory with no duplicates reports zero totals. | yes | VALID baseline-red | yes | block-2 |  | 2026-10-24T14:13:38+00:00 |
| ✓ PASS | FEATURE-Tests.md | test-harness | The test entry point runs vet and the Go test suite and succeeds when both succeed. | yes | VALID baseline-red | yes | block-2 |  | 2026-10-24T14:13:38+00:00 |
| ✓ PASS | VERIFICATION.md | full-grader | The complete supplied grader passes all 15 checks. | yes | VALID baseline-red | yes | block-2 |  | 2026-10-24T14:13:38+00:00 |
