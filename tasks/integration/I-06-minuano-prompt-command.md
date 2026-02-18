# I-06 — `minuano prompt` Command

## Repo

**minuano** (`/home/otavio/code/minuano`)

## Goal

New CLI subcommand that outputs self-contained prompts for Claude to execute.
This is the main bridge between Tramuntana and Minuano.

## Steps

1. Create `cmd/minuano/cmd_prompt.go`.

2. Add subcommand tree:
   ```
   minuano prompt
   ├── single <task-id>         # one task
   ├── auto --project <name>    # loop mode
   └── batch <id1> [id2] ...    # ordered list
   ```

3. **`minuano prompt single <task-id>`**:
   - Call `db.GetTaskWithContext(pool, taskID)`.
   - Output Markdown prompt:
     - Task title, ID, body.
     - Context entries (inherited, handoff, test_failure).
     - Instructions: work it, use `minuano-observe`/`minuano-handoff`/`minuano-done`.
     - Rule: do NOT loop. Return to interactive after this task.

4. **`minuano prompt auto --project <name>`**:
   - Verify project has ready tasks via `db.ListTasks(pool, &project)`.
   - Output Markdown prompt:
     - Loop instructions: `minuano-claim --project <name>`.
     - On empty output: stop, return to interactive.
     - On task JSON: work it, `minuano-done`, loop back.
     - Rules: never skip minuano-done, fix only what broke on test_failure.

5. **`minuano prompt batch <id1> <id2> ...`**:
   - Call `db.GetTaskWithContext` for each ID.
   - Output Markdown prompt:
     - Numbered task list with specs and context.
     - For each: `minuano-pick <id>` → work → `minuano-done`.
     - After all done: return to interactive.

6. All prompts include an environment section reminding the agent that
   `AGENT_ID`, `DATABASE_URL`, and scripts are already on PATH.

7. Output goes to stdout (no file writing — Tramuntana handles temp files).

## Acceptance

- `minuano prompt single design-auth` outputs a complete single-task prompt.
- `minuano prompt auto --project auth-system` outputs a loop prompt.
- `minuano prompt batch design-auth implement-auth` outputs a multi-task prompt.
- Prompts are valid Markdown, self-contained, and include all task context.
- An agent reading only the prompt would know exactly what to do.

## Depends on

- I-04 (minuano-pick script — referenced in batch prompt)
- I-05 (minuano-claim --project — referenced in auto prompt)
- Minuano Task 12 (`minuano show` / GetTaskWithContext exists)
