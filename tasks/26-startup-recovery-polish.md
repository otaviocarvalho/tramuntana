# Task 26 — Startup Recovery & Polish

## Goal

Implement startup state reconciliation, stale binding cleanup, multi-user
verification, and final error handling polish.

## Reference

- CCBot: `src/ccbot/session.py` — `resolve_stale_ids()`, window name → ID migration.
- PLAN.md: Startup Recovery section.

## Steps

1. Implement startup recovery in state management:
   - On boot, call `reconcileState()`:
     a. List live tmux windows: build `{windowName → windowID}` and `{windowID → exists}` maps.
     b. For each persisted window_id in `WindowStates`:
        - If alive in tmux: keep as-is.
        - If dead: try re-resolve by matching `WindowDisplayNames[id]` against live window names.
        - If re-resolved: update all references (bindings, offsets, etc.) to new window_id.
        - If unresolvable: drop the window state entry.
     c. Clean up stale thread bindings pointing to dropped windows.
     d. Clean up stale user_window_offsets for dropped windows.
     e. Clean up stale project_bindings for threads with no binding.
     f. Clean up stale session_map entries for dead windows.
   - Save state after reconciliation.

2. Multi-user verification:
   - Verify `FindUsersForWindow()` returns all users bound to a window.
   - Monitor routes messages to ALL users observing a window, not just the first.
   - Message queue is per-user — each user gets independent delivery.
   - Group chat IDs are correctly resolved for each user.

3. Error handling polish:
   - Stale windows: if `SendKeys` fails because window died, unbind and notify user.
   - Telegram API errors: log but don't crash. Retry transient errors.
   - JSONL parse errors: log and skip the line, continue processing.
   - tmux not running: clear error message on startup.
   - Session map file locked: retry with backoff.

4. Graceful shutdown:
   - On SIGINT/SIGTERM: cancel context, drain message queues, save all state.
   - Wait for in-flight Telegram API calls to complete (with timeout).

5. Wire everything together in `serve` command:
   - Start bot, monitor, status polling as goroutines.
   - All share the same context for cancellation.
   - Log startup: bot username, tmux session, number of recovered bindings.

## Acceptance

- Bot recovers correctly after restart (bindings survive, no duplicate messages).
- Dead windows are detected and cleaned up.
- Multiple users see messages from the same window.
- Graceful shutdown saves all state.
- No panics or unhandled errors during normal operation.

## Phase

6 — Polish

## Depends on

- Task 25
