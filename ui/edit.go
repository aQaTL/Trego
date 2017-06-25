package ui

import (
	"fmt"
	"github.com/aqatl/Trego/utils"
	"github.com/aqatl/go-trello"
	"github.com/fatih/color"
	. "github.com/jroimartin/gocui"
	"log"
	"time"
	"unicode/utf8"
	"github.com/aqatl/Trego/ui/dialog"
)

const (
	cardNameView     string = "CARD_EDITOR_VIEW_0"
	cardLabelsView   string = "CARD_EDITOR_VIEW_1"
	cardDescView     string = "CARD_EDITOR_VIEW_2"
	cardCommentsView string = "CARD_EDITOR_VIEW_3"
	cardListInfoView string = "CARD_EDITOR_VIEW_4"
)

type CardEditor struct {
	Mngr *TregoManager
	Card *trello.Card

	currView *View
}

func (cEdit *CardEditor) Layout(gui *Gui) error {
	utils.ErrCheck(
		cEdit.nameView(gui),
		cEdit.listInfoView(gui),
		cEdit.labelsView(gui),
		cEdit.descriptionView(gui),
		cEdit.commentsView(gui),
	)
	if _, err := gui.SetCurrentView(cEdit.currView.Name()); err != nil {
		log.Printf("Error with %v view: %v", cEdit.currView.Name(), err)
		return err
	}
	return nil
}

func (cEdit *CardEditor) nameView(gui *Gui) (err error) {
	w, _ := gui.Size()
	if nameView, err := gui.SetView(cardNameView, 0, 3, w*2/3-1, 5); err != nil {
		if err != ErrUnknownView {
			return err
		}

		nameView.Title = "Card name"
		nameView.Wrap = false
		nameView.Autoscroll = false
		nameView.Editable = true
		nameView.FgColor = ColorRed

		fmt.Fprint(nameView, cEdit.Card.Name)

		utils.ErrCheck(
			addEditorViewSwitching(gui, nameView, cEdit),
			addEditorClosing(gui, nameView, cEdit),
			addChangesSaving(gui, nameView, cEdit),
		)

		cEdit.Mngr.BotBar.CurrBotBarKey = cardNameView
		cEdit.currView = nameView
		gui.DeleteView(BottomBar) //Error intentionally ignored
	}
	return
}

func (cEdit *CardEditor) listInfoView(gui *Gui) (err error) {
	w, _ := gui.Size()
	if listInfoView, err := gui.SetView(cardListInfoView, w*2/3, 3, w-1, 5); err != nil {
		if err != ErrUnknownView {
			return err
		}

		listInfoView.Title = "Card in list:"
		listInfoView.Editable = false

		fmt.Fprint(listInfoView, yell(cEdit.Mngr.Lists[cEdit.Mngr.currListIdx].Name))
	}
	return
}

