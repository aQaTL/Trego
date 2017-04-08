package ui

import (
	. "github.com/jroimartin/gocui"
	"github.com/aqatl/go-trello"
	"github.com/aqatl/Trego/utils"
	"github.com/aqatl/Trego/ui/dialog"
	"log"
	"math"
	"github.com/aqatl/Trego/conn"
)

func SetKeyBindings(gui *Gui, mngr *TregoManager) (err error) {
	if err = gui.SetKeybinding("", KeyCtrlC, ModNone, quit); err != nil {
		return
	}

	for _, list := range mngr.Lists {
		utils.ErrCheck(
			gui.SetKeybinding(list.Name, KeyArrowUp, ModNone, utils.CursorUp),
			gui.SetKeybinding(list.Name, KeyArrowDown, ModNone, utils.CursorDown),
			addListSwitchingFunc(gui, list.Name, mngr),
			addListAddingFunc(gui, list.Name, mngr),
			addCardAddingFunc(gui, list.Name, mngr),
			addCardMovingFunc(gui, list.Name, mngr),
			addCardDeletingFunc(gui, list.Name, mngr),
			addBoardSwitchingFunc(gui, list.Name, mngr),
		)
	}
	return
}

func addCardDeletingFunc(gui *Gui, listName string, mngr *TregoManager) error {
	return gui.SetKeybinding(listName, 'd', ModNone, func(gui *Gui, view *View) error {
		delModeC := make(chan int)
		utils.DelNonGlobalKeyBinds(gui)

		utils.ErrCheck(
			mngr.SelectView(
				gui,
				dialog.SelectDialog(
					"Choose action",
					gui,
					delModeC,
					[]string{"Archive", "Delete"}).
						Name()))

		go func() {
			if delMode, ok := <-delModeC; ok {
				cards, err := mngr.Lists[mngr.currListIdx].Cards()
				utils.ErrCheck(err)
				cardIdx := SelectedItemIdx(view)
				if delMode == 0 { //Archive
					utils.ErrCheck(cards[cardIdx].Archive(true))
					log.Printf("Card %v archived successfully", cards[cardIdx].Name)
				} else if delMode == 1 { //Delete
					utils.ErrCheck(cards[cardIdx].Delete())
					log.Printf("Card %v deleted successfully", cards[cardIdx].Name)
				}
			}

			gui.Execute(func(gui *Gui) error {
				utils.ErrCheck(gui.DeleteView(mngr.Lists[mngr.currListIdx].Name))
				return nil
			})
			SetKeyBindings(gui, mngr)
		}()
		return nil
	})
}

func addCardMovingFunc(gui *Gui, listName string, mngr *TregoManager) error {
	return gui.SetKeybinding(listName, 'm', ModNone, func(gui *Gui, view *View) error {
		destListC := make(chan int)
		for _, list := range gui.Views() {
			gui.DeleteKeybindings(list.Name())
		}

		listNames := make([]string, len(mngr.Lists))
		for idx, list := range mngr.Lists {
			listNames[idx] = list.Name
		}

		//Used to determine card, that will be moved
		cardIdx := SelectedItemIdx(view)

		utils.ErrCheck(
			mngr.SelectView(
				gui,
				dialog.SelectDialog(
					"Choose destination",
					gui,
					destListC,
					listNames).
						Name()))

		go func() {
			if listIdx, ok := <-destListC; ok {
				cards, err := mngr.Lists[mngr.currListIdx].Cards()
				utils.ErrCheck(err)

				movedCard, err := cards[cardIdx].Move(mngr.Lists[listIdx])
				utils.ErrCheck(err)

				log.Printf("Card %v moved to list: %v", movedCard.Name, mngr.Lists[listIdx].Name)
				gui.Execute(func(gui *Gui) error {
					utils.ErrCheck(
						gui.DeleteView(mngr.Lists[listIdx].Name),
						gui.DeleteView(mngr.Lists[mngr.currListIdx].Name)) //Forces view update
					return nil
				})
			}

			SetKeyBindings(gui, mngr)
		}()

		return nil
	})
}

