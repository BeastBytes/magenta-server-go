package main

import (
	"log"
	"net"
	"os"
)

const (
	// Default port for the server to listen on
	DEFAULT_PORT = ":2500"
	// The maximum size of a message accepted from a client
	MAX_MESSAGE_LENGTH = 512
	// The Maximum size of a clients name
	MAX_NAME_LENGTH = 16
)

var client Client

func main() {
	var port string

	if len(os.Args) >= 2 {
		log.Printf("Attempting to listen on port: %s\n", os.Args[1])
		port = ":" + os.Args[1]
	} else {
		log.Printf("Attempting to listen on port: %s", DEFAULT_PORT)
		port = DEFAULT_PORT
	}

	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("SUCCESS: Server listening on port: %s", port)
	}

	connection, err := ln.Accept()
	checkError(err)

	handleClient(connection)

	for {

	}

	log.Printf("Exiting...")
}

func handleClient(conn net.Conn) {
	var buffer [MAX_MESSAGE_LENGTH]byte

	//Prompt the client for a name
	// TODO: Eventually we will have the client send their name
	// upon connection without having the server prompt them for it
	promptClientForName(conn)

	for {
		buflen, err := conn.Read(buffer[0:])
		// if for some reason we cannot read the clients input
		// we should log the err and close the clients connection
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}

		// Print the clients input to the console and also return it to the client
		sendToClient(conn, string(buffer[:buflen]))

	}
}

// Send message to the client
func sendToClient(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Println(err)
		conn.Close()
	}
}

// Prompt the client for their name and set it in the client struct.
func promptClientForName(conn net.Conn) {
	conn.Write([]byte("What is your name? "))

	var buffer [MAX_NAME_LENGTH]byte
	buflen, err := conn.Read(buffer[0:])
	if err != nil {
		log.Println(err)
		conn.Write([]byte("There was an error with your input, please reconnect and try again.\n"))
		conn.Close()
		return
	}

	client.Name = string(buffer[:buflen])
}

// Check for an error.  If there is an error, log it, and exit the program
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
