package command

import (
	oktaapi "github.com/duaraghav8/okta-admin/okta"
	"net/http"
)

type ResetUserMultifactorsCommand struct {
	*Command
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

	return c.Command.prepareHelpMessage(
		helpText,
		map[string]interface{}{
			"GlobalOptionsHelpText": c.Meta.GlobalOptionsHelpText,
		},
	)
}

func (c *ResetUserMultifactorsCommand) ParseArgs(args []string) (*ResetUserMultifactorsCommandConfig, error) {
	var cfg ResetUserMultifactorsCommandConfig

	flags := c.Meta.FlagSet
	flags.StringVar(&cfg.EmailID, "email", "", "")

	if err := flags.Parse(args); err != nil {
		return &cfg, err
	}
	err := c.Command.validateParameters(
		&parameter{Name: "api-token", Required: true, Value: c.Meta.GlobalOptions.ApiToken},
		&parameter{Name: "email", Required: true, Value: cfg.EmailID, ValidationFunc: ValidateEmailID},
		&parameter{Name: "org-url", Required: true, Value: c.Meta.GlobalOptions.OrgUrl, ValidationFunc: ValidateUrl},
	)
	return &cfg, err
}

func (c *ResetUserMultifactorsCommand) Run(args []string) int {
	cfg, err := c.ParseArgs(args)
	if err != nil {
		c.Logger.Printf("Failed to parse arguments: %v\n", err)
		return 1
	}

	client, err := c.OktaClient()
	if err != nil {
		c.Logger.Printf("Failed to initialize Okta client: %v\n", err)
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
		c.Logger.Printf("Failed to resolve user ID: %v\n", err)
		return 1
	}

	// Reset all Multifactors
	resp, err := client.User.ResetAllFactors(user["id"].(string))
	if err != nil {
		c.Logger.Printf("Failed to reset member's multifactors: %v\n", err)
		return 1
	}
	if resp.StatusCode != http.StatusOK {
		c.Logger.Printf("Failed to reset member's multifactors: %s\n", resp.Status)
		return 1
	}

	c.Logger.Printf("All multifactors for %s have been reset\n", cfg.EmailID)
	return 0
}
