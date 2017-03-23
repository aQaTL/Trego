package ui

import (
	. "github.com/jroimartin/gocui"
	"github.com/fatih/color"
	"math"
	"github.com/VojtechVitek/go-trello"
	"github.com/aqatl/Trego/utils"
)

const (
	BOTTOM_BAR string = "botbar"
	TOP_BAR string = "topbar"
	BOARD_SELECT string = "boardselectview"
	LIST_WIDTH int = 24
)

func (manager *TregoManager) Layout(gui *Gui) error {
	if err := bottomBarLayout(gui); err != nil {
		return err
	}
	if err := topBarLayout(gui); err != nil {
		return err
	}

	//loops through user's trello lists and adds them to gui
	for idx, list := range (manager.Lists) {
		utils.ErrCheck(AddList(gui, list, idx))
	}

	manager.CheckCurrView(gui, TOP_BAR)

	if _, err := gui.SetCurrentView(manager.currView.Name()); err != nil {
		manager.currView = nil
		manager.CheckCurrView(gui, TOP_BAR)
		utils.ErrCheck(manager.SelectView(gui, manager.currView.Name()))
		return nil
	}

	return nil
}

func AddList(gui *Gui, list trello.List, index int) error {
	_, maxY := gui.Size()
	if v, err := gui.SetView(list.Name,
		index * LIST_WIDTH + int(math.Abs(sign(index))), 3,
		index * LIST_WIDTH + LIST_WIDTH, maxY - 5);
			err != nil {

		if err != ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Highlight = true
		v.Wrap = true
		v.BgColor = ColorBlack
		v.SelBgColor = ColorGreen
		v.SelFgColor = ColorBlack
		v.FgColor = ColorWhite
		v.Title = list.Name
		gui.Cursor = true

		color.Output = v

		cards, err := list.Cards()
		if err != nil {
			return err
		}
		for idx, card := range (cards) {
			color.New(color.BgBlack).Add(color.FgWhite).Printf("%d.%v\n", idx, card.Name)
		}
	}

	return nil
}

func topBarLayout(gui *Gui) error {
	maxX, _ := gui.Size()
	if v, err := gui.SetView(TOP_BAR, 0, 0, maxX - 1, 2); err != nil {
		if err != ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Highlight = false
		v.BgColor = ColorBlack

		color.Output = v
		color.New(color.FgYellow).Add(color.Bold).Printf("Board: %v", "Trego")
	}

	return nil
}

//bottom bar with shortcuts
func bottomBarLayout(gui *Gui) error {
	maxX, maxY := gui.Size()
	if v, err := gui.SetView(BOTTOM_BAR, 0, maxY - 4, maxX - 1, maxY - 1); err != nil {
		if err != ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Highlight = false
		v.BgColor = ColorBlack

		color.Output = v
		color := color.New(color.FgYellow).Add(color.Bold)
		color.Printf("%-22s", "\xE2\x87\x84 move between lists")
		color.Printf("%s\n", "^C exit")
		color.Printf("%-22s", "\xE2\x87\x85 move inside list")
		color.Printf("%s\n", "b change board")

	}
	return nil
}

func sign(x int) float64 {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	} else {
		return 0
	}
}