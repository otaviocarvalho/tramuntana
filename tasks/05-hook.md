# Task 05 — Hook

## Goal

Implement the Claude Code SessionStart hook and the `--install` flag to register it.

## Reference

- CCBot: `src/ccbot/hook.py` — reads stdin JSON, gets tmux info, writes session_map.json.

## Steps

1. Create `hook/hook.go` with `Run() error` function:
   - Read JSON from stdin: `{"session_id": "...", "cwd": "...", "hook_event_name": "SessionStart"}`.
   - Validate: `session_id` matches UUID regex, `cwd` is absolute path, `hook_event_name == "SessionStart"`.
   - Read `$TMUX_PANE` env var. If empty, exit silently (not in tmux).
   - Exec `tmux display-message -t $PANE -p "#{session_name}:#{window_id}:#{window_name}"`.
   - Parse output as `session_name:window_id:window_name`.
   - Build key: `"session_name:window_id"` (e.g. `"tramuntana:@12"`).
   - File-locked read-modify-write of `session_map.json` (using state package).
   - Write entry: `{key: {session_id, cwd, window_name}}`.
   - Important: do NOT import config package (hook runs inside Claude's tmux pane where bot env vars are unavailable). Use `$TRAMUNTANA_DIR` or `~/.tramuntana` directly.

2. Implement `Install() error` function:
   - Find `tramuntana` binary path via `os.Executable()`.
   - Read `~/.claude/settings.json` (create if not found).
   - Check if hook already installed (scan for entries containing `"tramuntana hook"`).
   - Add to `hooks.SessionStart[].hooks[]`: `{"type": "command", "command": "/path/to/tramuntana hook", "timeout": 5}`.
   - Write back with atomic write.

## Acceptance

- `echo '{"session_id":"test-uuid","cwd":"/tmp","hook_event_name":"SessionStart"}' | TMUX_PANE=%0 tramuntana hook` writes to session_map.json.
- `tramuntana hook --install` adds the hook entry to `~/.claude/settings.json`.
- Hook does not import config package or require `TELEGRAM_BOT_TOKEN`.
- Duplicate install is a no-op.

## Phase

1 — Foundation

## Depends on

- Task 04
