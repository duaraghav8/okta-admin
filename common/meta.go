package common

import (
	"flag"
	"log"
)

// CommandConfig contains options that are made available
// to all actions. It is used to pass down global
// configuration.
type CommandConfig struct {
	OrgUrl, ApiToken string
}

// CommandMetadata is used to pass metadata to all actions.
// This ensures that any changes in the structure of
// information passed to actions doesn't force the actions
// to adapt to them.
type CommandMetadata struct {
	Logger                *log.Logger
	FlagSet               *flag.FlagSet
	GlobalOptions         *CommandConfig
	GlobalOptionsHelpText string
}
