# I-04 — `minuano-pick` Script

## Repo

**minuano** (`/home/otavio/code/minuano`)

## Goal

Create `scripts/minuano-pick` — a bash script that claims a specific task by
ID, used by Tramuntana's pick and batch modes.

## Steps

1. Create `scripts/minuano-pick`:
   ```bash
   #!/usr/bin/env bash
   set -euo pipefail

   TASK_ID="${1:?Usage: minuano-pick <task-id>}"
   AGENT_ID="${AGENT_ID:?AGENT_ID not set}"
   DB="${DATABASE_URL:?DATABASE_URL not set}"
   ```

2. SQL: same CTE structure as `minuano-claim` but with `WHERE id = '$TASK_ID'`
   instead of the priority-ordered queue SELECT.

3. Should:
   - Fail with clear error if task is not ready.
   - Print JSON output (same format as `minuano-claim`) on success.
   - Inject inherited context from done dependencies.
   - Update agent status.

4. Make executable: `chmod +x scripts/minuano-pick`.

5. Use `quote_literal()` for any interpolated values (SQL safety).

## Acceptance

- `AGENT_ID=x DATABASE_URL=... minuano-pick design-auth` claims that task.
- Prints JSON with task spec + context on success.
- Fails with clear message if task isn't ready.
- Same output format as `minuano-claim`.

## Depends on

- I-02 (ClaimByID query exists to verify the pattern)
- Minuano Task 24 (agent scripts pattern exists)
