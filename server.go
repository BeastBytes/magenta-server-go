package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	clients        []*Client
	newConnections chan net.Conn
	incoming       chan Message
	outgoing       chan string
	port           string
	listener       net.Listener
}

func NewServer(port string) *Server {
	ln := NewListener(port)
	server := &Server{
		clients:        make([]*Client, 0),
		newConnections: make(chan net.Conn),
		incoming:       make(chan Message),
		outgoing:       make(chan string),
		port:           port,
		listener:       ln,
	}

	return server
}

func (s *Server) SendToClients(message string) {
	for _, client := range s.clients {
		client.output <- message
	}
}

func (s *Server) connect(conn net.Conn) {
	name, _ := promptForNickName(conn)
	client := NewClient(name, conn, s.incoming)
	s.clients = append(s.clients, client)
}

// Listen can be exported so that we can stop the server from
// main or anywhere else and restart it later
func (s *Server) Start() {
	go s.loopThruIncoming()

	go func() {
		for {
			conn, err := s.listener.Accept()
			checkError(err)
			s.newConnections <- conn
		}
	}()

}

func (s *Server) loopThruIncoming() {
	for {
		select {
		case chat := <-s.incoming:
			s.SendToClients(fmt.Sprintf("%s: %s\n", chat.client.NickName, chat.msg))
		case newConn := <-s.newConnections:
			s.connect(newConn)
		}
	}
}

func NewListener(port string) net.Listener {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("SUCCESS: Server listening on port: %s", port)
	}

	return ln
}

// Prompt the client for their name and set it in the client struct.
func promptForNickName(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	writer.WriteString("What is your name?")
	writer.Flush()

	name, err := reader.ReadString('\n')

	name = strings.TrimSpace(name)
	return name, err
}

// Check for an error.  If there is an error, log it, and exit the program
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
