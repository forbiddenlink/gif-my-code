package render

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/forbiddenlink/gif-my-code/internal/highlight"
)

const sample = "package main\n\nfunc main() {\n\tprintln(\"hi\")\n}\n"

func frameHash(t *testing.T, r *Renderer) string {
	t.Helper()
	code, err := highlight.Highlight(sample, "go", "dracula")
	if err != nil {
		t.Fatal(err)
	}
	img, err := r.RenderFrame(code.Tokens, len([]rune(sample)), false, 1.0)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(img.Pix)
	return hex.EncodeToString(sum[:])
}

// Hash of the default frame after tab expansion landed on main; --padding must not change it.
const defaultFrameHash = "51233ad63127818db8e8f7659b75a6f0eb4c193f92e2602da94a5b0eedd1d189"

func TestDefaultPaddingMatchesPreviousOutput(t *testing.T) {
	r, err := NewRenderer(800, 16, nil, "macos", "dracula", false, true, "go", true, DefaultPadding)
	if err != nil {
		t.Fatal(err)
	}
	if got := frameHash(t, r); got != defaultFrameHash {
		t.Errorf("default frame hash = %s, want %s", got, defaultFrameHash)
	}
}

func TestPaddingChangesOutput(t *testing.T) {
	base, _ := NewRenderer(800, 16, nil, "macos", "dracula", false, true, "go", true, DefaultPadding)
	zero, err := NewRenderer(800, 16, nil, "macos", "dracula", false, true, "go", true, 0)
	if err != nil {
		t.Fatal(err)
	}
	if frameHash(t, base) == frameHash(t, zero) {
		t.Error("padding 0 should differ from default padding")
	}
}

func TestPaddingScalesWithHiDPI(t *testing.T) {
	r, err := NewRenderer(800, 16, nil, "none", "dracula", true, false, "go", true, 10)
	if err != nil {
		t.Fatal(err)
	}
	if r.config.Padding != 20 {
		t.Errorf("hidpi padding = %d, want 20", r.config.Padding)
	}
}

func TestPaddingRejectsNegative(t *testing.T) {
	_, err := NewRenderer(800, 16, nil, "none", "dracula", false, false, "go", true, -1)
	if err == nil {
		t.Fatal("expected error for negative padding")
	}
}

func TestPaddingRejectsTooLarge(t *testing.T) {
	for _, p := range []int{300, 400, 1000000} {
		if _, err := NewRenderer(800, 16, nil, "none", "dracula", false, false, "go", true, p); err == nil {
			t.Errorf("padding %d should be rejected", p)
		}
	}
	if _, err := NewRenderer(800, 16, nil, "none", "dracula", false, false, "go", true, 299); err != nil {
		t.Errorf("padding 299 should be accepted: %v", err)
	}
}

func tokensOf(t *testing.T, src string) []highlight.Token {
	t.Helper()
	code, err := highlight.Highlight(src, "go", "dracula")
	if err != nil {
		t.Fatal(err)
	}
	return code.Tokens
}

func TestOverflowLines(t *testing.T) {
	long := strings.Repeat("x := 1\n", 40)
	cases := []struct {
		name    string
		src     string
		window  string
		padding int
		want    int
	}{
		{"short file at default fits", sample, "macos", DefaultPadding, 0},
		{"long file at default overflows", long, "macos", DefaultPadding, 20},
		{"short file at large padding overflows", strings.Repeat("x := 1\n", 14), "macos", 200, 7},
		{"no window chrome leaves more room", strings.Repeat("x := 1\n", 14), "none", 200, 6},
		{"padding 0 fits the same file", strings.Repeat("x := 1\n", 14), "macos", 0, 0},
	}
	for _, c := range cases {
		got := OverflowLines(tokensOf(t, c.src), 16, c.window, c.padding)
		if got != c.want {
			t.Errorf("%s: OverflowLines = %d, want %d", c.name, got, c.want)
		}
	}
}
