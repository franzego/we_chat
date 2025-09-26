package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Manager struct {
	Clients ClientList
	sync.RWMutex
	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		Clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.SetupEventHandlers()
	return m

}

// the function to setup the event handlers
func (m *Manager) SetupEventHandlers() {
	m.handlers[EventtoSendMessage] = SendMessage
}

// the function to create the Eventhandler for the SendMessage Event
func SendMessage(event Event, c *Client) error {
	fmt.Println(event)
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
