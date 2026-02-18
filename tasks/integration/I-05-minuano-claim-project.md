# I-05 — `--project` Flag on `minuano-claim`

## Repo

**minuano** (`/home/otavio/code/minuano`)

## Goal

Add `--project` flag to `scripts/minuano-claim` so auto mode only pulls tasks
from the agent's project scope.

## Steps

1. Modify `scripts/minuano-claim` to accept a `--project` argument:
   ```bash
   PROJECT=""
   while [[ $# -gt 0 ]]; do
     case "$1" in
       --project) PROJECT="$2"; shift 2 ;;
       *) shift ;;
     esac
   done
   ```

2. When `PROJECT` is set, add to the inner SELECT:
   ```sql
   AND project_id = '$PROJECT'
   ```

3. When `PROJECT` is empty, behavior is unchanged (claims from any project).

4. Test both paths:
   - `minuano-claim` — claims any ready task.
   - `minuano-claim --project auth-system` — claims only from that project.

## Acceptance

- `minuano-claim --project auth-system` only claims tasks with `project_id = 'auth-system'`.
- `minuano-claim` (no flag) works as before.
- No SQL injection risk from the project parameter.

## Depends on

- I-03 (project filter concept validated in Go)
- Minuano Task 24 (minuano-claim exists)
