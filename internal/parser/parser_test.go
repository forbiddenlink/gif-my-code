package parser

import (
	"reflect"
	"strconv"
	"testing"
)

func TestValidLines(t *testing.T) {
	g, e := ParseHighlightLines("5,7-9")
	if e != nil || !reflect.DeepEqual(g, []int{5, 7, 8, 9}) {
		t.Fatal(g, e)
	}
}
func TestReversedRange(t *testing.T) {
	if _, e := ParseHighlightLines("9-7"); e == nil {
		t.Fatal("expected error")
	}
}
func TestZeroLine(t *testing.T) {
	if _, e := ParseHighlightLines("0"); e == nil {
		t.Fatal("zero line accepted")
	}
}
func TestZeroRange(t *testing.T) {
	if _, e := ParseHighlightLines("0-2"); e == nil {
		t.Fatal("zero range accepted")
	}
}

func TestHighlightExpansionLimit(t *testing.T) {
	lines, err := ParseHighlightLines("1-100000")
	if err != nil || len(lines) != MaxHighlightEntries {
		t.Fatal(len(lines), err)
	}
	for _, spec := range []string{"1-100001", "1-100000,1", "1,1-100000"} {
		if _, err := ParseHighlightLines(spec); err == nil {
			t.Fatalf("accepted %s", spec)
		}
	}
}
func TestMaxIntEndpointDoesNotOverflow(t *testing.T) {
	max := int(^uint(0) >> 1)
	spec := strconv.Itoa(max-1) + "-" + strconv.Itoa(max)
	lines, err := ParseHighlightLines(spec)
	if err != nil || !reflect.DeepEqual(lines, []int{max - 1, max}) {
		t.Fatal(lines, err)
	}
	if _, err := ParseHighlightLines("1-" + strconv.Itoa(max)); err == nil {
		t.Fatal("accepted unbounded expansion")
	}
}
func TestInvalidHighlightNumbers(t *testing.T) {
	for _, spec := range []string{"1.5", "1-2.5", "99999999999999999999999999999"} {
		if _, err := ParseHighlightLines(spec); err == nil {
			t.Fatalf("accepted %s", spec)
		}
	}
}
