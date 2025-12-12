package main

import (
	"collab-editor/internal/delivery/websocket"
	"collab-editor/internal/repository"
	"collab-editor/internal/usecase"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	log.Println("Connecting to Redis...")
	repo, err := repository.NewRedisRepo(redisURL)
	if err != nil {
		log.Fatal("Could not connect to Redis:", err)
	}
	log.Println("Connected to Redis successfully!")

	hub := websocket.NewHub()
	go hub.Run()

	editorService := usecase.NewEditorService(repo, hub)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, editorService, w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

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

	http.HandleFunc("/document", func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("token")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		roomID := r.URL.Query().Get("room")
		if roomID == "" {
			http.Error(w, "Room ID required", http.StatusBadRequest)
			return
		}

		doc, err := editorService.GetDocument(roomID)
		if err != nil {
			http.Error(w, "Failed to get doc", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(doc)
	})

	log.Println("Server started on :8000")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
