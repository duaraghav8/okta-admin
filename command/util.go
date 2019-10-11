package command

import (
	"errors"
	"fmt"
)

// requiredArgs takes a series of arguments and returns an
// error upon encountering the first empty argument.
// This function ensures that all specified arguments are
// non-empty.
func requiredArgs(args ...string) error {
	for _, arg := range args {
		if arg == "" {
			return errors.New(fmt.Sprintf("%s is a required argument", arg))
		}
	}
	return nil
}
