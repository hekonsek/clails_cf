package util

import (
	"fmt"
	"os"
)

const UnixExitCodeGeneralError = 1

func CliError(err error) bool {
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return true
	}
	return false
}

func ExitOnCliError(err error) {
	if CliError(err) {
		os.Exit(UnixExitCodeGeneralError)
	}
}
