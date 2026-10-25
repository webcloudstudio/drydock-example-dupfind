// Package cli implements the dupfind command line: argument parsing, directory
// validation, pipeline orchestration, and exit codes.
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"dupfind/internal/dupes"
	"dupfind/internal/report"
	"dupfind/internal/walk"
)

// Config is the parsed command line.
type Config struct {
	Summary bool
	MinSize int64
	Dir     string
}

const usage = "usage: dupfind [--summary] [--min-size <bytes>] <directory>"

// Run executes dupfind with args and returns the process exit code: 0 on success,
// 1 on operational failure, 2 on usage error. Standard output carries only CSV.
func Run(args []string, stdout, stderr io.Writer) int {
	cfg, err := parseArgs(args)
	if err != nil {
		fmt.Fprintf(stderr, "dupfind: %v; %s\n", err, usage)
		return 2
	}
	root, err := validateDirectory(cfg.Dir)
	if err != nil {
		fmt.Fprintf(stderr, "dupfind: %v\n", err)
		return 1
	}
	entries, err := walk.Files(root, cfg.MinSize, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "dupfind: %v\n", err)
		return 1
	}
	groups := dupes.Order(dupes.GroupByDigest(root, dupes.GroupBySize(entries), stderr))
	write := report.WriteListing
	if cfg.Summary {
		write = report.WriteSummary
	}
	if err := write(stdout, groups); err != nil {
		fmt.Fprintf(stderr, "dupfind: %v\n", err)
		return 1
	}
	return 0
}

// parseArgs interprets [--summary] [--min-size <bytes>] <directory>. Options must precede the
// directory; any other positional argument counts toward the single-directory requirement.
func parseArgs(args []string) (Config, error) {
	fs := flag.NewFlagSet("dupfind", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var cfg Config
	fs.BoolVar(&cfg.Summary, "summary", false, "")
	minSize := fs.String("min-size", "1", "")
	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}
	if fs.NArg() != 1 {
		return Config{}, fmt.Errorf("expected exactly one directory argument, got %d", fs.NArg())
	}
	n, err := parseMinSize(*minSize)
	if err != nil {
		return Config{}, err
	}
	cfg.MinSize = n
	cfg.Dir = fs.Arg(0)
	return cfg, nil
}

// parseMinSize accepts only a string of ASCII decimal digits whose value is at least 1 and fits
// in an int64. Signs, zero, hexadecimal, exponent forms, and overflow are rejected.
func parseMinSize(s string) (int64, error) {
	invalid := fmt.Errorf("--min-size must be a positive decimal integer, got %q", s)
	if s == "" {
		return 0, invalid
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, invalid
		}
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil || n < 1 {
		return 0, invalid
	}
	return n, nil
}

// validateDirectory confirms that dir names an existing directory and returns its path with
// symbolic links resolved, so that a link to a directory can be walked.
func validateDirectory(dir string) (string, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", dir)
	}
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return "", err
	}
	return resolved, nil
}
