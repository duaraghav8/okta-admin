package command

import (
	"context"
	"fmt"
	"github.com/duaraghav8/okta-admin/common"
	oktaapi "github.com/duaraghav8/okta-admin/okta"
	"github.com/okta/okta-sdk-golang/okta"
	"net/http"
)

type DeactivateUserCommand struct {
	Meta *common.CommandMetadata
}

type DeactivateUserCommandConfig struct {
	EmailID string
}

func (c *DeactivateUserCommand) Synopsis() string {
	return "Deactivate an organization member"
}

func (c *DeactivateUserCommand) Help() string {
	helpText := `
Usage: okta-admin deactivate-user [options]

  Deactivates an organization member but doesn't delete
  their account.
{{.GlobalOptionsHelpText}}
Options:

  -email Email ID of the user to deactivate
`

	res, err := common.PrepareMessage(
		helpText,
		map[string]interface{}{
			"GlobalOptionsHelpText": c.Meta.GlobalOptionsHelpText,
		},
	)
	if err != nil {
		return fmt.Sprintf("Failed to render help message: %v\n", err)
	}
	return res
}

func (c *DeactivateUserCommand) ParseArgs(args []string) (*DeactivateUserCommandConfig, error) {
	var cfg DeactivateUserCommandConfig

	flags := c.Meta.FlagSet
	flags.StringVar(&cfg.EmailID, "email", "", "")

	if err := flags.Parse(args); err != nil {
		return &cfg, err
	}
	return &cfg, common.RequiredArgs(map[string]string{
		"email":     cfg.EmailID,
		"org url":   c.Meta.GlobalOptions.OrgUrl,
		"api token": c.Meta.GlobalOptions.ApiToken,
	})
}

func (c *DeactivateUserCommand) Run(args []string) int {
	cfg, err := c.ParseArgs(args)
	if err != nil {
		fmt.Printf("Failed to parse arguments: %v\n", err)
		return 1
	}

	client, err := okta.NewClient(
		context.Background(),
		okta.WithOrgUrl(c.Meta.GlobalOptions.OrgUrl),
		okta.WithToken(c.Meta.GlobalOptions.ApiToken),
	)
	if err != nil {
		fmt.Printf("Failed to initialize Okta client: %v\n", err)
		return 1
	}

	// Fetch user ID
	user, _, err := oktaapi.GetUserByEmail(
		&oktaapi.Credentials{
			OrgUrl:   c.Meta.GlobalOptions.OrgUrl,
			ApiToken: c.Meta.GlobalOptions.ApiToken,
		},
		cfg.EmailID,
	)
	if err != nil {
		fmt.Printf("Failed to resolve user ID: %v\n", err)
		return 1
	}

	// Deactivate user
	resp, err := client.User.DeactivateUser(user["id"].(string), nil)
	if err != nil {
		fmt.Printf("Failed to deactivate member: %v\n", err)
		return 1
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to deactivate member: %v\n", resp)
		return 1
	}

	fmt.Printf("Successfully deactivated %s (ID: %s)\n", cfg.EmailID, user["id"])
	return 0
}
