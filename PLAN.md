# Tramuntana — Design Plan

A Go application that bridges Telegram group topics to Claude Code sessions via tmux,
with first-class Minuano task coordination integration. Spiritual successor to CCBot,
rewritten from scratch in Go with agentic capabilities.

---

## Core Invariant

```
1 Telegram Topic = 1 tmux Window = 1 Claude Code process
```

Optionally, a topic can also be bound to a Minuano project, enabling task-driven work
modes alongside interactive chat.

---

## Principles

- **Single binary.** One `tramuntana` CLI handles bot, hook, and utilities.
- **Minuano is external.** Communication via CLI commands only. No shared database.
- **Stateless restarts.** All durable state in JSON files. Bot recovers from crashes
  by reconciling persisted state against live tmux windows.
- **Go idioms.** Goroutines for concurrency, channels for message routing, context
  for cancellation. No Python async/await port — redesign for Go's model.

---

## Repository Layout

```
tramuntana/
├── cmd/
│   └── tramuntana/
│       └── main.go                # entry point, cobra root command
├── internal/
│   ├── bot/
│   │   ├── bot.go                 # telegram bot setup, handler registration
│   │   ├── handlers.go            # text, command, callback handlers
│   │   ├── commands.go            # /pick, /auto, /batch, /tasks, /project
│   │   ├── directory_browser.go   # inline keyboard directory picker
│   │   ├── window_picker.go       # inline keyboard for unbound windows
│   │   └── interactive.go         # AskUserQuestion/ExitPlanMode UI
│   ├── tmux/
│   │   └── tmux.go                # session/window management (shared patterns with minuano)
│   ├── monitor/
│   │   ├── monitor.go             # JSONL session monitor (poll loop)
│   │   ├── transcript.go          # JSONL parser: entries, tool pairing, content blocks
│   │   └── terminal.go            # terminal parser: status line, interactive detection
│   ├── state/
│   │   ├── state.go               # thread bindings, window states, offsets, project bindings
│   │   ├── session_map.go         # session_map.json read/write (hook writes, monitor reads)
│   │   └── monitor_state.go       # byte offsets per tracked session
│   ├── queue/
│   │   ├── queue.go               # per-user message queue with worker goroutines
│   │   └── flood.go               # telegram rate limit / flood control
│   ├── render/
│   │   ├── markdown.go            # markdown → telegram MarkdownV2 conversion
│   │   ├── screenshot.go          # terminal ANSI → PNG rendering
│   │   └── format.go              # tool result formatting, expandable quotes
│   ├── minuano/
│   │   ├── bridge.go              # exec minuano CLI commands, parse output
│   │   └── prompt.go              # generate prompts for pick/auto/batch modes
│   └── config/
│       └── config.go              # env vars, .env loading, defaults
├── hook/
│   └── hook.go                    # SessionStart hook: stdin JSON → session_map.json
├── .env.example
├── go.mod
├── go.sum
└── PLAN.md                        # this file
```

---

## CLI Commands

```
tramuntana
├── serve                    Run the Telegram bot (main mode)
│   └── --config PATH        Path to .env override
├── hook                     SessionStart hook (called by Claude Code)
│   └── --install            Install hook into ~/.claude/settings.json
└── version                  Print version
```

---

## Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `TELEGRAM_BOT_TOKEN` | yes | — | Bot token from @BotFather |
| `ALLOWED_USERS` | yes | — | Comma-separated Telegram user IDs |
| `TRAMUNTANA_DIR` | no | `~/.tramuntana` | Config/state directory |
| `TMUX_SESSION_NAME` | no | `tramuntana` | Tmux session name |
| `CLAUDE_COMMAND` | no | `claude` | Command to start Claude Code in new windows |
| `MONITOR_POLL_INTERVAL` | no | `2.0` | Seconds between JSONL poll cycles |
| `MINUANO_BIN` | no | `minuano` | Path to minuano binary |
| `MINUANO_DB` | no | — | DATABASE_URL for minuano (passed via --db flag) |

### State Files (in `$TRAMUNTANA_DIR`)

| File | Written by | Read by | Content |
|------|-----------|---------|---------|
| `state.json` | bot | bot | Thread bindings, window states, offsets, project bindings, group chat IDs |
| `session_map.json` | hook | monitor | `{session:window_id → {session_id, cwd, window_name}}` |
| `monitor_state.json` | monitor | monitor | `{session_id → {file_path, last_byte_offset}}` |

