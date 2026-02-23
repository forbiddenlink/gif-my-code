package render

import (
	"image"
	"image/color"

	"github.com/elizabethstein/gif-my-code/internal/highlight"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gomonobold"
)

// Config holds rendering configuration
type Config struct {
	Width           int
	Height          int
	FontSize        float64
	Padding         int
	LineHeight      float64
	BgColor         color.Color
	CursorColor     color.Color
	HighlightColor  color.Color
	HighlightLines  map[int]bool
}

// Renderer handles image rendering
type Renderer struct {
	config Config
	font   *truetype.Font
}

// NewRenderer creates a new renderer with default config
func NewRenderer(width int, fontSize float64, highlightLines []int) (*Renderer, error) {
	// Parse the embedded font
	font, err := truetype.Parse(gomonobold.TTF)
	if err != nil {
		return nil, err
	}

	// Convert highlight lines to map for O(1) lookup
	highlightMap := make(map[int]bool)
	for _, line := range highlightLines {
		highlightMap[line] = true
	}

	config := Config{
		Width:          width,
		Height:         600, // Will be calculated based on content
		FontSize:       fontSize,
		Padding:        40,
		LineHeight:     1.4,
		BgColor:        color.RGBA{40, 42, 54, 255}, // Dracula background
		CursorColor:    color.RGBA{255, 255, 255, 200}, // White with slight transparency
		HighlightColor: color.RGBA{255, 255, 0, 60}, // Yellow with transparency
		HighlightLines: highlightMap,
	}

	return &Renderer{
		config: config,
		font:   font,
	}, nil
}

// RenderFrame renders a single frame with the given tokens and cursor position
func (r *Renderer) RenderFrame(tokens []highlight.Token, cursorPos int, showCursor bool) (*image.RGBA, error) {
	// Create context
	dc := gg.NewContext(r.config.Width, r.config.Height)

	// Set background
	dc.SetColor(r.config.BgColor)
	dc.Clear()

	// Load font face
	face := truetype.NewFace(r.font, &truetype.Options{
		Size: r.config.FontSize,
	})
	dc.SetFontFace(face)

	// First pass: draw line highlights
	if len(r.config.HighlightLines) > 0 {
		r.drawLineHighlights(dc, tokens, cursorPos)
	}

	// Track position
	x := float64(r.config.Padding)
	y := float64(r.config.Padding) + r.config.FontSize

	charCount := 0
	
	// Draw tokens
	for _, token := range tokens {
		for _, ch := range token.Text {
			if charCount >= cursorPos {
				break
			}

			// Handle newlines
			if ch == '\n' {
				x = float64(r.config.Padding)
				y += r.config.FontSize * r.config.LineHeight
				charCount++
				continue
			}

			// Set color from token style
			textColor := tokenColor(token)
			dc.SetColor(textColor)

			// Draw character
			dc.DrawString(string(ch), x, y)
			
			// Move x position
			w, _ := dc.MeasureString(string(ch))
			x += w

			charCount++
		}

		if charCount >= cursorPos {
			break
		}
	}

	// Draw cursor
	if showCursor && cursorPos <= totalChars(tokens) {
		dc.SetColor(r.config.CursorColor)
		cursorWidth := 10.0
		cursorHeight := r.config.FontSize
		dc.DrawRectangle(x, y-r.config.FontSize+5, cursorWidth, cursorHeight)
		dc.Fill()
	}

	return dc.Image().(*image.RGBA), nil
}

// drawLineHighlights draws highlight backgrounds for specified lines
func (r *Renderer) drawLineHighlights(dc *gg.Context, tokens []highlight.Token, cursorPos int) {
	currentLine := 1
	charCount := 0
	y := float64(r.config.Padding) + r.config.FontSize

	// Calculate which lines are visible based on cursor position
	for _, token := range tokens {
		for _, ch := range token.Text {
			if charCount >= cursorPos {
				return
			}

			if ch == '\n' {
				currentLine++
				y += r.config.FontSize * r.config.LineHeight
			}
			charCount++
		}
	}

	// Reset and draw highlights
	currentLine = 1
	charCount = 0
	y = float64(r.config.Padding) + r.config.FontSize - r.config.FontSize + 5

	for _, token := range tokens {
		for _, ch := range token.Text {
			if charCount >= cursorPos {
				return
			}

			if ch == '\n' {
				// Check if we should highlight this line
				if r.config.HighlightLines[currentLine] {
					dc.SetColor(r.config.HighlightColor)
					highlightHeight := r.config.FontSize * r.config.LineHeight
					dc.DrawRectangle(
						0,
						y,
						float64(r.config.Width),
						highlightHeight,
					)
					dc.Fill()
				}
				
				currentLine++
				y += r.config.FontSize * r.config.LineHeight
			}
			charCount++
		}
	}

	// Highlight the current line if needed
	if r.config.HighlightLines[currentLine] {
		dc.SetColor(r.config.HighlightColor)
		highlightHeight := r.config.FontSize * r.config.LineHeight
		dc.DrawRectangle(
			0,
			y,
			float64(r.config.Width),
			highlightHeight,
		)
		dc.Fill()
	}
}

// tokenColor converts a chroma style to a color
func tokenColor(token highlight.Token) color.Color {
	style := token.Style
	
	// Use foreground color if available
	if style.Colour.IsSet() {
		c := style.Colour
		return color.RGBA{
			R: uint8((c >> 16) & 0xFF),
			G: uint8((c >> 8) & 0xFF),
			B: uint8(c & 0xFF),
			A: 255,
		}
	}
	
	// Default to white
	return color.RGBA{248, 248, 242, 255} // Dracula foreground
}

// totalChars counts total characters in tokens
func totalChars(tokens []highlight.Token) int {
	count := 0
	for _, token := range tokens {
		count += len([]rune(token.Text))
	}
	return count
}

// CalculateHeight calculates the required height based on content
func (r *Renderer) CalculateHeight(tokens []highlight.Token) int {
	lines := 1
	for _, token := range tokens {
		for _, ch := range token.Text {
			if ch == '\n' {
				lines++
			}
		}
	}
	
	height := int(float64(r.config.Padding)*2 + float64(lines)*r.config.FontSize*r.config.LineHeight)
	return height
}
