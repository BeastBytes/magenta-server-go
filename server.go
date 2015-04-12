package main

import (
	"log"
	"net"
)

const (
	LISTEN_PORT = ":2500"
)

func main() {
	ln, err := net.Listen("tcp", LISTEN_PORT)
	checkError(err)

	connection, err := ln.Accept()
	checkError(err)

	connection.Write([]byte("Hello"))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
