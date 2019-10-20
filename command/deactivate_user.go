package command

import (
	oktaapi "github.com/duaraghav8/okta-admin/okta"
	"net/http"
)

type DeactivateUserCommand struct {
	*Command
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

	return c.Command.prepareHelpMessage(
		helpText,
		map[string]interface{}{
			"GlobalOptionsHelpText": c.Meta.GlobalOptionsHelpText,
		},
	)
}

func (c *DeactivateUserCommand) ParseArgs(args []string) (*DeactivateUserCommandConfig, error) {
	var cfg DeactivateUserCommandConfig

	flags := c.Meta.FlagSet
	flags.StringVar(&cfg.EmailID, "email", "", "")

	if err := flags.Parse(args); err != nil {
		return &cfg, err
	}
	return &cfg, c.Command.requiredArgs(map[string]string{
		"email":     cfg.EmailID,
		"org url":   c.Meta.GlobalOptions.OrgUrl,
		"api token": c.Meta.GlobalOptions.ApiToken,
	})
}

func (c *DeactivateUserCommand) Run(args []string) int {
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

	// Deactivate user
	resp, err := client.User.DeactivateUser(user["id"].(string), nil)
	if err != nil {
		c.Logger.Printf("Failed to deactivate member: %v\n", err)
		return 1
	}
	if resp.StatusCode != http.StatusOK {
		c.Logger.Printf("Failed to deactivate member: %v\n", resp)
		return 1
	}

	c.Logger.Printf("Successfully deactivated %s (ID: %s)\n", cfg.EmailID, user["id"])
	return 0
}
