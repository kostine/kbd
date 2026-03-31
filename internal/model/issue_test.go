package model

import (
	"testing"
	"time"

	"github.com/kostine/kbd/internal/bd"
)

func TestIssueRow_ColumnCountAndValues(t *testing.T) {
	issue := bd.Issue{
		ID:              "ABC-123",
		Type:            "Bug",
		Status:          "Open",
		Title:           "Fix the thing",
		Assignee:        "Alice",
		DependencyCount: 3,
		CreatedAt:       time.Now().Add(-48 * time.Hour),
	}

	row := IssueRow(issue)

	if len(row) != len(IssueHeaders) {
		t.Fatalf("expected %d columns, got %d", len(IssueHeaders), len(row))
	}

	// Check fixed fields (not age, which is time-dependent).
	checks := []struct {
		index int
		name  string
		want  string
	}{
		{0, "ID", "ABC-123"},
		{1, "Type", "Bug"},
		{2, "Priority", "P0"},
		{3, "Status", "Open"},
		{4, "Title", "Fix the thing"},
		{5, "Assignee", "Alice"},
		{6, "Deps", "3"},
	}
	for _, c := range checks {
		if row[c.index] != c.want {
			t.Errorf("column %s (index %d): expected %q, got %q", c.name, c.index, c.want, row[c.index])
		}
	}

	// Age column should be non-empty for a non-zero CreatedAt.
	if row[7] == "" {
		t.Error("expected non-empty age column")
	}
}

func TestIssueRow_ZeroDependencies(t *testing.T) {
	issue := bd.Issue{
		ID:              "X-1",
		DependencyCount: 0,
	}
	row := IssueRow(issue)
	if row[6] != "" {
		t.Errorf("expected empty deps column, got %q", row[6])
	}
}

func TestIssueRow_ZeroCreatedAt(t *testing.T) {
	issue := bd.Issue{ID: "X-2"}
	row := IssueRow(issue)
	if row[7] != "" {
		t.Errorf("expected empty age column for zero time, got %q", row[7])
	}
}

func TestFormatAge(t *testing.T) {
	tests := []struct {
		name string
		ago  time.Duration
		want string
	}{
		{"30 seconds", 30 * time.Second, "30s"},
		{"5 minutes", 5 * time.Minute, "5m"},
		{"90 minutes", 90 * time.Minute, "1h"},
		{"3 hours", 3 * time.Hour, "3h"},
		{"23 hours", 23 * time.Hour, "23h"},
		{"2 days", 48 * time.Hour, "2d"},
		{"15 days", 15 * 24 * time.Hour, "15d"},
		{"45 days", 45 * 24 * time.Hour, "1mo"},
		{"90 days", 90 * 24 * time.Hour, "3mo"},
		{"400 days", 400 * 24 * time.Hour, "1y"},
		{"730 days", 730 * 24 * time.Hour, "2y"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := time.Now().Add(-tt.ago)
			got := FormatAge(input)
			if got != tt.want {
				t.Errorf("FormatAge(%v ago) = %q, want %q", tt.ago, got, tt.want)
			}
		})
	}
}

func TestFormatAge_ZeroTime(t *testing.T) {
	got := FormatAge(time.Time{})
	if got != "" {
		t.Errorf("FormatAge(zero) = %q, want empty string", got)
	}
}

func TestEpicRow_ColumnCount(t *testing.T) {
	issue := bd.Issue{ID: "bd-1", Type: "epic", Status: "open"}
	row := EpicRow(issue, nil)
	if len(row) != len(EpicHeaders) {
		t.Fatalf("row has %d columns, EpicHeaders has %d", len(row), len(EpicHeaders))
	}
}

func TestEpicRow_WithChildCounts(t *testing.T) {
	issue := bd.Issue{ID: "bd-epic1", Type: "epic", Status: "open", Title: "Big project", Priority: 1}
	counts := map[string]bd.ChildCounts{
		"bd-epic1": {Closed: 5, Total: 12},
	}
	row := EpicRow(issue, counts)
	if row[2] != "P1" {
		t.Errorf("PRI column = %q, want 'P1'", row[2])
	}
	// ISSUES column is at index 4
	if row[4] != "5/12" {
		t.Errorf("ISSUES column = %q, want '5/12'", row[4])
	}
}

func TestEpicRow_NoChildCounts(t *testing.T) {
	issue := bd.Issue{ID: "bd-epic2", Type: "epic", Status: "open"}
	row := EpicRow(issue, map[string]bd.ChildCounts{})
	if row[4] != "" {
		t.Errorf("ISSUES column = %q, want empty", row[4])
	}
}

func TestEpicRow_NilCounts(t *testing.T) {
	issue := bd.Issue{ID: "bd-epic3", Type: "epic", Status: "open"}
	row := EpicRow(issue, nil)
	if row[4] != "" {
		t.Errorf("ISSUES column = %q, want empty", row[4])
	}
}
