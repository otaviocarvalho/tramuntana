# Task 16 — Message Queue

## Goal

Implement `internal/queue/queue.go` and `internal/queue/flood.go` — per-user
message queue with goroutine workers, message merging, tool message editing,
status conversion, and flood control.

## Reference

- CCBot: `src/ccbot/handlers/message_queue.py` — per-user asyncio queues,
  `MessageTask`, merging, tool_msg_ids, status conversion, flood control.
- CCBot: `src/ccbot/handlers/message_sender.py` — safe_reply, safe_send, send_with_fallback.

## Steps

1. Create `internal/queue/queue.go`:
   ```go
   type MessageTask struct {
       UserID      int64
       ThreadID    int64
       ChatID      int64
       Parts       []string
       ContentType string  // "content", "tool_use", "tool_result", "status_update", "status_clear"
       ToolUseID   string  // for tool_result editing
       WindowID    string
   }
   type Queue struct {
       queues map[int64]chan MessageTask  // user_id → channel
       // ... tracking maps
   }
   ```
2. Implement per-user worker goroutines:
   - `Enqueue(task MessageTask)` — create worker goroutine on first message for user.
   - Worker reads from channel, processes tasks sequentially.
3. Implement message merging:
   - Before processing a content task, peek at channel for consecutive mergeable tasks.
   - Merge if: same window_id, both "content" type, neither is tool_use/tool_result,
     combined length < 3800 chars.
4. Implement tool_use/tool_result editing:
   - Track `toolMsgIDs map[string]int` — tool_use_id → Telegram message_id.
   - When `tool_use` is sent, store its message_id.
   - When `tool_result` arrives with matching tool_use_id, edit the message in-place.
5. Implement status conversion:
   - Track `statusMsgInfo map[userThread]StatusInfo` — (user_id, thread_id) → (msg_id, window_id, text).
   - When first content arrives and a status message exists, edit it in-place
     instead of deleting + sending new.
6. Create `internal/queue/flood.go`:
   - Track `floodUntil map[int64]time.Time` — user_id → flood ban expiry.
   - On Telegram 429 with retry > 10s: set flood ban.
   - During flood: drop status tasks, delay content tasks (max 10s wait).
   - Clear flood ban on expiry.
7. Sending: use MarkdownV2 with plain text fallback (from render package).

## Acceptance

- Each user gets a dedicated worker goroutine.
- Consecutive text messages are merged up to 3800 chars.
- tool_result edits the tool_use message in-place.
- Status messages are converted to content on first real message.
- Flood control drops status and delays content during 429.

## Phase

3 — Session Monitor

## Depends on

- Task 15