All writes use atomic temp-file + rename pattern.

---

## Telegram Bot Features

### Replicated from CCBot

| Feature | Description |
|---------|-------------|
| **Topic binding** | Each forum topic binds to a tmux window. Unbound topics trigger window picker or directory browser. |
| **Message forwarding** | User text → `tmux send-keys` to Claude's window. |
| **Session monitoring** | Poll Claude Code JSONL files, parse entries, route to correct topic. |
| **Tool result formatting** | Bash output, Read/Write summaries, Edit diffs, Grep/Glob counts, expandable quotes. |
| **Status line** | Detect Claude's spinner/status from terminal, show as editable message in topic. |
| **Interactive UI** | Detect AskUserQuestion/ExitPlanMode prompts, render as inline keyboard with navigation. |
| **Screenshots** | `/screenshot` — capture pane with ANSI colors, render to PNG, send with control keyboard. |
| **History** | `/history` — paginated display of session JSONL content. |
| **Claude commands** | `/clear`, `/compact`, `/cost`, `/help`, `/memory` forwarded to tmux. `/clear` resets session tracking. |
| **`/esc`** | Send Escape to interrupt Claude. |
| **Bash capture** | `!command` messages: capture terminal output for 30s, send as editable message. |
| **Topic close** | Kills tmux window, unbinds thread, cleans up state. |
| **Multi-user** | Multiple users can bind different topics to different windows. All observers notified. |
| **Rate limiting** | Per-user queues, flood control, status deduplication. |
| **Stale recovery** | On startup, reconcile persisted window IDs against live tmux windows. |
| **Directory browser** | Inline keyboard for navigating filesystem to create new sessions. |
| **Window picker** | Inline keyboard for binding to existing unbound windows. |
| **Topic renaming** | Auto-rename topic to match window name on bind. |

### New: Minuano Integration

| Command | Description |
|---------|-------------|
| `/project <name>` | Bind this topic to a Minuano project. Stored in `state.json` as `project_bindings`. |
| `/tasks` | Show ready tasks for the topic's bound project. Calls `minuano status --project X`. |
| `/pick <task-id>` | Single-task mode: claim one task, work it, exit. |
| `/auto` | Autonomous mode: loop claiming from project until queue empty. |
| `/batch <id1> [id2] ...` | Batch mode: work through specified tasks in order. |

---

## Minuano Bridge (`internal/minuano/`)

Tramuntana communicates with Minuano exclusively via its CLI. No shared DB.

### bridge.go — Command Execution

```go
type Bridge struct {
    Bin    string // path to minuano binary
    DBFlag string // optional --db flag value
}

func (b *Bridge) Status(project string) ([]Task, error)        // minuano status --project X (parse table output)
func (b *Bridge) Show(taskID string) (*TaskDetail, error)       // minuano show <id> (parse output)
func (b *Bridge) Tree(project string) (string, error)           // minuano tree --project X (raw text)
func (b *Bridge) Prompt(mode string, args ...string) (string, error) // minuano prompt <mode> <args...>
```

All functions exec the binary, capture stdout, parse the output.
Errors from minuano stderr are wrapped and returned.

### prompt.go — Prompt Generation

For each mode, generate a self-contained prompt that gets sent to the Claude
process in the topic's tmux window.

**Single mode** (`/pick <task-id>`):
1. Call `minuano show <task-id>` to get full task spec + context.
2. Build prompt: task spec + instructions (work it, call `minuano-done`, exit).
3. Write prompt to temp file.
4. Send to tmux: the prompt text directly (Claude reads it as user message).

**Auto mode** (`/auto`):
1. Call `minuano status --project X` to verify tasks exist.
2. Build prompt: loop instructions (claim from project, work, done, repeat, exit on empty).
3. Include project scope and env setup.
4. Send to tmux as user message.

**Batch mode** (`/batch <id1> <id2> ...`):
1. Call `minuano show` for each task ID to get specs.
2. Build prompt: ordered task list with specs, instructions to work each in order.
3. Send to tmux as user message.

### Environment Setup

When Tramuntana creates a tmux window (or sends a Minuano command), it ensures:
- `MINUANO_BIN` is on PATH
- `DATABASE_URL` is exported (from `MINUANO_DB` config)
- Minuano scripts directory is on PATH
- `AGENT_ID` is set (using topic slug or window ID)

