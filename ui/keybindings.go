package ui

import (
	. "github.com/jroimartin/gocui"
	"strings"
)

func SetKeyBindings(gui *Gui, manager *TregoManager) error {
	if err := gui.SetKeybinding("", KeyCtrlC, ModNone, quit); err != nil {
		return err
	}

	//Keybinding for switching list on tab keypress
	//I used anonymous function for manager variable access
	if err := gui.SetKeybinding("", KeyTab, ModNone, func(gui *Gui, v *View) error {
		for idx, list := range (manager.Lists) {
			if list.Name == manager.currentView.Name() {
				view, err := gui.SetCurrentView(manager.Lists[(idx + 1) % len(manager.Lists)].Name)
				if err != nil {
					return err
				}
				manager.currentView = view
				if _, err := gui.SetViewOnTop(view.Name()); err != nil {
					return err
				}
				break
			}
		}
		return nil
	}); err != nil {
		return err
	}

	for _, list := range (manager.Lists) {
		if err := gui.SetKeybinding(list.Name, KeyArrowUp, ModNone, cursorUp); err != nil {
			return err
		}
		if err := gui.SetKeybinding(list.Name, KeyArrowDown, ModNone, cursorDown); err != nil {
			return err
		}
	}

	return nil
}

//Moves cursor in list one line up
func cursorUp(g *Gui, v *View) error {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		if cy > 0 {
			if err := v.SetCursor(cx, cy - 1); err != nil {
				return err
			}
		}
		if oy > 0 && cy == 0 {
			if err := v.SetOrigin(ox, oy - 1); err != nil {
				return err
			}
		}
	}

	return nil
}

//Moves cursor in list one line down
func cursorDown(g *Gui, v *View) error {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		if cy + oy < (len(strings.Split(v.ViewBuffer(), "\n")) - 3) {
			if err := v.SetCursor(cx, cy + 1); err != nil {
				if err := v.SetOrigin(ox, oy + 1); err != nil {
					return nil
				}
			}
		}
	}
	return nil
}

func quit(gui *Gui, v *View) error {
	return ErrQuit
}
