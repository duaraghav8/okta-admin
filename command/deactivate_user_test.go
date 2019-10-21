package command

import (
	"testing"
)

func createTestDeactivateUserCommand(globalOptsHelpText string) *DeactivateUserCommand {
	return &DeactivateUserCommand{
		Command: createTestCommand(globalOptsHelpText, "test_deactivate_user_cmd"),
	}
}

func TestDeactivateUserCommand_Help(t *testing.T) {
	t.Parallel()
	c := createTestDeactivateUserCommand(testHelpMessage)
	testCommandHelp(t, c.Help())
}

func TestDeactivateUserCommand_ParseArgs(t *testing.T) {
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
