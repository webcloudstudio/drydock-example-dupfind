# Evidence: Block 2 · Service (block-2)

- block type: block
- date: 2026-10-25
- resulting state: closed/verified
- story points (combined assembled cost): 27304
- execution id: 20261025.175058.104Z-fca82587

## Stories built
- Parse --summary, --min-size, and exactly one directory argument; reject malformed usage with exit 2. (arg-parsing) [story]
- Provide the cmd/dupfind main package that builds with go build and delegates to cli.Run. (module-build) [story]
- Exit 1 with empty stdout when the directory is missing or is not a directory. (directory-validation) [story]
- Walk the directory recursively and yield eligible regular files at or above the minimum size. (walk) [story]
- Group candidate files by exact size and keep only sizes shared by two or more files. (size-grouping) [story]
- Hash same-size files with SHA-256 and form duplicate groups, skipping unreadable files with a warning. (digest-grouping) [story]
- Compute reclaimable bytes, order groups, number them from 1, and order paths within each group. (group-ordering) [story]
- Write the CSV duplicate listing with header group,size,sha256,path using encoding/csv. (listing) [story]
- Write the one-row summary of groups, duplicate files, and reclaimable bytes with --summary. (summary) [story]
- Provide POSIX bin/test.sh and a Go test suite covering every rule, output form, and exit code. (test-gate) [story]

## Acceptance tooling authorization
- FEATURE-TestGate.md#test-script-exits-zero: executable=sh; scope=test; authorization=existing Target environment
- FEATURE-TestGate.md#test-script-exits-zero: executable=go; scope=test; authorization=existing Target environment
- FEATURE-TestGate.md#test-script-posix-syntax: executable=sh; scope=test; authorization=existing Target environment
- FEATURE-TestGate.md#test-script-fails-on-failing-test: executable=sh; scope=test; authorization=existing Target environment
- FEATURE-TestGate.md#test-script-fails-on-failing-test: executable=go; scope=test; authorization=existing Target environment
- FEATURE-TestGate.md#test-script-independent-of-cwd: executable=sh; scope=test; authorization=existing Target environment
- FEATURE-TestGate.md#test-script-independent-of-cwd: executable=go; scope=test; authorization=existing Target environment
- FEATURE-TestGate.md#suite-has-tests-for-cli: executable=go; scope=test; authorization=existing Target environment
- FEATURE-TestGate.md#built-binary-exit-codes: executable=go; scope=test; authorization=existing Target environment

## Stacked context
- stack: go_compact.md (SP 2105)
- compass: COMPASS.md (SP 982)
- implements: FEATURE-ArgParsing.md (SP 1450)
- context: ARCHITECTURE.md (SP 2007)
- implements: FEATURE-ModuleBuild.md (SP 1031)
- implements: FEATURE-DirectoryValidation.md (SP 1270)
- implements: FEATURE-Walk.md (SP 1844)
- implements: FEATURE-SizeGrouping.md (SP 1327)
- implements: FEATURE-DigestGrouping.md (SP 2579)
- implements: FEATURE-GroupOrdering.md (SP 2014)
- implements: FEATURE-Listing.md (SP 2440)
- implements: FEATURE-Summary.md (SP 1953)
- implements: FEATURE-TestGate.md (SP 3219)
- stack: common_compact.md (SP 1653)

## Build directory changes
- .gitignore
- bin/test.sh
- data/logs/test_2026-10-25_135414.log
- data/logs/test_2026-10-25_135420.log
- internal/cli/cli.go
- internal/cli/cli_test.go
- internal/dupes/dupes.go
- internal/dupes/dupes_test.go
- internal/report/report.go
- internal/report/report_test.go
- internal/walk/walk.go
- internal/walk/walk_test.go

## Pre-build acceptance observation
- RED: argparse-usage-errors-exit-2 (FEATURE-ArgParsing.md)
  intent: Every malformed command line exits 2 with empty standard output.
  return code: 1
  stdout:
    [] 0 '' ''
  stderr:
    --- drydock: values at failure ---
      args = []
      code = 0
    --- drydock: end values ---
    Traceback (most recent call last):
      File "argparse-usage-errors-exit-2.py", line 28, in <module>
        assert code == 2, args
               ^^^^^^^^^
    AssertionError: []
- RED: argparse-diagnostic-on-stderr (FEATURE-ArgParsing.md)
  intent: A usage error writes exactly one diagnostic line to standard error.
  return code: 1
  stdout:
    0 '' ''
  stderr:
    --- drydock: values at failure ---
      code = 0
    --- drydock: end values ---
    Traceback (most recent call last):
      File "argparse-diagnostic-on-stderr.py", line 5, in <module>
        assert code == 2
               ^^^^^^^^^
    AssertionError
- GREEN (prepassed): argparse-valid-forms-not-usage-errors (FEATURE-ArgParsing.md)
  intent: Well-formed command lines are not rejected as usage errors.
  return code: 0
  stdout:
    ['/tmp/tmpfnakgsrk'] 0 '' ''
    ['--summary', '/tmp/tmpfnakgsrk'] 0 '' ''
    ['--min-size', '5', '/tmp/tmpfnakgsrk'] 0 '' ''
    ['--summary', '--min-size', '5', '/tmp/tmpfnakgsrk'] 0 '' ''
    ['--min-size', '5', '--summary', '/tmp/tmpfnakgsrk'] 0 '' ''
- GREEN (prepassed): modbuild-go-build (FEATURE-ModuleBuild.md)
  intent: go build ./cmd/dupfind succeeds and produces an executable.
  return code: 0
- GREEN (prepassed): modbuild-standard-library-only (FEATURE-ModuleBuild.md)
  intent: The executable's dependency graph contains only standard library and in-module packages.
  return code: 0
  stdout:
    dupfind/internal/cli
    dupfind/cmd/dupfind
- RED: modbuild-executable-runs (FEATURE-ModuleBuild.md)
  intent: The built executable runs and rejects an empty command line with exit 2.
  return code: 1
  stdout:
    0 '' ''
  stderr:
    --- drydock: values at failure ---
      code = 0
    --- drydock: end values ---
    Traceback (most recent call last):
      File "modbuild-executable-runs.py", line 5, in <module>
        assert code == 2
               ^^^^^^^^^
    AssertionError
- RED: dirval-missing-directory (FEATURE-DirectoryValidation.md)
  intent: A nonexistent directory exits 1 with empty standard output and one stderr line.
  return code: 1
  stdout:
    ['/tmp/tmpckhpptef/no-such-directory'] 0 '' ''
  stderr:
    --- drydock: values at failure ---
      args = ['/tmp/tmpckhpptef/no-such-directory']
      code = 0
    --- drydock: end values ---
    Traceback (most recent call last):
      File "dirval-missing-directory.py", line 10, in <module>
        assert code == 1, args
               ^^^^^^^^^
    AssertionError: ['/tmp/tmpckhpptef/no-such-directory']
