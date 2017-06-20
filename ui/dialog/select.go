package dialog

import (
	"fmt"
	"github.com/aqatl/Trego/utils"
	"github.com/jroimartin/gocui"
)

func SelectDialog(title string, gui *gocui.Gui, selIdxC chan int, values []string) *gocui.View {
	x1, y1, x2, y2 := calcDialogBounds(0, gui)
	y2 += len(values)

	dialogView, err := setUpDialogView(gui, selectDialog, title, x1, y1, x2, y2)
	utils.ErrCheck(err)

	for idx, item := range values {
		fmt.Fprintf(dialogView, "%d.%v\n", idx, item)
	}

	utils.ErrCheck(
		gui.SetKeybinding(selectDialog, gocui.KeyArrowDown, gocui.ModNone, utils.CursorDown),
		gui.SetKeybinding(selectDialog, gocui.KeyArrowUp, gocui.ModNone, utils.CursorUp))

	utils.ErrCheck(
		gui.SetKeybinding(
			selectDialog,
			gocui.KeyEnter,
			gocui.ModNone,
			func(gui *gocui.Gui, view *gocui.View) error {
				cleanUp(gui, selectDialog)
				_, cy := view.Cursor()
				selIdxC <- cy
				close(selIdxC)
				return nil
			}))

	utils.ErrCheck(
		gui.SetKeybinding(
			selectDialog,
			gocui.KeyCtrlQ,
			gocui.ModNone,
			func(gui *gocui.Gui, view *gocui.View) error {
				close(selIdxC)
				cleanUp(gui, selectDialog)
				return nil
			}))

	return dialogView
}
