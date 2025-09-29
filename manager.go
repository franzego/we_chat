package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/gorilla/websocket"
)

var ctx = context.Background()

type Manager struct {
	Clients ClientList
	sync.RWMutex
	handlers    map[string]EventHandler
	redisClient *redis.Client
	userIndex   map[string]*Client //map of users wiht thier connections
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

func NewManager() *Manager {
	m := &Manager{
		Clients:   make(ClientList),
		handlers:  make(map[string]EventHandler),
		userIndex: make(map[string]*Client),
	}
	m.SetupEventHandlers()
	return m

}

// the function to setup the event handlers
func (m *Manager) SetupEventHandlers() {
	m.handlers[EventtoSendMessage] = SendMessage
	m.handlers[EventtoJoinMEssage] = JoinMessage
}

// the function to create the Eventhandler for the SendMessage Event
func SendMessage(event Event, c *Client) error {
	var sendmsg SendMessageEvent
	if err := json.Unmarshal(event.Payload, &sendmsg); err != nil {
		log.Print("error in unmarshalling sent message")
	}
	//the sent messages are added here to the redis db for persistence
	msgJson, _ := json.Marshal(&sendmsg)
	err := c.manager.redisClient.RPush(ctx, "message-history", msgJson).Err()
	if err != nil {
		log.Println("there was an error writing to redis")
	}

	// have to wrap it with the Event that the c.outgoing can understand
	var broadcastmessage NewMessageEvent
	broadcastmessage.Sent = time.Now()
	broadcastmessage.Message = sendmsg.Message
	broadcastmessage.Sender = sendmsg.Sender
	data, err := json.Marshal(broadcastmessage)
	if err != nil {
		return fmt.Errorf("there was an error marshalling the broadcast message")
	}
	// we need to put the payload as an Event
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventtoSendMessage
	// try and implement a kind of switch to handle group messages and DMs
	switch {
	case sendmsg.Receipient == "":
		for client := range c.manager.Clients {
			client.Outgoing <- outgoingEvent
		}
	default:
		receiver, ok := c.manager.userIndex[sendmsg.Receipient]
		if ok && receiver != nil {
			receiver.Outgoing <- outgoingEvent
		} else {
			for user := range c.manager.userIndex {
				if c.manager.userIndex[user] == nil {
					msgJson2, _ := json.Marshal(outgoingEvent)
					err := c.manager.redisClient.RPush(ctx, "message-history"+c.Username, msgJson2).Err()
					if err != nil {
						log.Println("there was a problem with the offline redis")
					}
				}
			}
		}

		c.Outgoing <- outgoingEvent
	}

	//for offline message queues

	return nil

}

// func to create the Eventhandler for the joinmessage event
func JoinMessage(event Event, c *Client) error {
	//redclient := NewRedisService()
	var joinmsg JoinMessageEvent
	if err := json.Unmarshal(event.Payload, &joinmsg); err != nil {
		log.Print("error in unmarshalling payload")
	}

	c.Username = joinmsg.Username
	fmt.Printf("%v has joined!!", c.Username)
	//wheneber a user sends a joinmessage, it the previous messages should be sent first if offline
	c.manager.userIndex[c.Username] = c
	history := c.manager.redisClient.LRange(ctx, "message-history"+c.Username, 0, -1).Val()
	c.manager.redisClient.Del(ctx, "message-history"+c.Username) //dont forget to clear after delivery
	for _, hist := range history {
		var ev Event
		err := json.Unmarshal([]byte(hist), &ev)
		if err != nil {
			log.Println(err)
		}
		c.Outgoing <- ev
	}
	goingOut := Event{
		Type:    EventtoJoinMEssage,
		Payload: event.Payload,
	}
	for client := range c.manager.Clients {
		client.Outgoing <- goingOut
	}
	return nil
}

// the function to route the events to their respective handlers based on the type of event e.g a 'send-message'and actually execute them
func (m *Manager) RouteEvents(event Event, c *Client) error {
	//check if the event even exists in the first place
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("no such event type exists")
	}
}

var (
	WsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func (m *Manager) wsHandler(w http.ResponseWriter, r *http.Request) {
	///upgrade the regular http to websockets
	wsConn, err := WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("there was an error in upgrading the request: %v", err)
		return
	}

	client := NewClientService(wsConn, m)

	m.AddClient(client)

	go client.ReadMessages()

	go client.WriteMessage()
	//defer wsConn.Close()

}
func (m *Manager) AddClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	// we add the client to the list of clients
	m.Clients[client] = true
}
func (m *Manager) DeleteClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.Clients[client]; ok {
		client.conn.Close()
		close(client.Outgoing)
		delete(m.Clients, client)
	}

}

// this broadcasts the message to every client in the client map
/*func (m *Manager) Broadcast(msg Event) {
	m.RLock()
	defer m.RUnlock()
	//for client := range m.Clients {
	///	client.Outgoing <- msg
	//}
	for client := range m.Clients {
		select {
		case client.Outgoing <- msg:
			//log.Print("message enqued successfully")
		default:
			log.Println("dropping message for client")
		}
	}

}*/
