# Task 24 — Minuano Bridge

## Goal

Implement `internal/minuano/bridge.go` — execute Minuano CLI commands and parse
their output.

## Reference

- PLAN.md: Minuano Bridge section.
- Minuano CLI: `minuano status`, `minuano show`, `minuano tree`.

## Steps

1. Create `internal/minuano/bridge.go`:
   ```go
   type Bridge struct {
       Bin    string // path to minuano binary
       DBFlag string // optional --db flag value
   }
   ```
2. Implement command execution helper:
   - Build command with optional `--db` flag.
   - Exec via `os/exec`, capture stdout and stderr.
   - On non-zero exit: wrap stderr in error.
3. Implement `Status(project string) ([]Task, error)`:
   - Run `minuano status --project <project>`.
   - Parse table output into `[]Task` structs (id, subject, status, etc.).
   - Handle both `--json` (if available) and table format.
4. Implement `Show(taskID string) (*TaskDetail, error)`:
   - Run `minuano show <id>`.
   - Parse output into detailed task struct.
5. Implement `Tree(project string) (string, error)`:
   - Run `minuano tree --project <project>`.
   - Return raw text output (tree view).
6. Define task structs:
   ```go
   type Task struct {
       ID      string
       Subject string
       Status  string
       Project string
   }
   type TaskDetail struct {
       Task
       Description string
       DependsOn   []string
       Blocks      []string
   }
   ```

## Acceptance

- Bridge can execute minuano commands and parse output.
- `Status()` returns structured task lists.
- `Show()` returns detailed task info.
- Errors from minuano are properly wrapped.
- Works with or without `--db` flag.

## Phase

5 — Minuano Integration

## Depends on

- Task 02
