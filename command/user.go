package command

import (
	oktaapi "github.com/duaraghav8/okta-admin/okta"
	"net/http"
)

// getUserResult contains the result of an async HTTP request
// made to Okta API to fetch a single User.
type getUserResult struct {
	User oktaapi.ApiResponse
	Resp *http.Response
	Err  error
}

// getUser fetches a single User from Okta API asynchronously
func getUser(email, orgUrl, apiToken string, ch chan<- *getUserResult) {
	user, resp, err := oktaapi.GetUserByEmail(
		&oktaapi.Credentials{
			OrgUrl: orgUrl, ApiToken: apiToken,
		},
		email,
	)
	ch <- &getUserResult{User: user, Resp: resp, Err: err}
}
