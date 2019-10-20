package main

import (
	"log"
	"os"

	cmd "github.com/duaraghav8/okta-admin/command"
	"github.com/duaraghav8/okta-admin/version"
	"github.com/mitchellh/cli"
)

func main() {
	logger := log.New(os.Stdout, "", 0)
	meta, err := createMeta(logger)
	if err != nil {
		logger.Printf("Failed to create metadata for actions: %v\n", err)
		os.Exit(1)
	}

	c := cli.CLI{
		Name:    version.AppName,
		Version: version.FormattedVersion(),
		Commands: map[string]cli.CommandFactory{
			"create-user": func() (command cli.Command, err error) {
				return &cmd.CreateUserCommand{
					Meta: meta,
				}, nil
			},
			"deactivate-user": func() (command cli.Command, err error) {
				return &cmd.DeactivateUserCommand{
					Meta: meta,
				}, nil
			},
			"reset-user-password": func() (command cli.Command, err error) {
				return &cmd.ResetUserPasswordCommand{
					Meta: meta,
				}, nil
			},
			"reset-user-mfa": func() (command cli.Command, err error) {
				return &cmd.ResetUserMultifactorsCommand{
					Meta: meta,
				}, nil
			},
			"list-groups": func() (command cli.Command, err error) {
				return &cmd.ListGroupsCommand{
					Meta: meta,
				}, nil
			},
			"assign-groups": func() (command cli.Command, err error) {
				return &cmd.AssignUserGroupsCommand{
					Meta: meta,
				}, nil
			},
		},
		Args: os.Args[1:],
	}

	exitStatus, err := c.Run()
	if err != nil {
		logger.Println(err)
	}

	os.Exit(exitStatus)
}
