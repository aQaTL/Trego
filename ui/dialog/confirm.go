package dialog

import (
	"github.com/jroimartin/gocui"
	"log"
	"fmt"
	"strings"
	"math"
	"io"
	"github.com/aqatl/Trego/utils"
)

const (
	CONFIRM_DIALOG = "confirmdialogview"
	INPUT_DIALOG = "inputdialogview"
	INPUT_FIELD = "inputdialogfield"
)

func calcDialogBounds(msgL int, gui *gocui.Gui) (x1, y1, x2, y2 int) {
	wi, hi := gui.Size()
	w, h := float64(wi), float64(hi)
	x1 = int(w / 2.0 * 0.7)
	x2 = int(w / 2 * 1.3)
	y1 = int(h / 2 * 0.8)
	y2 = int(float64(y1) + math.Ceil(float64(msgL) / float64(x2 - x1))) + 1

	return
}

//Remember! When user chooses an option, it deletes returned view, so make sure
//to handle it properly. If you don't, it may cause Unknown view error
func ConfirmDialog(msg, title string, gui *gocui.Gui, choice chan bool) (view *gocui.View) {
	msgL := len(msg)
	x1, y1, x2, y2 := calcDialogBounds(msgL, gui)
	winW := x2 - x1
	currView := gui.CurrentView()

	confirmView, err := setUpDialogView(gui, CONFIRM_DIALOG, title, x1, y1, x2, y2)
	utils.ErrCheck(err)
	confirmView.Highlight = false
	printCentered(confirmView, msg, winW)
	view = confirmView

	if err := gui.SetKeybinding(CONFIRM_DIALOG, 'y', gocui.ModNone,
		func(gui *gocui.Gui, view *gocui.View) (err error) {
			dialogCleanUp(gui, currView, CONFIRM_DIALOG)
			choice <- true
			close(choice)
			return
		}); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding(CONFIRM_DIALOG, 'n', gocui.ModNone,
		func(gui *gocui.Gui, view *gocui.View) (err error) {
			dialogCleanUp(gui, currView, CONFIRM_DIALOG)
			choice <- false
			close(choice)
			return
		}); err != nil {
		log.Panicln(err)
	}
	return
}

func dialogCleanUp(gui *gocui.Gui, previousView *gocui.View, dialogTypes ...string) {
	if _, err := gui.SetCurrentView(previousView.Name()); err != nil {
		log.Panicln(err)
	}
	for _, dialogType := range dialogTypes {
		gui.DeleteKeybindings(dialogType)
		utils.ErrCheck(gui.DeleteView(dialogType))
	}
}

func InputDialog(msg, title, initValue string, gui *gocui.Gui, input chan string) *gocui.View {
	msgL := len(msg)
	x1, y1, x2, y2 := calcDialogBounds(msgL, gui)
	currView := gui.CurrentView()

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
			dialogCleanUp(gui, currView, INPUT_DIALOG, INPUT_FIELD)
			input <- inputView.Buffer()
			close(input)
			return err
	})

	return inputView
}

func printCentered(w io.Writer, text string, viewWidth int) {
	msgL := len(text)
	if msgL < viewWidth {
			fmt.Fprintf(w, "%s%s\n", strings.Repeat(" ", (viewWidth / 2 - (msgL / 2)) - 1), text)
		} else {
			fmt.Fprintf(w, "%s\n", text)
		}
}

func setUpDialogView(gui *gocui.Gui, viewName, viewTitle string, x1, y1, x2, y2 int) (view *gocui.View, err error) {
	if view, err = gui.SetView(viewName, x1, y1, x2, y2); err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

		view.Highlight = true
		view.Wrap = true
		view.Editable = false
		view.FgColor = gocui.ColorBlack
		view.FgColor = gocui.ColorGreen
		view.Title = viewTitle
	}
	return view, nil
}