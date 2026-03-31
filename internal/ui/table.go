package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Table is an enhanced table widget with selection tracking and vim navigation.
type Table struct {
	*tview.Table
	theme    *Theme
	headers  []string
	onSelect func(row int)
	keyBinds *KeyBindings
	lastSel  int
}

// NewTable creates a new table with the given headers.
func NewTable(theme *Theme) *Table {
	t := &Table{
		Table:    tview.NewTable(),
		theme:    theme,
		keyBinds: NewKeyBindings(),
		lastSel:  -1,
	}

	t.SetSelectable(true, false)
	t.SetFixed(1, 0)
	t.SetBorder(true)
	t.SetBorderColor(theme.BorderColor)

	// Disable tview's built-in selection coloring — we handle it manually
	// in onSelectionChanged to preserve per-cell text colors.
	t.SetSelectedStyle(tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorDefault))

	t.SetSelectionChangedFunc(t.onSelectionChanged)
	t.SetInputCapture(t.keyboard)
	return t
}

// onSelectionChanged updates cell backgrounds to highlight the selected row.
func (t *Table) onSelectionChanged(row, col int) {
	// Restore previous row
	if t.lastSel >= 1 {
		colCount := t.GetColumnCount()
		for c := 0; c < colCount; c++ {
			if cell := t.GetCell(t.lastSel, c); cell != nil {
				cell.SetBackgroundColor(tcell.ColorDefault)
				cell.SetAttributes(tcell.AttrNone)
			}
		}
	}

	// Highlight new row
	if row >= 1 {
		colCount := t.GetColumnCount()
		for c := 0; c < colCount; c++ {
			if cell := t.GetCell(row, c); cell != nil {
				cell.SetBackgroundColor(t.theme.TableSelected)
				cell.SetAttributes(tcell.AttrBold)
			}
		}
	}

	t.lastSel = row
}

// SetHeaders sets the column headers.
func (t *Table) SetHeaders(headers ...string) {
	t.headers = headers
	for col, h := range headers {
		cell := tview.NewTableCell(h).
			SetSelectable(false).
			SetTextColor(t.theme.TableHeader).
			SetAttributes(tcell.AttrBold)
		t.SetCell(0, col, cell)
	}
}

// SetRow sets a data row at the given index (1-based, row 0 is headers).
func (t *Table) SetRow(row int, values ...string) {
	for col, v := range values {
		cell := tview.NewTableCell(v).
			SetExpansion(1)
		t.SetCell(row, col, cell)
	}
}

// SetRowWithColor sets a data row with a specific color.
func (t *Table) SetRowWithColor(row int, color tcell.Color, values ...string) {
	for col, v := range values {
		cell := tview.NewTableCell(v).
			SetExpansion(1).
			SetTextColor(color)
		t.SetCell(row, col, cell)
	}
}

// SetCellColor changes the text color of a specific cell.
func (t *Table) SetCellColor(row, col int, color tcell.Color) {
	if cell := t.GetCell(row, col); cell != nil {
		cell.SetTextColor(color)
	}
}

// ClearData removes all rows except headers.
func (t *Table) ClearData() {
	rowCount := t.GetRowCount()
	for r := rowCount - 1; r >= 1; r-- {
		t.RemoveRow(r)
	}
}

// OnSelect sets the callback when Enter is pressed on a row.
func (t *Table) OnSelect(fn func(row int)) {
	t.onSelect = fn
}

// KeyBindings returns the table's key bindings for additional binding.
func (t *Table) KeyBindings() *KeyBindings {
	return t.keyBinds
}

func (t *Table) keyboard(ev *tcell.EventKey) *tcell.EventKey {
	if t.keyBinds.Handle(ev) {
		return nil
	}

	row, _ := t.GetSelection()
	maxRow := t.GetRowCount() - 1

	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'j':
			if row < maxRow {
				t.Select(row+1, 0)
			}
			return nil
		case 'k':
			if row > 1 {
				t.Select(row-1, 0)
			}
			return nil
		case 'g':
			t.Select(1, 0)
			return nil
		case 'G':
			t.Select(maxRow, 0)
			return nil
		}
	}

	if ev.Key() == tcell.KeyEnter && t.onSelect != nil {
		t.onSelect(row)
		return nil
	}

	return ev
}
