package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/aliskhannn/gosort/internal/sorter"
)

func main() {
	// Define command-line flags.
	var (
		flagCol    = flag.IntP("k", "k", 0, "Column number (1..N) to sort by; default is the whole line. Delimiter is tab.")
		flagNum    = flag.BoolP("numeric", "n", false, "Numeric sort (compare by numeric value).")
		flagRev    = flag.BoolP("reverse", "r", false, "Reverse the sorting order.")
		flagUniq   = flag.BoolP("unique", "u", false, "Output only unique lines (after sorting).")
		flagMonth  = flag.BoolP("month-sort", "M", false, "Sort by month name (Jan..Dec).")
		flagTrimTB = flag.BoolP("ignore-tb", "b", false, "Ignore trailing spaces when comparing.")
		flagCheck  = flag.BoolP("check", "c", false, "Check whether the input is sorted; do not sort.")
		flagHuman  = flag.BoolP("human-numeric", "h", false, "Numeric sort with suffixes (K, M, G, T, P, E).")
	)

	// Custom usage/help message.
	flag.Usage = func() {
		_, _ = fmt.Fprint(os.Stderr, "gosort — simplified sort\n")
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [FILE ...]\n\n", os.Args[0])
		flag.PrintDefaults()
		_, _ = fmt.Fprintln(os.Stderr, "\nIf FILE is not specified, input is read from STDIN.")
	}

	flag.Parse()

	// Validate incompatible flags.
	if *flagNum && *flagHuman {
		_, _ = fmt.Fprintln(os.Stderr, "Warning: both -n and -h specified; using -h (human-readable numbers).")
	}
	if *flagMonth && (*flagNum || *flagHuman) {
		_, _ = fmt.Fprintln(os.Stderr, "Warning: -M is incompatible with -n/-h; using -M.")
	}

	// Read input from files or STDIN.
	inputs := flag.Args()
	lines, err := readAll(inputs)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Read error:", err)
		os.Exit(1)
	}

	// Configure sort options.
	cfg := sorter.Config{
		Column:        *flagCol,
		Reverse:       *flagRev,
		Unique:        *flagUniq,
		IgnoreTrailWS: *flagTrimTB,
		Numeric:       *flagNum,
		HumanNumeric:  *flagHuman,
		Month:         *flagMonth,
		Delimiter:     "\t",
	}

	// Check mode (-c).
	if *flagCheck {
		var ok bool
		var i int
		var a, b string

		if ok, i, a, b = sorter.IsSorted(lines, cfg); ok {
			// Compatible with GNU sort: produce no output, exit code 0
			return
		}
		_, _ = fmt.Fprintf(os.Stderr, "sort: disorder at line %d: %q > %q\n", i, a, b)
		os.Exit(1)
	}

	// Perform sorting.
	out, err := sorter.Sort(lines, cfg)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Sort error:", err)
		os.Exit(1)
	}

	// Print result to STDOUT with line breaks.
	w := bufio.NewWriter(os.Stdout)
	for i, s := range out {
		if i > 0 {
			_, _ = w.WriteString("\n")
		}
		_, _ = w.WriteString(s)
	}

	_, _ = w.WriteString("\n")
	_ = w.Flush()
}

// readAll reads lines from a list of files.
// If files are empty, reads from STDIN.
// Special name "-" is supported for STDIN.
func readAll(files []string) ([]string, error) {
	if len(files) == 0 {
		return readFrom(os.Stdin)
	}

	var all []string
	for _, path := range files {
		var r io.ReadCloser

		if path == "-" {
			r = os.Stdin
		} else {
			f, err := os.Open(path)
			if err != nil {
				return nil, err
			}

			r = f
		}

		lines, err := readFrom(r)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}

		all = append(all, lines...)
		if r != os.Stdin {
			_ = r.Close()
		}
	}

	return all, nil
}

// readFrom reads all lines from io.Reader and returns a slice of strings.
// Carriage returns (\r) are stripped.
// If the input is empty, returns "no input" error.
func readFrom(r io.Reader) ([]string, error) {
	// Scanner with an increased buffer size to handle very long lines.
	s := bufio.NewScanner(r)

	const maxCap = 10 * 1024 * 1024 // 10MB per line
	s.Buffer(make([]byte, 64*1024), maxCap)

	var res []string
	for s.Scan() {
		res = append(res, strings.TrimRight(s.Text(), "\r"))
	}

	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("no input")
	}

	return res, nil
}
