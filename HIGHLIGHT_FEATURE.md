# Line Highlighting Feature 📍

## Overview
Draw attention to specific lines of code by adding a yellow highlight background. Perfect for tutorials, bug explanations, and feature demos.

## Usage

### Basic Syntax
```bash
gif-my-code example.go --highlight 5
```
Highlights line 5.

### Multiple Lines
```bash
gif-my-code example.go --highlight 5,7,10
```
Highlights lines 5, 7, and 10.

### Line Ranges
```bash
gif-my-code example.go --highlight 5-8
```
Highlights lines 5, 6, 7, and 8.

### Combined
```bash
gif-my-code example.go --highlight 3,5-8,12
```
Highlights lines 3, 5, 6, 7, 8, and 12.

## Examples

### Tutorial: Explain a Bug
```bash
gif-my-code buggy-code.js \
  --highlight 12 \
  --theme dracula \
  --speed 1.5 \
  --output bug-explanation.gif
```
Caption: "The bug is on line 12 - we're not checking for null! 🐛"

### Feature Demo: Show Changes
```bash
gif-my-code new-feature.py \
  --highlight 45-52 \
  --theme monokai \
  --speed 2.0 \
  --output new-feature-demo.gif
```
Caption: "Lines 45-52 show the new authentication logic ✨"

### Code Review: Point Out Issues
```bash
gif-my-code pull-request.go \
  --highlight 8,15,23 \
  --theme github \
  --output code-review-feedback.gif
```
Caption: "Three potential improvements on lines 8, 15, and 23 📝"

## Visual Style

- **Highlight Color:** Semi-transparent yellow (`rgba(255, 255, 0, 60)`)
- **Spans:** Full width of the frame
- **Animation:** Highlights appear as the typing animation reveals each line
- **Theme Compatible:** Works with all 50+ themes

## Tips

### For Tutorials
- Highlight the most important line
- Keep highlights to 1-3 lines max
- Pair with descriptive captions

### For Bug Explanations
- Highlight the problematic line
- Show before/after in separate GIFs
- Use `--speed 1.0` (slower) for clarity

### For Feature Demos
- Highlight the new/changed code
- Use ranges for multi-line features
- Use `--speed 1.5-2.0` (faster) to keep it snappy

## Limitations (v1.0)

- Fixed highlight color (yellow)
- Cannot customize highlight opacity
- No support for multiple highlight colors
- No annotations/comments (yet)

## Coming Soon (v1.1)

- [ ] Custom highlight colors (`--highlight-color "#ff00ff"`)
- [ ] Multiple colors (`--highlight "5:yellow,7:red"`)
- [ ] Per-line opacity
- [ ] Annotations (`--annotate "5:This is the bug"`)
- [ ] Highlight animations (fade in, pulse, etc.)

## Why This Feature Matters

**Differentiation:**
- Snappify ($12-30/mo) has this
- Carbon.now.sh doesn't (static only)
- Your tool: FREE + CLI + animated + highlights ✨

**Use Cases:**
1. Educational content creators
2. Code reviewers
3. Tutorial writers
4. Documentation authors
5. Social media posts
6. Conference slides

**Portfolio Impact:**
- Shows attention to detail
- Demonstrates UX thinking
- Proves you understand user needs
- More valuable than basic tool

---

**Status:** ✅ Implemented in v1.0
**Tested:** ✅ Working with `examples/example.go`
**Output:** `highlighted.gif` (1.48 MB, lines 5,7-9 highlighted)
