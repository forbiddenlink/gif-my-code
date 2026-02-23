package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/forbiddenlink/gif-my-code/internal/animator"
	"github.com/forbiddenlink/gif-my-code/internal/encoder"
	"github.com/forbiddenlink/gif-my-code/internal/highlight"
	"github.com/forbiddenlink/gif-my-code/internal/parser"
	"github.com/spf13/cobra"
)

var (
	theme        string
	speed        float64
	output       string
	width        int
	fontSize     float64
	language     string
	noCursor     bool
	fps          int
	highlightStr string
	windowStyle  string
)

var rootCmd = &cobra.Command{
	Use:   "gif-my-code [file]",
	Short: "Convert code snippets into beautiful animated GIFs",
	Long: `gif-my-code creates animated GIFs of your code with syntax highlighting and typing effects.
	
Perfect for:
  • Twitter/LinkedIn posts
  • README files
  • Documentation
  • Tutorials
  • Portfolio demos`,
	Args: cobra.MaximumNArgs(1),
	RunE: run,
}

func init() {
	rootCmd.Flags().StringVarP(&theme, "theme", "t", "dracula", "Color theme")
	rootCmd.Flags().Float64VarP(&speed, "speed", "s", 1.0, "Typing speed multiplier")
	rootCmd.Flags().StringVarP(&output, "output", "o", "code.gif", "Output file path")
	rootCmd.Flags().IntVarP(&width, "width", "w", 800, "Image width in pixels")
	rootCmd.Flags().Float64VarP(&fontSize, "font-size", "f", 16, "Font size")
	rootCmd.Flags().StringVarP(&language, "lang", "l", "", "Force language (auto-detect if not provided)")
	rootCmd.Flags().BoolVar(&noCursor, "no-cursor", false, "Disable cursor animation")
	rootCmd.Flags().IntVar(&fps, "fps", 30, "Frames per second")
	rootCmd.Flags().StringVar(&highlightStr, "highlight", "", "Lines to highlight (e.g., '5,7-9')")
	rootCmd.Flags().StringVar(&windowStyle, "window", "none", "Window style: macos, windows, or none")
}

func run(cmd *cobra.Command, args []string) error {
	var filePath string
	var code string
	var lang string
	var err error

	// Read input (file or stdin)
	if len(args) == 0 {
		// TODO: Read from stdin
		return fmt.Errorf("stdin support coming soon - please provide a file path")
	} else {
		filePath = args[0]
		code, err = parser.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		// Detect language from file extension if not provided
		if language == "" {
			lang = parser.DetectLanguage(filePath)
		} else {
			lang = language
		}
	}

	fmt.Printf("📖 Reading %s (%s)\n", filepath.Base(filePath), lang)
	fmt.Printf("🎨 Theme: %s\n", theme)
	fmt.Printf("⚡ Speed: %.1fx\n", speed)

	// Syntax highlight
	fmt.Println("✨ Applying syntax highlighting...")
	highlighted, err := highlight.Highlight(code, lang, theme)
	if err != nil {
		return fmt.Errorf("failed to highlight code: %w", err)
	}

	// Parse highlight lines
	var highlightLines []int
	if highlightStr != "" {
		highlightLines, err = parser.ParseHighlightLines(highlightStr)
		if err != nil {
			return fmt.Errorf("invalid highlight format: %w", err)
		}
		fmt.Printf("📍 Highlighting lines: %v\n", highlightLines)
	}

	// Generate frames
	fmt.Println("🎬 Generating animation frames...")
	if windowStyle != "none" && windowStyle != "" {
		fmt.Printf("🪟 Window style: %s\n", windowStyle)
	}
	config := animator.Config{
		Width:          width,
		FontSize:       fontSize,
		Speed:          speed,
		FPS:            fps,
		ShowCursor:     !noCursor,
		HighlightLines: highlightLines,
		WindowStyle:    windowStyle,
		Theme:          theme,
	}
	frames, err := animator.GenerateFrames(highlighted, config)
	if err != nil {
		return fmt.Errorf("failed to generate frames: %w", err)
	}

	fmt.Printf("   Generated %d frames\n", len(frames))

	// Encode GIF
	fmt.Println("🎁 Encoding GIF...")
	if err := encoder.EncodeGIF(frames, output, fps); err != nil {
		return fmt.Errorf("failed to encode GIF: %w", err)
	}

	// Get file size
	info, _ := os.Stat(output)
	sizeMB := float64(info.Size()) / 1024 / 1024

	fmt.Printf("\n✅ Done! Saved to: %s (%.2f MB)\n", output, sizeMB)
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
