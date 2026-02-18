# Task 19 — Interactive UI

## Goal

Detect AskUserQuestion, ExitPlanMode, PermissionPrompt, and other interactive UI
prompts from Claude Code's terminal, and render them as inline keyboards in Telegram.

## Reference

- CCBot: `src/ccbot/terminal_parser.py` — `UIPattern`, `extract_interactive_content()`,
  `is_interactive_ui()`.
- CCBot: `src/ccbot/handlers/interactive_ui.py` — `handle_interactive_ui()`,
  keyboard building, navigation callbacks.
- PLAN.md: Interactive UI Detection table.

## Steps

1. Add to `internal/monitor/terminal.go` — interactive UI detection:
   - Define UI patterns with top/bottom markers (ordered by priority):
     | Pattern | Top marker | Bottom marker |
     |---------|-----------|--------------|
     | ExitPlanMode | `"Would you like to proceed?"` or `"Claude has written up a plan"` | `"ctrl-g to edit"` or `"Esc to"` |
     | AskUserQuestion (multi-tab) | `"← [☐✔☒]"` | last non-empty line |
     | AskUserQuestion (single-tab) | `"[☐✔☒]"` | `"Enter to select"` |
     | PermissionPrompt | `"Do you want to proceed?"` | `"Esc to cancel"` |
     | RestoreCheckpoint | `"Restore the code"` | `"Enter to continue"` |
     | Settings | `"Settings:.*tab to cycle"` | `"Esc to cancel"` or `"Type to filter"` |
   - `IsInteractiveUI(paneText string) bool` — quick check.
   - `ExtractInteractiveContent(paneText string) (UIContent, bool)` — extract type + text between markers.

2. Create `internal/bot/interactive.go`:
   - `handleInteractiveUI(userID, threadID int64, windowID string)`:
     - Capture pane (plain text).
     - Detect and extract interactive content.
     - Build inline keyboard: ↑ ↓ ← → Space Tab Esc Enter Refresh.
     - For RestoreCheckpoint: vertical-only layout (no ← →).
     - Send or edit interactive message.
   - Track state per (user, thread): `interactiveMsgs map[key]int` (message_id),
     `interactiveMode map[key]string` (window_id).

3. Two triggers:
   - **JSONL trigger**: when monitor sees `tool_use` with `AskUserQuestion` or
     `ExitPlanMode` tool name → call `handleInteractiveUI()`.
   - **Polling trigger**: in status polling loop, check `IsInteractiveUI()` →
     call `handleInteractiveUI()` (catches PermissionPrompt which doesn't fire via JSONL).

4. Handle navigation callbacks:
   - Arrow keys → `tmux send-keys Up/Down/Left/Right`.
   - Space → `tmux send-keys Space`.
   - Tab → `tmux send-keys Tab`.
   - Esc → `tmux send-keys Escape`.
   - Enter → `tmux send-keys Enter`.
   - Refresh → re-capture and re-render.
   - After each key: wait 300ms, refresh the keyboard display.

## Acceptance

- All UI pattern types are detected from terminal output.
- Interactive keyboard appears when Claude prompts for input.
- Navigation keys work correctly.
- Keyboard updates after each action.
- Both JSONL and polling triggers work.

## Phase

4 — Rich Features

## Depends on

- Task 18
