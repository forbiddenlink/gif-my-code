package render

import (
	"bytes"
	"testing"

	"github.com/forbiddenlink/gif-my-code/internal/highlight"
)

func tabFrame(t *testing.T, src string) []byte {
	t.Helper()
	code, err := highlight.Highlight(src, "go", "dracula")
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, tok := range code.Tokens {
		n += len([]rune(tok.Text))
	}
	r, err := NewRenderer(800, 16, []int{2}, "macos", "dracula", false, true, "go", true)
	if err != nil {
		t.Fatal(err)
	}
	img, err := r.RenderFrame(code.Tokens, n, false, 1.0)
	if err != nil {
		t.Fatal(err)
	}
	return img.Pix
}

// A tab-indented file must render exactly like the same file indented with TabWidth spaces,
// including highlight spans and line numbers, instead of drawing a missing-glyph box.
func TestTabIndentRendersLikeSpaces(t *testing.T) {
	withTabs := tabFrame(t, "func f() {\n\treturn\n}\n")
	withSpaces := tabFrame(t, "func f() {\n    return\n}\n")
	if !bytes.Equal(withTabs, withSpaces) {
		t.Error("tab-indented frame differs from space-indented frame")
	}
}
