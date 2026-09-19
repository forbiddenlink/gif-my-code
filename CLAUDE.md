# gif-my-code

CLI that converts a code file into an animated GIF or MP4 with syntax highlighting, line
highlighting, and window chrome (macOS/Windows styling). Free/open-source alternative to
Snappify/Carbon/Ray.so, positioned on being scriptable and offline-capable.

## Stack

Go 1.25. Cobra for the CLI, chroma for syntax highlighting (250+ languages, 50+ themes),
fogleman/gg + golang/freetype for rendering, golang.org/x/image for image encoding.

## Commands

```bash
go build ./...
go test ./...
go build -o gif-my-code .
./gif-my-code example.go --highlight "5,7-9" --line-numbers
```

## CLI flags (`cmd/root.go`)

`gif-my-code [file]` with: `--theme` (`-t`, default dracula), `--speed` (`-s`, default 1.0),
`--output` (`-o`, default code.gif), `--format` (gif or mp4, auto-detected from output
extension), `--width` (`-w`, default 800), `--font-size` (`-f`, default 16), `--lang` (`-l`,
auto-detect if unset), `--no-cursor`, `--fps` (default 30), `--highlight` (e.g. `5,7-9`),
`--window` (macos, windows, or none), `--hidpi` (2x/Retina render), `--line-numbers`, `--laser`
(fluid laser reveal instead of typing animation, default true).

## Layout

- `main.go` - entry point, calls `cmd.Execute()`
- `cmd/root.go` - Cobra command definition and flag parsing
- `internal/parser/` - source file parsing
- `internal/highlight/` - chroma-based syntax + line highlighting
- `internal/animator/` - frame sequencing for the typing/laser animations
- `internal/render/` - frame rendering (fonts, window chrome)
- `internal/encoder/` - GIF encoding (`encoder.go`) and MP4 encoding (`mp4.go`)
- `demos/` - example output GIFs referenced from README

## CI

`.github/workflows/ci.yml` runs `go build ./...` and `go test ./...` on push/PR to `main`.
`dependabot-automerge.yml` waits for all checks on a dependabot PR's head SHA (there is no
branch protection on this repo) before squash-merging patch/minor updates; majors stay manual.
