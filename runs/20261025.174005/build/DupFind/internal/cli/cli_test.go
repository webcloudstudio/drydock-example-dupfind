package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// runCLI executes Run and returns its exit code and both output streams.
func runCLI(args ...string) (int, string, string) {
	var stdout, stderr bytes.Buffer
	code := Run(args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

// writeFiles creates each named file under root with the given contents.
func writeFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for name, body := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// parseCSV reads every record of text with encoding/csv.
func parseCSV(t *testing.T, text string) [][]string {
	t.Helper()
	records, err := csv.NewReader(strings.NewReader(text)).ReadAll()
	if err != nil {
		t.Fatalf("stdout is not valid CSV: %v\n%s", err, text)
	}
	return records
}

// sha256Hex returns the lowercase hexadecimal SHA-256 of body.
func sha256Hex(body string) string {
	sum := sha256.Sum256([]byte(body))
	return hex.EncodeToString(sum[:])
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

// snapshot maps every path under root to its contents, with directories mapped to "" and marked
// by a trailing slash.
func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, p)
		if d.IsDir() {
			out[filepath.ToSlash(rel)+"/"] = ""
			return nil
		}
		body, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		out[filepath.ToSlash(rel)] = string(body)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestRunUsageErrorsExitTwoWithEmptyStdout(t *testing.T) {
	dir := t.TempDir()
	other := filepath.Join(dir, "other")
	if err := os.Mkdir(other, 0o755); err != nil {
		t.Fatal(err)
	}
	cases := map[string][]string{
		"no arguments":                 {},
		"option without directory":     {"--summary"},
		"min-size without directory":   {"--min-size", "5"},
		"two directories":              {dir, other},
		"unknown option":               {"--frobnicate", dir},
		"short help":                   {"-h", dir},
		"long help":                    {"--help", dir},
		"min-size non-numeric":         {"--min-size", "abc", dir},
		"min-size zero":                {"--min-size", "0", dir},
		"min-size negative":            {"--min-size", "-5", dir},
		"min-size plus sign":           {"--min-size", "+5", dir},
		"min-size hexadecimal":         {"--min-size", "0x10", dir},
		"min-size exponent":            {"--min-size", "1e3", dir},
		"min-size overflow":            {"--min-size", "99999999999999999999999", dir},
		"min-size consumes directory":  {"--min-size", dir},
		"option after directory":       {dir, "--summary"},
		"option after directory, size": {dir, "--min-size", "5"},
	}
	for name, args := range cases {
		t.Run(name, func(t *testing.T) {
			code, out, errOut := runCLI(args...)
			if code != 2 {
				t.Errorf("exit = %d, want 2 (stderr %q)", code, errOut)
			}
			if out != "" {
				t.Errorf("stdout = %q, want empty", out)
			}
			if n := strings.Count(errOut, "\n"); n != 1 {
				t.Errorf("stderr lines = %d, want 1: %q", n, errOut)
			}
		})
	}
}

func TestRunValidFormsAreNotUsageErrors(t *testing.T) {
	dir := t.TempDir()
	cases := map[string][]string{
		"directory only":            {dir},
		"summary first":             {"--summary", dir},
		"min-size before directory": {"--min-size", "5", dir},
		"summary and min-size":      {"--summary", "--min-size", "5", dir},
		"min-size then summary":     {"--min-size", "5", "--summary", dir},
		"min-size with equals sign": {"--min-size=5", dir},
	}
	for name, args := range cases {
		t.Run(name, func(t *testing.T) {
			if code, _, errOut := runCLI(args...); code != 0 {
				t.Errorf("exit = %d, want 0 (stderr %q)", code, errOut)
			}
		})
	}
}

func TestRunMissingDirectoryExitsOneWithEmptyStdout(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "no-such-directory")
	for _, args := range [][]string{{missing}, {"--summary", missing}, {"--min-size", "3", missing}} {
		code, out, errOut := runCLI(args...)
		if code != 1 {
			t.Errorf("%v: exit = %d, want 1", args, code)
		}
		if out != "" {
			t.Errorf("%v: stdout = %q, want empty", args, out)
		}
		if n := strings.Count(errOut, "\n"); n != 1 {
			t.Errorf("%v: stderr lines = %d, want 1", args, n)
		}
	}
}

func TestRunRegularFileExitsOneWithEmptyStdout(t *testing.T) {
	path := filepath.Join(t.TempDir(), "plain.txt")
	if err := os.WriteFile(path, []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{path}, {"--summary", path}} {
		code, out, _ := runCLI(args...)
		if code != 1 {
			t.Errorf("%v: exit = %d, want 1", args, code)
		}
		if out != "" {
			t.Errorf("%v: stdout = %q, want empty", args, out)
		}
	}
}

func TestRunUsageErrorTakesPrecedenceOverMissingDirectory(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "absent")

	if code, _, _ := runCLI("--min-size", "0", missing); code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

func TestRunEmptyDirectorySucceedsWithHeaderOnly(t *testing.T) {
	code, out, _ := runCLI(t.TempDir())
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if got := parseCSV(t, out); !reflect.DeepEqual(got, [][]string{{"group", "size", "sha256", "path"}}) {
		t.Errorf("listing = %v, want header only", got)
	}
}

func TestRunListingReportsEveryDuplicateWithCorrectColumns(t *testing.T) {
	root := t.TempDir()
	content := "quarterly report\n"
	writeFiles(t, root, map[string]string{"a/one.txt": content, "b/two.txt": content, "c/deep/three.txt": content})

	code, out, errOut := runCLI(root)
	if code != 0 {
		t.Fatalf("exit = %d (stderr %q)", code, errOut)
	}
	got := parseCSV(t, out)
	if len(got) != 4 {
		t.Fatalf("rows = %d, want header plus 3", len(got))
	}
	var paths []string
	for _, rec := range got[1:] {
		if rec[0] != "1" || rec[1] != strconv.Itoa(len(content)) || rec[2] != sha256Hex(content) {
			t.Errorf("row = %v, want group 1, size %d, digest of content", rec, len(content))
		}
		paths = append(paths, rec[3])
	}
	if want := []string{"a/one.txt", "b/two.txt", "c/deep/three.txt"}; !reflect.DeepEqual(paths, want) {
		t.Errorf("paths = %v, want %v", paths, want)
	}
}

func TestRunListingOrdersGroupsByReclaimableBytesThenDigest(t *testing.T) {
	root := t.TempDir()
	big := "0123456789"
	small := "xyz"
	writeFiles(t, root, map[string]string{
		"big/a": big, "big/b": big,
		"s/1": small, "s/2": small, "s/3": small, "s/4": small, "s/5": small,
	})

	code, out, _ := runCLI(root)
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	rows := parseCSV(t, out)[1:]
	if len(rows) != 7 {
		t.Fatalf("rows = %d, want 7", len(rows))
	}
	// small: 3 bytes x 4 reclaimable copies = 12; big: 10 x 1 = 10.
	if rows[0][0] != "1" || rows[0][2] != sha256Hex(small) {
		t.Errorf("first row = %v, want group 1 for the small content", rows[0])
	}
	if rows[5][0] != "2" || rows[5][2] != sha256Hex(big) {
		t.Errorf("row 5 = %v, want group 2 for the big content", rows[5])
	}
}

func TestRunListingPathsAreSortedByByteValue(t *testing.T) {
	root := t.TempDir()
	content := "order me\n"
	names := []string{"b.txt", "B.txt", "a b.txt", "a/z.txt", ".hid/x.txt", "a.txt"}
	files := map[string]string{}
	for _, n := range names {
		files[n] = content
	}
	writeFiles(t, root, files)

	code, out, _ := runCLI(root)
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	var got []string
	for _, rec := range parseCSV(t, out)[1:] {
		got = append(got, rec[3])
	}
	want := append([]string(nil), names...)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("paths = %v, want %v", got, want)
	}
}