- RED: dirval-not-a-directory (FEATURE-DirectoryValidation.md)
  intent: A regular file given as the directory exits 1 with empty standard output.
  return code: 1
  stdout:
    ['/tmp/tmpyhexhc9w/plain.txt'] 0 '' ''
  stderr:
    --- drydock: values at failure ---
      args = ['/tmp/tmpyhexhc9w/plain.txt']
      code = 0
    --- drydock: end values ---
    Traceback (most recent call last):
      File "dirval-not-a-directory.py", line 12, in <module>
        assert code == 1, args
               ^^^^^^^^^
    AssertionError: ['/tmp/tmpyhexhc9w/plain.txt']
- RED: dirval-usage-precedes-validation (FEATURE-DirectoryValidation.md)
  intent: A usage error on a nonexistent directory still exits 2, not 1.
  return code: 1
  stdout:
    0 '' ''
  stderr:
    --- drydock: values at failure ---
      code = 0
    --- drydock: end values ---
    Traceback (most recent call last):
      File "dirval-usage-precedes-validation.py", line 9, in <module>
        assert code == 2
               ^^^^^^^^^
    AssertionError
- GREEN (prepassed): dirval-existing-directory-succeeds (FEATURE-DirectoryValidation.md)
  intent: An existing empty directory is accepted and the run exits 0.
  return code: 0
  stdout:
    0 '' ''
- RED: walk-rules (FEATURE-Walk.md)
  intent: Walk returns hidden and nested regular files, omits symbolic links and zero-byte files, and honors the minimum size.
  return code: 1
  stdout:
    FAIL	dupfind/internal/walk [build failed]
    FAIL
  stderr:
    # dupfind/internal/walk_test [dupfind/internal/walk.test]
    internal/walk/zz_acceptance_test.go:28:19: undefined: walk.Files
    
    --- drydock: values at failure ---
      result = CompletedProcess(args=['go', 'test', './internal/walk', '-run', 'TestAcceptanceWalk', '-count=1', '-v'], returncode=1, stdout='FAIL\tdupfind/internal/walk [build failed]\nFAIL\n', stderr='# dupfind/internal/walk_test [dupfind/internal/walk.test]\ninternal/walk/zz_acceptance_test.go:28:19: undefined: walk.Files\n')
    --- drydock: end values ---
    Traceback (most recent call last):
      File "walk-rules.py", line 97, in <module>
        assert result.returncode == 0
               ^^^^^^^^^^^^^^^^^^^^^^
    AssertionError
- GREEN (prepassed): walk-vet (FEATURE-Walk.md)
  intent: The walk package compiles and passes go vet.
  return code: 0
- RED: size-unique-and-shared (FEATURE-SizeGrouping.md)
  intent: Files with a unique size are never reported and a shared size with identical contents is reported.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      paths = []
    --- drydock: end values ---
    Traceback (most recent call last):
      File "size-unique-and-shared.py", line 24, in <module>
        assert paths == ["d1.bin", "sub/d2.bin"]
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
    AssertionError
- RED: size-same-size-different-content (FEATURE-SizeGrouping.md)
  intent: Files that share a size but differ in content are not reported as duplicates.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      rows = []
    --- drydock: end values ---
    Traceback (most recent call last):
      File "size-same-size-different-content.py", line 20, in <module>
        assert len(rows) == 1
               ^^^^^^^^^^^^^^
    AssertionError
