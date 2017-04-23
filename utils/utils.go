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
