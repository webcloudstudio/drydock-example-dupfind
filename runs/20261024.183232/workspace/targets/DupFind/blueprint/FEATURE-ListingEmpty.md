# FEATURE: Listing Empty

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Emits only the listing header when no duplicate groups exist. |
| Depends On  | FEATURE-ListingOrder.md |
| Provides    | output.WriteListing empty case |
| Consumes    | output.Order |

## Purpose

When the scan finds no duplicate groups, the listing is the header alone. The output is exactly the
header followed by one newline, with no data rows and no trailing blank line.

## Behavior

- A nil or empty group list writes exactly `group,size,sha256,path\n`.
- The empty case writes nothing else to the writer and returns no error.

## Interface

- `output.WriteListing(w io.Writer, groups []output.OrderedGroup) error` with an empty group slice.

## Test Contract

The `internal/output` package carries one Go test per criterion below, named exactly as stated.

## Programmatic Acceptance

=== AC listing-empty-header-only ===
Intent: With no duplicate groups the listing writer emits only the header row.

from drydock_go import run_test

result = run_test("./internal/output", "TestWriteListingEmptyHeaderOnly")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC listing-empty-header-only ===

=== AC listing-empty-exact-bytes ===
Intent: With no duplicate groups the listing output is exactly the header followed by one newline.

from drydock_go import run_test

result = run_test("./internal/output", "TestWriteListingEmptyExactBytes")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC listing-empty-exact-bytes ===

=== AC listing-empty-nil-equals-empty ===
Intent: A nil group slice and an empty group slice produce identical listing output.

from drydock_go import run_test

result = run_test("./internal/output", "TestWriteListingEmptyNilEqualsEmpty")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC listing-empty-nil-equals-empty ===

## User Acceptance

- None.

## Guardrails

- The empty listing never writes a blank line or a data row.
