# FEATURE: Listing CSV

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Writes the duplicate listing as RFC 4180 CSV with the exact header `group,size,sha256,path`. |
| Depends On  | FEATURE-Grouping.md, ARCHITECTURE.md |
| Provides    | output.WriteListing |
| Consumes    | scan.Group |

## Purpose

Render the duplicate groups produced by `scan.Group` as the default CSV listing on standard output.
This story owns the header, one row per file in every duplicate group, the lowercase hexadecimal
digest, and RFC 4180 quoting. Group ordering, path rendering, and group numbering belong to
FEATURE-ListingOrder.md. The empty listing belongs to FEATURE-ListingEmpty.md.

## Behavior

- `output.WriteListing` writes to an `io.Writer` using `encoding/csv`, with `\n` line endings.
- The first row is exactly `group,size,sha256,path`.
- Each file in each duplicate group produces one data row: group number, size in bytes, digest, path.
- The digest is lowercase hexadecimal SHA-256.
- A field containing a comma, a double quote, or a newline is quoted as RFC 4180 specifies.
- The writer accepts groups already ordered and numbered by FEATURE-ListingOrder.md and does not reorder them.

## Interface

- `output.WriteListing(w io.Writer, groups []scan.Group) error`

## Test Contract

The `internal/output` package carries one Go test per criterion below, named exactly as stated. Each
test builds real groups through the public interface and asserts the observed bytes or parsed fields.

## Programmatic Acceptance

=== AC listing-csv-header ===
Intent: The listing writer emits the exact header row group,size,sha256,path.

from drydock_go import run_test

result = run_test("./internal/output", "TestWriteListingHeader")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC listing-csv-header ===

=== AC listing-csv-rows ===
Intent: The listing writer emits one data row per file in every duplicate group.

from drydock_go import run_test

result = run_test("./internal/output", "TestWriteListingRowsPerFile")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC listing-csv-rows ===

=== AC listing-csv-digest-lowercase ===
Intent: The listing writer renders the sha256 field as lowercase hexadecimal.

from drydock_go import run_test

result = run_test("./internal/output", "TestWriteListingDigestLowercase")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC listing-csv-digest-lowercase ===

=== AC listing-csv-rfc4180-quoting ===
Intent: Fields containing commas, double quotes, or newlines are quoted per RFC 4180 and parse back to the original path.

from drydock_go import run_test

result = run_test("./internal/output", "TestWriteListingRFC4180Quoting")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC listing-csv-rfc4180-quoting ===

=== AC listing-csv-line-endings ===
Intent: The listing writer terminates every record with a single LF and never emits CR.

from drydock_go import run_test

result = run_test("./internal/output", "TestWriteListingLineEndings")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC listing-csv-line-endings ===

## User Acceptance

- None.

## Guardrails

- The writer reads no files and writes nothing except to the supplied `io.Writer`.
- The writer never writes diagnostics to standard output.
- The writer does not reorder or renumber groups; ordering belongs to FEATURE-ListingOrder.md.
