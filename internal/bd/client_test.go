package bd

import (
	"encoding/json"
	"testing"
)

// These tests verify the JSON parsing logic used by each Client method.
// Since the Client methods all follow the pattern of running bd CLI then
// json.Unmarshal into the target type, we test the parsing directly to
// avoid requiring the bd CLI binary in test environments.

func TestParseListIssues(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		wantCount int
		wantFirst string // ID of first issue
		wantErr   bool
	}{
		{
			name:      "multiple issues",
			json:      `[{"id":"bd-001","title":"First","description":"","status":"open","priority":1,"issue_type":"task","assignee":"","created_at":"2026-03-31T10:00:00Z","updated_at":"2026-03-31T10:00:00Z","dependency_count":0,"dependent_count":0,"comment_count":0},{"id":"bd-002","title":"Second","description":"","status":"open","priority":2,"issue_type":"bug","assignee":"iouri","created_at":"2026-03-31T11:00:00Z","updated_at":"2026-03-31T11:00:00Z","dependency_count":0,"dependent_count":0,"comment_count":0}]`,
			wantCount: 2,
			wantFirst: "bd-001",
		},
		{
			name:      "empty result",
			json:      `[]`,
			wantCount: 0,
		},
		{
			name:    "invalid json",
			json:    `not json`,
			wantErr: true,
		},
		{
			name:    "object instead of array",
			json:    `{"id":"bd-001"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var issues []Issue
			err := json.Unmarshal([]byte(tt.json), &issues)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(issues) != tt.wantCount {
				t.Fatalf("got %d issues, want %d", len(issues), tt.wantCount)
			}
			if tt.wantFirst != "" && issues[0].ID != tt.wantFirst {
				t.Errorf("first issue ID = %q, want %q", issues[0].ID, tt.wantFirst)
			}
		})
	}
}

func TestParseShowIssue(t *testing.T) {
	// bd show returns a single-element array
	tests := []struct {
		name    string
		json    string
		wantID  string
		wantErr bool
	}{
		{
			name:   "single issue in array",
			json:   `[{"id":"bd-abc123","title":"Fix login bug","description":"Users cannot log in","status":"open","priority":1,"issue_type":"bug","assignee":"iouri","created_at":"2026-03-31T17:49:16Z","updated_at":"2026-03-31T15:11:48Z","dependencies":[],"dependency_count":0,"dependent_count":0,"comment_count":2,"parent":"bd-epic1"}]`,
			wantID: "bd-abc123",
		},
		{
			name:    "empty array means not found",
			json:    `[]`,
			wantErr: true,
		},
		{
			name:    "invalid json",
			json:    `{broken`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Replicate ShowIssue parsing logic
			var issues []Issue
			err := json.Unmarshal([]byte(tt.json), &issues)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("Unmarshal() unexpected error: %v", err)
				}
				return
			}
			if len(issues) == 0 {
				if !tt.wantErr {
					t.Fatal("got empty array, expected issue")
				}
				return
			}
			if tt.wantErr {
				t.Fatal("expected error, got none")
			}
			if issues[0].ID != tt.wantID {
				t.Errorf("ID = %q, want %q", issues[0].ID, tt.wantID)
			}
		})
	}
}

func TestParseShowIssueRaw(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantKey string // a key to check exists in the map
		wantErr bool
	}{
		{
			name:    "raw issue with extra fields",
			json:    `[{"id":"bd-abc123","title":"Fix login bug","custom_field":"extra_data","nested":{"key":"value"}}]`,
			wantKey: "custom_field",
		},
		{
			name:    "empty array means not found",
			json:    `[]`,
			wantErr: true,
		},
		{
			name:    "invalid json",
			json:    `[{broken}]`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Replicate ShowIssueRaw parsing logic
			var arr []map[string]any
			err := json.Unmarshal([]byte(tt.json), &arr)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("Unmarshal() unexpected error: %v", err)
				}
				return
			}
			if len(arr) == 0 {
				if !tt.wantErr {
					t.Fatal("got empty array, expected data")
				}
				return
			}
			if tt.wantErr {
				t.Fatal("expected error, got none")
			}
			if _, ok := arr[0][tt.wantKey]; !ok {
				t.Errorf("key %q not found in parsed map", tt.wantKey)
			}
		})
	}
}

func TestParseChildren(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "two children",
			json:      `[{"id":"bd-child1","title":"Subtask 1","description":"","status":"open","priority":1,"issue_type":"task","assignee":"","parent":"bd-epic1","created_at":"2026-03-31T10:00:00Z","updated_at":"2026-03-31T10:00:00Z","dependency_count":0,"dependent_count":0,"comment_count":0},{"id":"bd-child2","title":"Subtask 2","description":"","status":"closed","priority":2,"issue_type":"task","assignee":"iouri","parent":"bd-epic1","created_at":"2026-03-31T11:00:00Z","updated_at":"2026-03-31T11:00:00Z","dependency_count":0,"dependent_count":0,"comment_count":0}]`,
			wantCount: 2,
		},
		{
			name:      "no children",
			json:      `[]`,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var issues []Issue
			err := json.Unmarshal([]byte(tt.json), &issues)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(issues) != tt.wantCount {
				t.Fatalf("got %d children, want %d", len(issues), tt.wantCount)
			}
			// Verify parent field is populated on children
			for i, issue := range issues {
				if issue.Parent == "" {
					t.Errorf("issues[%d].Parent is empty, expected parent ID", i)
				}
			}
		})
	}
}

func TestParseComments(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "two comments",
			json:      `[{"id":"c1","author":"iouri","body":"Looking into this","created_at":"2026-03-31T17:50:00Z"},{"id":"c2","author":"bot","body":"Automated update: status changed","created_at":"2026-03-31T18:00:00Z"}]`,
			wantCount: 2,
		},
		{
			name:      "no comments",
			json:      `[]`,
			wantCount: 0,
		},
		{
			name:    "invalid json",
			json:    `{"not":"an array"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var comments []Comment
			err := json.Unmarshal([]byte(tt.json), &comments)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(comments) != tt.wantCount {
				t.Fatalf("got %d comments, want %d", len(comments), tt.wantCount)
			}
		})
	}
}

