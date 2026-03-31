package view

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/kostine/kbd/internal/bd"
	"github.com/kostine/kbd/internal/ui"
)

// Detail shows full information for a single issue.
type Detail struct {
	*tview.Flex
	text   *tview.TextView
	theme  *ui.Theme
	client *bd.Client
	nav    Navigator
	id     string
}

// NewDetail creates a detail view for the given issue ID.
func NewDetail(client *bd.Client, theme *ui.Theme, nav Navigator, id string) *Detail {
	d := &Detail{
		Flex:   tview.NewFlex().SetDirection(tview.FlexRow),
		text:   tview.NewTextView(),
		theme:  theme,
		client: client,
		nav:    nav,
		id:     id,
	}

	d.text.SetDynamicColors(true)
	d.text.SetScrollable(true)
	d.text.SetWrap(true)
	d.text.SetBorder(true)
	d.text.SetBorderColor(theme.BorderColor)
	d.text.SetTitle(fmt.Sprintf(" %s ", id))
	d.text.SetTitleColor(theme.TitleColor)
	d.text.SetInputCapture(d.keyboard)

	d.AddItem(d.text, 0, 1, true)
	return d
}

// Refresh loads the issue detail from bd.
func (d *Detail) Refresh() {
	d.nav.SetStatus("Loading " + d.id + "...")
	go func() {
		raw, err := d.client.ShowIssueRaw(d.id)
		d.nav.App().QueueUpdateDraw(func() {
			if err != nil {
				d.nav.SetStatus("Error: " + err.Error())
				d.text.SetText(fmt.Sprintf("[red]Error: %s", err.Error()))
				return
			}
			d.renderDetail(raw)
			d.nav.SetStatus("")
		})
	}()
}

func (d *Detail) renderDetail(raw map[string]any) {
	var b strings.Builder

	writeField := func(label string, value any) {
		if value == nil || value == "" {
			return
		}
		fmt.Fprintf(&b, "[aqua]%-14s[white] %v\n", label+":", value)
	}

	writeField("ID", raw["id"])
	writeField("Title", raw["title"])
	writeField("Type", raw["type"])
	writeField("Status", raw["status"])
	writeField("Priority", raw["priority"])
	writeField("Assignee", raw["assignee"])
	writeField("Parent", raw["parent"])

	if labels, ok := raw["labels"].([]any); ok && len(labels) > 0 {
		strs := make([]string, len(labels))
		for i, l := range labels {
			strs[i] = fmt.Sprintf("%v", l)
		}
		writeField("Labels", strings.Join(strs, ", "))
	}

	if created, ok := raw["created_at"].(string); ok {
		writeField("Created", created)
	}
	if updated, ok := raw["updated_at"].(string); ok {
		writeField("Updated", updated)
	}

	b.WriteString("\n")

	if desc, ok := raw["description"].(string); ok && desc != "" {
		fmt.Fprintf(&b, "[yellow]── Description ──────────────────────────[white]\n\n")
		b.WriteString(desc)
		b.WriteString("\n")
	}

	d.text.SetText(b.String())
}

func (d *Detail) keyboard(ev *tcell.EventKey) *tcell.EventKey {
	if ev.Key() == tcell.KeyEscape || (ev.Key() == tcell.KeyRune && ev.Rune() == 'q') {
		d.nav.Pop()
		return nil
	}
	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'j':
			row, col := d.text.GetScrollOffset()
			d.text.ScrollTo(row+1, col)
			return nil
		case 'k':
			row, col := d.text.GetScrollOffset()
			if row > 0 {
				d.text.ScrollTo(row-1, col)
			}
			return nil
		case 'r':
			d.Refresh()
			return nil
		}
	}
	return ev
}

// Hints returns key binding hints for the menu.
func (d *Detail) Hints() []ui.Hint {
	return []ui.Hint{
		{Key: "j/k", Description: "Scroll"},
		{Key: "r", Description: "Refresh"},
		{Key: "Esc/q", Description: "Back"},
	}
}

// compile-time check
var _ tview.Primitive = (*Detail)(nil)
