package dialog

import (
	"github.com/jroimartin/gocui"
	"log"
	"fmt"
	"strings"
	"math"
)

const (
	DIALOG_VIEW = "dialogview"
)

//Remeber! When user chooses an option, it deletes returned view, so make sure
//to handle it properly. If you don't, it may cause Unknown view error
func ConfirmDialog(msg, title string, gui *gocui.Gui, choice chan bool) (view *gocui.View) {
	wi, hi := gui.Size()
	w, h := float64(wi), float64(hi)
	x1 := int(w / 2.0 * 0.7)
	x2 := int(w / 2 * 1.3)
	y1 := int(h / 2 * 0.8)

	msgL := len(msg)
	winW := x2 - x1
	y2 := int(float64(y1) + math.Ceil(float64(msgL) / float64(winW))) + 1

	if v, err := gui.SetView(DIALOG_VIEW, x1, y1, x2, y2); err != nil {
		if err != gocui.ErrUnknownView {
			log.Panicln(err)
		}

		v.Wrap = true
		v.FgColor = gocui.ColorBlack
		v.BgColor = gocui.ColorGreen
		v.Highlight = false
		if title != "" {
			v.Title = title
		}

		if msgL < winW {
			fmt.Fprintf(v, "%s%s\n", strings.Repeat(" ", (winW / 2 - (msgL / 2)) - 1), msg)
		} else {
			fmt.Fprintf(v, "%s\n", msg)
		}
		view = v
	}

	if err := gui.SetKeybinding(DIALOG_VIEW, 'y', gocui.ModNone,
		func(gui *gocui.Gui, view *gocui.View) (err error) {
			gui.DeleteKeybindings(DIALOG_VIEW)
			err = gui.DeleteView(DIALOG_VIEW)
			choice <- true
			return
		}); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding(DIALOG_VIEW, 'n', gocui.ModNone,
		func(gui *gocui.Gui, view *gocui.View) (err error) {
			gui.DeleteKeybindings(DIALOG_VIEW)
			err = gui.DeleteView(DIALOG_VIEW)
			choice <- false
			return
		}); err != nil {
		log.Panicln(err)
	}
	return
}
