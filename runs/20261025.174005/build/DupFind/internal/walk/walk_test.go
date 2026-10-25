package walk

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// writeFiles creates each named file under a fresh temporary root with the given contents.
func writeFiles(t *testing.T, files map[string]string) string {
	t.Helper()
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

// paths returns the Path of each entry in order.
func paths(entries []Entry) []string {
	out := []string{}
	for _, e := range entries {
		out = append(out, e.Path)
	}
	return out
}

// lockDir makes dir unreadable and restores its mode before the temporary directory is removed.
// It reports false, and the caller should skip, when permissions are not enforced here.
func lockDir(t *testing.T, dir string) bool {
	t.Helper()
	if err := os.Chmod(dir, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o755) })
	if _, err := os.ReadDir(dir); err == nil {
		return false
	}
	return true
}

func TestFilesReturnsNestedAndHiddenRegularFiles(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"top.txt":           "abc",
		"sub/deep/leaf.txt": "leafy",
		".hidden/.env":      "KEY=1",
	})

	got, err := Files(root, 1, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{".hidden/.env", "sub/deep/leaf.txt", "top.txt"}
	if !reflect.DeepEqual(paths(got), want) {
		t.Errorf("paths = %v, want %v", paths(got), want)
	}
}

func TestFilesReportsSizeOfEachFile(t *testing.T) {
	root := writeFiles(t, map[string]string{"a": "abc", "d/b": "abcdef"})

	got, err := Files(root, 1, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	sizes := map[string]int64{}
	for _, e := range got {
		sizes[e.Path] = e.Size
	}
	if sizes["a"] != 3 || sizes["d/b"] != 6 {
		t.Errorf("sizes = %v, want a=3 d/b=6", sizes)
	}
}

func TestFilesUsesSlashSeparatedRelativePaths(t *testing.T) {
	root := writeFiles(t, map[string]string{"one/two/three.txt": "x"})

	got, err := Files(root, 1, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"one/two/three.txt"}; !reflect.DeepEqual(paths(got), want) {
		t.Errorf("paths = %v, want %v", paths(got), want)
	}
}

func TestFilesOmitsSymlinkToFile(t *testing.T) {
	root := writeFiles(t, map[string]string{"real.txt": "content"})
	if err := os.Symlink("real.txt", filepath.Join(root, "link.txt")); err != nil {
		t.Skip("symbolic links not supported here:", err)
	}

	got, err := Files(root, 1, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"real.txt"}; !reflect.DeepEqual(paths(got), want) {
		t.Errorf("paths = %v, want %v", paths(got), want)
	}
}

func TestFilesDoesNotTraverseSymlinkToDirectory(t *testing.T) {
	root := writeFiles(t, map[string]string{"sub/inner.txt": "inner"})
	if err := os.Symlink("sub", filepath.Join(root, "dirlink")); err != nil {
		t.Skip("symbolic links not supported here:", err)
	}

	got, err := Files(root, 1, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"sub/inner.txt"}; !reflect.DeepEqual(paths(got), want) {
		t.Errorf("paths = %v, want %v", paths(got), want)
	}
}

func TestFilesOmitsZeroByteFilesEvenWhenMinSizeIsZero(t *testing.T) {
	root := writeFiles(t, map[string]string{"empty": "", "d/also-empty": "", "full": "x"})

	got, err := Files(root, 0, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"full"}; !reflect.DeepEqual(paths(got), want) {
		t.Errorf("paths = %v, want %v", paths(got), want)
	}
}

func TestFilesIncludesFileExactlyAtMinSize(t *testing.T) {
	root := writeFiles(t, map[string]string{"two": "ab", "three": "abc"})

	got, err := Files(root, 3, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"three"}; !reflect.DeepEqual(paths(got), want) {
		t.Errorf("paths = %v, want %v", paths(got), want)
	}
}

func TestFilesIgnoresFilesBelowMinSize(t *testing.T) {
	root := writeFiles(t, map[string]string{"one": "a", "four": "abcd"})

	got, err := Files(root, 2, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"four"}; !reflect.DeepEqual(paths(got), want) {
		t.Errorf("paths = %v, want %v", paths(got), want)
	}
}

func TestFilesEmptyDirectoryReturnsNoEntries(t *testing.T) {
	got, err := Files(t.TempDir(), 1, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("entries = %v, want none", paths(got))
	}
}

func TestFilesMissingRootIsError(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "absent")
	if _, err := Files(missing, 1, &bytes.Buffer{}); err == nil {
		t.Error("Files on a missing root returned no error")
	}
}

func TestFilesSkipsUnreadableSubdirectoryWithOneWarning(t *testing.T) {
	root := writeFiles(t, map[string]string{"ok.txt": "fine", "locked/hidden.txt": "secret"})
	if !lockDir(t, filepath.Join(root, "locked")) {
		t.Skip("directory permissions are not enforced here")
	}

	var warn bytes.Buffer
	got, err := Files(root, 1, &warn)
	if err != nil {
		t.Fatalf("Files returned error for unreadable subdirectory: %v", err)
	}
	if want := []string{"ok.txt"}; !reflect.DeepEqual(paths(got), want) {
		t.Errorf("paths = %v, want %v", paths(got), want)
	}
	if n := strings.Count(warn.String(), "\n"); n != 1 {
		t.Errorf("warning lines = %d, want 1: %q", n, warn.String())
	}
}
