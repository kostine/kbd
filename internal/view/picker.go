package view

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/kostine/kbd/internal/config"
	"github.com/kostine/kbd/internal/ui"
)

const (
	modeSavedContexts = iota
	modeBrowser
)

// Picker lets the user select a beads database by browsing folders
// or choosing from saved contexts.
type Picker struct {
	*tview.Flex
	ctxTable *tview.Table
	list     *tview.List
	theme    *ui.Theme
	nav      Navigator
	onSelect func(dbPath string)
	contexts []config.Context
	filtered []config.Context
	filter   string
	cwd      string
	entries  []dirEntry
	mode     int
}

type dirEntry struct {
	name    string
	path    string
	isBeads bool
}

// NewPicker creates a folder picker view.
func NewPicker(theme *ui.Theme, nav Navigator, onSelect func(dbPath string)) *Picker {
	p := &Picker{
		Flex:     tview.NewFlex().SetDirection(tview.FlexRow),
		ctxTable: tview.NewTable(),
		list:     tview.NewList(),
		theme:    theme,
		nav:      nav,
		onSelect: onSelect,
	}

	// Context table setup
	p.ctxTable.SetBorder(true)
	p.ctxTable.SetBorderColor(theme.BorderColor)
	p.ctxTable.SetSelectable(true, false)
	p.ctxTable.SetFixed(1, 0)
	p.ctxTable.SetSelectedStyle(tcell.StyleDefault.
		Background(theme.TableSelected).
		Foreground(theme.TableSelectedFg))
	p.ctxTable.SetInputCapture(p.keyboardContexts)

	// Browser list setup
	p.list.SetBorder(true)
	p.list.SetBorderColor(theme.BorderColor)
	p.list.SetHighlightFullLine(true)
	p.list.SetSelectedBackgroundColor(theme.TableSelected)
	p.list.SetSelectedTextColor(theme.TableSelectedFg)
	p.list.ShowSecondaryText(false)
	p.list.SetInputCapture(p.keyboardBrowser)

	return p
}

// Show initializes and displays the picker.
func (p *Picker) Show() {
	p.contexts, _ = config.LoadContexts()
	if len(p.contexts) > 0 {
		p.showSavedContexts()
	} else {
		p.startBrowser()
	}
}

func (p *Picker) showSavedContexts() {
	p.mode = modeSavedContexts
	p.Clear()
	p.AddItem(p.ctxTable, 0, 1, true)

	p.applyContextFilter()
	p.updateShortcuts()
	p.nav.App().SetFocus(p.ctxTable)
}

// SetFilter applies a search filter to the contexts list.
func (p *Picker) SetFilter(filter string) {
	p.filter = filter
	if p.mode == modeSavedContexts {
		p.applyContextFilter()
		p.updateShortcuts()
	}
}

func (p *Picker) applyContextFilter() {
	f := strings.ToLower(p.filter)
	p.filtered = nil
	for _, c := range p.contexts {
		if f != "" &&
			!strings.Contains(strings.ToLower(c.Name), f) &&
			!strings.Contains(strings.ToLower(c.Path), f) {
			continue
		}
		p.filtered = append(p.filtered, c)
	}

	p.ctxTable.Clear()
	title := " Contexts"
	if p.filter != "" {
		title += " /" + p.filter
	}
	title += " "
	p.ctxTable.SetTitle(title)
	p.ctxTable.SetTitleColor(p.theme.TitleColor)

	p.ctxTable.SetCell(0, 0, tview.NewTableCell("NAME").
		SetSelectable(false).
		SetTextColor(p.theme.TableHeader).
		SetAttributes(tcell.AttrBold).
		SetExpansion(1))
	p.ctxTable.SetCell(0, 1, tview.NewTableCell("PATH").
		SetSelectable(false).
		SetTextColor(p.theme.TableHeader).
		SetAttributes(tcell.AttrBold).
		SetExpansion(2))

	for i, c := range p.filtered {
		name := c.Name
		if c.Last {
			name += " (last)"
		}
		nameCell := tview.NewTableCell(name).SetExpansion(1)
		pathCell := tview.NewTableCell(shortenPath(c.Path)).
			SetTextColor(tcell.ColorGray).
			SetExpansion(2)
		if c.Last {
			nameCell.SetTextColor(tcell.ColorAqua)
		}
		p.ctxTable.SetCell(i+1, 0, nameCell)
		p.ctxTable.SetCell(i+1, 1, pathCell)
	}

	// Pre-select last-used context
	for i, c := range p.filtered {
		if c.Last {
			p.ctxTable.Select(i+1, 0)
			break
		}
	}
}

func (p *Picker) startBrowser() {
	p.mode = modeBrowser
	p.Clear()
	p.AddItem(p.list, 0, 1, true)
	p.nav.App().SetFocus(p.list)

	home, _ := os.UserHomeDir()
	if home == "" {
		home = "/"
	}
	p.browseTo(home)
}

