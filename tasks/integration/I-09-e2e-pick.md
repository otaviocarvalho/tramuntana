# I-09 — End-to-End Test: Pick Mode

## Repos

**Both** minuano and tramuntana

## Goal

Verify the full `/pick` flow works end-to-end:
Telegram → Tramuntana → minuano prompt → temp file → tmux → Claude → minuano-done

## Steps

1. **Setup**:
   - `minuano up && minuano migrate`
   - `minuano add "Test task" --body "Create a file called /tmp/test-pick.txt with the content 'hello'" --project e2e-test`
   - Start tramuntana: `tramuntana serve`
   - Create a Telegram topic, bind to a directory, bind project: `/project e2e-test`

2. **Execute**:
   - Send `/tasks` — verify the test task appears.
   - Send `/pick <task-id>`.

3. **Verify**:
   - Claude receives the prompt (visible in tmux).
   - Claude works on the task.
   - Claude calls `minuano-done`.
   - `minuano status --project e2e-test` shows the task as done.
   - Telegram topic shows Claude's progress via session monitor.

## Acceptance

- Full round-trip works without manual intervention.
- Task moves from ready → claimed → done.
- Telegram topic shows work progress.

## Depends on

- I-07, I-08 (Tramuntana integration complete)
- I-04, I-06 (Minuano pick/prompt commands exist)
