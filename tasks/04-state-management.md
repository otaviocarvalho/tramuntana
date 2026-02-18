# Task 04 — State Management

## Goal

Implement `internal/state/` — all three state files, atomic writes, and file locking.

## Reference

- CCBot: `src/ccbot/session.py` (SessionManager), `src/ccbot/monitor_state.py` (MonitorState),
  `src/ccbot/utils.py` (atomic_write_json), `src/ccbot/hook.py` (session_map flock).

## Steps

1. Create `internal/state/state.go` — main state manager:
   ```go
   type State struct {
       ThreadBindings     map[string]map[string]string  // user_id → thread_id → window_id
       WindowStates       map[string]WindowState        // window_id → {session_id, cwd, window_name}
       WindowDisplayNames map[string]string             // window_id → display_name
       UserWindowOffsets  map[string]map[string]int64   // user_id → window_id → byte_offset
       GroupChatIDs       map[string]int64              // "user_id:thread_id" → group_chat_id
       ProjectBindings    map[string]string             // thread_id → minuano_project_id
   }
   ```
   - `Load(path string) (*State, error)` — read from `state.json`, return empty state if not found.
   - `Save(path string) error` — atomic write.
   - `BindThread(userID, threadID, windowID string)`
   - `UnbindThread(userID, threadID string)`
   - `GetWindowForThread(userID, threadID string) (string, bool)`
   - `FindUsersForWindow(windowID string) []UserThread` — return all (userID, threadID) pairs bound to this window.
   - `SetWindowState(windowID string, ws WindowState)`
   - `SetGroupChatID(userID, threadID string, chatID int64)`
   - `GetGroupChatID(userID, threadID string) (int64, bool)`
   - `BindProject(threadID, projectID string)`
   - `GetProject(threadID string) (string, bool)`

2. Create `internal/state/session_map.go`:
   ```go
   type SessionMapEntry struct {
       SessionID  string `json:"session_id"`
       CWD        string `json:"cwd"`
       WindowName string `json:"window_name"`
   }
   ```
   - `LoadSessionMap(path string) (map[string]SessionMapEntry, error)` — read `session_map.json`.
   - `WriteSessionMap(path string, data map[string]SessionMapEntry) error` — file-locked write using `syscall.Flock`.

3. Create `internal/state/monitor_state.go`:
   ```go
   type TrackedSession struct {
       SessionID      string `json:"session_id"`
       FilePath       string `json:"file_path"`
       LastByteOffset int64  `json:"last_byte_offset"`
   }
   type MonitorState struct {
       TrackedSessions map[string]TrackedSession `json:"tracked_sessions"`
       dirty           bool
   }
   ```
   - `Load`, `SaveIfDirty`, `UpdateOffset`, `RemoveSession`.

4. Implement `atomicWriteJSON(path string, data any) error` as a shared helper:
   - Write to temp file in same directory (`.filename.RANDOM.tmp`).
   - `f.Sync()` then `os.Rename()`.

## Acceptance

- State loads from JSON, saves atomically.
- Session map uses file locking (flock) for concurrent hook writes.
- MonitorState tracks byte offsets and only saves when dirty.
- Empty/missing files return zero-value state (no error).

## Phase

1 — Foundation

## Depends on

- Task 02
