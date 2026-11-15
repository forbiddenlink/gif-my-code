package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/forbiddenlink/gif-my-code/internal/animator"
	"github.com/forbiddenlink/gif-my-code/internal/encoder"
	"github.com/forbiddenlink/gif-my-code/internal/highlight"
	"github.com/forbiddenlink/gif-my-code/internal/parser"
	"github.com/forbiddenlink/gif-my-code/internal/render"
	"github.com/spf13/cobra"
)

var (
	theme        string
	speed        float64
	output       string
	format       string
	width        int
	fontSize     float64
	language     string
	noCursor     bool
	fps          int
	highlightStr string
	windowStyle  string
	hiDPI        bool
	lineNumbers  bool
	laser        bool
	padding      int
	listThemes   bool
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
	rootCmd.Flags().StringVar(&format, "format", "", "Output format: gif or mp4 (auto-detect from extension if not provided)")
	rootCmd.Flags().IntVarP(&width, "width", "w", 800, "Image width in pixels")
	rootCmd.Flags().Float64VarP(&fontSize, "font-size", "f", 16, "Font size")
	rootCmd.Flags().StringVarP(&language, "lang", "l", "", "Force language (auto-detect if not provided)")
	rootCmd.Flags().BoolVar(&noCursor, "no-cursor", false, "Disable cursor animation")
	rootCmd.Flags().IntVar(&fps, "fps", 30, "Frames per second")
	rootCmd.Flags().StringVar(&highlightStr, "highlight", "", "Lines to highlight (e.g., '5,7-9')")
	rootCmd.Flags().StringVar(&windowStyle, "window", "none", "Window style: macos, windows, or none")
	rootCmd.Flags().BoolVar(&hiDPI, "hidpi", false, "Render at 2x resolution (Retina scale)")
	rootCmd.Flags().BoolVar(&lineNumbers, "line-numbers", false, "Show line numbers")
	rootCmd.Flags().BoolVar(&laser, "laser", true, "Use fluid laser reveal animation instead of typing")
	rootCmd.Flags().BoolVar(&listThemes, "list-themes", false, "Print the available color themes and exit")
	rootCmd.Flags().IntVar(&padding, "padding", render.DefaultPadding, "Inner padding in pixels (at 1x scale)")
}

func run(cmd *cobra.Command, args []string) error {
	if listThemes {
		printThemes(cmd.OutOrStdout())
		return nil
	}
	if err := checkTheme(theme); err != nil {
		return err
	}
	if fps <= 0 {
		return fmt.Errorf("fps must be positive")
	}
	var err error

	// Read input (file or piped stdin)
	stdinIsTerminal := true
	if fi, statErr := os.Stdin.Stat(); statErr == nil {
		stdinIsTerminal = fi.Mode()&os.ModeCharDevice != 0
	}
	code, lang, filePath, err := resolveInput(args, os.Stdin, stdinIsTerminal, language)
	if err != nil {
		return err
	}

	fmt.Printf("📖 Reading %s (%s)\n", filePath, lang)
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

	if msg := overflowWarning(render.OverflowLines(highlighted.Tokens, fontSize, windowStyle, padding), padding); msg != "" {
		fmt.Fprintln(os.Stderr, msg)
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
		HiDPI:          hiDPI,
		LineNumbers:    lineNumbers,
		Language:       lang,
		LaserReveal:    laser,
		Padding:        padding,
	}
	frames, err := animator.GenerateFrames(highlighted, config)
	if err != nil {
		return fmt.Errorf("failed to generate frames: %w", err)
	}

	fmt.Printf("   Generated %d frames\n", len(frames))

	// Determine output format
	outputFormat := format
	if outputFormat == "" {
		outputFormat = encoder.DetectFormat(output)
	}

	// Encode to the appropriate format
	switch outputFormat {
	case "mp4":
		fmt.Println("🎬 Encoding MP4...")
		if err := encoder.EncodeMP4Pipe(frames, output, fps); err != nil {
			return fmt.Errorf("failed to encode MP4: %w", err)
		}
	default:
		fmt.Println("🎁 Encoding GIF...")
		if err := encoder.EncodeGIF(frames, output, fps); err != nil {
			return fmt.Errorf("failed to encode GIF: %w", err)
		}
	}

	// Get file size
	info, _ := os.Stat(output)
	sizeMB := float64(info.Size()) / 1024 / 1024

	fmt.Printf("\n✅ Done! Saved to: %s (%.2f MB)\n", output, sizeMB)
	return nil
}

// resolveInput returns the code, language and display name for the run.
// With a file argument it reads the file; otherwise it reads piped stdin.
func resolveInput(args []string, stdin io.Reader, stdinIsTerminal bool, langFlag string) (code, lang, name string, err error) {
	if len(args) > 0 {
		code, err = parser.ReadFile(args[0])
		if err != nil {
			return "", "", "", fmt.Errorf("failed to read file: %w", err)
		}
		lang = langFlag
		if lang == "" {
			lang = parser.DetectLanguage(args[0])
		}
		return code, lang, filepath.Base(args[0]), nil
	}

	if stdinIsTerminal {
		return "", "", "", fmt.Errorf("no input: provide a file path or pipe code on stdin (for example: cat main.go | gif-my-code --lang go)")
	}
	data, err := io.ReadAll(stdin)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read stdin: %w", err)
	}
	if len(data) == 0 {
		return "", "", "", fmt.Errorf("no input: stdin was empty")
	}
	lang = langFlag
	if lang == "" {
		lang = "text"
	}
	return string(data), lang, "stdin", nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// overflowWarning explains lines that fall outside the fixed-height canvas, or returns "" when none do.
func overflowWarning(lines, padding int) string {
	if lines == 0 {
		return ""
	}
	return fmt.Sprintf("warning: %d line(s) at the bottom do not fit the canvas at --padding %d and may be cut off. Lower --padding or shorten the input.", lines, padding)
}

// printThemes writes the available theme names, one per line.
func printThemes(w io.Writer) {
	for _, name := range highlight.ListThemes() {
		fmt.Fprintln(w, name)
	}
}

// checkTheme rejects unknown themes instead of letting chroma silently fall back.
func checkTheme(name string) error {
	if highlight.ValidateTheme(name) {
		return nil
	}
	msg := fmt.Sprintf("unknown theme %q", name)
	if near := nearestThemes(name, 3); len(near) > 0 {
		msg += fmt.Sprintf(" (did you mean: %s?)", strings.Join(near, ", "))
	}
	return fmt.Errorf("%s. Run with --list-themes to see all options", msg)
}

// nearestThemes returns up to n theme names closest to name by edit distance.
func nearestThemes(name string, n int) []string {
	name = strings.ToLower(name)
	type scored struct {
		name string
		d    int
	}
	var all []scored
	for _, t := range highlight.ListThemes() {
		d := editDistance(name, t)
		if d <= 3 || strings.Contains(t, name) {
			all = append(all, scored{t, d})
		}
	}
	sort.SliceStable(all, func(i, j int) bool { return all[i].d < all[j].d })
	var out []string
	for i := 0; i < len(all) && i < n; i++ {
		out = append(out, all[i].name)
	}
	return out
}

func editDistance(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	prev := make([]int, len(rb)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(ra); i++ {
		cur := make([]int, len(rb)+1)
		cur[0] = i
		for j := 1; j <= len(rb); j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
		}
		prev = cur
	}
	return prev[len(rb)]
}
