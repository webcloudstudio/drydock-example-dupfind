// Package dupes groups files by size and SHA-256 digest and computes reclaimable bytes.
package dupes

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"dupfind/internal/walk"
)

// Group is two or more files with identical size and contents.
type Group struct {
	Number int      // 1-based position, assigned by Order
	Size   int64    // bytes per file
	Digest string   // lowercase hexadecimal SHA-256 of the contents
	Paths  []string // relative to the examined directory, sorted by Order
}

// Reclaimable returns the bytes recovered by removing every copy but one.
func (g Group) Reclaimable() int64 {
	return g.Size * int64(len(g.Paths)-1)
}

// GroupBySize buckets entries by exact size and returns only the buckets holding two or more
// files, ordered by ascending size. File contents are not read.
func GroupBySize(entries []walk.Entry) [][]walk.Entry {
	bySize := make(map[int64][]walk.Entry)
	for _, e := range entries {
		bySize[e.Size] = append(bySize[e.Size], e)
	}
	sizes := make([]int64, 0, len(bySize))
	for size, bucket := range bySize {
		if len(bucket) >= 2 {
			sizes = append(sizes, size)
		}
	}
	sort.Slice(sizes, func(i, j int) bool { return sizes[i] < sizes[j] })
	buckets := make([][]walk.Entry, 0, len(sizes))
	for _, size := range sizes {
		buckets = append(buckets, bySize[size])
	}
	return buckets
}

// GroupByDigest hashes the files of each size bucket with SHA-256 and returns a group for every
// digest shared by two or more files. A file that cannot be read yields one line on warn and is
// excluded from grouping. Files are opened read-only.
func GroupByDigest(root string, buckets [][]walk.Entry, warn io.Writer) []Group {
	var groups []Group
	for _, bucket := range buckets {
		byDigest := make(map[string][]string)
		for _, e := range bucket {
			digest, err := digestFile(filepath.Join(root, filepath.FromSlash(e.Path)))
			if err != nil {
				fmt.Fprintf(warn, "dupfind: warning: cannot read %s: %v\n", e.Path, err)
				continue
			}
			byDigest[digest] = append(byDigest[digest], e.Path)
		}
		for digest, paths := range byDigest {
			if len(paths) >= 2 {
				groups = append(groups, Group{Size: bucket[0].Size, Digest: digest, Paths: paths})
			}
		}
	}
	return groups
}

// Order sorts paths within each group by byte value, orders groups by reclaimable bytes
// descending and then by digest ascending, and numbers the groups from 1.
func Order(groups []Group) []Group {
	for i := range groups {
		sort.Strings(groups[i].Paths)
	}
	sort.Slice(groups, func(i, j int) bool {
		ri, rj := groups[i].Reclaimable(), groups[j].Reclaimable()
		if ri != rj {
			return ri > rj
		}
		return groups[i].Digest < groups[j].Digest
	})
	for i := range groups {
		groups[i].Number = i + 1
	}
	return groups
}

// digestFile returns the lowercase hexadecimal SHA-256 of the file at path.
func digestFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
