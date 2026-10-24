package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"

	"dupfind/internal/discovery"
)

func WriteListing(w io.Writer, groups []discovery.Group) error {
	ordered := append([]discovery.Group(nil), groups...)
	sort.Slice(ordered, func(i, j int) bool {
		a, b := ordered[i], ordered[j]
		ra, rb := a.Size*int64(len(a.Paths)-1), b.Size*int64(len(b.Paths)-1)
		if ra != rb {
			return ra > rb
		}
		return a.Digest < b.Digest
	})
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"group", "size", "sha256", "path"}); err != nil {
		return err
	}
	for i, g := range ordered {
		paths := append([]string(nil), g.Paths...)
		sort.Strings(paths)
		for _, p := range paths {
			if err := cw.Write([]string{fmt.Sprint(i + 1), fmt.Sprint(g.Size), g.Digest, p}); err != nil {
				return err
			}
		}
	}
	cw.Flush()
	return cw.Error()
}

func WriteSummary(w io.Writer, groups []discovery.Group) error {
	var files, bytes int64
	for _, g := range groups {
		n := int64(len(g.Paths) - 1)
		files += n
		bytes += g.Size * n
	}
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"groups", "duplicate_files", "reclaimable_bytes"}); err != nil {
		return err
	}
	if err := cw.Write([]string{fmt.Sprint(len(groups)), fmt.Sprint(files), fmt.Sprint(bytes)}); err != nil {
		return err
	}
	cw.Flush()
	return cw.Error()
}
