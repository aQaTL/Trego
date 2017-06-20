package main

import (
	"encoding/json"
	"github.com/aqatl/Trego/conn"
	"github.com/aqatl/Trego/ui"
	"github.com/aqatl/Trego/utils"
	. "github.com/jroimartin/gocui"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	gui, err := NewGui(OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	logF, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	utils.ErrCheck(err)
	log.SetOutput(logF)

	user := conn.Connect(gui)
	board := conn.BoardByName(user, "Testing board")
	lists := conn.Lists(board)

	var shortcutsBar ui.ShortcutsBar
	stringsJson, err := ioutil.ReadFile("strings.json")
	utils.ErrCheck(err)
	err = json.Unmarshal(stringsJson, &shortcutsBar)
	utils.ErrCheck(err)

	infoBar := &ui.InfoBar{BoardName: board.Name}

	gui.Mouse = false
	gui.Highlight = true
	gui.SelFgColor = ColorGreen

	mngr := &ui.TregoManager{
		Member:    user,
		Lists:     lists,
		CurrBoard: board,
		BotBar:    &shortcutsBar,
		TopBar:    infoBar,
	}

	gui.SetManager(mngr, &shortcutsBar, infoBar)

	ui.SetKeyBindings(gui, mngr)

	if err := gui.MainLoop(); err != nil && err != ErrQuit {
		log.Panicln(err)
	}
}
