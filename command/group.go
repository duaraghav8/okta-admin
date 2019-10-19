package command

import (
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

// SanitizeGroupNames takes raw user-supplied group names
// as input and prepares them for further processing.
func SanitizeGroupNames(names []string) []string {
	n := make([]string, len(names), len(names))
	for i := 0; i < len(names); i++ {
		n[i] = strings.TrimSpace(names[i])
	}
	return n
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
