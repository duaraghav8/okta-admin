package command

import (
	"testing"
)

func createTestCreateUserCommand(globalOptsHelpText string) *CreateUserCommand {
	return &CreateUserCommand{
		Command: createTestCommand(globalOptsHelpText, "test_create_user_cmd"),
	}
}

func TestCreateUserCommand_Help(t *testing.T) {
	t.Parallel()
	c := createTestCreateUserCommand(testHelpMessage)
	testCommandHelp(t, c.Help())
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
		t.Errorf("Expected first name to be %s, received %s", args[1], cfg.FirstName)
	}
	if cfg.LastName != args[3] {
		t.Errorf("Expected last name to be %s, received %s", args[1], cfg.LastName)
	}
	if cfg.Team != args[5] {
		t.Errorf("Expected team to be %s, received %s", args[1], cfg.Team)
	}
	if cfg.EmailID != args[7] {
		t.Errorf("Expected email id to be %s, received %s", args[1], cfg.EmailID)
	}
}
