package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func testMetaGlobalOptValues(t *testing.T, args []string, expected map[string]string) {
	t.Helper()

	meta, err := createMeta(log.New(ioutil.Discard, "", 0))
	if err != nil {
		t.Fatalf("Failed to create metadata: %v", err)
	}

	flags := meta.FlagSet
	if err := flags.Parse(args); err != nil {
		t.Fatalf("Failed to parse args: %v", err)
	}

	if meta.GlobalOptions.OrgUrl != expected["org_url"] {
		t.Errorf("Org URL: expected %s, received %s", expected["org_url"], meta.GlobalOptions.OrgUrl)
	}
	if meta.GlobalOptions.ApiToken != expected["api_token"] {
		t.Errorf("API token: expected %s, received %s", expected["api_token"], meta.GlobalOptions.ApiToken)
	}
}

func TestCreateMeta(t *testing.T) {

	const (
		orgUrl   = "https://foo.okta.com/"
		apiToken = "123456789abcdxyz"
	)
	expected := map[string]string{
		"org_url":   orgUrl,
		"api_token": apiToken,
	}

	t.Run("parses global options supplied as args", func(t *testing.T) {
		t.Parallel()
		args := []string{
			"-org-url", orgUrl,
			"-api-token", apiToken,
		}
		testMetaGlobalOptValues(t, args, expected)
	})

	// This test shouldn't be run in parallel. Because it manipulates
	// process environment for the duration of its run, it must
	// run in isolation.
	t.Run("parses global options supplied as env vars", func(t *testing.T) {
		var (
			originalEnvOrgUrl   = os.Getenv("OKTA_ORG_URL")
			originalEnvApiToken = os.Getenv("OKTA_API_TOKEN")
		)

		if err := os.Setenv("OKTA_ORG_URL", orgUrl); err != nil {
			t.Fatal("Unable to set org url env var")
		}
		if err := os.Setenv("OKTA_API_TOKEN", apiToken); err != nil {
			t.Fatal("Unable to set api token env var")
		}

		testMetaGlobalOptValues(t, []string{}, expected)

		// cleanup
		if err := os.Setenv("OKTA_ORG_URL", originalEnvOrgUrl); err != nil {
			t.Fatal("Unable to set org url env var back to original value")
		}
		if err := os.Setenv("OKTA_API_TOKEN", originalEnvApiToken); err != nil {
			t.Fatal("Unable to set api token env var back to original value")
		}
	})

}
