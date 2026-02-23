# gif-my-code - Design Document

## Overview
CLI tool that converts code snippets into beautiful animated GIFs with syntax highlighting and typing effects.

**Target MVP:** Ship in 1-2 weekends  
**Language:** Go  
**Portfolio Impact:** Visual tool that demonstrates itself

---

## Core Features (MVP)

### Must Have
- [x] Read code from file or stdin
- [x] Syntax highlighting (multiple languages)
- [x] Typing animation effect
- [x] Export as GIF
- [x] Multiple color themes
- [x] Configurable speed
- [x] Font selection

### Nice to Have (v2)
- [ ] Line highlighting (draw attention to specific lines)
- [ ] Annotations/comments that pop in
- [ ] Export as MP4 (for higher quality)
- [ ] Window chrome (make it look like a code editor)
- [ ] Custom backgrounds
- [ ] Gradients/patterns

---

## Technical Architecture

```
┌─────────────────┐
│   CLI Input     │  (file path, options)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Code Parser    │  Read file, detect language
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Syntax Highlight│  Use chroma library
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Frame Generator │  Create image for each "frame"
└────────┬────────┘  (progressive reveal)
         │
         ▼
┌─────────────────┐
│  GIF Encoder    │  Combine frames → GIF
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  File Output    │  Save output.gif
└─────────────────┘
```

---

## Tech Stack

