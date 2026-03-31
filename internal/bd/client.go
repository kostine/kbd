package bd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/kostine/kbd/internal/logging"
)

// Client wraps the bd CLI for all data operations.
type Client struct {
	// DBPath overrides the --db flag. Empty means auto-discover.
	DBPath string
}

// NewClient creates a bd client.
func NewClient(dbPath string) *Client {
	return &Client{DBPath: dbPath}
}

// run executes a bd command and returns stdout.
func (c *Client) run(args ...string) ([]byte, error) {
	if c.DBPath != "" {
		args = append([]string{"--db", c.DBPath}, args...)
	}
	args = append(args, "--json")

	cmd := exec.Command("bd", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		err := fmt.Errorf("bd %s: %s", strings.Join(args, " "), errMsg)
		logging.Error("%v", err)
		return nil, err
	}
	return stdout.Bytes(), nil
}

// runText executes a bd command and returns raw text stdout (no --json).
func (c *Client) runText(args ...string) (string, error) {
	if c.DBPath != "" {
		args = append([]string{"--db", c.DBPath}, args...)
	}

	cmd := exec.Command("bd", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return "", fmt.Errorf("bd %s: %s", strings.Join(args, " "), errMsg)
	}
	return stdout.String(), nil
}

// ListIssues returns issues matching the given filters.
func (c *Client) ListIssues(filters ...string) ([]Issue, error) {
	args := append([]string{"list", "--all", "--limit", "0"}, filters...)
	data, err := c.run(args...)
	if err != nil {
		return nil, err
	}
	var issues []Issue
	if err := json.Unmarshal(data, &issues); err != nil {
		return nil, fmt.Errorf("parse issues: %w", err)
	}
	return issues, nil
}

// ShowIssue returns detailed info for a single issue.
func (c *Client) ShowIssue(id string) (*Issue, error) {
	data, err := c.run("show", id)
	if err != nil {
		return nil, err
	}
	// bd show returns an array even for a single issue
	var issues []Issue
	if err := json.Unmarshal(data, &issues); err != nil {
		return nil, fmt.Errorf("parse issue: %w", err)
	}
	if len(issues) == 0 {
		return nil, fmt.Errorf("issue not found: %s", id)
	}
	return &issues[0], nil
}

// ShowIssueRaw returns the full bd show output as raw JSON for detail view.
func (c *Client) ShowIssueRaw(id string) (map[string]any, error) {
	data, err := c.run("show", id, "--long")
	if err != nil {
		return nil, err
	}
	// bd show returns an array even for a single issue
	var arr []map[string]any
	if err := json.Unmarshal(data, &arr); err != nil {
		return nil, fmt.Errorf("parse issue raw: %w", err)
	}
	if len(arr) == 0 {
		return nil, fmt.Errorf("issue not found: %s", id)
	}
	return arr[0], nil
}

// Children returns child issues of a parent.
func (c *Client) Children(id string) ([]Issue, error) {
	data, err := c.run("children", id)
	if err != nil {
		return nil, err
	}
	var issues []Issue
	if err := json.Unmarshal(data, &issues); err != nil {
		return nil, fmt.Errorf("parse children: %w", err)
	}
	return issues, nil
}

// Comments returns comments for an issue.
func (c *Client) Comments(id string) ([]Comment, error) {
	data, err := c.run("comments", id)
	if err != nil {
		return nil, err
	}
	var comments []Comment
	if err := json.Unmarshal(data, &comments); err != nil {
		return nil, fmt.Errorf("parse comments: %w", err)
	}
	return comments, nil
}

// SQL executes a raw SQL query and returns the result as a slice of maps.
func (c *Client) SQL(query string) ([]map[string]any, error) {
	data, err := c.run("sql", query)
	if err != nil {
		return nil, err
	}
	var rows []map[string]any
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, fmt.Errorf("parse sql: %w", err)
	}
	return rows, nil
}

// ChildCounts holds closed/total child issue counts for an epic.
type ChildCounts struct {
	Closed int
	Total  int
}

// EpicChildCounts returns a map of epic ID → child issue counts.
func (c *Client) EpicChildCounts() (map[string]ChildCounts, error) {
	query := `SELECT d.depends_on_id as epic_id, COUNT(*) as total, SUM(CASE WHEN i.status = "closed" THEN 1 ELSE 0 END) as closed FROM dependencies d JOIN issues i ON d.issue_id = i.id WHERE d.type = "parent-child" GROUP BY d.depends_on_id`
	rows, err := c.SQL(query)
	if err != nil {
		return nil, err
	}
	result := make(map[string]ChildCounts, len(rows))
	for _, row := range rows {
		id, _ := row["epic_id"].(string)
		total, _ := row["total"].(float64)
		closed, _ := row["closed"].(float64)
		result[id] = ChildCounts{Closed: int(closed), Total: int(total)}
	}
	return result, nil
}

// Graph returns the dependency graph for an issue as text.
func (c *Client) Graph(id string) (string, error) {
	return c.runText("graph", id)
}

// CloseIssue closes an issue.
func (c *Client) CloseIssue(id string) error {
	args := []string{"close", id}
	if c.DBPath != "" {
		args = append([]string{"--db", c.DBPath}, args...)
	}
	cmd := exec.Command("bd", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return fmt.Errorf("bd close: %s", errMsg)
	}
	return nil
}

// ReopenIssue reopens an issue.
func (c *Client) ReopenIssue(id string) error {
	args := []string{"reopen", id}
	if c.DBPath != "" {
		args = append([]string{"--db", c.DBPath}, args...)
	}
	cmd := exec.Command("bd", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return fmt.Errorf("bd reopen: %s", errMsg)
	}
	return nil
}
