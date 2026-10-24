#!/bin/sh
# run.sh — DupFind grader. Do not filter, skip, or reinterpret.
#
# Builds ./cmd/dupfind with `go build`, creates a fixed directory tree, and compares every
# command's exit code and standard output with the known answer below. Standard error is never
# compared.
#
#   sh sources/run.sh                from the application root (how `drydock uat` runs it)
#   sh uat/DupFind/bin/run.sh APP    from anywhere, grading the application root APP
#
# Exit 0 only when the build succeeds and every check passes.
set -u

APP=$(cd "${1:-.}" && pwd)
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT INT TERM
BIN="$WORK/dupfind"

if ! (cd "$APP" && go build -o "$BIN" ./cmd/dupfind); then
    echo "error: go build ./cmd/dupfind failed at the application root $APP." >&2
    echo "The deliverable is a Go module whose main package is ./cmd/dupfind." >&2
    exit 1
fi

# make_tree <dir> — the exam tree. Every byte is fixed, so every hash is known.
make_tree() {
    T=$1
    mkdir -p "$T/a" "$T/b" "$T/c/deep/x" "$T/photos" "$T/music" "$T/notes" "$T/.hidden" "$T/big" \
        "$T/empty" "$T/unique"
    printf 'quarterly report\n' > "$T/a/report.txt"
    printf 'quarterly report\n' > "$T/b/report-copy.txt"
    printf 'quarterly report\n' > "$T/c/deep/x/report.txt"
    : > "$T/photos/img1.jpg"
    i=0
    while [ $i -lt 16 ]; do printf 'JPEG%04d' $i >> "$T/photos/img1.jpg"; i=$((i + 1)); done
    cp "$T/photos/img1.jpg" "$T/photos/img1 (copy).jpg"
    printf 'la la la\n' > "$T/music/song.mp3"
    printf 'la la la\n' > "$T/music/song,live.mp3"
    printf 'abcdefgh\n' > "$T/notes/one.txt"
    printf 'abcdefgX\n' > "$T/notes/two.txt"
    printf 'KEY=1\n' > "$T/.hidden/.env"
    printf 'KEY=1\n' > "$T/env.bak"
    head -c 1000 /dev/zero | tr '\0' z > "$T/big/blob-a.bin"
    head -c 1000 /dev/zero | tr '\0' z > "$T/big/blob-b.bin"
    head -c 1000 /dev/zero | tr '\0' y > "$T/big/blob-c.bin"
    : > "$T/empty/one"
    : > "$T/empty/two"
    ln -s ../a/report.txt "$T/notes/report-link.txt"
    printf 'only one of me\n' > "$T/unique/solo.txt"
}

TREE="$WORK/tree"
make_tree "$TREE"

PASSED=0
FAILED=0

# check <id> <expected exit> <dupfind args...>  — expected standard output on stdin.
check() {
    id=$1
    want_rc=$2
    shift 2
    cat > "$WORK/want"
    (cd "$WORK" && "$BIN" "$@") > "$WORK/got" 2> "$WORK/err"
    rc=$?
    if [ "$rc" -eq "$want_rc" ] && cmp -s "$WORK/want" "$WORK/got"; then
        PASSED=$((PASSED + 1))
        echo "PASS $id: dupfind $*"
    else
        FAILED=$((FAILED + 1))
        echo "FAIL $id: dupfind $*"
        echo "  exit: expected $want_rc, got $rc"
        diff "$WORK/want" "$WORK/got" | sed 's/^/  /'
        sed 's/^/  stderr: /' "$WORK/err"
    fi
}

