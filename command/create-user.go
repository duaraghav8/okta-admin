package command

import (
	"context"
	"flag"
	"fmt"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
	"strings"
)

type CreateUserCommand struct{}

type CreateUserCommandConfig struct {
	Team                         string
	Domain, ApiToken             string
	FirstName, LastName, EmailID string
}

func (c *CreateUserCommand) Synopsis() string {
	return "Create a new user in the organization"
}

func (c *CreateUserCommand) Help() string {
	helpText := `
Usage: okta-admin create-user [options]

  Invites a new user to the Organization.
  Okta sends out an invite to the specified Email ID.

Options:

  -domain    Okta organization domain (eg- https://foo.okta.com/)
  -api-token Token to authenticate with Okta API
  -email     Email ID of the user to invite
  -fname     First name of the user to invite (Default: Default)
  -lname     Last name of the user to invite (Default: User)
  -team      The team in the organization the user is part of
`

	return strings.TrimSpace(helpText)
}

func (c *CreateUserCommand) ParseArgs(args []string) (*CreateUserCommandConfig, error) {
	var cfg CreateUserCommandConfig
	flags := flag.NewFlagSet("create-user", flag.ContinueOnError)

	flags.StringVar(&cfg.Domain, "domain", "", "")
	flags.StringVar(&cfg.ApiToken, "api-token", "", "")
	flags.StringVar(&cfg.EmailID, "email", "", "")
	flags.StringVar(&cfg.Team, "team", "", "")
	flags.StringVar(&cfg.FirstName, "fname", "Default", "")
	flags.StringVar(&cfg.LastName, "lname", "User", "")

	if err := flags.Parse(args); err != nil {
		return &cfg, err
	}
	return &cfg, requiredArgs(cfg.Domain, cfg.ApiToken, cfg.EmailID, cfg.Team)
}

func (c *CreateUserCommand) Run(args []string) int {
	cfg, err := c.ParseArgs(args)
	if err != nil {
		fmt.Printf("Failed to parse arguments: %v\n", err)
		return 1
	}

	client, err := okta.NewClient(context.Background(), okta.WithOrgUrl(cfg.Domain), okta.WithToken(cfg.ApiToken))
	if err != nil {
		fmt.Printf("Failed to initialize Okta client: %v\n", err)
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
	user, _, err := client.User.CreateUser(okta.User{Profile: &profile}, queries)
	if err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
		return 1
	}

	fmt.Printf("ID: %s\n", user.Id)
	return 0
}
