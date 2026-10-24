// Package scan finds byte-identical regular files in a directory tree.
package scan

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Group describes files with the same size and content digest.
type Group struct {
	Size   int64
	SHA256 string
	Paths  []string
}

// Scan walks root without following symbolic links. Files that cannot be read
// are skipped with a warning; traversal errors are returned.
func Scan(root string, minSize int64, warnings io.Writer) ([]Group, error) {
	bySize := make(map[int64][]string)
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() || info.Size() < minSize {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		bySize[info.Size()] = append(bySize[info.Size()], filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, err
	}
	var groups []Group
	for size, paths := range bySize {
		if len(paths) < 2 {
			continue
		}
		byDigest := make(map[string][]string)
		for _, rel := range paths {
			f, err := os.Open(filepath.Join(root, filepath.FromSlash(rel)))
			if err != nil {
				if warnings != nil {
					fmt.Fprintf(warnings, "warning: cannot read %s: %v\n", rel, err)
				}
				continue
			}
			h := sha256.New()
			_, readErr := io.Copy(h, f)
			closeErr := f.Close()
			if readErr != nil || closeErr != nil {
				if warnings != nil {
					cause := readErr
					if cause == nil {
						cause = closeErr
					}
					fmt.Fprintf(warnings, "warning: cannot read %s: %v\n", rel, cause)
				}
				continue
			}
			byDigest[fmt.Sprintf("%x", h.Sum(nil))] = append(byDigest[fmt.Sprintf("%x", h.Sum(nil))], rel)
		}
		for digest, matches := range byDigest {
			if len(matches) >= 2 {
				groups = append(groups, Group{Size: size, SHA256: digest, Paths: matches})
			}
		}
	}
	return groups, nil
}

// NormalizePaths is retained as a small utility for consumers constructing groups.
func NormalizePaths(paths []string) {
	sort.Slice(paths, func(i, j int) bool { return strings.Compare(paths[i], paths[j]) < 0 })
}
