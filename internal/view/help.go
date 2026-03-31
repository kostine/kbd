package view

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/kostine/kbd/internal/ui"
	tmux "github.com/kostine/kbd/internal/tmux"
)

// Help displays all available keyboard shortcuts.
type Help struct {
	*tview.Flex
	text  *tview.TextView
	theme *ui.Theme
	nav   Navigator
}

// NewHelp creates the help view.
func NewHelp(theme *ui.Theme, nav Navigator) *Help {
	h := &Help{
		Flex:  tview.NewFlex().SetDirection(tview.FlexRow),
		text:  tview.NewTextView(),
		theme: theme,
		nav:   nav,
	}

	h.text.SetDynamicColors(true)
	h.text.SetScrollable(true)
	h.text.SetBorder(true)
	h.text.SetBorderColor(theme.BorderColor)
	h.text.SetTitle(" Help ")
	h.text.SetTitleColor(theme.TitleColor)
	h.text.SetInputCapture(h.keyboard)

	h.render()
	h.AddItem(h.text, 0, 1, true)
	return h
}

func (h *Help) render() {
	var b strings.Builder

	section := func(title string) {
		fmt.Fprintf(&b, "\n[aqua::b]%s[white]\n", title)
	}
	key := func(k, desc string) {
		fmt.Fprintf(&b, "  [yellow]%-14s[white] %s\n", k, desc)
	}

	section("Navigation")
	key("j / ↓", "Move down")
	key("k / ↑", "Move up")
	key("g", "Jump to top")
	key("G", "Jump to bottom")
	key("Enter", "View issue detail")
	key("Esc", "Back / pop filter")
	key("q", "Quit")

	section("Filtering")
	key("/", "Search by title")
	key(":", "Filter by type (epic, task, ...)")
	key("0-9", "Quick filter from shortcuts panel")
	key("Esc", "Remove last filter")

	section("Actions")
	key("x", "Close issue")
	key("o", "Reopen issue")
	key("r", "Refresh data")
	key("h", "Show this help")

	if tmux.InTmux() {
		section("Tmux Integration")
		key("w", "Work on — send issue ID to Claude pane")
		key("c", "Chat — send freeform message to Claude pane")
	}

	section("Commands (via :)")
	key(":epic", "Show epics")
	key(":task", "Show tasks")
	key(":bug", "Show bugs")
	key(":all", "Show all types")
	key(":context", "Switch database context")
	key(":q", "Quit")

	section("Detail View")
	key("j / k", "Scroll up/down")
	key("r", "Refresh")
	key("Esc / q", "Back to list")

	section("Context Picker")
	key("Enter", "Select context")
	key("b", "Browse for database")
	key("d", "Delete saved context")
	key("~", "Jump to home directory")
	key("Backspace", "Parent directory")

	h.text.SetText(b.String())
}

func (h *Help) keyboard(ev *tcell.EventKey) *tcell.EventKey {
	if ev.Key() == tcell.KeyEscape || (ev.Key() == tcell.KeyRune && (ev.Rune() == 'q' || ev.Rune() == 'h')) {
		h.nav.Pop()
		return nil
	}
	// Scroll with j/k
	if ev.Key() == tcell.KeyRune {
		row, col := h.text.GetScrollOffset()
		switch ev.Rune() {
		case 'j':
			h.text.ScrollTo(row+1, col)
			return nil
		case 'k':
			if row > 0 {
				h.text.ScrollTo(row-1, col)
			}
			return nil
		}
	}
	return ev
}

// Hints returns key binding hints.
func (h *Help) Hints() []ui.Hint {
	return []ui.Hint{
		{Key: "j/k", Description: "Scroll"},
		{Key: "Esc/h/q", Description: "Close"},
	}
}
