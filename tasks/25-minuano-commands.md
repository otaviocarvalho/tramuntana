# Task 25 — Minuano Commands

## Goal

Implement `/project`, `/tasks`, `/pick`, `/auto`, `/batch` bot commands for
Minuano integration.

## Reference

- PLAN.md: Minuano Integration section, Minuano commands table, prompt.go spec.

## Steps

1. Add to `internal/bot/commands.go` or create `internal/bot/minuano_commands.go`.

2. Implement `/project <name>`:
   - Parse project name from command args.
   - Store binding: `state.BindProject(threadID, projectName)`.
   - Reply with confirmation: `"Bound to project: <name>"`.

3. Implement `/tasks`:
   - Look up project binding for this thread.
   - If no project bound, reply with error.
   - Call `bridge.Status(project)` to get ready tasks.
   - Format task list and reply.

4. Create `internal/minuano/prompt.go` — prompt generation:
   - **Single mode** (`/pick <task-id>`):
     - Call `bridge.Show(taskID)` to get full task spec.
     - Build prompt: task spec + instructions (work it, mark done, exit).
   - **Auto mode** (`/auto`):
     - Call `bridge.Status(project)` to verify tasks exist.
     - Build prompt: loop instructions (claim from project, work, done, repeat, exit on empty).
   - **Batch mode** (`/batch <id1> <id2> ...`):
     - Call `bridge.Show()` for each task ID.
     - Build prompt: ordered task list with specs, work each in order.

5. Implement `/pick <task-id>`:
   - Resolve window, get project binding.
   - Generate single-mode prompt.
   - Write prompt to temp file (long prompts avoid send-keys limits).
   - Send to tmux: `"Please read and follow the instructions in /tmp/tramuntana-task-XXXX.md"`.
   - Reply `"Working on task <id>..."`.

6. Implement `/auto`:
   - Resolve window, get project binding.
   - Generate auto-mode prompt.
   - Write to temp file, send reference to tmux.
   - Reply `"Starting autonomous mode for project <name>..."`.

7. Implement `/batch <id1> [id2] ...`:
   - Parse task IDs from args.
   - Generate batch-mode prompt.
   - Write to temp file, send reference to tmux.
   - Reply `"Working on batch: <ids>..."`.

8. Environment setup for Minuano commands:
   - Before sending prompt, ensure env vars are set in the tmux window:
     `MINUANO_BIN`, `DATABASE_URL`, `AGENT_ID`.

## Acceptance

- `/project` binds a topic to a Minuano project.
- `/tasks` shows ready tasks for the bound project.
- `/pick` sends a single-task prompt to Claude.
- `/auto` sends a loop prompt to Claude.
- `/batch` sends a multi-task prompt to Claude.
- Long prompts use temp file + reference (not send-keys).

## Phase

5 — Minuano Integration

## Depends on

- Task 24
- Task 11
