package main

import (
	"fmt"
	"strings"
)

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

	// join command tells the server that the user wishes to join, or create,
	// a certain channel
	RegisterCommand(NewCommand("join", func(c *Client, s *Server, words []string) {
		channel := words[1]
		if IsValidChannelName(channel) {
			if !s.HasChannel(channel) {
				s.AddChannel(channel)
			}
			s.AddUserToChannel(channel, c)
		} else {
			c.ChannelOut() <- "Invalid channel name"
		}
	}))

	// part commmand sends the server a message that the user is ready to leave
	// the channel passed
	RegisterCommand(NewCommand("part", func(c *Client, s *Server, words []string) {
		channel := words[1]
		if IsValidChannelName(channel) && s.HasChannel(channel) {
			s.RemoveUserFromChannel(channel, c)
		}
	}))

	// channel command is divided into subcommands.
	// The first subcommand is "users".
	// "users" sends a list of all users in a channel
	RegisterCommand(NewCommand("channel", func(c *Client, s *Server, words []string) {
		chanName := words[1]

		var subCmd string
		if len(words) > 2 {
			subCmd = words[2]
		}

		fmt.Println(subCmd)

		if subCmd != "" {
			switch subCmd {
			// get a list of all users in the channel
			case "users":
				channel, err := s.getChannel(chanName)

				if err == nil {
					channelStats := fmt.Sprintf("%s has %d users.\n-----------------\n", chanName, len(channel.Users()))
					users := make([]string, 0)
					for _, u := range channel.Users() {
						users = append(users, u.Nickname())
					}
					usersString := strings.Join(users, ", ")
					c.ChannelOut() <- channelStats + usersString + "\n"
				}
			}
		}
	}))
}
