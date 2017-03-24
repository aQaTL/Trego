package ui

import (
	. "github.com/jroimartin/gocui"
)

func SetKeyBindings(gui *Gui, mngr *TregoManager) (err error) {
	if err = gui.SetKeybinding("", KeyCtrlC, ModNone, quit); err != nil {
		return
	}

	for _, list := range mngr.Lists {
		if err = gui.SetKeybinding(list.Name, KeyArrowUp, ModNone, CursorUp); err != nil {
			return
		}
		if err = gui.SetKeybinding(list.Name, KeyArrowDown, ModNone, CursorDown); err != nil {
			return
		}
		if err = addListSwitchingFunc(gui, list.Name, mngr); err != nil {
			return
		}
	}
	return
}

//Keybinding for switching list on tab keypress
//I used anonymous function for manager variable access
func addListSwitchingFunc(gui *Gui, viewName string, mngr *TregoManager) (err error) {
	switchListRight := func(gui *Gui, v *View) (err error) {
		mngr.currListIdx = (mngr.currListIdx + 1) % len(mngr.Lists)
		err = mngr.SelectView(gui, mngr.Lists[mngr.currListIdx].Name)
		return
	}
	switchListLeft := func(gui *Gui, v *View) (err error) {
		if mngr.currListIdx == 0 {
			mngr.currListIdx = len(mngr.Lists)
		}
		mngr.currListIdx--
		err = mngr.SelectView(gui, mngr.Lists[mngr.currListIdx % len(mngr.Lists)].Name)
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
