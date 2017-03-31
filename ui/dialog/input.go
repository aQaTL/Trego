package dialog

import (
	"github.com/jroimartin/gocui"
	"github.com/aqatl/Trego/utils"
	"fmt"
	"strings"
)

//Provides dialog with an input field
func InputDialog(msg, title, initValue string, gui *gocui.Gui, input chan string) *gocui.View {
	msgL := len(msg)
	x1, y1, x2, y2 := calcDialogBounds(msgL, gui)

	dialogView, err := setUpDialogView(gui, INPUT_DIALOG, title, x1, y1, x2, y2)
	utils.ErrCheck(err)
	dialogView.Highlight = false
	printCentered(dialogView, msg, (x2 - x1))

	inputView, err := setUpDialogView(gui, INPUT_FIELD, "", x1, y1 + 3, x2, y2 + 3) //Place it a little lower
	utils.ErrCheck(err)
	inputView.Editable = true
	fmt.Fprint(inputView, initValue)

	gui.SetKeybinding(
		INPUT_FIELD,
		gocui.KeyEnter,
		gocui.ModNone,
		func(gui *gocui.Gui, view *gocui.View) error {
			dialogCleanUp(gui, INPUT_DIALOG, INPUT_FIELD)
			input <- strings.TrimSuffix(inputView.Buffer(), " \n")
			close(input)
			return nil
		})

	return inputView
}
