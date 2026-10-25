// Package walk traverses a directory tree and yields eligible regular files.
package walk

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"sort"
)

// Entry is a regular file eligible for duplicate detection.
type Entry struct {
	Path string // relative to root, "/" separated, no leading "./"
	Size int64
}

// Files walks root recursively and returns every regular file whose size is at least minSize.
// Symbolic links are neither followed nor reported, so a directory reachable only through a link
// is not traversed. Zero-byte files are never returned. A subdirectory that cannot be read yields
// one line on warn and is skipped. File contents are never opened.
func Files(root string, minSize int64, warn io.Writer) ([]Entry, error) {
	var entries []Entry
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if path == root {
				return err
			}
			fmt.Fprintf(warn, "dupfind: warning: cannot read directory %s: %v\n", path, err)
			return nil
		}
		if d.IsDir() || !d.Type().IsRegular() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			fmt.Fprintf(warn, "dupfind: warning: cannot stat %s: %v\n", path, err)
			return nil
		}
		if info.Size() == 0 || info.Size() < minSize {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		entries = append(entries, Entry{Path: filepath.ToSlash(rel), Size: info.Size()})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	return entries, nil
}
