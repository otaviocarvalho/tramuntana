# Task 03 — Tmux Package

## Goal

Implement `internal/tmux/tmux.go` with all tmux management functions. Reuse patterns
from minuano's tmux package where applicable.

## Reference

- Minuano: `internal/tmux/tmux.go` at `/home/otavio/code/minuano/internal/tmux/tmux.go`
- CCBot: `src/ccbot/tmux_manager.py` — async libtmux wrapper.

## Steps

1. Create `internal/tmux/tmux.go`.
2. Implement these functions (all shell out to the `tmux` binary via `os/exec`):
   - `SessionExists(name string) bool`
   - `EnsureSession(name string) error` — create if missing, no-op if exists.
   - `ListWindows(session string) ([]Window, error)` — return list of `{ID, Name, CWD}`.
   - `NewWindow(session, name, dir string, env map[string]string) (string, error)` — create window, return window ID. Set env vars and start Claude command.
   - `SendKeys(session, windowID, keys string) error` — send literal text to window.
   - `SendKeysWithDelay(session, windowID, text string, delayMs int) error` — send text, wait, then send Enter. CCBot uses 500ms delay between text and Enter.
   - `CapturePane(session, windowID string, withAnsi bool) (string, error)` — capture visible pane content. `withAnsi=true` for screenshots (`-e` flag), `false` for text parsing.
   - `KillWindow(session, windowID string) error`
   - `DisplayMessage(paneID, format string) (string, error)` — `tmux display-message -t PANE -p FORMAT`. Used by hook to get session:window_id:window_name.
   - `RenameWindow(session, windowID, newName string) error`
3. `Window` struct: `ID string` (e.g. `@12`), `Name string`, `CWD string`.
4. Handle errors from tmux cleanly — wrap with context.

## Acceptance

- All functions compile.
- `EnsureSession` + `NewWindow` + `SendKeys` can create a session, add a window, and send text.
- `CapturePane` returns visible output (with and without ANSI).
- `ListWindows` returns correctly parsed window info.

## Phase

1 — Foundation

## Depends on

- Task 01
