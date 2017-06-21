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

func RemoveLabel(labels []trello.Label, idx int) []trello.Label {
	copy(labels[idx:], labels[idx+1:])
	return labels[:len(labels)-1]
}

func RemoveInt(slice []int, idx int) []int {
	copy(slice[idx:], slice[idx+1:])
	return slice[:len(slice)-1]
}

var FgColors = [8]string{
	"-", "red", "green", "-", "-", "purple", "-", "-",
}

var HiFgColors = [8]string{
	"black", "orange", "lime", "yellow", "blue", "pink", "sky", "",
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
