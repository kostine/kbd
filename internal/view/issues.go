package view

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/kostine/kbd/internal/bd"
	"github.com/kostine/kbd/internal/model"
	tmux "github.com/kostine/kbd/internal/tmux"
	"github.com/kostine/kbd/internal/ui"
)

// RefreshInterval is how often the issue list auto-refreshes.
const RefreshInterval = 5 * time.Second

// Issues is the main issue list view.
type Issues struct {
	*tview.Flex
	table       *ui.Table
	theme       *ui.Theme
	client      *bd.Client
	issues      []bd.Issue
	filtered    []bd.Issue
	deltas      []model.RowKind
	filters     model.Filters
	childCounts map[string]bd.ChildCounts
	filterStack []string
	nav         Navigator
	cancel      context.CancelFunc
}

// Navigator is the interface for stack navigation.
type Navigator interface {
	Push(name string, view tview.Primitive, hints []ui.Hint)
	Pop()
	SetStatus(msg string)
	SetCounts(counts []ui.TypeStatusCounts)
	SetShortcuts(shortcuts []ui.Shortcut)
	App() *tview.Application
	Pages() ui.DialogPages
}

// NewIssues creates the issue list view.
func NewIssues(client *bd.Client, theme *ui.Theme, nav Navigator) *Issues {
	v := &Issues{
		Flex:   tview.NewFlex().SetDirection(tview.FlexRow),
		table:  ui.NewTable(theme),
		theme:  theme,
		client: client,
		nav:    nav,
	}

	v.table.SetTitle(" Issues ")
	v.table.SetTitleColor(theme.TitleColor)
	v.table.SetHeaders(model.IssueHeaders...)
	v.table.OnSelect(v.onSelect)

	v.table.KeyBindings().BindRune('x', "Close", v.closeIssue)
	v.table.KeyBindings().BindRune('o', "Reopen", v.reopenIssue)
	v.table.KeyBindings().BindRune('r', "Refresh", func() { v.Refresh() })
	v.table.KeyBindings().BindRune('h', "Help", v.showHelp)

	if tmux.InTmux() {
		v.table.KeyBindings().BindRune('w', "Work on", v.workOn)
		v.table.KeyBindings().BindRune('c', "Chat", v.chat)
	}

	v.AddItem(v.table, 0, 1, true)
	return v
}

// StartAutoRefresh begins the background refresh loop.
func (v *Issues) StartAutoRefresh() {
	if v.cancel != nil {
		v.cancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel
	go v.refreshLoop(ctx)
}

// StopAutoRefresh stops the background refresh loop.
func (v *Issues) StopAutoRefresh() {
	if v.cancel != nil {
		v.cancel()
		v.cancel = nil
	}
}

func (v *Issues) refreshLoop(ctx context.Context) {
	ticker := time.NewTicker(RefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			v.silentRefresh()
		}
	}
}

// silentRefresh reloads data without showing "Loading..." status.
func (v *Issues) silentRefresh() {
	issues, err := v.client.ListIssues()
	if err != nil {
		return // silently skip on error
	}
	childCounts, _ := v.client.EpicChildCounts()
	v.nav.App().QueueUpdateDraw(func() {
		oldFiltered := v.filtered
		v.issues = issues
		v.childCounts = childCounts
		v.applyFilterWithDeltas(oldFiltered)
	})
}

// Refresh reloads issues from bd (manual, shows status).
func (v *Issues) Refresh() {
	v.nav.SetStatus("Loading...")
	go func() {
		issues, err := v.client.ListIssues()
		if err != nil {
			v.nav.App().QueueUpdateDraw(func() {
				v.nav.SetStatus("Error: " + err.Error())
			})
			return
		}
		childCounts, _ := v.client.EpicChildCounts()
		v.nav.App().QueueUpdateDraw(func() {
			v.issues = issues
			v.childCounts = childCounts
			v.applyFilter()
			v.nav.SetStatus("")
			// Start auto-refresh after first successful load
			if v.cancel == nil {
				v.StartAutoRefresh()
			}
		})
	}()
}

