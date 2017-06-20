package utils

import (
	. "github.com/jroimartin/gocui"
	"strings"
	"strconv"
	"github.com/aqatl/KogoKolejka/utils"
)

func DelNonGlobalKeyBinds(gui *Gui) {
	for _, view := range gui.Views() {
		gui.DeleteKeybindings(view.Name())
	}
}

func SelectedItemIdx(view *View) int {
	if len(strings.Split(view.Buffer(), "\n")) <= 1 {
		return -1
	}

	_, cy := view.Cursor()
	currLine, err := view.Line(cy)
	dotIdx := strings.Index(currLine, ".")
	itemIdx64, err := strconv.ParseInt(currLine[:dotIdx], 10, 32)
	utils.ErrCheck(err)
	return int(itemIdx64)
}
