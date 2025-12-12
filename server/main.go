package main

import (
	"collab-editor/internal/delivery/websocket"
	"collab-editor/internal/repository"
	"collab-editor/internal/usecase"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	repo := repository.NewInMemoryRepo()
	hub := websocket.NewHub()

	go hub.Run()

	editorService := usecase.NewEditorService(repo, hub)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, editorService, w, r)
	})

	log.Println("Server started on :8000")

	var jwtSecret = []byte("my-secret-key")
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		userID := "User-" + r.URL.Query().Get("user")
		if userID == "User-" {
			userID = "Anonymous"
		}

		claims := jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			http.Error(w, "Error signing token", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(tokenString))
	})
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