func TestParseSQL(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "sql result rows",
			json:      `[{"id":"bd-001","title":"First","status":"open"},{"id":"bd-002","title":"Second","status":"closed"}]`,
			wantCount: 2,
		},
		{
			name:      "empty result",
			json:      `[]`,
			wantCount: 0,
		},
		{
			name:      "aggregation result",
			json:      `[{"count":42}]`,
			wantCount: 1,
		},
		{
			name:    "invalid json",
			json:    `not json at all`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rows []map[string]any
			err := json.Unmarshal([]byte(tt.json), &rows)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(rows) != tt.wantCount {
				t.Fatalf("got %d rows, want %d", len(rows), tt.wantCount)
			}
		})
	}
}

func TestParseSQLFieldAccess(t *testing.T) {
	data := `[{"id":"bd-001","title":"Test","priority":1,"count":5}]`

	var rows []map[string]any
	if err := json.Unmarshal([]byte(data), &rows); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(rows))
	}

	row := rows[0]

	// String field
	if id, ok := row["id"].(string); !ok || id != "bd-001" {
		t.Errorf("id = %v, want %q", row["id"], "bd-001")
	}

	// Numeric field (json.Unmarshal into any produces float64)
	if count, ok := row["count"].(float64); !ok || count != 5 {
		t.Errorf("count = %v, want 5", row["count"])
	}

	// Numeric priority
	if pri, ok := row["priority"].(float64); !ok || pri != 1 {
		t.Errorf("priority = %v, want 1", row["priority"])
	}
}

func TestRealisticBdOutput(t *testing.T) {
	// This is the exact sample JSON from bd, testing full round-trip parsing
	// as used by ShowIssue (array with single element).
	raw := `[{"id":"bd-abc123","title":"Fix login bug","description":"Users cannot log in","status":"open","priority":1,"issue_type":"bug","assignee":"iouri","created_at":"2026-03-31T17:49:16Z","updated_at":"2026-03-31T15:11:48Z","dependencies":[],"dependency_count":0,"dependent_count":0,"comment_count":2,"parent":"bd-epic1"}]`

	var issues []Issue
	if err := json.Unmarshal([]byte(raw), &issues); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("got %d issues, want 1", len(issues))
	}

	issue := issues[0]
	if issue.ID != "bd-abc123" {
		t.Errorf("ID = %q, want %q", issue.ID, "bd-abc123")
	}
	if issue.Title != "Fix login bug" {
		t.Errorf("Title = %q, want %q", issue.Title, "Fix login bug")
	}
	if issue.Description != "Users cannot log in" {
		t.Errorf("Description = %q, want %q", issue.Description, "Users cannot log in")
	}
	if issue.Status != "open" {
		t.Errorf("Status = %q, want %q", issue.Status, "open")
	}
	if issue.Priority != 1 {
		t.Errorf("Priority = %d, want 1", issue.Priority)
	}
	if issue.Type != "bug" {
		t.Errorf("Type = %q, want %q", issue.Type, "bug")
	}
	if issue.Assignee != "iouri" {
		t.Errorf("Assignee = %q, want %q", issue.Assignee, "iouri")
	}
	if issue.Parent != "bd-epic1" {
		t.Errorf("Parent = %q, want %q", issue.Parent, "bd-epic1")
	}
	if issue.CommentCount != 2 {
		t.Errorf("CommentCount = %d, want 2", issue.CommentCount)
	}
	if issue.DependencyCount != 0 {
		t.Errorf("DependencyCount = %d, want 0", issue.DependencyCount)
	}
	if issue.ClosedAt != nil {
		t.Errorf("ClosedAt = %v, want nil", issue.ClosedAt)
	}
	if issue.DeferUntil != nil {
		t.Errorf("DeferUntil = %v, want nil", issue.DeferUntil)
	}
	if issue.DueDate != nil {
		t.Errorf("DueDate = %v, want nil", issue.DueDate)
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name   string
		dbPath string
	}{
		{name: "empty path", dbPath: ""},
		{name: "with path", dbPath: "/tmp/test.db"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.dbPath)
			if c == nil {
				t.Fatal("NewClient returned nil")
			}
			if c.DBPath != tt.dbPath {
				t.Errorf("DBPath = %q, want %q", c.DBPath, tt.dbPath)
			}
		})
	}
}
