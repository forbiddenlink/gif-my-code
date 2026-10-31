package highlight

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// TabWidth is the tab stop width in columns. The embedded font has no tab glyph,
// so tabs are expanded to spaces before rendering.
const TabWidth = 4

// Token represents a syntax-highlighted token
type Token struct {
	Text  string
	Style chroma.StyleEntry
}

// HighlightedCode represents syntax-highlighted code as a list of tokens
type HighlightedCode struct {
	Tokens []Token
	Theme  string
}

// Highlight applies syntax highlighting to code
func Highlight(code, language, themeName string) (*HighlightedCode, error) {
	// Get lexer for language
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	// Get style/theme
	style := styles.Get(themeName)
	if style == nil {
		style = styles.Fallback
	}

	// Tokenize the code
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize: %w", err)
	}

	// Convert to our Token format
	tokens := []Token{}
	for _, token := range iterator.Tokens() {
		styleEntry := style.Get(token.Type)
		tokens = append(tokens, Token{
			Text:  token.Value,
			Style: styleEntry,
		})
	}

	return &HighlightedCode{
		Tokens: expandTabs(tokens),
		Theme:  themeName,
	}, nil
}

// ToPlainText returns the code as plain text (useful for debugging)
func (h *HighlightedCode) ToPlainText() string {
	var buf bytes.Buffer
	for _, token := range h.Tokens {
		buf.WriteString(token.Text)
	}
	return buf.String()
}

// ListThemes returns all available theme names
func ListThemes() []string {
	names := styles.Names()
	return names
}

// ValidateTheme checks if a theme exists
func ValidateTheme(name string) bool {
	return styles.Get(name) != nil
}

// expandTabs replaces each tab with spaces up to the next tab stop, tracking the
// column across tokens. Doing it here keeps the renderer's character counts,
// cursor position and line highlights in step with what is drawn.
func expandTabs(tokens []Token) []Token {
	col := 0
	out := make([]Token, len(tokens))
	for i, tok := range tokens {
		var b strings.Builder
		for _, ch := range tok.Text {
			switch ch {
			case '\t':
				n := TabWidth - col%TabWidth
				b.WriteString(strings.Repeat(" ", n))
				col += n
			case '\n':
				b.WriteRune(ch)
				col = 0
			default:
				b.WriteRune(ch)
				col++
			}
		}
		out[i] = Token{Text: b.String(), Style: tok.Style}
	}
	return out
}
