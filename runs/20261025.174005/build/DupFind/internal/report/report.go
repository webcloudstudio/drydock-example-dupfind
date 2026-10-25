// Package report writes duplicate listings and summaries as CSV using encoding/csv.
package report

import (
	"encoding/csv"
	"io"
	"strconv"

	"dupfind/internal/dupes"
)

// WriteListing writes the header group,size,sha256,path and one row per file of every group,
// in the order given.
func WriteListing(w io.Writer, groups []dupes.Group) error {
	cw := csv.NewWriter(w)
	cw.Write([]string{"group", "size", "sha256", "path"})
	for _, g := range groups {
		for _, p := range g.Paths {
			cw.Write([]string{strconv.Itoa(g.Number), strconv.FormatInt(g.Size, 10), g.Digest, p})
		}
	}
	cw.Flush()
	return cw.Error()
}

// WriteSummary writes the header groups,duplicate_files,reclaimable_bytes and one row holding
// the number of groups, the sum of files-1 over groups, and the total reclaimable bytes.
func WriteSummary(w io.Writer, groups []dupes.Group) error {
	var files, reclaimable int64
	for _, g := range groups {
		files += int64(len(g.Paths) - 1)
		reclaimable += g.Reclaimable()
	}
	cw := csv.NewWriter(w)
	cw.Write([]string{"groups", "duplicate_files", "reclaimable_bytes"})
	cw.Write([]string{
		strconv.Itoa(len(groups)),
		strconv.FormatInt(files, 10),
		strconv.FormatInt(reclaimable, 10),
	})
	cw.Flush()
	return cw.Error()
}
