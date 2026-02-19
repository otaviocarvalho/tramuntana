# Launching Tramuntana + Minuano

This guide explains how **Tramuntana** (Telegram bridge) and **Minuano** (task coordination) work together, what to start first, and how to operate agents from Telegram.

---

## Architecture Overview

```
Telegram Group Forum
   │
   ▼
Tramuntana (Go bot)
   │
   ├── Maps 1 Topic = 1 tmux Window = 1 Claude Code process
   │
   ├── CLI exec ──► Minuano (task queries, prompt generation)
   │                    │
   │                    ▼
   │               PostgreSQL (task DB)
   │
   └── tmux ──► Claude Code (sends prompts, reads output)
                    │
                    └── bash scripts ──► PostgreSQL (claim, done, observe)
```

**Tramuntana** is the user-facing layer. It runs a Telegram bot that maps group forum topics to tmux windows, each running a Claude Code session. Users interact with Claude through Telegram messages.

**Minuano** is the task coordination layer. It manages a PostgreSQL database of tasks, dependencies, context, and agent state. Tramuntana uses Minuano to give Claude structured work — rather than free-form chat, Claude follows task specs with test gates.

---

## Prerequisites

1. **Minuano binary** — built and accessible:
   ```bash
   cd /path/to/minuano
   go build ./cmd/minuano
   ```

2. **PostgreSQL running** (via Minuano):
   ```bash
   minuano up
   minuano migrate
   ```

3. **Tasks loaded** in Minuano:
   ```bash
   minuano add "Build auth system" --project myapp --priority 8 \
     --body "Implement OAuth2 authentication" \
     --test-cmd "go test ./..."
   ```

