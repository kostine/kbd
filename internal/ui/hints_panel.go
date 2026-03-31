package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HintsPanel renders key hints in a columnar grid layout.
type HintsPanel struct {
	*tview.TextView
	theme *Theme
	hints []Hint
	cols  int
}

// NewHintsPanel creates a hints panel with the given number of columns.
func NewHintsPanel(theme *Theme, cols int) *HintsPanel {
	p := &HintsPanel{
		TextView: tview.NewTextView(),
		theme:    theme,
		cols:     cols,
	}
	p.SetDynamicColors(true)
	p.SetBackgroundColor(tcell.ColorDefault)
	p.SetWrap(false)
	return p
}

// SetHints updates the displayed hints.
func (p *HintsPanel) SetHints(hints []Hint) {
	p.hints = hints
	p.render()
}

// Rows returns how many rows are needed for the current hints.
func (p *HintsPanel) Rows() int {
	n := len(p.hints)
	rows := (n + p.cols - 1) / p.cols
	if rows < 1 {
		return 1
	}
	return rows
}

func (p *HintsPanel) render() {
	p.Clear()

	n := len(p.hints)
	rows := p.Rows()

	for row := 0; row < rows; row++ {
		if row > 0 {
			fmt.Fprint(p, "\n")
		}
		for col := 0; col < p.cols; col++ {
			idx := col*rows + row
			if idx >= n {
				continue
			}
			if col > 0 {
				fmt.Fprint(p, "  ")
			}
			h := p.hints[idx]
			fmt.Fprintf(p, "[aqua]%-6s[white]%-10s", h.Key, h.Description)
		}
	}
}