- UNVERIFIED: size-unit-tests (FEATURE-SizeGrouping.md)
  intent: The project's own tests for the size bucketing pass and actually execute.
  return code: 1
  error: unverified acceptance: not-run
  stderr:
    DRYDOCK_NATIVE_RESULT={"status": "not-run"}
    --- drydock: values at failure ---
      result = GoResult(passed=False, status='not-run', return_code=0, stdout='{"Time":"2026-10-25T13:49:13.828258656-04:00","Action":"start","Package":"dupfind/internal/dupes"}\n{"Time":"2026-10-25T13:49:13.828313922-04:00","Action":"output","Package":"dupfind/internal/dupes","Output":"?   \\tdupfind/internal/dupes\\t[no test files]\\n"}\n{"Time":"2026-10-25T13:49:13.828318473-04:00","Action":"skip","Package":"… (+61 chars)
    --- drydock: end values ---
    Traceback (most recent call last):
      File "size-unit-tests.py", line 4, in <module>
        assert result.passed, (result.status, result.stdout, result.stderr)
               ^^^^^^^^^^^^^
    AssertionError: ('not-run', '{"Time":"2026-10-25T13:49:13.828258656-04:00","Action":"start","Package":"dupfind/internal/dupes"}\n{"Time":"2026-10-25T13:49:13.828313922-04:00","Action":"output","Package":"dupfind/internal/dupes","Output":"?   \\tdupfind/internal/dupes\\t[no test files]\\n"}\n{"Time":"2026-10-25T13:49:13.828318473-04:00","Action":"skip","Package":"dupfind/internal/dupes","Elapsed":0}\n', '')
- RED: digest-columns-match-contents (FEATURE-DigestGrouping.md)
  intent: Reported size and sha256 are the byte length and lowercase SHA-256 of the fixture contents.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      rows = []
    --- drydock: end values ---
    Traceback (most recent call last):
      File "digest-columns-match-contents.py", line 23, in <module>
        assert len(rows) == 3
               ^^^^^^^^^^^^^^
    AssertionError
- RED: digest-splits-same-size-groups (FEATURE-DigestGrouping.md)
  intent: Same-size files with different digests form separate groups and a lone digest is dropped.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      by_group = {}
    --- drydock: end values ---
    Traceback (most recent call last):
      File "digest-splits-same-size-groups.py", line 29, in <module>
        assert len(by_group) == 2
               ^^^^^^^^^^^^^^^^^^
    AssertionError
- RED: digest-unreadable-skipped (FEATURE-DigestGrouping.md)
  intent: An unreadable duplicate candidate is skipped without changing the exit code (skipped when permissions cannot be enforced).
  return code: 1
  stderr:
    --- drydock: values at failure ---
      paths = []
    --- drydock: end values ---
    Traceback (most recent call last):
      File "digest-unreadable-skipped.py", line 31, in <module>
        assert paths == ["ok1.txt", "ok2.txt"]
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
    AssertionError
- GREEN (prepassed): digest-tree-unmodified (FEATURE-DigestGrouping.md)
  intent: Running the tool leaves every file's bytes and the directory listing unchanged.
  return code: 0
- UNVERIFIED: digest-unit-tests (FEATURE-DigestGrouping.md)
  intent: The project's own digest grouping tests pass and actually execute.
  return code: 1
  error: unverified acceptance: not-run
  stderr:
    DRYDOCK_NATIVE_RESULT={"status": "not-run"}
    --- drydock: values at failure ---
      result = GoResult(passed=False, status='not-run', return_code=0, stdout='{"Time":"2026-10-25T13:49:26.194927537-04:00","Action":"start","Package":"dupfind/internal/dupes"}\n{"Time":"2026-10-25T13:49:26.195115533-04:00","Action":"output","Package":"dupfind/internal/dupes","Output":"?   \\tdupfind/internal/dupes\\t[no test files]\\n"}\n{"Time":"2026-10-25T13:49:26.195120766-04:00","Action":"skip","Package":"… (+61 chars)
    --- drydock: end values ---
    Traceback (most recent call last):
      File "digest-unit-tests.py", line 4, in <module>
        assert result.passed, (result.status, result.stdout, result.stderr)
               ^^^^^^^^^^^^^
    AssertionError: ('not-run', '{"Time":"2026-10-25T13:49:26.194927537-04:00","Action":"start","Package":"dupfind/internal/dupes"}\n{"Time":"2026-10-25T13:49:26.195115533-04:00","Action":"output","Package":"dupfind/internal/dupes","Output":"?   \\tdupfind/internal/dupes\\t[no test files]\\n"}\n{"Time":"2026-10-25T13:49:26.195120766-04:00","Action":"skip","Package":"dupfind/internal/dupes","Elapsed":0}\n', '')
- RED: order-by-reclaimable-then-digest (FEATURE-GroupOrdering.md)
  intent: Groups appear by reclaimable bytes descending then digest ascending and are numbered from 1.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      seen = []
      expected = [(12, '3608bca1e44ea6c4d268eb6db02260269892c0b42b86bbf1e77a6fa16c3c9282', 3, ['s/1', 's/2', 's/3', 's/4', 's/5']), (10, '84d89877f0d4041efb6bf91a16f0248f2fd573e6af05c19f96bedb9f882f7882', 10, ['big/a', 'big/b']), (6, '1367edb55bf2b6b8aa4f874a5df16ed409642f43c05a0e06dbace6f18a95e61f', 6, ['t1/a', 't1/b']), (6, '58f55ac2c0a64469a13b9abd3104e6c1749e258e1f8fccdc7c3ecb370cce704e', 6, ['t2/a', 't2/b']),… (+94 chars)
    --- drydock: end values ---
    Traceback (most recent call last):
      File "order-by-reclaimable-then-digest.py", line 39, in <module>
        assert [int(s[0]) for s in seen] == list(range(1, len(expected) + 1))
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
    AssertionError
- RED: order-paths-byte-ascending (FEATURE-GroupOrdering.md)
  intent: Paths within a group are in byte order and written relative with slash separators and no leading dot-slash.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      got = []
      paths = ['b.txt', 'B.txt', 'a b.txt', 'a/z.txt', '.hid/x.txt', 'a.txt']
    --- drydock: end values ---
    Traceback (most recent call last):
      File "order-paths-byte-ascending.py", line 22, in <module>
        assert got == sorted(paths, key=lambda s: s.encode("utf-8"))
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
    AssertionError
- RED: order-reclaimable-beats-size (FEATURE-GroupOrdering.md)
  intent: A smaller file with more copies precedes a larger file with fewer copies when it reclaims more bytes.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      rows = []
      small = b'sss'
    --- drydock: end values ---
    Traceback (most recent call last):
      File "order-reclaimable-beats-size.py", line 28, in <module>
        assert rows[0][2] == hashlib.sha256(small).hexdigest()
               ~~~~^^^
    IndexError: list index out of range
- UNVERIFIED: order-unit-tests (FEATURE-GroupOrdering.md)
  intent: The project's own ordering test passes and actually executes.
  return code: 1
  error: unverified acceptance: not-run
  stderr:
    DRYDOCK_NATIVE_RESULT={"status": "not-run"}
    --- drydock: values at failure ---
      result = GoResult(passed=False, status='not-run', return_code=0, stdout='{"Time":"2026-10-25T13:49:33.89061549-04:00","Action":"start","Package":"dupfind/internal/dupes"}\n{"Time":"2026-10-25T13:49:33.89066849-04:00","Action":"output","Package":"dupfind/internal/dupes","Output":"?   \\tdupfind/internal/dupes\\t[no test files]\\n"}\n{"Time":"2026-10-25T13:49:33.89067207-04:00","Action":"skip","Package":"dup… (+58 chars)
    --- drydock: end values ---
    Traceback (most recent call last):
      File "order-unit-tests.py", line 4, in <module>
        assert result.passed, (result.status, result.stdout, result.stderr)
               ^^^^^^^^^^^^^
    AssertionError: ('not-run', '{"Time":"2026-10-25T13:49:33.89061549-04:00","Action":"start","Package":"dupfind/internal/dupes"}\n{"Time":"2026-10-25T13:49:33.89066849-04:00","Action":"output","Package":"dupfind/internal/dupes","Output":"?   \\tdupfind/internal/dupes\\t[no test files]\\n"}\n{"Time":"2026-10-25T13:49:33.89067207-04:00","Action":"skip","Package":"dupfind/internal/dupes","Elapsed":0}\n', '')
- RED: listing-full-output (FEATURE-Listing.md)
  intent: The entire standard output equals the CSV the standard csv writer produces from the fixture-derived rows.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      out = ''
      expected = 'group,size,sha256,path\n1,12,398e46cf2cbf90b489c4a699321be3b124c2637deaf8d5a5ae11aac7459688c4,a/one.txt\n1,12,398e46cf2cbf90b489c4a699321be3b124c2637deaf8d5a5ae11aac7459688c4,b/two.txt\n1,12,398e46cf2cbf90b489c4a699321be3b124c2637deaf8d5a5ae11aac7459688c4,c/deep/three.txt\n2,5,f2c82decdd7181cf98945929a62598db7e6b477e11f6e0eb0ae97020eff151ad,.hidden/.env\n2,5,f2c82decdd7181cf98945929a62598db7e6b47… (+37 chars)
    --- drydock: end values ---
    Traceback (most recent call last):
      File "listing-full-output.py", line 39, in <module>
        assert out == expected
               ^^^^^^^^^^^^^^^
    AssertionError
- RED: listing-csv-quoting (FEATURE-Listing.md)
  intent: Paths containing a comma or a double quote survive a CSV round trip unchanged.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      paths = ['song,live.mp3', 'say "hi".txt', 'plain.txt']
      rows = []
    --- drydock: end values ---
    Traceback (most recent call last):
      File "listing-csv-quoting.py", line 23, in <module>
        assert sorted(r[3] for r in rows[1:]) == sorted(paths)
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
    AssertionError
- RED: listing-header-only-when-no-duplicates (FEATURE-Listing.md)
  intent: A directory without duplicates (or with only zero-byte files) yields the header row only, with exit 0.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      out = ''
    --- drydock: end values ---
    Traceback (most recent call last):
      File "listing-header-only-when-no-duplicates.py", line 20, in <module>
        assert out.splitlines() == ["group,size,sha256,path"]
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
    AssertionError
- RED: listing-min-size (FEATURE-Listing.md)
  intent: Files smaller than --min-size are excluded from the listing while larger duplicates remain.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      rows = []
    --- drydock: end values ---
    Traceback (most recent call last):
      File "listing-min-size.py", line 25, in <module>
        assert sorted(r[3] for r in rows) == ["l1", "l2"]
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
    AssertionError
- RED: listing-symlinks-and-hidden (FEATURE-Listing.md)
  intent: Hidden files are listed while symbolic links are neither followed nor reported.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      rows = []
    --- drydock: end values ---
    Traceback (most recent call last):
      File "listing-symlinks-and-hidden.py", line 23, in <module>
        assert sorted(r[3] for r in rows) == [".hidden/.env", "env.bak"]
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
    AssertionError
- RED: summary-totals-match-fixture (FEATURE-Summary.md)
  intent: The summary row equals totals computed from the fixture's known duplicate structure.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      rows = []
    --- drydock: end values ---
    Traceback (most recent call last):
      File "summary-totals-match-fixture.py", line 30, in <module>
        assert len(rows) == 2
               ^^^^^^^^^^^^^^
    AssertionError
- RED: summary-no-duplicates (FEATURE-Summary.md)
  intent: With no duplicates the summary row holds three zeros and exit status is 0.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      rows = []
    --- drydock: end values ---
    Traceback (most recent call last):
      File "summary-no-duplicates.py", line 20, in <module>
        assert len(rows) == 2
               ^^^^^^^^^^^^^^
    AssertionError
- RED: summary-agrees-with-listing (FEATURE-Summary.md)
  intent: The summary totals equal the totals derived from the listing produced by the same options, including --min-size.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      srows = []
    --- drydock: end values ---
    Traceback (most recent call last):
      File "summary-agrees-with-listing.py", line 38, in <module>
        assert len(srows) == 2
               ^^^^^^^^^^^^^^^
    AssertionError
- UNVERIFIED: summary-unit-tests (FEATURE-Summary.md)
  intent: The project's own summary tests pass and actually execute.
  return code: 1
  error: unverified acceptance: not-run
  stderr:
    DRYDOCK_NATIVE_RESULT={"status": "not-run"}
    --- drydock: values at failure ---
      result = GoResult(passed=False, status='not-run', return_code=0, stdout='{"Time":"2026-10-25T13:49:53.488952641-04:00","Action":"start","Package":"dupfind/internal/report"}\n{"Time":"2026-10-25T13:49:53.489019309-04:00","Action":"output","Package":"dupfind/internal/report","Output":"?   \\tdupfind/internal/report\\t[no test files]\\n"}\n{"Time":"2026-10-25T13:49:53.48902687-04:00","Action":"skip","Package"… (+64 chars)
    --- drydock: end values ---
    Traceback (most recent call last):
      File "summary-unit-tests.py", line 4, in <module>
        assert result.passed, (result.status, result.stdout, result.stderr)
               ^^^^^^^^^^^^^
    AssertionError: ('not-run', '{"Time":"2026-10-25T13:49:53.488952641-04:00","Action":"start","Package":"dupfind/internal/report"}\n{"Time":"2026-10-25T13:49:53.489019309-04:00","Action":"output","Package":"dupfind/internal/report","Output":"?   \\tdupfind/internal/report\\t[no test files]\\n"}\n{"Time":"2026-10-25T13:49:53.48902687-04:00","Action":"skip","Package":"dupfind/internal/report","Elapsed":0}\n', '')
- RED: test-script-exits-zero (FEATURE-TestGate.md)
  intent: bin/test.sh runs go vet and go test on the finished application and exits zero.
  return code: 1
  stderr:
    sh: 0: cannot open bin/test.sh: No such file
    
    --- drydock: values at failure ---
      result = CompletedProcess(args=['sh', 'bin/test.sh'], returncode=2, stdout='', stderr='sh: 0: cannot open bin/test.sh: No such file\n')
    --- drydock: end values ---
    Traceback (most recent call last):
      File "test-script-exits-zero.py", line 7, in <module>
        assert result.returncode == 0, result.returncode
               ^^^^^^^^^^^^^^^^^^^^^^
    AssertionError: 2
- RED: test-script-posix-syntax (FEATURE-TestGate.md)
  intent: bin/test.sh exists and parses as POSIX sh without executing.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      script = 'bin/test.sh'
    --- drydock: end values ---
    Traceback (most recent call last):
      File "test-script-posix-syntax.py", line 6, in <module>
        assert os.path.isfile(script)
               ^^^^^^^^^^^^^^^^^^^^^^
    AssertionError
- GREEN (prepassed): test-script-fails-on-failing-test (FEATURE-TestGate.md)
  intent: bin/test.sh exits non-zero when any Go test fails, proven on a temporary copy of the application.
  return code: 0
  stderr:
    sh: 0: cannot open bin/test.sh: No such file
- RED: test-script-independent-of-cwd (FEATURE-TestGate.md)
  intent: bin/test.sh locates the application root itself, so running it from another directory still exits zero.
  return code: 1
  stderr:
    sh: 0: cannot open /mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261025.174005/build/DupFind/bin/test.sh: No such file
    
    --- drydock: values at failure ---
      result = CompletedProcess(args=['sh', '/mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261025.174005/build/DupFind/bin/test.sh'], returncode=2, stdout='', stderr='sh: 0: cannot open /mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261025.174005/build/DupFind/bin/test.sh: No such file\n')
    --- drydock: end values ---
    Traceback (most recent call last):
      File "test-script-independent-of-cwd.py", line 11, in <module>
        assert result.returncode == 0, result.returncode
               ^^^^^^^^^^^^^^^^^^^^^^
    AssertionError: 2
- GREEN (prepassed): suite-has-tests-for-cli (FEATURE-TestGate.md)
  intent: The Go test suite contains test files for the cli package and they pass when run alone.
  return code: 0
  stdout:
    ok  	dupfind/internal/cli	0.005s
- RED: built-binary-exit-codes (FEATURE-TestGate.md)
  intent: The built executable exits 0 for a directory, 1 for a missing path, and 2 for each malformed argument form the suite covers.
  return code: 1
  stderr:
    --- drydock: values at failure ---
      missing_code = 0
      missing_out = ''
      missing_err = ''
    --- drydock: end values ---
    Traceback (most recent call last):
      File "built-binary-exit-codes.py", line 17, in <module>
        assert missing_code == 1, (missing_code, missing_out, missing_err)
               ^^^^^^^^^^^^^^^^^
    AssertionError: (0, '', '')

## Post-build programmatic acceptance
- PASS: argparse-usage-errors-exit-2 (FEATURE-ArgParsing.md)
  intent: Every malformed command line exits 2 with empty standard output.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    [] 2 '' 'dupfind: expected exactly one directory argument, got 0; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
    ['--summary'] 2 '' 'dupfind: expected exactly one directory argument, got 0; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
    ['--min-size', '5'] 2 '' 'dupfind: expected exactly one directory argument, got 0; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
    ['/tmp/tmp7uvhvtd0', '/tmp/tmp7uvhvtd0/other'] 2 '' 'dupfind: expected exactly one directory argument, got 2; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
    ['--frobnicate', '/tmp/tmp7uvhvtd0'] 2 '' 'dupfind: flag provided but not defined: -frobnicate; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
    ['-h', '/tmp/tmp7uvhvtd0'] 2 '' 'dupfind: flag: help requested; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
    ['--min-size', 'abc', '/tmp/tmp7uvhvtd0'] 2 '' 'dupfind: --min-size must be a positive decimal integer, got "abc"; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
    ['--min-size', '0', '/tmp/tmp7uvhvtd0'] 2 '' 'dupfind: --min-size must be a positive decimal integer, got "0"; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
    ['--min-size', '-5', '/tmp/tmp7uvhvtd0'] 2 '' 'dupfind: --min-size must be a positive decimal integer, got "-5"; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
    ['--min-size', '+5', '/tmp/tmp7uvhvtd0'] 2 '' 'dupfind: --min-size must be a positive decimal integer, got "+5"; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
    ['--min-size', '0x10', '/tmp/tmp7uvhvtd0'] 2 '' 'dupfind: --min-size must be a positive decimal integer, got "0x10"; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
    ['--min-size', '1e3', '/tmp/tmp7uvhvtd0'] 2 '' 'dupfind: --min-size must be a positive decimal integer, got "1e3"; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
    ['--min-size', '99999999999999999999999', '/tmp/tmp7uvhvtd0'] 2 '' 'dupfind: --min-size must be a positive decimal integer, got "99999999999999999999999"; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
    ['--min-size', '/tmp/tmp7uvhvtd0'] 2 '' 'dupfind: expected exactly one directory argument, got 0; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
    ['/tmp/tmp7uvhvtd0', '--summary'] 2 '' 'dupfind: expected exactly one directory argument, got 2; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
- PASS: argparse-diagnostic-on-stderr (FEATURE-ArgParsing.md)
  intent: A usage error writes exactly one diagnostic line to standard error.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    2 '' 'dupfind: --min-size must be a positive decimal integer, got "0"; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
- PASS: argparse-valid-forms-not-usage-errors (FEATURE-ArgParsing.md)
  intent: Well-formed command lines are not rejected as usage errors.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    ['/tmp/tmpvsf4xwjs'] 0 'group,size,sha256,path\n' ''
    ['--summary', '/tmp/tmpvsf4xwjs'] 0 'groups,duplicate_files,reclaimable_bytes\n0,0,0\n' ''
    ['--min-size', '5', '/tmp/tmpvsf4xwjs'] 0 'group,size,sha256,path\n' ''
    ['--summary', '--min-size', '5', '/tmp/tmpvsf4xwjs'] 0 'groups,duplicate_files,reclaimable_bytes\n0,0,0\n' ''
    ['--min-size', '5', '--summary', '/tmp/tmpvsf4xwjs'] 0 'groups,duplicate_files,reclaimable_bytes\n0,0,0\n' ''
- PASS: modbuild-go-build (FEATURE-ModuleBuild.md)
  intent: go build ./cmd/dupfind succeeds and produces an executable.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
- PASS: modbuild-standard-library-only (FEATURE-ModuleBuild.md)
  intent: The executable's dependency graph contains only standard library and in-module packages.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    dupfind/internal/walk
    dupfind/internal/dupes
    dupfind/internal/report
    dupfind/internal/cli
    dupfind/cmd/dupfind
- PASS: modbuild-executable-runs (FEATURE-ModuleBuild.md)
  intent: The built executable runs and rejects an empty command line with exit 2.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    2 '' 'dupfind: expected exactly one directory argument, got 0; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
- PASS: dirval-missing-directory (FEATURE-DirectoryValidation.md)
  intent: A nonexistent directory exits 1 with empty standard output and one stderr line.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    ['/tmp/tmpw26vh3u2/no-such-directory'] 1 '' 'dupfind: stat /tmp/tmpw26vh3u2/no-such-directory: no such file or directory\n'
    ['--summary', '/tmp/tmpw26vh3u2/no-such-directory'] 1 '' 'dupfind: stat /tmp/tmpw26vh3u2/no-such-directory: no such file or directory\n'
    ['--min-size', '3', '/tmp/tmpw26vh3u2/no-such-directory'] 1 '' 'dupfind: stat /tmp/tmpw26vh3u2/no-such-directory: no such file or directory\n'
- PASS: dirval-not-a-directory (FEATURE-DirectoryValidation.md)
  intent: A regular file given as the directory exits 1 with empty standard output.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    ['/tmp/tmppcjv2ht6/plain.txt'] 1 '' 'dupfind: /tmp/tmppcjv2ht6/plain.txt is not a directory\n'
    ['--summary', '/tmp/tmppcjv2ht6/plain.txt'] 1 '' 'dupfind: /tmp/tmppcjv2ht6/plain.txt is not a directory\n'
- PASS: dirval-usage-precedes-validation (FEATURE-DirectoryValidation.md)
  intent: A usage error on a nonexistent directory still exits 2, not 1.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    2 '' 'dupfind: --min-size must be a positive decimal integer, got "0"; usage: dupfind [--summary] [--min-size <bytes>] <directory>\n'
- PASS: dirval-existing-directory-succeeds (FEATURE-DirectoryValidation.md)
  intent: An existing empty directory is accepted and the run exits 0.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    0 'group,size,sha256,path\n' ''
- PASS: walk-rules (FEATURE-Walk.md)
  intent: Walk returns hidden and nested regular files, omits symbolic links and zero-byte files, and honors the minimum size.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    === RUN   TestAcceptanceWalkEligibility
    --- PASS: TestAcceptanceWalkEligibility (0.00s)
    === RUN   TestAcceptanceWalkMinSize
    --- PASS: TestAcceptanceWalkMinSize (0.00s)
    PASS
    ok  	dupfind/internal/walk	0.009s
- PASS: walk-vet (FEATURE-Walk.md)
  intent: The walk package compiles and passes go vet.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
- PASS: size-unique-and-shared (FEATURE-SizeGrouping.md)
  intent: Files with a unique size are never reported and a shared size with identical contents is reported.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    group,size,sha256,path
    1,6,c02d8e6211ef7dc5af42085cbb323a122b76905ac098467ea1f85465ec14429d,d1.bin
    1,6,c02d8e6211ef7dc5af42085cbb323a122b76905ac098467ea1f85465ec14429d,sub/d2.bin
- PASS: size-same-size-different-content (FEATURE-SizeGrouping.md)
  intent: Files that share a size but differ in content are not reported as duplicates.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    group,size,sha256,path
- PASS: size-unit-tests (FEATURE-SizeGrouping.md)
  intent: The project's own tests for the size bucketing pass and actually execute.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stderr:
    DRYDOCK_NATIVE_RESULT={"status": "pass"}
    DRYDOCK_NATIVE_RESULT={"status": "pass"}
- PASS: digest-columns-match-contents (FEATURE-DigestGrouping.md)
  intent: Reported size and sha256 are the byte length and lowercase SHA-256 of the fixture contents.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    group,size,sha256,path
    1,17,8a3c67892f82af22b58377b8e85bb41899053788d3aa3b81d2e6b4580e917c29,a/one.txt
    1,17,8a3c67892f82af22b58377b8e85bb41899053788d3aa3b81d2e6b4580e917c29,b/two.txt
    1,17,8a3c67892f82af22b58377b8e85bb41899053788d3aa3b81d2e6b4580e917c29,c/three.txt
- PASS: digest-splits-same-size-groups (FEATURE-DigestGrouping.md)
  intent: Same-size files with different digests form separate groups and a lone digest is dropped.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    group,size,sha256,path
    1,4,84e73dc50f2be9000ab2a87f8026c1f45e1fec954af502e9904031645b190d4f,y1
    1,4,84e73dc50f2be9000ab2a87f8026c1f45e1fec954af502e9904031645b190d4f,y2
    2,4,88d4266fd4e6338d13b845fcf289579d209c897823b9217da3e161936f031589,x1
    2,4,88d4266fd4e6338d13b845fcf289579d209c897823b9217da3e161936f031589,x2
- PASS: digest-unreadable-skipped (FEATURE-DigestGrouping.md)
  intent: An unreadable duplicate candidate is skipped without changing the exit code (skipped when permissions cannot be enforced).
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    group,size,sha256,path
    1,11,abb7f0ae43ba52cc56233a5ecb4dfa11765f26b1282a18346d811b6a85af19c1,ok1.txt
    1,11,abb7f0ae43ba52cc56233a5ecb4dfa11765f26b1282a18346d811b6a85af19c1,ok2.txt
    
    dupfind: warning: cannot read locked.txt: open /tmp/tmpftarbyvj/locked.txt: permission denied
- PASS: digest-tree-unmodified (FEATURE-DigestGrouping.md)
  intent: Running the tool leaves every file's bytes and the directory listing unchanged.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    group,size,sha256,path
    1,8,2b8425c4d20e743705f4787b4dda39344b4242bc8636228a00b7d65378aa7694,a/f1
    1,8,2b8425c4d20e743705f4787b4dda39344b4242bc8636228a00b7d65378aa7694,b/f2
- PASS: digest-unit-tests (FEATURE-DigestGrouping.md)
  intent: The project's own digest grouping tests pass and actually execute.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stderr:
    DRYDOCK_NATIVE_RESULT={"status": "pass"}
    DRYDOCK_NATIVE_RESULT={"status": "pass"}
- PASS: order-by-reclaimable-then-digest (FEATURE-GroupOrdering.md)
  intent: Groups appear by reclaimable bytes descending then digest ascending and are numbered from 1.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    group,size,sha256,path
    1,3,3608bca1e44ea6c4d268eb6db02260269892c0b42b86bbf1e77a6fa16c3c9282,s/1
    1,3,3608bca1e44ea6c4d268eb6db02260269892c0b42b86bbf1e77a6fa16c3c9282,s/2
    1,3,3608bca1e44ea6c4d268eb6db02260269892c0b42b86bbf1e77a6fa16c3c9282,s/3
    1,3,3608bca1e44ea6c4d268eb6db02260269892c0b42b86bbf1e77a6fa16c3c9282,s/4
    1,3,3608bca1e44ea6c4d268eb6db02260269892c0b42b86bbf1e77a6fa16c3c9282,s/5
    2,10,84d89877f0d4041efb6bf91a16f0248f2fd573e6af05c19f96bedb9f882f7882,big/a
    2,10,84d89877f0d4041efb6bf91a16f0248f2fd573e6af05c19f96bedb9f882f7882,big/b
    3,6,1367edb55bf2b6b8aa4f874a5df16ed409642f43c05a0e06dbace6f18a95e61f,t1/a
    3,6,1367edb55bf2b6b8aa4f874a5df16ed409642f43c05a0e06dbace6f18a95e61f,t1/b
    4,6,58f55ac2c0a64469a13b9abd3104e6c1749e258e1f8fccdc7c3ecb370cce704e,t2/a
    4,6,58f55ac2c0a64469a13b9abd3104e6c1749e258e1f8fccdc7c3ecb370cce704e,t2/b
    5,6,82a8129d6ecabf4c6eed5fdbdef815ca329ed0e4bb1da1ee742e984aa505ae2e,t3/a
    5,6,82a8129d6ecabf4c6eed5fdbdef815ca329ed0e4bb1da1ee742e984aa505ae2e,t3/b
- PASS: order-paths-byte-ascending (FEATURE-GroupOrdering.md)
  intent: Paths within a group are in byte order and written relative with slash separators and no leading dot-slash.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    group,size,sha256,path
    1,9,57eeb029403166b5608dbd105d998092ecf608c82b1f9a5930a34eb78a12b494,.hid/x.txt
    1,9,57eeb029403166b5608dbd105d998092ecf608c82b1f9a5930a34eb78a12b494,B.txt
    1,9,57eeb029403166b5608dbd105d998092ecf608c82b1f9a5930a34eb78a12b494,a b.txt
    1,9,57eeb029403166b5608dbd105d998092ecf608c82b1f9a5930a34eb78a12b494,a.txt
    1,9,57eeb029403166b5608dbd105d998092ecf608c82b1f9a5930a34eb78a12b494,a/z.txt
    1,9,57eeb029403166b5608dbd105d998092ecf608c82b1f9a5930a34eb78a12b494,b.txt
- PASS: order-reclaimable-beats-size (FEATURE-GroupOrdering.md)
  intent: A smaller file with more copies precedes a larger file with fewer copies when it reclaims more bytes.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    group,size,sha256,path
    1,3,a871c47a7f48a12b38a994e48a9659fab5d6376f3dbce37559bcb617efe8662d,s1
    1,3,a871c47a7f48a12b38a994e48a9659fab5d6376f3dbce37559bcb617efe8662d,s2
    1,3,a871c47a7f48a12b38a994e48a9659fab5d6376f3dbce37559bcb617efe8662d,s3
    1,3,a871c47a7f48a12b38a994e48a9659fab5d6376f3dbce37559bcb617efe8662d,s4
    1,3,a871c47a7f48a12b38a994e48a9659fab5d6376f3dbce37559bcb617efe8662d,s5
    2,10,0085eda52ea686d4ddb3b426dfa909c4fcd2159e3d98c1a7fa12cc3582ab79f7,l1
    2,10,0085eda52ea686d4ddb3b426dfa909c4fcd2159e3d98c1a7fa12cc3582ab79f7,l2
- PASS: order-unit-tests (FEATURE-GroupOrdering.md)
  intent: The project's own ordering test passes and actually executes.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stderr:
    DRYDOCK_NATIVE_RESULT={"status": "pass"}
- PASS: listing-full-output (FEATURE-Listing.md)
  intent: The entire standard output equals the CSV the standard csv writer produces from the fixture-derived rows.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    group,size,sha256,path
    1,12,398e46cf2cbf90b489c4a699321be3b124c2637deaf8d5a5ae11aac7459688c4,a/one.txt
    1,12,398e46cf2cbf90b489c4a699321be3b124c2637deaf8d5a5ae11aac7459688c4,b/two.txt
    1,12,398e46cf2cbf90b489c4a699321be3b124c2637deaf8d5a5ae11aac7459688c4,c/deep/three.txt
    2,5,f2c82decdd7181cf98945929a62598db7e6b477e11f6e0eb0ae97020eff151ad,.hidden/.env
    2,5,f2c82decdd7181cf98945929a62598db7e6b477e11f6e0eb0ae97020eff151ad,env.bak
- PASS: listing-csv-quoting (FEATURE-Listing.md)
  intent: Paths containing a comma or a double quote survive a CSV round trip unchanged.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    group,size,sha256,path
    1,9,5740c59f8dcbb2398fa201cb291e54be2941f9296b5b74cad76c719244035953,plain.txt
    1,9,5740c59f8dcbb2398fa201cb291e54be2941f9296b5b74cad76c719244035953,"say ""hi"".txt"
    1,9,5740c59f8dcbb2398fa201cb291e54be2941f9296b5b74cad76c719244035953,"song,live.mp3"
- PASS: listing-header-only-when-no-duplicates (FEATURE-Listing.md)
  intent: A directory without duplicates (or with only zero-byte files) yields the header row only, with exit 0.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    group,size,sha256,path
- PASS: listing-min-size (FEATURE-Listing.md)
  intent: Files smaller than --min-size are excluded from the listing while larger duplicates remain.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    group,size,sha256,path
    1,50,9bf8af326de6023135a2268f61a927f09b1856d7ba2140a1ef4e2d1ebf0282a7,l1
    1,50,9bf8af326de6023135a2268f61a927f09b1856d7ba2140a1ef4e2d1ebf0282a7,l2
- PASS: listing-symlinks-and-hidden (FEATURE-Listing.md)
  intent: Hidden files are listed while symbolic links are neither followed nor reported.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    group,size,sha256,path
    1,6,929c28f541260cd8690babe97ae61472829cdf36e059c7b14e6b963f0f9ee77b,.hidden/.env
    1,6,929c28f541260cd8690babe97ae61472829cdf36e059c7b14e6b963f0f9ee77b,env.bak
- PASS: summary-totals-match-fixture (FEATURE-Summary.md)
  intent: The summary row equals totals computed from the fixture's known duplicate structure.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    groups,duplicate_files,reclaimable_bytes
    2,3,43
- PASS: summary-no-duplicates (FEATURE-Summary.md)
  intent: With no duplicates the summary row holds three zeros and exit status is 0.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    groups,duplicate_files,reclaimable_bytes
    0,0,0
- PASS: summary-agrees-with-listing (FEATURE-Summary.md)
  intent: The summary totals equal the totals derived from the listing produced by the same options, including --min-size.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    groups,duplicate_files,reclaimable_bytes
    3,4,264
    
    
    groups,duplicate_files,reclaimable_bytes
    2,3,260
    
    
    groups,duplicate_files,reclaimable_bytes
    1,1,200
    
    
    groups,duplicate_files,reclaimable_bytes
    0,0,0
- PASS: summary-unit-tests (FEATURE-Summary.md)
  intent: The project's own summary tests pass and actually execute.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stderr:
    DRYDOCK_NATIVE_RESULT={"status": "pass"}
- PASS: test-script-exits-zero (FEATURE-TestGate.md)
  intent: bin/test.sh runs go vet and go test on the finished application and exits zero.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    timestamp: 2026-10-25_135534
    arguments: 
    working directory: /mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261025.174005/build/DupFind
    ?   	dupfind/cmd/dupfind	[no test files]
    ok  	dupfind/internal/cli	0.017s
    ok  	dupfind/internal/dupes	0.009s
    ok  	dupfind/internal/report	0.007s
    ok  	dupfind/internal/walk	0.011s
    completed: exit 0
- PASS: test-script-posix-syntax (FEATURE-TestGate.md)
  intent: bin/test.sh exists and parses as POSIX sh without executing.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
- PASS: test-script-fails-on-failing-test (FEATURE-TestGate.md)
  intent: bin/test.sh exits non-zero when any Go test fails, proven on a temporary copy of the application.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    timestamp: 2026-10-25_135538
    arguments: 
    working directory: /tmp/tmpll49ulh6/app
    ?   	dupfind/cmd/dupfind	[no test files]
    --- FAIL: TestForcedFailure (0.00s)
        zz_forced_failure_test.go:6: forced
    FAIL
    FAIL	dupfind/internal/cli	0.014s
    ok  	dupfind/internal/dupes	0.006s
    ok  	dupfind/internal/report	0.002s
    ok  	dupfind/internal/walk	0.007s
    FAIL
    completed: exit 1
- PASS: test-script-independent-of-cwd (FEATURE-TestGate.md)
  intent: bin/test.sh locates the application root itself, so running it from another directory still exits zero.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    timestamp: 2026-10-25_135539
    arguments: 
    working directory: /mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261025.174005/build/DupFind
    ?   	dupfind/cmd/dupfind	[no test files]
    ok  	dupfind/internal/cli	(cached)
    ok  	dupfind/internal/dupes	(cached)
    ok  	dupfind/internal/report	(cached)
    ok  	dupfind/internal/walk	(cached)
    completed: exit 0
- PASS: suite-has-tests-for-cli (FEATURE-TestGate.md)
  intent: The Go test suite contains test files for the cli package and they pass when run alone.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    ok  	dupfind/internal/cli	0.020s
- PASS: built-binary-exit-codes (FEATURE-TestGate.md)
  intent: The built executable exits 0 for a directory, 1 for a missing path, and 2 for each malformed argument form the suite covers.
  target interpreter: /mnt/c/Users/barlo/projects/drydock/.venv/bin/python
  provisioning: not required
  return code: 0
  stdout:
    [] 2 dupfind: expected exactly one directory argument, got 0; usage: dupfind [--summary] [--min-size <bytes>] <directory>
    
    ['/tmp/tmprudvtboz', '/tmp/tmprudvtboz'] 2 dupfind: expected exactly one directory argument, got 2; usage: dupfind [--summary] [--min-size <bytes>] <directory>
    
    ['--min-size', 'abc', '/tmp/tmprudvtboz'] 2 dupfind: --min-size must be a positive decimal integer, got "abc"; usage: dupfind [--summary] [--min-size <bytes>] <directory>
    
    ['--min-size', '0', '/tmp/tmprudvtboz'] 2 dupfind: --min-size must be a positive decimal integer, got "0"; usage: dupfind [--summary] [--min-size <bytes>] <directory>
    
    ['--frobnicate', '/tmp/tmprudvtboz'] 2 dupfind: flag provided but not defined: -frobnicate; usage: dupfind [--summary] [--min-size <bytes>] <directory>

## Build summary
I'll start by looking at the build directory and the staged sources.The grader in `sources/run.sh` is read-only, so I'll leave it alone. Next I'll read the existing scaffold from block 1.The block-1 scaffold is in place, and `cmd/dupfind/main.go` already delegates to `cli.Run`. I'll implement the packages to the specs now, starting with the walk and dupes layers.Next, the CLI: argument parsing, directory validation, and pipeline wiring.I left a stray `errors` placeholder at the bottom of that file. Removing it now.Before writing tests, I'll build and run the staged grader to see where the implementation stands.All 15 grader checks pass. Now the Go test suite, starting with the walk package.Next, the dupes package tests: size grouping, digest grouping, unreadable-file handling, and ordering.The hand-rolled insertion sort in that file is unnecessary, so I'll swap it for `sort.Strings`.Now the report tests (listing and summary, read back through `encoding/csv`) and the CLI tests.The feature spec names specific tests that the acceptance criteria invoke by exact name (`TestWriteSummary`, `TestGroupByDigest`, `TestGroupByDigestSkipsUnreadable`, `TestOrder`). My names don't match, so I'll rename them and add the missing `TestOrder`.Now the CLI tests. They replace the block-1 placeholder test, which asserted `Run(nil)` returns 0. That contradicts the new rule that no arguments is a usage error.Next I'll fix a stale comment in the dupes test helper, then run gofmt, vet, and the full test suite.The full suite passes. Gofmt rewrote two of my test files, which accounts for the on-disk change notices. Next I'll check for skipped tests, since a skip would hide an unverified path, and run the race detector.No skips and no failures under `-race`, and I'm not root, so the permission-based tests ran for real (76 passing). Now the POSIX `bin/test.sh` launcher, plus a `.gitignore` so the logs it writes stay out of git.Now I'll verify the gate as specified: POSIX syntax check, the test script from a different working directory, and the grader again.Block 2 is in place. The staged grader passes all 15 checks, and `sh bin/test.sh` exits 0 from the application root and from another directory.

Implemented per the specs:
- **arg-parsing / directory-validation:** `cli.Run` uses a local `flag.FlagSet` with `ContinueOnError` and discarded output. `--min-size` is validated by hand (digits only, ≥1, fits int64), since `flag.Int64` would accept `+5`, `0x10`, and similar. Usage errors return 2 with one stderr line. A missing or non-directory path returns 1.
- **walk:** `filepath.WalkDir`, which doesn't follow symlinks. Zero-byte and below-`--min-size` files are dropped. An unreadable subdirectory produces one warning and the walk continues.
- **size-grouping / digest-grouping:** same-size buckets, then SHA-256 over read-only opens. Unreadable files produce one warning each and are excluded.
- **group-ordering / listing / summary:** reclaimable bytes descending, then digest; paths in byte order; numbered from 1. Listing and summary are written with `encoding/csv`.
- **test-gate:** `bin/test.sh` is POSIX `sh` and runs `go vet` then `go test`.

Test suite: 76 tests across `walk`, `dupes`, `report`, and `cli`, with no skips. It runs as uid 1000, so the permission-based unreadable-file tests are enforced. The race detector is clean. Tests check the names the feature specs reference (`TestGroupBySizeOmitsUniqueSizes`, `TestGroupByDigest`, `TestGroupByDigestSkipsUnreadable`, `TestOrder`, `TestWriteSummary`).

Choices I made where the spec was open:
- **Replaced the block-1 placeholder test** in `internal/cli/cli_test.go`. It asserted `Run(nil)` returns 0, which contradicts the rule that no arguments is a usage error.
- **Symlinked directory argument:** `cli` resolves it with `filepath.EvalSymlinks` before walking, so a link to a directory is walked, as the validation spec requires.
- **`bin/test.sh` logging:** it writes to `data/logs/test_<stamp>.log` and prints the log at the end rather than live. Piping through `tee` needs `pipefail`, which is not POSIX.
- **Added `.gitignore`** for `data/`, `*.log`, and the built `/dupfind` binary, so the log output stays out of git.

I didn't run any Python AC containers (`drydock_go`); I verified each acceptance assertion's behavior through the grader and the Go tests. I did not commit, per the build instructions.

RESULT: SUCCESS

FILES CHANGED:
- internal/walk/walk.go
- internal/walk/walk_test.go
- internal/dupes/dupes.go
- internal/dupes/dupes_test.go
- internal/report/report.go
- internal/report/report_test.go
- internal/cli/cli.go
- internal/cli/cli_test.go
- bin/test.sh
- .gitignore

SUMMARY:
Implemented arg parsing, directory validation, recursive walk, size and SHA-256 digest grouping, ordering, CSV listing and summary, and the POSIX test gate. `sh sources/run.sh`: 15/15 passed (exit 0). `go vet ./...` clean; `go test -count=1 -race ./...`: 76 passed, 0 skipped, 0 failed. `sh bin/test.sh` exits 0 from the root and from another directory. `sources/run.sh` was not modified.

<blueprint-decisions>
[{"spec":"FEATURE-Walk.md","severity":"Low","subject":"Unreadable subdirectory warning","decision":"walk.Files emits one warning line per unreadable subdirectory and returns nil for it, so the walk continues. A failure to read the root itself is returned as an error (exit 1). Acceptable?"},{"spec":"FEATURE-DirectoryValidation.md","severity":"Low","subject":"Symlinked directory argument","decision":"cli resolves the directory with filepath.EvalSymlinks before walking so a symlink to a directory is walked; validation still uses os.Stat as specified."},{"spec":"FEATURE-TestGate.md","severity":"Low","subject":"bin/test.sh logging","decision":"Output goes to data/logs/test_<stamp>.log and is printed after completion rather than streamed, because pipefail/tee exit-status handling is not POSIX. Acceptable?"}]
</blueprint-decisions>

BLOCKERS:
- none
