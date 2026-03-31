package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Prompt is a command/filter input bar.
type Prompt struct {
	*tview.InputField
	theme    *Theme
	onSubmit func(text string)
	onCancel func()
}

// NewPrompt creates a new prompt bar.
func NewPrompt(theme *Theme) *Prompt {
	p := &Prompt{
		InputField: tview.NewInputField(),
		theme:      theme,
	}
	p.SetFieldBackgroundColor(theme.PromptBg)
	p.SetFieldTextColor(theme.PromptFg)
	p.SetInputCapture(p.keyboard)
	return p
}

// Activate shows the prompt with the given prefix (e.g., ":" or "/").
func (p *Prompt) Activate(prefix string, submit func(string), cancel func()) {
	p.onSubmit = submit
	p.onCancel = cancel
	p.SetLabel(prefix)
	p.SetText("")
}

func (p *Prompt) keyboard(ev *tcell.EventKey) *tcell.EventKey {
	switch ev.Key() {
	case tcell.KeyEnter:
		if p.onSubmit != nil {
			p.onSubmit(p.GetText())
		}
		return nil
	case tcell.KeyEscape:
		if p.onCancel != nil {
			p.onCancel()
		}
		return nil
	}
	return ev
}
