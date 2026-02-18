# Task 20 — Screenshot Rendering

## Goal

Implement `internal/render/screenshot.go` — render terminal ANSI output to PNG
images with embedded fonts.

## Reference

- CCBot: `src/ccbot/screenshot.py` — ANSI parsing, font tiers, `text_to_image()`.
- CCBot embeds: JetBrains Mono, Noto Sans Mono CJK, Symbola.

## Steps

1. Add dependencies: `go get golang.org/x/image github.com/golang/freetype`.
2. Embed fonts via `//go:embed`:
   - Download and embed JetBrains Mono Regular (Latin, box-drawing, symbols).
   - Download and embed Noto Sans Mono CJK SC Regular (CJK characters).
   - Download and embed Symbola (misc symbols like ⏵✔❌).
   - Place in `internal/render/fonts/` directory.
3. Create `internal/render/screenshot.go`:
   - Implement three-tier font selection:
     - Tier 0: JetBrains Mono (default for Latin + box-drawing).
     - Tier 1: Noto CJK (CJK Unicode ranges).
     - Tier 2: Symbola (specific symbol codepoints).
4. Implement ANSI SGR parsing:
   - Parse `\x1b[...m` escape sequences.
   - Handle: reset (0), fg 30-37, bg 40-47, bright fg 90-97, bright bg 100-107.
   - Handle extended: 256-color (38;5;N / 48;5;N) and RGB (38;2;R;G;B / 48;2;R;G;B).
   - Implement 256-color mapping (16 system, 6×6×6 cube, 24 grayscale).
5. Implement `RenderScreenshot(paneText string) ([]byte, error)`:
   - Parse each line for ANSI sequences → list of styled character runs.
   - Render to `image.RGBA`:
     - Background: `(30, 30, 30)` dark theme.
     - Default fg: `(212, 212, 212)`.
     - Padding: 16px.
     - Line height: `fontSize * 1.4`.
   - Two-pass: measure all lines first, then render to correctly sized image.
   - Encode as PNG via `image/png`.
   - Return PNG bytes.

## Acceptance

- Terminal output with ANSI colors renders to a readable PNG.
- All three font tiers work (Latin, CJK, symbols).
- 256-color and RGB ANSI sequences are handled.
- Box-drawing characters render correctly.
- Output is a valid PNG image.

## Phase

4 — Rich Features

## Depends on

- Task 07
