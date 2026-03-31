package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const dialogKey = "dialog"

// DialogPages is the interface needed by dialogs to show/dismiss themselves.
type DialogPages interface {
	AddPage(name string, item tview.Primitive, resize, visible bool) *tview.Pages
	RemovePage(name string) *tview.Pages
	ShowPage(name string) *tview.Pages
	HasPage(tag string) bool
}

// modal creates a centered modal overlay wrapping the given form.
// width/height: fixed size. Use 0 for height to fill proportionally.
func modal(form tview.Primitive, width, height int) tview.Primitive {
	if height > 0 {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(form, height, 0, true).
				AddItem(nil, 0, 1, false), width, 0, true).
			AddItem(nil, 0, 1, false)
	}
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, 0, 4, true).
			AddItem(nil, 0, 1, false), width, 0, true).
		AddItem(nil, 0, 1, false)
}

// ShowConfirm shows a Yes/No confirmation dialog.
func ShowConfirm(pages DialogPages, app *tview.Application, title, message string, onOK, onCancel func()) {
	f := tview.NewForm()
	f.SetBorder(true)
	f.SetTitle(" " + title + " ")
	f.SetTitleColor(tcell.ColorAqua)
	f.SetBorderColor(tcell.ColorGray)
	f.SetBackgroundColor(tcell.ColorDarkSlateGray)
	f.SetButtonsAlign(tview.AlignCenter)
	f.SetButtonBackgroundColor(tcell.ColorNavy)
	f.SetButtonTextColor(tcell.ColorWhite)
	f.SetButtonActivatedStyle(tcell.StyleDefault.
		Background(tcell.ColorAqua).
		Foreground(tcell.ColorBlack))

	f.AddTextView("", message, 0, 2, true, false)

	dismiss := func() {
		pages.RemovePage(dialogKey)
	}

	f.AddButton("OK", func() {
		dismiss()
		if onOK != nil {
			onOK()
		}
	})
	f.AddButton("Cancel", func() {
		dismiss()
		if onCancel != nil {
			onCancel()
		}
	})

	f.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		if ev.Key() == tcell.KeyEscape {
			dismiss()
			if onCancel != nil {
				onCancel()
			}
			return nil
		}
		return ev
	})

	m := modal(f, 60, 10)
	pages.AddPage(dialogKey, m, true, false)
	pages.ShowPage(dialogKey)
	app.SetFocus(f)
}

// ShowInput shows a dialog with an editable text area.
func ShowInput(pages DialogPages, app *tview.Application, title, defaultValue string, onOK func(string), onCancel func()) {
	f := tview.NewForm()
	f.SetBorder(true)
	f.SetTitle(" " + title + " ")
	f.SetTitleColor(tcell.ColorAqua)
	f.SetBorderColor(tcell.ColorGray)
	f.SetBackgroundColor(tcell.ColorDarkSlateGray)
	f.SetButtonsAlign(tview.AlignCenter)
	f.SetButtonBackgroundColor(tcell.ColorNavy)
	f.SetButtonTextColor(tcell.ColorWhite)
	f.SetButtonActivatedStyle(tcell.StyleDefault.
		Background(tcell.ColorAqua).
		Foreground(tcell.ColorBlack))
	f.SetFieldBackgroundColor(tcell.ColorBlack)
	f.SetFieldTextColor(tcell.ColorWhite)
	f.SetLabelColor(tcell.ColorAqua)

	f.AddTextArea("", defaultValue, 0, 3, 0, nil)

	dismiss := func() {
		pages.RemovePage(dialogKey)
	}

	f.AddButton("OK", func() {
		ta := f.GetFormItemByLabel("").(*tview.TextArea)
		text := ta.GetText()
		dismiss()
		if onOK != nil {
			onOK(text)
		}
	})
	f.AddButton("Cancel", func() {
		dismiss()
		if onCancel != nil {
			onCancel()
		}
	})

	f.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		if ev.Key() == tcell.KeyEscape {
			dismiss()
			if onCancel != nil {
				onCancel()
			}
			return nil
		}
		return ev
	})

	m := modal(f, 70, 10)
	pages.AddPage(dialogKey, m, true, false)
	pages.ShowPage(dialogKey)
	app.SetFocus(f)
}
