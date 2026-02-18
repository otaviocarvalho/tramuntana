# I-11 — End-to-End Test: Batch Mode

## Repos

**Both** minuano and tramuntana

## Goal

Verify `/batch` works: Claude works through a specific list of tasks in order.

## Steps

1. **Setup**:
   - Create 3 independent ready tasks in the same project:
     ```
     minuano add "Batch 1" --project e2e-test --body "Create /tmp/test-batch-1.txt"
     minuano add "Batch 2" --project e2e-test --body "Create /tmp/test-batch-2.txt"
     minuano add "Batch 3" --project e2e-test --body "Create /tmp/test-batch-3.txt"
     ```
   - Bind topic: `/project e2e-test`

2. **Execute**:
   - Send `/batch <id1> <id3>` (skip id2 — verify cherry-picking works).

3. **Verify**:
   - Claude claims and works task 1 via `minuano-pick`.
   - Claude claims and works task 3 via `minuano-pick`.
   - Task 2 is NOT touched (still ready).
   - Claude returns to interactive mode after the two tasks.
   - `minuano status` shows tasks 1 and 3 done, task 2 still ready.

## Acceptance

- Only the specified tasks are worked.
- Tasks are worked in the order specified.
- Non-specified tasks are untouched.
- Claude returns to interactive after the batch.

## Depends on

- I-09 (pick mode works, validates basic flow)
