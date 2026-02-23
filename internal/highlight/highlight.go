package highlight

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

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
		Tokens: tokens,
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
