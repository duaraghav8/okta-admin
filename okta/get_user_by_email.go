package okta

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// GetUserByEmail returns information about the user associated
// with the specified email ID.
func GetUserByEmail(c *Credentials, email string) (ApiResponse, *http.Response, error) {
	var user ApiResponse
	client := &http.Client{}
	endpoint := fmt.Sprintf("/api/v1/users/%s", email)

	reqUrl, err := CreateRequestUrl(c.OrgUrl, endpoint)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("unable to create request: %v", err))
	}

	for n, v := range CreateRequestHeaders(c.ApiToken) {
		req.Header.Set(n, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return ApiResponse{}, resp, errors.New(fmt.Sprintf("failed to fetch user (%s)", resp.Status))
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, resp, errors.New(fmt.Sprintf("failed to read API response: %v", err))
	}
	return user, resp, nil
}
