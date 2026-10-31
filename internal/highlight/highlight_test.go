package highlight

import (
	"strings"
	"testing"
)

func joined(h *HighlightedCode) string {
	var b strings.Builder
	for _, t := range h.Tokens {
		b.WriteString(t.Text)
	}
	return b.String()
}

func TestTabsExpandToNextTabStop(t *testing.T) {
	pad := strings.Repeat(" ", TabWidth)
	cases := []struct{ name, src, want string }{
		{"leading tab", "\tx\n", pad + "x\n"},
		{"two tabs", "\t\tx\n", pad + pad + "x\n"},
		{"mid-line tab aligns to stop", "a\tb\n", "a" + strings.Repeat(" ", TabWidth-1) + "b\n"},
		{"column resets after newline", "ab\n\tc\n", "ab\n" + pad + "c\n"},
		{"no tabs unchanged", "x := 1\n", "x := 1\n"},
	}
	for _, c := range cases {
		h, err := Highlight(c.src, "go", "dracula")
		if err != nil {
			t.Fatal(err)
		}
		if got := joined(h); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}

func TestTabsExpandedAcrossTokens(t *testing.T) {
	// A tab after a multi-token prefix still aligns to the next stop.
	h, err := Highlight("if x\t{\n}\n", "go", "dracula")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(joined(h), "\t") {
		t.Errorf("tab survived: %q", joined(h))
	}
	if !strings.HasPrefix(joined(h), "if x"+strings.Repeat(" ", TabWidth-len("if x")%TabWidth)+"{") {
		t.Errorf("wrong alignment: %q", joined(h))
	}
}