# ── A. Full listing ──────────────────────────────────────────────────────────────────────
check A1 0 tree <<'EOF'
group,size,sha256,path
1,1000,950f88b09cf1d5e2cdbc5660c77dce3962265c548797950095629a0ea2daea46,big/blob-a.bin
1,1000,950f88b09cf1d5e2cdbc5660c77dce3962265c548797950095629a0ea2daea46,big/blob-b.bin
2,128,5c633a59e48621c6ca2a9c2c254792627ed975a86228be1fdf9d8a177520c0f7,photos/img1 (copy).jpg
2,128,5c633a59e48621c6ca2a9c2c254792627ed975a86228be1fdf9d8a177520c0f7,photos/img1.jpg
3,17,8a3c67892f82af22b58377b8e85bb41899053788d3aa3b81d2e6b4580e917c29,a/report.txt
3,17,8a3c67892f82af22b58377b8e85bb41899053788d3aa3b81d2e6b4580e917c29,b/report-copy.txt
3,17,8a3c67892f82af22b58377b8e85bb41899053788d3aa3b81d2e6b4580e917c29,c/deep/x/report.txt
4,9,31f86092953fd687403924ef23270af8da18b58147a0055468ecb4dc00b1c5c5,"music/song,live.mp3"
4,9,31f86092953fd687403924ef23270af8da18b58147a0055468ecb4dc00b1c5c5,music/song.mp3
5,6,929c28f541260cd8690babe97ae61472829cdf36e059c7b14e6b963f0f9ee77b,.hidden/.env
5,6,929c28f541260cd8690babe97ae61472829cdf36e059c7b14e6b963f0f9ee77b,env.bak
EOF
check A2 0 "$TREE" <<'EOF'
group,size,sha256,path
1,1000,950f88b09cf1d5e2cdbc5660c77dce3962265c548797950095629a0ea2daea46,big/blob-a.bin
1,1000,950f88b09cf1d5e2cdbc5660c77dce3962265c548797950095629a0ea2daea46,big/blob-b.bin
2,128,5c633a59e48621c6ca2a9c2c254792627ed975a86228be1fdf9d8a177520c0f7,photos/img1 (copy).jpg
2,128,5c633a59e48621c6ca2a9c2c254792627ed975a86228be1fdf9d8a177520c0f7,photos/img1.jpg
3,17,8a3c67892f82af22b58377b8e85bb41899053788d3aa3b81d2e6b4580e917c29,a/report.txt
3,17,8a3c67892f82af22b58377b8e85bb41899053788d3aa3b81d2e6b4580e917c29,b/report-copy.txt
3,17,8a3c67892f82af22b58377b8e85bb41899053788d3aa3b81d2e6b4580e917c29,c/deep/x/report.txt
4,9,31f86092953fd687403924ef23270af8da18b58147a0055468ecb4dc00b1c5c5,"music/song,live.mp3"
4,9,31f86092953fd687403924ef23270af8da18b58147a0055468ecb4dc00b1c5c5,music/song.mp3
5,6,929c28f541260cd8690babe97ae61472829cdf36e059c7b14e6b963f0f9ee77b,.hidden/.env
5,6,929c28f541260cd8690babe97ae61472829cdf36e059c7b14e6b963f0f9ee77b,env.bak
EOF

# ── B. Summary ───────────────────────────────────────────────────────────────────────────
check B1 0 --summary tree <<'EOF'
groups,duplicate_files,reclaimable_bytes
5,6,1177
EOF

# ── C. Minimum size ──────────────────────────────────────────────────────────────────────
check C1 0 --min-size 20 tree <<'EOF'
group,size,sha256,path
1,1000,950f88b09cf1d5e2cdbc5660c77dce3962265c548797950095629a0ea2daea46,big/blob-a.bin
1,1000,950f88b09cf1d5e2cdbc5660c77dce3962265c548797950095629a0ea2daea46,big/blob-b.bin
2,128,5c633a59e48621c6ca2a9c2c254792627ed975a86228be1fdf9d8a177520c0f7,photos/img1 (copy).jpg
2,128,5c633a59e48621c6ca2a9c2c254792627ed975a86228be1fdf9d8a177520c0f7,photos/img1.jpg
EOF
check C2 0 --summary --min-size 20 tree <<'EOF'
groups,duplicate_files,reclaimable_bytes
2,2,1128
EOF
check C3 0 --min-size 200 tree <<'EOF'
group,size,sha256,path
1,1000,950f88b09cf1d5e2cdbc5660c77dce3962265c548797950095629a0ea2daea46,big/blob-a.bin
1,1000,950f88b09cf1d5e2cdbc5660c77dce3962265c548797950095629a0ea2daea46,big/blob-b.bin
EOF

# ── D. No duplicates ─────────────────────────────────────────────────────────────────────
check D1 0 tree/unique <<'EOF'
group,size,sha256,path
EOF
check D2 0 --summary tree/empty <<'EOF'
groups,duplicate_files,reclaimable_bytes
0,0,0
EOF

# ── E. Operational failures ──────────────────────────────────────────────────────────────
check E1 1 tree/no-such-directory < /dev/null
check E2 1 tree/env.bak < /dev/null

# ── U. Usage errors ──────────────────────────────────────────────────────────────────────
check U1 2 < /dev/null
check U2 2 tree tree/a < /dev/null
check U3 2 --min-size abc tree < /dev/null
check U4 2 --min-size 0 tree < /dev/null
check U5 2 --frobnicate tree < /dev/null

TOTAL=$((PASSED + FAILED))
echo "DupFind: $PASSED passed, $FAILED failed (of $TOTAL)"
[ "$FAILED" -eq 0 ]
