# FEATURE: Listing Order

| Field       | Value |
|-------------|-------|
| Version     | 20261024 V1 |
| Description | Orders duplicate groups and their paths, numbers the groups from 1, and renders paths relative to the examined directory. |
| Depends On  | FEATURE-ListingCSV.md |
| Provides    | output.Order |
| Consumes    | output.WriteListing |

## Purpose

Establish the deterministic order and the rendered path form that the listing and summary both
depend on. Groups are ranked by reclaimable bytes, then by digest. Paths are rendered relative to the
examined directory, so the listing is stable across runs and machines.

## Behavior

- Reclaimable bytes for a group equal `size × (files − 1)`.
- Groups are ordered by reclaimable bytes descending, then by digest ascending (lowercase hexadecimal, byte order).
- Groups are numbered from 1 in that order.
- Within a group, paths are ordered by byte value ascending.
- Paths are relative to the examined directory, use `/` as the separator, and carry no leading `./`.

## Interface

- `output.Order(groups []scan.Group, root string) []output.OrderedGroup`

## Test Contract

The `internal/output` package carries one Go test per criterion below, named exactly as stated.

## Programmatic Acceptance

=== AC listing-order-reclaimable-desc ===
Intent: Groups are ordered by reclaimable bytes descending.

from drydock_go import run_test

result = run_test("./internal/output", "TestOrderReclaimableDescending")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC listing-order-reclaimable-desc ===

=== AC listing-order-digest-tiebreak ===
Intent: Groups with equal reclaimable bytes are ordered by digest ascending.

from drydock_go import run_test

result = run_test("./internal/output", "TestOrderDigestAscendingTiebreak")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC listing-order-digest-tiebreak ===

=== AC listing-order-numbering ===
Intent: Groups are numbered from 1 in their final order.

from drydock_go import run_test

result = run_test("./internal/output", "TestOrderNumbersFromOne")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC listing-order-numbering ===

=== AC listing-order-path-bytes ===
Intent: Paths within a group are ordered by byte value ascending, including a path containing a comma.

from drydock_go import run_test

result = run_test("./internal/output", "TestOrderPathsByteOrder")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC listing-order-path-bytes ===

=== AC listing-order-relative-paths ===
Intent: Paths are relative to the examined directory with slash separators and no leading ./.

from drydock_go import run_test

result = run_test("./internal/output", "TestOrderRelativeSlashPaths")
print(result.status, result.stdout, result.stderr)
assert result.passed, (result.status, result.stdout, result.stderr)
=== END AC listing-order-relative-paths ===

## User Acceptance

- None.

## Guardrails

- Ordering is a pure function of the groups and the root; it reads no file contents.
- Ordering never changes group membership or the set of paths.