func TestRunListingHeaderOnlyWhenNoDuplicates(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{"solo.txt": "only one of me\n", "empty/one": "", "empty/two": ""})

	code, out, _ := runCLI(root)
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if got := parseCSV(t, out); !reflect.DeepEqual(got, [][]string{{"group", "size", "sha256", "path"}}) {
		t.Errorf("listing = %v, want header only", got)
	}
}

func TestRunListingQuotesCommaAndQuotePaths(t *testing.T) {
	root := t.TempDir()
	content := "quote me\n"
	paths := []string{"song,live.mp3", `say "hi".txt`, "plain.txt"}
	files := map[string]string{}
	for _, p := range paths {
		files[p] = content
	}
	writeFiles(t, root, files)

	code, out, _ := runCLI(root)
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	var got []string
	for _, rec := range parseCSV(t, out)[1:] {
		if len(rec) != 4 {
			t.Fatalf("row has %d fields, want 4: %v", len(rec), rec)
		}
		got = append(got, rec[3])
	}
	sort.Strings(paths)
	sort.Strings(got)
	if !reflect.DeepEqual(got, paths) {
		t.Errorf("paths round-trip = %v, want %v", got, paths)
	}
}

func TestRunListingIncludesHiddenAndExcludesSymlinks(t *testing.T) {
	root := t.TempDir()
	content := "KEY=1\n"
	writeFiles(t, root, map[string]string{".hidden/.env": content, "env.bak": content})
	if err := os.Symlink(filepath.Join(root, "env.bak"), filepath.Join(root, "link.bak")); err != nil {
		t.Skip("symbolic links not supported here:", err)
	}
	if err := os.Symlink(filepath.Join(root, ".hidden"), filepath.Join(root, "linkdir")); err != nil {
		t.Skip("symbolic links not supported here:", err)
	}

	code, out, _ := runCLI(root)
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	var got []string
	for _, rec := range parseCSV(t, out)[1:] {
		got = append(got, rec[3])
	}
	if want := []string{".hidden/.env", "env.bak"}; !reflect.DeepEqual(got, want) {
		t.Errorf("paths = %v, want %v", got, want)
	}
}

