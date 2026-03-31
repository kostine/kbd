package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HeaderGridRows is the height of the shortcuts/hints grid.
const HeaderGridRows = 5

// Header has two main parts: LEFT (key hints) and RIGHT (numbered shortcuts),
// plus a status line underneath that expands for errors.
type Header struct {
	*tview.Flex
	Hints     *HintsPanel
	Shortcuts *ShortcutPanel
	status    *tview.TextView
	grid      *tview.Flex
	theme     *Theme
	hasError  bool
}

// Rows returns the current height needed for the header.
func (h *Header) Rows() int {
	if h.hasError {
		return HeaderGridRows + 3
	}
	return HeaderGridRows + 1
}

// NewHeader creates a header bar.
func NewHeader(theme *Theme) *Header {
	h := &Header{
		Flex:      tview.NewFlex().SetDirection(tview.FlexRow),
		Hints:     NewHintsPanel(theme, 2),
		Shortcuts: NewShortcutPanel(theme),
		status:    tview.NewTextView().SetDynamicColors(true),
		theme:     theme,
	}

	h.status.SetBackgroundColor(tcell.ColorDefault)
	h.status.SetWrap(true)
	h.status.SetWordWrap(true)

	// Grid row: LEFT hints | RIGHT shortcuts
	h.grid = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(h.Hints, 0, 1, false).
		AddItem(h.Shortcuts, 0, 1, false)

	h.rebuild()

	return h
}

func (h *Header) rebuild() {
	h.Clear()
	h.AddItem(h.grid, HeaderGridRows, 0, false)
	if h.hasError {
		h.AddItem(h.status, 3, 0, false)
	} else {
		h.AddItem(h.status, 1, 0, false)
	}
}

// SetHints updates the left-side key hint panel.
func (h *Header) SetHints(hints []Hint) {
	h.Hints.SetHints(hints)
}

// SetContext is a no-op.
func (h *Header) SetContext(text string) {}

// SetStatus shows a temporary status/error message below the header grid.
// Empty string clears it. Error messages get extra rows for wrapping.
func (h *Header) SetStatus(msg string) {
	h.status.Clear()
	wasError := h.hasError
	if msg == "" {
		h.hasError = false
	} else {
		h.hasError = strings.Contains(msg, "Error")
		fmt.Fprintf(h.status, "[yellow]%s", msg)
	}
	if h.hasError != wasError {
		h.rebuild()
	}
}

// SetCrumbs is a no-op.
func (h *Header) SetCrumbs(crumbs []string) {}
