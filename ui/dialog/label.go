package dialog

import (
	"github.com/jroimartin/gocui"
	"github.com/aqatl/Trego/utils"
	"fmt"
	"sync"
	"strings"
)

type labelDialogMngr struct {
	labelChan            chan [2]string
	titleView, colorView *gocui.View
	gui                  *gocui.Gui

	once sync.Once
}

func (mngr *labelDialogMngr) Layout(gui *gocui.Gui) error {
	x1, y1, x2, y2 := calcDialogBounds(0, gui)
	y2 += 1

	titleView, err := setUpDialogView(gui, labelTitle, "Label title", x1, y1, x2, y2)
	utils.ErrCheck(err)
	mngr.titleView = titleView

	colorView, err := setUpDialogView(gui, labelColor, "Label color", x1, y1+3, x2, y2+10)
	utils.ErrCheck(err)
	mngr.colorView = colorView

	for i, col := range utils.FgColors {
		fmt.Fprintf(colorView, "%d.%v\n", i, col)
	}
	for i, col := range utils.HiFgColors {
		fmt.Fprintf(colorView, "%d.%v\n", i+len(utils.FgColors), col)
	}

	mngr.once.Do(mngr.configDialog)

	return nil
}

func (mngr *labelDialogMngr) configDialog() {
	switchViewFunc := func(gui *gocui.Gui, view *gocui.View) (err error) {
		if gui.CurrentView().Name() == labelTitle {
			_, err = gui.SetCurrentView(labelColor)
		} else {
			_, err = gui.SetCurrentView(labelTitle)
		}
		return
	}
	utils.ErrCheck(
		mngr.gui.SetKeybinding(labelTitle, gocui.KeyTab, gocui.ModNone, switchViewFunc),
		mngr.gui.SetKeybinding(labelColor, gocui.KeyTab, gocui.ModNone, switchViewFunc),
	)

	closeFunc := func(gui *gocui.Gui, v *gocui.View) error {
		err := closeLabelDialog(gui, mngr.titleView, mngr.colorView, mngr.labelChan)
		return err
	}
	utils.ErrCheck(
		mngr.gui.SetKeybinding(labelTitle, gocui.KeyEnter, gocui.ModNone, closeFunc),
		mngr.gui.SetKeybinding(labelColor, gocui.KeyEnter, gocui.ModNone, closeFunc),
	)

	utils.ErrCheck(
		mngr.gui.SetKeybinding(labelColor, gocui.KeyArrowUp, gocui.ModNone, utils.CursorUp),
		mngr.gui.SetKeybinding(labelColor, gocui.KeyArrowDown, gocui.ModNone, utils.CursorDown),
	)

	mngr.titleView.Editor = gocui.DefaultEditor

	_, err := mngr.gui.SetCurrentView(labelTitle)
	utils.ErrCheck(err)
}

func LabelDialog(gui *gocui.Gui, labelChan chan [2]string) gocui.Manager {
	mngr := &labelDialogMngr{
		labelChan: labelChan,
		gui:       gui,
	}

	return mngr
}

func closeLabelDialog(gui *gocui.Gui, titleView, colorView *gocui.View, labelChan chan [2] string) error {
	gui.DeleteKeybindings(labelTitle)
	gui.DeleteKeybindings(labelColor)

	title := strings.Trim(titleView.Buffer(), "\n ")
	colorIdx := utils.SelectedItemIdx(colorView)
	color := ""
	if colorIdx < 8 {
		color = utils.FgColors[colorIdx]
	} else {
		color = utils.HiFgColors[colorIdx-8]
	}

	utils.ErrCheck(
		gui.DeleteView(labelTitle),
		gui.DeleteView(labelColor),
	)

	labelChan <- [2]string{title, color}
	return nil
}
