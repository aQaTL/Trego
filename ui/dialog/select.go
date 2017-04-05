package dialog

import (
	"github.com/jroimartin/gocui"
	"github.com/aqatl/Trego/utils"
	"fmt"
)

func SelectDialog(title string, gui *gocui.Gui, selIdxC chan int, values []string) *gocui.View {
	x1, y1, x2, y2 := calcDialogBounds(0, gui)
	y2 += len(values)

	dialogView, err := setUpDialogView(gui, SELECT_DIALOG, title, x1, y1, x2, y2)
	utils.ErrCheck(err)

	for idx, item := range values {
		fmt.Fprintf(dialogView, "%d.%v\n", idx, item)
	}

	utils.ErrCheck(
		gui.SetKeybinding(SELECT_DIALOG, gocui.KeyArrowDown, gocui.ModNone, utils.CursorDown),
		gui.SetKeybinding(SELECT_DIALOG, gocui.KeyArrowUp, gocui.ModNone, utils.CursorUp))

	utils.ErrCheck(
		gui.SetKeybinding(
			SELECT_DIALOG,
			gocui.KeyEnter,
			gocui.ModNone,
			func(gui *gocui.Gui, view *gocui.View) error {
				dialogCleanUp(gui, SELECT_DIALOG)
				_, cy := view.Cursor()
				selIdxC <- cy
				close(selIdxC)
				return nil
			}))

	utils.ErrCheck(
		gui.SetKeybinding(
			SELECT_DIALOG,
			gocui.KeyCtrlQ,
			gocui.ModNone,
			func(gui *gocui.Gui, view *gocui.View) error {
				close(selIdxC)
				dialogCleanUp(gui, SELECT_DIALOG)
				return nil
			}))

	return dialogView
}
