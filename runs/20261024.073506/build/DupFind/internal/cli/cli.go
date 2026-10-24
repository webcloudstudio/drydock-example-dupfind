package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"dupfind/internal/discovery"
	"dupfind/internal/output"
)

// Run executes dupfind with caller-supplied streams.
func Run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("dupfind", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	summary := fs.Bool("summary", false, "write summary CSV")
	minSize := fs.Int64("min-size", 1, "minimum file size in bytes")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, "usage: dupfind [--summary] [--min-size bytes] <directory>")
		return 2
	}
	if *minSize <= 0 || fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: dupfind [--summary] [--min-size bytes] <directory>")
		return 2
	}
	root := fs.Arg(0)
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		fmt.Fprintf(stderr, "error: %s is not an accessible directory\n", root)
		return 1
	}
	groups, err := discovery.Find(root, *minSize, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "error: cannot scan directory: %v\n", err)
		return 1
	}
	if *summary {
		if err := output.WriteSummary(stdout, groups); err != nil {
			fmt.Fprintf(stderr, "error: cannot write output: %v\n", err)
			return 1
		}
		return 0
	}
	if err := output.WriteListing(stdout, groups); err != nil {
		fmt.Fprintf(stderr, "error: cannot write output: %v\n", err)
		return 1
	}
	return 0
}
