package command

import (
	"context"
	"fmt"
	"github.com/duaraghav8/okta-admin/common"
	oktaapi "github.com/duaraghav8/okta-admin/okta"
	"github.com/okta/okta-sdk-golang/okta"
	"net/http"
)

type ResetUserMultifactorsCommand struct {
	Meta *common.CommandMetadata
}

type ResetUserMultifactorsCommandConfig struct {
	EmailID string
}

func (c *ResetUserMultifactorsCommand) Synopsis() string {
	return "Reset all Multifactors of an organization member"
}

func (c *ResetUserMultifactorsCommand) Help() string {
	helpText := `
Usage: okta-admin reset-user-mfa [options]

  Resets all Multifactors of an organization member.
  Okta doesn't notify the user of this reset explicity,
  but rather lets them setup their Multifactors when they
  log into the domain post MFA reset.
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

func (c *ResetUserMultifactorsCommand) ParseArgs(args []string) (*ResetUserMultifactorsCommandConfig, error) {
	var cfg ResetUserMultifactorsCommandConfig

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

func (c *ResetUserMultifactorsCommand) Run(args []string) int {
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

	// Reset all Multifactors
	resp, err := client.User.ResetAllFactors(user["id"].(string))
	if err != nil {
		fmt.Printf("Failed to reset member's multifactors: %v\n", err)
		return 1
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to reset member's multifactors: %v\n", resp)
		return 1
	}

	fmt.Printf("All multifactors for %s have been reset\n", cfg.EmailID)
	return 0
}
