package command

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"text/template"
)

// FillTemplateMessage interpolates data into a complex string
// and makes it more polished. It abstracts away the nuances
// of templating from its users.
func FillTemplateMessage(msg string, filler map[string]interface{}) (string, error) {
	builder := &strings.Builder{}
	tpl := template.Must(template.New("").Parse(msg))

	if err := tpl.Execute(builder, filler); err != nil {
		return msg, err
	}
	return strings.TrimSpace(builder.String()), nil
}

// Coalesce returns the first non-empty string
// it encounters amongst the arguments supplied to the function.
func Coalesce(args ...string) string {
	for _, v := range args {
		if v != "" {
			return v
		}
	}
	return ""
}

// ValidateUrl returns an error if the Parameter supplied to it
// is not a valid URL. Because this function is only a wrapper
// around a standard library function, it doesn't need to be
// tested.
func ValidateUrl(u string) error {
	_, err := url.ParseRequestURI(u)
	return err
}

// ValidateEmailID returns an error if the Parameter supplied to
// it is not a valid Email ID.
func ValidateEmailID(email string) error {
	rxEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) > 254 || !rxEmail.MatchString(email) {
		return errors.New("invalid email id")
	}
	return nil
}
