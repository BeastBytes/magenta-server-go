/************************************************
* client.go
* Author: Jeramy Singleton
* Date: 12 April 2015
*
* Description:  A client is a remote user that
* has connected to the server.  Information for
* the client should be stored in a struct and
* be accessible to the server.

* Reference for Read and Write functionality:
* https://gist.github.com/drewolson/3950226
*************************************************/

package main

import (
	"bufio"
	"net"
)

type Client struct {
	NickName    string
	RealName    string
	conn        net.Conn
	reader      *bufio.Reader
	writer      *bufio.Writer
	serverInput chan Message
	output      chan string
}

func NewClient(name string, conn net.Conn, serverInput chan Message) *Client {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	client := &Client{
		NickName:    name,
		conn:        conn,
		reader:      reader,
		writer:      writer,
		serverInput: serverInput,
		output:      make(chan string),
	}

	// Start the loops that constantly monitor for input and output
	// from the client
	client.listen()

	return client
}

// listen runs two loops, concurrently, that monitor the Client's input and
// output channels
func (c *Client) listen() {
	go c.processInput()
	go c.SendOutput()
}

// GetInput receives data from the Client and sends it to the server
func (c *Client) processInput() {
	for {
		// The buffer will stop reading when it detects a new line chatacter
		input, _ := c.reader.ReadString('\n')

		// If the client sends an empty message just ignore it
		if input := trimMessage(input); !isEmpty(input) {
			msg := NewMessage(c, input)
			c.serverInput <- *msg
		}
	}
}

// SendOutput monitors the output channel.  When data is received from the server
// it is sent to the Client
func (c *Client) SendOutput() {
	for data := range c.output {
		c.writer.WriteString(data)
		c.writer.Flush()
	}
}
