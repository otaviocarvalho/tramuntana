# Task 14 — Tool Result Formatting

## Goal

Implement `internal/render/format.go` — format tool results, expandable quotes,
and thinking blocks for Telegram display.

## Reference

- CCBot: `src/ccbot/handlers/response_builder.py` — tool result formatting,
  expandable quotes, content splitting.
- PLAN.md: Tool Result Formatting table.

## Steps

1. Create `internal/render/format.go`.
2. Implement `FormatToolUse(name, input string) string`:
   - Return summary line: `**Read**(file.py)`, `**Bash**(git status)`, etc.
3. Implement `FormatToolResult(toolName, content string, isError bool) string`:
   - Per-tool formatting from PLAN.md:
     - `Read` → `"Read N lines"` (count lines, no content shown).
     - `Write` → `"Wrote N lines"` (count lines, no content shown).
     - `Bash` → `"Output N lines"` + expandable quote.
     - `Grep` → `"Found N matches"` + expandable quote.
     - `Glob` → `"Found N files"` + expandable quote.
     - `Edit` → `"Added X, removed Y"` + expandable quote (unified diff).
     - `Task` → `"Agent output N lines"` + expandable quote.
     - `WebFetch` → `"Fetched N characters"` + expandable quote.
     - `WebSearch` → `"N search results"` + expandable quote.
     - Error → first line (100 chars), expandable if multiline.
4. Implement expandable quote helper:
   - Telegram expandable quotes use `**>||` spoiler syntax or blockquote.
   - Truncate long content to fit Telegram's 4096 char message limit.
5. Implement `FormatThinking(text string) string`:
   - Truncate to 500 chars.
   - Wrap in expandable quote.
6. Implement `FormatText(text string) string`:
   - Strip system tags (e.g. `<system-reminder>` blocks).
   - Return clean text.

## Acceptance

- Each tool type produces correctly formatted output.
- Expandable quotes work for long content.
- Thinking blocks are truncated and quoted.
- System tags are stripped from text blocks.

## Phase

3 — Session Monitor

## Depends on

- Task 13