### Core Libraries
1. **[chroma](https://github.com/alecthomas/chroma)** - Syntax highlighting
   - Supports 250+ languages
   - 50+ built-in themes
   - Fast, mature, well-documented

2. **[gg](https://github.com/fogleman/gg)** - 2D graphics rendering
   - Simple Go API for drawing
   - Text rendering with custom fonts
   - Image composition

3. **[imaging](https://github.com/disintegration/imaging)** - Image manipulation
   - Resize, crop, effects
   - Fast and idiomatic Go

4. **[gift](https://github.com/disintegration/gift)** - GIF encoding
   - Create animated GIFs
   - Control frame delays
   - Optimize file size

5. **[cobra](https://github.com/spf13/cobra)** - CLI framework
   - Flag parsing
   - Subcommands
   - Help generation

### Fonts
- **JetBrains Mono** - Default monospace font
- Embed font files in binary for portability

---

## CLI Design

### Basic Usage
```bash
# Generate GIF from file
gif-my-code example.go

# With options
gif-my-code example.go \
  --theme dracula \
  --speed 2 \
  --output demo.gif \
  --width 800

# From stdin
cat example.go | gif-my-code --lang go

# List available themes
gif-my-code themes

# List supported languages
gif-my-code languages
```

### Flags
```
--theme, -t      Color theme (default: dracula)
--speed, -s      Typing speed multiplier (default: 1.0)
--output, -o     Output file path (default: code.gif)
--width, -w      Image width in pixels (default: 800)
--font-size, -f  Font size (default: 16)
--lang, -l       Force language (auto-detect if not provided)
--no-cursor      Disable cursor animation
--fps            Frames per second (default: 30)
```

---

## Animation Algorithm

### Typing Effect
1. Start with empty canvas
2. For each character in code:
   - Add character to visible text
   - Render frame with cursor
   - Add frame to GIF
3. Add final frame (no cursor, hold 2 seconds)

### Timing
- **Characters per frame:** Based on speed multiplier
  - Speed 0.5: 1 char per 2 frames (slow)
  - Speed 1.0: 2 chars per frame (normal)
  - Speed 2.0: 4 chars per frame (fast)
- **FPS:** 30 (smooth but not too large)
- **Final hold:** 60 frames (2 seconds)

### Cursor Animation
- Blinking cursor at current position
- Alternate visible/invisible every 15 frames
- Small rectangle after last character

---

## Code Structure

```
gif-my-code/
├── cmd/
│   └── root.go           # CLI entry point
├── internal/
│   ├── parser/
│   │   └── parser.go     # Read & detect language
│   ├── highlight/
│   │   └── highlight.go  # Syntax highlighting
│   ├── render/
│   │   ├── renderer.go   # Generate image frames
│   │   └── theme.go      # Theme definitions
│   ├── animator/
│   │   └── animator.go   # Animation logic
│   └── encoder/
│       └── encoder.go    # GIF encoding
├── assets/
│   └── fonts/
│       └── JetBrainsMono.ttf
├── examples/
│   ├── example.go
│   ├── example.py
│   └── example.js
├── go.mod
├── go.sum
├── main.go
├── DESIGN.md
└── README.md
```

---

## Implementation Steps

### Phase 1: Foundation (Weekend 1, Day 1)
1. Initialize Go project
2. Set up Cobra CLI structure
3. Install dependencies
4. Create basic file reader
5. Test syntax highlighting with chroma

### Phase 2: Rendering (Weekend 1, Day 2)
1. Set up gg for image rendering
2. Render static code image (no animation)
3. Test different themes
4. Embed font in binary

### Phase 3: Animation (Weekend 2, Day 1)
1. Implement frame generation
2. Add typing effect logic
3. Add cursor animation
4. Test with different speeds

### Phase 4: GIF Export (Weekend 2, Day 2)
1. Combine frames into GIF
2. Optimize file size
3. Add CLI flags
4. Test across languages
5. Create examples

### Phase 5: Polish
1. Error handling
2. Progress indicator
3. Documentation
4. Demo GIFs for README

---

## Testing Strategy

### Test Files
Create examples in multiple languages:
- `examples/example.go` - Go
- `examples/example.py` - Python
- `examples/example.js` - JavaScript
- `examples/example.rs` - Rust
- `examples/example.tsx` - TypeScript/React

### Quality Checks
- [ ] Works on Mac/Linux/Windows
- [ ] File sizes reasonable (<5MB for 50 lines)
- [ ] Smooth animation at 30fps
- [ ] All themes render correctly
- [ ] Language detection accurate
- [ ] Handles Unicode/emoji in code

---

## Launch Strategy

### Portfolio Integration
1. Generate GIF of MCP Server Studio code
2. Generate GIF of Kindred App logic
3. Add to portfolio projects with "Made with gif-my-code"

### Social Launch
1. Post demo GIF on Twitter
2. GitHub repo with examples
3. Product Hunt submission
4. Dev.to article: "I built a tool to make code GIFs"

### Viral Loop
- Add watermark option (optional, removable in paid version)
- "Made with gif-my-code" in corner
- Link to repo in every shared GIF

---

## Future Monetization (Post-MVP)

### Free Tier
- CLI tool (open source)
- Basic themes
- GIF export only
- Watermark

### Pro Tier ($5/mo or $2/GIF)
- Web version (no CLI needed)
- Premium themes
- No watermark
- MP4 export (higher quality)
- Custom fonts
- Line annotations
- Window chrome
- API access

### Enterprise ($49/mo)
- Batch processing
- Custom branding
- Priority support
- Team features

---

## Success Metrics

### Week 1
- 100 GitHub stars
- 50 GIFs generated
- 10 social shares

### Month 1
- 500 GitHub stars
- 1000 GIFs generated
- Featured on Product Hunt
- 10 paying customers

### Month 3
- 2000 GitHub stars
- Used in 50+ project READMEs
- $500 MRR
- Featured in newsletter (Go Weekly, etc.)

---

## Open Questions

1. **File size optimization:** How to keep GIFs small?
   - Use palette reduction
   - Skip redundant frames
   - Add compression

2. **Language detection:** Use file extension or content analysis?
   - Start with extension (simpler)
   - Add content analysis if needed

3. **Unicode handling:** How to render emoji/special chars?
   - Use proper Unicode font
   - Test extensively

---

## Similar Tools (Competition Analysis)

- **Carbon** (carbon.now.sh) - Static images only, no animation
- **Ray.so** - Static images, beautiful but not animated
- **Asciinema** - Terminal recording, not code-focused
- **GitHub Gists** - Static, no animation

**Our advantage:** Only tool that creates animated code GIFs with typing effect

---

## Timeline

**Weekend 1 (Feb 29-Mar 1):**
- Day 1: Project setup, syntax highlighting working
- Day 2: Static rendering working

**Weekend 2 (Mar 7-8):**
- Day 1: Animation working
- Day 2: GIF export + polish

**Launch:** Monday, March 9, 2026

---

## Notes

- Keep it simple for v1 - typing animation + syntax highlighting is enough
- Focus on making examples beautiful (portfolio pieces)
- Documentation is marketing - make README amazing
- Every feature should demo itself

