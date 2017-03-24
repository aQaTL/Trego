package ui

import (
	"github.com/VojtechVitek/go-trello"
	"github.com/jroimartin/gocui"
	"github.com/aqatl/Trego/utils"
)

type TregoManager struct {
	Member      *trello.Member
	Lists       []trello.List
	currListIdx int
	currView    *gocui.View
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
			utils.ErrCheck(mngr.SelectView(gui, mngr.Lists[mngr.currListIdx].Name))
		} else {
			utils.ErrCheck(mngr.SelectView(gui, replacementViewName))
		}
	}
}