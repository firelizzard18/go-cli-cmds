package subcmds

import (
	"github.com/urfave/cli"
	"gitlab.com/ayufan/golang-cli-helpers"
)

type Executor interface {
	Execute(c *cli.Context) error
}

type Collection interface {
	Commands() []cli.Command
}

type ExecutorWithCommands interface {
	Executor
	Collection
}

type RegisterableCollection interface {
	Collection
	RegisterCommand(cmd cli.Command)
	Register(name, usage string, flags ...cli.Flag) RegistrationContext
}

type RegistrationContext interface {
	Configure(cfg func(cli.Command)) RegistrationContext
	Executor(exec Executor)
	Commands(coll Collection)
	ExecutorWithCommands(execColl ExecutorWithCommands)
}

type collection struct {
	commands []cli.Command
}

func NewCollection() RegisterableCollection {
	return &collection{[]cli.Command{}}
}

func (c *collection) Commands() []cli.Command {
	return c.commands
}

func (c *collection) RegisterCommand(cmd cli.Command) {
	c.commands = append(c.commands, cmd)
}

func (c *collection) Register(name, usage string, flags ...cli.Flag) RegistrationContext {
	return &context{c, cli.Command{
		Name:  name,
		Usage: usage,
		Flags: flags,
	}}
}

type context struct {
	collection RegisterableCollection
	command    cli.Command
}

func (c *context) Configure(cfg func(cli.Command)) RegistrationContext {
	cfg(c.command)
	return c
}

func (c *context) Executor(exec Executor) {
	c.command.Flags = append(c.command.Flags, clihelpers.GetFlagsFromStruct(exec)...)
	c.command.Action = exec.Execute
	c.collection.RegisterCommand(c.command)
}

func (c *context) Commands(coll Collection) {
	c.command.Flags = append(c.command.Flags, clihelpers.GetFlagsFromStruct(coll)...)
	c.command.Subcommands = coll.Commands()
	c.collection.RegisterCommand(c.command)
}

func (c *context) ExecutorWithCommands(execColl ExecutorWithCommands) {
	c.command.Flags = append(c.command.Flags, clihelpers.GetFlagsFromStruct(execColl)...)
	c.command.Action = execColl.Execute
	c.command.Subcommands = execColl.Commands()
	c.collection.RegisterCommand(c.command)
}
