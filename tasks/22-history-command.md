# Task 22 — History Command

## Goal

Implement `/history` — paginated display of session JSONL transcript content.

## Reference

- CCBot: `src/ccbot/handlers/history.py` — paginated history, inline keyboard navigation.

## Steps

1. Create `internal/bot/history.go`.
2. Implement `/history` handler:
   - Resolve window for user's thread.
   - Find the JSONL file for the window's session (from monitor state or session_map).
   - Read and parse the full JSONL file using transcript parser.
   - Format entries into displayable text (abbreviated — show type, tool names,
     first line of text, timestamps).
3. Implement pagination:
   - `ENTRIES_PER_PAGE = 10` (or similar).
   - Show page N of M.
   - Inline keyboard: ◀ Previous | `page/total` | Next ▶.
   - Callback data: `hist_page:N`.
4. Handle pagination callbacks:
   - Re-read JSONL (or cache), format requested page, edit message.
5. Keep entries concise:
   - `assistant` text: first 100 chars + `...`.
   - `tool_use`: `Tool: Read(file.py)`.
   - `tool_result`: `Result: 42 lines`.
   - `user` text: first 100 chars.

## Acceptance

- `/history` shows paginated session transcript.
- Pagination keyboard navigates between pages.
- Entries are concise and readable.
- Works correctly even for long sessions.

## Phase

4 — Rich Features

## Depends on

- Task 17
