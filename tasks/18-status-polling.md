# Task 18 — Status Polling

## Goal

Implement `internal/monitor/terminal.go` (status detection) and status polling
logic — detect Claude's spinner/status line from the terminal, show as an
editable message in the topic.

## Reference

- CCBot: `src/ccbot/terminal_parser.py` — `strip_pane_chrome()`, status line detection.
- CCBot: `src/ccbot/handlers/status_polling.py` — `status_poll_loop()`, 1-second interval.

## Steps

1. Create `internal/monitor/terminal.go`:
   - Implement `StripPaneChrome(paneText string) string`:
     - Find the chrome separator: line of `─` chars (≥20 wide) in the last 10 lines.
     - Return text above the separator.
   - Implement `ExtractStatusLine(paneText string) (string, bool)`:
     - After stripping chrome, find the line containing a spinner character
       (`·✻✽✶✳✢`) in the last few lines above the separator.
     - Return the status text and whether one was found.
2. Implement status polling loop (can be in `internal/bot/` or `internal/monitor/`):
   - Run as a goroutine, 1-second interval.
   - For each thread-bound window:
     - Skip if user's message queue is non-empty (avoid status noise during content delivery).
     - Capture pane (plain text, no ANSI).
     - Extract status line.
     - If status found and different from last: send/edit status message in topic.
     - If no status and had one before: clear status message.
   - Deduplication: track last status text per (user, thread) to avoid redundant edits.
3. Implement topic existence probe (every 60 seconds):
   - For each bound thread, call a lightweight Telegram API method.
   - On error indicating deleted topic: clean up bindings (same as topic close).

## Acceptance

- Status line is detected from Claude's terminal output.
- Status message appears/updates in real-time in the Telegram topic.
- Status message is cleared when Claude finishes.
- Duplicate status texts don't trigger edits.
- Deleted topics are detected and cleaned up.

## Phase

3 — Session Monitor

## Depends on

- Task 17
