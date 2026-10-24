package main

import (
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestExecutablePackageBuilds(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "dupfind")
	cmd := exec.Command("go", "build", "-o", binary, "./cmd/dupfind")
	cmd.Dir = filepath.Join("..", "..")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build ./cmd/dupfind failed: %v\n%s", err, output)
	}
}

func TestApplicationPackageLayout(t *testing.T) {
	cmd := exec.Command("go", "list", "-f", "{{.ImportPath}}", "./...")
	cmd.Dir = filepath.Join("..", "..")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list ./... failed: %v", err)
	}
	got := strings.Fields(string(output))
	sort.Strings(got)
	want := []string{
		"dupfind/cmd/dupfind",
		"dupfind/internal/cli",
		"dupfind/internal/output",
		"dupfind/internal/scan",
	}
	if len(got) != len(want) {
		t.Fatalf("package count = %d, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("package %d = %q, want %q", i, got[i], want[i])
		}
	}
}
