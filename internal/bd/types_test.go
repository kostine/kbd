package bd

import (
	"encoding/json"
	"testing"
	"time"
)

func TestIssueUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    Issue
		wantErr bool
	}{
		{
			name: "full issue with all fields",
			json: `{
				"id":"bd-abc123",
				"title":"Fix login bug",
				"description":"Users cannot log in",
				"status":"open",
				"priority":1,
				"issue_type":"bug",
				"assignee":"iouri",
				"parent":"bd-epic1",
				"labels":["backend","auth"],
				"created_at":"2026-03-31T17:49:16Z",
				"updated_at":"2026-03-31T15:11:48Z",
				"closed_at":"2026-03-31T18:00:00Z",
				"defer_until":"2026-04-15T00:00:00Z",
				"due_date":"2026-04-30T00:00:00Z",
				"dependencies":[],
				"dependency_count":0,
				"dependent_count":0,
				"comment_count":2
			}`,
			want: Issue{
				ID:              "bd-abc123",
				Title:           "Fix login bug",
				Description:     "Users cannot log in",
				Status:          "open",
				Priority:        1,
				Type:            "bug",
				Assignee:        "iouri",
				Parent:          "bd-epic1",
				Labels:          []string{"backend", "auth"},
				CreatedAt:       time.Date(2026, 3, 31, 17, 49, 16, 0, time.UTC),
				UpdatedAt:       time.Date(2026, 3, 31, 15, 11, 48, 0, time.UTC),
				ClosedAt:        timePtr(time.Date(2026, 3, 31, 18, 0, 0, 0, time.UTC)),
				DeferUntil:      timePtr(time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC)),
				DueDate:         timePtr(time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)),
				DependencyCount: 0,
				DependentCount:  0,
				CommentCount:    2,
			},
		},
		{
			name: "missing optional time fields are nil",
			json: `{
				"id":"bd-xyz789",
				"title":"Add dark mode",
				"description":"",
				"status":"open",
				"priority":3,
				"issue_type":"feature",
				"assignee":"",
				"parent":"",
				"labels":[],
				"created_at":"2026-03-31T10:00:00Z",
				"updated_at":"2026-03-31T10:00:00Z",
				"dependency_count":1,
				"dependent_count":2,
				"comment_count":0
			}`,
			want: Issue{
				ID:              "bd-xyz789",
				Title:           "Add dark mode",
				Description:     "",
				Status:          "open",
				Priority:        3,
				Type:            "feature",
				Assignee:        "",
				Parent:          "",
				Labels:          []string{},
				CreatedAt:       time.Date(2026, 3, 31, 10, 0, 0, 0, time.UTC),
				UpdatedAt:       time.Date(2026, 3, 31, 10, 0, 0, 0, time.UTC),
				ClosedAt:        nil,
				DeferUntil:      nil,
				DueDate:         nil,
				DependencyCount: 1,
				DependentCount:  2,
				CommentCount:    0,
			},
		},
		{
			name: "null optional fields are nil",
			json: `{
				"id":"bd-null1",
				"title":"Test nulls",
				"description":"desc",
				"status":"closed",
				"priority":2,
				"issue_type":"task",
				"assignee":"iouri",
				"parent":"",
				"labels":null,
				"created_at":"2026-03-31T12:00:00Z",
				"updated_at":"2026-03-31T12:00:00Z",
				"closed_at":null,
				"defer_until":null,
				"due_date":null,
				"dependency_count":0,
				"dependent_count":0,
				"comment_count":0
			}`,
			want: Issue{
				ID:              "bd-null1",
				Title:           "Test nulls",
				Description:     "desc",
				Status:          "closed",
				Priority:        2,
				Type:            "task",
				Assignee:        "iouri",
				Parent:          "",
				Labels:          nil,
				CreatedAt:       time.Date(2026, 3, 31, 12, 0, 0, 0, time.UTC),
				UpdatedAt:       time.Date(2026, 3, 31, 12, 0, 0, 0, time.UTC),
				ClosedAt:        nil,
				DeferUntil:      nil,
				DueDate:         nil,
				DependencyCount: 0,
				DependentCount:  0,
				CommentCount:    0,
			},
		},
		{
			name: "high priority value",
			json: `{
				"id":"bd-pri5",
				"title":"Urgent fix",
				"description":"",
				"status":"open",
				"priority":5,
				"issue_type":"bug",
				"assignee":"",
				"parent":"",
				"labels":[],
				"created_at":"2026-03-31T12:00:00Z",
				"updated_at":"2026-03-31T12:00:00Z",
				"dependency_count":0,
				"dependent_count":0,
				"comment_count":0
			}`,
			want: Issue{
				ID:       "bd-pri5",
				Title:    "Urgent fix",
				Status:   "open",
				Priority: 5,
				Type:     "bug",
				Labels:   []string{},
				CreatedAt: time.Date(2026, 3, 31, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 3, 31, 12, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Issue
			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			assertIssueEqual(t, tt.want, got)
		})
	}
}

