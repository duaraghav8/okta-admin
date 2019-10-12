package main

import (
	"fmt"
	"os"

	cmd "github.com/duaraghav8/okta-admin/command"
	"github.com/mitchellh/cli"
)

func main() {
	meta, err := createMeta(os.Args[1:])
	if err != nil {
		fmt.Printf("Failed to create metadata for actions: %v\n", err)
		os.Exit(1)
	}

	c := cli.CLI{
		Name:    AppName,
		Version: Version,
		Commands: map[string]cli.CommandFactory{
			"create-user": func() (command cli.Command, err error) {
				return &cmd.CreateUserCommand{
					Meta: meta,
				}, nil
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
