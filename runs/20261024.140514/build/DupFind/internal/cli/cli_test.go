package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestRunListsExactDuplicateFiles(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a"), []byte("same"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b"), []byte("same"), 0600); err != nil {
		t.Fatal(err)
	}
	var out, errout bytes.Buffer
	if code := Run([]string{dir}, &out, &errout); code != 0 {
		t.Fatalf("code = %d", code)
	}
	digest := sha256.Sum256([]byte("same"))
	want := "group,size,sha256,path\n1,4," + hex.EncodeToString(digest[:]) + ",a\n1,4," + hex.EncodeToString(digest[:]) + ",b\n"
	if out.String() != want || errout.Len() != 0 {
		t.Fatalf("stdout=%q want=%q stderr=%q", out.String(), want, errout.String())
	}
}

func TestRunRejectsUsageAndOperationalErrors(t *testing.T) {
	for _, tc := range []struct {
		args []string
		want int
	}{{nil, 2}, {[]string{"--min-size", "0", "."}, 2}, {[]string{"--frobnicate", "."}, 2}, {[]string{"missing-dir-df"}, 1}} {
		var out, errout bytes.Buffer
		if got := Run(tc.args, &out, &errout); got != tc.want {
			t.Errorf("Run(%v) = %d, want %d", tc.args, got, tc.want)
		}
		if out.Len() != 0 || errout.Len() == 0 {
			t.Errorf("Run(%v) stdout=%q stderr=%q", tc.args, out.String(), errout.String())
		}
	}
}

func TestRunSummary(t *testing.T) {
	dir := t.TempDir()
	for _, n := range []string{"a", "b"} {
		if err := os.WriteFile(filepath.Join(dir, n), []byte("same"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	var out, errout bytes.Buffer
	if got := Run([]string{"--summary", dir}, &out, &errout); got != 0 {
		t.Fatalf("code = %d, stderr=%q", got, errout.String())
	}
	if got, want := out.String(), "groups,duplicate_files,reclaimable_bytes\n1,1,4\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}
