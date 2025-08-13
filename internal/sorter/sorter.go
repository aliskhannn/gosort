package sorter

import (
	"sort"
	"strings"
)

// Config содержит настройки сортировки.
// Column: номер колонки (1..N). 0 — сортировать по всей строке.
// Delimiter: разделитель колонок (по умолчанию табуляция).
// Numeric: сравнивать как числа; HumanNumeric: числа с суффиксами (K/M/...).
// Month: сравнивать как месяцы (Jan..Dec). Reverse: обратный порядок.
// Unique: убрать дубликаты после сортировки. IgnoreTrailWS: игнорировать хвостовые пробелы.
// Приоритет ключей: Month > HumanNumeric > Numeric > строковый.
// Если указаны несовместимые флаги, приоритет реализуется в KeyFrom.
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

// Sort сортирует копию входных строк согласно конфигурации.
func Sort(lines []string, cfg Config) ([]string, error) {
	if len(lines) == 0 {
		return []string{}, nil // Пустой вход
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
			// Для стабильности — вторичный ключ по исходному индексу
			return li < lj
		}

		return cmp < 0
	})

	out := make([]string, 0, len(lines))
	var last *key

	for _, id := range idx {
		k := ks[id]
		line := lines[id]

		// Обрезаем хвостовые пробелы для -b
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

// IsSorted проверяет, отсортирован ли срез согласно cfg.
// Возвращает ok, индекс(1-based) первой ошибки, левое и правое значения (для сообщения).
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
