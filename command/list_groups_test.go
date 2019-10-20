package command

import (
	"flag"
	"github.com/duaraghav8/okta-admin/common"
	"strings"
	"testing"
)

func createTestListGroupsCommand(globalOptsHelpText string) *ListGroupsCommand {
	m := &common.CommandMetadata{
		GlobalOptionsHelpText: globalOptsHelpText,
		GlobalOptions: &common.CommandConfig{
			OrgUrl:   "https://foo.okta.com/",
			ApiToken: "123abc",
		},
		FlagSet: flag.NewFlagSet("test_list_groups_cmd", flag.ContinueOnError),
	}
	return &ListGroupsCommand{Meta: m}
}

func TestListGroupsCommand_Help(t *testing.T) {
	t.Parallel()
	const globalHelpMsg = `
Welcome to Hogwarts!
`
	c := createTestListGroupsCommand(globalHelpMsg)
	if !strings.Contains(c.Help(), globalHelpMsg) {
		t.Errorf("Expected final help message to contain \"%s\"", globalHelpMsg)
	}
}

func TestListGroupsCommand_ParseArgs(t *testing.T) {
	t.Run("without groups", func(t *testing.T) {
		t.Parallel()

		c := createTestListGroupsCommand("")
		args := []string{"-detailed"}

		cfg, err := c.ParseArgs(args)
		if err != nil {
			t.Fatalf("Failed to parse arguments: %v", err)
		}

		if !cfg.Detailed {
			t.Errorf("Expected -detailed flag to be set")
		}
		if len(cfg.GroupNames) != 0 {
			t.Errorf("Expected group names slice to be empty, received %v", cfg.GroupNames)
		}
	})

	t.Run("with groups", func(t *testing.T) {
		t.Parallel()

		var groups = []string{"Tech", "office_staff", "Marketing", "talent-acquisition"}
		c := createTestListGroupsCommand("")
		args := []string{
			"-groups", "Tech\t,    office_staff  ,\t\tMarketing  ,  talent-acquisition",
		}

		cfg, err := c.ParseArgs(args)
		if err != nil {
			t.Fatalf("Failed to parse arguments: %v", err)
		}

		if cfg.Detailed {
			t.Errorf("Expected -detailed flag to be unset")
		}
		if len(cfg.GroupNames) != len(groups) {
			t.Fatalf("Expected %d group names, received %d", len(groups), len(cfg.GroupNames))
		}
		for i := 0; i < len(groups); i++ {
			if cfg.GroupNames[i] != groups[i] {
				t.Errorf("Expected GroupNames[%d] to be %s, received %s", i, groups[i], cfg.GroupNames[i])
			}
		}
	})
}
