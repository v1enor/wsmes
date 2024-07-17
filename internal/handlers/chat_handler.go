package handlers

import (
	"log"
	"net/http"
	"time"
	"wsmes/domain"
	"wsmes/internal/database"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ChatHandler(db *database.Database, msgTime int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		EnableCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to websocket:", err)
			return
		}
		defer ws.Close()

		roomName := r.URL.Query().Get("room")
		if roomName == "" {
			http.Error(w, "Room name is required", http.StatusBadRequest)
			return
		}

		room := db.GetRoom(roomName)
		if room == nil {
			room = db.CreateRoom(roomName)
			log.Println("Created new room:", roomName)
		}

		room.Lock.Lock()
		if room.Clients == nil {
			room.Clients = make(map[*websocket.Conn]bool)
		}
		room.Clients[ws] = true
		room.Lock.Unlock()
		log.Printf("Client %v added to room %s", ws.RemoteAddr(), roomName)

		lastMessages, err := db.GetLastMessages(roomName, msgTime)
		if err != nil {
			log.Println("Error getting last messages:", err)
		} else {
			for _, msg := range lastMessages {
				log.Println("Sending message to client:", msg)
				if err := ws.WriteJSON(msg); err != nil {
					log.Println("Error writing message to client:", err)
					break
				}
			}
		}

		go func() {
			for message := range room.Broadcast {
				room.Lock.Lock()
				for client := range room.Clients {
					// Убрать условие client != ws, чтобы отправлять сообщения обратно отправителю
					if err := client.WriteJSON(message); err != nil {
						log.Println("Error writing message to client:", err)
						client.Close()
						delete(room.Clients, client)
					}
				}
				room.Lock.Unlock()
			}
		}()

		for {
			var msg domain.Message
			if err := ws.ReadJSON(&msg); err != nil {
				log.Println("Error reading message from client:", err)
				break
			}
			msg.Time = time.Now()
			log.Println("Received message:", msg)

			room.Broadcast <- msg
			room.Messages = append(room.Messages, msg)
			if err := db.SaveMessage(roomName, msg); err != nil {
				log.Println("Error saving message:", err)
			}
		}

		room.Lock.Lock()
		delete(room.Clients, ws)
		room.Lock.Unlock()
		log.Printf("Client %v removed from room %s", ws.RemoteAddr(), roomName)
	}
}
