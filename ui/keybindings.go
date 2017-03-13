package ui

import (
	. "github.com/jroimartin/gocui"
	"github.com/aqatl/Trego/ui/dialog"
	"log"
	"io/ioutil"
	"strconv"
)

func SetKeyBindings(gui *Gui, manager *TregoManager) (err error) {
	if err = gui.SetKeybinding("", KeyCtrlC, ModNone, quit); err != nil {
		return
	}

	//Testing code, not meat to be used
	if err = gui.SetKeybinding("", 'p', ModNone, func(gui *Gui, v *View) error {
		choice := make(chan bool)
		manager.SelectView(gui, dialog.ConfirmDialog("Are you sure? [y/n]", "", gui, choice).Name())
		go func() {
			dialChoice := <-choice
			manager.currentView = nil
			//test, if choice is being registered correctly
			e := ioutil.WriteFile(
				"choice.txt",
				[]byte(strconv.FormatBool(dialChoice)),
				644)
			if e != nil {
				log.Panicln(e, "!@#")
			}
		}()
		return nil
	}); err != nil {
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
				err = mngr.SelectView(gui, mngr.Lists[(idx + 1) % len(mngr.Lists)].Name)
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
				err = mngr.SelectView(gui, mngr.Lists[(idx - 1) % len(mngr.Lists)].Name)
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