func (cEdit *CardEditor) labelsView(gui *Gui) (err error) {
	w, _ := gui.Size()
	if view, err := gui.SetView(cardLabelsView, 0, 6, w-1, 8); err != nil {
		if err != ErrUnknownView {
			return err
		}

		view.Title = "Labels"
		view.Editable = true

		labelsLens := make([]int, len(cEdit.Card.Labels))

		var currLabel *trello.Label
		if len(cEdit.Card.Labels) > 0 {
			currLabel = &cEdit.Card.Labels[0]
		}

		view.Editor = EditorFunc(func(view *View, key Key, ch rune, mod Modifier) {
			switch key {
			case KeyArrowRight:
				cx, cy := view.Cursor()
				bufferLen := utf8.RuneCountInString(view.Buffer())
				sum := 0

				for i, labelLen := range labelsLens {
					sum += labelLen + 1
					if sum == cx+labelLen+1 && sum < bufferLen-1 {
						utils.ErrCheck(view.SetCursor(sum, cy))
						currLabel = &cEdit.Card.Labels[i+1]
						break
					}
				}
			case KeyArrowLeft:
				cx, cy := view.Cursor()
				if cx == 0 {
					return
				}

				sum := 0
				for i, labelLen := range labelsLens {
					sum += labelLen + 1
					if sum == cx {
						utils.ErrCheck(view.SetCursor(cx-labelLen-1, cy))
						currLabel = &cEdit.Card.Labels[i]
						break
					}
				}
			}
		})

		for i, label := range cEdit.Card.Labels {
			if label.Name == "" {
				label.Name = "\u2588\u2588\u2588\u2588\u2588"
			}
			col, hi := utils.MapColor(label.Color)
			fmt.Fprintf(view, "\033[3%d;%dm%v\033[0m ", col, hi, label.Name)
			labelsLens[i] = utf8.RuneCountInString(label.Name)
		}

		utils.ErrCheck(
			addEditorViewSwitching(gui, view, cEdit),
			addEditorClosing(gui, view, cEdit),
		)

		gui.SetKeybinding(cardLabelsView, 'd', ModNone, func(gui *Gui, view *View) error {
			if *currLabel == *new(trello.Label) {
				return nil
			}
			utils.ErrCheck(currLabel.DeleteLabel())
			log.Printf("Successfully deleted label %v with color: %v", currLabel.Name, currLabel.Color)
			for i, label := range cEdit.Card.Labels {
				if label.Id == currLabel.Id {
					cEdit.Card.Labels = utils.RemoveLabel(cEdit.Card.Labels, i)
					labelsLens = utils.RemoveInt(labelsLens, i)
					if len(cEdit.Card.Labels) >= 2 {
						currLabel = &cEdit.Card.Labels[0]
					} else {
						currLabel = new(trello.Label)
					}
					break
				}
			}

			view.Clear()
			gui.DeleteKeybindings(cardLabelsView)
			utils.ErrCheck(
				gui.DeleteView(cardLabelsView),
				cEdit.labelsView(gui),
			)
			cEdit.currView = view

			return nil
		})

		gui.SetKeybinding(cardLabelsView, 'n', ModNone, func(gui *Gui, view *View) error {
			labelChan := make(chan [2]string)

			gui.SetManager(dialog.LabelDialog(gui, labelChan), cEdit.Mngr.BotBar, cEdit.Mngr.TopBar)
			cEdit.Mngr.BotBar.CurrBotBarKey = "label_dialog"

			go func(gui *Gui, cEdit *CardEditor, labelChan chan [2]string) {
				if labelArr, ok := <-labelChan; ok {
					label, err := cEdit.Card.AddNewLabel(labelArr[0], labelArr[1])
					utils.ErrCheck(err)

					cEdit.Card.Labels = append(cEdit.Card.Labels, *label)
					labelsLens = append(labelsLens, utf8.RuneCountInString(label.Name))
				}

				gui.Execute(func(gui *Gui) error {
					gui.SetManager(cEdit, cEdit.Mngr.BotBar, cEdit.Mngr.TopBar)
					return nil
				})
			}(gui, cEdit, labelChan)
			return nil
		})

		gui.SetKeybinding(cardLabelsView, 'r', ModNone, func(gui *Gui, view *View) error {
			inputChan := make(chan string)

			inputDialogViews := dialog.InputDialog(
				"Rename \""+currLabel.Name+"\"",
				"Rename",
				"",
				gui,
				inputChan,
			)
			cEdit.currView = inputDialogViews[0]

			go func(gui *Gui, inputChan <-chan string) {
				if newName, ok := <-inputChan; ok {
					utils.ErrCheck(currLabel.UpdateName(newName))
				}

				gui.Execute(func(gui *Gui) error {
					cEdit.currView = view
					dialog.DeleteDialog(gui, inputDialogViews[:]...)
					utils.ErrCheck(gui.DeleteView(cardLabelsView))
					return nil
				})
			}(gui, inputChan)

			return nil
		})

		gui.SetKeybinding(cardLabelsView, 'a', ModNone, func(gui *Gui, view *View) error {
			labels, err := cEdit.Mngr.CurrBoard.Labels()
			utils.ErrCheck(err)

			selIdxChan := make(chan int)
			values := make([]string, len(labels))

			for i, label := range labels {
				if label.Name == "" {
					label.Name = "\u2588\u2588\u2588\u2588\u2588"
				}
				col, hi := utils.MapColor(label.Color)
				values[i] = fmt.Sprintf("\033[3%d;%dm%v\033[0m", col, hi, label.Name)
			}

			selectDialog := dialog.SelectDialog(
				"Select label",
				gui,
				selIdxChan,
				values,
			)
			cEdit.currView = selectDialog

			go func(selIdxChan chan int, selectDialog *View) {
				if idx, ok := <-selIdxChan; ok {
					_, err := cEdit.Card.AddLabel(labels[idx].Id)
					utils.ErrCheck(err)
					cEdit.Card.Labels = append(cEdit.Card.Labels, labels[idx])
				}

				gui.Execute(func(gui *Gui) error {
					cEdit.currView = view
					dialog.DeleteDialog(gui, selectDialog)
					utils.ErrCheck(gui.DeleteView(cardLabelsView))
					return nil
				})
			}(selIdxChan, selectDialog)

			return nil
		})
	}
	return
}

