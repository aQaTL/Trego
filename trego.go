package main

import (
	"github.com/aqatl/Trego/conn"
	"github.com/aqatl/Trego/ui"
	. "github.com/jroimartin/gocui"
	"log"
)

func main() {
	gui, err := NewGui(OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	user := conn.Connect(gui)
	board := conn.BoardByName(user, "Testing board")
	lists := conn.Lists(board)

	gui.Mouse = false
	gui.Highlight = true
	gui.SelFgColor = ColorGreen
	mngr := &ui.TregoManager{Member: user, Lists: lists, CurrBoard: board}
	gui.SetManager(mngr)

	ui.SetKeyBindings(gui, mngr)

	if err := gui.MainLoop(); err != nil && err != ErrQuit {
		log.Panicln(err)
	}
}
