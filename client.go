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
	serverInput chan string
	output      chan string
}

func NewClient(name string, conn net.Conn, serverInput chan string) *Client {
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

// Run two loops, concurrently, that monitor the Client's input and
// output channels
func (c *Client) listen() {
	go c.GetInput()
	go c.SendOutput()
}

// When input is received from the Client send it to the server
func (c *Client) GetInput() {
	for {
		// The buffer will stop reading when it detects a new line chatacter
		input, _ := c.reader.ReadString('\n')

		// TODO: In the future, the server will be responsible for attaching
		// the Client's NickName to the message
		c.serverInput <- c.NickName + ": " + input
	}
}

// When output is detected from the server send it to the Client
func (c *Client) SendOutput() {
	for data := range c.output {
		c.writer.WriteString(data)
		c.writer.Flush()
	}
}
