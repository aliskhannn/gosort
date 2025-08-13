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
	// Определение флагов командной строки
	var (
		flagCol    = flag.IntP("k", "k", 0, "Номер колонки (1..N) для сортировки; по умолчанию вся строка. Разделитель — табуляция.")
		flagNum    = flag.BoolP("numeric", "n", false, "Числовая сортировка (по числовому значению).")
		flagRev    = flag.BoolP("reverse", "r", false, "Обратный порядок (reverse).")
		flagUniq   = flag.BoolP("unique", "u", false, "Выводить только уникальные строки (после сортировки)")
		flagMonth  = flag.BoolP("month-sort", "M", false, "Сортировка по месяцу (Jan..Dec)")
		flagTrimTB = flag.BoolP("ignore-tb", "b", false, "Игнорировать хвостовые пробелы при сравнении")
		flagCheck  = flag.BoolP("check", "c", false, "Проверить, отсортированы ли данные (ничего не сортировать)")
		flagHuman  = flag.BoolP("human-numeric", "h", false, "Числовая сортировка с суффиксами (K, M, G, T, P, E)")
	)

	// Кастомное отображение справки по флагам
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "gosort — упрощённый sort\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [FILE ...]\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nЕсли FILE не указан, читается STDIN.")
	}

	flag.Parse()

	// Валидация несовместимых флагов
	if *flagNum && *flagHuman {
		fmt.Fprintln(os.Stderr, "Предупреждение: указаны -n и -h одновременно; используется -h (человекочитаемые числа).")
	}
	if *flagMonth && (*flagNum || *flagHuman) {
		fmt.Fprintln(os.Stderr, "Предупреждение: -M несовместим с -n/-h; используется -M.")
	}

	// Чтение входных данных из файлов или STDIN
	inputs := flag.Args()
	lines, err := readAll(inputs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения:", err)
		os.Exit(1)
	}

	// Конфигурация сортировки
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

	// Режим проверки сортировки (-c)
	if *flagCheck {
		var ok bool
		var i int
		var a, b string

		if ok, i, a, b = sorter.IsSorted(lines, cfg); ok {
			// Совместимо с GNU sort: ничего не выводим, код 0
			return
		}
		fmt.Fprintf(os.Stderr, "sort: disorder at line %d: %q > %q\n", i, a, b)
		os.Exit(1)
	}

	// Выполнение сортировки
	out, err := sorter.Sort(lines, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка сортировки:", err)
		os.Exit(1)
	}

	// Вывод результата на STDOUT с переносами строк
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

// readAll читает строки из списка файлов.
// Если files пуст, читает из STDIN.
// Поддерживается специальное имя "-" для STDIN.
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

// readFrom читает все строки из io.Reader и возвращает срез строк.
// Игнорируются символы возврата каретки (\r).
// Если вход пуст, возвращается ошибка "no input".
func readFrom(r io.Reader) ([]string, error) {
	// Scanner с увеличенным буфером, чтобы работать с длинными строками.
	s := bufio.NewScanner(r)

	const maxCap = 10 * 1024 * 1024 // 10MB на строку
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
