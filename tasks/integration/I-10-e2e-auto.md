# I-10 — End-to-End Test: Auto Mode

## Repos

**Both** minuano and tramuntana

## Goal

Verify `/auto` works: Claude loops through multiple tasks and stops when empty.

## Steps

1. **Setup**:
   - `minuano add "Task A" --priority 8 --project e2e-test --body "Create /tmp/test-auto-a.txt"`
   - `minuano add "Task B" --after <task-a-id> --project e2e-test --body "Create /tmp/test-auto-b.txt"`
   - Bind topic: `/project e2e-test`

2. **Execute**:
   - Send `/auto`.

3. **Verify**:
   - Claude claims Task A (highest priority, no deps).
   - Claude completes Task A, calls minuano-done.
   - Trigger cascade makes Task B ready.
   - Claude claims Task B.
   - Claude completes Task B, calls minuano-done.
   - Claude runs minuano-claim --project, gets empty, stops.
   - Both tasks show as done in `minuano status`.
   - Claude returns to interactive mode.

## Acceptance

- Multiple tasks are worked in dependency order.
- Trigger cascade (done → ready) works.
- Agent exits loop when queue is empty.
- No token burn after completion.

## Depends on

- I-09 (pick mode works, validates basic flow)
