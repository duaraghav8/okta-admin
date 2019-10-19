// Package okta enables users to make API requests to Okta
// which are not supported by the Okta Go SDK.
package okta

import (
	"errors"
	"fmt"
	"github.com/okta/okta-sdk-golang/okta"
	"net/url"
	"path"
)

// GenericResult represents a struct that can hold results
// returned by any okta sdk function that passes response
// from the upstream API.
type GenericResult struct {
	Err  error
	Resp *okta.Response
}

// Credentials contains all information required to authenticate
// to and access an Okta domain.
type Credentials struct {
	OrgUrl, ApiToken string
}

// ApiResponse represents an arbitrary JSON response object
// from a REST API.
type ApiResponse map[string]interface{}

// CreateRequestUrl creates the URL to call when you
// want to make an Okta API request to a specific endpoint.
func CreateRequestUrl(orgUrl, endpoint string) (string, error) {
	u, err := url.Parse(orgUrl)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to parse organization url: %v", err))
	}
	u.Path = path.Join(u.Path, endpoint)
	return u.String(), nil
}

// CreateRequestHeaders creates HTTP headers required by the Okta
// REST API.
func CreateRequestHeaders(apiToken string) map[string]string {
	return map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("SSWS %s", apiToken),
	}
}
