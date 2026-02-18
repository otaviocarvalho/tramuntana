# Tramuntana ↔ Minuano Integration Plan

Tramuntana talks to Minuano exclusively via CLI commands. No shared database,
no library imports, no IPC. Minuano is a black box that accepts commands and
returns text/JSON.

---

## The Three Agent Modes

| Mode | Telegram command | What happens |
|------|-----------------|--------------|
| **Pick** | `/pick <task-id>` | Claim one specific task, work it, done, return to interactive |
| **Auto** | `/auto` | Loop: claim next ready from project → work → done → repeat until empty |
| **Batch** | `/batch <id1> [id2] ...` | Work through the specified tasks in order, then return |

All three modes generate a self-contained prompt that gets delivered to the
existing Claude Code process in the topic's tmux window.

---

## What Minuano Needs (new features)

Tramuntana's bridge calls these commands. They don't exist yet in Minuano.

### M1. `--json` flag on `minuano status` and `minuano show`

The bridge needs machine-parseable output. Human-readable tables are fragile
to parse. Add `--json` to emit structured JSON to stdout.

```bash
# Current (table)
minuano status --project auth
#   ◎  design-auth-a1b2c    Design auth system    ready    —    0/3

# New (JSON)
minuano status --project auth --json
# [{"id":"design-auth-a1b2c","title":"Design auth system","status":"ready",...}]

minuano show design-auth --json
# {"id":"design-auth-a1b2c","title":"...","body":"...","context":[...],...}
```

### M2. `minuano prompt` command

New subcommand that outputs a self-contained prompt for Claude to execute.
Tramuntana writes this to a temp file and tells Claude to read it.

```
minuano prompt single <task-id>           # one task, work it, done, exit
minuano prompt auto --project <name>      # loop: claim from project, work, done, repeat
minuano prompt batch <id1> [id2] ...      # work these tasks in order
```

Each mode outputs Markdown instructions that include:
- The task spec(s) with full body and context
- The scripts to call (`minuano-done`, `minuano-observe`, etc.)
- Mode-specific behavior (exit after one, loop until empty, follow list)
- Environment requirements (`AGENT_ID`, `DATABASE_URL`, scripts PATH)

### M3. `--project` filter on `minuano-claim` script

The `scripts/minuano-claim` bash script needs a `--project` argument so auto
mode only pulls tasks from the correct project.

```bash
# Current: claims any ready task
minuano-claim

# New: claims only from a specific project
minuano-claim --project auth-system
```

### M4. `minuano-pick <task-id>` script

New bash script that claims a specific task by ID (not queue-based).
Used by pick and batch modes. Fails if the task isn't claimable (not ready,
already claimed, etc.).

```bash
minuano-pick design-auth-a1b2c
# prints JSON with task spec + context, same format as minuano-claim
```

### M5. `ClaimByID` query in `queries.go`

New Go function backing `minuano-pick`. Claims a specific task by exact/partial
ID, injects inherited context, updates agent — same pattern as `AtomicClaim`
but targeting a specific task instead of the next-from-queue.

### M6. Project filter on `AtomicClaim` query

Add optional `projectID` parameter to the existing `AtomicClaim` function.
When set, the WHERE clause includes `AND project_id = $N`. This backs the
`--project` flag on `minuano-claim`.

---

## What Tramuntana Needs (already in its task plan, verified)

These are covered by existing Tramuntana tasks 24 and 25:

- **Task 24**: `internal/minuano/bridge.go` — exec minuano commands, parse JSON output
- **Task 25**: `/project`, `/tasks`, `/pick`, `/auto`, `/batch` commands + `prompt.go`

However, the existing tasks need updates to account for the integration
specifics:

### T1. Update Task 24 — Bridge must use `--json` flag

The bridge should always call minuano with `--json` for machine parsing.
Fall back to table parsing only if `--json` isn't available.

### T2. Update Task 25 — Prompts reference temp files

The `/pick`, `/auto`, `/batch` commands must:
1. Call `minuano prompt <mode> ...` to generate the prompt text
2. Write it to a temp file (`/tmp/tramuntana-task-XXXX.md`)
3. Send a short instruction to Claude via tmux: reference the temp file
4. Set environment in the tmux window before sending the prompt

### T3. Update Task 25 — Environment bootstrap

Before any Minuano command, the tmux window needs:
```bash
export DATABASE_URL="..."
export AGENT_ID="tramuntana-<topic-slug>"
export PATH="$PATH:/path/to/minuano/scripts"
```

This should happen once on window creation (not on every command).

---

## Integration Task Order

These tasks span both repos. Execute them in this order:

