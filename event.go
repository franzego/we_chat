package main

import (
	"encoding/json"
)

// all my event related stuff
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// i did this, so that the various events can be routed to the correct event handler
// the eventhandler is a function that handles the events attached to it
type EventHandler func(event Event, c *Client) error

const (
	EventtoSendMessage = "send-message" //whenever a new msg is sent from the client, this is the event that we will send
	EventtoJoinMEssage = "join-message" //whenever a join message is sent
)

type SendMessageEvent struct {
	Message   string `json:"message"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
}

type JoinMessageEvent struct {
	Username string `json:"username"`
}
