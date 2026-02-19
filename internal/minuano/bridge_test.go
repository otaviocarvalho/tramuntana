package minuano

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewBridge(t *testing.T) {
	b := NewBridge("/usr/bin/minuano", "/path/to/db")
	if b.Bin != "/usr/bin/minuano" {
		t.Errorf("Bin = %q", b.Bin)
	}
	if b.DBFlag != "/path/to/db" {
		t.Errorf("DBFlag = %q", b.DBFlag)
	}
}

func TestNewBridge_NoDBFlag(t *testing.T) {
	b := NewBridge("/usr/bin/minuano", "")
	if b.DBFlag != "" {
		t.Error("DBFlag should be empty")
	}
}

func TestTaskJSON(t *testing.T) {
	jsonStr := `{
		"id": "task-1",
		"title": "Fix bug",
		"body": "Detailed spec",
		"status": "ready",
		"priority": 5,
		"attempt": 0,
		"max_attempts": 3
	}`
	var task Task
	if err := json.Unmarshal([]byte(jsonStr), &task); err != nil {
		t.Fatal(err)
	}
	if task.ID != "task-1" {
		t.Errorf("ID = %q", task.ID)
	}
	if task.Title != "Fix bug" {
		t.Errorf("Title = %q", task.Title)
	}
	if task.Status != "ready" {
		t.Errorf("Status = %q", task.Status)
	}
	if task.Priority != 5 {
		t.Errorf("Priority = %d", task.Priority)
	}
}

func TestTaskJSON_WithOptionalFields(t *testing.T) {
	project := "proj-1"
	claimedBy := "agent-1"
	jsonStr := `{
		"id": "task-2",
		"title": "Another task",
		"body": "",
		"status": "claimed",
		"priority": 3,
		"capability": "code",
		"claimed_by": "agent-1",
		"project_id": "proj-1",
		"attempt": 1,
		"max_attempts": 3
	}`
	var task Task
	if err := json.Unmarshal([]byte(jsonStr), &task); err != nil {
		t.Fatal(err)
	}
	if task.Capability == nil || *task.Capability != "code" {
		t.Error("Capability should be 'code'")
	}
	if task.ClaimedBy == nil || *task.ClaimedBy != claimedBy {
		t.Error("ClaimedBy should be set")
	}
	if task.ProjectID == nil || *task.ProjectID != project {
		t.Error("ProjectID should be set")
	}
}

func TestTaskDetailJSON(t *testing.T) {
	jsonStr := `{
		"task": {
			"id": "task-1",
			"title": "Fix bug",
			"body": "Spec here",
			"status": "ready",
			"priority": 5,
			"attempt": 0,
			"max_attempts": 3
		},
		"context": [
			{
				"id": 1,
				"task_id": "task-1",
				"kind": "inherited",
				"content": "Found a bug in module X"
			}
		]
	}`
	var detail TaskDetail
	if err := json.Unmarshal([]byte(jsonStr), &detail); err != nil {
		t.Fatal(err)
	}
	if detail.Task.ID != "task-1" {
		t.Errorf("Task.ID = %q", detail.Task.ID)
	}
	if len(detail.Context) != 1 {
		t.Fatalf("expected 1 context entry, got %d", len(detail.Context))
	}
	if detail.Context[0].Kind != "inherited" {
		t.Errorf("Context kind = %q", detail.Context[0].Kind)
	}
}

func TestBridge_Run_NonExistentBinary(t *testing.T) {
	b := NewBridge("/nonexistent/binary", "")
	_, err := b.run("status")
	if err == nil {
		t.Error("should fail for nonexistent binary")
	}
}

func TestBridge_Status_NonExistentBinary(t *testing.T) {
	b := NewBridge("/nonexistent/binary", "")
	_, err := b.Status("project-1")
	if err == nil {
		t.Error("should fail for nonexistent binary")
	}
}

func TestBridge_Show_NonExistentBinary(t *testing.T) {
	b := NewBridge("/nonexistent/binary", "")
	_, err := b.Show("task-1")
	if err == nil {
		t.Error("should fail for nonexistent binary")
	}
}

func TestBridge_Tree_NonExistentBinary(t *testing.T) {
	b := NewBridge("/nonexistent/binary", "")
	_, err := b.Tree("project-1")
	if err == nil {
		t.Error("should fail for nonexistent binary")
	}
}

func TestBridge_Prompt_NonExistentBinary(t *testing.T) {
	b := NewBridge("/nonexistent/binary", "")
	_, err := b.PromptSingle("task-1")
	if err == nil {
		t.Error("should fail for nonexistent binary")
	}
}

// TestBridge_Status_MockScript tests Status parsing with a mock script.
func TestBridge_Status_MockScript(t *testing.T) {
	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "minuano")

	script := `#!/bin/bash
echo '[{"id":"task-1","title":"Fix bug","body":"","status":"ready","priority":5,"attempt":0,"max_attempts":3}]'
`
	os.WriteFile(scriptPath, []byte(script), 0755)

	b := NewBridge(scriptPath, "")
	tasks, err := b.Status("")
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].ID != "task-1" {
		t.Errorf("ID = %q", tasks[0].ID)
	}
	if tasks[0].Status != "ready" {
		t.Errorf("Status = %q", tasks[0].Status)
	}
}

// TestBridge_Show_MockScript tests Show parsing with a mock script.
func TestBridge_Show_MockScript(t *testing.T) {
	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "minuano")

	script := `#!/bin/bash
echo '{"task":{"id":"task-1","title":"Fix bug","body":"Spec here","status":"ready","priority":5,"attempt":0,"max_attempts":3},"context":[]}'
`
	os.WriteFile(scriptPath, []byte(script), 0755)

	b := NewBridge(scriptPath, "")
	detail, err := b.Show("task-1")
	if err != nil {
		t.Fatal(err)
	}
	if detail.Task.ID != "task-1" {
		t.Errorf("Task.ID = %q", detail.Task.ID)
	}
	if detail.Task.Body != "Spec here" {
		t.Errorf("Task.Body = %q", detail.Task.Body)
	}
}

// TestBridge_Tree_MockScript tests Tree output with a mock script.
func TestBridge_Tree_MockScript(t *testing.T) {
	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "minuano")

	script := `#!/bin/bash
echo "  ◎  task-1  Fix bug"
echo "    └── ○  task-2  Refactor"
`
	os.WriteFile(scriptPath, []byte(script), 0755)

	b := NewBridge(scriptPath, "")
	tree, err := b.Tree("project-1")
	if err != nil {
		t.Fatal(err)
	}
	if tree == "" {
		t.Error("tree should not be empty")
	}
	if !containsSubstr(tree, "task-1") {
		t.Error("tree should contain task-1")
	}
}

// TestBridge_DBFlag tests that --db flag is passed.
func TestBridge_DBFlag_MockScript(t *testing.T) {
	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "minuano")

	// Script echoes all args as JSON array
	script := `#!/bin/bash
echo "[\"$@\"]"
`
	os.WriteFile(scriptPath, []byte(script), 0755)

	b := NewBridge(scriptPath, "postgresql://localhost/test")
	out, err := b.run("status", "--json")
	if err != nil {
		t.Fatal(err)
	}
	if !containsSubstr(out, "--db") {
		t.Error("should include --db flag")
	}
	if !containsSubstr(out, "postgresql://localhost/test") {
		t.Error("should include DB connection string")
	}
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
