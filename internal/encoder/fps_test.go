package encoder

import (
	"image"
	"image/gif"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInvalidFPSPreservesOutput(t *testing.T) {
	for name, encode := range map[string]func([]*image.RGBA, string, int) error{"gif": EncodeGIF, "mp4": EncodeMP4, "pipe": EncodeMP4Pipe} {
		for _, fps := range []int{0, -1} {
			t.Run(name, func(t *testing.T) {
				p := filepath.Join(t.TempDir(), "synthetic-output")
				if err := os.WriteFile(p, []byte("KEEP"), 0600); err != nil {
					t.Fatal(err)
				}
				err := encode([]*image.RGBA{image.NewRGBA(image.Rect(0, 0, 1, 1))}, p, fps)
				if err == nil || !strings.Contains(err.Error(), "fps must be positive") {
					t.Fatalf("unexpected error: %v", err)
				}
				b, e := os.ReadFile(p)
				if e != nil || string(b) != "KEEP" {
					t.Fatal("destination changed")
				}
			})
		}
	}
}
func TestPositiveGIFFrameRatesRemainSupported(t *testing.T) {
	for _, fps := range []int{1, 25, 100, 120} {
		p := filepath.Join(t.TempDir(), "synthetic.gif")
		if e := EncodeGIF([]*image.RGBA{image.NewRGBA(image.Rect(0, 0, 1, 1))}, p, fps); e != nil {
			t.Fatal(e)
		}
		f, e := os.Open(p)
		if e != nil {
			t.Fatal(e)
		}
		g, e := gif.DecodeAll(f)
		f.Close()
		if e != nil {
			t.Fatal(e)
		}
		if len(g.Image) != 1 || g.Delay[0] != 100/fps {
			t.Fatal("GIF timing behavior changed")
		}
	}
}
