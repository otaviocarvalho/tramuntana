# Task 10 — Window Picker

## Goal

Implement the inline keyboard for binding a topic to an existing unbound tmux window.

## Reference

- CCBot: `src/ccbot/handlers/directory_browser.py` — `build_window_picker()`,
  `CB_WIN_BIND` callbacks, `UNBOUND_WINDOWS_KEY`.

## Steps

1. Add to `internal/bot/window_picker.go`.
2. Implement `buildWindowPicker(windows []tmux.Window) tgbotapi.InlineKeyboardMarkup`:
   - List unbound windows with their cwds (replace home dir with `~`).
   - Two windows per row (up to the number of windows).
   - Use numeric index in callback data (`win_bind:N`).
   - Cache window list in per-user state (keyed by user_id).
   - Bottom row: `New Session` (switches to directory browser), `Cancel`.
3. Handle callback queries:
   - `win_bind:N` — bind topic to Nth window from cached list.
   - `win_new` — switch to directory browser.
   - `win_cancel` — dismiss picker.
4. On bind:
   - `state.BindThread(userID, threadID, windowID)`.
   - Rename Telegram topic to match window display name.
   - Send pending text to the window.
5. Update the unbound-topic flow in `handleTextMessage()`:
   - Call `tmux.ListWindows()`, filter out already-bound windows.
   - If unbound windows exist → show window picker.
   - If no unbound windows → show directory browser.

## Acceptance

- Unbound topic with available windows shows window picker.
- Selecting a window binds the topic and forwards pending text.
- "New Session" switches to directory browser.
- Topic is renamed to match the window name.

## Phase

2 — Core Bot

## Depends on

- Task 09
