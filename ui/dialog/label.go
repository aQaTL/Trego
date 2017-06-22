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

	i := 0
	for _, col := range utils.FgColors {
		if col != "-" {
			fmt.Fprintf(mngr.colorView, "%d.%v\n", i, col)
			i++
		}
	}
	for _, col := range utils.HiFgColors {
		if col != "-" {
			fmt.Fprintf(mngr.colorView, "%d.%v\n", i, col)
			i++
		}
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

	utils.ErrCheck(
		mngr.gui.SetKeybinding(labelTitle, gocui.KeyEnter, gocui.ModNone, mngr.commitLabelDialog),
		mngr.gui.SetKeybinding(labelColor, gocui.KeyEnter, gocui.ModNone, mngr.commitLabelDialog),
		mngr.gui.SetKeybinding(labelTitle, gocui.KeyCtrlQ, gocui.ModNone, mngr.closeLabelDialog),
		mngr.gui.SetKeybinding(labelColor, gocui.KeyCtrlQ, gocui.ModNone, mngr.closeLabelDialog),
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

func (mngr *labelDialogMngr) closeLabelDialog(gui *gocui.Gui, view *gocui.View) error {
	DeleteDialog(gui, mngr.titleView, mngr.colorView)
	close(mngr.labelChan)

	return nil
}

func (mngr *labelDialogMngr) commitLabelDialog(gui *gocui.Gui, view *gocui.View) error {
	gui.DeleteKeybindings(labelTitle)
	gui.DeleteKeybindings(labelColor)

	title := strings.Trim(mngr.titleView.Buffer(), "\n ")

	colorIdx := utils.SelectedItemIdx(mngr.colorView)
	color, err := mngr.colorView.Line(colorIdx)
	if err != nil { //Pressed enter before selecting color
		return mngr.closeLabelDialog(gui, view)
	}
	color = color[2:]

	utils.ErrCheck(
		gui.DeleteView(labelTitle),
		gui.DeleteView(labelColor),
	)

	mngr.labelChan <- [2]string{title, color}
	close(mngr.labelChan)
	return nil
}
