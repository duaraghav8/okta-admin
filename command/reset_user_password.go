package command

import (
	oktaapi "github.com/duaraghav8/okta-admin/okta"
	"net/http"
)

type ResetUserPasswordCommand struct {
	*Command
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

	return c.Command.prepareHelpMessage(
		helpText,
		map[string]interface{}{
			"GlobalOptionsHelpText": c.Meta.GlobalOptionsHelpText,
		},
	)
}

func (c *ResetUserPasswordCommand) ParseArgs(args []string) (*ResetUserPasswordCommandConfig, error) {
	var cfg ResetUserPasswordCommandConfig

	flags := c.Meta.FlagSet
	flags.StringVar(&cfg.EmailID, "email", "", "")

	if err := flags.Parse(args); err != nil {
		return &cfg, err
	}

	err := c.Command.validateParameters(
		&Parameter{Name: "api-token", Required: true, Value: c.Meta.GlobalOptions.ApiToken},
		&Parameter{Name: "email", Required: true, Value: cfg.EmailID, ValidationFunc: ValidateEmailID},
		&Parameter{Name: "org-url", Required: true, Value: c.Meta.GlobalOptions.OrgUrl, ValidationFunc: ValidateUrl},
	)
	return &cfg, err
}

func (c *ResetUserPasswordCommand) Run(args []string) int {
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

	// Reset password
	_, resp, err := client.User.ResetPassword(user["id"].(string), nil)
	if err != nil {
		c.Logger.Printf("Failed to reset member's password: %v\n", err)
		return 1
	}
	if resp.StatusCode != http.StatusOK {
		c.Logger.Printf("Failed to reset member's password: %s\n", resp.Status)
		return 1
	}

	c.Logger.Printf("Reset link sent to %s\n", cfg.EmailID)
	return 0
}