//Keybinding for switching list on tab keypress
//I used anonymous function for mngr variable access
func addListSwitchingFunc(gui *Gui, viewName string, mngr *TregoManager) (err error) {
	switchListRight := func(gui *Gui, v *View) (err error) {
		mngr.currListIdx = (mngr.currListIdx + 1) % len(mngr.Lists)
		nextViewName := mngr.Lists[mngr.currListIdx].Name

		_, _, x2, _, err := gui.ViewPosition(nextViewName)
		w, _ := gui.Size()
		if x2 > w {
			mngr.listViewOffset -= 1
		} else if mngr.currListIdx == 0 {
			mngr.listViewOffset = 0
		}

		err = mngr.SelectView(gui, nextViewName)
		return
	}
	switchListLeft := func(gui *Gui, v *View) (err error) {
		if mngr.currListIdx == 0 {
			mngr.currListIdx = len(mngr.Lists)
		}
		mngr.currListIdx--
		previousViewName := mngr.Lists[mngr.currListIdx%len(mngr.Lists)].Name

		x1, _, _, _, err := gui.ViewPosition(previousViewName)
		if x1 < 0 {
			mngr.listViewOffset += 1
		} else if mngr.currListIdx == len(mngr.Lists)-1 {
			//Scrolls board to the end
			for mngr.currListIdx = 0; mngr.currListIdx != len(mngr.Lists)-1; {
				utils.ErrCheck(switchListRight(gui, mngr.currView))
			}
		}

		return mngr.SelectView(gui, previousViewName)
	}

	gui.SetKeybinding(viewName, KeyTab, ModNone, switchListRight)
	gui.SetKeybinding(viewName, KeyArrowRight, ModNone, switchListRight)
	gui.SetKeybinding(viewName, KeyArrowLeft, ModNone, switchListLeft)

	return
}

func addCardAddingFunc(gui *Gui, viewName string, mngr *TregoManager) error {
	return gui.SetKeybinding(viewName, 'n', ModNone, func(gui *Gui, view *View) error {
		cardNameC := make(chan string)
		utils.DelNonGlobalKeyBinds(gui)

		utils.ErrCheck(mngr.SelectView(
			gui,
			dialog.InputDialog(
				"Name your card",
				"New card",
				"",
				gui,
				cardNameC).Name()))

		go func() {
			if cardName, ok := <-cardNameC; ok {
				list := mngr.Lists[mngr.currListIdx]
				card, err := list.AddCard(trello.Card{
					IdList: list.Id,
					Name:   cardName,
					Pos:    math.MaxFloat64, //end of the list
				})
				if err != nil {
					log.Panicf("Card add: %v in list %v", err, list.Name)
				}
				log.Printf("Successfully added new card: %v", card.Name)

				gui.Execute(func(gui *Gui) error {
					utils.ErrCheck(gui.DeleteView(list.Name)) //Forces view update
					return nil
				})
			}
			SetKeyBindings(gui, mngr)
		}()
		return nil
	})
}

func addListAddingFunc(gui *Gui, viewName string, mngr *TregoManager) error {
	return gui.SetKeybinding(viewName, 'l', ModNone, func(gui *Gui, view *View) error {
		viewNameC := make(chan string)
		utils.DelNonGlobalKeyBinds(gui)

		utils.ErrCheck(mngr.SelectView(
			gui,
			dialog.InputDialog(
				"Name your list",
				"New list",
				"",
				gui,
				viewNameC).Name()))

		go func() {
			if viewName, ok := <-viewNameC; ok {
				list, err := mngr.CurrBoard.AddList(trello.List{
					Name:    viewName,
					IdBoard: mngr.CurrBoard.Id,
					Pos:     math.MaxFloat32,
				})
				if err != nil {
					log.Printf("List add: %v", err)
				}

				mngr.Lists = append(mngr.Lists, *list)
				utils.ErrCheck(AddList(gui, *list, len(mngr.Lists)-1, mngr.listViewOffset))
			}

			SetKeyBindings(gui, mngr)
		}()
		return nil
	})
}

func addBoardSwitchingFunc(gui *Gui, listName string, mngr *TregoManager) error {
	return gui.SetKeybinding(listName, 'b', ModNone, func(gui *Gui, view *View) error {
		boardSelC := make(chan int)
		utils.DelNonGlobalKeyBinds(gui)

		boards, err := mngr.Member.Boards()
		utils.ErrCheck(err)
		boardNames := make([]string, len(boards))
		for i, board := range boards {
			boardNames[i] = board.Name
		}

		utils.ErrCheck(
			mngr.SelectView(
				gui,
				dialog.SelectDialog(
					"Select board",
					gui,
					boardSelC,
					boardNames,
				).Name()))

		go func() {
			if boardIdx, ok := <-boardSelC; ok {
				mngr.CurrBoard = &boards[boardIdx]
				log.Printf("Changing board to: %v", mngr.CurrBoard.Name)

				for _, list := range mngr.Lists {
					utils.ErrCheck(gui.DeleteView(list.Name))
				}
				utils.ErrCheck(gui.DeleteView(TOP_BAR))
				mngr.Lists = conn.Lists(mngr.CurrBoard)
				mngr.currListIdx = 0
				mngr.listViewOffset = 0
				mngr.currView = nil

				gui.Execute(func(gui *Gui) error {
					mngr.Layout(gui)
					return nil
				})
			}
			SetKeyBindings(gui, mngr)
		}()
		return nil
	})
}

func quit(gui *Gui, v *View) error {
	return ErrQuit
}
