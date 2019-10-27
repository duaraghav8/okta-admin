package command

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

const testHelpMessage = `
Welcome to Hogwarts!
`

func createTestCommand(globalOptsHelpText, flagSetName string) *Command {
	m := &Metadata{
		GlobalOptionsHelpText: globalOptsHelpText,
		GlobalOptions: &Config{
			OrgUrl:   "https://foo.okta.com/",
			ApiToken: "123abc",
		},
		FlagSet: flag.NewFlagSet(flagSetName, flag.ContinueOnError),
	}
	return &Command{
		Meta:   m,
		Logger: log.New(ioutil.Discard, "", 0),
	}
}

func testCommandHelp(t *testing.T, commandHelpMsg string) {
	t.Helper()
	if !strings.Contains(commandHelpMsg, testHelpMessage) {
		t.Errorf("Expected final help message to contain \"%s\"", testHelpMessage)
	}
}

func testEq(a, b []string) bool {
	// Either both slices must be nil, or both must have the
	// same length
	if ((a == nil) != (b == nil)) || (len(a) != len(b)) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestCommand_parseListOfValues(t *testing.T) {
	t.Parallel()

	cmd := &Command{}
	testCases := []struct {
		input    string
		expected []string
	}{
		{"", []string{}},
		{"      ", []string{}},
		{"\t  \t  ", []string{}},
		{"hello,world", []string{"hello", "world"}},
		{"a ,\tb, c,  d", []string{"a", "b", "c", "d"}},
		{"Harry_Potter,  Ron-Weasley", []string{"Harry_Potter", "Ron-Weasley"}},
		{"a@b.com,c#$641&^*, :-90dg", []string{"a@b.com", "c#$641&^*", ":-90dg"}},
	}

	for _, tc := range testCases {
		if res := cmd.parseListOfValues(tc.input, ","); !testEq(res, tc.expected) {
			t.Errorf("Unexpected output %v for raw input \"%s\"", res, tc.input)
		}
	}
}

func TestCommand_validateParameters(t *testing.T) {
	cmd := &Command{}

	t.Run("required and optional params", func(t *testing.T) {
		t.Parallel()
		param := &parameter{Required: true, Name: "foobar", Value: "", ValidationFunc: nil}

		if err := cmd.validateParameters(param); err == nil {
			t.Error("Expected param to be invalid as required field is not set")
		}

		param.Required = false
		if err := cmd.validateParameters(param); err != nil {
			t.Error("Expected param to be valid as it is not required")
		}
	})
}
