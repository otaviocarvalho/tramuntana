# Task 07 — Bot Setup

## Goal

Implement `internal/bot/bot.go` — Telegram bot connection, authorization, handler
registration, and the main serve loop.

## Reference

- CCBot: `src/ccbot/bot.py` — `create_bot()`, `post_init`, handler registration, `run_polling()`.

## Steps

1. Add dependency: `go get github.com/go-telegram-bot-api/telegram-bot-api/v5`.
2. Create `internal/bot/bot.go`:
   ```go
   type Bot struct {
       api     *tgbotapi.BotAPI
       config  *config.Config
       state   *state.State
       tmux    *tmux.Manager  // or just the package functions
       updates tgbotapi.UpdatesChannel
   }
   ```
3. Implement `New(cfg *config.Config) (*Bot, error)`:
   - Create BotAPI with token.
   - Set up update config: `AllowedUpdates: ["message", "callback_query"]`.
   - Initialize state (load from disk).
   - Ensure tmux session exists.
4. Implement `Run(ctx context.Context) error`:
   - Start polling for updates.
   - Route updates in a main goroutine loop:
     - `update.Message` → `handleMessage()`
     - `update.CallbackQuery` → `handleCallback()`
   - Respect context cancellation for graceful shutdown.
5. Implement authorization check in message/callback routing:
   - Check `config.IsAllowedUser(userID)`.
   - For group chats, also check `config.IsAllowedGroup(chatID)`.
   - Silently ignore unauthorized messages.
6. Wire `serve` command in `cmd/tramuntana/main.go` to create and run the bot.

## Acceptance

- Bot connects to Telegram and receives updates.
- Unauthorized users are silently ignored.
- Bot shuts down gracefully on SIGINT/SIGTERM.
- `tramuntana serve` starts the bot.

## Phase

2 — Core Bot

## Depends on

- Task 06
