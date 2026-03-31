package model

import "github.com/kostine/kbd/internal/bd"

// RowKind indicates what happened to a row between refreshes.
type RowKind int

const (
	RowUnchanged RowKind = iota
	RowAdded
	RowUpdated
	RowDeleted
)

// DeltaResult holds the filtered issues with per-row change tracking.
type DeltaResult struct {
	Issues []bd.Issue
	Kinds  []RowKind // parallel to Issues
}

// ComputeDeltas compares old and new issue lists (by ID) and returns
// the new list annotated with change kinds.
func ComputeDeltas(oldIssues, newIssues []bd.Issue) DeltaResult {
	oldMap := make(map[string]bd.Issue, len(oldIssues))
	for _, issue := range oldIssues {
		oldMap[issue.ID] = issue
	}

	result := DeltaResult{
		Issues: newIssues,
		Kinds:  make([]RowKind, len(newIssues)),
	}

	for i, issue := range newIssues {
		old, existed := oldMap[issue.ID]
		if !existed {
			result.Kinds[i] = RowAdded
		} else if issueChanged(old, issue) {
			result.Kinds[i] = RowUpdated
		} else {
			result.Kinds[i] = RowUnchanged
		}
	}

	return result
}

// issueChanged returns true if any visible field changed.
// Ignores time fields (UpdatedAt) to avoid noise.
func issueChanged(old, new bd.Issue) bool {
	return old.Title != new.Title ||
		old.Status != new.Status ||
		old.Type != new.Type ||
		old.Assignee != new.Assignee ||
		old.Priority != new.Priority ||
		old.Parent != new.Parent ||
		old.DependencyCount != new.DependencyCount ||
		old.DependentCount != new.DependentCount ||
		old.CommentCount != new.CommentCount
}
