package ui

import (
	"fmt"
	"github.com/aqatl/Trego/utils"
	. "github.com/jroimartin/gocui"
	"log"
	"github.com/aqatl/go-trello"
	"github.com/fatih/color"
	"time"
)

const (
	CARD_NAME_VIEW      string = "CARD_EDITOR_VIEW_0"
	CARD_LABELS_VIEW    string = "CARD_EDITOR_VIEW_1"
	CARD_DESC_VIEW      string = "CARD_EDITOR_VIEW_2"
	CARD_COMMENTS_VIEW  string = "CARD_EDITOR_VIEW_3"
	CARD_LIST_INFO_VIEW string = "CARD_EDITOR_VIEW_4"
)

func CardEditorLayout(listView *View, gui *Gui, mngr *TregoManager) {
	mngr.Mode = CARD_EDITOR
	cardIdx := SelectedItemIdx(listView)
	cards, err := mngr.Lists[mngr.currListIdx].Cards()
	utils.ErrCheck(err)
	card := cards[cardIdx]

	w, h := gui.Size()
	if cardNameView, err := gui.SetView(CARD_NAME_VIEW, 0, 3, w*2/3-1, 5); err != nil {
		if err != ErrUnknownView {
			utils.ErrCheck(err)
		}

		cardNameView.Title = "Card name"
		cardNameView.Wrap = false
		cardNameView.Autoscroll = false
		cardNameView.Editable = true
		cardNameView.FgColor = ColorRed

		fmt.Fprint(cardNameView, card.Name)

		utils.ErrCheck(
			addEditorViewSwitching(gui, cardNameView, mngr),
			addEditorClosing(gui, cardNameView, mngr),
			addChangesSaving(gui, cardNameView, mngr, &card),
		)

		utils.ErrCheck(mngr.SelectView(gui, CARD_NAME_VIEW))
	}

	if listInfoView, err := gui.SetView(CARD_LIST_INFO_VIEW, w*2/3, 3, w-1, 5); err != nil {
		if err != ErrUnknownView {
			utils.ErrCheck(err)
		}

		listInfoView.Title = "Card in list:"
		listInfoView.Editable = false

		fmt.Fprint(listInfoView, yell(mngr.Lists[mngr.currListIdx].Name))
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

		descriptionView.Title = "Description"
		descriptionView.Editable = true
		descriptionView.Wrap = true

		fmt.Fprint(descriptionView, card.Desc)

		utils.ErrCheck(
			addEditorViewSwitching(gui, descriptionView, mngr),
			addEditorClosing(gui, descriptionView, mngr),
			addChangesSaving(gui, descriptionView, mngr, &card),
		)
	}

	if commentsView, err := gui.SetView(CARD_COMMENTS_VIEW, int(w/3)+1, 9, w-1, h-5); err != nil {
		if err != ErrUnknownView {
			utils.ErrCheck(err)
		}

		commentsView.Title = "Comments"
		commentsView.Wrap = true
		commentsView.Editable = false

		actions, reqErr := card.Actions()
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
			addEditorViewSwitching(gui, commentsView, mngr),
			addEditorClosing(gui, commentsView, mngr),
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
		gui.DeleteKeybindings(CARD_LIST_INFO_VIEW)
		gui.DeleteKeybindings(CARD_LABELS_VIEW)
		gui.DeleteKeybindings(CARD_DESC_VIEW)
		gui.DeleteKeybindings(CARD_COMMENTS_VIEW)
		utils.ErrCheck(
			gui.DeleteView(CARD_NAME_VIEW),
			gui.DeleteView(CARD_LIST_INFO_VIEW),
			gui.DeleteView(CARD_LABELS_VIEW),
			gui.DeleteView(CARD_DESC_VIEW),
			gui.DeleteView(CARD_COMMENTS_VIEW),
		)
		mngr.Mode = BOARD_VIEW
		mngr.currView = nil
		return nil
	})
}

func addChangesSaving(gui *Gui, view *View, mngr *TregoManager, card *trello.Card) error {
	return gui.SetKeybinding(view.Name(), KeyCtrlS, ModNone, func(gui *Gui, view *View) (err error) {
		switch view.Name() {
		case CARD_NAME_VIEW:
			_, err = card.SetName(view.Buffer()[:len(view.Buffer())-2])
		case CARD_DESC_VIEW:
			_, err = card.SetDescription(view.Buffer())
		default:
			//Unsupported view
		}
		utils.ErrCheck(err)
		return nil
	})
}
