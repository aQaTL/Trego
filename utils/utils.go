package utils

import "log"

func ErrCheck(errs ...error) {
	for _, err := range errs {
		if err != nil {
			log.Panicln(err)
		}
	}
}
