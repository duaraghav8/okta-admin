package command

import (
	"testing"
)

func createTestResetUserPasswordCommand(globalOptsHelpText string) *ResetUserPasswordCommand {
	return &ResetUserPasswordCommand{
		Command: createTestCommand(globalOptsHelpText, "test_reset_user_password_cmd"),
	}
}

func TestResetUserPasswordCommand_Help(t *testing.T) {
	t.Parallel()
	c := createTestResetUserPasswordCommand(testHelpMessage)
	testCommandHelp(t, c.Help())
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
