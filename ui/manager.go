package ui

import (
	"github.com/aqatl/Trego/utils"
	"github.com/aqatl/go-trello"
	"github.com/jroimartin/gocui"
)

type TregoManager struct {
	Member         *trello.Member
	Lists          []trello.List
	CurrBoard      *trello.Board
	currListIdx    int
	currView       *gocui.View
	listViewOffset int
}

func (mngr *TregoManager) SelectView(gui *gocui.Gui, viewName string) error {
	if view, err := gui.SetCurrentView(viewName); err == nil {
		mngr.currView = view
		_, err = gui.SetViewOnTop(view.Name())
		return err
	} else {
		return err
	}
}

func (mngr *TregoManager) CheckCurrView(gui *gocui.Gui, replacementViewName string) {
	if mngr.currView == nil {
		if len(mngr.Lists) > 0 {
			utils.ErrCheck(mngr.SelectView(gui, mngr.Lists[mngr.currListIdx].Id))
		} else {
			utils.ErrCheck(mngr.SelectView(gui, replacementViewName))
		}
	}
}
