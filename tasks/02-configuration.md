# Task 02 — Configuration

## Goal

Implement `internal/config/config.go` — the config struct, `.env` loading, and defaults.

## Reference

- CCBot: `src/ccbot/config.py` — singleton config, env var loading, path defaults.

## Steps

1. Add dependency: `go get github.com/joho/godotenv`.
2. Create `internal/config/config.go` with a `Config` struct containing all fields from PLAN.md:
   - `TelegramBotToken string` (required)
   - `AllowedUsers []int64` (required, parsed from comma-separated)
   - `AllowedGroups []int64` (optional, parsed from comma-separated)
   - `TramuntanaDir string` (default `~/.tramuntana`)
   - `TmuxSessionName string` (default `tramuntana`)
   - `ClaudeCommand string` (default `claude`)
   - `MonitorPollInterval float64` (default `2.0`)
   - `MinuanoBin string` (default `minuano`)
   - `MinuanoDB string` (optional)
3. Implement `Load() (*Config, error)`:
   - Call `godotenv.Load()` (ignore error if `.env` missing).
   - Read each env var with `os.Getenv`.
   - Validate required fields, return error if missing.
   - Parse `ALLOWED_USERS` and `ALLOWED_GROUPS` as comma-separated int64 slices.
   - Expand `~` in `TramuntanaDir` to home directory.
   - Ensure `TramuntanaDir` exists (create if needed).
4. Implement `IsAllowedUser(userID int64) bool` on Config.
5. Implement `IsAllowedGroup(groupID int64) bool` on Config.

## Acceptance

- `Load()` reads from env vars and `.env` file.
- Missing `TELEGRAM_BOT_TOKEN` or `ALLOWED_USERS` returns an error.
- Default values are applied for optional fields.
- `TramuntanaDir` directory is created if it doesn't exist.

## Phase

1 — Foundation

## Depends on

- Task 01
