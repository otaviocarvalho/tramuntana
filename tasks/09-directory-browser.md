# Task 09 — Directory Browser

## Goal

Implement the inline keyboard directory browser for creating new Claude sessions
from a chosen filesystem directory.

## Reference

- CCBot: `src/ccbot/handlers/directory_browser.py` — `build_directory_browser()`,
  navigation callbacks, session creation on confirm.

## Steps

1. Create `internal/bot/directory_browser.go`.
2. Implement directory browser state (per-user, in-memory):
   ```go
   type BrowseState struct {
       CurrentPath string
       Page        int
       Dirs        []string  // cached subdirectory names
       PendingText string    // text to send after binding
   }
   ```
3. Implement `buildDirectoryBrowser(path string, page int) tgbotapi.InlineKeyboardMarkup`:
   - List non-hidden subdirectories, sorted.
   - `DIRS_PER_PAGE = 6`, two dirs per row.
   - Truncate display names to 13 chars with `…`.
   - Use numeric index in callback data (avoid 64-byte limit).
   - Pagination row: ◀ `page/total` ▶.
   - Action row: `..` (parent dir), `Select` (confirm), `Cancel`.
4. Handle callback queries:
   - `dir_select:N` — navigate into Nth subdirectory.
   - `dir_page:N` — change page.
   - `dir_up` — navigate to parent directory.
   - `dir_confirm` — create new tmux window at current path.
   - `dir_cancel` — dismiss browser.
5. On confirm:
   - `tmux.NewWindow()` with selected path and Claude command.
   - Wait up to 5s for session_map entry (poll `session_map.json`).
   - Bind thread → window.
   - Rename Telegram topic to match window name.
   - Send pending text to the new window.

## Acceptance

- Unbound topic with no available windows shows directory browser.
- Navigation works: enter subdirectories, go up, paginate.
- Confirming creates a new Claude session in the selected directory.
- Topic is renamed and pending text is forwarded.

## Phase

2 — Core Bot

## Depends on

- Task 08
