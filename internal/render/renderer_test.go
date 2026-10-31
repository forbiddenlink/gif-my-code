package render

import (
	"crypto/sha256"
	"encoding/hex"
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

// Hash captured from origin/main (8144944) before --padding existed.
const defaultFrameHash = "ede8347dbeb587e954c0596afb1e35ed18c99052e8c0751c43df1a4edd176e61"

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
