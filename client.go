//everything related to individual clients
// a new person joinig, will be created here in the backend

package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type Client struct {
	conn     *websocket.Conn
	manager  *Manager
	Outgoing chan Event //the queue for each client
	//Message
}

func NewClientService(co *websocket.Conn, man *Manager) *Client {
	return &Client{
		conn:     co,
		manager:  man,
		Outgoing: make(chan Event),
	}
}

// read messages from the client

func (c *Client) ReadMessages() {
	defer func() {
		//we will clean up the connection after the client has run into an err or closes the connection while reading is going on
		c.manager.DeleteClient(c)
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("there was an error in reading messages: %v", err)
			}
			break
		}
		var request Event
		if err = json.Unmarshal(msg, &request); err != nil {
			log.Print(err)
		}

		//c.manager.Broadcast(request)
		if err := c.manager.RouteEvents(request, c); err != nil {
			log.Panicf("error handling messages: %v", err)
		}
	}

}

func (c *Client) WriteMessage() {
	defer func() {
		//c.conn.Close()
		c.manager.DeleteClient(c)
	}()
	for wmsg := range c.Outgoing {

		data, err := json.Marshal(wmsg)
		if err != nil {
			log.Println(err)
			return
		}
		if err := c.conn.WriteMessage(1, data); err != nil {
			log.Printf("error writing message to client: %v", err)
			break
		}
		log.Println("message sent")
	}

}
