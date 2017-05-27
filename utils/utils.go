package utils

import (
	"github.com/aqatl/go-trello"
	"log"
)

func ErrCheck(errs ...error) {
	for _, err := range errs {
		if err != nil {
			log.Panicln(err)
		}
	}
}

func RemoveList(lists []trello.List, idx int) []trello.List {
	copy(lists[idx:], lists[idx+1:])
	return lists[:len(lists)-1]
}

var FgColors = []string{
	"null", "red", "dark green", "brown", "dark blue", "purple", "cyan", "",
}

var HiFgColors = []string{
	"black", "orange", "green", "yellow", "blue", "pink", "sky", "white",
}

func MapColor(colorStr string) (color, hi int) {
	for i, col := range FgColors {
		if col == colorStr {
			return i, 2
		}
	}
	for i, col := range HiFgColors {
		if col == colorStr {
			return i, 1
		}
	}
	return 0, 0
}
