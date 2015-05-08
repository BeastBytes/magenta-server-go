package main

import "fmt"

type Channel struct {
	users []User
	join  chan User
	part  chan User
}

// NewChannel initializes a Channel object and returns it
// to the caller
func NewChannel() *Channel {
	return &Channel{
		users: make([]User, 0),
		join:  make(chan User),
		part:  make(chan User),
	}
}

// listen montiors the Channels channels for activity
func (c *Channel) listen() {
	for {
		select {
		case j := <-c.join:
			c.addUser(j)
		case p := <-c.part:
			c.removeUser(p)
		}
	}
}

// sendToChannel sends a message to all clients in the Channel
func (c *Channel) sendToChannel(msg Message) {
	for _, user := range c.users {
		user.ChannelOut() <- fmt.Sprintf("%s: %s", msg.client.Nickname()+msg.msg)
	}
}

// removeClient removes a Client from the Channels client list
func (c *Channel) removeUser(user User) {
	for i, u := range c.users {
		if u.Nickname() == user.Nickname() {
			// delete the client out of the slice.  We do not want to
			// set the value to nil here because the client may still
			// exist in other channels
			c.users = append(c.users[:i], c.users[i+1:]...)
		}
	}
}

// addClient adds a Client to the list of Clients that are currently
// in the Channel
func (c *Channel) addUser(user User) {
	c.users = append(c.users, user)
}

// Users returns the users in the channel
func (c *Channel) Users() []User {
	return c.users
}
