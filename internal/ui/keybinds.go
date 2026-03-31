package ui

import "github.com/gdamore/tcell/v2"

// KeyAction represents a bindable action.
type KeyAction struct {
	Key         tcell.Key
	Rune        rune
	Description string
	Action      func()
}

// KeyBindings manages key-to-action mappings.
type KeyBindings struct {
	actions map[tcell.Key]KeyAction
	runes   map[rune]KeyAction
}

// NewKeyBindings creates an empty key binding set.
func NewKeyBindings() *KeyBindings {
	return &KeyBindings{
		actions: make(map[tcell.Key]KeyAction),
		runes:   make(map[rune]KeyAction),
	}
}

// BindKey registers a special key action.
func (kb *KeyBindings) BindKey(key tcell.Key, desc string, action func()) {
	kb.actions[key] = KeyAction{Key: key, Description: desc, Action: action}
}

// BindRune registers a rune key action.
func (kb *KeyBindings) BindRune(r rune, desc string, action func()) {
	kb.runes[r] = KeyAction{Rune: r, Description: desc, Action: action}
}

// Handle processes a key event. Returns true if handled.
func (kb *KeyBindings) Handle(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyRune {
		if a, ok := kb.runes[ev.Rune()]; ok {
			a.Action()
			return true
		}
		return false
	}
	if a, ok := kb.actions[ev.Key()]; ok {
		a.Action()
		return true
	}
	return false
}

// Hints returns all visible key bindings as label pairs for the menu.
func (kb *KeyBindings) Hints() []Hint {
	var hints []Hint
	for _, a := range kb.runes {
		hints = append(hints, Hint{Key: string(a.Rune), Description: a.Description})
	}
	for _, a := range kb.actions {
		hints = append(hints, Hint{Key: tcell.KeyNames[a.Key], Description: a.Description})
	}
	return hints
}

// Hint is a key-description pair for menu display.
type Hint struct {
	Key         string
	Description string
}
