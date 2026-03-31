package app

import (
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/kostine/kbd/internal/bd"
	"github.com/kostine/kbd/internal/ui"
	"github.com/kostine/kbd/internal/view"
)

// Padding in characters/rows around the entire UI.
const (
	padLeft   = 2
	padRight  = 2
	padTop    = 1
	padBottom = 1
)

// App is the main kbd application.
type App struct {
	tapp          *tview.Application
	root          *tview.Flex
	pages         *tview.Pages
	layout        *tview.Flex
	header        *ui.Header
	statusBar     *ui.StatusBar
	prompt        *ui.Prompt
	theme         *ui.Theme
	dbPath        string
	client        *bd.Client
	issues        *view.Issues
	picker        *view.Picker
	stack         []stackEntry
	promptMode    bool
	statusBarRows int
}

type stackEntry struct {
	name  string
	view  tview.Primitive
	hints []ui.Hint
}

// New creates a new kbd application.
func New(dbPath string) *App {
	theme := ui.DefaultTheme()

	a := &App{
		tapp:      tview.NewApplication(),
		pages:     tview.NewPages(),
		header:    ui.NewHeader(theme),
		statusBar: ui.NewStatusBar(theme),
		prompt:    ui.NewPrompt(theme),
		theme:     theme,
		dbPath:    dbPath,
	}

	a.rebuildLayout()

	inner := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, padTop, 0, false).
		AddItem(a.layout, 0, 1, true).
		AddItem(nil, padBottom, 0, false)

	a.root = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nil, padLeft, 0, false).
		AddItem(inner, 0, 1, true).
		AddItem(nil, padRight, 0, false)

	a.tapp.SetInputCapture(a.keyboard)

	return a
}

func (a *App) rebuildLayout() {
	if a.layout == nil {
		a.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	}
	a.layout.Clear()
	a.layout.AddItem(a.header, a.header.Rows(), 0, false)
	if a.promptMode {
		a.layout.AddItem(a.prompt, 1, 0, false)
	}
	a.layout.AddItem(a.pages, 0, 1, true)
	rows := a.statusBarRows
	if rows < 1 {
		rows = 1
	}
	a.layout.AddItem(a.statusBar, rows, 0, false)
}

// Run starts the application.
func (a *App) Run() error {
	resolved := a.resolveDB()
	if resolved != "" {
		a.startWithDB(resolved)
	} else {
		a.showPicker()
	}
	a.tapp.SetRoot(a.root, true)
	err := a.tapp.Run()
	if a.issues != nil {
		a.issues.StopAutoRefresh()
	}
	return err
}

// resolveDB tries to find a beads database.
// Priority: --db flag > .beads in cwd > empty (show picker).
func (a *App) resolveDB() string {
	if a.dbPath != "" {
		return a.dbPath
	}
	// Check current working directory
	cwd, err := os.Getwd()
	if err == nil {
		candidates := []string{
			filepath.Join(cwd, ".beads", "dolt"),
			filepath.Join(cwd, ".beads"),
		}
		for _, c := range candidates {
			if info, err := os.Stat(c); err == nil && info.IsDir() {
				return c
			}
		}
		// Check for .db files in .beads
		beadsDir := filepath.Join(cwd, ".beads")
		if entries, err := os.ReadDir(beadsDir); err == nil {
			for _, e := range entries {
				if filepath.Ext(e.Name()) == ".db" {
					return filepath.Join(beadsDir, e.Name())
				}
			}
		}
	}
	return ""
}

// startWithDB initializes the issue view with the given database path.
func (a *App) startWithDB(dbPath string) {
	a.dbPath = dbPath
	a.client = bd.NewClient(dbPath)

	// Clear any existing views
	for len(a.stack) > 0 {
		top := a.stack[len(a.stack)-1]
		a.stack = a.stack[:len(a.stack)-1]
		a.pages.RemovePage(top.name)
	}

	if a.issues != nil {
		a.issues.StopAutoRefresh()
	}

	a.issues = view.NewIssues(a.client, a.theme, a)
	a.Push("issues", a.issues, a.issues.Hints())
	a.issues.Refresh()
}

// showPicker displays the folder picker / context selector.
func (a *App) showPicker() {
	a.picker = view.NewPicker(a.theme, a, func(dbPath string) {
		a.picker = nil
		a.startWithDB(dbPath)
	})
	a.Push("picker", a.picker, a.picker.Hints())
	a.picker.Show()
}

func (a *App) isPickerActive() bool {
	if len(a.stack) == 0 {
		return false
	}
	return a.stack[len(a.stack)-1].name == "picker"
}

// Push adds a view to the navigation stack.
func (a *App) Push(name string, v tview.Primitive, hints []ui.Hint) {
	a.stack = append(a.stack, stackEntry{name: name, view: v, hints: hints})
	a.pages.AddAndSwitchToPage(name, v, true)
	a.updateHeaderHints(hints)
}

