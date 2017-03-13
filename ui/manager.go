package ui

import (
	"github.com/VojtechVitek/go-trello"
	"github.com/jroimartin/gocui"
)

type TregoManager struct {
	Member      *trello.Member
	Lists       []trello.List
	currentView *gocui.View
}

func (manager *TregoManager) SelectView(gui *gocui.Gui, viewName string) (err error) {
	if view, err := gui.SetCurrentView(viewName); err == nil {
		manager.currentView = view
		_, err = gui.SetViewOnTop(view.Name())
	}
	return
}