// Package output formats deterministic CSV results.
package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"

	"dupfind/internal/scan"
)

// Ordered returns a sorted copy of groups, with paths sorted by byte order.
func Ordered(groups []scan.Group) []scan.Group {
	out := make([]scan.Group, len(groups))
	copy(out, groups)
	for i := range out {
		out[i].Paths = append([]string(nil), out[i].Paths...)
		sort.Slice(out[i].Paths, func(a, b int) bool { return out[i].Paths[a] < out[i].Paths[b] })
	}
	sort.Slice(out, func(i, j int) bool {
		a := out[i].Size * int64(len(out[i].Paths)-1)
		b := out[j].Size * int64(len(out[j].Paths)-1)
		if a != b {
			return a > b
		}
		return out[i].SHA256 < out[j].SHA256
	})
	return out
}

// Listing writes the duplicate listing CSV.
func Listing(w io.Writer, groups []scan.Group) error {
	c := csv.NewWriter(w)
	if err := c.Write([]string{"group", "size", "sha256", "path"}); err != nil {
		return err
	}
	for i, g := range Ordered(groups) {
		for _, p := range g.Paths {
			if err := c.Write([]string{fmt.Sprint(i + 1), fmt.Sprint(g.Size), g.SHA256, p}); err != nil {
				return err
			}
		}
	}
	c.Flush()
	return c.Error()
}

// Summary writes group, duplicate-file, and reclaimable-byte totals.
func Summary(w io.Writer, groups []scan.Group) error {
	c := csv.NewWriter(w)
	ordered := Ordered(groups)
	var duplicates, reclaimable int64
	for _, g := range ordered {
		copies := int64(len(g.Paths) - 1)
		duplicates += copies
		reclaimable += g.Size * copies
	}
	if err := c.Write([]string{"groups", "duplicate_files", "reclaimable_bytes"}); err != nil {
		return err
	}
	if err := c.Write([]string{fmt.Sprint(len(ordered)), fmt.Sprint(duplicates), fmt.Sprint(reclaimable)}); err != nil {
		return err
	}
	c.Flush()
	return c.Error()
}