4. **Telegram bot token** from [@BotFather](https://t.me/BotFather)

5. **tmux** installed and available

---

## Startup Order

### Step 1: Start Minuano infrastructure

```bash
cd /path/to/minuano
minuano up       # PostgreSQL via Docker
minuano migrate  # Create tables
```

### Step 2: Configure Tramuntana

Create a `.env` file:

```bash
# Required
TELEGRAM_BOT_TOKEN=123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
ALLOWED_USERS=123456789

# Minuano integration
MINUANO_BIN=/path/to/minuano/minuano
MINUANO_DB=postgres://minuano:minuano@localhost:5432/minuanodb?sslmode=disable

# tmux
TMUX_SESSION_NAME=tramuntana

# Optional
WORKING_DIR=/path/to/your/codebase
LOG_DIR=/tmp/tramuntana-logs
```

### Step 3: Start Tramuntana

```bash
cd /path/to/tramuntana
go build ./cmd/tramuntana
tramuntana serve
```

Tramuntana will:
- Connect to Telegram
- Create (or attach to) the tmux session
- Start listening for commands in the forum topics

---

## Telegram Commands

### Session Management

| Command | Description |
|---------|-------------|
| `/start` | Create a Claude Code session for this topic |
| `/stop` | Kill the Claude session |
| `/status` | Show session status |

### Minuano Integration

| Command | Description |
|---------|-------------|
| `/project <name>` | Bind this topic to a Minuano project |
| `/project` | Show current project binding |
| `/tasks` | List tasks for the bound project |
| `/pick <task-id>` | Send a single task to Claude |
| `/auto` | Start autonomous mode (loop through all ready tasks) |
| `/batch <id1> [id2] ...` | Send multiple tasks to Claude in order |

---

## How the Integration Works

### The Bridge

Tramuntana communicates with Minuano through `internal/minuano/bridge.go`. It executes the `minuano` binary as a subprocess and parses its output:

| Operation | Command executed | Output |
|-----------|-----------------|--------|
| List tasks | `minuano status --project X --json` | JSON array |
| Task detail | `minuano show <id> --json` | JSON object |
| Dependency tree | `minuano tree --project X` | Plain text |
| Single prompt | `minuano prompt single <id>` | Markdown |
| Auto prompt | `minuano prompt auto --project X` | Markdown |
| Batch prompt | `minuano prompt batch <id1> <id2>` | Markdown |

### Prompt Delivery

When you send `/pick`, `/auto`, or `/batch`:

1. Tramuntana calls `minuano prompt <mode> ...` to generate a self-contained Markdown prompt
2. The prompt is written to a temp file (`/tmp/tramuntana-task-XXXX.md`)
3. Tramuntana sends a reference to Claude via tmux: `"Please read and follow the instructions in /tmp/tramuntana-task-XXXX.md"`
4. Claude reads the file, claims the task(s), and works autonomously
5. Claude calls `minuano-done <id> "summary"` when finished (runs tests as the completion gate)
6. Tramuntana's session monitor captures Claude's output and routes it back to the Telegram topic

### Environment Bootstrap

When Tramuntana creates a tmux window for a topic, it injects the Minuano environment:

```bash
export DATABASE_URL="postgres://minuano:minuano@localhost:5432/minuanodb?sslmode=disable"
export AGENT_ID="tramuntana-<topic-slug>"
export PATH="$PATH:/path/to/minuano/scripts"
```

This happens before Claude launches, so the Claude process inherits the variables. The agent scripts (`minuano-claim`, `minuano-done`, etc.) use `DATABASE_URL` to talk directly to PostgreSQL.

---

## Example: Full Workflow

### 1. Set up tasks in Minuano

```bash
# Add tasks with dependencies
minuano add "Design API schema" --project backend --priority 10 \
  --body "Design REST API schema for users and posts" \
  --test-cmd "go test ./..."

minuano add "Implement users endpoint" --project backend --priority 7 \
  --after design-api \
  --body "CRUD endpoints for /api/users" \
  --test-cmd "go test ./internal/api/users/..."

minuano add "Implement posts endpoint" --project backend --priority 7 \
  --after design-api \
  --body "CRUD endpoints for /api/posts" \
  --test-cmd "go test ./internal/api/posts/..."

# Check the tree
minuano tree --project backend
#   ◎  design-api-x1       Design API schema (ready)
#   ├── ○  implement-users-y2  Implement users endpoint (pending)
#   └── ○  implement-posts-z3  Implement posts endpoint (pending)
```

### 2. Work tasks from Telegram

```
You:    /project backend
Bot:    Bound to project: backend

You:    /tasks
Bot:    Tasks [backend]:
          ◎ design-api-x1 — Design API schema [ready]
          ○ implement-users-y2 — Implement users endpoint [pending]
          ○ implement-posts-z3 — Implement posts endpoint [pending]

You:    /pick design-api-x1
Bot:    Working on task design-api-x1...

        [Claude reads the prompt, designs the schema, runs tests]

Claude: Task design-api-x1 completed: designed REST API schema with OpenAPI spec

You:    /tasks
Bot:    Tasks [backend]:
          ✓ design-api-x1 — Design API schema [done]
          ◎ implement-users-y2 — Implement users endpoint [ready]
          ◎ implement-posts-z3 — Implement posts endpoint [ready]

You:    /auto
Bot:    Starting autonomous mode for project backend...

        [Claude loops: claims users → works → done → claims posts → works → done → empty → stops]
```

### 3. Or run agents directly with Minuano (no Telegram)

```bash
# Single agent working through the project
minuano run --project backend

# Multiple agents in parallel
minuano spawn --project backend --count 2

# Monitor
minuano watch
```

---

## Troubleshooting

- **`MINUANO_BIN not set`**: Add it to your `.env` or set the environment variable. Must point to the built `minuano` binary.
- **`minuano: command not found`** in agent scripts: The `PATH` wasn't set up. Check that `MINUANO_BIN` is correct and the scripts directory is adjacent to it.
- **`/tasks` returns empty**: Make sure you bound the topic to the right project with `/project <name>` and that tasks exist (`minuano status --project <name>`).
- **Claude doesn't pick up tasks**: Check that the tmux window has the right environment (`tmux show-environment -t <window>`). Verify `DATABASE_URL` and `AGENT_ID` are set.
- **Tasks stuck in `pending`**: Dependencies aren't met. Use `minuano tree --project X` to see the dependency graph. A task becomes `ready` only when all its dependencies are `done`.
- **Tramuntana can't find Minuano**: Ensure `MINUANO_BIN` is an absolute path to the built binary. Run `$MINUANO_BIN status` manually to test.
