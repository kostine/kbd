package model

import (
	"testing"

	"github.com/kostine/kbd/internal/bd"
)

func sampleIssues() []bd.Issue {
	return []bd.Issue{
		{ID: "1", Title: "Fix login bug", Type: "Bug", Assignee: "Alice", Status: "Open"},
		{ID: "2", Title: "Add search feature", Type: "Feature", Assignee: "Bob", Status: "In Progress"},
		{ID: "3", Title: "Update README", Type: "Task", Assignee: "Alice", Status: "Closed"},
		{ID: "4", Title: "Fix search performance", Type: "Bug", Assignee: "Carol", Status: "Open"},
	}
}

func TestFilterIssues_EmptyFilters(t *testing.T) {
	issues := sampleIssues()
	got := FilterIssues(issues, Filters{})
	if len(got) != len(issues) {
		t.Fatalf("expected %d issues, got %d", len(issues), len(got))
	}
	for i := range issues {
		if got[i].ID != issues[i].ID {
			t.Errorf("issue[%d]: expected ID %q, got %q", i, issues[i].ID, got[i].ID)
		}
	}
}

func TestFilterIssues_TitleSubstringCaseInsensitive(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  []string // expected IDs
	}{
		{"lowercase match", "fix", []string{"1", "4"}},
		{"uppercase match", "FIX", []string{"1", "4"}},
		{"mixed case match", "Search", []string{"2", "4"}},
		{"exact substring", "README", []string{"3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterIssues(sampleIssues(), Filters{Title: tt.title})
			if len(got) != len(tt.want) {
				t.Fatalf("expected %d issues, got %d", len(tt.want), len(got))
			}
			for i, id := range tt.want {
				if got[i].ID != id {
					t.Errorf("result[%d]: expected ID %q, got %q", i, id, got[i].ID)
				}
			}
		})
	}
}

func TestFilterIssues_TypeExact(t *testing.T) {
	tests := []struct {
		name string
		typ  string
		want []string
	}{
		{"bug", "Bug", []string{"1", "4"}},
		{"bug lowercase", "bug", []string{"1", "4"}},
		{"feature", "Feature", []string{"2"}},
		{"task", "Task", []string{"3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterIssues(sampleIssues(), Filters{Type: tt.typ})
			if len(got) != len(tt.want) {
				t.Fatalf("expected %d issues, got %d", len(tt.want), len(got))
			}
			for i, id := range tt.want {
				if got[i].ID != id {
					t.Errorf("result[%d]: expected ID %q, got %q", i, id, got[i].ID)
				}
			}
		})
	}
}

func TestFilterIssues_AssigneeExact(t *testing.T) {
	tests := []struct {
		name     string
		assignee string
		want     []string
	}{
		{"Alice", "Alice", []string{"1", "3"}},
		{"alice lowercase", "alice", []string{"1", "3"}},
		{"Bob", "Bob", []string{"2"}},
		{"Carol", "Carol", []string{"4"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterIssues(sampleIssues(), Filters{Assignee: tt.assignee})
			if len(got) != len(tt.want) {
				t.Fatalf("expected %d issues, got %d", len(tt.want), len(got))
			}
			for i, id := range tt.want {
				if got[i].ID != id {
					t.Errorf("result[%d]: expected ID %q, got %q", i, id, got[i].ID)
				}
			}
		})
	}
}

func TestFilterIssues_StatusExact(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   []string
	}{
		{"Open", "Open", []string{"1", "4"}},
		{"open lowercase", "open", []string{"1", "4"}},
		{"In Progress", "In Progress", []string{"2"}},
		{"Closed", "Closed", []string{"3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterIssues(sampleIssues(), Filters{Status: tt.status})
			if len(got) != len(tt.want) {
				t.Fatalf("expected %d issues, got %d", len(tt.want), len(got))
			}
			for i, id := range tt.want {
				if got[i].ID != id {
					t.Errorf("result[%d]: expected ID %q, got %q", i, id, got[i].ID)
				}
			}
		})
	}
}

func TestFilterIssues_MultipleFilters(t *testing.T) {
	tests := []struct {
		name    string
		filters Filters
		want    []string
	}{
		{
			name:    "type and status",
			filters: Filters{Type: "Bug", Status: "Open"},
			want:    []string{"1", "4"},
		},
		{
			name:    "type and assignee",
			filters: Filters{Type: "Bug", Assignee: "Alice"},
			want:    []string{"1"},
		},
		{
			name:    "title and type",
			filters: Filters{Title: "search", Type: "Bug"},
			want:    []string{"4"},
		},
		{
			name:    "all filters",
			filters: Filters{Title: "fix", Type: "Bug", Assignee: "Carol", Status: "Open"},
			want:    []string{"4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterIssues(sampleIssues(), tt.filters)
			if len(got) != len(tt.want) {
				t.Fatalf("expected %d issues, got %d", len(tt.want), len(got))
			}
			for i, id := range tt.want {
				if got[i].ID != id {
					t.Errorf("result[%d]: expected ID %q, got %q", i, id, got[i].ID)
				}
			}
		})
	}
}

func TestFilterIssues_NoMatches(t *testing.T) {
	got := FilterIssues(sampleIssues(), Filters{Title: "nonexistent"})
	if len(got) != 0 {
		t.Fatalf("expected 0 issues, got %d", len(got))
	}
}

func TestFilterIssues_EmptyIssueList(t *testing.T) {
	got := FilterIssues(nil, Filters{Title: "anything"})
	if len(got) != 0 {
		t.Fatalf("expected 0 issues, got %d", len(got))
	}

	got = FilterIssues([]bd.Issue{}, Filters{})
	if len(got) != 0 {
		t.Fatalf("expected 0 issues, got %d", len(got))
	}
}
