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

func (manager *TregoManager) SelectList(gui *gocui.Gui, listName string) (err error) {
	if view, err := gui.SetCurrentView(listName); err == nil {
		manager.currentView = view
		_, err = gui.SetViewOnTop(view.Name())
	}
	return
}