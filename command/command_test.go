package command

import (
	"flag"
	"io/ioutil"
	"log"
)

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
