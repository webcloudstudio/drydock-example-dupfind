package output

import (
	"bytes"
	"testing"

	"dupfind/internal/discovery"
)

func TestWriteListingOrdersGroupsAndCSVPaths(t *testing.T) {
	groups := []discovery.Group{
		{Size: 2, Digest: "b", Paths: []string{"z", "a"}},
		{Size: 2, Digest: "a", Paths: []string{"comma,name", "b"}},
	}
	var got bytes.Buffer
	if err := WriteListing(&got, groups); err != nil {
		t.Fatal(err)
	}
	want := "group,size,sha256,path\n1,2,a,b\n1,2,a,\"comma,name\"\n2,2,b,a\n2,2,b,z\n"
	if got.String() != want {
		t.Fatalf("listing = %q, want %q", got.String(), want)
	}
}

func TestWriteSummaryTotalsAndEmptyGroups(t *testing.T) {
	for _, tc := range []struct {
		name   string
		groups []discovery.Group
		want   string
	}{
		{"empty", nil, "groups,duplicate_files,reclaimable_bytes\n0,0,0\n"},
		{"totals", []discovery.Group{{Size: 5, Paths: []string{"a", "b", "c"}}, {Size: 7, Paths: []string{"d", "e"}}}, "groups,duplicate_files,reclaimable_bytes\n2,3,17\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got bytes.Buffer
			if err := WriteSummary(&got, tc.groups); err != nil {
				t.Fatal(err)
			}
			if got.String() != tc.want {
				t.Fatalf("summary = %q, want %q", got.String(), tc.want)
			}
		})
	}
}
