package ui

import "github.com/gdamore/tcell/v2"

// Theme defines the color palette for the TUI.
type Theme struct {
	BgColor     tcell.Color
	FgColor     tcell.Color
	BorderColor tcell.Color
	TitleColor  tcell.Color

	TableHeader    tcell.Color
	TableSelected  tcell.Color
	TableSelectedFg tcell.Color

	StatusOpen       tcell.Color
	StatusInProgress tcell.Color
	StatusBlocked    tcell.Color
	StatusClosed     tcell.Color

	MenuKey  tcell.Color
	MenuDesc tcell.Color

	PromptFg tcell.Color
	PromptBg tcell.Color

	ErrorFg tcell.Color

	DeltaAdd    tcell.Color
	DeltaUpdate tcell.Color
	DeltaDelete tcell.Color
}

// DefaultTheme returns a sensible dark theme.
func DefaultTheme() *Theme {
	return &Theme{
		BgColor:     tcell.ColorDefault,
		FgColor:     tcell.ColorWhite,
		BorderColor: tcell.ColorGray,
		TitleColor:  tcell.ColorAqua,

		TableHeader:     tcell.ColorAqua,
		TableSelected:   tcell.ColorDarkSlateGray,
		TableSelectedFg: tcell.ColorDefault,

		StatusOpen:       tcell.ColorGreen,
		StatusInProgress: tcell.ColorYellow,
		StatusBlocked:    tcell.ColorRed,
		StatusClosed:     tcell.ColorGray,

		MenuKey:  tcell.ColorAqua,
		MenuDesc: tcell.ColorDefault,

		PromptFg: tcell.ColorDefault,
		PromptBg: tcell.ColorDefault,

		ErrorFg: tcell.ColorRed,

		DeltaAdd:    tcell.ColorGreen,
		DeltaUpdate: tcell.ColorYellow,
		DeltaDelete: tcell.ColorRed,
	}
}

// StatusColor returns the color for a given issue status.
func (t *Theme) StatusColor(status string) tcell.Color {
	switch status {
	case "open":
		return t.FgColor
	case "in_progress":
		return t.StatusInProgress
	case "blocked":
		return t.StatusBlocked
	case "closed":
		return t.StatusClosed
	case "on_hold", "on-hold":
		return tcell.ColorOrange
	case "deferred", "awaiting_response":
		return tcell.ColorBlue
	case "pinned":
		return tcell.ColorPurple
	default:
		return t.FgColor
	}
}

// PriorityColor returns the color for a given priority level.
func (t *Theme) PriorityColor(priority int) tcell.Color {
	switch priority {
	case 0:
		return tcell.ColorRed
	case 1:
		return tcell.ColorOrange
	case 2:
		return t.FgColor
	default:
		return tcell.ColorGray
	}
}
