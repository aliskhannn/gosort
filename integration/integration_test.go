//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const binName = "gosort"

func buildBinary(t *testing.T) string {
	t.Helper()
	binPath := filepath.Join(t.TempDir(), binName)
	cmd := exec.Command("go", "build", "-o", binPath, "github.com/aliskhannn/gosort/cmd/gosort")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, string(out))
	}
	return binPath
}

func runCmd(t *testing.T, bin string, args ...string) string {
	t.Helper()
	cmd := exec.Command(bin, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out.String())
	}
	return strings.TrimRight(out.String(), "\n")
}

func TestSort_Small(t *testing.T) {
	bin := buildBinary(t)
	out := runCmd(t, bin, "../testdata/small.txt")
	expected := strings.TrimRight(readFile(t, "../testdata/expected_small.txt"), "\n")
	if out != expected {
		t.Errorf("unexpected output:\nGot:\n%s\nWant:\n%s", out, expected)
	}
}

func TestSort_ByColumn(t *testing.T) {
	bin := buildBinary(t)
	out := runCmd(t, bin, "-k", "2", "../testdata/table.txt")
	expected := "2\ta\n1\tb\n3\tc"
	if out != expected {
		t.Errorf("unexpected output:\nGot:\n%s\nWant:\n%s", out, expected)
	}
}

func TestSort_Numeric(t *testing.T) {
	bin := buildBinary(t)
	out := runCmd(t, bin, "-n", "../testdata/nums.txt")
	expected := "1\n2\n10\n20"
	if out != expected {
		t.Errorf("unexpected output:\nGot:\n%s\nWant:\n%s", out, expected)
	}
}

func TestSort_Months(t *testing.T) {
	bin := buildBinary(t)
	out := runCmd(t, bin, "-M", "../testdata/months.txt")
	expected := "Jan\nFeb\nMar\nDec"
	if out != expected {
		t.Errorf("unexpected output:\nGot:\n%s\nWant:\n%s", out, expected)
	}
}

func TestSort_HumanSizes(t *testing.T) {
	bin := buildBinary(t)
	out := runCmd(t, bin, "-h", "../testdata/humansizes.txt")
	expected := "512\n1K\n128K\n2M"
	if out != expected {
		t.Errorf("unexpected output:\nGot:\n%s\nWant:\n%s", out, expected)
	}
}

func TestSort_Blanks_Unique(t *testing.T) {
	bin := buildBinary(t)
	out := runCmd(t, bin, "-u", "-b", "../testdata/blanks.txt")
	expected := "apple\nbanana"
	if out != expected {
		t.Errorf("unexpected output:\nGot:\n%s\nWant:\n%s", out, expected)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	return string(data)
}
