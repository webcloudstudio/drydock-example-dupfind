# BIG ERRORS — action required

- Command: `build`
- Phase: build step
- State: Failed
- Timestamp: 2026-10-24T07:45:01+00:00
- Execution ID: 20261024.074358.665Z-bbda77ab
- Challenge Execution ID: -
- Classification: programmatic acceptance failed: summary-no-duplicates, summary-filtered, listing-empty, listing-relative-paths, listing-csv-path
- Evidence / logs: /mnt/c/Users/barlo/projects/drydock/uat/DupFind/runs/20261024.073506/workspace/targets/DupFind/evidence/block-2.md

## Diagnostic

Block "Block 2 · Service" [block-2] failed its acceptance criteria.
  Story "Emit duplicate summary totals as CSV." [summary_output] does not meet its own acceptance criteria:
    - AC summary-no-duplicates — Summary succeeds with zero totals when no non-empty duplicate groups exist.
        assertion: PermissionError: [Errno 13] Permission denied: './cmd/dupfind'
        raised at: subprocess.py:1955
        process exit code: 1
        error: PermissionError: [Errno 13] Permission denied: './cmd/dupfind'
        values at failure:
          errno_num = 13
          err_msg = 'Permission denied'
          err_filename = './cmd/dupfind'
        check stderr:
          Traceback (most recent call last):
            File "summary-no-duplicates.py", line 11, in <module>
… (99 more lines truncated)

## Recovery

Run: drydock build DupFind
  to continue the build. This story resumes where it left off and is retried against the checks it failed.
