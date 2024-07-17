package domain

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	Name      string                   `json:"name"`
	Clients   map[*websocket.Conn]bool `json:"-"`
	Broadcast chan Message             `json:"-"`
	Messages  []Message                `json:"messages"`
	Lock      sync.Mutex               `json:"-"`
}

func NewRoom(name string) *Room {
	return &Room{
		Name:      name,
		Messages:  []Message{},
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan Message),
	}
}