// SetTitleFilter applies a title search filter (from /).
func (v *Issues) SetTitleFilter(filter string) {
	v.filters.Title = filter
	v.pushFilterStack("title")
	v.applyFilter()
	v.updateTitle()
}

// ApplyShortcut applies a filter from a shortcut.
func (v *Issues) ApplyShortcut(category, value string) {
	switch category {
	case "assignee":
		v.filters.Assignee = value
	case "type":
		v.filters.Type = value
	case "status":
		v.filters.Status = value
	}
	v.pushFilterStack(category)
	v.applyFilter()
	v.updateTitle()
}

// PopFilter removes the most recently applied filter layer.
// Returns true if a filter was removed.
func (v *Issues) PopFilter() bool {
	if len(v.filterStack) == 0 {
		return false
	}
	cat := v.filterStack[len(v.filterStack)-1]
	v.filterStack = v.filterStack[:len(v.filterStack)-1]
	switch cat {
	case "title":
		v.filters.Title = ""
	case "type":
		v.filters.Type = ""
	case "assignee":
		v.filters.Assignee = ""
	case "status":
		v.filters.Status = ""
	}
	v.applyFilter()
	v.updateTitle()
	return true
}

// SetTypeFilter applies a type filter (from :).
func (v *Issues) SetTypeFilter(typeFilter string) {
	if typeFilter == "all" || typeFilter == "issues" || typeFilter == "list" {
		typeFilter = ""
	}
	v.filters.Type = typeFilter
	v.pushFilterStack("type")
	v.applyFilter()
	v.updateTitle()
}

func (v *Issues) pushFilterStack(category string) {
	filtered := v.filterStack[:0]
	for _, c := range v.filterStack {
		if c != category {
			filtered = append(filtered, c)
		}
	}
	v.filterStack = append(filtered, category)
}

func (v *Issues) updateTitle() {
	title := " Issues"
	if v.filters.Type != "" {
		title = fmt.Sprintf(" %s", v.filters.Type)
	}
	if v.filters.Assignee != "" {
		title += fmt.Sprintf(" @%s", v.filters.Assignee)
	}
	if v.filters.Status != "" {
		title += fmt.Sprintf(" [%s]", v.filters.Status)
	}
	if v.filters.Title != "" {
		title += fmt.Sprintf(" /%s", v.filters.Title)
	}
	title += fmt.Sprintf(" (%d) ", len(v.filtered))
	v.table.SetTitle(title)
}

// applyFilter filters, sorts, and renders without delta tracking.
func (v *Issues) applyFilter() {
	v.filtered = model.FilterIssues(v.issues, v.filters)
	model.SortDefault(v.filtered, v.filters.Type)
	v.deltas = nil
	v.renderTable()
	v.updateCounts()
	v.updateShortcuts()
}

// applyFilterWithDeltas filters, sorts, and renders with delta tracking.
func (v *Issues) applyFilterWithDeltas(oldFiltered []bd.Issue) {
	v.filtered = model.FilterIssues(v.issues, v.filters)
	model.SortDefault(v.filtered, v.filters.Type)
	dr := model.ComputeDeltas(oldFiltered, v.filtered)
	v.deltas = dr.Kinds
	v.renderTable()
	v.updateTitle()
	v.updateCounts()
	v.updateShortcuts()
}

func (v *Issues) updateCounts() {
	type statusMap struct {
		order  []string
		counts map[string]int
	}
	types := make(map[string]*statusMap)
	var typeOrder []string

	for _, issue := range v.issues {
		t := issue.Type
		if t == "" {
			t = "other"
		}
		s := issue.Status
		if s == "" {
			s = "open"
		}
		sm, ok := types[t]
		if !ok {
			sm = &statusMap{counts: make(map[string]int)}
			types[t] = sm
			typeOrder = append(typeOrder, t)
		}
		if _, exists := sm.counts[s]; !exists {
			sm.order = append(sm.order, s)
		}
		sm.counts[s]++
	}

	counts := make([]ui.TypeStatusCounts, 0, len(typeOrder))
	for _, t := range typeOrder {
		sm := types[t]
		tc := ui.TypeStatusCounts{Type: t}
		for _, s := range sm.order {
			tc.Statuses = append(tc.Statuses, ui.StatusCount{Status: s, Count: sm.counts[s]})
		}
		counts = append(counts, tc)
	}
	v.nav.SetCounts(counts)
}

