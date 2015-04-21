package main

import "strings"

type Message struct {
	client *Client
	msg    string
}

func NewMessage(client *Client, msg string) *Message {
	return &Message{client, msg}
}

// Remove the return and newline characters from a string
func trimMessage(msg string) string {
	if !isEmpty(msg) {
		// Thanks to hyphenated (#go-nuts) for pointing out that I should also
		// trim off the \r from the input
		msg = strings.TrimRight(msg, "\r\n")
	}
	return msg
}

// isEmpty is a helper function to check for an empty string
func isEmpty(msg string) bool {
	if msg == "" {
		return true
	}

	return false
}
