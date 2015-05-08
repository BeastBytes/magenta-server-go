package main

import "fmt"

var (
	Commands = make([]*Command, 0)
)

// RegisterCommand adds a new command to the Commands slice
func RegisterCommand(cmd *Command) {
	Commands = append(Commands, cmd)
}

// GetCommand attempts to find a command in the Commands slice
func GetCommand(name string) *Command {
	for _, c := range Commands {
		if c.Name == name {
			return c
		}
	}

	return nil
}

// InitCommands registers commands with the Commands slice
func InitCommands() {
	RegisterCommand(NewCommand("join", func(c *Client, words []string) {
		fmt.Printf("%s wants to join %s\n", c.Nickname(), words[1])
	}))

	RegisterCommand(NewCommand("part", func(c *Client, words []string) {
		fmt.Printf("%s wants to part %s\n", c.Nickname(), words[1])
	}))
}
