package main

import "net"

type User interface {
	Nickname() string
	Realname() string
	Connection() net.Conn
	ChannelOut() chan string
}
