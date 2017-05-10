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
	"log"
	"time"
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
	utils.ErrCheck(
		bottomBarLayout(gui),
		topBarLayout(gui, mngr),
	)

	switch mngr.Mode {
	case BoardView:
		//loops through user's trello lists and adds them to gui
		for idx, list := range mngr.Lists {
			log.Printf("Adding view: %v", list.Name)
			utils.ErrCheck(AddList(gui, list, idx, mngr.listViewOffset))
		}
	case CardEditor:
		currView, err := gui.View(mngr.Lists[mngr.currListIdx].Id)
		utils.ErrCheck(err)
		CardEditorLayout(currView, gui, mngr)
		return nil
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

func addListEditor(gui *Gui, view *View) {
	dstNum := ""
	lastKey := time.Now()

	view.Editor = EditorFunc(func(v *View, key Key, ch rune, mod Modifier) {
		if ch >= 0x30 && ch <= 0x39 {
			if time.Since(lastKey).Seconds() > 1 {
				dstNum = ""
			}
			dstNum = dstNum + string(ch)

			lastKey = time.Now()
		} else if ch == 'g' {
			dst, err := strconv.Atoi(dstNum)
			if err != nil {
				return
			}
			viewLines := strings.Split(v.ViewBuffer(), "\n")
			if len(viewLines) <= dst {
				return
			}

			_, cy := v.Cursor()
			dst -= cy
			if dst < 0 {
				for i := 0; i > dst; i-- {
					utils.CursorUp(gui, v)
				}
			} else {
				for i := 0; i < dst; i++ {
					utils.CursorDown(gui, v)
				}
			}
		}
	})
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

		addListEditor(gui, v)

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

func topBarLayout(gui *Gui, mngr *TregoManager) error {
	maxX, _ := gui.Size()
	if v, err := gui.SetView(TopBar, 0, 0, maxX-1, 2); err != nil {
		if err != ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Highlight = false
		v.BgColor = ColorBlack

		color.Output = v
		color.New(color.FgYellow).Add(color.Bold).Printf("Board: %v", mngr.CurrBoard.Name)
	}

	return nil
}

//bottom bar with shortcuts
func bottomBarLayout(gui *Gui) error {
	maxX, maxY := gui.Size()
	if v, err := gui.SetView(BottomBar, 0, maxY-4, maxX-1, maxY-1); err != nil {
		if err != ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Highlight = false
		v.BgColor = ColorBlack

		fmt.Fprint(v,
			cyan("\xE2\x87\x84"), yell(" change list "),
			cyan("^C"), yell(" exit "),
			cyan("n"), yell(" add card "),
			cyan("l"), yell(" add list "),
			cyan("d"), yell(" delete "),
			cyan("s"), yell(" card search "),
			cyan("m"), yell(" move card\n"),
			cyan("\xE2\x87\x85"), yell(" change card "),
			cyan("^d"), yell(" move card down "),
			cyan("^u"), yell(" move card up "),
			cyan("b"), yell(" change board"))
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
