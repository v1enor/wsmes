package main

import (
	"log"
	"net/http"
	"wsmes/config"
	"wsmes/internal/database"
	"wsmes/internal/handlers"
)

func main() {
	db := database.NewDatabase()
	cnfg := config.LoadConfig()
	http.HandleFunc("/ws", handlers.ChatHandler(db, cnfg.MsgTime))
	http.HandleFunc("/rooms", handlers.RoomsHandler(db))

	log.Println("Starting server on : " + cnfg.Port)
	err := http.ListenAndServe(":"+cnfg.Port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
