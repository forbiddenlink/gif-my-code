# gif-my-code 🎬

> Convert code snippets into beautiful animated GIFs with syntax highlighting and typing effects.

Perfect for Twitter posts, README files, documentation, tutorials, and portfolio demos.

## ✨ Features

- 🎨 **Syntax highlighting** for 250+ languages
- 🌈 **50+ color themes** (Dracula, Monokai, Nord, and more)
- ⚡ **Customizable speed** - Control typing animation speed
- 📐 **Flexible sizing** - Set width and font size
- 🎯 **Smart language detection** - Automatic from file extension
- 💾 **Optimized output** - Reasonable file sizes
- 📍 **Line highlighting** - Draw attention to specific lines
- 🪟 **Window chrome** - macOS or Windows style (NEW!)
- 🎨 **Gradient backgrounds** - Professional visual polish (NEW!)
- ✨ **Drop shadows** - Depth and dimension (NEW!)

## 🚀 Installation

### From Source
```bash
git clone https://github.com/elizabethstein/gif-my-code.git
cd gif-my-code
go build -o gif-my-code
```

### Using Go Install
```bash
go install github.com/elizabethstein/gif-my-code@latest
```

## 📖 Usage

### Basic Usage
```bash
# Generate GIF from a file
gif-my-code example.go

# Specify output path
gif-my-code example.py --output demo.gif

# Use different theme
gif-my-code example.js --theme monokai

# Adjust speed (higher = faster)
gif-my-code example.tsx --speed 2.0

# Combine options
gif-my-code example.rs \
  --theme dracula \
  --speed 1.5 \
  --width 1000 \
  --font-size 18 \
  --output rust-demo.gif
```

### Options
```
  -t, --theme string       Color theme (default "dracula")
  -s, --speed float        Typing speed multiplier (default 1.0)
  -o, --output string      Output file path (default "code.gif")
  -w, --width int          Image width in pixels (default 800)
  -f, --font-size float    Font size (default 16)
  -l, --lang string        Force language (auto-detect if not provided)
      --highlight string   Lines to highlight (e.g., '5,7-9')
      --window string      Window style: macos, windows, or none (default "none")
      --no-cursor          Disable cursor animation
      --fps int            Frames per second (default 30)
```

### List Available Themes
```bash
gif-my-code themes
```

### Supported Languages
Auto-detection works for 50+ languages including:
- Go, Python, JavaScript, TypeScript, Rust, Ruby
- Java, C, C++, C#, Swift, Kotlin
- HTML, CSS, JSON, YAML, Markdown
- And many more...

## 🎨 Examples

### React Component
```bash
gif-my-code examples/example.tsx \
  --theme dracula \
  --speed 1.5 \
  --output react-component.gif
```

### Python Script
```bash
gif-my-code examples/example.py \
  --theme monokai \
  --speed 2.0 \
  --output python-script.gif
```

### Go Program
```bash
gif-my-code examples/example.go \
  --theme nord \
  --output go-program.gif
```

### With Line Highlighting
```bash
# Highlight specific lines (great for tutorials!)
gif-my-code examples/example.py \
  --highlight 5,7-9 \
  --theme dracula \
  --output tutorial-demo.gif

# Perfect for bug explanations
gif-my-code buggy-code.js \
  --highlight 12 \
  --theme monokai \
  --output bug-fix-demo.gif
```

### With Window Chrome (Professional Look!)
```bash
# macOS style window
gif-my-code examples/example.go \
  --window macos \
  --theme dracula \
  --output macos-demo.gif

# Combine window + highlighting
gif-my-code examples/example.py \
  --window macos \
  --highlight 3,5-7 \
  --theme monokai \
  --speed 1.5 \
  --output professional-demo.gif

# Windows style
gif-my-code examples/example.tsx \
  --window windows \
  --theme github \
  --output windows-demo.gif
```

## 🎯 Use Cases

### Twitter/LinkedIn Posts
Create eye-catching code demos for social media.

### README Files
Add animated examples to your project documentation.

### Tutorials
Show code snippets coming to life in educational content.

### Portfolio
Demonstrate your projects with cinematic code reveals.

## 🛠️ Development

### Project Structure
```
gif-my-code/
├── cmd/                  # CLI commands
├── internal/
│   ├── parser/          # File reading & language detection
│   ├── highlight/       # Syntax highlighting
│   ├── render/          # Image rendering
│   ├── animator/        # Frame generation
│   └── encoder/         # GIF encoding
├── examples/            # Example code files
└── assets/              # Fonts and resources
```

### Building
```bash
go build -o gif-my-code
```

### Testing
```bash
# Test with examples
./gif-my-code examples/example.go
./gif-my-code examples/example.py
./gif-my-code examples/example.tsx
```

## 🎨 Popular Themes

- `dracula` - Dark theme with vibrant colors
- `monokai` - Classic Sublime Text theme
- `nord` - Cool, calm Nordic palette
- `github` - GitHub's light theme
- `solarized-dark` - Low contrast dark
- `gruvbox` - Retro groove colors
- `tokyo-night` - Modern dark theme

Run `gif-my-code themes` for the full list.

## 🚧 Roadmap

### v1.0 ✅ COMPLETE
- [x] Basic typing animation
- [x] Syntax highlighting
- [x] Multiple themes
- [x] CLI interface
- [x] Line highlighting
- [x] Window chrome (macOS & Windows)
- [x] Gradient backgrounds
- [x] Drop shadows

### v1.1 (Coming Soon - Based on Feedback!)
- [ ] **Line numbers toggle** - Most requested
- [ ] **Custom window titles** - Easy polish
- [ ] **Export resolution** (2x/3x for retina)
- [ ] **Watermark option** (optional branding)
- [ ] **No-shadow toggle** (for flat look)

### v1.2 (Roadmap)
- [ ] Diff mode (added/removed lines in green/red)
- [ ] Custom fonts (JetBrains Mono, Fira Code)
- [ ] MP4 export (smaller files, higher quality)
- [ ] Stdin support (pipe code directly)

### v2.0 (Future)
- [ ] Annotations (arrows, boxes, comments)
- [ ] Web version (no CLI needed)
- [ ] Batch processing
- [ ] API access

**Want a feature? [Open an issue](https://github.com/elizabethstein/gif-my-code/issues)!**

## 🤝 Contributing

Contributions welcome! Feel free to:
- Report bugs
- Suggest features
- Submit pull requests
- Share your GIFs made with this tool

## 📝 License

MIT License - see LICENSE file for details.

## 🙏 Acknowledgments

Built with:
- [chroma](https://github.com/alecthomas/chroma) - Syntax highlighting
- [gg](https://github.com/fogleman/gg) - 2D graphics
- [cobra](https://github.com/spf13/cobra) - CLI framework

## 💬 Questions?

- Open an issue on GitHub
- Tweet [@elizabethstein](https://twitter.com/elizabethstein)
- Check out my [portfolio](https://elizabethannstein.com)

---

**Made with 🐾 by [Elizabeth Stein](https://github.com/elizabethstein)**

*If this tool helped you, consider starring the repo!* ⭐
