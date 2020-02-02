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
	meta, err := createMeta()
	if err != nil {
		logger.Printf("Failed to create metadata for actions: %v\n", err)
		os.Exit(1)
	}

	globalCommand := &cmd.Command{
		Meta:   meta,
		Logger: logger,
	}

	c := cli.CLI{
		Name:    version.AppName,
		Version: version.FormattedVersion(),
		Commands: map[string]cli.CommandFactory{
			"create-user": func() (command cli.Command, err error) {
				return &cmd.CreateUserCommand{Command: globalCommand}, nil
			},
			"deactivate-user": func() (command cli.Command, err error) {
				return &cmd.DeactivateUserCommand{Command: globalCommand}, nil
			},
			"reset-user-password": func() (command cli.Command, err error) {
				return &cmd.ResetUserPasswordCommand{Command: globalCommand}, nil
			},
			"reset-user-mfa": func() (command cli.Command, err error) {
				return &cmd.ResetUserMultifactorsCommand{Command: globalCommand}, nil
			},
			"list-groups": func() (command cli.Command, err error) {
				return &cmd.ListGroupsCommand{Command: globalCommand}, nil
			},
			"assign-groups": func() (command cli.Command, err error) {
				return &cmd.AssignUserGroupsCommand{Command: globalCommand}, nil
			},
			"webpage": func() (command cli.Command, err error) {
				return &cmd.OpenWebpageCommand{Command: globalCommand}, nil
			},
		},
		Args:       os.Args[1:],
		HelpWriter: os.Stdout,
	}

	exitStatus, err := c.Run()
	if err != nil {
		logger.Println(err)
	}

	os.Exit(exitStatus)
}
