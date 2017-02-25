package ui

import (
	. "github.com/jroimartin/gocui"
	"github.com/fatih/color"
	"math"
	"github.com/VojtechVitek/go-trello"
	"strings"
)

const (
	BOTTOM_BAR string = "botbar"
	TOP_BAR string = "topbar"
	LIST_WIDTH int = 24
)

func (manager *TregoManager) Layout(gui *Gui) error {
	if err := bottomBar(gui); err != nil {
		return err
	}
	if err := topBar(gui); err != nil {
		return err
	}

	for idx, list := range (manager.Lists) {
		if err := AddList(gui, list, idx); err != nil {
			return err
		}
	}

	if manager.currentView == nil && len(manager.Lists) > 0 {
		manager.selectList(gui, manager.Lists[0].Name)
	}

	if _, err := gui.SetCurrentView(manager.currentView.Name()); err != nil {
		return err
	}

	return nil
}

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
		if err := gui.SetKeybinding(list.Name, KeyArrowUp, ModNone, listUp); err != nil {
			return err
		}
		if err := gui.SetKeybinding(list.Name, KeyArrowDown, ModNone, listDown); err != nil {
			return err
		}
	}

	return nil
}

func AddList(gui *Gui, list trello.List, index int) error {
	_, maxY := gui.Size()
	if v, err := gui.SetView(list.Name,
		index * LIST_WIDTH + int(math.Abs(sign(index))), 3,
		index * LIST_WIDTH + LIST_WIDTH, maxY - 5);
			err != nil {

		if err != ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Highlight = true
		v.Wrap = true
		v.BgColor = ColorBlack
		v.SelBgColor = ColorGreen
		v.SelFgColor = ColorBlack
		v.Title = list.Name
		gui.Cursor = true

		color.Output = v

		cards, err := list.Cards()
		if err != nil {
			return err
		}
		for idx, card := range (cards) {
			color.New(color.BgBlack).Add(color.FgWhite).Printf("%d. %v\n", idx, card.Name)
		}
	}

	return nil
}

func topBar(gui *Gui) error {
	maxX, _ := gui.Size()
	if v, err := gui.SetView(TOP_BAR, 0, 0, maxX - 1, 2); err != nil {
		if err != ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Highlight = false
		v.BgColor = ColorBlack

		color.Output = v
		color.New(color.FgYellow).Add(color.Bold).Printf("Board: %v", "Trego")
	}

	return nil
}

//bottom bar with shortcuts
func bottomBar(gui *Gui) error {
	maxX, maxY := gui.Size()
	if v, err := gui.SetView(BOTTOM_BAR, 0, maxY - 4, maxX - 1, maxY - 1); err != nil {
		if err != ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Highlight = false
		v.BgColor = ColorBlack

		color.Output = v
		color.New(color.FgYellow).Add(color.Bold).Println("Ala nie ma kota")

	}
	return nil
}

func listUp(g *Gui, v *View) error {
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

func listDown(g *Gui, v *View) error {
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

func sign(x int) float64 {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	} else {
		return 0
	}
}