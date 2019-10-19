package command

import (
	"github.com/duaraghav8/okta-admin/common"
	oktaapi "github.com/duaraghav8/okta-admin/okta"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
	"strings"
)

type listGroupsResult struct {
	oktaapi.GenericResult
	Groups []*okta.Group
}

type addUserToGroupResult struct {
	oktaapi.GenericResult
	GroupName, GroupId string
}

type NumberOfExistingGroups uint32
type OktaGroups []*okta.Group
type FilterGroupsEvalFunc func(group *okta.Group, i int) bool

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

func GetGroupDetailsPretty(g *okta.Group) string {
	tpl := `
Name:        {{.Name}}
ID:          {{.Id}}
Description: {{.Description}}

Links
  Users: {{.LinkUsers}}
  Apps:  {{.LinkApps}}
`

	res, _ := common.PrepareMessage(tpl, map[string]interface{}{
		"Id":          g.Id,
		"Name":        g.Profile.Name,
		"LinkUsers":   getLink(g, "users"),
		"LinkApps":    getLink(g, "apps"),
		"Description": common.FirstNonEmptyStr(g.Profile.Description, "[None]"),
	})
	return res
}

func getLink(g *okta.Group, linkType string) string {
	links := g.Links.(map[string]interface{})
	return links[linkType].(map[string]interface{})["href"].(string)
}

func listGroups(client *okta.Client, qp *query.Params, ch chan<- *listGroupsResult) {
	groups, resp, err := client.Group.ListGroups(qp)
	ch <- &listGroupsResult{
		Groups:        groups,
		GenericResult: oktaapi.GenericResult{Resp: resp, Err: err},
	}
}

func addUserToGroup(client *okta.Client, uid, gid, gname string, ch chan<- *addUserToGroupResult) {
	resp, err := client.Group.AddUserToGroup(gid, uid)
	ch <- &addUserToGroupResult{
		GroupId:       gid,
		GroupName:     gname,
		GenericResult: oktaapi.GenericResult{Err: err, Resp: resp},
	}
}
