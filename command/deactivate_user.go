package command

type DeactivateUserCommand struct {}

func (c *DeactivateUserCommand) Help() string {
	return ""
}

func (c *DeactivateUserCommand) Run(args []string) int {
	return 0
}

func (c *DeactivateUserCommand) Synopsis() string {
	return "Deactivate an organization member"
}
