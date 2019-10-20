package command

import (
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
	"net/http"
)

type CreateUserCommand struct {
	*Command
}

type CreateUserCommandConfig struct {
	EmailID             string
	Team                string
	FirstName, LastName string
}

func (c *CreateUserCommand) Synopsis() string {
	return "Create a new user in the organization"
}

func (c *CreateUserCommand) Help() string {
	helpText := `
Usage: okta-admin create-user [options]

  Invites a new user to the Organization.
  Okta sends out an invite to the specified Email ID.
{{.GlobalOptionsHelpText}}
Options:

  -email Email ID of the user to invite
  -fname First name of the user to invite (Default: Default)
  -lname Last name of the user to invite (Default: User)
  -team  The team in the organization the user should be part of
`

	return c.Command.prepareHelpMessage(
		helpText,
		map[string]interface{}{
			"GlobalOptionsHelpText": c.Meta.GlobalOptionsHelpText,
		},
	)
}

func (c *CreateUserCommand) ParseArgs(args []string) (*CreateUserCommandConfig, error) {
	var cfg CreateUserCommandConfig
	flags := c.Meta.FlagSet

	flags.StringVar(&cfg.EmailID, "email", "", "")
	flags.StringVar(&cfg.Team, "team", "", "")
	flags.StringVar(&cfg.FirstName, "fname", "Default", "")
	flags.StringVar(&cfg.LastName, "lname", "User", "")

	if err := flags.Parse(args); err != nil {
		return &cfg, err
	}
	return &cfg, c.Command.requiredArgs(map[string]string{
		"email":     cfg.EmailID,
		"team":      cfg.Team,
		"org url":   c.Meta.GlobalOptions.OrgUrl,
		"api token": c.Meta.GlobalOptions.ApiToken,
	})
}

func (c *CreateUserCommand) Run(args []string) int {
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

	queries := query.NewQueryParams(query.WithActivate(true))
	profile := okta.UserProfile{
		"team":      cfg.Team,
		"email":     cfg.EmailID,
		"login":     cfg.EmailID,
		"firstName": cfg.FirstName,
		"lastName":  cfg.LastName,
	}
	user, resp, err := client.User.CreateUser(okta.User{Profile: &profile}, queries)
	if err != nil {
		c.Logger.Printf("Failed to create user: %v\n", err)
		return 1
	}
	if resp.StatusCode != http.StatusOK {
		c.Logger.Printf("Failed to create user: %s\n", resp.Status)
		return 1
	}

	c.Logger.Printf("ID: %s\n", user.Id)
	return 0
}
