package main

import (
	"fmt"
	"os"

	cmd "github.com/duaraghav8/okta-admin/command"
	"github.com/mitchellh/cli"
)

func main() {
	meta, err := createMeta()
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
				return &cmd.ResetUserPasswordCommand{
					Meta: meta,
				}, nil
			},
			"deactivate-user": func() (command cli.Command, err error) {
				return &cmd.DeactivateUserCommand{
					Meta: meta,
				}, nil
			},
			"reset-user-mfa": func() (command cli.Command, err error) {
				return &cmd.ResetUserMultifactorsCommand{
					Meta: meta,
				}, nil
			},
			"assign-groups": func() (command cli.Command, err error) {
				return &cmd.AssignUserGroupsCommand{
					Meta: meta,
				}, nil
			},
			"assign-apps": func() (command cli.Command, err error) {
				return &cmd.AssignApplicationsCommand{}, nil
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
