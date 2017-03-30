package utils

import (
	"strings"
	. "github.com/jroimartin/gocui"
)

//Moves cursor in list one line up
func CursorUp(g *Gui, v *View) (err error) {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		if cy > 0 {
			err = v.SetCursor(cx, cy - 1)
		}
		if oy > 0 && cy == 0 {
			err = v.SetOrigin(ox, oy - 1)
		}
	}
	return
}

//Moves cursor in list one line down
func CursorDown(g *Gui, v *View) (err error) {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		if cy + oy < (len(strings.Split(v.ViewBuffer(), "\n")) - 3) {
			if err = v.SetCursor(cx, cy + 1); err != nil {
				err = v.SetOrigin(ox, oy + 1)
			}
		}
	}
	return
}

