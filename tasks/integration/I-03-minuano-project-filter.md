# I-03 — Project Filter on AtomicClaim

## Repo

**minuano** (`/home/otavio/code/minuano`)

## Goal

Add an optional `projectID` parameter to `AtomicClaim` so agents can be scoped
to a specific project's task queue.

## Steps

1. Change `AtomicClaim` signature:
   ```go
   func AtomicClaim(pool *pgxpool.Pool, agentID string, capability *string, projectID *string) (*Task, error)
   ```

2. When `projectID` is not nil, add `AND project_id = $N` to the inner SELECT
   that finds the next ready task.

3. Update all callers of `AtomicClaim` (currently in `cmd/minuano/` if any,
   and in `scripts/minuano-claim` via SQL).

4. The existing behavior (no project filter) must be preserved when `projectID`
   is nil.

## Acceptance

- `AtomicClaim(pool, "agent-1", nil, &project)` only claims tasks from that project.
- `AtomicClaim(pool, "agent-1", nil, nil)` claims from any project (unchanged behavior).
- No regressions in existing claim behavior.

## Depends on

- Minuano Task 05 (queries.go exists)
