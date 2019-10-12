package command

type ResetUserMultifactorCommand struct {}

func (c *ResetUserMultifactorCommand) Help() string {
	return ""
}

func (c *ResetUserMultifactorCommand) Run(args []string) int {
	return 0
}

func (c *ResetUserMultifactorCommand) Synopsis() string {
	return "Reset organization member's Multifactor"
}
