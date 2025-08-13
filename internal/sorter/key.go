package sorter

import (
	"math"
	"strconv"
	"strings"
)

// key представляет предобработанное представление строки для сравнения.
// Поля взаимно исключают доминирующие режимы сравнения (month, human, num),
// но строковый ключ присутствует всегда для fallback и вторичных сравнений.
type key struct {
	month int     // 1..12, 0 если не распознано/не используется
	human float64 // значение для -h, NaN если не используется/ошибка
	num   float64 // значение для -n, NaN если не используется/ошибка
	text  string  // строковый ключ (возможна обрезка хвостовых пробелов)
}

// buildKeyspace создает срез ключей для сортировки на основе строк lines и конфигурации cfg.
// Для каждой строки формируется ключ key с заполнением полей month, human, num и text
// в зависимости от включенных режимов сортировки. Колонки выбираются через cfg.Column и cfg.Delimiter.
// Если cfg.IgnoreTrailWS=true, хвостовые пробелы и табуляции удаляются.
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

// compareKeys сравнивает два ключа a и b и возвращает:
// -1, если a < b
//
//	1, если a > b
//	0, если a == b
//
// Приоритет сравнения: month > human > num > text. Неопределенные/NaN значения отправляются в конец.
func compareKeys(a, b key) int {
	// Приоритет: month > human > num > text
	if a.month != 0 || b.month != 0 {
		am := a.month
		bm := b.month
		if am == 0 && bm != 0 {
			return 1 // нераспознанные в конец
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

// column извлекает указанную колонку col из строки s, используя разделитель delim.
// Колонки нумеруются с 1. Если колонка не найдена, возвращается пустая строка.
// Для col<=1 возвращается либо вся строка, либо первая колонка до разделителя.
func column(s string, col int, delim string) string {
	if col <= 1 {
		if col == 1 {
			// первая колонка
			if i := strings.Index(s, delim); i >= 0 {
				return s[:i]
			}
		}

		return s
	}

	// Быстрый проход без alloc: считаем разделители
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

// parseNumber пытается интерпретировать строку s как число с плавающей точкой.
// Игнорируются пробелы вокруг числа. Если парсинг не удался или строка пустая, возвращается NaN и false.
// Иначе возвращается число и true.
func parseNumber(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return math.NaN(), false
	}

	// Используем ParseFloat для универсальности
	// (десятичная точка, экспоненты поддерживаются).
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}

	return v, true
}
