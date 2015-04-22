package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

type Server struct {
	clients        []*Client
	newConnections chan net.Conn
	incoming       chan Message
	outgoing       chan string
	quit           chan bool
	port           string
	listener       net.Listener
}

// NewServer creates a new server object
func NewServer(port string, quit chan bool) *Server {
	ln := newListener(port)
	server := &Server{
		clients:        make([]*Client, 0),
		newConnections: make(chan net.Conn),
		incoming:       make(chan Message),
		outgoing:       make(chan string),
		quit:           quit,
		port:           port,
		listener:       ln,
	}

	return server
}

// SendToClients passes a message to all the clients
// connected
func (s *Server) SendToClients(message string) {
	for _, client := range s.clients {
		client.output <- message
	}
}

// connect prompts a user for a nickname, creates a new
// client and adds that client to the list of clients
func (s *Server) connect(conn net.Conn) {
	name := promptForNickName(conn)
	client := NewClient(name, conn, s.incoming)
	s.clients = append(s.clients, client)
}

// Run listens for incoming connections and starts
// the loop that processes incoming messages from
// clients
func (s *Server) Run() {
	go s.listen()
	go s.processIncoming()
}

// listen listens for new connections and passes
// those new connections on
func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		checkError(err)
		s.newConnections <- conn
		log.Print("Server: new connection from IP", conn.RemoteAddr().String())
	}
}

// Stop stops the server and frees up resources
func (s *Server) Stop() {
	log.Println("Server: stopping")
	close(s.newConnections)
	for _, c := range s.clients {
		c.Close("Server has shutdown, please come back later")
	}
	s.quit <- true
}

// processIncoming monitors all server channels
// for incoming messages.  Those messages can be
// from clients, or even the server itself
func (s *Server) processIncoming() {
	for {
		select {
		case msg := <-s.incoming:
			s.SendToClients(fmt.Sprintf("%s: %s\n", msg.client.NickName, msg.msg))
		case newConn := <-s.newConnections:
			s.connect(newConn)
		}
	}
}

// newListener returns a new net.Listener
func newListener(port string) net.Listener {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("SUCCESS: Server listening on port: %s", port)

	return ln
}

// Prompt the client for their name and set it in the client struct.
func promptForNickName(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	var name string
	var err error
	validName := false
	for !validName {

		writer.WriteString("What is your name?")
		writer.Flush()

		name, err = reader.ReadString('\n')
		name = strings.TrimSpace(name)

		if isValidName(name) && err == nil {
			validName = true
			log.Println(name, "has connected")
		}
	}
	return name
}

// isValidName checks a name for invalid characters.
func isValidName(name string) bool {
	validRgx := regexp.MustCompile(`(^[A-Za-z]\w+\S*$)`)

	return validRgx.MatchString(name)
}

// Check for an error.  If there is an error, log it, and exit the program
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
