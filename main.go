package main

import (
	"log"
	"os"
)

func main() {
	quit := make(chan bool)
	port := getPort()

	InitCommands()
	InitStatusCodeMessages()

	server := NewServer(port, quit)
	server.Run()

	for {
		select {
		case <-quit:
			server.Stop()
			os.Exit(1)
		}
	}
}

// getPort attempts to parse a port from any arguments passed during
// startup.  If no port is included then the server uses the default
// port.
func getPort() string {
	var port string
	if len(os.Args) >= 2 {
		log.Printf("Attempting to listen on port: %s\n", os.Args[1])
		port = ":" + os.Args[1]
	} else {
		log.Printf("Attempting to listen on port: %s", DefaultPort)
		port = DefaultPort
	}

	return port
}
