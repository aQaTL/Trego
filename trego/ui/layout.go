package ui

import (
	. "github.com/jroimartin/gocui"
	"github.com/fatih/color"
	"math"
)

const (
	BOTTOM_BAR string = "botbar"
	TOP_BAR string = "topbar"
	LIST_WIDTH int = 24
)

type List struct {
	Name  string
	Cards []Card
}

type Card struct {
	Name string
}

var listCounter int

func Layout(gui *Gui) error {
	if err := bottomBar(gui); err != nil {
		return err
	}
	if err := topBar(gui); err != nil {
		return err
	}

	return nil
}

func AddList(gui *Gui, list List) error {
	_, maxY := gui.Size()
	if v, err := gui.SetView(list.Name,
		listCounter * LIST_WIDTH + int(math.Abs(sign(listCounter))), 3,
		listCounter * LIST_WIDTH + LIST_WIDTH, maxY - 5);
			err != nil && err == ErrUnknownView {

		listCounter++
		v.Editable = false
		v.Highlight = true
		v.Wrap = true
		v.BgColor = ColorBlack
		v.Title = list.Name

		color.Output = v

		for idx, card := range (list.Cards) {
			color.New(color.BgBlack).Add(color.FgWhite).Printf("%d. %v\n", idx, card.Name)
		}
	} else {
		return err
	}

	return nil
}

func topBar(gui *Gui) error {
	maxX, _ := gui.Size()
	if v, err := gui.SetView(TOP_BAR, 0, 0, maxX - 1, 2); err != nil {
		if err != ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Highlight = false
		v.BgColor = ColorBlack

		color.Output = v
		color.New(color.FgYellow).Add(color.Bold).Printf("Board: %v", "Trego")
	}

	return nil
}

//bottom bar with shortcuts
func bottomBar(gui *Gui) error {
	maxX, maxY := gui.Size()
	if v, err := gui.SetView(BOTTOM_BAR, 0, maxY - 4, maxX - 1, maxY - 1); err != nil {
		if err != ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Highlight = false
		v.BgColor = ColorBlack

		color.Output = v
		color.New(color.FgYellow).Add(color.Bold).Println("Ala nie ma kota")

		if _, err := gui.SetCurrentView(BOTTOM_BAR); err != nil {
			return err
		}
	}
	return nil
}

func SetKeyBindings(gui *Gui) error {
	if err := gui.SetKeybinding("", KeyCtrlC, ModNone, quit); err != nil {
		return err
	}

	return nil
}

func quit(gui *Gui, v *View) error {
	return ErrQuit
}

func sign(x int) float64 {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	} else {
		return 0
	}
}