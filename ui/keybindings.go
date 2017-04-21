package ui

import (
	. "github.com/jroimartin/gocui"
	"github.com/aqatl/go-trello"
	"github.com/aqatl/Trego/utils"
	"github.com/aqatl/Trego/ui/dialog"
	"log"
	"math"
	"github.com/aqatl/Trego/conn"
	"strings"
	"fmt"
)

func SetKeyBindings(gui *Gui, mngr *TregoManager) (err error) {
	if err = gui.SetKeybinding("", KeyCtrlC, ModNone, quit); err != nil {
		return
	}

	utils.ErrCheck(addListAddingFunc(gui, TOP_BAR, mngr))

	for _, list := range mngr.Lists {
		utils.ErrCheck(
			gui.SetKeybinding(list.Id, KeyArrowUp, ModNone, utils.CursorUp),
			gui.SetKeybinding(list.Id, KeyArrowDown, ModNone, utils.CursorDown),
			addListSwitchingFunc(gui, list.Id, mngr),
			addListAddingFunc(gui, list.Id, mngr),
			addCardAddingFunc(gui, list.Id, mngr),
			addCardMovingFunc(gui, list.Id, mngr),
			addDeletingFunc(gui, list.Id, mngr),
			addBoardSwitchingFunc(gui, list.Id, mngr),
			addCardSearchingFunc(gui, list.Id, mngr),
		)
	}
	return
}
func addCardSearchingFunc(gui *Gui, listName string, mngr *TregoManager) error {
	return gui.SetKeybinding(listName, 's', ModNone, func(gui *Gui, listView *View) error {

		x1, _, x2, _, err1 := gui.ViewPosition(listName)
		_, y1, _, _, err2 := gui.ViewPosition(BOTTOM_BAR)
		utils.ErrCheck(err1, err2)
		if searchView, err := gui.SetView(SEARCH_VIEW, x1, y1-3, x2, y1-1); err != nil {
			if err != ErrUnknownView {
				return err
			}

			searchView.Highlight = true
			searchView.Wrap = false
			searchView.Editable = true

			list := mngr.Lists[mngr.currListIdx]
			cards, err := list.Cards()
			utils.ErrCheck(err)

			searchView.Editor = EditorFunc(func(v *View, key Key, ch rune, mod Modifier) {
				switch {
				case ch != 0 && mod == 0:
					v.EditWrite(ch)
				case key == KeySpace:
					v.EditWrite(' ')
				case key == KeyBackspace || key == KeyBackspace2:
					v.EditDelete(true)
				case key == KeyDelete:
					v.EditDelete(false)
				case key == KeyInsert:
					v.Overwrite = !v.Overwrite
					return
				case key == KeyArrowDown:
					v.MoveCursor(0, 1, false)
					return
				case key == KeyArrowUp:
					v.MoveCursor(0, -1, false)
					return
				case key == KeyArrowLeft:
					v.MoveCursor(-1, 0, false)
					return
				case key == KeyArrowRight:
					v.MoveCursor(1, 0, false)
					return
				}
				if len(v.Buffer()) != 0 {
					listView.Clear()
					for idx, card := range cards {
						if strings.Contains(card.Name, v.Buffer()[:len(v.Buffer())-2]) {
							fmt.Fprintf(listView, "%d.%v\n", idx, card.Name)
						}
					}
				}
			})

			utils.ErrCheck(gui.SetKeybinding(searchView.Name(), KeyEnter, ModNone, func(gui *Gui, v *View) error {
				gui.DeleteKeybindings(SEARCH_VIEW)
				utils.ErrCheck(gui.DeleteView(SEARCH_VIEW))
				return nil
			}))
			utils.ErrCheck(mngr.SelectView(gui, SEARCH_VIEW))
		}

		return nil
	})
}

func addDeletingFunc(gui *Gui, listName string, mngr *TregoManager) error {
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
					[]string{"Archive", "Delete", "Archive list"}).
						Name()))

		go func() {
			currList := mngr.Lists[mngr.currListIdx]
			if delMode, ok := <-delModeC; ok {
				cards, err := currList.Cards()
				utils.ErrCheck(err)
				cardIdx := SelectedItemIdx(view)
				if delMode == 0 && cardIdx >= 0 { //Archive
					utils.ErrCheck(cards[cardIdx].Archive(true))
					log.Printf("Card %v archived successfully", cards[cardIdx].Name)
				} else if delMode == 1 && cardIdx >= 0 { //Delete
					utils.ErrCheck(cards[cardIdx].Delete())
					log.Printf("Card %v deleted successfully", cards[cardIdx].Name)
				} else if delMode == 2 { //Archive list
					utils.ErrCheck(currList.Archive(true))
					log.Printf("List %v archived successfully", currList.Name)
					mngr.Lists = utils.RemoveList(mngr.Lists, mngr.currListIdx)
					mngr.currListIdx = 0
					mngr.listViewOffset = 0
				}
			}

			gui.Execute(func(gui *Gui) error {
				utils.ErrCheck(gui.DeleteView(currList.Id))
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
						gui.DeleteView(mngr.Lists[listIdx].Id),
						gui.DeleteView(mngr.Lists[mngr.currListIdx].Id)) //Forces view update
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
		nextViewId := mngr.Lists[mngr.currListIdx].Id

		_, _, x2, _, err := gui.ViewPosition(nextViewId)
		w, _ := gui.Size()
		if x2 > w {
			mngr.listViewOffset -= 1
		} else if mngr.currListIdx == 0 {
			mngr.listViewOffset = 0
		}

		err = mngr.SelectView(gui, nextViewId)
		return
	}
	switchListLeft := func(gui *Gui, v *View) (err error) {
		if mngr.currListIdx == 0 {
			mngr.currListIdx = len(mngr.Lists)
		}
		mngr.currListIdx--
		previousViewId := mngr.Lists[mngr.currListIdx%len(mngr.Lists)].Id

		x1, _, _, _, err := gui.ViewPosition(previousViewId)
		if x1 < 0 {
			mngr.listViewOffset += 1
		} else if mngr.currListIdx == len(mngr.Lists)-1 {
			//Scrolls board to the end
			for mngr.currListIdx = 0; mngr.currListIdx != len(mngr.Lists)-1; {
				utils.ErrCheck(switchListRight(gui, mngr.currView))
			}
		}

		return mngr.SelectView(gui, previousViewId)
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
					utils.ErrCheck(gui.DeleteView(list.Id)) //Forces view update
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
				gui.Execute(func(gui *Gui) error {
					utils.ErrCheck(AddList(gui, *list, len(mngr.Lists)-1, mngr.listViewOffset))
					mngr.SelectView(gui, list.Id)
					mngr.currListIdx = len(mngr.Lists) - 1
					return nil
				})
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
					utils.ErrCheck(gui.DeleteView(list.Id))
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