```
Minuano side (can be done in parallel with Tramuntana phases 1-4):
  I-01  --json flag on status/show
  I-02  ClaimByID query
  I-03  Project filter on AtomicClaim
  I-04  minuano-pick script
  I-05  --project flag on minuano-claim
  I-06  minuano prompt command

Tramuntana side (after Tramuntana task 23, requires Minuano I-01..I-06):
  I-07  Update bridge to use --json (Tramuntana task 24 extension)
  I-08  Update commands to use minuano prompt (Tramuntana task 25 extension)
  I-09  End-to-end test: pick mode
  I-10  End-to-end test: auto mode
  I-11  End-to-end test: batch mode
```

### Dependency Graph

```
Minuano:     I-01 ──────────────────────────────┐
             I-02 → I-04 ──────────────────────┐│
             I-03 → I-05 ──────────────────────┤│
             I-06 (depends on I-04, I-05) ──────┤│
                                                ││
Tramuntana:  Tasks 01-23 (independent) ─────────┤│
             I-07 (depends on I-01, Task 24) ───┤│
             I-08 (depends on I-06, Task 25) ───┤│
             I-09 (depends on I-07, I-08) ──────┘│
             I-10 (depends on I-09) ─────────────┘
             I-11 (depends on I-09) ──────────────
```

---

## Communication Protocol

### Tramuntana → Minuano

All via CLI exec. No sockets, no DB sharing, no files (except temp prompt files).

| Operation | Command | Output format |
|-----------|---------|---------------|
| List tasks | `minuano status --project X --json` | JSON array of task objects |
| Task detail | `minuano show <id> --json` | JSON task object with context |
| Dependency tree | `minuano tree --project X` | Plain text (for display only) |
| Generate prompt | `minuano prompt <mode> [args]` | Markdown text (prompt content) |

### Minuano → Tramuntana

None. Minuano doesn't know Tramuntana exists. The Claude agent inside the tmux
window calls `minuano-done`, `minuano-observe`, `minuano-handoff` directly
against the Minuano database. Tramuntana observes the results through Claude's
JSONL output (which the session monitor already reads).

### Feedback Loop

```
User sends /pick design-auth
  → Tramuntana: minuano prompt single design-auth → writes /tmp/prompt.md
  → Tramuntana: tmux send-keys "read /tmp/prompt.md and follow the instructions"
  → Claude reads prompt, starts working
  → Claude calls minuano-done design-auth "implemented auth"
  → Claude's JSONL output captured by Tramuntana's session monitor
  → Monitor routes output to the Telegram topic
  → User sees progress in Telegram
```

---

## Environment Setup Detail

When Tramuntana creates a new tmux window (directory browser → confirm), it should
inject Minuano environment if `MINUANO_DB` is configured:

```go
// In tmux window creation flow
if config.MinuanoDB != "" {
    tmux.SendKeys(wid, fmt.Sprintf(`export DATABASE_URL=%q`, config.MinuanoDB))
    tmux.SendKeys(wid, fmt.Sprintf(`export AGENT_ID="tramuntana-%s"`, windowName))
    if config.MinuanoBin != "" {
        scriptsDir := filepath.Join(filepath.Dir(config.MinuanoBin), "..", "scripts")
        tmux.SendKeys(wid, fmt.Sprintf(`export PATH="$PATH:%s"`, scriptsDir))
    }
}
```

This happens before `claude` is launched, so the Claude process inherits the env.

---

## Prompt Templates

### Single Mode (`minuano prompt single <task-id>`)

```markdown
# Task: <title>

## Specification
<body>

## Context
<inherited/handoff/test_failure entries>

## Instructions
1. Work on this task following the specification above.
2. Write observations as you discover things: `minuano-observe <id> "note"`
3. Write handoffs before risky changes: `minuano-handoff <id> "note"`
4. When complete: `minuano-done <id> "summary of what was built"`
5. Do NOT loop or claim another task. Return to interactive mode after this task.
```

### Auto Mode (`minuano prompt auto --project <name>`)

```markdown
# Autonomous Task Mode — Project: <name>

Work through all ready tasks in this project until the queue is empty.

## Loop
1. Run `minuano-claim --project <name>` to claim the next ready task.
2. If no output: queue is empty. Stop and return to interactive mode.
3. Read the JSON output for task spec and context.
4. Work on the task following its specification.
5. Write observations: `minuano-observe <id> "note"`
6. When complete: `minuano-done <id> "summary"`
7. Go back to step 1.

## Rules
- Never skip `minuano-done`. It runs tests.
- If tests fail, fix only what broke.
- If stuck, write a handoff note and call `minuano-done` to record cleanly.
```

### Batch Mode (`minuano prompt batch <id1> <id2> ...`)

```markdown
# Batch Task Mode

Work through these tasks in order:

## Task 1: <title1>
<spec + context>

## Task 2: <title2>
<spec + context>

...

## Instructions
For each task:
1. Run `minuano-pick <id>` to claim it.
2. Work on it following its specification.
3. Write observations: `minuano-observe <id> "note"`
4. When complete: `minuano-done <id> "summary"`
5. Move to the next task.

After all tasks are complete, return to interactive mode.
```
