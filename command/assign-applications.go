package command

type AssignApplicationsCommand struct {}

func (c *AssignApplicationsCommand) Help() string {
	return ""
}

func (c *AssignApplicationsCommand) Run(args []string) int {
	return 0
}

func (c *AssignApplicationsCommand) Synopsis() string {
	return "Assign applications to an organization member"
}
