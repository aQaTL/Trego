package dialog

import (
	"fmt"
	"github.com/aqatl/Trego/utils"
	"github.com/jroimartin/gocui"
	"strings"
)

//Provides dialog with an input field
func InputDialog(msg, title, initValue string, gui *gocui.Gui, input chan string) *gocui.View {
	msgL := len(msg)
	x1, y1, x2, y2 := calcDialogBounds(msgL, gui)

	dialogView, err := setUpDialogView(gui, inputDialog, title, x1, y1, x2, y2)
	utils.ErrCheck(err)
	dialogView.Highlight = false
	printCentered(dialogView, msg, (x2 - x1))

	inputView, err := setUpDialogView(gui, inputField, "", x1, y1+3, x2, y2+3) //Place it a little lower
	utils.ErrCheck(err)
	inputView.Editable = true
	fmt.Fprint(inputView, initValue)

	utils.ErrCheck(
		gui.SetKeybinding(
			inputField,
			gocui.KeyEnter,
			gocui.ModNone,
			func(gui *gocui.Gui, view *gocui.View) error {
				dialogCleanUp(gui, inputDialog, inputField)
				input <- strings.TrimSuffix(inputView.Buffer(), " \n")
				close(input)
				return nil
			}))

	utils.ErrCheck(
		gui.SetKeybinding(
			inputField,
			gocui.KeyCtrlQ,
			gocui.ModNone,
			func(gui *gocui.Gui, view *gocui.View) error {
				close(input)
				dialogCleanUp(gui, inputDialog, inputField)
				return nil
			}))

	return inputView
}
