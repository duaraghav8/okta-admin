package common

import (
	"errors"
	"fmt"
	"strings"
	"text/template"
)

// RequiredArgs takes a series of arguments and returns an
// error upon encountering the first empty argument.
// This function ensures that all specified arguments are
// non-empty.
func RequiredArgs(args map[string]string) error {
	for k, v := range args {
		if v == "" {
			return errors.New(fmt.Sprintf("%s is a required argument", k))
		}
	}
	return nil
}

// PrepareMessage interpolates data into a complex string
// and makes it more polished. It abstracts away the
// nuances of templating from its users.
func PrepareMessage(msg string, filler map[string]interface{}) (string, error) {
	builder := &strings.Builder{}
	tpl := template.Must(template.New("").Parse(msg))

	if err := tpl.Execute(builder, filler); err != nil {
		return msg, err
	}
	return strings.TrimSpace(builder.String()), nil
}