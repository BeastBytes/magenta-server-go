package main

type Command struct {
	Name string
	Cmd  func(c *Client, s *Server, words []string)
}

// NewCommand returns a new Command
func NewCommand(name string, cmd func(c *Client, s *Server, words []string)) *Command {
	return &Command{Name: name, Cmd: cmd}
}
