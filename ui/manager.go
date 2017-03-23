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

func (manager *TregoManager) SelectView(gui *gocui.Gui, viewName string) error {
	if view, err := gui.SetCurrentView(viewName); err == nil {
		manager.currView = view
		_, err = gui.SetViewOnTop(view.Name())
		return err
	} else {
		return err
	}
}

func (manager *TregoManager) CheckCurrView(gui *gocui.Gui, replacementViewName string) {
	if manager.currView == nil {
		if len(manager.Lists) > 0 {
			utils.ErrCheck(manager.SelectView(gui, manager.Lists[0].Name))
		} else {
			utils.ErrCheck(manager.SelectView(gui, replacementViewName))
		}
	}
}