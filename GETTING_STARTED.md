# Getting Started with gif-my-code

## Quick Start (5 minutes)

### Step 1: Build the Project
```bash
cd /Volumes/LizsDisk/gif-my-code
go build -o gif-my-code
```

### Step 2: Test with Examples
```bash
# Generate your first GIF
./gif-my-code examples/example.go

# You should see output like:
# 📖 Reading example.go (go)
# 🎨 Theme: dracula
# ⚡ Speed: 1.0x
# ✨ Applying syntax highlighting...
# 🎬 Generating animation frames...
#    Generated 156 frames
# 🎁 Encoding GIF...
# ✅ Done! Saved to: code.gif (2.34 MB)
```

### Step 3: View Your GIF
```bash
open code.gif
# or
xdg-open code.gif  # Linux
```

## Next Steps

### Try Different Themes
```bash
# List all available themes
./gif-my-code themes

# Try some popular ones
./gif-my-code examples/example.py --theme monokai -o monokai.gif
./gif-my-code examples/example.py --theme nord -o nord.gif
./gif-my-code examples/example.py --theme github -o github.gif
```

### Adjust Speed
```bash
# Slower (0.5x)
./gif-my-code examples/example.tsx --speed 0.5 -o slow.gif

# Normal (1.0x) - default
./gif-my-code examples/example.tsx -o normal.gif

# Fast (2.0x)
./gif-my-code examples/example.tsx --speed 2.0 -o fast.gif

# Very fast (3.0x)
./gif-my-code examples/example.tsx --speed 3.0 -o veryfast.gif
```

### Customize Appearance
```bash
# Larger size
./gif-my-code examples/example.go --width 1200 --font-size 20 -o large.gif

# Compact
./gif-my-code examples/example.go --width 600 --font-size 14 -o compact.gif

# No cursor
./gif-my-code examples/example.go --no-cursor -o no-cursor.gif
```

## Create Your Own

### 1. Create a Code Snippet
```bash
cat > hello.py << 'EOF'
def greet(name):
    print(f"Hello, {name}!")
    print("Welcome to gif-my-code! 🎨")

greet("World")
EOF
```

### 2. Generate GIF
```bash
./gif-my-code hello.py --theme dracula --speed 1.5 --output hello.gif
```

### 3. Share!
- Add to your README
- Post on Twitter/LinkedIn
- Include in documentation
- Show off in your portfolio

## Common Workflows

### For README Examples
```bash
# Make it quick and compact
./gif-my-code my-feature.js \
  --speed 2.0 \
  --width 700 \
  --theme github \
  --output readme-demo.gif
```

### For Social Media
```bash
# Make it eye-catching
./gif-my-code cool-algorithm.py \
  --speed 1.5 \
  --width 900 \
  --font-size 18 \
  --theme dracula \
  --output social-demo.gif
```

### For Tutorials
```bash
# Make it slow and clear
./gif-my-code tutorial-example.tsx \
  --speed 0.8 \
  --width 1000 \
  --font-size 17 \
  --theme nord \
  --output tutorial-step-1.gif
```

### For Portfolio
```bash
# Make it cinematic
./gif-my-code portfolio-project.go \
  --speed 1.2 \
  --width 1100 \
  --font-size 19 \
  --theme tokyo-night \
  --output portfolio-demo.gif
```

## Troubleshooting

### "command not found: gif-my-code"
Make sure you built the binary and you're in the right directory:
```bash
cd /Volumes/LizsDisk/gif-my-code
./gif-my-code examples/example.go
```

### GIF is too large
Try:
- Reducing the width: `--width 600`
- Increasing speed: `--speed 2.0`
- Using a shorter code snippet

### Animation is too fast/slow
Adjust speed:
- Too fast? Lower speed: `--speed 0.5`
- Too slow? Increase speed: `--speed 2.0`

### Colors look wrong
Try a different theme:
```bash
./gif-my-code examples/example.go --theme monokai
```

## Tips & Tricks

### 1. Keep Snippets Short
15-30 lines work best. Longer code = larger GIF files.

### 2. Choose the Right Theme
- Light themes: `github`, `solarized-light`
- Dark themes: `dracula`, `monokai`, `nord`
- Match your brand colors

### 3. Speed Guidelines
- Documentation: 1.0-1.5x (readable)
- Social media: 1.5-2.5x (engaging)
- Quick demos: 2.0-3.0x (snappy)

### 4. File Size
Typical sizes:
- 10 lines: ~1-2 MB
- 20 lines: ~2-4 MB
- 30 lines: ~4-6 MB

### 5. Test Before Sharing
Always view your GIF before posting to ensure:
- Text is readable
- Speed feels right
- Colors look good
- No awkward pauses

## Next: Building Your Portfolio

Once you're comfortable, use gif-my-code to create demos for:
1. Your GitHub README files
2. Portfolio project showcases
3. Social media posts about your work
4. Technical blog posts
5. Documentation examples

Happy coding! 🎬✨
