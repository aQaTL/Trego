package dialog

import (
	"github.com/aqatl/Trego/utils"
	"github.com/jroimartin/gocui"
)

//Provides option dialog. User can make decision with 'y' and 'n' key.
func ConfirmDialog(msg, title string, gui *gocui.Gui, choice chan bool) (view *gocui.View) {
	msgL := len(msg)
	x1, y1, x2, y2 := calcDialogBounds(msgL, gui)
	winW := x2 - x1

	confirmView, err := setUpDialogView(gui, confirmDialog, title, x1, y1, x2, y2)
	utils.ErrCheck(err)
	confirmView.Highlight = false
	printCentered(confirmView, msg, winW)
	view = confirmView

	utils.ErrCheck(gui.SetKeybinding(confirmDialog, 'y', gocui.ModNone,
		func(gui *gocui.Gui, view *gocui.View) (err error) {
			cleanUp(gui, confirmDialog)
			choice <- true
			close(choice)
			return
		}))

	utils.ErrCheck(gui.SetKeybinding(confirmDialog, 'n', gocui.ModNone,
		func(gui *gocui.Gui, view *gocui.View) (err error) {
			cleanUp(gui, confirmDialog)
			choice <- false
			close(choice)
			return
		}))

	utils.ErrCheck(
		gui.SetKeybinding(
			confirmDialog,
			gocui.KeyCtrlQ,
			gocui.ModNone,
			func(gui *gocui.Gui, view *gocui.View) error {
				close(choice)
				cleanUp(gui, confirmDialog)
				return nil
			}))

	return
}
