package sorter

import "strings"

var months = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

// parseMonth распознаёт трёхбуквенное имя месяца (регистронезависимо).
func parseMonth(s string) int {
	s = strings.ToLower(s)
	if len(s) < 3 {
		return 0
	}

	key := strings.ToLower(s[:3])
	if m, ok := months[key]; ok {
		return m
	}

	return 0
}
