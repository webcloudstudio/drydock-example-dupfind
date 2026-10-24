package discovery

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestFindGroupsNestedFilesByContent(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"a/file", ".hidden/deep/file", "single"} {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("identical"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	groups, err := Find(root, 1, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	wantDigest := fmt.Sprintf("%x", sha256.Sum256([]byte("identical")))
	if len(groups) != 1 || groups[0].Size != 9 || groups[0].Digest != wantDigest || fmt.Sprint(groups[0].Paths) != "[.hidden/deep/file a/file single]" {
		t.Fatalf("Find() = %#v", groups)
	}
}

func TestFindIgnoresEmptyFilesAndSymlinks(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"empty-a", "empty-b"} {
		if err := os.WriteFile(filepath.Join(root, name), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "target"), []byte("value"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "target"), filepath.Join(root, "link")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	groups, err := Find(root, 1, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 0 {
		t.Fatalf("Find() = %#v, want no groups", groups)
	}
}

func TestFindAppliesMinimumSize(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"a", "b"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("1234"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	groups, err := Find(root, 5, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 0 {
		t.Fatalf("Find() = %#v, want no groups", groups)
	}
}
