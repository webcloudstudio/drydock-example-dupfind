#!/bin/sh
# CommandCenter Operation
# Name: Test
# Type: batch
#
# Runs go vet and go test over the application root. Exits 0 only when both pass.
# POSIX sh: no bashisms, so it runs under dash as well as bash.
set -eu

ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$ROOT"

STAMP=$(date +%Y-%m-%d_%H%M%S)
LOG_DIR="$ROOT/data/logs"
LOG="$LOG_DIR/test_$STAMP.log"
mkdir -p "$LOG_DIR"

{
    echo "timestamp: $STAMP"
    echo "arguments: $*"
    echo "working directory: $ROOT"
} > "$LOG"

rc=0
( go vet ./... && go test ./... ) >> "$LOG" 2>&1 || rc=$?

echo "completed: exit $rc" >> "$LOG"
cat "$LOG"
exit "$rc"
