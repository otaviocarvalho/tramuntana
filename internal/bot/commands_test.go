package bot

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/otaviocarvalho/tramuntana/internal/config"
	"github.com/otaviocarvalho/tramuntana/internal/state"
)

func TestResolveWindow_Bound(t *testing.T) {
	b := &Bot{
		config: &config.Config{
			AllowedUsers: []int64{100},
		},
		state: state.NewState(),
	}
	b.state.BindThread("100", "42", "@5")

	// Simulate thread ID cache
	threadCacheMu.Lock()
	threadIDCache[1001] = 42
	threadCacheMu.Unlock()
	defer func() {
		threadCacheMu.Lock()
		delete(threadIDCache, 1001)
		threadCacheMu.Unlock()
	}()

	msg := &tgbotapi.Message{
		MessageID: 1001,
		From:      &tgbotapi.User{ID: 100},
	}

	windowID, bound := b.resolveWindow(msg)
	if !bound {
		t.Fatal("expected window to be bound")
	}
	if windowID != "@5" {
		t.Errorf("got window %q, want @5", windowID)
	}
}

func TestResolveWindow_Unbound(t *testing.T) {
	b := &Bot{
		config: &config.Config{
			AllowedUsers: []int64{100},
		},
		state: state.NewState(),
	}

	// Simulate thread ID cache
	threadCacheMu.Lock()
	threadIDCache[1002] = 99
	threadCacheMu.Unlock()
	defer func() {
		threadCacheMu.Lock()
		delete(threadIDCache, 1002)
		threadCacheMu.Unlock()
	}()

	msg := &tgbotapi.Message{
		MessageID: 1002,
		From:      &tgbotapi.User{ID: 100},
	}

	_, bound := b.resolveWindow(msg)
	if bound {
		t.Error("expected window to be unbound")
	}
}

func TestWindowIDFromKey(t *testing.T) {
	tests := []struct {
		key  string
		want string
	}{
		{"tramuntana:@5", "@5"},
		{"session:@12", "@12"},
		{"nocolon", "nocolon"},
		{"a:b:@3", "@3"},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := windowIDFromKey(tt.key)
			if got != tt.want {
				t.Errorf("windowIDFromKey(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestHandleCommand_Routes(t *testing.T) {
	tests := []struct {
		cmd     string
		forward bool // true if it should use forwardCommand
	}{
		{"clear", true},
		{"compact", true},
		{"cost", true},
		{"help", true},
		{"memory", true},
	}

	for _, tt := range tests {
		t.Run(tt.cmd, func(t *testing.T) {
			msg := &tgbotapi.Message{
				MessageID: 1,
				Text:      "/" + tt.cmd,
				Entities: []tgbotapi.MessageEntity{
					{Type: "bot_command", Offset: 0, Length: len(tt.cmd) + 1},
				},
			}
			if msg.Command() != tt.cmd {
				t.Errorf("Command() = %q, want %q", msg.Command(), tt.cmd)
			}
		})
	}
}

func TestHandleCommand_EscRoute(t *testing.T) {
	msg := &tgbotapi.Message{
		MessageID: 1,
		Text:      "/esc",
		Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: 4},
		},
	}
	if msg.Command() != "esc" {
		t.Errorf("Command() = %q, want esc", msg.Command())
	}
}
