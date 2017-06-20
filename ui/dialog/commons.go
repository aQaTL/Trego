/*
	Package that contains variety of dialogs

	Example of using InputDialog

		input := make(chan string)
		//Prevents nested dialogs and other glitches (like double handler call)
		for _, view := range gui.Views() {
			gui.DeleteKeybindings(view.Name())
		}
		utils.ErrCheck(
			mngr.SelectView(
				gui,
				dialog.InputDialog("Are you sure? [y/n]", "", "", gui, input).Name()))
		go func() {
			if userInput, ok := <-input; ok {
				log.Print(userInput)
			}
			SetKeyBindings(gui, mngr)
		}()

	Similar example for ConfirmDialog

		option := make(chan bool)
		previousView := gui.CurrentView()
		log.Println("callback: ", previousView.Name())
		//Prevents nested dialogs and other glitches (like double handler call)
		for _, view := range gui.Views() {
			gui.DeleteKeybindings(view.Name())
		}
		utils.ErrCheck(
			mngr.SelectView(
				gui,
				dialog.ConfirmDialog("message", "title", gui, option).Name()))

		go func() {
			if choice, ok = <-option; ok {
				log.Printf("Choosen option: %v", choice)
			}

			mngr.currView = previousView
			log.Println("thread: ", mngr.currView.Name(), previousView.Name())
			SetKeyBindings(gui, mngr)
		}()

	Remember! When user chooses an option, it deletes returned view, so make sure
	to handle it properly. If you don't, it may cause Unknown view error
*/
package dialog

import (
	"fmt"
	"github.com/aqatl/Trego/utils"
	"github.com/jroimartin/gocui"
	"io"
	"math"
	"strings"
)

const (
	confirmDialog = "confirmdialogview"
	inputDialog   = "inputdialogview"
	inputField    = "inputdialogfield"
	selectDialog  = "selectdialogview"
	labelTitle    = "labeltitledialogview"
	labelColor    = "labelcolordialogview"
)

func calcDialogBounds(msgL int, gui *gocui.Gui) (x1, y1, x2, y2 int) {
	wi, hi := gui.Size()
	w, h := float64(wi), float64(hi)
	x1 = int(w / 2.0 * 0.7)
	x2 = int(w / 2 * 1.3)
	y1 = int(h / 2 * 0.8)
	y2 = int(float64(y1)+math.Ceil(float64(msgL)/float64(x2-x1))) + 1

	return
}

func cleanUp(gui *gocui.Gui, dialogTypes ...string) {
	for _, dialogType := range dialogTypes {
		gui.DeleteKeybindings(dialogType)
		utils.ErrCheck(gui.DeleteView(dialogType))
	}
}

func printCentered(w io.Writer, text string, viewWidth int) {
	msgL := len(text)
	if msgL < viewWidth {
		fmt.Fprintf(w, "%s%s\n", strings.Repeat(" ", (viewWidth/2-(msgL/2))-1), text)
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
		view.Editable = true
		view.FgColor = gocui.ColorBlack
		view.FgColor = gocui.ColorGreen
		view.Title = viewTitle

		utils.AddNumericSelectEditor(gui, view)
	}
	return view, nil
}
