package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCommandBuilds(t *testing.T) {
	root := filepath.Clean(filepath.Join("..", ".."))
	cmd := exec.Command("go", "build", "-o", filepath.Join(t.TempDir(), "dupfind"), "./cmd/dupfind")
	cmd.Dir = root
	if err := cmd.Run(); err != nil {
		t.Fatalf("go build ./cmd/dupfind: %v", err)
	}
}

func TestApplicationPackageDirectoriesExist(t *testing.T) {
	for _, name := range []string{"cli", "discovery", "output"} {
		if stat, err := os.Stat(filepath.Join("..", name)); err != nil || !stat.IsDir() {
			t.Errorf("internal/%s directory is unavailable: %v", name, err)
		}
	}
}
