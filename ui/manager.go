package ui

import (
	"github.com/aqatl/Trego/utils"
	"github.com/aqatl/go-trello"
	"github.com/jroimartin/gocui"
)

type TregoManager struct {
	Mode TregoMode //Initially BoardView

	Member         *trello.Member
	Lists          []trello.List
	CurrBoard      *trello.Board
	currListIdx    int
	currView       *gocui.View
	listViewOffset int
}

type TregoMode int

const (
	BoardView  TregoMode = iota
	CardEditor
)

func (mngr *TregoManager) SelectView(gui *gocui.Gui, viewName string) error {
	view, err := gui.SetCurrentView(viewName)
	if err != nil {
		return err
	}
	mngr.currView = view
	_, err = gui.SetViewOnTop(view.Name())
	return err
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

func (mngr *TregoManager) SwitchListRight(gui *gocui.Gui, v *gocui.View) (err error) {
	mngr.currListIdx = (mngr.currListIdx + 1) % len(mngr.Lists)
	nextViewId := mngr.Lists[mngr.currListIdx].Id

	_, _, x2, _, err := gui.ViewPosition(nextViewId)
	w, _ := gui.Size()
	if x2 > w {
		mngr.listViewOffset -= 1
	} else if mngr.currListIdx == 0 {
		mngr.listViewOffset = 0
	}

	err = mngr.SelectView(gui, nextViewId)
	return
}

func (mngr *TregoManager) SwitchListLeft(gui *gocui.Gui, v *gocui.View) (err error) {
	if mngr.currListIdx == 0 {
		mngr.currListIdx = len(mngr.Lists)
	}
	mngr.currListIdx--
	previousViewId := mngr.Lists[mngr.currListIdx%len(mngr.Lists)].Id

	x1, _, _, _, err := gui.ViewPosition(previousViewId)
	if x1 < 0 {
		mngr.listViewOffset += 1
	} else if mngr.currListIdx == len(mngr.Lists)-1 {
		//Scrolls board to the end
		for mngr.currListIdx = 0; mngr.currListIdx != len(mngr.Lists)-1; {
			utils.ErrCheck(mngr.SwitchListRight(gui, mngr.currView))
		}
	}

	return mngr.SelectView(gui, previousViewId)
}
