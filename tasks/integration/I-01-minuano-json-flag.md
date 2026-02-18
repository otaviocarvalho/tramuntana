# I-01 — `--json` Flag on Minuano Status and Show

## Repo

**minuano** (`/home/otavio/code/minuano`)

## Goal

Add `--json` flag to `minuano status` and `minuano show` so Tramuntana's bridge
can parse output reliably instead of scraping tables.

## Steps

1. **`cmd/minuano/cmd_status.go`**: Add `--json` bool flag. When set, serialize
   the `[]*db.Task` slice to JSON and print to stdout instead of the table.
   Use `encoding/json` with `json.MarshalIndent`.

2. **`cmd/minuano/cmd_show.go`**: Add `--json` bool flag. When set, serialize
   a struct containing the task + its context entries to JSON.
   ```go
   type ShowOutput struct {
       Task    *db.Task          `json:"task"`
       Context []*db.TaskContext  `json:"context"`
   }
   ```

3. Add JSON struct tags to `db.Task` and `db.TaskContext` in `internal/db/queries.go`:
   ```go
   type Task struct {
       ID          string          `json:"id"`
       Title       string          `json:"title"`
       Body        string          `json:"body"`
       Status      string          `json:"status"`
       Priority    int             `json:"priority"`
       Capability  *string         `json:"capability,omitempty"`
       // ... etc
   }
   ```

4. **`cmd/minuano/cmd_tree.go`**: Optionally add `--json` flag here too
   (lower priority, tree output is display-only in Tramuntana).

## Acceptance

- `minuano status --json` outputs a JSON array of task objects.
- `minuano status --project X --json` filters and outputs JSON.
- `minuano show <id> --json` outputs a JSON object with task + context.
- Existing table output is unchanged when `--json` is not set.
- JSON struct tags on all DB types.

## Depends on

- Minuano Task 10 (`minuano status` exists)
- Minuano Task 12 (`minuano show` exists)