func TestRunSymlinkedDirectoryArgumentIsWalked(t *testing.T) {
	real := t.TempDir()
	content := "same\n"
	writeFiles(t, real, map[string]string{"a": content, "b": content})
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(real, link); err != nil {
		t.Skip("symbolic links not supported here:", err)
	}

	code, out, errOut := runCLI(link)
	if code != 0 {
		t.Fatalf("exit = %d (stderr %q)", code, errOut)
	}
	if rows := parseCSV(t, out); len(rows) != 3 {
		t.Errorf("rows = %d, want header plus 2", len(rows))
	}
}

func TestRunMinSizeExcludesSmallFilesFromListing(t *testing.T) {
	root := t.TempDir()
	large := strings.Repeat("L", 50)
	writeFiles(t, root, map[string]string{"s1": "tiny", "s2": "tiny", "l1": large, "l2": large})

	code, out, _ := runCLI("--min-size", "5", root)
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	var got []string
	for _, rec := range parseCSV(t, out)[1:] {
		got = append(got, rec[3])
	}
	if want := []string{"l1", "l2"}; !reflect.DeepEqual(got, want) {
		t.Errorf("paths = %v, want %v", got, want)
	}
}

func TestRunSummaryTotalsMatchFixture(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"a/r.txt": "quarterly report\n", "b/r2.txt": "quarterly report\n", "c/r3.txt": "quarterly report\n",
		"m/s.mp3": "la la la\n", "m/s2.mp3": "la la la\n",
		"u/solo": "unique content here\n",
	})

	code, out, errOut := runCLI("--summary", root)
	if code != 0 {
		t.Fatalf("exit = %d (stderr %q)", code, errOut)
	}
	want := [][]string{
		{"groups", "duplicate_files", "reclaimable_bytes"},
		{"2", "3", strconv.Itoa(17*2 + 9*1)},
	}
	if got := parseCSV(t, out); !reflect.DeepEqual(got, want) {
		t.Errorf("summary = %v, want %v", got, want)
	}
}

