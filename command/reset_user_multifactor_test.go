package command

import (
	"testing"
)

func createTestResetUserMultifactorsCommand(globalOptsHelpText string) *ResetUserMultifactorsCommand {
	return &ResetUserMultifactorsCommand{
		Command: createTestCommand(globalOptsHelpText, "test_reset_user_mfa_cmd"),
	}
}

func TestResetUserMultifactorsCommand_Help(t *testing.T) {
	t.Parallel()
	c := createTestResetUserMultifactorsCommand(testHelpMessage)
	testCommandHelp(t, c.Help())
}

func TestResetUserMultifactorsCommand_ParseArgs(t *testing.T) {
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
