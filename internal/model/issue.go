package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/kostine/kbd/internal/bd"
)

// IssueHeaders are the default columns for the issue list.
var IssueHeaders = []string{"ID", "TYPE", "PRI", "STATUS", "TITLE", "ASSIGNEE", "DEPS", "AGE"}

// EpicHeaders includes an ISSUES column between STATUS and TITLE.
var EpicHeaders = []string{"ID", "TYPE", "PRI", "STATUS", "ISSUES", "TITLE", "ASSIGNEE", "DEPS", "AGE"}

// DisplayStatus formats a status for display, replacing _ with -.
func DisplayStatus(s string) string {
	return strings.ReplaceAll(s, "_", "-")
}

// FormatPriority formats a priority int as P0-P4.
func FormatPriority(p int) string {
	return fmt.Sprintf("P%d", p)
}

// IssueRow converts an Issue into a string slice for table display.
func IssueRow(issue bd.Issue) []string {
	deps := ""
	if issue.DependencyCount > 0 {
		deps = fmt.Sprintf("%d", issue.DependencyCount)
	}
	return []string{
		issue.ID,
		issue.Type,
		FormatPriority(issue.Priority),
		DisplayStatus(issue.Status),
		issue.Title,
		issue.Assignee,
		deps,
		FormatAge(issue.CreatedAt),
	}
}

// EpicRow converts an Issue into a string slice with the ISSUES column.
func EpicRow(issue bd.Issue, childCounts map[string]bd.ChildCounts) []string {
	deps := ""
	if issue.DependencyCount > 0 {
		deps = fmt.Sprintf("%d", issue.DependencyCount)
	}
	issues := ""
	if cc, ok := childCounts[issue.ID]; ok {
		issues = fmt.Sprintf("%d/%d", cc.Closed, cc.Total)
	}
	return []string{
		issue.ID,
		issue.Type,
		FormatPriority(issue.Priority),
		DisplayStatus(issue.Status),
		issues,
		issue.Title,
		issue.Assignee,
		deps,
		FormatAge(issue.CreatedAt),
	}
}

// FormatAge converts a timestamp to a human-readable age string.
func FormatAge(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	case d < 365*24*time.Hour:
		return fmt.Sprintf("%dmo", int(d.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%dy", int(d.Hours()/(24*365)))
	}
}
