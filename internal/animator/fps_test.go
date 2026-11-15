package animator

import (
	"github.com/forbiddenlink/gif-my-code/internal/highlight"
	"strings"
	"testing"
)

func TestInvalidFPSRejectedBeforeRendering(t *testing.T) {
	for _, fps := range []int{0, -1} {
		_, err := GenerateFrames(nil, Config{FPS: fps})
		if err == nil || !strings.Contains(err.Error(), "fps must be positive") {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}
func TestOneFPSAnimation(t *testing.T) {
	frames, err := GenerateFrames(&highlight.HighlightedCode{Tokens: []highlight.Token{{Text: "x"}}}, Config{FPS: 1, Width: 80, FontSize: 10, Speed: 1, WindowStyle: "none", Theme: "dracula", ShowCursor: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(frames) != 3 {
		t.Fatalf("want one typing plus two hold frames, got %d", len(frames))
	}
}
