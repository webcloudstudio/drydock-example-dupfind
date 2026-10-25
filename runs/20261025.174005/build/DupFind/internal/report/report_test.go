package report

import (
	"bytes"
	"encoding/csv"
	"reflect"
	"testing"

	"dupfind/internal/dupes"
)

// parseCSV reads every record of text with encoding/csv.
func parseCSV(t *testing.T, text string) [][]string {
	t.Helper()
	records, err := csv.NewReader(bytes.NewBufferString(text)).ReadAll()
	if err != nil {
		t.Fatalf("output is not valid CSV: %v\n%s", err, text)
	}
	return records
}

// sampleGroups returns two ordered groups: a three-copy group and a two-copy group.
func sampleGroups() []dupes.Group {
	return []dupes.Group{
		{Number: 1, Size: 10, Digest: "d1", Paths: []string{"a/one", "b/two", "c/three"}},
		{Number: 2, Size: 4, Digest: "d2", Paths: []string{"x", "y"}},
	}
}

func TestWriteListingHeaderAndRows(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteListing(&buf, sampleGroups()); err != nil {
		t.Fatal(err)
	}

	got := parseCSV(t, buf.String())
	want := [][]string{
		{"group", "size", "sha256", "path"},
		{"1", "10", "d1", "a/one"},
		{"1", "10", "d1", "b/two"},
		{"1", "10", "d1", "c/three"},
		{"2", "4", "d2", "x"},
		{"2", "4", "d2", "y"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("listing = %v, want %v", got, want)
	}
}

func TestWriteListingWithNoGroupsWritesHeaderOnly(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteListing(&buf, nil); err != nil {
		t.Fatal(err)
	}

	got := parseCSV(t, buf.String())
	if want := [][]string{{"group", "size", "sha256", "path"}}; !reflect.DeepEqual(got, want) {
		t.Errorf("listing = %v, want %v", got, want)
	}
}

func TestWriteListingUsesLFLineEndings(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteListing(&buf, sampleGroups()); err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(buf.Bytes(), []byte("\r")) {
		t.Errorf("listing contains carriage return: %q", buf.String())
	}
}

func TestWriteListingQuotesCommaAndQuoteInPaths(t *testing.T) {
	paths := []string{"song,live.mp3", `say "hi".txt`, "plain.txt"}
	groups := []dupes.Group{{Number: 1, Size: 8, Digest: "d", Paths: paths}}

	var buf bytes.Buffer
	if err := WriteListing(&buf, groups); err != nil {
		t.Fatal(err)
	}

	got := parseCSV(t, buf.String())
	for i, rec := range got[1:] {
		if len(rec) != 4 {
			t.Fatalf("row %d has %d fields, want 4", i, len(rec))
		}
		if rec[3] != paths[i] {
			t.Errorf("path %d round-trip = %q, want %q", i, rec[3], paths[i])
		}
	}
}

func TestWriteSummary(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteSummary(&buf, sampleGroups()); err != nil {
		t.Fatal(err)
	}

	got := parseCSV(t, buf.String())
	want := [][]string{
		{"groups", "duplicate_files", "reclaimable_bytes"},
		{"2", "3", "24"}, // 10*2 + 4*1 reclaimable; 2 + 1 duplicate files
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("summary = %v, want %v", got, want)
	}
}

func TestWriteSummaryNoGroupsIsZeroRow(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteSummary(&buf, nil); err != nil {
		t.Fatal(err)
	}

	got := parseCSV(t, buf.String())
	want := [][]string{
		{"groups", "duplicate_files", "reclaimable_bytes"},
		{"0", "0", "0"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("summary = %v, want %v", got, want)
	}
}

func TestWriteSummaryCountsOnlyDuplicateCopies(t *testing.T) {
	groups := []dupes.Group{{Number: 1, Size: 5, Digest: "d", Paths: []string{"a", "b", "c", "d"}}}

	var buf bytes.Buffer
	if err := WriteSummary(&buf, groups); err != nil {
		t.Fatal(err)
	}

	got := parseCSV(t, buf.String())
	if want := []string{"1", "3", "15"}; !reflect.DeepEqual(got[1], want) {
		t.Errorf("summary row = %v, want %v", got[1], want)
	}
}
