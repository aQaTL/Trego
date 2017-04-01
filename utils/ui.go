package utils

import (
	. "github.com/jroimartin/gocui"
)

func DelNonGlobalKeyBinds(gui *Gui) {
	for _, view := range gui.Views() {
		gui.DeleteKeybindings(view.Name())
	}
}
