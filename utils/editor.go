package utils

import (
	. "github.com/jroimartin/gocui"
	"strings"
	"time"
	"strconv"
)

//Moves cursor in the view one line up
func CursorUp(g *Gui, v *View) (err error) {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		if cy > 0 {
			err = v.SetCursor(cx, cy-1)
		}
		if oy > 0 && cy == 0 {
			err = v.SetOrigin(ox, oy-1)
		}

		viewBuffer := strings.Split(v.ViewBuffer(), "\n")
		if len(viewBuffer) > 0 && cy-1 >= 0 {
			if line := viewBuffer[cy-1]; len(line) > 0 && (line[0] < 0x30 || line[0] > 0x39) {
				return CursorUp(g, v)
			}
		}
	}
	return
}

//Moves cursor in the view one line down
func CursorDown(g *Gui, v *View) (err error) {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		viewBuffer := strings.Split(v.ViewBuffer(), "\n")
		if cy+oy < (len(viewBuffer) - 3) {
			if err = v.SetCursor(cx, cy+1); err != nil {
				err = v.SetOrigin(ox, oy+1)
			}
		}
		if len(viewBuffer) > 0 {
			line := viewBuffer[cy+1]
			if cy+1 != len(viewBuffer) && len(line) > 0 && (line[0] < 0x30 || line[0] > 0x39) {
				return CursorDown(g, v)
			}
		}
	}
	return
}

func AddNumericSelectEditor(gui *Gui, view *View) {
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
					CursorUp(gui, v)
				}
			} else {
				for i := 0; i < dst; i++ {
					CursorDown(gui, v)
				}
			}
		}
	})
}
