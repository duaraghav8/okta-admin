package command

import (
	"context"
	"fmt"
	"github.com/duaraghav8/okta-admin/common"
	"github.com/okta/okta-sdk-golang/okta"
	"net/http"
	"strings"
)

type ListGroupsCommand struct {
	Command
}

type ListGroupsCommandConfig struct {
	Detailed   bool
	GroupNames []string
}

func (c *ListGroupsCommand) Synopsis() string {
	return "List groups in the organization"
}

func (c *ListGroupsCommand) Help() string {
	helpText := `
Usage: okta-admin list-groups [options]

  Lists existing groups in the organization.
  If no arguments are specified, this subcommand lists the
  names of all groups.
{{.GlobalOptionsHelpText}}
Options:

  -groups   Comma-separated list of group names to return info
            about. This is usually combined with -detailed to
            get more information about the specified groups.
            If left unspecified, all groups are listed.
  -detailed Whether to display detailed information about the
            groups. If unspecified, only Group Names are returned.
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

func (c *ListGroupsCommand) ParseArgs(args []string) (*ListGroupsCommandConfig, error) {
	var cfg ListGroupsCommandConfig
	var groupNames string

	flags := c.Meta.FlagSet
	flags.StringVar(&groupNames, "groups", "", "")
	flags.BoolVar(&cfg.Detailed, "detailed", false, "")

	if err := flags.Parse(args); err != nil {
		return &cfg, err
	}
	cfg.GroupNames = GetGroupNames(groupNames, GroupNameSep)

	return &cfg, common.RequiredArgs(map[string]string{
		"org url":   c.Meta.GlobalOptions.OrgUrl,
		"api token": c.Meta.GlobalOptions.ApiToken,
	})
}

func (c *ListGroupsCommand) Run(args []string) int {
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

	groups, resp, err := client.Group.ListGroups(nil)
	if err != nil {
		c.Meta.Logger.Printf("Failed to fetch groups list: %v\n", err)
		return 1
	}
	if resp.StatusCode != http.StatusOK {
		c.Meta.Logger.Printf("Failed to fetch groups list: %s\n", resp.Status)
		return 1
	}

	// Filter groups if names are supplied
	if len(cfg.GroupNames) > 0 {
		groups = FilterGroups(groups, func(g *okta.Group, i int) bool {
			for _, n := range cfg.GroupNames {
				if g.Profile.Name == n {
					return true
				}
			}
			return false
		})
	}

	for _, g := range groups {
		if cfg.Detailed {
			c.Meta.Logger.Println(GetGroupDetailsPretty(g))
			c.Meta.Logger.Println(strings.Repeat("=", 40))
		} else {
			c.Meta.Logger.Println(g.Profile.Name)
		}
	}

	return 0
}
