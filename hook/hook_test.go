package hook

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestIsHookInstalled(t *testing.T) {
	tests := []struct {
		name     string
		settings map[string]any
		command  string
		want     bool
	}{
		{
			name:     "empty settings",
			settings: map[string]any{},
			command:  "/usr/bin/tramuntana hook",
			want:     false,
		},
		{
			name: "already installed",
			settings: map[string]any{
				"hooks": map[string]any{
					"SessionStart": []any{
						map[string]any{
							"type":    "command",
							"command": "/usr/bin/tramuntana hook",
							"timeout": 5,
						},
					},
				},
			},
			command: "/usr/bin/tramuntana hook",
			want:    true,
		},
		{
			name: "different hook",
			settings: map[string]any{
				"hooks": map[string]any{
					"SessionStart": []any{
						map[string]any{
							"type":    "command",
							"command": "some-other-hook",
						},
					},
				},
			},
			command: "/usr/bin/tramuntana hook",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isHookInstalled(tt.settings, tt.command)
			if got != tt.want {
				t.Errorf("isHookInstalled = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstall_CreatesSettingsFile(t *testing.T) {
	// Create a fake executable
	tmpDir := t.TempDir()
	fakeExe := filepath.Join(tmpDir, "tramuntana")
	os.WriteFile(fakeExe, []byte("#!/bin/sh"), 0755)

	// Override HOME to temp dir
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Create .claude dir
	os.MkdirAll(filepath.Join(tmpDir, ".claude"), 0755)

	// We can't easily test Install() since it uses os.Executable(),
	// but we can test the settings file manipulation logic
	settingsPath := filepath.Join(tmpDir, ".claude", "settings.json")

	// Start with empty settings
	settings := map[string]any{}
	hookCommand := fakeExe + " hook"

	// Verify not installed
	if isHookInstalled(settings, hookCommand) {
		t.Error("should not be installed initially")
	}

	// Add hook
	hooks := make(map[string]any)
	hookEntry := map[string]any{
		"type":    "command",
		"command": hookCommand,
		"timeout": 5,
	}
	hooks["SessionStart"] = []any{hookEntry}
	settings["hooks"] = hooks

	out, _ := json.MarshalIndent(settings, "", "  ")
	os.WriteFile(settingsPath, out, 0644)

	// Verify installed
	data, _ := os.ReadFile(settingsPath)
	var loaded map[string]any
	json.Unmarshal(data, &loaded)
	if !isHookInstalled(loaded, hookCommand) {
		t.Error("should be installed after adding")
	}
}

func TestUUIDRegex(t *testing.T) {
	valid := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"00000000-0000-0000-0000-000000000000",
	}
	invalid := []string{
		"not-a-uuid",
		"550e8400-e29b-41d4-a716",
		"",
	}

	for _, s := range valid {
		if !uuidRegex.MatchString(s) {
			t.Errorf("%q should match UUID regex", s)
		}
	}
	for _, s := range invalid {
		if uuidRegex.MatchString(s) {
			t.Errorf("%q should not match UUID regex", s)
		}
	}
}
