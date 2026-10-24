// Package cli implements dupfind's command-line contract.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"dupfind/internal/output"
	"dupfind/internal/scan"
)

// Run executes dupfind with args and writes command output to stdout and stderr.
func Run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("dupfind", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	summary := fs.Bool("summary", false, "write duplicate totals instead of a file listing")
	minText := fs.String("min-size", "1", "minimum file size in bytes (positive integer)")
	printHelp := func() {
		fmt.Fprintln(stdout, "dupfind finds exact duplicate regular files in a directory tree.")
		fmt.Fprintln(stdout, "Usage: dupfind [--summary] [--min-size <bytes>] <directory>")
		fmt.Fprintln(stdout, "Options:")
		fmt.Fprintln(stdout, "  --summary           write duplicate group and reclaimable-space totals")
		fmt.Fprintln(stdout, "  --min-size <bytes>  ignore files smaller than this positive size (default: 1)")
	}
	fs.Usage = func() {}
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printHelp()
			return 0
		}
		fmt.Fprintln(stderr, "error: invalid command options")
		return 2
	}
	minSize, err := strconv.ParseInt(*minText, 10, 64)
	if err != nil || minSize <= 0 {
		fmt.Fprintln(stderr, "error: --min-size must be a positive integer")
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "error: exactly one directory argument is required")
		return 2
	}
	root := fs.Arg(0)
	info, err := os.Stat(root)
	if err != nil {
		fmt.Fprintf(stderr, "error: cannot access directory %s: %v\n", root, err)
		return 1
	}
	if !info.IsDir() {
		fmt.Fprintf(stderr, "error: path is not a directory: %s\n", root)
		return 1
	}
	groups, err := scan.Scan(root, minSize, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "error: cannot scan directory %s: %v\n", root, err)
		return 1
	}
	if *summary {
		if err := output.Summary(stdout, groups); err != nil {
			fmt.Fprintf(stderr, "error: cannot write output: %v\n", err)
			return 1
		}
	} else if err := output.Listing(stdout, groups); err != nil {
		fmt.Fprintf(stderr, "error: cannot write output: %v\n", err)
		return 1
	}
	return 0
}