func TestRunSummaryNoDuplicatesIsZeroRow(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{"empty/one": "", "empty/two": ""})

	code, out, _ := runCLI("--summary", root)
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	want := [][]string{{"groups", "duplicate_files", "reclaimable_bytes"}, {"0", "0", "0"}}
	if got := parseCSV(t, out); !reflect.DeepEqual(got, want) {
		t.Errorf("summary = %v, want %v", got, want)
	}
}

func TestRunSummaryAgreesWithListingForEveryMinSize(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"s1": strings.Repeat("S", 4), "s2": strings.Repeat("S", 4),
		"m1": strings.Repeat("M", 30), "m2": strings.Repeat("M", 30), "m3": strings.Repeat("M", 30),
		"b1": strings.Repeat("B", 200), "b2": strings.Repeat("B", 200),
	})

	for _, opts := range [][]string{{}, {"--min-size", "20"}, {"--min-size", "100"}, {"--min-size", "1000"}} {
		code, listing, _ := runCLI(append(append([]string{}, opts...), root)...)
		if code != 0 {
			t.Fatalf("%v listing exit = %d", opts, code)
		}
		groups := map[string][]int{}
		for _, rec := range parseCSV(t, listing)[1:] {
			size, _ := strconv.Atoi(rec[1])
			groups[rec[0]] = append(groups[rec[0]], size)
		}
		wantFiles, wantBytes := 0, 0
		for _, sizes := range groups {
			wantFiles += len(sizes) - 1
			wantBytes += sizes[0] * (len(sizes) - 1)
		}
		want := [][]string{
			{"groups", "duplicate_files", "reclaimable_bytes"},
			{strconv.Itoa(len(groups)), strconv.Itoa(wantFiles), strconv.Itoa(wantBytes)},
		}

		code, summary, _ := runCLI(append(append([]string{"--summary"}, opts...), root)...)
		if code != 0 {
			t.Fatalf("%v summary exit = %d", opts, code)
		}
		if got := parseCSV(t, summary); !reflect.DeepEqual(got, want) {
			t.Errorf("%v: summary = %v, want %v", opts, got, want)
		}
	}
}

func TestRunDoesNotModifyTree(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"a/f1": "keep me\n", "b/f2": "keep me\n", "c/other": "different!",
		"d/.hidden": "keep me\n", "e/empty": "",
	})
	before := snapshot(t, root)

	for _, args := range [][]string{{root}, {"--summary", root}, {"--min-size", "3", root}} {
		if code, _, errOut := runCLI(args...); code != 0 {
			t.Fatalf("%v exit = %d (stderr %q)", args, code, errOut)
		}
	}
	if after := snapshot(t, root); !reflect.DeepEqual(after, before) {
		t.Errorf("tree changed:\nbefore %v\nafter  %v", before, after)
	}
}

func TestRunUnreadableDuplicateIsSkippedWithoutChangingExitCode(t *testing.T) {
	root := t.TempDir()
	content := "same bytes\n"
	writeFiles(t, root, map[string]string{"ok1.txt": content, "ok2.txt": content, "locked.txt": content})
	if !lockFile(t, filepath.Join(root, "locked.txt")) {
		t.Skip("file permissions are not enforced here")
	}

	code, out, errOut := runCLI(root)
	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr %q)", code, errOut)
	}
	var got []string
	for _, rec := range parseCSV(t, out)[1:] {
		got = append(got, rec[3])
	}
	if want := []string{"ok1.txt", "ok2.txt"}; !reflect.DeepEqual(got, want) {
		t.Errorf("paths = %v, want %v", got, want)
	}
	if n := strings.Count(errOut, "\n"); n != 1 {
		t.Errorf("warning lines = %d, want 1: %q", n, errOut)
	}
}
