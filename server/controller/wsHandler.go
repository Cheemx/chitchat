package controller

import (
	"api/middleware"
	"api/model"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var m = model.NewManager()

func ServeWS(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authorized or not
	user, ok := r.Context().Value(middleware.UserContextKey).(model.User)

	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}
	username := user.Username

	log.Println("New Connection")

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	c := model.NewUser(conn, m)
	c.Username = username
	m.AddUser(c)

	go c.ReadMessages()
	go c.WriteMessages()
}
