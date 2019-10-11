package main

import (
	"fmt"
	"os"

	cmd "github.com/duaraghav8/okta-admin/command"
	"github.com/mitchellh/cli"
)

func main() {
	c := cli.CLI{
		Name:    "okta-admin",
		Version: "1.0.0",
		Commands: map[string]cli.CommandFactory{
			"create-user": func() (command cli.Command, err error) {
				return &cmd.CreateUserCommand{}, nil
			},
			"reset-user-password": func() (command cli.Command, err error) {
				return &cmd.ResetUserPasswordCommand{}, nil
			},
			"deactivate-user": func() (command cli.Command, err error) {
				return &cmd.DeactivateUserCommand{}, nil
			},
			"list-user-accesses": func() (command cli.Command, err error) {
				return &cmd.ListUserAccessesCommand{}, nil
			},
			"assign-applications": func() (command cli.Command, err error) {
				return &cmd.AssignApplicationsCommand{}, nil
			},
			"reset-user-multifactor": func() (command cli.Command, err error) {
				return &cmd.ResetUserMultifactorCommand{}, nil
			},
		},
		Args: os.Args[1:],
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(exitStatus)
}
