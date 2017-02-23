package ui

import (

)
import (
	"github.com/VojtechVitek/go-trello"
	"github.com/jroimartin/gocui"
	"log"
)

type TregoManager struct {
	Lists []trello.List
	currentView *gocui.View
}

func (manager *TregoManager) selectList(gui *gocui.Gui, listName string) {
	if view, err := gui.SetCurrentView(listName); err != nil {
		log.Panicln(err)
	} else {
		manager.currentView = view
	}
}

