package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRunWritesSummaryForValidOptions(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"a", "b"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("same"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"--summary", "--min-size", "2", root}, &stdout, &stderr); code != 0 {
		t.Fatalf("Run() = %d, stderr %q", code, stderr.String())
	}
	if stdout.String() != "groups,duplicate_files,reclaimable_bytes\n1,1,4\n" {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunRejectsMalformedArgumentsWithoutStdout(t *testing.T) {
	for _, args := range [][]string{{}, {"one", "two"}, {"--min-size", "0", "dir"}, {"--min-size", "bad", "dir"}, {"--unknown", "dir"}} {
		t.Run(filepath.Join(args...), func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := Run(args, &stdout, &stderr); code != 2 {
				t.Fatalf("Run(%v) = %d", args, code)
			}
			if stdout.Len() != 0 || stderr.Len() == 0 {
				t.Fatalf("stdout=%q stderr=%q", stdout.String(), stderr.String())
			}
		})
	}
}

func TestRunRejectsMissingAndNonDirectoryPaths(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{filepath.Join(root, "missing"), file} {
		var stdout, stderr bytes.Buffer
		if code := Run([]string{path}, &stdout, &stderr); code != 1 {
			t.Fatalf("Run(%q) = %d", path, code)
		}
		if stdout.Len() != 0 || stderr.Len() == 0 {
			t.Fatalf("stdout=%q stderr=%q", stdout.String(), stderr.String())
		}
	}
}
