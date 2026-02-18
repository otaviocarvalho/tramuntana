# I-08 — Update Tramuntana Commands to Use `minuano prompt`

## Repo

**tramuntana** (`/home/otavio/code/tramuntana`)

## Goal

Update `/pick`, `/auto`, `/batch` handlers (Tramuntana task 25) to use
`minuano prompt` for generating prompts instead of building them in Go.

## Steps

1. **`/pick <task-id>`**:
   - Call `bridge.Prompt("single", taskID)` → get Markdown text.
   - Write to temp file: `/tmp/tramuntana-pick-XXXX.md`.
   - Ensure env vars are set in tmux window (one-time bootstrap).
   - Send to tmux: `"Please read and follow the instructions in /tmp/tramuntana-pick-XXXX.md"`.
   - Reply to Telegram: `"Working on task <id>..."`.

2. **`/auto`**:
   - Look up project binding for this thread.
   - Call `bridge.Prompt("auto", "--project", project)` → get Markdown text.
   - Write to temp file: `/tmp/tramuntana-auto-XXXX.md`.
   - Send to tmux.
   - Reply: `"Starting autonomous mode for project <name>..."`.

3. **`/batch <id1> [id2] ...`**:
   - Parse task IDs from args.
   - Call `bridge.Prompt("batch", ids...)` → get Markdown text.
   - Write to temp file: `/tmp/tramuntana-batch-XXXX.md`.
   - Send to tmux.
   - Reply: `"Working on batch: <ids>..."`.

4. **Environment bootstrap** (shared helper):
   ```go
   func (b *Bot) ensureMinuanoEnv(windowID string) error {
       // Check if already bootstrapped (track in state)
       // If not: send export commands via tmux
       // Mark as bootstrapped
   }
   ```
   Called before the first minuano command in any window.

5. **Temp file cleanup**: delete temp files after 1 hour via a background
   goroutine, or use `os.CreateTemp` and let the OS handle it.

## Acceptance

- `/pick` generates prompt via minuano, writes to temp file, sends reference.
- `/auto` scopes to the bound project.
- `/batch` handles multiple task IDs.
- Environment is bootstrapped before first command.
- Claude receives the prompt and can act on it.

## Depends on

- I-06 (minuano prompt command exists)
- I-07 (bridge can call minuano prompt)
- Tramuntana Task 25 (command handlers skeleton)
