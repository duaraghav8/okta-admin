package command

import (
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
