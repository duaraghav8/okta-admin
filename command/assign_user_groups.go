package command

import (
	"context"
	"fmt"
	"github.com/duaraghav8/okta-admin/common"
	oktaapi "github.com/duaraghav8/okta-admin/okta"
	"github.com/okta/okta-sdk-golang/okta"
	"net/http"
	"strings"
)

const GroupNameSep = ","

type AssignUserGroupsCommand struct {
	Meta *common.CommandMetadata
}

type AssignUserGroupsCommandConfig struct {
	EmailID    string
	GroupNames []string
}

func (c *AssignUserGroupsCommand) Synopsis() string {
	return "Assign groups to user"
}

func (c *AssignUserGroupsCommand) Help() string {
	helpText := `
Usage: okta-admin assign-groups [options]

  Adds an organization member to one or more groups.
  This assumes that the specified group(s) already exist
  in the organization. If no groups are specified, this
  subcommand does nothing.
{{.GlobalOptionsHelpText}}
Options:

  -email  Email ID of the user to deactivate
  -groups Comma-separated list of groups to assign to the user
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

func (c *AssignUserGroupsCommand) ParseArgs(args []string) (*AssignUserGroupsCommandConfig, error) {
	var cfg AssignUserGroupsCommandConfig
	var groupNames string

	flags := c.Meta.FlagSet
	flags.StringVar(&cfg.EmailID, "email", "", "")
	flags.StringVar(&groupNames, "groups", "", "")

	if err := flags.Parse(args); err != nil {
		return &cfg, err
	}

	// Extract group names from the single comma-separated string
	// received by the -groups flag.
	if groupNames == "" {
		cfg.GroupNames = []string{}
	} else {
		cfg.GroupNames = SanitizeGroupNames(strings.Split(groupNames, GroupNameSep))
	}

	return &cfg, common.RequiredArgs(map[string]string{
		"email":     cfg.EmailID,
		"org url":   c.Meta.GlobalOptions.OrgUrl,
		"api token": c.Meta.GlobalOptions.ApiToken,
	})
}

func (c *AssignUserGroupsCommand) Run(args []string) int {
	var (
		user   oktaapi.ApiResponse
		groups = OktaGroups{}

		getUserCh        = make(chan *getUserResult)
		listGroupsCh     = make(chan *listGroupsResult)
		addUserToGroupCh = make(chan *addUserToGroupResult)
	)
	var neg NumberOfExistingGroups

	cfg, err := c.ParseArgs(args)
	if err != nil {
		c.Meta.Logger.Printf("Failed to parse arguments: %v\n", err)
		return 1
	}
	if len(cfg.GroupNames) == 0 {
		c.Meta.Logger.Println("No groups were specified, nothing to do")
		return 0
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

	// Fetch User info and list of groups in the organization
	go listGroups(client, nil, listGroupsCh)
	go getUser(cfg.EmailID,
		c.Meta.GlobalOptions.OrgUrl, c.Meta.GlobalOptions.ApiToken, getUserCh)

	// The first issue encountered should stop further execution
	for i := 0; i < 2; i++ {
		select {
		case u := <-getUserCh:
			if u.Err != nil {
				c.Meta.Logger.Printf("Failed to resolve u ID: %v\n", err)
				return 1
			}
			user = u.User
		case g := <-listGroupsCh:
			if g.Err != nil {
				c.Meta.Logger.Printf("Failed to fetch list of groups: %v\n", err)
				return 1
			}
			if g.Resp.StatusCode != http.StatusOK {
				c.Meta.Logger.Printf("Failed to fetch list of groups: %v\n", g.Resp)
				return 1
			}
			groups = g.Groups
		}
	}

	neg = NumberOfExistingGroups(len(cfg.GroupNames))
	for _, n := range cfg.GroupNames {
		gid := groups.GetID(n)
		if gid == "" {
			c.Meta.Logger.Printf("%s does not exist\n", n)
			neg--
			continue
		}
		go addUserToGroup(client, user["id"].(string), gid, n, addUserToGroupCh)
	}

	for i := 0; i < int(neg); i++ {
		added := <-addUserToGroupCh
		if added.Err != nil {
			c.Meta.Logger.Printf("Failed to add user to %s: %v\n", added.GroupName, added.Err)
		} else if added.Resp.StatusCode != http.StatusNoContent {
			c.Meta.Logger.Printf("Failed to add user to %s: %s\n", added.GroupName, added.Resp.Status)
		} else {
			c.Meta.Logger.Printf("Added to %s\n", added.GroupName)
		}
	}

	return 0
}
