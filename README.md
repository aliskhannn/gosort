# gosort

`gosort` — a simplified analogue of the UNIX utility `sort`.
The program sorts lines from files or standard input and outputs them sorted to standard output.

---

## Features

Supported flags, similar to GNU `sort`:

| Flag   | Description                                                        |
| ------ | ------------------------------------------------------------------ |
| `-k N` | Sort by column N (1..N), delimiter — tab by default                |
| `-n`   | Numeric sort                                                       |
| `-r`   | Reverse order                                                      |
| `-u`   | Output only unique lines                                           |
| `-M`   | Sort by month name (Jan, Feb … Dec)                                |
| `-b`   | Ignore trailing whitespace when comparing                          |
| `-c`   | Check whether the input is sorted; if not — print an error message |
| `-h`   | Numeric sort with human-readable suffixes (K, M, G, T, P, E)       |

Flags can be combined, for example: `-nr` — numeric sort in reverse order.

---

## Project structure

```
.
├── cmd                     
│   └── gosort              
│       └── main.go         # Main program file (flag parsing, sorting entry point)
├── internal                
│   └── sorter              # Sorting logic and helper functions
│       ├── sorter.go       # Core sorting implementation
│       ├── key.go          # Column sorting logic (-k)
│       ├── month.go        # Support for month sorting (-M)
│       └── hsize.go        # Support for human-readable size sorting (-h)
├── testdata                # Test data (files for integration tests)
├── go.mod                  
├── go.sum                  
├── Makefile                
└── README.md               
```

---

## Installation

```bash
git clone https://github.com/aliskhannn/gosort.git
cd gosort
```

## Build

```bash
make build
```

The executable will appear in the `bin/` directory:

```text
bin/gosort
```

## Run

```bash
./bin/gosort [OPTIONS] [FILE ...]
```

If no file is specified, STDIN is used.

---

## Usage examples

Sort plain strings:

```bash
echo -e "orange\napple\nbanana" | ./bin/gosort
# Output:
# apple
# banana
# orange
```

Sort by the second column (tab-delimited):

```bash
echo -e "2\tapple\n1\tbanana\n3\tcherry" | ./bin/gosort -k 2
# Output:
# 2	apple
# 1	banana
# 3	cherry
```

Numeric sort:

```bash
echo -e "10\n2\n30" | ./bin/gosort -n
# Output:
# 2
# 10
# 30
```

Remove duplicates while ignoring trailing spaces:

```bash
echo -e "apple  \napple\nbanana" | ./bin/gosort -b -u
# Output:
# apple
# banana
```

Sort months:

```bash
echo -e "Mar\nJan\nFeb" | ./bin/gosort -M
# Output:
# Jan
# Feb
# Mar
```

Check sorting (`-c`):

```bash
echo -e "a\nc\nb" | ./bin/gosort -c
# Output to stderr:
# sort: disorder at line 3: "c" > "b"
```

---

## Testing & Linting

Run tests:

```bash
make test
```

Run linting:

```bash
make lint
```

---

## Cleaning build artifacts

```bash
make clean
```

Removes binaries and temporary directories.

---

## Notes

* `gosort` handles large files efficiently, without fully loading them into memory for key sorting.
* If conflicting flags are specified (`-n` and `-h`, or `-M` with `-n/-h`), a warning is shown, and one mode takes priority.