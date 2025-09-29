//everything related to individual clients
// a new person joining, will be created here in the backend

package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	//the time we wait before we close the connection if there is no reply
	pongWait = 10 * time.Second
	// how often we send pings
	pingInterval = (pongWait * 9) / 10
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
	Username string
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
	//to make sure that extremely large messages are not sent
	c.conn.SetReadLimit(512)
	//the initial pong timer
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	//for us to handle pong messages that will be sent after the ping has been received by the browser and responded to
	c.conn.SetPongHandler(c.pongHandler)

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
			log.Printf("error handling messages: %v", err)
		}
	}

}

// the ponghandler function that actually handles the pong and then resets the original timer
func (c *Client) pongHandler(appdata string) error {
	log.Println("pong")
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}

/*
func (c *Client) WriteMessage() {

	ticker := time.NewTicker(pingInterval)

	defer func() {
		ticker.Stop()
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
*/
func (c *Client) WriteMessage() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.manager.DeleteClient(c)
	}()
	for {
		select {
		case wmsg, ok := <-c.Outgoing:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed")
				}
				return
			}

			data, err := json.Marshal(wmsg)
			if err != nil {
				log.Print(err)
				return
			}
			if err := c.conn.WriteMessage(1, data); err != nil {
				log.Println(err)
			}
			log.Println("sent message")
		case <-ticker.C:
			log.Println("ping")
			//we will send the ping to the client
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("msg: ", err)
				return
			}
		}
	}

}
