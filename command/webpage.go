package command

import (
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

type OpenWebpageCommand struct {
	*Command
}

type OpenWebpageCommandConfig struct {
	Type string
}

var Webpages = map[string]string{
	"home":    "/admin/dashboard",
	"api":     "/admin/access/api/tokens",
	"apps":    "/admin/apps/active",
	"logs":    "/report/system_log_2",
	"people":  "/admin/users",
	"groups":  "/admin/groups",
	"admins":  "/admin/access/admins",
	"account": "/admin/settings/account",
}

func getWebpageTypes() []string {
	types := make([]string, 0, len(Webpages))
	for t := range Webpages {
		types = append(types, t)
	}
	return types
}

func getOpenCommand(os string) (string, error) {
	var (
		cmd string
		err error
	)

	// TODO: Add support for windows os
	switch os {
	case "linux":
		cmd = "xdg-open"
	case "darwin":
		cmd = "open"
	default:
		err = fmt.Errorf("unsupported os: %s", os)
	}

	return cmd, err
}

func (c *OpenWebpageCommand) Synopsis() string {
	return "Open a specific webpage of your organization"
}

func (c *OpenWebpageCommand) Help() string {
	helpText := `
Usage: okta-admin webpage <type>

  Open a specific webpage of your Okta organization.
  The page is opened in the default app on your system configured
  to handle okta.com URIs. Usually this is your default web browser.
  This command does not require the org API token to be set.

  Possible values for type are:
  {{.Webpages}}

  If type is omitted, the admin dashboard is opened. Note that you
  should already be authenticated with Okta and authorized to access
  these pages.
{{.GlobalOptionsHelpText}}
`

	return c.Command.prepareHelpMessage(
		helpText,
		map[string]interface{}{
			"Webpages":              strings.Join(getWebpageTypes(), ", "),
			"GlobalOptionsHelpText": c.Meta.GlobalOptionsHelpText,
		},
	)
}

func (c *OpenWebpageCommand) ParseArgs(args []string) (*OpenWebpageCommandConfig, error) {
	var cfg OpenWebpageCommandConfig
	flags := c.Meta.FlagSet

	if err := flags.Parse(args); err != nil {
		return &cfg, err
	}
	err := c.Command.validateParameters(
		&parameter{Name: "org-url", Required: true, Value: c.Meta.GlobalOptions.OrgUrl, ValidationFunc: ValidateUrl},
	)

	args = flags.Args()
	if len(args) == 0 {
		cfg.Type = "home"
	} else {
		cfg.Type = args[0]
	}

	return &cfg, err
}

func (c *OpenWebpageCommand) Run(args []string) int {
	var openCmd string

	cfg, err := c.ParseArgs(args)
	if err != nil {
		c.Logger.Printf("Failed to parse arguments: %v\n", err)
		return 1
	}

	pageType := strings.ToLower(cfg.Type)
	route, ok := Webpages[pageType]
	if !ok {
		c.Logger.Printf("Invalid page type '%s'", pageType)
		return 1
	}

	uri := path.Join(c.Meta.GlobalOptions.OrgUrl, route)
	if openCmd, err = getOpenCommand(runtime.GOOS); err != nil {
		c.Logger.Printf("Cannot proceed to open %s: %s", uri, err)
		return 1
	}
	if err = exec.Command(openCmd, uri).Start(); err != nil {
		c.Logger.Printf("Failed to open %s: %s", uri, err)
	}
	c.Logger.Println(uri)

	return 0
}