// RefreshShortcuts re-publishes the issues shortcuts to the header.
func (v *Issues) RefreshShortcuts() {
	v.updateShortcuts()
}

func (v *Issues) updateShortcuts() {
	assigneeCounts := countField(v.filtered, func(i bd.Issue) string { return i.Assignee })
	typeCounts := countField(v.filtered, func(i bd.Issue) string { return i.Type })
	statusCounts := countField(v.filtered, func(i bd.Issue) string { return i.Status })

	var assignees, types, statuses []ui.ShortcutSource
	for _, kv := range assigneeCounts {
		assignees = append(assignees, ui.ShortcutSource{
			Category: "assignee", Value: kv.value, Label: kv.value, Count: kv.count,
		})
	}
	for _, kv := range typeCounts {
		types = append(types, ui.ShortcutSource{
			Category: "type", Value: kv.value, Label: kv.value, Count: kv.count,
		})
	}
	for _, kv := range statusCounts {
		statuses = append(statuses, ui.ShortcutSource{
			Category: "status", Value: kv.value, Label: kv.value, Count: kv.count,
		})
	}

	commands := []ui.ShortcutSource{
		{Category: "command", Value: "context", Label: "context"},
	}
	shortcuts := ui.BuildShortcuts(assignees, types, statuses, commands)
	v.nav.SetShortcuts(shortcuts)
}

type fieldCount struct {
	value string
	count int
}

func countField(issues []bd.Issue, extract func(bd.Issue) string) []fieldCount {
	m := make(map[string]int)
	for _, issue := range issues {
		v := extract(issue)
		if v != "" {
			m[v]++
		}
	}
	result := make([]fieldCount, 0, len(m))
	for k, c := range m {
		result = append(result, fieldCount{value: k, count: c})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].count > result[j].count
	})
	return result
}

func (v *Issues) isEpicView() bool {
	return strings.EqualFold(v.filters.Type, "epic")
}

func (v *Issues) renderTable() {
	// Remember current selection to restore after re-render
	selectedID := v.selectedID()

	v.table.Clear()
	if v.isEpicView() {
		v.table.SetHeaders(model.EpicHeaders...)
	} else {
		v.table.SetHeaders(model.IssueHeaders...)
	}
	for i, issue := range v.filtered {
		var row []string
		if v.isEpicView() {
			row = model.EpicRow(issue, v.childCounts)
		} else {
			row = model.IssueRow(issue)
		}
		color := v.rowColor(i, issue)
		v.table.SetRowWithColor(i+1, color, row...)
	}

	// Restore selection by ID, or stay at top if not found
	if selectedID != "" {
		for i, issue := range v.filtered {
			if issue.ID == selectedID {
				v.table.Select(i+1, 0)
				return
			}
		}
	}
	v.table.ScrollToBeginning()
}

// selectedID returns the ID of the currently selected issue, or empty.
func (v *Issues) selectedID() string {
	row, _ := v.table.GetSelection()
	idx := row - 1
	if idx >= 0 && idx < len(v.filtered) {
		return v.filtered[idx].ID
	}
	return ""
}

// rowColor returns the color for a row, using delta color if available,
// otherwise falling back to status color.
func (v *Issues) rowColor(idx int, issue bd.Issue) tcell.Color {
	if v.deltas != nil && idx < len(v.deltas) {
		switch v.deltas[idx] {
		case model.RowAdded:
			return v.theme.DeltaAdd
		case model.RowUpdated:
			return v.theme.DeltaUpdate
		}
	}
	return v.theme.StatusColor(issue.Status)
}

