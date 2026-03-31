package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Shortcut is a numbered quick-filter entry.
type Shortcut struct {
	Key      int    // 0-9
	Label    string // display label, e.g. "iouri" or "epic"
	Category string // "assignee", "type", or "status"
	Value    string // filter value
}

// ShortcutPanel renders numbered shortcuts in a 5-row × 2-column grid.
type ShortcutPanel struct {
	*tview.TextView
	theme     *Theme
	shortcuts [10]Shortcut // slots 0-9
	count     int
}

// NewShortcutPanel creates a shortcut panel.
func NewShortcutPanel(theme *Theme) *ShortcutPanel {
	p := &ShortcutPanel{
		TextView: tview.NewTextView(),
		theme:    theme,
	}
	p.SetDynamicColors(true)
	p.SetBackgroundColor(tcell.ColorDefault)
	p.SetWrap(false)
	return p
}

// ShortcutSource holds a value and its count for ranking.
type ShortcutSource struct {
	Category string
	Value    string
	Label    string
	Count    int
}

// BuildShortcuts computes the top-10 shortcuts from issue data counts
// plus fixed command shortcuts appended at the end.
// Groups are laid out in fixed order: Assignee, Type, Status, Commands.
// Within each group, entries are sorted by count descending.
// Each input slice should already be sorted by count descending.
func BuildShortcuts(assignees, types, statuses, commands []ShortcutSource) []Shortcut {
	var all []ShortcutSource
	all = append(all, assignees...)
	all = append(all, types...)
	all = append(all, statuses...)
	all = append(all, commands...)

	if len(all) > 10 {
		all = all[:10]
	}

	shortcuts := make([]Shortcut, len(all))
	for i, s := range all {
		shortcuts[i] = Shortcut{
			Key:      i,
			Label:    s.Label,
			Category: s.Category,
			Value:    s.Value,
		}
	}
	return shortcuts
}

// SetShortcuts updates the panel display.
func (p *ShortcutPanel) SetShortcuts(shortcuts []Shortcut) {
	p.count = len(shortcuts)
	for i := 0; i < 10; i++ {
		if i < len(shortcuts) {
			p.shortcuts[i] = shortcuts[i]
		} else {
			p.shortcuts[i] = Shortcut{}
		}
	}
	p.render()
}

// Get returns the shortcut for a given key (0-9), or nil if empty.
func (p *ShortcutPanel) Get(key int) *Shortcut {
	if key < 0 || key >= 10 || key >= p.count {
		return nil
	}
	s := p.shortcuts[key]
	return &s
}

var catColor = map[string]string{
	"assignee": "green",
	"type":     "aqua",
	"status":   "yellow",
	"context":  "purple",
	"command":  "gray",
}

func (p *ShortcutPanel) render() {
	p.Clear()

	// 5 rows × 2 columns: left col = 0-4, right col = 5-9
	for row := 0; row < 5; row++ {
		if row > 0 {
			fmt.Fprint(p, "\n")
		}
		p.writeSlot(row)
		fmt.Fprint(p, "  ")
		p.writeSlot(row + 5)
	}
}

func (p *ShortcutPanel) writeSlot(idx int) {
	if idx >= p.count {
		fmt.Fprintf(p, "%-22s", "")
		return
	}
	s := p.shortcuts[idx]
	color := catColor[s.Category]
	if color == "" {
		color = "white"
	}
	label := s.Label
	if len(label) > 18 {
		label = label[:17] + "…"
	}
	fmt.Fprintf(p, "[white::b]%d[gray]:[%s]%-19s", s.Key, color, label)
}
