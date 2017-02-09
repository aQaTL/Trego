package ui

import (
	. "github.com/jroimartin/gocui"
	"github.com/fatih/color"
)

const (
	SHORTCUTS string = "shortcuts"
)

func Layout(gui *Gui) error {
	if err := shortcutsView(gui); err != nil {
		return  err
	}

	return nil
}

//bottom bar with shortcuts
func shortcutsView(gui *Gui) error {
	maxX, maxY := gui.Size()
	if v, err := gui.SetView(SHORTCUTS, 0, maxY - 4, maxX - 1, maxY - 1); err != nil {
		if err != ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Highlight = false
		v.BgColor = ColorBlack

		color.Output = v
		color.New(color.FgYellow).Add(color.Bold).Println("Ala nie ma kota")

		if _, err := gui.SetCurrentView(SHORTCUTS); err != nil {
			return err
		}
	}
	return nil
}

func SetKeyBindings(gui *Gui) error {
	if err := gui.SetKeybinding("", KeyCtrlC, ModNone, quit); err != nil {
		return err
	}

	return nil
}

func quit(gui *Gui, v *View) error {
	return ErrQuit
}
