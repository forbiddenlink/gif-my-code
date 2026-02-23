package render

import (
	"fmt"
	"image"
	"image/color"
	"math"

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
	HiDPI          bool
	LineNumbers    bool
	ScaleFactor    float64
}

// Renderer handles image rendering
type Renderer struct {
	config Config
	font   *truetype.Font
}

// NewRenderer creates a new renderer with enhanced visual config
func NewRenderer(width int, fontSize float64, highlightLines []int, windowStyle string, theme string, hiDPI bool, lineNumbers bool) (*Renderer, error) {
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

	scaleFactor := 1.0
	if hiDPI {
		scaleFactor = 2.0
	}

	// Enhanced config with professional styling - Cinematic Deep Space theme
	config := Config{
		Width:          int(float64(width) * scaleFactor),
		Height:         int(600 * scaleFactor), // Will be calculated based on content
		FontSize:       fontSize * scaleFactor,
		Padding:        int(36.0 * scaleFactor),        // Tighter, more professional padding
		LineHeight:     1.5,                            // Better readability
		BgColor:        color.RGBA{15, 17, 26, 255},    // Deep Space Window Base (#0F111A)
		CursorColor:    color.RGBA{255, 255, 255, 220}, // Slightly more opaque
		HighlightColor: color.RGBA{255, 255, 255, 15},  // Subtle white wash for highlight background
		HighlightLines: highlightMap,
		WindowStyle:    windowStyle,
		Theme:          theme,
		CornerRadius:   16.0 * scaleFactor, // Smoother, larger rounded corners
		ShadowEnabled:  true,               // Drop shadow
		HiDPI:          hiDPI,
		LineNumbers:    lineNumbers,
		ScaleFactor:    scaleFactor,
	}

	return &Renderer{
		config: config,
		font:   font,
	}, nil
}

// RenderFrame renders a single frame with the given tokens and cursor position
func (r *Renderer) RenderFrame(tokens []highlight.Token, cursorPos int, showCursor bool, progress float64) (*image.RGBA, error) {
	// Create context with extra space for shadow
	shadowOffset := 20.0 * r.config.ScaleFactor
	dc := gg.NewContext(r.config.Width+int(shadowOffset*2), r.config.Height+int(shadowOffset*2))

	// Clear with transparent background
	dc.SetColor(color.RGBA{0, 0, 0, 0})
	dc.Clear()

	// Draw gradient background with rounded corners
	r.drawGradientBackground(dc, shadowOffset, progress)

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

	gutterWidth := 0.0
	if r.config.LineNumbers {
		gutterWidth = 50.0 * r.config.ScaleFactor
	}

	// First pass: draw line highlights and line numbers
	if len(r.config.HighlightLines) > 0 || r.config.LineNumbers {
		r.drawLineHighlights(dc, tokens, cursorPos, shadowOffset, gutterWidth)
	}

	// Track position (adjusted for shadow offset)
	chromeHeight := 0.0
	if r.config.WindowStyle == "macos" || r.config.WindowStyle == "windows" {
		chromeHeight = 40.0 * r.config.ScaleFactor // Title bar height
	}

	x := float64(r.config.Padding) + shadowOffset + gutterWidth
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
				x = float64(r.config.Padding) + shadowOffset + gutterWidth
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
		cursorWidth := 10.0 * r.config.ScaleFactor
		cursorHeight := r.config.FontSize
		dc.DrawRectangle(x, y-r.config.FontSize+5*r.config.ScaleFactor, cursorWidth, cursorHeight)
		dc.Fill()
	}

	return dc.Image().(*image.RGBA), nil
}

// drawGradientBackground draws a gradient background with rounded corners and shadow
func (r *Renderer) drawGradientBackground(dc *gg.Context, offset float64, progress float64) {
	// Draw multi-layered soft shadow for cinematic depth
	if r.config.ShadowEnabled {
		// Glow tint based on theme (simplified syntax-aware glow)
		glowColor := color.RGBA{0, 240, 255, 10} // default neon cyan glow
		if r.config.Theme == "dracula" {
			glowColor = color.RGBA{255, 121, 198, 10} // pink glow
		} else if r.config.Theme == "monokai" {
			glowColor = color.RGBA{253, 151, 31, 10} // orange glow
		} else if r.config.Theme == "nord" {
			glowColor = color.RGBA{136, 192, 208, 10} // ice blue glow
		}

		// Layer 1: Ambient large glow (colored)
		dc.SetColor(glowColor)
		dc.DrawRoundedRectangle(
			offset-4*r.config.ScaleFactor,
			offset+12*r.config.ScaleFactor,
			float64(r.config.Width)+8*r.config.ScaleFactor,
			float64(r.config.Height)+8*r.config.ScaleFactor,
			r.config.CornerRadius+4*r.config.ScaleFactor,
		)
		dc.Fill()

		// Layer 2: Mid shadow
		dc.SetColor(color.RGBA{0, 0, 0, 25})
		dc.DrawRoundedRectangle(
			offset-2*r.config.ScaleFactor,
			offset+8*r.config.ScaleFactor,
			float64(r.config.Width)+4*r.config.ScaleFactor,
			float64(r.config.Height)+4*r.config.ScaleFactor,
			r.config.CornerRadius+2*r.config.ScaleFactor,
		)
		dc.Fill()

		// Layer 3: Tight contact shadow
		dc.SetColor(color.RGBA{0, 0, 0, 50})
		dc.DrawRoundedRectangle(
			offset,
			offset+4*r.config.ScaleFactor,
			float64(r.config.Width),
			float64(r.config.Height),
			r.config.CornerRadius,
		)
		dc.Fill()
	}

	// Draw gradient background
	// Simulate deep space gradient with animated breathing
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

	// Breathing effect: vary the height and position of the inner glow
	breathOffset := math.Sin(progress*math.Pi*2) * 10 * r.config.ScaleFactor

	// Draw subtle overlay for gradient effect
	dc.SetColor(gradient2)
	dc.DrawRoundedRectangle(
		offset,
		offset+float64(r.config.Height)*0.15+breathOffset,
		float64(r.config.Width),
		float64(r.config.Height)*0.85-breathOffset,
		r.config.CornerRadius,
	)
	dc.Fill()

	// The Ghost Outline (1px inner ring)
	dc.SetColor(color.RGBA{255, 255, 255, 20}) // 8% opacity white
	dc.DrawRoundedRectangle(
		offset+0.5,
		offset+0.5,
		float64(r.config.Width)-1*r.config.ScaleFactor,
		float64(r.config.Height)-1*r.config.ScaleFactor,
		r.config.CornerRadius,
	)
	dc.SetLineWidth(1.0 * r.config.ScaleFactor)
	dc.Stroke()
}

