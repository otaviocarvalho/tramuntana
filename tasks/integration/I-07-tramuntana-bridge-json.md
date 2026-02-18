# I-07 — Update Tramuntana Bridge to Use `--json`

## Repo

**tramuntana** (`/home/otavio/code/tramuntana`)

## Goal

Update `internal/minuano/bridge.go` (Tramuntana task 24) to use Minuano's
`--json` flag for reliable output parsing.

## Steps

1. `Status(project)`: call `minuano status --project X --json`, unmarshal
   JSON array into `[]Task` structs.

2. `Show(taskID)`: call `minuano show <id> --json`, unmarshal JSON into
   `TaskDetail` struct (task + context).

3. `Prompt(mode, args...)`: call `minuano prompt <mode> <args>`, capture
   raw stdout as string (Markdown, no JSON parsing needed).

4. `Tree(project)`: call `minuano tree --project X`, capture raw stdout
   (plain text, for display only — no JSON needed).

5. Define Go structs matching Minuano's JSON output:
   ```go
   type Task struct {
       ID          string  `json:"id"`
       Title       string  `json:"title"`
       Status      string  `json:"status"`
       Priority    int     `json:"priority"`
       ProjectID   *string `json:"project_id,omitempty"`
       // ...
   }
   type TaskContext struct {
       Kind      string  `json:"kind"`
       Content   string  `json:"content"`
       AgentID   *string `json:"agent_id,omitempty"`
       // ...
   }
   type ShowOutput struct {
       Task    Task          `json:"task"`
       Context []TaskContext  `json:"context"`
   }
   ```

6. Error handling: if `--json` fails (older minuano version), return a clear
   error suggesting the user update minuano.

## Acceptance

- Bridge parses JSON from minuano status/show without error.
- `Prompt()` returns raw Markdown text.
- Struct fields match Minuano's JSON output.
- Errors from minuano are wrapped with command context.

## Depends on

- I-01 (minuano --json flag exists)
- Tramuntana Task 24 (bridge skeleton)
