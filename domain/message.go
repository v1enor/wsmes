package domain

import "time"

type Message struct {
	Username string    `json:"Username"`
	Content  string    `json:"Content"`
	Time     time.Time `json:"Time"`
}
