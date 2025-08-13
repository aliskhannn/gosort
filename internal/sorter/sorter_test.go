package sorter

import (
	"math"
	"reflect"
	"testing"
)

func TestParseHumanSize(t *testing.T) {
	tests := []struct {
		in   string
		want float64
		ok   bool
	}{
		{"10", 10, true},
		{"1K", 1024, true},
		{"1KB", 1024, true},
		{"2M", 2 * 1024 * 1024, true},
		{"1.5G", 1.5 * 1024 * 1024 * 1024, true},
		{"   3T  ", 3 * math.Pow(1024, 4), true},
		{"5P", 5 * math.Pow(1024, 5), true},
		{"6E", 6 * math.Pow(1024, 6), true},
		{"bad", math.NaN(), false},
		{"", math.NaN(), false},
		{"123XB", math.NaN(), false},
	}
	for _, tt := range tests {
		got, ok := parseHumanSize(tt.in)
		if ok != tt.ok {
			t.Errorf("%q: ok=%v, want %v", tt.in, ok, tt.ok)
		}
		if tt.ok && got != tt.want {
			t.Errorf("%q: got=%v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestParseMonth(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"Jan", 1}, {"feb", 2}, {"Mar", 3},
		{"JUN", 6}, {"Dec", 12}, {"Xxx", 0}, {"", 0},
	}
	for _, tt := range tests {
		if got := parseMonth(tt.in); got != tt.want {
			t.Errorf("parseMonth(%q)=%v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestParseNumber(t *testing.T) {
	tests := []struct {
		in   string
		want float64
		ok   bool
	}{
		{"123", 123, true},
		{" 3.14 ", 3.14, true},
		{"-2.5e2", -250, true},
		{"", math.NaN(), false},
		{"bad", 0, false},
	}
	for _, tt := range tests {
		got, ok := parseNumber(tt.in)
		if ok != tt.ok {
			t.Errorf("%q: ok=%v, want %v", tt.in, ok, tt.ok)
		}
		if tt.ok && got != tt.want {
			t.Errorf("%q: got=%v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestSort_String(t *testing.T) {
	lines := []string{"b", "a", "c"}
	got, err := Sort(lines, Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSort_Numeric(t *testing.T) {
	lines := []string{"10", "2", "1"}
	got, _ := Sort(lines, Config{Numeric: true})
	want := []string{"1", "2", "10"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSort_HumanNumeric(t *testing.T) {
	lines := []string{"2K", "1K", "3K"}
	got, _ := Sort(lines, Config{HumanNumeric: true})
	want := []string{"1K", "2K", "3K"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSort_Month(t *testing.T) {
	lines := []string{"Mar", "Jan", "Feb"}
	got, _ := Sort(lines, Config{Month: true})
	want := []string{"Jan", "Feb", "Mar"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestIsSorted(t *testing.T) {
	lines := []string{"a", "b", "c"}
	if ok, _, _, _ := IsSorted(lines, Config{}); !ok {
		t.Error("expected sorted")
	}
	lines = []string{"b", "a"}
	if ok, idx, left, right := IsSorted(lines, Config{}); ok || idx != 2 || left != "b" || right != "a" {
		t.Errorf("unexpected result: ok=%v, idx=%v, left=%q, right=%q", ok, idx, left, right)
	}
}
