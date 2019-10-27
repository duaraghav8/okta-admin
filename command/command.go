package command

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/okta/okta-sdk-golang/okta"
	"log"
	"strings"
)

// Config contains options that are made available
// to all actions. It is used to pass down global
// configuration.
type Config struct {
	OrgUrl, ApiToken string
}

// Metadata is used to pass metadata to all actions.
// This makes making structural changes in data passed
// to commands easy.
type Metadata struct {
	FlagSet               *flag.FlagSet
	GlobalOptions         *Config
	GlobalOptionsHelpText string
}

// Command contains objects passed to all CLI commands
type Command struct {
	Meta       *Metadata
	Logger     *log.Logger
	oktaClient *okta.Client
}

// Parameter represents a commandline Parameter with full
// context.
type Parameter struct {
	Required       bool
	Name, Value    string
	ValidationFunc func(value string) error
}

// ValueSep is the string separating individual values in a
// raw string.
const ValueSep = ","

// OktaClient returns an instance of Okta Client initialized
// with organization-specific API credentials. This method
// only creates the client the first time it is called.
// Subsequent calls return the cached client.
func (c *Command) OktaClient() (*okta.Client, error) {
	if c.oktaClient != nil {
		return c.oktaClient, nil
	}

	if c.Meta.GlobalOptions.OrgUrl == "" {
		return nil, errors.New("org URL cannot be empty")
	}
	if c.Meta.GlobalOptions.ApiToken == "" {
		return nil, errors.New("api token cannot be empty")
	}

	client, err := okta.NewClient(context.Background(),
		okta.WithOrgUrl(c.Meta.GlobalOptions.OrgUrl), okta.WithToken(c.Meta.GlobalOptions.ApiToken))
	if err != nil {
		// Cache the newly created client
		c.oktaClient = client
	}
	return client, err
}

func (c *Command) prepareHelpMessage(helpText string, filler map[string]interface{}) string {
	res, err := FillTemplateMessage(helpText, filler)
	if err != nil {
		return fmt.Sprintf("Failed to render help message: %v\n", err)
	}
	return res
}

// parseListOfValues takes a raw string containing multiple
// values for a single commandline option and returns a list
// of those individual values.
func (c *Command) parseListOfValues(rawInput, sep string) []string {
	if strings.TrimSpace(rawInput) == "" {
		return []string{}
	}

	l := strings.Split(rawInput, sep)
	res := make([]string, len(l), len(l))
	for i := 0; i < len(l); i++ {
		res[i] = strings.TrimSpace(l[i])
	}

	return res
}

func (c *Command) validateParameters(params ...*Parameter) error {
	for _, p := range params {
		// Return error if a required param is not set
		if p.Required && p.Value == "" {
			return errors.New(fmt.Sprintf("%s is required", p.Name))
		}
		// If validation func is supplied and param is set,
		// validate it.
		if p.ValidationFunc != nil && p.Value != "" {
			if err := p.ValidationFunc(p.Value); err != nil {
				return err
			}
		}
	}
	return nil
}
