# Tramuntana — Development Guide

## Project

Tramuntana is a Go application that bridges Telegram group topics to Claude Code sessions
via tmux, with first-class Minuano task coordination integration. Spiritual successor to
CCBot, rewritten from scratch in Go. See `PLAN.md` for the full specification.

## Architecture

- **Language**: Go
- **CLI framework**: cobra (`github.com/spf13/cobra`)
- **Telegram**: go-telegram-bot-api/v5 (`github.com/go-telegram-bot-api/telegram-bot-api/v5`)
- **Tmux**: exec `tmux` binary directly via `os/exec`
- **Screenshots**: ANSI → PNG via `golang.org/x/image` + `github.com/golang/freetype`
- **Fonts**: JetBrains Mono + Noto CJK + Symbola embedded via `//go:embed`
- **MarkdownV2**: Custom converter (ported from CCBot approach)
- **Auth**: `ALLOWED_USERS` + `ALLOWED_GROUPS` (comma-separated IDs in env)

## Core Invariant

```
1 Telegram Topic = 1 tmux Window = 1 Claude Code process
```

## Repository Layout

```
cmd/tramuntana/main.go          Entry point, cobra root command
internal/bot/                    Telegram bot, handlers, commands, UI
internal/tmux/                   Tmux session/window management
internal/monitor/                JSONL session monitor, transcript parser, terminal parser
internal/state/                  State files (state.json, session_map.json, monitor_state.json)
internal/queue/                  Per-user message queue, flood control
internal/render/                 Markdown conversion, tool formatting, screenshot rendering
internal/minuano/                Minuano CLI bridge, prompt generation
internal/config/                 Environment config loading
hook/                            Claude Code SessionStart hook
tasks/                           Ordered implementation tasks
```

## Reference Implementation

CCBot (Python predecessor) is at `/home/otavio/code/ccbot` on branch `fix/group-chat-id-routing`.
Key reference files:

| Tramuntana package | CCBot reference |
|-------------------|-----------------|
| `internal/bot/` | `src/ccbot/bot.py`, `src/ccbot/handlers/` |
| `internal/monitor/` | `src/ccbot/session_monitor.py`, `src/ccbot/transcript_parser.py`, `src/ccbot/terminal_parser.py` |
| `internal/state/` | `src/ccbot/session.py`, `src/ccbot/monitor_state.py`, `src/ccbot/hook.py` |
| `internal/queue/` | `src/ccbot/handlers/message_queue.py`, `src/ccbot/handlers/message_sender.py` |
| `internal/render/` | `src/ccbot/screenshot.py`, `src/ccbot/markdown_v2.py`, `src/ccbot/handlers/response_builder.py` |
| `internal/config/` | `src/ccbot/config.py` |
| `hook/` | `src/ccbot/hook.py` |

## Task Execution Protocol

Implementation follows the ordered tasks in `tasks/`. **Always work sequentially.**

### How to proceed

1. **Read the current task file** from `tasks/` in numeric order (01, 02, 03, ...).
2. **Check "Depends on"** — if a dependency task is not yet complete, do that one first.
3. **Read PLAN.md** for the detailed specification of whatever you're building. The task
   file tells you *what* to build; PLAN.md tells you *how*.
4. **Read the CCBot reference files** listed in the task for implementation details.
   The Go version should follow Go idioms (goroutines, channels, context) but replicate
   CCBot's behavior and edge-case handling.
5. **Implement** the task following its Steps and Acceptance criteria.
6. **Verify** the Acceptance criteria are met before moving on.
7. **Move to the next task** in numeric order.

### Task phases (in order)

| Phase | Tasks | What it covers |
|-------|-------|----------------|
| 1 — Foundation | 01–05 | Project init, config, tmux, state management, hook |
| 2 — Core Bot | 06–12 | CLI, bot setup, text handler, directory browser, window picker, commands, topic close |
| 3 — Session Monitor | 13–18 | Transcript parser, formatting, MarkdownV2, message queue, monitor loop, status polling |
| 4 — Rich Features | 19–23 | Interactive UI, screenshots, history, bash capture |
| 5 — Minuano Integration | 24–25 | Bridge, commands (/project, /tasks, /pick, /auto, /batch) |
| 6 — Polish | 26 | Startup recovery, multi-user, error handling, graceful shutdown |
| 7 — Cross-Repo Integration | I-01–I-11 | See `tasks/integration/` and `INTEGRATION.md` |

### Current progress

Track which task you're on. When you complete a task, note it here:

- [ ] 01 — Project Initialization
- [ ] 02 — Configuration
- [ ] 03 — Tmux Package
- [ ] 04 — State Management
- [ ] 05 — Hook
- [ ] 06 — CLI Skeleton
- [ ] 07 — Bot Setup
- [ ] 08 — Text Handler
- [ ] 09 — Directory Browser
- [ ] 10 — Window Picker
- [ ] 11 — Claude Commands
- [ ] 12 — Topic Close
- [ ] 13 — JSONL Transcript Parser
- [ ] 14 — Tool Result Formatting
- [ ] 15 — Markdown to MarkdownV2
- [ ] 16 — Message Queue
- [ ] 17 — Session Monitor
- [ ] 18 — Status Polling
- [ ] 19 — Interactive UI
- [ ] 20 — Screenshot Rendering
- [ ] 21 — Screenshot Command
- [ ] 22 — History Command
- [ ] 23 — Bash Capture
- [ ] 24 — Minuano Bridge
- [ ] 25 — Minuano Commands
- [ ] 26 — Startup Recovery & Polish

**Integration tasks** (see `tasks/integration/` and `INTEGRATION.md`):

Minuano side (in `/home/otavio/code/minuano`):
- [ ] I-01 — `--json` flag on `minuano status` and `minuano show`
- [ ] I-02 — `ClaimByID` query in queries.go
- [ ] I-03 — Project filter on `AtomicClaim`
- [ ] I-04 — `minuano-pick` script
- [ ] I-05 — `--project` flag on `minuano-claim`
- [ ] I-06 — `minuano prompt` command

Tramuntana side (depends on Minuano I-01..I-06):
- [ ] I-07 — Update bridge to use `--json`
- [ ] I-08 — Update commands to use `minuano prompt`
- [ ] I-09 — E2E test: pick mode
- [ ] I-10 — E2E test: auto mode
- [ ] I-11 — E2E test: batch mode

## Conventions

- **No ORM, no database.** All state in JSON files with atomic writes.
- **No extra dependencies** beyond those listed in PLAN.md without a clear reason.
- **Go idioms.** Goroutines for concurrency, channels for message routing, context for cancellation.
- **Error handling** — return errors to the caller, print user-friendly messages at the command level.
- **tmux interaction** — always via `os/exec`, never via a Go binding library.
- **Test with** `go build ./cmd/tramuntana` as the minimum bar. Run the binary manually against a real Telegram bot + tmux session for integration testing.

## Building & Running

```bash
go build ./cmd/tramuntana          # build the binary
./tramuntana serve                 # start the Telegram bot
./tramuntana hook --install        # install Claude Code hook
./tramuntana version               # print version
```
