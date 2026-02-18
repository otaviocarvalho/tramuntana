# Task 15 — Markdown to MarkdownV2

## Goal

Implement `internal/render/markdown.go` — convert standard Markdown to Telegram
MarkdownV2 format.

## Reference

- CCBot: `src/ccbot/markdown_v2.py` — uses `telegramify-markdown` library.
- Telegram MarkdownV2 spec: special chars that must be escaped:
  `_`, `*`, `[`, `]`, `(`, `)`, `~`, `` ` ``, `>`, `#`, `+`, `-`, `=`, `|`, `{`, `}`, `.`, `!`

## Steps

1. Create `internal/render/markdown.go`.
2. Implement `ToMarkdownV2(text string) string`:
   - Escape all MarkdownV2 special characters outside of formatted regions.
   - Preserve code blocks (``` and inline `): content inside code is only escaped for `` ` `` and `\`.
   - Convert bold (`**text**` → `*text*` in MarkdownV2).
   - Convert italic (`_text_` → `_text_` — same but ensure proper escaping).
   - Convert links (`[text](url)` → `[text](url)` — same format but escape inside).
   - Handle nested formatting carefully.
3. Implement `ToPlainText(text string) string`:
   - Strip all markdown formatting, return raw text.
   - Used as fallback when MarkdownV2 send fails.
4. Implement `SendWithFallback` pattern:
   - Try sending as MarkdownV2.
   - On parse error, retry as plain text.
   - This can be a helper function or part of the message sending logic.

## Acceptance

- Standard markdown converts correctly to MarkdownV2.
- Code blocks are preserved without double-escaping.
- Special characters outside formatted regions are escaped.
- Plain text fallback strips all formatting.

## Phase

3 — Session Monitor

## Depends on

- Task 07
