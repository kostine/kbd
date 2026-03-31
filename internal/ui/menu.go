package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Menu renders keyboard hints at the bottom of the screen.
type Menu struct {
	*tview.TextView
	theme *Theme
}

// NewMenu creates a new menu bar.
func NewMenu(theme *Theme) *Menu {
	m := &Menu{
		TextView: tview.NewTextView(),
		theme:    theme,
	}
	m.SetDynamicColors(true)
	m.SetBackgroundColor(tcell.ColorDefault)
	return m
}

// SetHints updates the menu with the given hints.
func (m *Menu) SetHints(hints []Hint) {
	m.Clear()
	for i, h := range hints {
		if i > 0 {
			fmt.Fprint(m, "  ")
		}
		fmt.Fprintf(m, "[aqua]<%s>[white] %s", h.Key, h.Description)
	}
}
