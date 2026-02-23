package render

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/forbiddenlink/gif-my-code/internal/highlight"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gomonobold"
)

// Config holds rendering configuration
type Config struct {
	Width          int
	Height         int
	FontSize       float64
	Padding        int
	LineHeight     float64
	BgColor        color.Color
	CursorColor    color.Color
	HighlightColor color.Color
	HighlightLines map[int]bool
	WindowStyle    string // "macos", "windows", "none"
	Theme          string
	CornerRadius   float64
	ShadowEnabled  bool
}

// Renderer handles image rendering
type Renderer struct {
	config Config
	font   *truetype.Font
}

// NewRenderer creates a new renderer with enhanced visual config
func NewRenderer(width int, fontSize float64, highlightLines []int, windowStyle string, theme string) (*Renderer, error) {
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

	// Enhanced config with professional styling - Cinematic Deep Space theme
	config := Config{
		Width:          width,
		Height:         600, // Will be calculated based on content
		FontSize:       fontSize,
		Padding:        36,                             // Tighter, more professional padding
		LineHeight:     1.5,                            // Better readability
		BgColor:        color.RGBA{15, 17, 26, 255},    // Deep Space Window Base (#0F111A)
		CursorColor:    color.RGBA{255, 255, 255, 220}, // Slightly more opaque
		HighlightColor: color.RGBA{255, 255, 255, 15},  // Subtle white wash for highlight background
		HighlightLines: highlightMap,
		WindowStyle:    windowStyle,
		Theme:          theme,
		CornerRadius:   16,   // Smoother, larger rounded corners
		ShadowEnabled:  true, // Drop shadow
	}

	return &Renderer{
		config: config,
		font:   font,
	}, nil
}

