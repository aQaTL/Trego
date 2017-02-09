package main

import (
	."github.com/jroimartin/gocui"
	"log"
	"github.com/aqatl/Trego/trego/ui"
)


func main() {

	gui, err := NewGui(OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	gui.Mouse = false
	gui.SetManagerFunc(ui.Layout)

	ui.SetKeyBindings(gui)

	if err := gui.MainLoop(); err != nil && err != ErrQuit {
		log.Panicln(err)
	}
}
