package command

import (
	"flag"
	"github.com/duaraghav8/okta-admin/common"
	"strings"
	"testing"
)

func createTestCreateUserCommand(globalOptsHelpText string) *CreateUserCommand {
	m := &common.CommandMetadata{
		GlobalOptionsHelpText: globalOptsHelpText,
		GlobalOptions: &common.CommandConfig{
			OrgUrl:   "https://foo.okta.com/",
			ApiToken: "123abc",
		},
		FlagSet: flag.NewFlagSet("test_create_user_cmd", flag.ContinueOnError),
	}
	return &CreateUserCommand{Meta: m}
}

func TestCreateUserCommand_Help(t *testing.T) {
	t.Parallel()
	const globalHelpMsg = `
Welcome to Hogwarts!
`
	c := createTestCreateUserCommand(globalHelpMsg)
	if !strings.Contains(c.Help(), globalHelpMsg) {
		t.Errorf("Expected final help message to contain \"%s\"", globalHelpMsg)
	}
}

func TestCreateUserCommand_ParseArgs(t *testing.T) {
	t.Parallel()

	c := createTestCreateUserCommand("")
	args := []string{
		"-fname", "Harry",
		"-lname", "Potter",
		"-team", "Gryffindor",
		"-email", "harry.potter@hogwarts.co.uk",
	}

	cfg, err := c.ParseArgs(args)
	if err != nil {
		t.Fatalf("Failed to parse arguments: %v", err)
	}

	if cfg.FirstName != args[1] {
		t.Errorf("Expected first name to be %s, received %s", args[1], cfg.Team)
	}
	if cfg.LastName != args[3] {
		t.Errorf("Expected last name to be %s, received %s", args[1], cfg.Team)
	}
	if cfg.Team != args[5] {
		t.Errorf("Expected team to be %s, received %s", args[1], cfg.Team)
	}
	if cfg.EmailID != args[7] {
		t.Errorf("Expected email id to be %s, received %s", args[1], cfg.Team)
	}
}
