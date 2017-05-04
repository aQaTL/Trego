package ui

import (
	"fmt"
	"github.com/aqatl/Trego/utils"
	. "github.com/jroimartin/gocui"
	"log"
)

const (
	CARD_NAME_VIEW    string = "CARD_EDITOR_VIEW_0"
	CARD_LABELS_VIEW  string = "CARD_EDITOR_VIEW_1"
	CARD_DESC_VIEW    string = "CARD_EDITOR_VIEW_2"
	CARD_ACTIONS_VIEW string = "CARD_EDITOR_VIEW_3"
)

func cardEditorLayout(listView *View, gui *Gui, mngr *TregoManager) {
	mngr.Mode = CARD_EDITOR
	cardIdx := SelectedItemIdx(listView)
	cards, err := mngr.Lists[mngr.currListIdx].Cards()
	utils.ErrCheck(err)
	card := cards[cardIdx]

	w, h := gui.Size()
	if cardNameView, err := gui.SetView(CARD_NAME_VIEW, 0, 3, w-1, 5); err != nil {
		if err != ErrUnknownView {
			utils.ErrCheck(err)
		}

		cardNameView.Wrap = true
		cardNameView.Autoscroll = false
		cardNameView.Editable = false

		fmt.Fprintf(cardNameView, "%v %v %v",
			cyan(card.Name),
			yell("in listView"),
			cyan(mngr.Lists[mngr.currListIdx].Name),
		)

		utils.ErrCheck(
			addEditorViewSwitching(gui, cardNameView, mngr),
			addEditorClosing(gui, cardNameView, mngr),
		)

		utils.ErrCheck(mngr.SelectView(gui, CARD_NAME_VIEW))
	}

	if labelsView, err := gui.SetView(CARD_LABELS_VIEW, 0, 6, w-1, 8); err != nil {
		if err != ErrUnknownView {
			utils.ErrCheck(err)
		}

		labelsView.Title = "Labels"
		labelsView.Editable = false

		for _, label := range card.Labels {
			//TODO Colored labels
			fmt.Fprint(labelsView, label.Name)
		}

		utils.ErrCheck(
			addEditorViewSwitching(gui, labelsView, mngr),
			addEditorClosing(gui, labelsView, mngr),
		)
	}

	if descriptionView, err := gui.SetView(CARD_DESC_VIEW, 0, 9, int(w/3), h-5); err != nil {
		if err != ErrUnknownView {
			utils.ErrCheck(err)
		}

		descriptionView.Wrap = true

		fmt.Fprint(descriptionView, "hejo")

		utils.ErrCheck(
			addEditorViewSwitching(gui, descriptionView, mngr),
			addEditorClosing(gui, descriptionView, mngr),
		)
	}

	if actionsView, err := gui.SetView(CARD_ACTIONS_VIEW, int(w/3)+1, 9, w-1, h-5); err != nil {
		if err != ErrUnknownView {
			utils.ErrCheck(err)
		}

		actionsView.Wrap = true
		actionsView.Title = "Activities"

		actions, reqErr := card.Actions()
		utils.ErrCheck(reqErr)
		for _, action := range actions {
			if action.Data.Text != "" {
				fmt.Fprintln(actionsView, action.Date, action.Type, action.Data.Text)
			}
		}

		utils.ErrCheck(
			addEditorViewSwitching(gui, actionsView, mngr),
			addEditorClosing(gui, actionsView, mngr),
		)
	}
}

func addEditorViewSwitching(gui *Gui, view *View, mngr *TregoManager) error {
	return gui.SetKeybinding(view.Name(), KeyTab, ModNone, func(gui *Gui, view *View) error {
		idx := (view.Name()[17] - 48 + 1) % 4
		nextViewName := view.Name()[:17] + string(rune(idx+48))
		log.Printf("switching editor view to %v", nextViewName)
		utils.ErrCheck(mngr.SelectView(gui, nextViewName))
		return nil
	})
}

func addEditorClosing(gui *Gui, view *View, mngr *TregoManager) error {
	return gui.SetKeybinding(view.Name(), KeyCtrlQ, ModNone, func(gui *Gui, view *View) error {
		gui.DeleteKeybindings(CARD_NAME_VIEW)
		gui.DeleteKeybindings(CARD_DESC_VIEW)
		gui.DeleteKeybindings(CARD_ACTIONS_VIEW)
		gui.DeleteKeybindings(CARD_LABELS_VIEW)
		utils.ErrCheck(
			gui.DeleteView(CARD_NAME_VIEW),
			gui.DeleteView(CARD_DESC_VIEW),
			gui.DeleteView(CARD_ACTIONS_VIEW),
			gui.DeleteView(CARD_LABELS_VIEW),
		)
		mngr.Mode = BOARD_VIEW
		mngr.currView = nil
		return nil
	})
}
