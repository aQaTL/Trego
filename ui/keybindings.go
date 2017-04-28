package ui

import (
	"fmt"
	"github.com/aqatl/Trego/conn"
	"github.com/aqatl/Trego/ui/dialog"
	"github.com/aqatl/Trego/utils"
	"github.com/aqatl/go-trello"
	"github.com/fatih/color"
	. "github.com/jroimartin/gocui"
	"log"
	"math"
	"strings"
)

func SetKeyBindings(gui *Gui, mngr *TregoManager) (err error) {
	if err = gui.SetKeybinding("", KeyCtrlC, ModNone, quit); err != nil {
		return
	}

	utils.ErrCheck(addListAdding(gui, TOP_BAR, mngr))

	for _, list := range mngr.Lists {
		utils.ErrCheck(
			gui.SetKeybinding(list.Id, KeyArrowUp, ModNone, utils.CursorUp),
			gui.SetKeybinding(list.Id, KeyArrowDown, ModNone, utils.CursorDown),
			addListSwitching(gui, list.Id, mngr),
			addListAdding(gui, list.Id, mngr),
			addListMoving(gui, list.Id, mngr),
			addCardAdding(gui, list.Id, mngr),
			addCardToListMoving(gui, list.Id, mngr),
			addDeleting(gui, list.Id, mngr),
			addBoardSwitching(gui, list.Id, mngr),
			addCardSearching(gui, list.Id, mngr),
			addCardMoving(gui, list.Id, mngr),
		)

		gui.SetKeybinding(list.Id, 'q', ModNone, func(gui *Gui, view *View) error {
			log.Print("=========")
			cards, err := mngr.Lists[mngr.currListIdx].FreshCards()
			utils.ErrCheck(err)
			for _, card := range cards {
				log.Printf("%v: %f", card.Name, card.Pos)
			}
			return nil
		})
	}
	return
}

