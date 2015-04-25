package main

type Command struct {
	Name string
	Cmd  func(name string, params ...string)
}

// NewCommand returns a new Command
func NewCommand(name string, cmd func(params ...string)) *Command {
	return &Command(name, cmd)
}
