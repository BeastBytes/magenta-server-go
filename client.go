/************************************************
* client.go
* Author: Jeramy Singleton
* Date: 12 April 2015
*
* Description:  A client is a remote user that
* has connected to the server.  Information for
* the client should be stored in a struct and
* be accessible to the server.
*************************************************/

package main

import "net"

type Client struct {
	Name string
	Conn net.Conn
}