func (cEdit *CardEditor) descriptionView(gui *Gui) (err error) {
	w, h := gui.Size()
	if descriptionView, err := gui.SetView(cardDescView, 0, 9, int(w/3), h-5); err != nil {
		if err != ErrUnknownView {
			return err
		}

		descriptionView.Title = "Description"
		descriptionView.Editable = true
		descriptionView.Wrap = true

		fmt.Fprint(descriptionView, cEdit.Card.Desc)

		utils.ErrCheck(
			addEditorViewSwitching(gui, descriptionView, cEdit),
			addEditorClosing(gui, descriptionView, cEdit),
			addChangesSaving(gui, descriptionView, cEdit),
		)
	}
	return
}

func (cEdit *CardEditor) commentsView(gui *Gui) (err error) {
	w, h := gui.Size()
	if commentsView, err := gui.SetView(cardCommentsView, int(w/3)+1, 9, w-1, h-5); err != nil {
		if err != ErrUnknownView {
			return err
		}

		commentsView.Title = "Comments"
		commentsView.Wrap = true
		commentsView.Editable = false

		actions, reqErr := cEdit.Card.Actions()
		utils.ErrCheck(reqErr)
		for _, action := range actions {
			if action.Type == trello.CommentCard {
				creationTime, err := time.Parse(time.RFC3339, action.Date)
				utils.ErrCheck(err)
				fmt.Fprintf(
					commentsView,
					"%v %v\n> %v\n",
					color.BlueString("%v", creationTime.Format("2.01.2006 15:04")),
					yell(action.MemberCreator.Username),
					cyan(action.Data.Text),
				)
			}
		}

		utils.ErrCheck(
			addEditorViewSwitching(gui, commentsView, cEdit),
			addEditorClosing(gui, commentsView, cEdit),
		)
	}
	return
}

func addEditorViewSwitching(gui *Gui, view *View, cEdit *CardEditor) error {
	return gui.SetKeybinding(view.Name(), KeyTab, ModNone, func(gui *Gui, view *View) error {
		idx := (view.Name()[17] - 48 + 1) % 4
		nextViewName := view.Name()[:17] + string(rune(idx+48))

		//log.Printf("switching editor view to %v", nextViewName)

		nextView, err := gui.View(nextViewName)
		utils.ErrCheck(err)
		cEdit.currView = nextView

		cEdit.Mngr.BotBar.CurrBotBarKey = nextViewName

		gui.DeleteView(BottomBar)
		return nil
	})
}

func addEditorClosing(gui *Gui, view *View, cEdit *CardEditor) error {
	return gui.SetKeybinding(view.Name(), KeyCtrlQ, ModNone, func(gui *Gui, view *View) error {
		gui.DeleteKeybindings(cardNameView)
		gui.DeleteKeybindings(cardListInfoView)
		gui.DeleteKeybindings(cardLabelsView)
		gui.DeleteKeybindings(cardDescView)
		gui.DeleteKeybindings(cardCommentsView)

		gui.DeleteView(cardNameView)
		gui.DeleteView(cardListInfoView)
		gui.DeleteView(cardLabelsView)
		gui.DeleteView(cardDescView)
		gui.DeleteView(cardCommentsView)

		cEdit.Mngr.BotBar.CurrBotBarKey = cEdit.Mngr.BotBar.DefaultBotBarKey
		gui.SetManager(cEdit.Mngr, cEdit.Mngr.BotBar, cEdit.Mngr.TopBar)
		SetKeyBindings(gui, cEdit.Mngr)
		return nil
	})
}

func addChangesSaving(gui *Gui, view *View, cEdit *CardEditor) error {
	return gui.SetKeybinding(view.Name(), KeyCtrlS, ModNone, func(gui *Gui, view *View) (err error) {
		switch view.Name() {
		case cardNameView:
			_, err = cEdit.Card.SetName(view.Buffer()[:len(view.Buffer())-1])
		case cardDescView:
			_, err = cEdit.Card.SetDescription(view.Buffer())
		default:
			//Unsupported view
		}
		utils.ErrCheck(err)
		return nil
	})
}
