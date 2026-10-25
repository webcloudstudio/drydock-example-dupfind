# Sea Trials: DupFind

## Policy

| Consequence | On FAIL | On INCONCLUSIVE |
|---|---|---|
| blocks  | fail   | attest |
| scores  | score  | score  |
| attests | report | report |

## st-001: The supplied grader passes
Type: technical
Required: yes
Criterion: The completed dupfind shall make sh sources/run.sh exit zero; that script's exit status is the sole acceptance verdict.
Testability: deterministic
Consequence: blocks
Verification: proof
Command: ["sh", "sources/run.sh"]
Pattern: ubiquitous

## st-002: Test suite passes
Type: technical
Required: yes
Criterion: The application shall provide a POSIX-compatible bin/test.sh that exits zero when go vet and every automated test pass.
Testability: deterministic
Consequence: blocks
Verification: proof
Pattern: ubiquitous

## st-003: The executable builds with go build
Type: technical
Required: yes
Criterion: The application root shall be a Go module from which go build ./cmd/dupfind produces the dupfind executable using only the standard library.
Testability: deterministic
Consequence: blocks
Verification: proof
Pattern: ubiquitous

## st-004: Duplicate groups are exact and ordered
Type: behavioral
Required: yes
Criterion: When dupfind lists a directory, the application shall report every group of two or more regular files with identical size and SHA-256 digest, ordered by reclaimable bytes descending then digest ascending, with paths relative to the directory in byte order.
Testability: deterministic
Consequence: blocks
Verification: proof
Pattern: event

## st-005: Empty files and symbolic links are never reported
Type: behavioral
Required: yes
Criterion: The application shall ignore zero-byte files and shall neither follow nor report symbolic links.
Testability: deterministic
Consequence: blocks
Verification: proof
Pattern: ubiquitous

## st-006: The summary totals the listing
Type: behavioral
Required: yes
Criterion: When --summary is given, the application shall write the group count, the duplicate file count, and the reclaimable bytes of the listing that the same options would produce.
Testability: deterministic
Consequence: blocks
Verification: proof
Pattern: event

## st-007: Errors map to exit codes
Type: behavioral
Required: yes
Criterion: If the directory is missing or is not a directory, then the application shall exit 1; if the arguments are malformed, then the application shall exit 2; in both cases it shall write nothing to standard output.
Testability: deterministic
Consequence: blocks
Verification: proof
Pattern: unwanted

## st-008: Files are never modified
Type: technical
Required: yes
Criterion: The application shall only read the files it examines and shall never delete, move, or write any file under the directory.
Testability: judgeable
Consequence: scores
Verification: llm
Pattern: ubiquitous
