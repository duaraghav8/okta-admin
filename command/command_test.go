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
