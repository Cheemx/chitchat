package model

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type UserBase map[*User]bool

type User struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username   string             `json:"username,omitempty"`
	Password   string             `json:"password,omitempty"`
	connection *websocket.Conn
	manager    *Manager
	egress     chan Event
	// egress is used to avoid concurrent writes on the websocket connection
}

type Event struct {
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type Manager struct {
	users UserBase
	sync.RWMutex
	// create map handler for events
}

func NewManager() *Manager {
	m := &Manager{
		users: make(UserBase),
	}
	return m
}

func NewUser(conn *websocket.Conn, m *Manager) *User {
	return &User{
		connection: conn,
		manager:    m,
		egress:     make(chan Event),
	}
}

func (m *Manager) AddUser(c *User) {
	m.Lock()
	defer m.Unlock()

	m.users[c] = true
	log.Println("user Added:", c)

	// Notify all users about the new users
	event := Event{
		Sender:    "System",
		Content:   c.Username + " joined the chat",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	m.broadcast(event, nil)
}

func (m *Manager) RemoveUser(c *User) {
	m.Lock()
	defer m.Unlock()

	_, ok := m.users[c]
	if ok {
		c.connection.Close()
		delete(m.users, c)
		log.Println("user removed: ", c.Username)

		//Notify all users about the user leaving
		event := Event{
			Sender:    "System",
			Content:   c.Username + " left the chat",
			Timestamp: time.Now().Format(time.RFC3339),
		}
		m.broadcast(event, nil)
	}
}

func (m *Manager) broadcast(event Event, ignore *User) {
	for user := range m.users {
		if user != ignore {
			user.egress <- event
		}
	}
}

func (c *User) ReadMessages() {
	defer func() {
		c.manager.RemoveUser(c)
	}()

	for {
		// msgTypes : ping, pong, data, control, binary etc.
		var event Event
		err := c.connection.ReadJSON(&event)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}
		event.Sender = c.Username
		event.Timestamp = time.Now().Format(time.RFC3339)
		c.manager.broadcast(event, c)
	}
}

func (c *User) WriteMessages() {
	defer func() {
		c.manager.RemoveUser(c)
	}()

	for {
		select {
		case event, ok := <-c.egress:
			if !ok {
				e := c.connection.WriteMessage(websocket.CloseMessage, nil)
				if e != nil {
					log.Println("Connection Closed: ", e)
				}
				return
			}
			data, err := json.Marshal(event)

			if err != nil {
				log.Println(err)
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
				return
			}
			log.Println("message sent")
		}
	}
}
