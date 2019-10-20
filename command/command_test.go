package command

import "flag"

func createTestCommand(globalOptsHelpText, flagSetName string) *Command {
	m := &Metadata{
		GlobalOptionsHelpText: globalOptsHelpText,
		GlobalOptions: &Config{
			OrgUrl:   "https://foo.okta.com/",
			ApiToken: "123abc",
		},
		FlagSet: flag.NewFlagSet(flagSetName, flag.ContinueOnError),
	}
	return &Command{Meta: m}
}
