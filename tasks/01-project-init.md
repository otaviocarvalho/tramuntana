# Task 01 — Project Initialization

## Goal

Bootstrap the Go module and create the full directory skeleton from PLAN.md.

## Steps

1. Run `go mod init github.com/otaviocarvalho/tramuntana`.
2. Create every directory in the repository layout:
   - `cmd/tramuntana/`
   - `internal/bot/`
   - `internal/tmux/`
   - `internal/monitor/`
   - `internal/state/`
   - `internal/queue/`
   - `internal/render/`
   - `internal/minuano/`
   - `internal/config/`
   - `hook/`
   - `tasks/` (already exists)
3. Create placeholder `main.go` at `cmd/tramuntana/main.go` with `package main` and an empty `func main()`.
4. Create `.env.example` with all env vars from PLAN.md (Configuration section):
   - `TELEGRAM_BOT_TOKEN`, `ALLOWED_USERS`, `ALLOWED_GROUPS`, `TRAMUNTANA_DIR`,
     `TMUX_SESSION_NAME`, `CLAUDE_COMMAND`, `MONITOR_POLL_INTERVAL`, `MINUANO_BIN`, `MINUANO_DB`.
5. Create `.gitignore` with Go defaults: binary, `.env`, `vendor/`, `*.exe`, `*.test`, `*.out`.

## Acceptance

- `go build ./cmd/tramuntana` succeeds (even if the binary does nothing).
- Directory structure matches PLAN.md layout.
- `.env` is gitignored.

## Phase

1 — Foundation
