package command

import (
	"flag"
	"fmt"
	"github.com/duaraghav8/okta-admin/common"
	"log"
)

// Command contains data passed to all CLI commands
type Command struct {
	Meta *Metadata
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
	Logger                *log.Logger
	FlagSet               *flag.FlagSet
	GlobalOptions         *Config
	GlobalOptionsHelpText string
}

func (c *Command) prepareHelpMessage(helpText string, filler map[string]interface{}) string {
	res, err := common.PrepareMessage(helpText, filler)
	if err != nil {
		return fmt.Sprintf("Failed to render help message: %v\n", err)
	}
	return res
}
