package dupes

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"dupfind/internal/walk"
)

// writeFiles creates each named file under a fresh temporary root with the given contents.
func writeFiles(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, body := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// entriesOf builds walk entries for the named files, sizing each from its contents.
func entriesOf(files map[string]string) []walk.Entry {
	var out []walk.Entry
	for name, body := range files {
		out = append(out, walk.Entry{Path: name, Size: int64(len(body))})
	}
	return out
}

// sha256Hex returns the lowercase hexadecimal SHA-256 of body.
func sha256Hex(body string) string {
	sum := sha256.Sum256([]byte(body))
	return hex.EncodeToString(sum[:])
}

// bucketPaths returns the sorted paths of each bucket, in bucket order.
func bucketPaths(buckets [][]walk.Entry) [][]string {
	var out [][]string
	for _, b := range buckets {
		var ps []string
		for _, e := range b {
			ps = append(ps, e.Path)
		}
		out = append(out, sortedStrings(ps))
	}
	return out
}

// sortedStrings returns a sorted copy of in.
func sortedStrings(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}

// lockFile makes path unreadable and restores its mode before cleanup. It reports false, and the
// caller should skip, when permissions are not enforced here.
func lockFile(t *testing.T, path string) bool {
	t.Helper()
	if err := os.Chmod(path, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(path, 0o644) })
	if f, err := os.Open(path); err == nil {
		f.Close()
		return false
	}
	return true
}

func TestGroupBySizeOmitsUniqueSizes(t *testing.T) {
	entries := entriesOf(map[string]string{"a": "xx", "b": "yyy", "c": "zzzz"})

	if got := GroupBySize(entries); len(got) != 0 {
		t.Errorf("buckets = %v, want none", bucketPaths(got))
	}
}

