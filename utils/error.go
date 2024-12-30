package utils

import "errors"

func Assert(err error) {
	if err != nil {
		panic(err)
	}
}

func ErrorIf(condition bool) {
	if condition {
		panic(errors.New("conditional error"))
	}
}
