package command

import (
	"context"
	"fmt"
	"github.com/duaraghav8/okta-admin/common"
	oktaapi "github.com/duaraghav8/okta-admin/okta"
	"github.com/okta/okta-sdk-golang/okta"
	"net/http"
)

type ResetUserPasswordCommand struct {
	Command
}

type ResetUserPasswordCommandConfig struct {
	EmailID string
}

func (c *ResetUserPasswordCommand) Synopsis() string {
	return "Reset organization member's password"
}

func (c *ResetUserPasswordCommand) Help() string {
	helpText := `
Usage: okta-admin reset-user-password [options]

  Resets password of an organization member.
  Okta emails a password reset link to the specified member.
{{.GlobalOptionsHelpText}}
Options:

  -email Email ID of the organization member
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

func (c *ResetUserPasswordCommand) ParseArgs(args []string) (*ResetUserPasswordCommandConfig, error) {
	var cfg ResetUserPasswordCommandConfig

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

func (c *ResetUserPasswordCommand) Run(args []string) int {
	cfg, err := c.ParseArgs(args)
	if err != nil {
		c.Meta.Logger.Printf("Failed to parse arguments: %v\n", err)
		return 1
	}

	client, err := okta.NewClient(
		context.Background(),
		okta.WithOrgUrl(c.Meta.GlobalOptions.OrgUrl),
		okta.WithToken(c.Meta.GlobalOptions.ApiToken),
	)
	if err != nil {
		c.Meta.Logger.Printf("Failed to initialize Okta client: %v\n", err)
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
		c.Meta.Logger.Printf("Failed to resolve user ID: %v\n", err)
		return 1
	}

	// Reset password
	_, resp, err := client.User.ResetPassword(user["id"].(string), nil)
	if err != nil {
		c.Meta.Logger.Printf("Failed to reset member's password: %v\n", err)
		return 1
	}
	if resp.StatusCode != http.StatusOK {
		c.Meta.Logger.Printf("Failed to reset member's password: %v\n", resp)
		return 1
	}

	c.Meta.Logger.Printf("Reset link sent to %s\n", cfg.EmailID)
	return 0
}
