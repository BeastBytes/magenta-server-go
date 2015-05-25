package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

type Server struct {
	clients        []*Client
	channels       map[string]*Channel
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
		channels:       make(map[string]*Channel),
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
		client.Receive(message)
	}
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
			isCommand := s.parseMessage(msg)
			if !isCommand {
				s.SendToClients(fmt.Sprintf("%s: %s\n", msg.client.Nickname(), msg.msg))
			}
		case newConn := <-s.newConnections:
			go s.connect(newConn)
		}
	}
}

// parseMessage checks to see if a msg contains a command
// from the client. If it does contain a command that command
// is run and the method returns true.  If no command is found
// the method returns false
func (s *Server) parseMessage(msg Message) bool {
	words := strings.Split(msg.msg, " ")
	cmd := words[0]

	cmdRgx := regexp.MustCompile(`^[\/]\w+`)

	if cmdRgx.MatchString(cmd) {
		cmd, err := parseCommand(words)
		if err != nil {
			msg.client.Receive(fmt.Sprintf("%s\n", err.Error()))
		} else {
			cmd.Cmd(msg.client, s, words)
		}

		return true
	}

	return false
}

// connect prompts a user for a nickname, creates a new
// client and adds that client to the list of clients
func (s *Server) connect(conn net.Conn) {
	name := promptForNickName(conn)
	client := NewClient(name, conn, s.incoming)
	s.clients = append(s.clients, client)
}

func (s *Server) AddChannel(channel string) {
	if !s.HasChannel(channel) {
		s.channels[channel] = NewChannel(channel)
		go s.channels[channel].listen()
		fmt.Printf("%s channel added\n", channel)
	}
}

func (s *Server) HasChannel(channel string) bool {
	_, ok := s.channels[channel]
	return ok
}

func (s *Server) AddUserToChannel(channel string, user User) {
	if s.HasChannel(channel) {
		c := s.channels[channel]
		c.join <- user
	} else {
		fmt.Printf("%s not a valid channel\n", channel)
	}
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

// parseCommand attempts to match a command request sent by the client
// to any commands contained in the Commands slice.  If found, the
// function returns that command and a nil error.  If not found, it
// returns nil and an invalid command error.
func parseCommand(words []string) (*Command, error) {
	cmd := words[0]
	cmd = strings.TrimLeft(cmd, "/")

	command := GetCommand(strings.ToLower(cmd))

	if command == nil {
		return nil, errors.New("cmd: invalid command")
	}

	if len(words) <= 1 {
		return nil, errors.New("cmd: invalid number of parameters")
	}

	return command, nil
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

// Check for an error.  If there is an error, log it, and exit the program
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
