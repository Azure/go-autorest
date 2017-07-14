package utils

import (
	"fmt"
	"os"
)

// OnErrorFail prints a failure message and exits the program if err is not nil.
func OnErrorFail(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %s", message, err)
		os.Exit(1)
	}
}
