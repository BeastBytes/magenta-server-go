package main

import (
	"net"
	"testing"
)

type client struct {
	nickname   string
	realname   string
	conn       net.Conn
	channelOut chan string
}

// Connection returns the clients connection
func (c *client) Connection() net.Conn {
	return c.conn
}

// ChannelOut returns the clients output channel
func (c *client) ChannelOut() chan string {
	return c.channelOut
}

// Nickname returns a string containing the Clients nickname
func (c *client) Nickname() string {
	return c.nickname
}

// RealName returns a string containing the Clients real name
func (c *client) Realname() string {
	return c.realname
}

// Begin the tests

var channel *Channel
var User1 = &client{"Joe", "Joe Bob", nil, nil}
var User2 = &client{"Chuck", "Chuck Popagopalus", nil, nil}
var User3 = &client{"Alexis", "Alexis Hurly", nil, nil}

func setup() {
	channel = NewChannel("#test")
}

func tearDown() {
	channel = nil
}

// TestAddClient adds a single client to a new channel
// and then checks to make sure that user was added correctly
func TestAddClient(t *testing.T) {
	setup()

	channel.addUser(User1)

	if len(channel.users) != 1 {
		t.Error("The channel should contain \"1\" user.")
	}

	for _, c := range channel.users {
		if c.Nickname() != User1.Nickname() {
			t.Error("Client was not added properly")
		}
	}

	tearDown()
}

// TestRemoveClient adds three clients to the Channel, then removes
// the second one to test that our removeClient method indeed removes
// the correct User
func TestRemoveClient(t *testing.T) {
	setup()
	channel.addUser(User1)
	channel.addUser(User2)
	channel.addUser(User3)

	if len(channel.users) != 3 {
		t.Error("channel.users should contain \"3\" users")
	}

	channel.removeUser(User2)

	if len(channel.users) != 2 {
		t.Error("channel.users should contain \"2\" users after removeUser")
	}

	for _, c := range channel.users {
		if c.Nickname() == User2.Nickname() {
			t.Error("User2 was not removed correctly")
		}
	}

	tearDown()
}

func TestValidChannelName(t *testing.T) {
	validNames := []string{"#test", "#test1", "#tes_t"}
	invalidNames := []string{"t#est", "1test", "_test"}

	for _, validName := range validNames {
		if !IsValidChannelName(validName) {
			t.Errorf("%s should be a valid channel name", validName)
		}
	}

	for _, invalidName := range invalidNames {
		if IsValidChannelName(invalidName) {
			t.Errorf("%s should not be a valid channel name", invalidName)
		}
	}
}
