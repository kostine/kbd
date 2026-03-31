package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TypeStatusCounts holds per-status counts for a single issue type.
type TypeStatusCounts struct {
	Type     string
	Statuses []StatusCount // in discovery order
}

// StatusCount is a status name and its count.
type StatusCount struct {
	Status string
	Count  int
}

// StatusBar renders a summary grid: one line per type, counts per status.
type StatusBar struct {
	*tview.TextView
	theme *Theme
}

// NewStatusBar creates a new status bar.
func NewStatusBar(theme *Theme) *StatusBar {
	s := &StatusBar{
		TextView: tview.NewTextView(),
		theme:    theme,
	}
	s.SetDynamicColors(true)
	s.SetBackgroundColor(tcell.ColorDefault)
	s.SetWrap(false)
	return s
}

// statusColors maps known statuses to tview color tags.
var statusColors = map[string]string{
	"open":        "green",
	"in_progress": "yellow",
	"blocked":     "red",
	"closed":      "gray",
	"deferred":    "blue",
	"pinned":      "purple",
	"hooked":      "yellow",
}

func colorFor(status string) string {
	if c, ok := statusColors[status]; ok {
		return c
	}
	return "white"
}

// SetCounts updates the status bar. One line per type.
func (s *StatusBar) SetCounts(counts []TypeStatusCounts) {
	s.Clear()
	for i, tc := range counts {
		if i > 0 {
			fmt.Fprint(s, "\n")
		}
		fmt.Fprintf(s, "[aqua::b]%-12s[white]", tc.Type)
		for j, sc := range tc.Statuses {
			if j > 0 {
				fmt.Fprint(s, " ")
			}
			fmt.Fprintf(s, "[%s]%s:%d", colorFor(sc.Status), sc.Status, sc.Count)
		}
	}
}

// LineCount returns how many lines the data needs.
func LineCount(counts []TypeStatusCounts) int {
	if len(counts) == 0 {
		return 1
	}
	return len(counts)
}
