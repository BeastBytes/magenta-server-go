package main

const (
	// DefaultPort for the server to listen on
	DefaultPort = ":2500"
	// MaxMessageLength is the maximum size of a message
	// accepted from a client
	MaxMessageLength = 512
	// MaxNameLength is the maximum size of a clients name
	MaxNameLength = 16
	// IdleTime is the amount of time a client is inactive
	// before it is considered idle.
	IdleTime = 5 // in minutes
)