This is done via `tmux send-keys "export ..."` before the Claude process starts.

---

## Message Flow

### User → Claude (text message)

```
Telegram topic
  → bot.handlers.textHandler
    → resolve window for (user_id, thread_id)
    → [unbound?] show window picker or directory browser
    → [bound?] tmux.SendKeys(window_id, text)
    → start status polling for this window
```

### Claude → Telegram (session monitor)

```
monitor.Run() goroutine (poll loop every 2s)
  → load session_map.json
  → for each active session:
      → read JSONL file from last byte offset
      → transcript.Parse(lines) → []Entry
      → for each entry:
          → state.FindUsersForSession(session_id) → [(user_id, thread_id)]
          → for each user:
              → render.Format(entry) → text parts
              → queue.Enqueue(user_id, thread_id, parts)

queue.Worker(user_id) goroutine
  → dequeue message task
  → merge consecutive text messages (up to 3800 chars)
  → try edit tool_use message for tool_result
  → convert status message to content (reduce message count)
  → send via Telegram API (MarkdownV2, fallback to plain)
  → update status message from terminal
```

### Minuano commands (/pick, /auto, /batch)

```
Telegram topic
  → bot.commands.pickHandler (or autoHandler, batchHandler)
    → minuano.Bridge.Prompt("single", taskID)  → prompt text
    → tmux.SendKeys(window_id, prompt)
    → reply "Working on task <id>..."
    → [monitor picks up Claude's output normally]
```

---

## JSONL Transcript Parsing (`internal/monitor/transcript.go`)

### Entry Types

| JSONL type | Handling |
|-----------|---------|
| `user` | Parse content blocks. Detect tool_result (pair with pending tool_use). Detect local commands. |
| `assistant` | Parse content blocks. Extract text, tool_use, thinking. |
| `summary` | Extract display name. |
| other | Ignore. |

### Content Block Types

| Block type | Formatting |
|-----------|-----------|
| `text` | Direct text. Strip system tags. Detect local command XML. |
| `tool_use` | Summary line: `**Read**(file.py)`, `**Bash**(git status)`, etc. Store as pending for pairing. |
| `tool_result` | Pair with pending tool_use by ID. Format per tool (see below). |
| `thinking` | Truncate to 500 chars, wrap in expandable quote. |

### Tool Result Formatting

| Tool | Format |
|------|--------|
| Read | `"Read N lines"` (no content) |
| Write | `"Wrote N lines"` (no content) |
| Bash | `"Output N lines"` + expandable quote |
| Grep | `"Found N matches"` + expandable quote |
| Glob | `"Found N files"` + expandable quote |
| Edit | Unified diff, `"Added X, removed Y"` + expandable quote |
| Task | `"Agent output N lines"` + expandable quote |
| WebFetch | `"Fetched N characters"` + expandable quote |
| WebSearch | `"N search results"` + expandable quote |
| Error | First line (100 chars), expandable if multiline |

### Tool Pairing Across Poll Cycles

The monitor keeps a `map[string]PendingTool` (keyed by tool_use_id) that persists
across poll iterations. A tool_use in one poll cycle may get its tool_result in the
next. Flush unpaired tools only on session end or timeout.

---

## Terminal Parsing (`internal/monitor/terminal.go`)

### Status Line Detection

Find the chrome separator (line of `─` chars, ≥20 wide) in the last 10 lines of
captured pane output. The line above it containing a spinner character
(`·✻✽✶✳✢`) is the status text.

### Interactive UI Detection

Ordered pattern matching on captured pane:

| Pattern | Top marker | Bottom marker |
|---------|-----------|--------------|
| ExitPlanMode | `"Would you like to proceed?"` or `"Claude has written up a plan"` | `"ctrl-g to edit"` or `"Esc to"` |
| AskUserQuestion (multi) | `"← [☐✔☒]"` | last non-empty line |
| AskUserQuestion (single) | `"[☐✔☒]"` | `"Enter to select"` |
| PermissionPrompt | `"Do you want to proceed?"` | `"Esc to cancel"` |
| RestoreCheckpoint | `"Restore the code"` | `"Enter to continue"` |
| Settings | `"Settings:.*tab to cycle"` | `"Esc to cancel"` or `"Type to filter"` |

