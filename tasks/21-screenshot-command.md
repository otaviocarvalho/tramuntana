# Task 21 — Screenshot Command

## Goal

Implement the `/screenshot` command — capture the tmux pane, render to PNG,
send as document with a control keyboard.

## Reference

- CCBot: `src/ccbot/bot.py` — `/screenshot` handler, refresh callback,
  quick key callbacks with screenshot update.

## Steps

1. Add to `internal/bot/commands.go` or create `internal/bot/screenshot.go`.
2. Implement `/screenshot` handler:
   - Resolve window for user's thread.
   - Capture pane with ANSI: `tmux.CapturePane(session, windowID, true)`.
   - Render to PNG: `render.RenderScreenshot(paneText)`.
   - Send as Telegram document (photo compresses too much) with inline keyboard.
3. Build screenshot control keyboard:
   - Row 1: ↑ ↓ ← → (arrow keys)
   - Row 2: Space, Tab, Esc, Enter
   - Row 3: Refresh 🔄
   - Callback data: `ss_up`, `ss_down`, `ss_left`, `ss_right`, `ss_space`,
     `ss_tab`, `ss_esc`, `ss_enter`, `ss_refresh`.
4. Handle screenshot callbacks:
   - For key buttons: send the key to tmux, wait 500ms, capture + re-render + edit document.
   - For refresh: just capture + re-render + edit document.
   - Editing a document in Telegram: use `editMessageMedia` with new PNG.
5. Track screenshot message per (user, thread) to know which message to edit.

## Acceptance

- `/screenshot` sends a PNG image of the terminal.
- Control keyboard allows sending keys and auto-refreshes the image.
- Refresh button updates the screenshot without sending keys.
- ANSI colors are preserved in the screenshot.

## Phase

4 — Rich Features

## Depends on

- Task 20
