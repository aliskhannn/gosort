package sorter

import (
	"sort"
	"strings"
)

// Config defines sorting options.
// Column: column index (1..N). 0 means sort by the whole line.
// Delimiter: column separator (default tab).
// Numeric: compare as numbers; HumanNumeric: numbers with suffixes (K/M/...).
// Month: compare as months (Jan..Dec). Reverse: reverse order.
// Unique: remove duplicates after sorting. IgnoreTrailWS: trim trailing spaces.
// Key priority: Month > HumanNumeric > Numeric > Text.
// If conflicting flags are set, priority is resolved in KeyFrom.
type Config struct {
	Column        int
	Delimiter     string
	Numeric       bool
	HumanNumeric  bool
	Month         bool
	Reverse       bool
	Unique        bool
	IgnoreTrailWS bool
}

// Sort returns a sorted copy of lines according to cfg.
func Sort(lines []string, cfg Config) ([]string, error) {
	if len(lines) == 0 {
		return []string{}, nil // empty entrance
	}

	ks := buildKeyspace(lines, cfg)
	idx := make([]int, len(lines))

	for i := range idx {
		idx[i] = i
	}

	sort.Slice(idx, func(i, j int) bool {
		li, lj := idx[i], idx[j]
		cmp := compareKeys(ks[li], ks[lj])

		if cfg.Reverse {
			cmp = -cmp
		}

		if cmp == 0 {
			// For stability — secondary key on the original index.
			return li < lj
		}

		return cmp < 0
	})

	out := make([]string, 0, len(lines))
	var last *key

	for _, id := range idx {
		k := ks[id]
		line := lines[id]

		// Trim trailing spaces for -b.
		if cfg.IgnoreTrailWS {
			line = strings.TrimRight(line, " \t")
		}

		if cfg.Unique {
			if last != nil && compareKeys(*last, k) == 0 {
				continue
			}
			copyK := k
			last = &copyK
		}

		out = append(out, line)
	}

	return out, nil
}

// IsSorted checks if lines are sorted according to cfg.
// Returns ok, index (1-based) of the first error, and the left/right values for reporting.
func IsSorted(lines []string, cfg Config) (bool, int, string, string) {
	if len(lines) <= 1 {
		return true, 0, "", ""
	}

	ks := buildKeyspace(lines, cfg)
	prev := ks[0]

	for i := 1; i < len(ks); i++ {
		cmp := compareKeys(prev, ks[i])

		if cfg.Reverse {
			cmp = -cmp
		}

		if cmp > 0 || (cfg.Unique && cmp == 0) {
			return false, i + 1, lines[i-1], lines[i]
		}

		prev = ks[i]
	}

	return true, 0, "", ""
}
