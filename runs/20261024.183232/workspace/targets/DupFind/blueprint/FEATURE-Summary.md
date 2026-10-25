# FEATURE: Summary

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Writes the --summary CSV totals computed from the same grouping the listing would produce. |
| Depends On  | FEATURE-Grouping.md, ARCHITECTURE.md |
| Provides    | output.WriteSummary |
| Consumes    | scan.Group |

## Purpose

Report the number of duplicate groups, the number of duplicate files, and the reclaimable bytes for
the same options that produce the listing. The totals are derived from the same groups, so the
summary and listing never disagree.

## Behavior

- `output.WriteSummary` writes the header `groups,duplicate_files,reclaimable_bytes` and exactly one totals row, using `encoding/csv` with `\n` line endings.
- `groups` is the number of duplicate groups.
- `duplicate_files` is the sum of `files − 1` over all groups.
- `reclaimable_bytes` is the sum of every group's `size × (files − 1)`.
- With no duplicate groups the totals row is `0,0,0`.

## Interface

- `output.WriteSummary(w io.Writer, groups []scan.Group) error`

## Test Contract

The `internal/output` package carries one Go test per criterion below, named exactly as stated.

## Programmatic Acceptance

=== AC summary-header ===
Intent: The summary writer emits the exact header groups,duplicate_files,reclaimable_bytes.

from drydock_go import run_test

result = run_test("./internal/output", "TestWriteSummaryHeader")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC summary-header ===

=== AC summary-totals ===
Intent: The totals row reports group count, duplicate file count, and reclaimable bytes computed from the supplied groups.

from drydock_go import run_test

result = run_test("./internal/output", "TestWriteSummaryTotals")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC summary-totals ===

=== AC summary-no-duplicates ===
Intent: With no duplicate groups the summary totals row is 0,0,0.

from drydock_go import run_test

result = run_test("./internal/output", "TestWriteSummaryNoDuplicatesZeroRow")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC summary-no-duplicates ===

=== AC summary-matches-listing ===
Intent: The summary totals equal the sums taken over the same grouping that WriteListing renders.

from drydock_go import run_test

result = run_test("./internal/output", "TestWriteSummaryMatchesListing")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC summary-matches-listing ===

=== AC summary-csv-format ===
Intent: The summary is a two-line CSV document with LF line endings and no trailing blank row.

from drydock_go import run_test

result = run_test("./internal/output", "TestWriteSummaryCSVFormat")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC summary-csv-format ===

## User Acceptance

- None.

## Guardrails

- The summary is computed only from the supplied groups; it reads no file contents and no paths on disk.
- The summary never writes a row other than the header and the single totals row.
