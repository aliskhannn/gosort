package sorter

import (
	"math"
	"strconv"
	"strings"
)

// parseHumanSize парсит значения вида 10K, 2M, 3G, 4T, 5P, 6E (в двоичных степенях 1024).
// Допускается суффикс 'B' (например, 10KB). Регистр суффикса не важен. Пробелы вокруг числа игнорируются.
// Возвращает (значение, true) при успехе; иначе (NaN, false).
func parseHumanSize(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return math.NaN(), false
	}

	// Отделяем конечные буквы (до 2-х: например, KB)
	base := s
	suf := ""

	for i := len(s) - 1; i >= 0; i-- {
		if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') {
			continue
		}

		base = s[:i+1]
		suf = s[i+1:]
		break
	}

	if base == s {
		// Нет букв — суффикс пустой
		suf = ""
	}

	val, err := strconv.ParseFloat(base, 64)
	if err != nil {
		return math.NaN(), false
	}

	suf = strings.ToUpper(suf)
	if strings.HasSuffix(suf, "B") && len(suf) > 1 {
		suf = strings.TrimSuffix(suf, "B")
	}

	multipliers := map[string]float64{
		"":  1,
		"K": 1024,
		"M": 1024 * 1024,
		"G": 1024 * 1024 * 1024,
		"T": math.Pow(1024, 4),
		"P": math.Pow(1024, 5),
		"E": math.Pow(1024, 6),
	}

	mul, ok := multipliers[suf]
	if !ok {
		return math.NaN(), false // Неизвестный суффикс
	}

	return val * mul, true
}
