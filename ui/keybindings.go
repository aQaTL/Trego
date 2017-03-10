package ui

import (
	. "github.com/jroimartin/gocui"
)

func SetKeyBindings(gui *Gui, manager *TregoManager) (err error) {
	if err = gui.SetKeybinding("", KeyCtrlC, ModNone, quit); err != nil {
		return
	}

	for _, list := range (manager.Lists) {
		if err = gui.SetKeybinding(list.Name, KeyArrowUp, ModNone, CursorUp); err != nil {
			return
		}
		if err = gui.SetKeybinding(list.Name, KeyArrowDown, ModNone, CursorDown); err != nil {
			return
		}
		if err = addListSwitchingFunc(gui, list.Name, manager); err != nil {
			return
		}
	}
	return
}


//Keybinding for switching list on tab keypress
//I used anonymous function for manager variable access
func addListSwitchingFunc(gui *Gui, viewName string, mngr *TregoManager) (err error) {
	switchListRight := func(gui *Gui, v *View) (err error) {
		for idx, list := range (mngr.Lists) {
			if list.Name == mngr.currentView.Name() {
				err = mngr.SelectList(gui, mngr.Lists[(idx + 1) % len(mngr.Lists)].Name)
				break
			}
		}
		return
	}
	switchListLeft := func(gui *Gui, v *View) (err error) {
		for idx, list := range (mngr.Lists) {
			if list.Name == mngr.currentView.Name() {
				if idx == 0 {
					idx = len(mngr.Lists)
				}
				err = mngr.SelectList(gui, mngr.Lists[(idx - 1) % len(mngr.Lists)].Name)
				break
			}
		}
		return
	}

	gui.SetKeybinding(viewName, KeyTab, ModNone, switchListRight)
	gui.SetKeybinding(viewName, KeyArrowRight, ModNone, switchListRight)
	gui.SetKeybinding(viewName, KeyArrowLeft, ModNone, switchListLeft)

	return
}

func quit(gui *Gui, v *View) error {
	return ErrQuit
}