// RenderFrame renders a single frame with the given tokens and cursor position
func (r *Renderer) RenderFrame(tokens []highlight.Token, cursorPos int, showCursor bool) (*image.RGBA, error) {
	// Create context with extra space for shadow
	shadowOffset := 20.0
	dc := gg.NewContext(r.config.Width+int(shadowOffset*2), r.config.Height+int(shadowOffset*2))

	// Clear with transparent background
	dc.SetColor(color.RGBA{0, 0, 0, 0})
	dc.Clear()

	// Draw gradient background with rounded corners
	r.drawGradientBackground(dc, shadowOffset)

	// Draw window chrome if enabled
	if r.config.WindowStyle == "macos" {
		r.drawMacOSChrome(dc, shadowOffset)
	} else if r.config.WindowStyle == "windows" {
		r.drawWindowsChrome(dc, shadowOffset)
	}

	// Load font face
	face := truetype.NewFace(r.font, &truetype.Options{
		Size: r.config.FontSize,
	})
	dc.SetFontFace(face)

	// First pass: draw line highlights
	if len(r.config.HighlightLines) > 0 {
		r.drawLineHighlights(dc, tokens, cursorPos, shadowOffset)
	}

	// Track position (adjusted for shadow offset)
	chromeHeight := 0.0
	if r.config.WindowStyle == "macos" || r.config.WindowStyle == "windows" {
		chromeHeight = 40.0 // Title bar height
	}

	x := float64(r.config.Padding) + shadowOffset
	y := float64(r.config.Padding) + r.config.FontSize + shadowOffset + chromeHeight

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

// drawGradientBackground draws a gradient background with rounded corners and shadow
func (r *Renderer) drawGradientBackground(dc *gg.Context, offset float64) {
	// Draw multi-layered soft shadow for cinematic depth
	if r.config.ShadowEnabled {
		// Layer 1: Ambient large glow
		dc.SetColor(color.RGBA{0, 0, 0, 15})
		dc.DrawRoundedRectangle(
			offset-4,
			offset+12,
			float64(r.config.Width)+8,
			float64(r.config.Height)+8,
			r.config.CornerRadius+4,
		)
		dc.Fill()

		// Layer 2: Mid shadow
		dc.SetColor(color.RGBA{0, 0, 0, 25})
		dc.DrawRoundedRectangle(
			offset-2,
			offset+8,
			float64(r.config.Width)+4,
			float64(r.config.Height)+4,
			r.config.CornerRadius+2,
		)
		dc.Fill()

		// Layer 3: Tight contact shadow
		dc.SetColor(color.RGBA{0, 0, 0, 50})
		dc.DrawRoundedRectangle(
			offset,
			offset+4,
			float64(r.config.Width),
			float64(r.config.Height),
			r.config.CornerRadius,
		)
		dc.Fill()
	}

	// Draw gradient background
	// Simulate deep space gradient with two-tone approach
	gradient1 := color.RGBA{13, 14, 21, 255} // #0D0E15 Base Canvas
	gradient2 := color.RGBA{22, 24, 33, 255} // #161821 Inner Glow

	// Draw base rectangle
	dc.SetColor(gradient1)
	dc.DrawRoundedRectangle(
		offset,
		offset,
		float64(r.config.Width),
		float64(r.config.Height),
		r.config.CornerRadius,
	)
	dc.Fill()

	// Draw subtle overlay for gradient effect
	dc.SetColor(gradient2)
	dc.DrawRoundedRectangle(
		offset,
		offset+float64(r.config.Height)*0.15,
		float64(r.config.Width),
		float64(r.config.Height)*0.85,
		r.config.CornerRadius,
	)
	dc.Fill()

	// The Ghost Outline (1px inner ring)
	dc.SetColor(color.RGBA{255, 255, 255, 20}) // 8% opacity white
	dc.DrawRoundedRectangle(
		offset+0.5,
		offset+0.5,
		float64(r.config.Width)-1,
		float64(r.config.Height)-1,
		r.config.CornerRadius,
	)
	dc.SetLineWidth(1.0)
	dc.Stroke()
}

// drawMacOSChrome draws macOS-style window controls
func (r *Renderer) drawMacOSChrome(dc *gg.Context, offset float64) {
	y := offset + 20.0
	x := offset + 20.0
	spacing := 8.0
	dotSize := 12.0 // Slightly larger/standard dots

	// Red dot
	dc.SetColor(color.RGBA{236, 106, 94, 255}) // #EC6A5E Premium Red
	dc.DrawCircle(x, y, dotSize/2)
	dc.Fill()

	// Yellow dot
	x += dotSize + spacing
	dc.SetColor(color.RGBA{244, 191, 79, 255}) // #F4BF4F Premium Yellow
	dc.DrawCircle(x, y, dotSize/2)
	dc.Fill()

	// Green dot
	x += dotSize + spacing
	dc.SetColor(color.RGBA{97, 197, 84, 255}) // #61C554 Premium Green
	dc.DrawCircle(x, y, dotSize/2)
	dc.Fill()
}

// drawWindowsChrome draws Windows-style window controls
func (r *Renderer) drawWindowsChrome(dc *gg.Context, offset float64) {
	// Draw title bar
	dc.SetColor(color.RGBA{30, 30, 30, 255})
	dc.DrawRectangle(offset, offset, float64(r.config.Width), 40)
	dc.Fill()

	// Draw window controls (simplified)
	x := offset + float64(r.config.Width) - 40
	y := offset + 20.0
	dc.SetColor(color.RGBA{200, 200, 200, 255})
	dc.DrawRectangle(x-30, y-6, 12, 12)
	dc.StrokePreserve()
}

// drawLineHighlights draws highlight backgrounds for specified lines
func (r *Renderer) drawLineHighlights(dc *gg.Context, tokens []highlight.Token, cursorPos int, offset float64) {
	currentLine := 1
	charCount := 0

	// Pre-calculate line heights for positioning
	chromeHeight := 0.0
	if r.config.WindowStyle == "macos" || r.config.WindowStyle == "windows" {
		chromeHeight = 40.0
	}

	// Start y position aligned with the text baseline, adjusted back up to bounds
	y := offset + float64(r.config.Padding) + chromeHeight - 5
	accentColor := color.RGBA{0, 240, 255, 255} // Neon Cyan
	highlightHeight := r.config.FontSize * r.config.LineHeight

	for _, token := range tokens {
		for _, ch := range token.Text {
			if charCount >= cursorPos {
				return
			}

			if ch == '\n' {
				// Check if we should highlight this line
				if r.config.HighlightLines[currentLine] {
					// 1. Draw subtle background wash
					dc.SetColor(r.config.HighlightColor)
					dc.DrawRectangle(
						offset,
						y,
						float64(r.config.Width),
						highlightHeight,
					)
					dc.Fill()

					// 2. Draw vibrant left anchor border
					dc.SetColor(accentColor)
					dc.DrawRectangle(
						offset,
						y,
						4.0, // 4px wide accent
						highlightHeight,
					)
					dc.Fill()
				}

				currentLine++
				y += highlightHeight
			}
			charCount++
		}
	}

	// Highlight the current line if needed (last line case)
	if r.config.HighlightLines[currentLine] {
		// 1. Draw subtle background wash
		dc.SetColor(r.config.HighlightColor)
		dc.DrawRectangle(
			offset,
			y,
			float64(r.config.Width),
			highlightHeight,
		)
		dc.Fill()

		// 2. Draw vibrant left anchor border
		dc.SetColor(accentColor)
		dc.DrawRectangle(
			offset,
			y,
			4.0, // 4px wide accent
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
