package command

import (
	oktaapi "github.com/duaraghav8/okta-admin/okta"
	"net/http"
)

type getUserResult struct {
	User oktaapi.ApiResponse
	Resp *http.Response
	Err  error
}

func getUser(email, orgUrl, apiToken string, ch chan<- *getUserResult) {
	user, resp, err := oktaapi.GetUserByEmail(
		&oktaapi.Credentials{
			OrgUrl: orgUrl, ApiToken: apiToken,
		},
		email,
	)
	ch <- &getUserResult{User: user, Resp: resp, Err: err}
}
