# Task 23 — Bash Capture

## Goal

Implement `!command` bash capture — send the command to Claude's bash mode,
capture terminal output for 30 seconds, and stream it as an editable message.

## Reference

- CCBot: `src/ccbot/bot.py` — `_capture_bash_output()`, 30s poll loop, message editing.
- CCBot: `src/ccbot/terminal_parser.py` — `extract_bash_output()`.

## Steps

1. Add to `internal/monitor/terminal.go`:
   - Implement `ExtractBashOutput(paneText, command string) (string, bool)`:
     - Strip bottom chrome via `StripPaneChrome()`.
     - Search from bottom upward for `"! <cmd_prefix>"` or `"!<cmd_prefix>"`
       (match on first 10 chars of command to handle terminal truncation).
     - Return command echo + output below it.

2. Implement bash capture in `internal/bot/handlers.go` (or separate file):
   - When text starts with `!` (already detected in Task 08):
     - Send `!` key to tmux (enters bash mode).
     - Wait 1 second.
     - Send the rest of the command + Enter.
     - Launch capture goroutine.
   - Capture goroutine:
     - Wait 2 seconds for command to start.
     - Poll up to 30 times (1-second intervals, 30-second max):
       - Capture pane (plain text).
       - Call `ExtractBashOutput()`.
       - First output: send new message.
       - Changed output: edit message in-place.
       - Unchanged: skip.
     - Truncate output to 3800 chars if too long (prepend `"… "`).
   - Track active captures per (user, thread).
   - Cancel capture when new user message arrives in same topic.

3. Use context.Context for cancellation support.

## Acceptance

- `!git status` sends the command and shows output in Telegram.
- Output updates in real-time as it streams (message edits).
- Long output is truncated.
- New message in same topic cancels the capture.
- 30-second timeout stops the capture.

## Phase

4 — Rich Features

## Depends on

- Task 18
