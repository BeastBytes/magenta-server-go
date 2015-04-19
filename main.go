package main

import (
	"log"
	"os"
)

var quit chan bool

func main() {
	quit = make(chan bool)
	port := getPort()

	server := NewServer(port)
	server.Start()

	for {
		select {
		case quit := <-quit:
			if quit == true {
				os.Exit(1)
			}
		}
	}
}

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
