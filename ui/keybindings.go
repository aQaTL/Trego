package ui

import (
	. "github.com/jroimartin/gocui"
	"github.com/aqatl/go-trello"
	"github.com/aqatl/Trego/utils"
	"github.com/aqatl/Trego/ui/dialog"
	"log"
	"math"
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

		if err = addListAddingFunc(gui, list.Name, mngr); err != nil {
			return
		}

		if err = addCardAddingFunc(gui, list.Name, mngr); err != nil {
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

func addCardAddingFunc(gui *Gui, viewName string, mngr *TregoManager) error {
	return gui.SetKeybinding(viewName, 'n', ModNone, func(gui *Gui, view *View) error {
		cardNameC := make(chan string)
		for _, view := range gui.Views() {
			gui.DeleteKeybindings(view.Name())
		}

		utils.ErrCheck(mngr.SelectView(
			gui,
			dialog.InputDialog(
				"Name your card",
				"New card",
				"",
				gui,
				cardNameC).Name()))

		go func() {
			cardName := <-cardNameC
			SetKeyBindings(gui, mngr)
			list := mngr.Lists[mngr.currListIdx]
			card, err := list.AddCard(trello.Card{
				IdList: list.Id,
				Name: cardName,
				Pos: math.MaxFloat64, //end of the list
			})

			if err != nil {
				log.Panicf("Card add: %v in list %v", err, list.Name)
			}

			gui.Execute(func(gui *Gui) error {
				utils.ErrCheck(gui.DeleteView(list.Name)) //Forces view update
				return nil
			})
			log.Printf("Successfully added new card: %v", card.Name)
		}()
		return nil
	})
}

func addListAddingFunc(gui *Gui, viewName string, mngr *TregoManager) error {
	return gui.SetKeybinding(viewName, 'l', ModNone, func(gui *Gui, view *View) error {
		viewNameC := make(chan string)
		for _, view := range gui.Views() {
			gui.DeleteKeybindings(view.Name())
		}

		utils.ErrCheck(mngr.SelectView(
			gui,
			dialog.InputDialog(
				"Name your list",
				"New list",
				"",
				gui,
				viewNameC).Name()))

		go func() {
			viewName := <-viewNameC

			list, err := mngr.CurrBoard.AddList(trello.List{
				Name: viewName,
				IdBoard: mngr.CurrBoard.Id,
				Pos: math.MaxFloat32,
			})
			if err != nil {
				log.Printf("List add: %v", err)
			}

			mngr.Lists = append(mngr.Lists, *list)
			utils.ErrCheck(AddList(gui, *list, len(mngr.Lists) - 1))
			SetKeyBindings(gui, mngr)
		}()
		return nil
	})
}

func quit(gui *Gui, v *View) error {
	return ErrQuit
}
