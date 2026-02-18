# Task 17 — Session Monitor

## Goal

Implement `internal/monitor/monitor.go` — the poll loop that reads Claude Code
JSONL transcript files and routes parsed entries to the message queue.

## Reference

- CCBot: `src/ccbot/session_monitor.py` — `_monitor_loop()`, `scan_projects()`,
  `_read_new_lines()`, `_detect_and_cleanup_changes()`, mtime optimization.

## Steps

1. Create `internal/monitor/monitor.go`:
   ```go
   type Monitor struct {
       config       *config.Config
       state        *state.State
       monitorState *state.MonitorState
       queue        *queue.Queue
       pendingTools map[string]PendingTool
       fileMtimes   map[string]float64
       lastSessionMap map[string]state.SessionMapEntry
       pollInterval time.Duration
   }
   ```
2. Implement `Run(ctx context.Context)` — the main poll loop:
   - Every `pollInterval` seconds:
     a. Load `session_map.json`.
     b. Detect changes vs last session map (new/removed/changed sessions).
     c. Clean up stale tracked sessions.
     d. For each active session: check for new JSONL content.
3. Implement JSONL file discovery — replicate CCBot's full scanning:
   - Scan `~/.claude/projects/` for directories.
   - Read `sessions-index.json` in each project dir to get session_id → file path.
   - Match against active tmux window cwds (from session_map).
   - Fall back to globbing `*.jsonl` files in matching project dirs.
4. Implement incremental reading:
   - Open JSONL file, seek to `lastByteOffset`.
   - Read new lines, parse each via `transcript.ParseLine()`.
   - Advance offset only past successfully parsed lines (partial line = retry next cycle).
   - Detect file truncation (offset > file size → reset to 0, for `/clear`).
5. Implement mtime optimization:
   - Cache file mtime per session_id.
   - Skip reading if mtime hasn't changed.
6. Route parsed entries:
   - For each entry, call `state.FindUsersForWindow(windowID)`.
   - For each user: format entry via render package, enqueue via queue package.
   - Pass pending tools map to transcript parser for cross-cycle tool pairing.
7. Save monitor state periodically (`monitorState.SaveIfDirty()`).

## Acceptance

- Monitor polls JSONL files at the configured interval.
- New content is parsed and routed to correct Telegram topics.
- Byte offsets are tracked — no duplicate messages after restart.
- File truncation (after `/clear`) resets tracking.
- Mtime optimization avoids unnecessary reads.

## Phase

3 — Session Monitor

## Depends on

- Task 16
- Task 14
