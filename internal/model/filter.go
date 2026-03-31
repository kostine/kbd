package model

import (
	"strings"

	"github.com/kostine/kbd/internal/bd"
)

// Filters holds all active filter criteria.
type Filters struct {
	Title    string
	Type     string
	Assignee string
	Status   string
}

// IsEmpty returns true if no filters are set.
func (f Filters) IsEmpty() bool {
	return f.Title == "" && f.Type == "" && f.Assignee == "" && f.Status == ""
}

// FilterIssues applies all active filters to the issue list.
func FilterIssues(issues []bd.Issue, f Filters) []bd.Issue {
	if f.IsEmpty() {
		return issues
	}
	title := strings.ToLower(f.Title)
	typ := strings.ToLower(f.Type)
	assignee := strings.ToLower(f.Assignee)
	status := strings.ToLower(f.Status)

	var result []bd.Issue
	for _, issue := range issues {
		if typ != "" && strings.ToLower(issue.Type) != typ {
			continue
		}
		if assignee != "" && strings.ToLower(issue.Assignee) != assignee {
			continue
		}
		if status != "" && strings.ToLower(issue.Status) != status {
			continue
		}
		if title != "" && !strings.Contains(strings.ToLower(issue.Title), title) {
			continue
		}
		result = append(result, issue)
	}
	return result
}
