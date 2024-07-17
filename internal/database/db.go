package database

import (
	"log"
	"sync"
	"time"
	"wsmes/domain"
)

type Database struct {
	messages map[string][]domain.Message
	rooms    map[string]*domain.Room
	lock     sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		messages: make(map[string][]domain.Message),
		rooms:    make(map[string]*domain.Room),
	}
}

func (db *Database) SaveMessage(roomName string, msg domain.Message) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.messages[roomName] = append(db.messages[roomName], msg)
	log.Println("Saving message to database:", msg)
	return nil
}

func (db *Database) GetLastMessages(roomName string, msgTime int) ([]domain.Message, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	allMessages, ok := db.messages[roomName]
	if !ok {
		return nil, nil
	}

	var filteredMessages []domain.Message
	for _, msg := range allMessages {
		if msg.Time.After(time.Now().Add(-time.Duration(msgTime) * time.Minute)) {
			filteredMessages = append(filteredMessages, msg)
		}
	}

	return filteredMessages, nil
}

func (db *Database) GetRoom(roomName string) *domain.Room {
	db.lock.RLock()
	defer db.lock.RUnlock()

	room, ok := db.rooms[roomName]
	if !ok {
		return nil
	}
	return room
}

func (db *Database) CreateRoom(roomName string) *domain.Room {
	db.lock.Lock()
	defer db.lock.Unlock()

	room, exists := db.rooms[roomName]
	if exists {
		return room
	}

	room = domain.NewRoom(roomName)
	db.rooms[roomName] = room
	log.Println("Created new room:", roomName)
	return room
}

func (db *Database) SaveRoom(room *domain.Room) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.rooms[room.Name] = room
	return nil
}

func (db *Database) GetAllRooms() []string {
	db.lock.RLock()
	defer db.lock.RUnlock()

	var rooms []string
	for name := range db.rooms {
		rooms = append(rooms, name)
	}
	return rooms
}
