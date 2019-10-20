package command

import (
	oktaapi "github.com/duaraghav8/okta-admin/okta"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
	"strings"
)

// listGroupsResult contains the result of an async HTTP request
// made to Okta API to fetch list of Groups.
type listGroupsResult struct {
	oktaapi.GenericResult
	Groups []*okta.Group
}

// addUserToGroup contains the result of an async HTTP request
// made to Okta API to add a User to a Group.
type addUserToGroupResult struct {
	oktaapi.GenericResult
	GroupName, GroupId string
}

type NumberOfExistingGroups uint32
type OktaGroups []*okta.Group
type FilterGroupsEvalFunc func(group *okta.Group, i int) bool

// GroupNameSep is the string by which the group name list
// supplied as raw input is split into individual names.
const GroupNameSep = ","

// GetID returns the ID of the Group whose name is specified.
// If the Group with that name doesn't exist, this function
// simply returns an empty string.
func (groups OktaGroups) GetID(name string) string {
	for _, g := range groups {
		if g.Profile.Name == name {
			return g.Id
		}
	}
	return ""
}

// FilterGroups filters Okta Groups based on a user-supplied
// evaluation function.
func FilterGroups(groups OktaGroups, eval FilterGroupsEvalFunc) OktaGroups {
	res := make(OktaGroups, 0, len(groups))
	for i, g := range groups {
		if eval(g, i) {
			res = append(res, g)
		}
	}
	return res
}

// SanitizeGroupNames takes raw user-supplied group names
// as input and prepares them for further processing.
func SanitizeGroupNames(names []string) []string {
	n := make([]string, len(names), len(names))
	for i := 0; i < len(names); i++ {
		n[i] = strings.TrimSpace(names[i])
	}
	return n
}

// GetGroupNames takes a list of group names as a raw string
// and returns a structured list of the same.
func GetGroupNames(rawInput, sep string) []string {
	if strings.TrimSpace(rawInput) == "" {
		return []string{}
	}
	return SanitizeGroupNames(strings.Split(rawInput, sep))
}

// GetGroupDetailsPretty returns a pretty string describing
// the group passed to it.
func GetGroupDetailsPretty(g *okta.Group) string {
	tpl := `
Name:        {{.Name}}
ID:          {{.Id}}
Description: {{.Description}}

Links
  Users: {{.LinkUsers}}
  Apps:  {{.LinkApps}}
`

	res, _ := PrepareMessage(tpl, map[string]interface{}{
		"Id":          g.Id,
		"Name":        g.Profile.Name,
		"LinkUsers":   getLinkFromGroup(g, "users"),
		"LinkApps":    getLinkFromGroup(g, "apps"),
		"Description": FirstNonEmptyStr(g.Profile.Description, "[None]"),
	})
	return res
}

// getLinkFromGroup returns a specific type of link from the
// group passed to it. It abstracts away the nuances of
// typecasting Links to retrieve data.
func getLinkFromGroup(g *okta.Group, linkType string) string {
	links := g.Links.(map[string]interface{})
	return links[linkType].(map[string]interface{})["href"].(string)
}

// listGroups fetches the list of Groups from Okta API asynchronously
func listGroups(client *okta.Client, qp *query.Params, ch chan<- *listGroupsResult) {
	groups, resp, err := client.Group.ListGroups(qp)
	ch <- &listGroupsResult{
		Groups:        groups,
		GenericResult: oktaapi.GenericResult{Resp: resp, Err: err},
	}
}

// addUserToGroup adds a user to a group using the Okta API asynchronously
func addUserToGroup(client *okta.Client, uid, gid, gname string, ch chan<- *addUserToGroupResult) {
	resp, err := client.Group.AddUserToGroup(gid, uid)
	ch <- &addUserToGroupResult{
		GroupId:       gid,
		GroupName:     gname,
		GenericResult: oktaapi.GenericResult{Err: err, Resp: resp},
	}
}
