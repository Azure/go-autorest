package utils

import (
	"fmt"
)

// OnErrorPanic prints a failure message and exits the program if err is not nil.
func OnErrorPanic(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", message, err))
	}
}
