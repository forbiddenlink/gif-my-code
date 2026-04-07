package encoder

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// EncodeMP4 encodes frames into an MP4 video using ffmpeg
func EncodeMP4(frames []*image.RGBA, outputPath string, fps int) error {
	if len(frames) == 0 {
		return fmt.Errorf("no frames to encode")
	}

	// Check if ffmpeg is available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH - please install ffmpeg")
	}

	// Create a temporary directory for PNG frames
	tempDir, err := os.MkdirTemp("", "gif-my-code-frames-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Write frames as PNG files
	for i, frame := range frames {
		framePath := filepath.Join(tempDir, fmt.Sprintf("frame_%06d.png", i))
		if err := writePNG(frame, framePath); err != nil {
			return fmt.Errorf("failed to write frame %d: %w", i, err)
		}
	}

	// Build ffmpeg command
	// -framerate: input frame rate
	// -i: input pattern
	// -c:v libx264: H.264 codec
	// -pix_fmt yuv420p: pixel format for compatibility
	// -crf 23: quality (0-51, lower is better, 23 is default)
	// -preset medium: encoding speed/quality tradeoff
	// -y: overwrite output without asking
	args := []string{
		"-framerate", fmt.Sprintf("%d", fps),
		"-i", filepath.Join(tempDir, "frame_%06d.png"),
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-crf", "23",
		"-preset", "medium",
		"-y",
		outputPath,
	}

	cmd := exec.Command("ffmpeg", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w\n%s", err, stderr.String())
	}

	return nil
}

// EncodeMP4Pipe encodes frames by piping raw video data to ffmpeg
// This is more efficient for large videos as it avoids writing temp files
func EncodeMP4Pipe(frames []*image.RGBA, outputPath string, fps int) error {
	if len(frames) == 0 {
		return fmt.Errorf("no frames to encode")
	}

	// Check if ffmpeg is available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH - please install ffmpeg")
	}

	// Get dimensions from first frame
	bounds := frames[0].Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Build ffmpeg command for piped input
	// -f rawvideo: input format
	// -pix_fmt rgba: input pixel format (matches image.RGBA)
	// -s WxH: input dimensions
	args := []string{
		"-f", "rawvideo",
		"-pix_fmt", "rgba",
		"-s", fmt.Sprintf("%dx%d", width, height),
		"-framerate", fmt.Sprintf("%d", fps),
		"-i", "-", // read from stdin
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-crf", "23",
		"-preset", "medium",
		"-y",
		outputPath,
	}

	cmd := exec.Command("ffmpeg", args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Write raw RGBA data for each frame
	for _, frame := range frames {
		if _, err := stdin.Write(frame.Pix); err != nil {
			stdin.Close()
			cmd.Wait()
			return fmt.Errorf("failed to write frame data: %w", err)
		}
	}

	stdin.Close()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w\n%s", err, stderr.String())
	}

	return nil
}

// writePNG writes an RGBA image to a file as PNG
func writePNG(img *image.RGBA, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// DetectFormat returns the output format based on file extension
func DetectFormat(outputPath string) string {
	ext := filepath.Ext(outputPath)
	switch ext {
	case ".mp4", ".m4v":
		return "mp4"
	case ".gif":
		return "gif"
	case ".webm":
		return "webm"
	default:
		return "gif" // default to GIF
	}
}

// Encode encodes frames to the appropriate format based on output path
func Encode(frames []*image.RGBA, outputPath string, fps int) error {
	format := DetectFormat(outputPath)
	switch format {
	case "mp4":
		return EncodeMP4Pipe(frames, outputPath, fps)
	case "gif":
		return EncodeGIF(frames, outputPath, fps)
	default:
		return EncodeGIF(frames, outputPath, fps)
	}
}

// Placeholder to satisfy compiler if io is imported but not used
var _ = io.EOF
