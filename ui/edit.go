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
	cardNameView     string = "CARD_EDITOR_VIEW_0"
	cardLabelsView   string = "CARD_EDITOR_VIEW_1"
	cardDescView     string = "CARD_EDITOR_VIEW_2"
	cardCommentsView string = "CARD_EDITOR_VIEW_3"
	cardListInfoView string = "CARD_EDITOR_VIEW_4"
)

func CardEditorLayout(listView *View, gui *Gui, mngr *TregoManager) {
	mngr.Mode = CardEditor
	cardIdx := SelectedItemIdx(listView)
	cards, err := mngr.Lists[mngr.currListIdx].Cards()
	utils.ErrCheck(err)
	card := cards[cardIdx]

	utils.ErrCheck(
		nameView(gui, mngr, &card),
		listInfoView(gui, mngr),
		labelsView(gui, mngr, &card),
		descriptionView(gui, mngr, &card),
		commentsView(gui, mngr, &card),
	)
}

func nameView(gui *Gui, mngr *TregoManager, card *trello.Card) (err error) {
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

		fmt.Fprint(nameView, card.Name)

		utils.ErrCheck(
			addEditorViewSwitching(gui, nameView, mngr),
			addEditorClosing(gui, nameView, mngr),
			addChangesSaving(gui, nameView, mngr, card),
		)

		utils.ErrCheck(mngr.SelectView(gui, cardNameView))
		mngr.CurrBotBarKey = cardNameView
		utils.ErrCheck(gui.DeleteView(BottomBar))
	}
	return
}

func listInfoView(gui *Gui, mngr *TregoManager) (err error) {
	w, _ := gui.Size()
	if listInfoView, err := gui.SetView(cardListInfoView, w*2/3, 3, w-1, 5); err != nil {
		if err != ErrUnknownView {
			return err
		}

		listInfoView.Title = "Card in list:"
		listInfoView.Editable = false

		fmt.Fprint(listInfoView, yell(mngr.Lists[mngr.currListIdx].Name))
	}
	return
}

func labelsView(gui *Gui, mngr *TregoManager, card *trello.Card) (err error) {
	w, _ := gui.Size()
	if labelsView, err := gui.SetView(cardLabelsView, 0, 6, w-1, 8); err != nil {
		if err != ErrUnknownView {
			return err
		}

		labelsView.Title = "Labels"
		labelsView.Editable = false

		for _, label := range card.Labels {
			//TODO Colored labels
			fmt.Fprint(labelsView, label.Name, label.Color)

		}

		utils.ErrCheck(
			addEditorViewSwitching(gui, labelsView, mngr),
			addEditorClosing(gui, labelsView, mngr),
		)
	}
	return
}

func descriptionView(gui *Gui, mngr *TregoManager, card *trello.Card) (err error) {
	w, h := gui.Size()
	if descriptionView, err := gui.SetView(cardDescView, 0, 9, int(w/3), h-5); err != nil {
		if err != ErrUnknownView {
			return err
		}

		descriptionView.Title = "Description"
		descriptionView.Editable = true
		descriptionView.Wrap = true

		fmt.Fprint(descriptionView, card.Desc)

		utils.ErrCheck(
			addEditorViewSwitching(gui, descriptionView, mngr),
			addEditorClosing(gui, descriptionView, mngr),
			addChangesSaving(gui, descriptionView, mngr, card),
		)
	}
	return
}

func commentsView(gui *Gui, mngr *TregoManager, card *trello.Card) (err error) {
	w, h := gui.Size()
	if commentsView, err := gui.SetView(cardCommentsView, int(w/3)+1, 9, w-1, h-5); err != nil {
		if err != ErrUnknownView {
			return err
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
	return
}

func addEditorViewSwitching(gui *Gui, view *View, mngr *TregoManager) error {
	return gui.SetKeybinding(view.Name(), KeyTab, ModNone, func(gui *Gui, view *View) error {
		idx := (view.Name()[17] - 48 + 1) % 4
		nextViewName := view.Name()[:17] + string(rune(idx+48))
		log.Printf("switching editor view to %v", nextViewName)
		utils.ErrCheck(mngr.SelectView(gui, nextViewName))
		mngr.CurrBotBarKey = nextViewName
		utils.ErrCheck(gui.DeleteView(BottomBar))
		return nil
	})
}

func addEditorClosing(gui *Gui, view *View, mngr *TregoManager) error {
	return gui.SetKeybinding(view.Name(), KeyCtrlQ, ModNone, func(gui *Gui, view *View) error {
		gui.DeleteKeybindings(cardNameView)
		gui.DeleteKeybindings(cardListInfoView)
		gui.DeleteKeybindings(cardLabelsView)
		gui.DeleteKeybindings(cardDescView)
		gui.DeleteKeybindings(cardCommentsView)
		utils.ErrCheck(
			gui.DeleteView(cardNameView),
			gui.DeleteView(cardListInfoView),
			gui.DeleteView(cardLabelsView),
			gui.DeleteView(cardDescView),
			gui.DeleteView(cardCommentsView),
		)
		mngr.Mode = BoardView
		mngr.currView = nil
		mngr.CurrBotBarKey = mngr.DefaultBotBarKey
		utils.ErrCheck(gui.DeleteView(BottomBar))
		return nil
	})
}

func addChangesSaving(gui *Gui, view *View, mngr *TregoManager, card *trello.Card) error {
	return gui.SetKeybinding(view.Name(), KeyCtrlS, ModNone, func(gui *Gui, view *View) (err error) {
		switch view.Name() {
		case cardNameView:
			_, err = card.SetName(view.Buffer()[:len(view.Buffer())-2])
		case cardDescView:
			_, err = card.SetDescription(view.Buffer())
		default:
			//Unsupported view
		}
		utils.ErrCheck(err)
		return nil
	})
}
