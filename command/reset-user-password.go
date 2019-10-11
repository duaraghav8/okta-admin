package command

type ResetUserPasswordCommand struct {}

func (c *ResetUserPasswordCommand) Help() string {
	return ""
}

func (c *ResetUserPasswordCommand) Run(args []string) int {
	return 0
}

func (c *ResetUserPasswordCommand) Synopsis() string {
	return "Reset organization member's password"
}
