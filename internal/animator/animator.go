package animator

import (
	"fmt"
	"image"
	"math"

	"github.com/forbiddenlink/gif-my-code/internal/highlight"
	"github.com/forbiddenlink/gif-my-code/internal/render"
)

// Config holds animation configuration
type Config struct {
	Width          int
	FontSize       float64
	Speed          float64
	FPS            int
	ShowCursor     bool
	HighlightLines []int
	WindowStyle    string
	Theme          string
	HiDPI          bool
	LineNumbers    bool
	Language       string
	LaserReveal    bool
}

// GenerateFrames creates all animation frames
func GenerateFrames(code *highlight.HighlightedCode, config Config) ([]*image.RGBA, error) {
	// Create renderer with highlight config and visual enhancements
	renderer, err := render.NewRenderer(config.Width, config.FontSize, config.HighlightLines, config.WindowStyle, config.Theme, config.HiDPI, config.LineNumbers, config.Language, config.LaserReveal)
	if err != nil {
		return nil, fmt.Errorf("failed to create renderer: %w", err)
	}

	// Calculate total characters
	totalChars := 0
	for _, token := range code.Tokens {
		totalChars += len([]rune(token.Text))
	}

	// Calculate characters per frame based on speed
	charsPerFrame := int(math.Max(1, 2*config.Speed))

	// Calculate cursor blink interval (blink every 15 frames = 0.5 seconds at 30fps)
	cursorBlinkInterval := config.FPS / 2

	// Calculate total frames to estimate animation progress
	typingFrames := (totalChars / charsPerFrame) + 1
	finalFrameCount := config.FPS * 2
	totalFrames := typingFrames + finalFrameCount

	frames := []*image.RGBA{}
	cursorVisible := true
	frameCount := 0

	// Generate typing frames
	for charPos := 0; charPos <= totalChars; charPos += charsPerFrame {
		if charPos > totalChars {
			charPos = totalChars
		}

		// Toggle cursor visibility
		if frameCount%cursorBlinkInterval == 0 {
			cursorVisible = !cursorVisible
		}

		showCursor := config.ShowCursor && cursorVisible
		progress := float64(frameCount) / float64(totalFrames)

		frame, err := renderer.RenderFrame(code.Tokens, charPos, showCursor, progress)
		if err != nil {
			return nil, fmt.Errorf("failed to render frame: %w", err)
		}

		frames = append(frames, frame)
		frameCount++
	}

	// Add final frames (hold for 2 seconds with no cursor)
	for i := 0; i < finalFrameCount; i++ {
		progress := float64(frameCount) / float64(totalFrames)
		frame, err := renderer.RenderFrame(code.Tokens, totalChars, false, progress)
		if err != nil {
			return nil, fmt.Errorf("failed to render final frame: %w", err)
		}
		frames = append(frames, frame)
		frameCount++
	}

	return frames, nil
}
