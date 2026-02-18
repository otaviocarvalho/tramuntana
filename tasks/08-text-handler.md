# Task 08 — Text Handler

## Goal

Implement the text message handler — forward user text to the Claude process in the
bound tmux window.

## Reference

- CCBot: `src/ccbot/bot.py:text_handler` — resolve window, handle unbound, send keys.

## Steps

1. Create `internal/bot/handlers.go`.
2. Implement `handleTextMessage(msg *tgbotapi.Message)`:
   - Extract `userID`, `threadID` (message_thread_id for forum topics), `chatID`.
   - Store group chat ID: `state.SetGroupChatID(userID, threadID, chatID)`.
   - Look up window binding: `state.GetWindowForThread(userID, threadID)`.
   - **If bound**: send text to tmux window via `tmux.SendKeysWithDelay()` (500ms delay before Enter, matching CCBot).
   - **If unbound**: trigger window picker or directory browser (placeholder for now — just reply "topic not bound").
3. Handle the send-keys strategy from PLAN.md:
   - Short messages: direct `tmux send-keys` with literal text + Enter.
   - The 500ms delay between text and Enter avoids TUI misinterpretation.
4. Handle `!` prefix detection:
   - If text starts with `!` and len > 1, send `!` first, wait 1s, then send the rest.
   - Placeholder for bash capture (Task 23).
5. Save state after group_chat_id updates.

## Acceptance

- Text in a bound topic is forwarded to the correct tmux window.
- Unbound topics show a "not bound" message (picker comes in Task 09-10).
- 500ms delay is applied between text and Enter.
- `!` commands send the `!` separately with 1s delay.
- Group chat IDs are persisted.

## Phase

2 — Core Bot

## Depends on

- Task 07
