# Task 12 — Topic Close

## Goal

Detect topic close events and clean up: kill the tmux window, unbind the thread,
remove state entries.

## Reference

- CCBot: `src/ccbot/bot.py` — `FORUM_TOPIC_CLOSED` filter, cleanup handler.
- CCBot: `src/ccbot/handlers/cleanup.py` — `cleanup_thread_state()`.

## Steps

1. Add topic close detection in the message handler:
   - In `go-telegram-bot-api`, forum topic closed events arrive as messages with
     `ForumTopicClosed` field set.
   - Detect this in the message routing logic.
2. Implement `handleTopicClose(msg *tgbotapi.Message)`:
   - Look up window binding for this thread.
   - If bound:
     - Kill the tmux window: `tmux.KillWindow()`.
     - Unbind the thread: `state.UnbindThread()`.
     - Remove window state, display name, offsets.
     - Remove project binding if any.
     - Remove from session_map.json if applicable.
     - Save state.
3. Handle gracefully if window is already dead (tmux kill fails).

## Acceptance

- Closing a Telegram forum topic kills the associated tmux window.
- All state entries for the thread/window are cleaned up.
- Closing an unbound topic is a no-op.
- Already-dead windows don't cause errors.

## Phase

2 — Core Bot

## Depends on

- Task 08
