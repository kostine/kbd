package bd

import "time"

// Issue represents a bd issue as returned by bd --json commands.
type Issue struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Type            string     `json:"issue_type"`
	Status          string     `json:"status"`
	Assignee        string     `json:"assignee"`
	Description     string     `json:"description"`
	Parent          string     `json:"parent"`
	Priority        int        `json:"priority"`
	Labels          []string   `json:"labels"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ClosedAt        *time.Time `json:"closed_at"`
	DeferUntil      *time.Time `json:"defer_until"`
	DueDate         *time.Time `json:"due_date"`
	DependencyCount int        `json:"dependency_count"`
	DependentCount  int        `json:"dependent_count"`
	CommentCount    int        `json:"comment_count"`
}

// Comment represents a comment on an issue.
type Comment struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// StatusSummary represents the output of bd status --json.
type StatusSummary struct {
	Open       int `json:"open"`
	InProgress int `json:"in_progress"`
	Blocked    int `json:"blocked"`
	Closed     int `json:"closed"`
	Total      int `json:"total"`
}
