package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"wsmes/internal/database"
)

func RoomsHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		EnableCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		rooms := db.GetAllRooms()
		roomsJSON, err := json.Marshal(rooms)
		if err != nil {
			log.Println("Error marshalling rooms:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(roomsJSON)
	}
}
