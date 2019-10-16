package command

import (
	"flag"
	"github.com/duaraghav8/okta-admin/common"
	"strings"
	"testing"
)

func createTestResetUserPasswordCommand(globalOptsHelpText string) *ResetUserPasswordCommand {
	m := &common.CommandMetadata{
		GlobalOptionsHelpText: globalOptsHelpText,
		GlobalOptions: &common.CommandConfig{
			OrgUrl:   "https://foo.okta.com/",
			ApiToken: "123abc",
		},
		FlagSet: flag.NewFlagSet("test_reset_user_password_cmd", flag.ContinueOnError),
	}
	return &ResetUserPasswordCommand{Meta: m}
}

func TestResetUserPasswordCommand_Help(t *testing.T) {
	t.Parallel()
	const globalHelpMsg = `
Welcome to Hogwarts!
`
	c := createTestResetUserPasswordCommand(globalHelpMsg)
	if !strings.Contains(c.Help(), globalHelpMsg) {
		t.Errorf("Expected final help message to contain \"%s\"", globalHelpMsg)
	}
}

func TestResetUserPasswordCommand_ParseArgs(t *testing.T) {
	t.Parallel()

	c := createTestDeactivateUserCommand("")
	args := []string{
		"-email", "harry.potter@hogwarts.co.uk",
	}

	cfg, err := c.ParseArgs(args)
	if err != nil {
		t.Fatalf("Failed to parse arguments: %v", err)
	}

	if cfg.EmailID != args[1] {
		t.Errorf("Expected email id to be %s, received %s", args[1], cfg.EmailID)
	}
}