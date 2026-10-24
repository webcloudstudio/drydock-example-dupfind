package discovery

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// File is a regular file relative to the scan root.
type File struct {
	Path string
	Size int64
}

// Group contains files with identical size and SHA-256 digest.
type Group struct {
	Size   int64
	Digest string
	Paths  []string
}

// Find walks root and groups duplicate regular files. Read failures are warnings.
func Find(root string, minSize int64, warning io.Writer) ([]Group, error) {
	bySize := make(map[int64][]string)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 || !entry.Type().IsRegular() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Size() == 0 || info.Size() < minSize {
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
	groups := make([]Group, 0)
	for size, paths := range bySize {
		if len(paths) < 2 {
			continue
		}
		byDigest := make(map[string][]string)
		for _, rel := range paths {
			full := filepath.Join(root, filepath.FromSlash(rel))
			f, err := os.Open(full)
			if err != nil {
				fmt.Fprintf(warning, "warning: cannot read %s: %v\n", rel, err)
				continue
			}
			h := sha256.New()
			_, readErr := io.Copy(h, f)
			closeErr := f.Close()
			if readErr != nil {
				fmt.Fprintf(warning, "warning: cannot read %s: %v\n", rel, readErr)
				continue
			}
			if closeErr != nil {
				fmt.Fprintf(warning, "warning: cannot read %s: %v\n", rel, closeErr)
				continue
			}
			byDigest[fmt.Sprintf("%x", h.Sum(nil))] = append(byDigest[fmt.Sprintf("%x", h.Sum(nil))], rel)
		}
		for digest, members := range byDigest {
			if len(members) < 2 {
				continue
			}
			sort.Slice(members, func(i, j int) bool { return members[i] < members[j] })
			groups = append(groups, Group{Size: size, Digest: digest, Paths: members})
		}
	}
	return groups, nil
}