// Pop removes the top view from the stack.
func (a *App) Pop() {
	if len(a.stack) <= 1 {
		return
	}
	top := a.stack[len(a.stack)-1]
	a.stack = a.stack[:len(a.stack)-1]
	a.pages.RemovePage(top.name)

	// Clean up picker reference if we're popping it
	if top.name == "picker" {
		a.picker = nil
	}

	current := a.stack[len(a.stack)-1]
	a.pages.SwitchToPage(current.name)
	a.updateHeaderHints(current.hints)

	// Restore issues shortcuts when returning to issues view
	if current.name == "issues" && a.issues != nil {
		a.issues.RefreshShortcuts()
	}
}

// SetStatus updates the status message in the header.
func (a *App) SetStatus(msg string) {
	oldRows := a.header.Rows()
	a.header.SetStatus(msg)
	if a.header.Rows() != oldRows {
		a.rebuildLayout()
	}
}

// SetCounts updates the footer status bar with type/status counts.
func (a *App) SetCounts(counts []ui.TypeStatusCounts) {
	a.statusBar.SetCounts(counts)
	needed := ui.LineCount(counts)
	if needed != a.statusBarRows {
		a.statusBarRows = needed
		a.rebuildLayout()
	}
}

// SetShortcuts updates the header shortcut panel.
func (a *App) SetShortcuts(shortcuts []ui.Shortcut) {
	a.header.Shortcuts.SetShortcuts(shortcuts)
}

// App returns the underlying tview.Application.
func (a *App) App() *tview.Application {
	return a.tapp
}

// Pages returns the pages container for dialog management.
func (a *App) Pages() ui.DialogPages {
	return a.pages
}

func (a *App) updateHeaderHints(hints []ui.Hint) {
	a.header.SetHints(hints)
}

func (a *App) keyboard(ev *tcell.EventKey) *tcell.EventKey {
	if a.promptMode {
		return ev
	}

	// Let dialogs handle their own input
	if a.pages.HasPage("dialog") {
		return ev
	}

	switch ev.Key() {
	case tcell.KeyCtrlC:
		a.tapp.Stop()
		return nil
	case tcell.KeyEscape:
		// Pop view stack first if not on root view
		if len(a.stack) > 1 {
			a.Pop()
			return nil
		}
		// On root view, pop filters
		if a.issues != nil && a.issues.PopFilter() {
			return nil
		}
	}

	if ev.Key() == tcell.KeyRune {
		r := ev.Rune()

		// Number keys 0-9: quick shortcuts
		if r >= '0' && r <= '9' {
			idx := int(r - '0')
			if a.picker != nil && a.isPickerActive() {
				a.picker.SelectByIndex(idx)
			} else if s := a.header.Shortcuts.Get(idx); s != nil {
				switch s.Category {
				case "command":
					switch s.Value {
					case "context":
						if a.issues != nil {
							a.issues.StopAutoRefresh()
						}
						a.showPicker()
					}
				default:
					if a.issues != nil {
						a.issues.ApplyShortcut(s.Category, s.Value)
					}
				}
			}
			return nil
		}

		switch r {
		case 'q':
			if len(a.stack) <= 1 {
				a.tapp.Stop()
				return nil
			}
		case '/':
			a.showPrompt("/", a.onSearchSubmit)
			return nil
		case ':':
			a.showPrompt(":", a.onCommandSubmit)
			return nil
		}
	}

	return ev
}

func (a *App) showPrompt(prefix string, onSubmit func(string)) {
	a.promptMode = true
	a.prompt.Activate(prefix, func(text string) {
		a.hidePrompt()
		onSubmit(text)
	}, func() {
		a.hidePrompt()
	})
	a.rebuildLayout()
	a.tapp.SetFocus(a.prompt)
}

func (a *App) hidePrompt() {
	a.promptMode = false
	a.rebuildLayout()
	if len(a.stack) > 0 {
		a.tapp.SetFocus(a.stack[len(a.stack)-1].view)
	}
}

func (a *App) onSearchSubmit(text string) {
	if a.picker != nil && a.isPickerActive() {
		a.picker.SetFilter(text)
	} else if a.issues != nil {
		a.issues.SetTitleFilter(text)
	}
}

func (a *App) onCommandSubmit(cmd string) {
	switch cmd {
	case "q", "quit":
		a.tapp.Stop()
	case "context", "ctx":
		if a.issues != nil {
			a.issues.StopAutoRefresh()
		}
		a.showPicker()
	default:
		if a.issues != nil {
			for len(a.stack) > 1 {
				a.Pop()
			}
			a.issues.SetTypeFilter(cmd)
		}
	}
}
