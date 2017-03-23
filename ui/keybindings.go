package ui

import (
	. "github.com/jroimartin/gocui"
	"github.com/aqatl/Trego/ui/dialog"
	"github.com/aqatl/Trego/utils"
)

func SetKeyBindings(gui *Gui, manager *TregoManager) (err error) {
	if err = gui.SetKeybinding("", KeyCtrlC, ModNone, quit); err != nil {
		return
	}

	////Testing code, not meant to be used
	//if err = gui.SetKeybinding("", KeyCtrlP, ModNone, func(gui *Gui, v *View) error {
	//	input := make(chan string)
	//	manager.SelectView(gui, dialog.InputDialog("Are you sure? [y/n]", "", "", gui, input).Name())
	//	go func() {
	//		userInput := <-input
	//		manager.currentView = nil
	//		//test, if choice is being registered correctly
	//		//e := ioutil.WriteFile(
	//		//	"choice.txt",
	//		//	[]byte(userInput),
	//		//	644)
	//		fmt.Fprint(os.Stderr, userInput)
	//		//if e != nil {
	//		//	log.Panicln(e, "!@#")
	//		//}
	//	}()
	//	return nil
	//}); err != nil {
	//	return
	//}

	//Testing testing testing
	if err = gui.SetKeybinding("", KeyCtrlP, ModNone, func(gui *Gui, v *View) error {
		option := make(chan bool)
		currView := gui.CurrentView()
		//Prevents nested dialogs and other glitches (like double handler call)
		for _, view := range gui.Views() {
			gui.DeleteKeybindings(view.Name())
		}
		utils.ErrCheck(
			manager.SelectView(
				gui,
				dialog.ConfirmDialog("message", "title", gui, option).Name()))

		go func() {
			_ = <-option
			manager.currView = currView

			gui.Execute(func(gui *Gui) error {
				utils.ErrCheck(manager.SelectView(gui, manager.currView.Name()))
				SetKeyBindings(gui, manager)
				return nil
			})
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
		mngr.currListIdx++
		err = mngr.SelectView(gui, mngr.Lists[(mngr.currListIdx) % len(mngr.Lists)].Name)
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
