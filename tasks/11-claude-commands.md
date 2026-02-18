# Task 11 — Claude Commands

## Goal

Implement bot commands that are forwarded to Claude Code: `/clear`, `/compact`,
`/cost`, `/help`, `/memory`, `/esc`.

## Reference

- CCBot: `src/ccbot/bot.py` — command handlers for forwarded commands,
  `/esc` sends Escape key.

## Steps

1. Create `internal/bot/commands.go`.
2. Implement a generic command forwarder:
   - Resolve window for the user's thread.
   - If not bound, reply with error message.
   - Send the command text (e.g. `/compact`) as keys to the tmux window.
3. Implement specific commands:
   - `/clear` — send `/clear` to tmux. Also reset session tracking: clear the
     window's entry from monitor state (byte offset) so the monitor starts fresh
     after Claude clears its session.
   - `/compact` — send `/compact` to tmux.
   - `/cost` — send `/cost` to tmux.
   - `/help` — send `/help` to tmux.
   - `/memory` — send `/memory` to tmux.
   - `/esc` — send Escape key to tmux (`tmux send-keys Escape`). Do NOT send literal text.
4. Register all command handlers in bot setup.

## Acceptance

- Each command forwards correctly to the bound tmux window.
- `/esc` sends the Escape key (not literal text "/esc").
- `/clear` additionally resets session monitoring state.
- Commands on unbound topics reply with an error message.

## Phase

2 — Core Bot

## Depends on

- Task 08
