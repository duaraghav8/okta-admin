package command

type CreateUserCommand struct {}

func (c *CreateUserCommand) Help() string {
	return ""
}

func (c *CreateUserCommand) Run(args []string) int {
	return 0
}

func (c *CreateUserCommand) Synopsis() string {
	return "Create a new user in the organization"
}