package main

import (
	. "github.com/jroimartin/gocui"
	"log"
	"github.com/aqatl/Trego/trego/ui"
	"github.com/aqatl/Trego/trego/conn"
)

func main() {
	user := conn.Connect()
	lists := conn.GetLists(conn.GetBoard(user, "Trego"))

	gui, err := NewGui(OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	gui.Mouse = false
	manager := &ui.TregoManager{Lists: lists}
	gui.SetManager(manager)

	ui.SetKeyBindings(gui, manager)

	if err := gui.MainLoop(); err != nil && err != ErrQuit {
		log.Panicln(err)
	}
}
