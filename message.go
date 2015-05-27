package main

import "strings"

type Message struct {
	client *Client
	msg    string
}

// NewMessage returns a new Message
func NewMessage(client *Client, msg string) *Message {
	return &Message{client, msg}
}

// Remove trailing spaces, tabs, return, and newline characters from a string
// Thanks to hyphenated (#go-nuts) for pointing out that I
// should also trim off the \r from the input
func trimMessage(msg string) string {
	return strings.TrimRight(msg, " \t\r\n")
}

// isEmpty is a helper function to check for an empty string
func isEmpty(msg string) bool {
	if msg == "" {
		return true
	}

	return false
}