func TestGroupBySizeKeepsSharedSizes(t *testing.T) {
	entries := entriesOf(map[string]string{
		"p": "aaaaa", "q": "bbbbb",
		"r": "ccccccccc", "s": "ddddddddd", "t": "eeeeeeeee",
		"u": "f",
	})

	got := bucketPaths(GroupBySize(entries))
	want := [][]string{{"p", "q"}, {"r", "s", "t"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("buckets = %v, want %v", got, want)
	}
}

func TestGroupBySizeEmptyInputReturnsNoBuckets(t *testing.T) {
	if got := GroupBySize(nil); len(got) != 0 {
		t.Errorf("buckets = %v, want none", bucketPaths(got))
	}
}

func TestGroupByDigest(t *testing.T) {
	files := map[string]string{"a/one": "same\n", "b/two": "same\n"}
	root := writeFiles(t, files)

	got := GroupByDigest(root, GroupBySize(entriesOf(files)), &bytes.Buffer{})
	if len(got) != 1 {
		t.Fatalf("groups = %d, want 1", len(got))
	}
	g := got[0]
	if g.Size != 5 || g.Digest != sha256Hex("same\n") {
		t.Errorf("group size/digest = %d/%s, want 5/%s", g.Size, g.Digest, sha256Hex("same\n"))
	}
	if !reflect.DeepEqual(sortedStrings(g.Paths), []string{"a/one", "b/two"}) {
		t.Errorf("paths = %v, want a/one and b/two", g.Paths)
	}
}

func TestGroupByDigestSameSizeDifferentContentsIsNotDuplicate(t *testing.T) {
	files := map[string]string{"a": "aaaa", "b": "bbbb"}
	root := writeFiles(t, files)

	if got := GroupByDigest(root, GroupBySize(entriesOf(files)), &bytes.Buffer{}); len(got) != 0 {
		t.Errorf("groups = %+v, want none", got)
	}
}

func TestGroupByDigestSplitsSameSizeIntoSeparateGroups(t *testing.T) {
	files := map[string]string{"x1": "abcd", "x2": "abcd", "y1": "abce", "y2": "abce", "z1": "abcf"}
	root := writeFiles(t, files)

	got := GroupByDigest(root, GroupBySize(entriesOf(files)), &bytes.Buffer{})
	if len(got) != 2 {
		t.Fatalf("groups = %d, want 2", len(got))
	}
	byDigest := map[string][]string{}
	for _, g := range got {
		byDigest[g.Digest] = sortedStrings(g.Paths)
	}
	want := map[string][]string{
		sha256Hex("abcd"): {"x1", "x2"},
		sha256Hex("abce"): {"y1", "y2"},
	}
	if !reflect.DeepEqual(byDigest, want) {
		t.Errorf("groups = %v, want %v", byDigest, want)
	}
}

func TestGroupByDigestSkipsUnreadable(t *testing.T) {
	files := map[string]string{"ok1": "same", "ok2": "same", "locked": "same"}
	root := writeFiles(t, files)
	if !lockFile(t, filepath.Join(root, "locked")) {
		t.Skip("file permissions are not enforced here")
	}

	var warn bytes.Buffer
	got := GroupByDigest(root, GroupBySize(entriesOf(files)), &warn)
	if len(got) != 1 {
		t.Fatalf("groups = %d, want 1", len(got))
	}
	if want := []string{"ok1", "ok2"}; !reflect.DeepEqual(sortedStrings(got[0].Paths), want) {
		t.Errorf("paths = %v, want %v", got[0].Paths, want)
	}
	if n := strings.Count(warn.String(), "\n"); n != 1 {
		t.Errorf("warning lines = %d, want 1: %q", n, warn.String())
	}
}

func TestGroupByDigestNeverOpensUniqueSizeFiles(t *testing.T) {
	files := map[string]string{"solo": "unique-size", "twin1": "pair", "twin2": "pair"}
	root := writeFiles(t, files)
	if !lockFile(t, filepath.Join(root, "solo")) {
		t.Skip("file permissions are not enforced here")
	}

	var warn bytes.Buffer
	GroupByDigest(root, GroupBySize(entriesOf(files)), &warn)
	if warn.Len() != 0 {
		t.Errorf("warning written for a unique-size file: %q", warn.String())
	}
}

func TestGroupByDigestDoesNotModifyFiles(t *testing.T) {
	files := map[string]string{"a": "twin", "b": "twin"}
	root := writeFiles(t, files)

	GroupByDigest(root, GroupBySize(entriesOf(files)), &bytes.Buffer{})
	for name, want := range files {
		got, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}
}

func TestGroupReclaimableForTwoAndThreeCopies(t *testing.T) {
	two := Group{Size: 10, Paths: []string{"a", "b"}}
	three := Group{Size: 7, Paths: []string{"a", "b", "c"}}

	if got := two.Reclaimable(); got != 10 {
		t.Errorf("two-copy reclaimable = %d, want 10", got)
	}
	if got := three.Reclaimable(); got != 14 {
		t.Errorf("three-copy reclaimable = %d, want 14", got)
	}
}

func TestOrder(t *testing.T) {
	groups := []Group{
		{Size: 1, Digest: "cc", Paths: []string{"t3/b", "t3/a"}},      // 1 reclaimable
		{Size: 3, Digest: "bb", Paths: []string{"s/2", "s/1", "s/3"}}, // 6 reclaimable
		{Size: 5, Digest: "aa", Paths: []string{"big/b", "big/a"}},    // 5 reclaimable
		{Size: 1, Digest: "dd", Paths: []string{"t1/a", "t1/b"}},      // 1 reclaimable
	}

	got := Order(groups)
	wantDigests := []string{"bb", "aa", "cc", "dd"}
	wantPaths := [][]string{
		{"s/1", "s/2", "s/3"},
		{"big/a", "big/b"},
		{"t3/a", "t3/b"},
		{"t1/a", "t1/b"},
	}
	for i, g := range got {
		if g.Digest != wantDigests[i] {
			t.Errorf("group %d digest = %s, want %s", i+1, g.Digest, wantDigests[i])
		}
		if g.Number != i+1 {
			t.Errorf("group %d number = %d, want %d", i+1, g.Number, i+1)
		}
		if !reflect.DeepEqual(g.Paths, wantPaths[i]) {
			t.Errorf("group %d paths = %v, want %v", i+1, g.Paths, wantPaths[i])
		}
	}
}

func TestOrderPutsReclaimableBytesFirst(t *testing.T) {
	small := Group{Size: 3, Digest: "bb", Paths: []string{"s1", "s2", "s3", "s4", "s5"}} // 12
	large := Group{Size: 10, Digest: "aa", Paths: []string{"l1", "l2"}}                  // 10

	got := Order([]Group{large, small})
	if got[0].Digest != "bb" || got[1].Digest != "aa" {
		t.Errorf("order = %s, %s; want bb then aa", got[0].Digest, got[1].Digest)
	}
}

func TestOrderBreaksReclaimableTiesByDigest(t *testing.T) {
	later := Group{Size: 10, Digest: "dd", Paths: []string{"a", "b"}}       // 10
	earlier := Group{Size: 5, Digest: "cc", Paths: []string{"c", "d", "e"}} // 10

	got := Order([]Group{later, earlier})
	if got[0].Digest != "cc" || got[1].Digest != "dd" {
		t.Errorf("order = %s, %s; want cc then dd", got[0].Digest, got[1].Digest)
	}
}

func TestOrderSortsPathsByByteValue(t *testing.T) {
	g := Group{Size: 1, Digest: "x", Paths: []string{"b.txt", "B.txt", "a b.txt", "a/z.txt", ".hid/x.txt", "a.txt"}}

	got := Order([]Group{g})
	want := []string{".hid/x.txt", "B.txt", "a b.txt", "a.txt", "a/z.txt", "b.txt"}
	if !reflect.DeepEqual(got[0].Paths, want) {
		t.Errorf("paths = %v, want %v", got[0].Paths, want)
	}
}

func TestOrderNumbersGroupsFromOne(t *testing.T) {
	groups := []Group{
		{Size: 1, Digest: "a", Paths: []string{"p", "q"}},
		{Size: 2, Digest: "b", Paths: []string{"r", "s"}},
		{Size: 3, Digest: "c", Paths: []string{"t", "u"}},
	}

	got := Order(groups)
	for i, g := range got {
		if g.Number != i+1 {
			t.Errorf("group %d number = %d, want %d", i, g.Number, i+1)
		}
	}
}

func TestOrderEmptyReturnsEmpty(t *testing.T) {
	if got := Order(nil); len(got) != 0 {
		t.Errorf("order = %v, want empty", got)
	}
}
