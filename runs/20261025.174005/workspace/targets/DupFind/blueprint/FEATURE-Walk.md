# FEATURE: Walk

| Field       | Value |
|-------------|-------|
| Version     | 20261025 V1 |
| Description | Walk the directory recursively and yield eligible regular files at or above the minimum size. |
| Depends On  | ARCHITECTURE.md, FEATURE-DirectoryValidation.md |
| Provides    | walk.Files |
| Consumes    | cli.Config |

## Purpose

`internal/walk` enumerates the files that may be considered for duplicate detection. It never opens
file contents.

## Interface

```go
package walk

type Entry struct {
	Path string // relative to root, "/" separated, no leading "./"
	Size int64
}

func Files(root string, minSize int64, warn io.Writer) ([]Entry, error)
```

## Behavior

- Traversal is recursive and includes hidden files and directories (names beginning `.`).
- Only regular files are returned. Entries are classified without following links (`os.Lstat` or the `fs.DirEntry` type).
- Symbolic links are neither followed nor reported, whether they point at files or directories. A directory reachable only through a symbolic link is not traversed.
- Files whose size is below `minSize` are ignored. Zero-byte files are never returned, even if `minSize` is below 1.
- Paths use `/` as the separator and no leading `./`.
- A subdirectory that cannot be read yields one warning line on `warn`, is skipped, and the walk continues. This does not change the exit code.
- Return order is unspecified; ordering is the responsibility of later stories.

## Reads

Directory listings and file metadata under the root.

## Writes

Warning lines on the `warn` writer only.

## Programmatic Acceptance

=== AC walk-rules ===
Intent: Walk returns hidden and nested regular files, omits symbolic links and zero-byte files, and honors the minimum size.

import os
import pathlib
import subprocess
import sys

mod = ""
for line in pathlib.Path("go.mod").read_text().splitlines():
    if line.startswith("module "):
        mod = line.split()[1]
assert mod

GO_TEST = """package walk_test

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"MODPATH/internal/walk"
)

func build(t *testing.T, files map[string]string) string {
	root := t.TempDir()
	for name, body := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func check(t *testing.T, root string, files map[string]string, min int) {
	got, err := walk.Files(root, int64(min), io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	var gotPaths []string
	for _, e := range got {
		gotPaths = append(gotPaths, e.Path)
		if int(e.Size) != len(files[e.Path]) {
			t.Errorf("size of %s = %d", e.Path, e.Size)
		}
	}
	var want []string
	for name, body := range files {
		if len(body) >= min && len(body) > 0 {
			want = append(want, name)
		}
	}
	sort.Strings(gotPaths)
	sort.Strings(want)
	if len(gotPaths) != len(want) {
		t.Fatalf("got %v want %v", gotPaths, want)
	}
	for i := range want {
		if gotPaths[i] != want[i] {
			t.Fatalf("got %v want %v", gotPaths, want)
		}
	}
}

func TestAcceptanceWalkEligibility(t *testing.T) {
	files := map[string]string{"a.txt": "abc", ".hid/.x": "hidden", "sub/deep/f": "ff", "empty": "", "b": "z"}
	root := build(t, files)
	if err := os.Symlink("a.txt", filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("sub", filepath.Join(root, "dirlink")); err != nil {
		t.Fatal(err)
	}
	check(t, root, files, 1)
}

func TestAcceptanceWalkMinSize(t *testing.T) {
	files := map[string]string{"one": "a", "two": "ab", "three": "abc", "d/four": "abcd"}
	root := build(t, files)
	check(t, root, files, 3)
	check(t, root, files, 1)
}
""".replace("MODPATH", mod)

target = pathlib.Path("internal/walk/zz_acceptance_test.go")
target.write_text(GO_TEST)
try:
    result = subprocess.run(
        ["go", "test", "./internal/walk", "-run", "TestAcceptanceWalk", "-count=1", "-v"],
        capture_output=True, text=True,
    )
    print(result.stdout)
    print(result.stderr, file=sys.stderr)
    assert result.returncode == 0
finally:
    if target.exists():
        os.remove(target)
=== END AC walk-rules ===

=== AC walk-vet ===
Intent: The walk package compiles and passes go vet.

import subprocess
import sys

result = subprocess.run(["go", "vet", "./internal/walk"], capture_output=True, text=True)
print(result.stdout)
print(result.stderr, file=sys.stderr)
assert result.returncode == 0
=== END AC walk-vet ===

## User Acceptance

- None.

## Guardrails

- Never follow or report a symbolic link, and never traverse a directory reached only through one.
- Never open file contents in the walk.
- Never delete, move, or write any file under the directory examined.
- Zero-byte files are never reported.