func (v *Issues) showHelp() {
	help := NewHelp(v.theme, v.nav)
	v.nav.Push("help", help, help.Hints())
}

func (v *Issues) onSelect(row int) {
	idx := row - 1
	if idx < 0 || idx >= len(v.filtered) {
		return
	}
	issue := v.filtered[idx]
	detail := NewDetail(v.client, v.theme, v.nav, issue.ID)
	v.nav.Push("detail:"+issue.ID, detail, detail.Hints())
	detail.Refresh()
}

func (v *Issues) workOn() {
	row, _ := v.table.GetSelection()
	idx := row - 1
	if idx < 0 || idx >= len(v.filtered) {
		return
	}
	issue := v.filtered[idx]
	defaultMsg := fmt.Sprintf("lets work on %s", issue.ID)

	ui.ShowInput(v.nav.Pages(), v.nav.App(), "Send to Claude", defaultMsg,
		func(msg string) {
			if msg == "" {
				return
			}
			if err := tmux.SendToOtherPane(msg); err != nil {
				v.nav.SetStatus("Error: " + err.Error())
				return
			}
			v.nav.SetStatus("Sent: " + msg)
		},
		nil,
	)
}

func (v *Issues) chat() {
	ui.ShowInput(v.nav.Pages(), v.nav.App(), "Send to Claude", "",
		func(msg string) {
			if msg == "" {
				return
			}
			if err := tmux.SendToOtherPane(msg); err != nil {
				v.nav.SetStatus("Error: " + err.Error())
				return
			}
			v.nav.SetStatus("Sent: " + msg)
		},
		nil,
	)
}

func (v *Issues) closeIssue() {
	row, _ := v.table.GetSelection()
	idx := row - 1
	if idx < 0 || idx >= len(v.filtered) {
		return
	}
	issue := v.filtered[idx]
	go func() {
		err := v.client.CloseIssue(issue.ID)
		v.nav.App().QueueUpdateDraw(func() {
			if err != nil {
				v.nav.SetStatus("Error: " + err.Error())
				return
			}
			v.nav.SetStatus("Closed " + issue.ID)
			v.Refresh()
		})
	}()
}

func (v *Issues) reopenIssue() {
	row, _ := v.table.GetSelection()
	idx := row - 1
	if idx < 0 || idx >= len(v.filtered) {
		return
	}
	issue := v.filtered[idx]
	go func() {
		err := v.client.ReopenIssue(issue.ID)
		v.nav.App().QueueUpdateDraw(func() {
			if err != nil {
				v.nav.SetStatus("Error: " + err.Error())
				return
			}
			v.nav.SetStatus("Reopened " + issue.ID)
			v.Refresh()
		})
	}()
}

// Hints returns the key binding hints for the menu.
func (v *Issues) Hints() []ui.Hint {
	hints := []ui.Hint{
		{Key: "Enter", Description: "View"},
		{Key: "j/k", Description: "Navigate"},
		{Key: "g/G", Description: "Top/Bottom"},
		{Key: "/", Description: "Search"},
		{Key: ":", Description: "Type"},
		{Key: "0-9", Description: "Quick filter"},
	}
	if tmux.InTmux() {
		hints = append(hints, ui.Hint{Key: "w", Description: "Work on"})
		hints = append(hints, ui.Hint{Key: "c", Description: "Chat"})
	}
	hints = append(hints,
		ui.Hint{Key: "x", Description: "Close"},
		ui.Hint{Key: "o", Description: "Reopen"},
		ui.Hint{Key: "r", Description: "Refresh"},
		ui.Hint{Key: "h", Description: "Help"},
		ui.Hint{Key: "Esc", Description: "Pop filter"},
		ui.Hint{Key: "q", Description: "Quit"},
	)
	return hints
}

// compile-time check
var _ tview.Primitive = (*Issues)(nil)