// drawMacOSChrome draws macOS-style window controls
func (r *Renderer) drawMacOSChrome(dc *gg.Context, offset float64) {
	y := offset + 20.0*r.config.ScaleFactor
	x := offset + 20.0*r.config.ScaleFactor
	spacing := 8.0 * r.config.ScaleFactor
	dotSize := 12.0 * r.config.ScaleFactor

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
	dc.DrawRectangle(offset, offset, float64(r.config.Width), 40*r.config.ScaleFactor)
	dc.Fill()

	// Draw window controls (simplified)
	x := offset + float64(r.config.Width) - 40*r.config.ScaleFactor
	y := offset + 20.0*r.config.ScaleFactor
	dc.SetColor(color.RGBA{200, 200, 200, 255})
	dc.DrawRectangle(x-30*r.config.ScaleFactor, y-6*r.config.ScaleFactor, 12*r.config.ScaleFactor, 12*r.config.ScaleFactor)
	dc.StrokePreserve()
}

// drawLineHighlights draws highlight backgrounds and line numbers for specified lines
func (r *Renderer) drawLineHighlights(dc *gg.Context, tokens []highlight.Token, cursorPos int, offset float64, gutterWidth float64) {
	currentLine := 1
	charCount := 0

	// Pre-calculate line heights for positioning
	chromeHeight := 0.0
	if r.config.WindowStyle == "macos" || r.config.WindowStyle == "windows" {
		chromeHeight = 40.0 * r.config.ScaleFactor
	}

	// Start y position aligned with the text baseline, adjusted back up to bounds
	y := offset + float64(r.config.Padding) + chromeHeight - 5*r.config.ScaleFactor
	accentColor := color.RGBA{0, 240, 255, 255} // Neon Cyan
	highlightHeight := r.config.FontSize * r.config.LineHeight

	drawNumbersAndHighlights := func(line int, currentY float64) {
		// Draw highlight if enabled
		if r.config.HighlightLines[line] {
			// 1. Draw subtle background wash
			dc.SetColor(r.config.HighlightColor)
			dc.DrawRectangle(
				offset,
				currentY,
				float64(r.config.Width),
				highlightHeight,
			)
			dc.Fill()

			// 2. Draw vibrant left anchor border
			dc.SetColor(accentColor)
			dc.DrawRectangle(
				offset,
				currentY,
				4.0*r.config.ScaleFactor, // 4px wide accent
				highlightHeight,
			)
			dc.Fill()
		}

		// Draw line numbers if enabled
		if r.config.LineNumbers {
			numStr := fmt.Sprintf("%2d", line)

			// Load slightly smaller font for line numbers
			numberSize := r.config.FontSize * 0.8
			face := truetype.NewFace(r.font, &truetype.Options{
				Size: numberSize,
			})
			dc.SetFontFace(face)

			// Faint white for line numbers
			dc.SetColor(color.RGBA{255, 255, 255, 100})

			// Position number in the gutter
			numX := offset + float64(r.config.Padding)
			numY := currentY + r.config.FontSize
			dc.DrawString(numStr, numX, numY)

			// Draw 1px vertical separator line at the right edge of gutter
			if line == 1 {
				dc.SetColor(color.RGBA{255, 255, 255, 15})
				sepX := offset + float64(r.config.Padding) + gutterWidth - (15.0 * r.config.ScaleFactor)
				dc.DrawLine(sepX, offset+chromeHeight+float64(r.config.Padding), sepX, float64(r.config.Height)-float64(r.config.Padding))
				dc.SetLineWidth(1.0 * r.config.ScaleFactor)
				dc.Stroke()
			}

			// Restore normal font size
			normalFace := truetype.NewFace(r.font, &truetype.Options{
				Size: r.config.FontSize,
			})
			dc.SetFontFace(normalFace)
		}
	}

	// Draw the first line immediately
	drawNumbersAndHighlights(currentLine, y)

	for _, token := range tokens {
		for _, ch := range token.Text {
			if charCount >= cursorPos {
				return
			}

			if ch == '\n' {
				currentLine++
				y += highlightHeight
				drawNumbersAndHighlights(currentLine, y)
			}
			charCount++
		}
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

	chromeHeight := 0.0
	if r.config.WindowStyle == "macos" || r.config.WindowStyle == "windows" {
		chromeHeight = 40.0 * r.config.ScaleFactor
	}

	height := int(float64(r.config.Padding)*2 + float64(lines)*r.config.FontSize*r.config.LineHeight + chromeHeight)
	return height
}