func TestIssueArrayUnmarshal(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		wantCount int
		wantIDs   []string
	}{
		{
			name:      "single item array (bd show format)",
			json:      `[{"id":"bd-abc123","title":"Fix login bug","description":"Users cannot log in","status":"open","priority":1,"issue_type":"bug","assignee":"iouri","created_at":"2026-03-31T17:49:16Z","updated_at":"2026-03-31T15:11:48Z","dependencies":[],"dependency_count":0,"dependent_count":0,"comment_count":2,"parent":"bd-epic1"}]`,
			wantCount: 1,
			wantIDs:   []string{"bd-abc123"},
		},
		{
			name: "multiple items array (bd list format)",
			json: `[
				{"id":"bd-001","title":"First","description":"","status":"open","priority":1,"issue_type":"task","assignee":"","created_at":"2026-03-31T10:00:00Z","updated_at":"2026-03-31T10:00:00Z","dependency_count":0,"dependent_count":0,"comment_count":0},
				{"id":"bd-002","title":"Second","description":"","status":"closed","priority":2,"issue_type":"bug","assignee":"iouri","created_at":"2026-03-31T11:00:00Z","updated_at":"2026-03-31T11:00:00Z","dependency_count":0,"dependent_count":0,"comment_count":1},
				{"id":"bd-003","title":"Third","description":"","status":"open","priority":3,"issue_type":"feature","assignee":"","created_at":"2026-03-31T12:00:00Z","updated_at":"2026-03-31T12:00:00Z","dependency_count":0,"dependent_count":0,"comment_count":0}
			]`,
			wantCount: 3,
			wantIDs:   []string{"bd-001", "bd-002", "bd-003"},
		},
		{
			name:      "empty array",
			json:      `[]`,
			wantCount: 0,
			wantIDs:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var issues []Issue
			if err := json.Unmarshal([]byte(tt.json), &issues); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			if len(issues) != tt.wantCount {
				t.Fatalf("got %d issues, want %d", len(issues), tt.wantCount)
			}
			for i, wantID := range tt.wantIDs {
				if issues[i].ID != wantID {
					t.Errorf("issues[%d].ID = %q, want %q", i, issues[i].ID, wantID)
				}
			}
		})
	}
}

func TestCommentUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    Comment
		wantErr bool
	}{
		{
			name: "standard comment",
			json: `{
				"id":"comment-001",
				"author":"iouri",
				"body":"This needs to be fixed ASAP",
				"created_at":"2026-03-31T18:00:00Z"
			}`,
			want: Comment{
				ID:        "comment-001",
				Author:    "iouri",
				Body:      "This needs to be fixed ASAP",
				CreatedAt: time.Date(2026, 3, 31, 18, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "comment with empty body",
			json: `{
				"id":"comment-002",
				"author":"bot",
				"body":"",
				"created_at":"2026-03-31T19:00:00Z"
			}`,
			want: Comment{
				ID:        "comment-002",
				Author:    "bot",
				Body:      "",
				CreatedAt: time.Date(2026, 3, 31, 19, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "comment with multiline body",
			json: `{
				"id":"comment-003",
				"author":"iouri",
				"body":"Line one\nLine two\nLine three",
				"created_at":"2026-03-31T20:00:00Z"
			}`,
			want: Comment{
				ID:        "comment-003",
				Author:    "iouri",
				Body:      "Line one\nLine two\nLine three",
				CreatedAt: time.Date(2026, 3, 31, 20, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Comment
			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("ID = %q, want %q", got.ID, tt.want.ID)
			}
			if got.Author != tt.want.Author {
				t.Errorf("Author = %q, want %q", got.Author, tt.want.Author)
			}
			if got.Body != tt.want.Body {
				t.Errorf("Body = %q, want %q", got.Body, tt.want.Body)
			}
			if !got.CreatedAt.Equal(tt.want.CreatedAt) {
				t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, tt.want.CreatedAt)
			}
		})
	}
}

func TestCommentArrayUnmarshal(t *testing.T) {
	data := `[
		{"id":"c1","author":"alice","body":"First comment","created_at":"2026-03-31T10:00:00Z"},
		{"id":"c2","author":"bob","body":"Second comment","created_at":"2026-03-31T11:00:00Z"}
	]`

	var comments []Comment
	if err := json.Unmarshal([]byte(data), &comments); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(comments) != 2 {
		t.Fatalf("got %d comments, want 2", len(comments))
	}
	if comments[0].ID != "c1" {
		t.Errorf("comments[0].ID = %q, want %q", comments[0].ID, "c1")
	}
	if comments[1].Author != "bob" {
		t.Errorf("comments[1].Author = %q, want %q", comments[1].Author, "bob")
	}
}

func TestStatusSummaryUnmarshal(t *testing.T) {
	data := `{"open":5,"in_progress":3,"blocked":1,"closed":10,"total":19}`

	var got StatusSummary
	if err := json.Unmarshal([]byte(data), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got.Open != 5 {
		t.Errorf("Open = %d, want 5", got.Open)
	}
	if got.InProgress != 3 {
		t.Errorf("InProgress = %d, want 3", got.InProgress)
	}
	if got.Blocked != 1 {
		t.Errorf("Blocked = %d, want 1", got.Blocked)
	}
	if got.Closed != 10 {
		t.Errorf("Closed = %d, want 10", got.Closed)
	}
	if got.Total != 19 {
		t.Errorf("Total = %d, want 19", got.Total)
	}
}

// helpers

func timePtr(t time.Time) *time.Time {
	return &t
}

func assertIssueEqual(t *testing.T, want, got Issue) {
	t.Helper()
	if got.ID != want.ID {
		t.Errorf("ID = %q, want %q", got.ID, want.ID)
	}
	if got.Title != want.Title {
		t.Errorf("Title = %q, want %q", got.Title, want.Title)
	}
	if got.Description != want.Description {
		t.Errorf("Description = %q, want %q", got.Description, want.Description)
	}
	if got.Status != want.Status {
		t.Errorf("Status = %q, want %q", got.Status, want.Status)
	}
	if got.Priority != want.Priority {
		t.Errorf("Priority = %d, want %d", got.Priority, want.Priority)
	}
	if got.Type != want.Type {
		t.Errorf("Type = %q, want %q", got.Type, want.Type)
	}
	if got.Assignee != want.Assignee {
		t.Errorf("Assignee = %q, want %q", got.Assignee, want.Assignee)
	}
	if got.Parent != want.Parent {
		t.Errorf("Parent = %q, want %q", got.Parent, want.Parent)
	}
	if got.DependencyCount != want.DependencyCount {
		t.Errorf("DependencyCount = %d, want %d", got.DependencyCount, want.DependencyCount)
	}
	if got.DependentCount != want.DependentCount {
		t.Errorf("DependentCount = %d, want %d", got.DependentCount, want.DependentCount)
	}
	if got.CommentCount != want.CommentCount {
		t.Errorf("CommentCount = %d, want %d", got.CommentCount, want.CommentCount)
	}
	if !got.CreatedAt.Equal(want.CreatedAt) {
		t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, want.CreatedAt)
	}
	if !got.UpdatedAt.Equal(want.UpdatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, want.UpdatedAt)
	}

	// Compare optional time pointers
	assertTimeEqual(t, "ClosedAt", want.ClosedAt, got.ClosedAt)
	assertTimeEqual(t, "DeferUntil", want.DeferUntil, got.DeferUntil)
	assertTimeEqual(t, "DueDate", want.DueDate, got.DueDate)

	// Compare labels
	if len(got.Labels) != len(want.Labels) {
		t.Errorf("Labels length = %d, want %d", len(got.Labels), len(want.Labels))
	} else {
		for i := range want.Labels {
			if got.Labels[i] != want.Labels[i] {
				t.Errorf("Labels[%d] = %q, want %q", i, got.Labels[i], want.Labels[i])
			}
		}
	}
}

func assertTimeEqual(t *testing.T, field string, want, got *time.Time) {
	t.Helper()
	if want == nil && got == nil {
		return
	}
	if want == nil && got != nil {
		t.Errorf("%s = %v, want nil", field, got)
		return
	}
	if want != nil && got == nil {
		t.Errorf("%s = nil, want %v", field, want)
		return
	}
	if !got.Equal(*want) {
		t.Errorf("%s = %v, want %v", field, got, want)
	}
}