### Bash Output Extraction

Strip bottom chrome. Search upward for `"! <command>"` echo. Return command + output below.

---

## State Management (`internal/state/`)

### state.json

```json
{
  "thread_bindings": {
    "123456": {"42": "@12"}
  },
  "window_states": {
    "@12": {"session_id": "uuid", "cwd": "/path", "window_name": "project"}
  },
  "window_display_names": {
    "@12": "project"
  },
  "user_window_offsets": {
    "123456": {"@12": 54321}
  },
  "group_chat_ids": {
    "123456:42": -1001234567890
  },
  "project_bindings": {
    "42": "auth-system"
  }
}
```

`project_bindings` maps `thread_id → minuano_project_id`. This is the new
field that CCBot doesn't have.

### Atomic Writes

All JSON state files: write to temp file in same dir → fsync → os.Rename.

### Startup Recovery

On boot, reconcile persisted window IDs against live tmux:
1. List live windows: build `{window_name → window_id}` map.
2. For each persisted window_id: if alive, keep. If dead, try re-resolve by display name. If unresolvable, drop.
3. Clean up stale thread bindings, offsets, project bindings.

---

## Per-User Message Queue (`internal/queue/`)

Each user gets a dedicated goroutine processing a channel:

```go
type MessageTask struct {
    UserID    int64
    ThreadID  int64
    ChatID    int64
    Parts     []string
    ContentType string  // "text", "tool_use", "tool_result", "status"
    ToolUseID   string  // for tool_result editing
    WindowID    string
}
```

### Processing Rules

1. **Merge** consecutive text tasks up to 3800 chars. Tool entries break merge chains.
2. **Edit** tool_use message in-place when tool_result arrives (tracked by tool_use_id → msg_id).
3. **Convert** existing status message to first content message (edit instead of delete+send).
4. **Fallback** from MarkdownV2 to plain text on any send error.
5. **Flood control**: on 429 with retry > 10s, drop status tasks, delay content tasks.

---

## Screenshot Rendering (`internal/render/screenshot.go`)

Capture pane with ANSI (`tmux capture-pane -e -p`) → parse escape sequences → render
to image.

Go options:
- Use `golang.org/x/image/font` + `github.com/golang/freetype` for text rendering.
- Embed a monospace font (JetBrains Mono or similar).
- Parse ANSI SGR sequences for foreground/background colors (16 + 256 + RGB).
- Output PNG via `image/png`.

---

## Hook (`hook/hook.go`)

Called by Claude Code on `SessionStart` events.

1. Read JSON from stdin: `{session_id, cwd, hook_event_name}`.
2. Validate: UUID format, absolute path, event == "SessionStart".
3. Read `$TMUX_PANE` env var.
4. Exec `tmux display-message -t $PANE -p "#{session_name}:#{window_id}:#{window_name}"`.
5. Build key: `"tramuntana:@N"` (session_name:window_id).
6. File-locked write to `session_map.json` (flock).

### Installation (`tramuntana hook --install`)

1. Find `tramuntana` binary path.
2. Read `~/.claude/settings.json`.
3. Add to `hooks.SessionStart[].hooks[]`: `{"type": "command", "command": "/path/to/tramuntana hook", "timeout": 5}`.
4. Write back.

---

## Key Differences from CCBot

| Aspect | CCBot (Python) | Tramuntana (Go) |
|--------|---------------|-----------------|
| Language | Python 3.12, async/await | Go, goroutines/channels |
| Telegram lib | python-telegram-bot v21 | telebot/v4 or go-telegram-bot-api/v5 |
| tmux | libtmux (Python bindings) | exec `tmux` binary directly |
| JSONL parsing | json module + custom parser | encoding/json + custom parser |
| State | dict + json.dump | struct + encoding/json |
| Screenshot | Pillow + freetype | image/png + x/image/font |
| Concurrency | asyncio + per-user queues | goroutines + per-user channels |
| Minuano | — | First-class bridge via CLI |
| Agent modes | — | /pick, /auto, /batch |
| Project binding | — | topic → minuano project |

---

## Go Dependencies

```
github.com/spf13/cobra              CLI framework
github.com/joho/godotenv            .env loading
github.com/go-telegram-bot-api/telegram-bot-api/v5   Telegram Bot API
golang.org/x/image                  Font rendering for screenshots
github.com/golang/freetype          TrueType font rendering
```

