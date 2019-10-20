package command

import (
	"errors"
	"flag"
	"fmt"
	"github.com/duaraghav8/okta-admin/common"
	"log"
)

// Command contains objects passed to all CLI commands
type Command struct {
	Meta   *Metadata
	Logger *log.Logger
}

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

// requiredArgs takes a series of arguments and returns an
// error upon encountering the first empty argument.
// This function ensures that all specified arguments are
// non-empty.
func (c *Command) requiredArgs(args map[string]string) error {
	for k, v := range args {
		if v == "" {
			return errors.New(fmt.Sprintf("%s is a required argument", k))
		}
	}
	return nil
}

func (c *Command) prepareHelpMessage(helpText string, filler map[string]interface{}) string {
	res, err := common.PrepareMessage(helpText, filler)
	if err != nil {
		return fmt.Sprintf("Failed to render help message: %v\n", err)
	}
	return res
}