func (p *Picker) browseTo(dir string) {
	p.cwd = dir
	p.list.Clear()
	p.list.SetTitle(" " + shortenPath(dir) + " ")
	p.list.SetTitleColor(p.theme.TitleColor)

	p.entries = nil

	if dir != "/" {
		p.entries = append(p.entries, dirEntry{name: "..", path: filepath.Dir(dir)})
	}

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		p.nav.SetStatus("Error: " + err.Error())
		return
	}

	var dirs []dirEntry
	for _, e := range dirEntries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, ".") && name != ".beads" {
			continue
		}
		full := filepath.Join(dir, name)
		hasBeads := hasBeadsDir(full)
		dirs = append(dirs, dirEntry{name: name, path: full, isBeads: hasBeads})
	}

	sort.SliceStable(dirs, func(i, j int) bool {
		if dirs[i].isBeads != dirs[j].isBeads {
			return dirs[i].isBeads
		}
		return dirs[i].name < dirs[j].name
	})

	p.entries = append(p.entries, dirs...)

	for _, e := range p.entries {
		label := e.name
		if e.isBeads {
			label = "[green::b]" + e.name + " [green][db]"
		}
		if e.name == ".." {
			label = "[gray].."
		}
		p.list.AddItem(label, "", 0, nil)
	}

}

func (p *Picker) keyboardContexts(ev *tcell.EventKey) *tcell.EventKey {
	if ev.Key() == tcell.KeyEnter {
		row, _ := p.ctxTable.GetSelection()
		idx := row - 1
		if idx >= 0 && idx < len(p.filtered) {
			p.selectContext(p.filtered[idx].Path)
		}
		return nil
	}
	if ev.Key() == tcell.KeyEscape {
		if p.filter != "" {
			p.filter = ""
			p.applyContextFilter()
			return nil
		}
		return ev
	}
	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'b':
			p.startBrowser()
			return nil
		case 'd':
			row, _ := p.ctxTable.GetSelection()
			idx := row - 1
			if idx >= 0 && idx < len(p.filtered) {
				config.RemoveContext(p.filtered[idx].Path)
				p.contexts, _ = config.LoadContexts()
				p.filter = ""
				if len(p.contexts) > 0 {
					p.showSavedContexts()
				} else {
					p.startBrowser()
				}
			}
			return nil
		case 'j':
			row, _ := p.ctxTable.GetSelection()
			if row < p.ctxTable.GetRowCount()-1 {
				p.ctxTable.Select(row+1, 0)
			}
			return nil
		case 'k':
			row, _ := p.ctxTable.GetSelection()
			if row > 1 {
				p.ctxTable.Select(row-1, 0)
			}
			return nil
		}
	}
	return ev
}

func (p *Picker) keyboardBrowser(ev *tcell.EventKey) *tcell.EventKey {
	switch ev.Key() {
	case tcell.KeyEnter:
		idx := p.list.GetCurrentItem()
		if idx < 0 || idx >= len(p.entries) {
			return nil
		}
		e := p.entries[idx]
		if e.isBeads {
			dbPath := findBeadsDB(e.path)
			if dbPath != "" {
				p.selectContext(dbPath)
			}
		} else {
			p.browseTo(e.path)
		}
		return nil
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if p.cwd != "/" {
			p.browseTo(filepath.Dir(p.cwd))
		}
		return nil
	case tcell.KeyEscape:
		if len(p.contexts) > 0 {
			p.showSavedContexts()
		}
		return nil
	}

	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'j':
			cur := p.list.GetCurrentItem()
			if cur < p.list.GetItemCount()-1 {
				p.list.SetCurrentItem(cur + 1)
			}
			return nil
		case 'k':
			cur := p.list.GetCurrentItem()
			if cur > 0 {
				p.list.SetCurrentItem(cur - 1)
			}
			return nil
		case '~':
			home, _ := os.UserHomeDir()
			if home != "" {
				p.browseTo(home)
			}
			return nil
		}
	}

	return ev
}

func (p *Picker) selectContext(dbPath string) {
	config.AddContext(dbPath)
	if p.onSelect != nil {
		p.onSelect(dbPath)
	}
}

func (p *Picker) updateShortcuts() {
	n := len(p.filtered)
	if n > 10 {
		n = 10
	}
	shortcuts := make([]ui.Shortcut, n)
	for i := 0; i < n; i++ {
		shortcuts[i] = ui.Shortcut{
			Key:      i,
			Label:    p.filtered[i].Name,
			Category: "context",
			Value:    p.filtered[i].Path,
		}
	}
	p.nav.SetShortcuts(shortcuts)
}

// SelectByIndex selects a context by shortcut index (0-9).
func (p *Picker) SelectByIndex(idx int) {
	if p.mode != modeSavedContexts {
		return
	}
	if idx >= 0 && idx < len(p.filtered) {
		p.selectContext(p.filtered[idx].Path)
	}
}

// Hints returns key binding hints.
func (p *Picker) Hints() []ui.Hint {
	return []ui.Hint{
		{Key: "Enter", Description: "Select"},
		{Key: "j/k", Description: "Navigate"},
		{Key: "/", Description: "Search"},
		{Key: "b", Description: "Browse"},
		{Key: "d", Description: "Delete"},
		{Key: "Bksp", Description: "Parent"},
		{Key: "~", Description: "Home"},
		{Key: "Esc", Description: "Back"},
		{Key: "q", Description: "Quit"},
	}
}

func hasBeadsDir(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, ".beads"))
	return err == nil && info.IsDir()
}

func findBeadsDB(dir string) string {
	beadsDir := filepath.Join(dir, ".beads")
	dolt := filepath.Join(beadsDir, "dolt")
	if info, err := os.Stat(dolt); err == nil && info.IsDir() {
		return dolt
	}
	entries, _ := os.ReadDir(beadsDir)
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".db") {
			return filepath.Join(beadsDir, e.Name())
		}
	}
	return beadsDir
}

func shortenPath(path string) string {
	home, _ := os.UserHomeDir()
	if home != "" && strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}
