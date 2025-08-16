package sorter

import (
	"math"
	"strconv"
	"strings"
)

// key represents a preprocessed form of a string for comparison.
// Fields are mutually exclusive for dominant comparison modes (month, human, num),
// but the text field is always present as a fallback and for secondary comparisons.
type key struct {
	month int     // 1..12, 0 if not recognized/not used
	human float64 // value for -h, NaN if not used/error
	num   float64 // value for -n, NaN if not used/error
	text  string  // string key (trailing spaces may be trimmed)
}

// buildKeyspace builds a slice of keys for sorting based on lines and cfg.
// For each line, a key is constructed filling month, human, num, and text
// depending on enabled modes. Columns are extracted via cfg.Column and cfg.Delimiter.
// If cfg.IgnoreTrailWS=true, trailing spaces and tabs are trimmed.
func buildKeyspace(lines []string, cfg Config) []key {
	ks := make([]key, len(lines))

	for i, s := range lines {
		col := s

		if cfg.Column > 0 {
			col = column(s, cfg.Column, cfg.Delimiter)
		}
		if cfg.IgnoreTrailWS {
			col = strings.TrimRight(col, " \t")
		}

		k := key{text: col}
		if cfg.Month {
			k.month = parseMonth(col)
		}

		if cfg.HumanNumeric {
			if v, ok := parseHumanSize(col); ok {
				k.human = v
			} else {
				k.human = math.NaN()
			}
		} else if cfg.Numeric {
			if v, ok := parseNumber(col); ok {
				k.num = v
			} else {
				k.num = math.NaN()
			}
		}

		ks[i] = k
	}

	return ks
}

// compareKeys compares two keys a and b and returns:
// -1 if a < b
//
//	1 if a > b
//	0 if a == b
//
// Priority order: month > human > num > text.
// Undefined/NaN values are placed at the end.
func compareKeys(a, b key) int {
	// Priority: month > human > num > text.
	if a.month != 0 || b.month != 0 {
		am := a.month
		bm := b.month
		if am == 0 && bm != 0 {
			return 1 // unrecognized to the end
		}
		if am != 0 && bm == 0 {
			return -1
		}
		if am < bm {
			return -1
		}
		if am > bm {
			return 1
		}
	}
	if !math.IsNaN(a.human) || !math.IsNaN(b.human) {
		ah := a.human
		bh := b.human
		if math.IsNaN(ah) && !math.IsNaN(bh) {
			return 1
		}
		if !math.IsNaN(ah) && math.IsNaN(bh) {
			return -1
		}
		if ah < bh {
			return -1
		}
		if ah > bh {
			return 1
		}
	}
	if !math.IsNaN(a.num) || !math.IsNaN(b.num) {
		an := a.num
		bn := b.num
		if math.IsNaN(an) && !math.IsNaN(bn) {
			return 1
		}
		if !math.IsNaN(an) && math.IsNaN(bn) {
			return -1
		}
		if an < bn {
			return -1
		}
		if an > bn {
			return 1
		}
	}
	if a.text < b.text {
		return -1
	}
	if a.text > b.text {
		return 1
	}
	return 0
}

// column extracts the col-th column from string s, using delim as separator.
// Columns are numbered starting from 1. If not found, returns "".
// For col<=1 it returns either the whole line or the first column up to the delimiter.
func column(s string, col int, delim string) string {
	if col <= 1 {
		if col == 1 {
			// First column.
			if i := strings.Index(s, delim); i >= 0 {
				return s[:i]
			}
		}

		return s
	}

	// Quick pass without alloc: counting separators.
	start := 0
	seen := 1

	for {
		idx := strings.Index(s[start:], delim)
		if idx < 0 {
			if seen == col {
				return s[start:]
			}

			return ""
		}

		if seen == col {
			return s[start : start+idx]
		}

		start += idx + len(delim)
		seen++
	}
}

// parseNumber tries to parse s as a floating-point number.
// Leading/trailing spaces are ignored. On failure or empty string, returns NaN and false.
// Otherwise returns the number and true.
func parseNumber(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return math.NaN(), false
	}

	// Use ParseFloat for versatility
	// (decimal point, exponents supported).
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}

	return v, true
}