func addCardSearching(gui *Gui, listName string, mngr *TregoManager) error {
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

func addDeleting(gui *Gui, listName string, mngr *TregoManager) error {
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

func addCardMoving(gui *Gui, listName string, mngr *TregoManager) error {
	queue := make(chan trello.Card)
	go func(queue chan trello.Card) {
		for {
			card := <-queue
			_, err := card.Move(fmt.Sprintf("%f", card.Pos))
			utils.ErrCheck(err)
		}
	}(queue)

	err := gui.SetKeybinding(listName, KeyCtrlU, ModNone, func(gui *Gui, view *View) error {
		cardIdx := SelectedItemIdx(view)
		cards, err := mngr.Lists[mngr.currListIdx].Cards()
		utils.ErrCheck(err)

		if cardIdx < 0 || len(cards) < 2 {
			return nil
		}
		prevCardIdx := cardIdx - 1
		if prevCardIdx < 0 {
			prevCardIdx = len(cards) - 1
		}

		if prevCardIdx == len(cards)-1 {
			cards[cardIdx].Pos = cards[prevCardIdx].Pos + (1 << 16)

			currCard := cards[cardIdx]
			copy(cards, cards[1:])
			cards[prevCardIdx] = currCard
		} else if prevCardIdx == 0 {
			cards[cardIdx].Pos = cards[0].Pos / 2

			cards[cardIdx], cards[prevCardIdx] = cards[prevCardIdx], cards[cardIdx]
		} else {
			cards[cardIdx].Pos =
					cards[prevCardIdx].Pos -
							(cards[prevCardIdx].Pos-cards[cardIdx-2].Pos)/2

			cards[cardIdx], cards[prevCardIdx] = cards[prevCardIdx], cards[cardIdx]
		}

		queue <- cards[prevCardIdx]

		color.Output = view
		view.Clear()
		for idx, card := range cards {
			defaultCardColor.Printf("%d.%v\n", idx, card.Name)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return gui.SetKeybinding(listName, KeyCtrlD, ModNone, func(gui *Gui, view *View) error {
		cardIdx := SelectedItemIdx(view)
		cards, err := mngr.Lists[mngr.currListIdx].Cards()
		utils.ErrCheck(err)

		if cardIdx < 0 || len(cards) < 2 {
			return nil
		}
		nextCardIdx := (cardIdx + 1) % len(cards)

		if nextCardIdx == 0 {
			cards[cardIdx].Pos = math.Ceil(cards[0].Pos / 2)

			currCard := cards[cardIdx]
			copy(cards[1:], cards[:len(cards)-1])
			cards[0] = currCard
		} else if nextCardIdx == len(cards)-1 {
			cards[cardIdx].Pos = cards[nextCardIdx].Pos + (1 << 16)

			cards[cardIdx], cards[nextCardIdx] = cards[nextCardIdx], cards[cardIdx]
		} else {
			cards[cardIdx].Pos =
					cards[nextCardIdx].Pos +
							(cards[cardIdx+2%len(cards)].Pos -
									cards[nextCardIdx].Pos) / 2

			cards[cardIdx], cards[nextCardIdx] = cards[nextCardIdx], cards[cardIdx]
		}

		queue <- cards[nextCardIdx]

		color.Output = view
		view.Clear()
		for idx, card := range cards {
			defaultCardColor.Printf("%d.%v\n", idx, card.Name)
		}
		return nil
	})
}

func addCardToListMoving(gui *Gui, listName string, mngr *TregoManager) error {
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

				movedCard, err := cards[cardIdx].MoveToList(mngr.Lists[listIdx])
				utils.ErrCheck(err)
				movedCard, err = movedCard.Move("bottom")
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
func addListSwitching(gui *Gui, viewName string, mngr *TregoManager) (err error) {
	utils.ErrCheck(gui.SetKeybinding(viewName, KeyTab, ModNone, mngr.SwitchListRight),
		gui.SetKeybinding(viewName, KeyArrowRight, ModNone, mngr.SwitchListRight),
		gui.SetKeybinding(viewName, KeyArrowLeft, ModNone, mngr.SwitchListLeft))
	return
}

func addListMoving(gui *Gui, viewName string, mngr *TregoManager) error {
	queue := make(chan trello.List)
	go func(queue chan trello.List) {
		for {
			list := <-queue
			list.Move(fmt.Sprintf("%f", list.Pos))
		}
	}(queue)

	err := gui.SetKeybinding(viewName, KeyCtrlN, ModNone, func(gui *Gui, view *View) error {
		if len(mngr.Lists) < 2 {
			return nil
		}

		currIdx := mngr.currListIdx
		nextListIdx := (currIdx + 1) % len(mngr.Lists)
		if nextListIdx == 0 {
			mngr.Lists[currIdx].Pos = float32(math.Ceil(float64(mngr.Lists[0].Pos) / 2))

			currList := mngr.Lists[currIdx]
			copy(mngr.Lists[1:], mngr.Lists[:len(mngr.Lists)-1])
			mngr.Lists[0] = currList
		} else if nextListIdx == len(mngr.Lists)-1 {
			mngr.Lists[currIdx].Pos = mngr.Lists[nextListIdx].Pos + (1 << 16)

			mngr.Lists[currIdx], mngr.Lists[nextListIdx] =
					mngr.Lists[nextListIdx], mngr.Lists[currIdx]
		} else {
			mngr.Lists[currIdx].Pos =
					mngr.Lists[nextListIdx].Pos +
							(mngr.Lists[(currIdx+2)%len(mngr.Lists)].Pos -
									mngr.Lists[nextListIdx].Pos) / 2

			mngr.Lists[currIdx], mngr.Lists[nextListIdx] =
					mngr.Lists[nextListIdx], mngr.Lists[currIdx]
		}

		queue <- mngr.Lists[nextListIdx]

		utils.ErrCheck(mngr.SwitchListRight(gui, view))
		if nextListIdx == 0 {
			mngr.listViewOffset = 0
		}

		return nil
	})

	if err != nil {
		return err
	}

	return gui.SetKeybinding(viewName, KeyCtrlP, ModNone, func(gui *Gui, view *View) error {
		if len(mngr.Lists) < 2 {
			return nil
		}
		prevListIdx := mngr.currListIdx - 1
		if prevListIdx < 0 {
			prevListIdx = len(mngr.Lists) - 1
		}

		currIdx := mngr.currListIdx
		if prevListIdx == len(mngr.Lists)-1 {
			mngr.Lists[currIdx].Pos = mngr.Lists[prevListIdx].Pos + (1 << 16)

			currList := mngr.Lists[currIdx]
			copy(mngr.Lists, mngr.Lists[1:])
			mngr.Lists[prevListIdx] = currList
		} else if prevListIdx == 0 {
			mngr.Lists[currIdx].Pos = mngr.Lists[prevListIdx].Pos / 2

			mngr.Lists[currIdx], mngr.Lists[prevListIdx] =
					mngr.Lists[prevListIdx], mngr.Lists[currIdx]
		} else {
			mngr.Lists[currIdx].Pos =
					mngr.Lists[prevListIdx].Pos -
							(mngr.Lists[prevListIdx].Pos-mngr.Lists[currIdx-2].Pos)/2

			mngr.Lists[currIdx], mngr.Lists[prevListIdx] =
					mngr.Lists[prevListIdx], mngr.Lists[currIdx]
		}

		queue <- mngr.Lists[prevListIdx]

		utils.ErrCheck(mngr.Layout(gui), mngr.SwitchListLeft(gui, view))

		return nil
	})
}

func addCardAdding(gui *Gui, viewName string, mngr *TregoManager) error {
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

func addListAdding(gui *Gui, viewName string, mngr *TregoManager) error {
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

func addBoardSwitching(gui *Gui, listName string, mngr *TregoManager) error {
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
