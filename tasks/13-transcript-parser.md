# Task 13 — JSONL Transcript Parser

## Goal

Implement `internal/monitor/transcript.go` — parse Claude Code JSONL transcript
entries, extract content blocks, and pair tool_use with tool_result.

## Reference

- CCBot: `src/ccbot/transcript_parser.py` — `TranscriptParser`, `parse_line()`,
  content block extraction, tool pairing.

## Steps

1. Create `internal/monitor/transcript.go`.
2. Define entry types:
   ```go
   type Entry struct {
       Type       string         // "user", "assistant", "summary"
       Message    json.RawMessage
       Blocks     []ContentBlock
   }
   type ContentBlock struct {
       Type      string // "text", "tool_use", "tool_result", "thinking"
       Text      string // for text/thinking blocks
       ToolName  string // for tool_use
       ToolInput string // for tool_use (first arg or summary)
       ToolUseID string // for tool_use and tool_result
       Content   string // for tool_result
       IsError   bool   // for tool_result
   }
   ```
3. Implement `ParseLine(line []byte) (*Entry, error)`:
   - Parse JSON: expect `{"type": "...", "message": {"content": [...]}}`.
   - Handle entry types: `user`, `assistant`, `summary`. Ignore others.
   - For `assistant` and `user`: extract content blocks from `message.content` array.
   - Content block types: `text`, `tool_use`, `tool_result`, `thinking`.
4. Implement tool_use input extraction:
   - `Read` → file path.
   - `Bash` → command (first 100 chars).
   - `Write` → file path.
   - `Edit` → file path.
   - `Grep`/`Glob` → pattern.
   - `Task` → description.
   - `WebFetch` → URL.
   - `WebSearch` → query.
5. Implement tool pairing:
   ```go
   type PendingTool struct {
       ToolUseID string
       ToolName  string
       Input     string
       WindowID  string
   }
   ```
   - Maintain `map[string]PendingTool` keyed by tool_use_id.
   - On `tool_use`: store in pending map.
   - On `tool_result`: look up by tool_use_id, return paired result.
   - Pending tools persist across parse calls (carried by caller).

## Acceptance

- Parses all JSONL entry types from Claude Code sessions.
- Extracts text, tool_use, tool_result, and thinking blocks.
- Tool_use input is extracted per tool type.
- Tool pairing works within and across parse calls.

## Phase

3 — Session Monitor

## Depends on

- Task 04
