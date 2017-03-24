package dialog

import (
	"github.com/jroimartin/gocui"
	"github.com/aqatl/Trego/utils"
	"log"
)

//Provides option dialog. User can make decision with 'y' and 'n' key.
func ConfirmDialog(msg, title string, gui *gocui.Gui, choice chan bool) (view *gocui.View) {
	msgL := len(msg)
	x1, y1, x2, y2 := calcDialogBounds(msgL, gui)
	winW := x2 - x1

	confirmView, err := setUpDialogView(gui, CONFIRM_DIALOG, title, x1, y1, x2, y2)
	utils.ErrCheck(err)
	confirmView.Highlight = false
	printCentered(confirmView, msg, winW)
	view = confirmView

	if err := gui.SetKeybinding(CONFIRM_DIALOG, 'y', gocui.ModNone,
		func(gui *gocui.Gui, view *gocui.View) (err error) {
			dialogCleanUp(gui, CONFIRM_DIALOG)
			choice <- true
			close(choice)
			return
		}); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding(CONFIRM_DIALOG, 'n', gocui.ModNone,
		func(gui *gocui.Gui, view *gocui.View) (err error) {
			dialogCleanUp(gui, CONFIRM_DIALOG)
			choice <- false
			close(choice)
			return
		}); err != nil {
		log.Panicln(err)
	}
	return
}
