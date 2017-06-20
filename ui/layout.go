package ui

import (
	"github.com/aqatl/Trego/utils"
	"github.com/aqatl/go-trello"
	"github.com/fatih/color"
	. "github.com/jroimartin/gocui"
	"math"
	"strconv"
	"strings"
	"fmt"
)

const (
	BottomBar  string = "botbar"
	TopBar     string = "topbar"
	SearchView string = "searchview"
	ListWidth  int    = 24
)

var (
	defaultCardColor = color.New(color.BgBlack).Add(color.FgWhite)
	yell             = color.New(color.FgYellow).Add(color.Bold).SprintFunc()
	cyan             = color.New(color.FgCyan).Add(color.Bold).SprintFunc()
)

func (mngr *TregoManager) Layout(gui *Gui) error {
	//loops through user's trello lists and adds them to gui
	for idx, list := range mngr.Lists {
		utils.ErrCheck(AddList(gui, list, idx, mngr.listViewOffset))
	}

	mngr.CheckCurrView(gui, TopBar)

	if _, err := gui.SetCurrentView(mngr.currView.Name()); err != nil {
		mngr.currView = nil
		mngr.CheckCurrView(gui, TopBar)
		utils.ErrCheck(mngr.SelectView(gui, mngr.currView.Name()))
		return nil
	}

	return nil
}

func AddList(gui *Gui, list trello.List, index, offset int) error {
	_, maxY := gui.Size()
	if v, err := gui.SetView(list.Id,
		(index+offset)*ListWidth+int(math.Abs(sign(index))), 3,
		(index+offset)*ListWidth+ListWidth, maxY-5); err != nil {

		if err != ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Highlight = true
		v.Wrap = true
		v.BgColor = ColorBlack
		v.SelBgColor = ColorGreen
		v.SelFgColor = ColorBlack
		v.FgColor = ColorWhite
		v.Title = list.Name
		gui.Cursor = true

		utils.AddNumericSelectEditor(gui, v)

		color.Output = v

		cards, err := list.FreshCards()
		if err != nil {
			return err
		}
		v.Clear()
		for idx, card := range cards {
			defaultCardColor.Printf("%d.%v\n", idx, card.Name)
		}
	}

	return nil
}

type InfoBar struct{
	BoardName string
}

func (iBar *InfoBar) Layout(gui *Gui) error {
	maxX, _ := gui.Size()
	if v, err := gui.SetView(TopBar, 0, 0, maxX-1, 2); err != nil {
		if err != ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Highlight = false
		v.BgColor = ColorBlack

		color.Output = v
		color.New(color.FgYellow).Add(color.Bold).Printf("Board: %v", iBar.BoardName)
	}

	return nil
}

type ShortcutsBar struct {
	DefaultBotBarKey string `json:"default_bottom_bar_key"`
	Data             map[string][][]string `json:"bottom_bar"`

	CurrBotBarKey string
}

//bottom bar with shortcuts
func (botBar *ShortcutsBar) Layout(gui *Gui) error {
	maxX, maxY := gui.Size()
	if v, err := gui.SetView(BottomBar, 0, maxY-4, maxX-1, maxY-1); err != nil {
		if err != ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Highlight = false
		v.BgColor = ColorBlack

		if botBar.CurrBotBarKey == "" {
			botBar.CurrBotBarKey = botBar.DefaultBotBarKey
		}

		shortcuts := botBar.Data[botBar.CurrBotBarKey]
		for i := 0; i < 2; i++ {
			for j := 0; j < len(shortcuts[i]); j += 2 {
				fmt.Fprintf(v,
					"%v %v ",
					cyan(shortcuts[i][j]),
					yell(shortcuts[i][j+1]),
				)
			}
			fmt.Fprintln(v)
		}
	}
	return nil
}

func SelectedItemIdx(view *View) int {
	if len(strings.Split(view.Buffer(), "\n")) <= 1 {
		return -1
	}

	_, cy := view.Cursor()
	currLine, err := view.Line(cy)
	dotIdx := strings.Index(currLine, ".")
	itemIdx64, err := strconv.ParseInt(currLine[:dotIdx], 10, 32)
	utils.ErrCheck(err)
	return int(itemIdx64)
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
