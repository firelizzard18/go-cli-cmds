package subcmds

import (
	"github.com/urfave/cli"
)

var commands = NewCollection()

func Commands() []cli.Command {
	return commands.Commands()
}

func RegisterCommand(cmd cli.Command) {
	commands.RegisterCommand(cmd)
}

func Register(name, usage string, flags ...cli.Flag) RegistrationContext {
	return commands.Register(name, usage, flags...)
}
