# I-02 — `ClaimByID` Query in Minuano

## Repo

**minuano** (`/home/otavio/code/minuano`)

## Goal

Add a `ClaimByID` function to `internal/db/queries.go` that claims a specific
task by ID — same as `AtomicClaim` but targeting a known task instead of
pulling from the queue.

## Steps

1. Add `ClaimByID(pool, taskID, agentID string) (*Task, error)` to `queries.go`.

2. The function should:
   - Resolve partial ID via `ResolvePartialID`.
   - Verify the task is claimable (status = 'ready', attempt < max_attempts).
   - In a transaction:
     a. `UPDATE tasks SET status='claimed', claimed_by=$1, claimed_at=NOW(), attempt=attempt+1 WHERE id=$2 AND status='ready'`
     b. Inject inherited context from done dependencies (same CTE as AtomicClaim).
     c. Update agent status if agent is registered.
   - Return the claimed task, or an error if not claimable.

3. Return a clear error for each failure case:
   - Task not found.
   - Task not in 'ready' status (already claimed, pending, done, failed).
   - Max attempts reached.

## Acceptance

- `ClaimByID(pool, "design-auth", "agent-1")` claims that specific task.
- Returns an error if the task isn't ready.
- Inherited context is injected from done dependencies.
- Partial ID matching works.

## Depends on

- Minuano Task 05 (queries.go exists)
