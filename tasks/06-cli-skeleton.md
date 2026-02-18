# Task 06 — CLI Skeleton

## Goal

Set up the cobra-based CLI framework with `serve`, `hook`, and `version` subcommands.

## Reference

- Minuano: `cmd/minuano/main.go` — cobra root command pattern.

## Steps

1. Add dependency: `go get github.com/spf13/cobra`.
2. Update `cmd/tramuntana/main.go` with a cobra root command named `tramuntana`.
3. Add subcommands:
   - `serve` — run the Telegram bot (placeholder for now, just prints "starting...").
     - Flag: `--config PATH` — path to `.env` override (calls `godotenv.Load(path)` before config).
   - `hook` — call `hook.Run()` from the hook package.
     - Flag: `--install` — call `hook.Install()` instead.
   - `version` — print version string (hardcode `"tramuntana v0.1.0"` for now).
4. Root command should print usage/help when called with no subcommand.
5. Wire up config loading in `serve` command's `PreRunE`.

## Acceptance

- `go build ./cmd/tramuntana && ./tramuntana --help` prints the command tree.
- `./tramuntana serve --help` shows the `--config` flag.
- `./tramuntana hook --install` triggers the hook installation.
- `./tramuntana version` prints the version.

## Phase

2 — Core Bot

## Depends on

- Task 05
