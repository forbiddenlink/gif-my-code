package encoder

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"
)

// EncodeGIF encodes frames into an animated GIF
func EncodeGIF(frames []*image.RGBA, outputPath string, fps int) error {
	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Convert frames to paletted images
	palettedFrames := make([]*image.Paletted, len(frames))
	delays := make([]int, len(frames))
	
	// Calculate delay (in 100ths of a second)
	delay := 100 / fps

	for i, frame := range frames {
		bounds := frame.Bounds()
		
		// Create paletted image with Plan9 palette (256 colors)
		palettedImage := image.NewPaletted(bounds, palette.Plan9)
		
		// Use draw to properly convert RGBA to paletted with dithering
		draw.Draw(palettedImage, bounds, frame, bounds.Min, draw.Src)
		
		palettedFrames[i] = palettedImage
		delays[i] = delay
	}

	// Encode GIF
	return gif.EncodeAll(outFile, &gif.GIF{
		Image: palettedFrames,
		Delay: delays,
	})
}
