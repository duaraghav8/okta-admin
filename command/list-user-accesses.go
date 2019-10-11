package command

type ListUserAccessesCommand struct {}

func (c *ListUserAccessesCommand) Help() string {
	return ""
}

func (c *ListUserAccessesCommand) Run(args []string) int {
	return 0
}

func (c *ListUserAccessesCommand) Synopsis() string {
	return "List applications assigned to an organization member"
}