Minimal dependency surface. tmux interaction via `os/exec`. JSONL parsing via stdlib.

---

## Implementation Phases

### Phase 1 — Foundation

| Task | Description |
|------|-------------|
| Project init | Go module, directory skeleton, config loading |
| tmux package | Session/window management (reuse patterns from minuano) |
| State management | state.json, session_map.json, atomic writes, startup recovery |
| Hook | SessionStart hook + install command |

### Phase 2 — Core Bot

| Task | Description |
|------|-------------|
| Bot setup | Telegram bot connection, handler registration, authorization |
| Text handler | Forward text to tmux window, handle unbound topics |
| Directory browser | Inline keyboard filesystem navigator |
| Window picker | Inline keyboard for existing windows |
| Claude commands | /clear (with session reset), /compact, /cost, /help, /memory, /esc |
| Topic close | Kill window, unbind, cleanup |

### Phase 3 — Session Monitor

| Task | Description |
|------|-------------|
| JSONL parser | Parse Claude Code transcript entries, tool pairing |
| Monitor loop | Poll JSONL files, byte-offset tracking, session_map watching |
| Message formatting | Tool results, expandable quotes, thinking blocks |
| Markdown converter | Markdown → Telegram MarkdownV2 |
| Message queue | Per-user goroutine, merging, flood control, status management |
| Status polling | Detect Claude status line, show/update in topic |

### Phase 4 — Rich Features

| Task | Description |
|------|-------------|
| Interactive UI | Detect AskUserQuestion/ExitPlanMode, inline keyboard navigation |
| Screenshots | ANSI capture → PNG rendering with control keyboard |
| History | /history with pagination |
| Bash capture | ! command output capture and streaming |

### Phase 5 — Minuano Integration

| Task | Description |
|------|-------------|
| Bridge | Exec minuano CLI, parse output |
| /project | Bind topic to minuano project |
| /tasks | Show ready tasks for project |
| /pick | Single-task mode: generate prompt, send to Claude |
| /auto | Autonomous mode: generate loop prompt, send to Claude |
| /batch | Batch mode: generate multi-task prompt, send to Claude |
| Prompt templates | Self-contained instructions for each mode |

### Phase 6 — Polish

| Task | Description |
|------|-------------|
| Multi-user | Verify all state is per-user, notifications fan out correctly |
| Error handling | Stale windows, topic deletion detection, graceful degradation |
| Startup recovery | Full reconciliation of persisted state vs live tmux |
| Rate limiting | Telegram flood control, status deduplication |

---

## What Minuano Needs (prerequisites from the other project)

For Phase 5 to work, Minuano needs these additions:

1. `minuano prompt single <task-id>` — output a self-contained prompt for one task.
2. `minuano prompt auto --project X` — output a loop prompt scoped to project.
3. `minuano prompt batch <id1> <id2> ...` — output a multi-task prompt.
4. `--project` filter on `minuano-claim` script.
5. `minuano-pick <task-id>` script — claim a specific task by ID.

These are independent of Tramuntana and can be implemented in parallel.

---

## Send-Keys Strategy

The critical UX detail: how to deliver prompts to an existing Claude session.

**Short messages** (user text, commands): Direct `tmux send-keys` with literal text + Enter.
Same 500ms delay as CCBot between text and Enter to avoid TUI misinterpretation.

**Long prompts** (Minuano task specs): Write to temp file, send a short instruction
to Claude referencing the file:

```
Please read and follow the instructions in /tmp/tramuntana-task-XXXX.md
```

This avoids tmux send-keys character limits and keeps the terminal readable.

**`!` commands**: Send `!` first, wait 1s for bash mode, then send the rest.

---

## Open Questions

1. **Telegram bot library**: `go-telegram-bot-api/v5` is simpler but lower-level.
   `telebot/v4` has more middleware. Pick based on forum topic support quality.
2. **Font embedding for screenshots**: Embed JetBrains Mono in the binary via `//go:embed`
   or require it installed? Embedding is more portable.
3. **Session monitor: fsnotify vs polling?** CCBot polls. fsnotify would be more
   efficient but adds complexity for cross-filesystem edge cases. Start with polling,
   optimize later.
4. **Minuano output format**: Currently human-readable tables. May need `--json` flag
   on minuano commands for reliable parsing by the bridge.
