package main

import (
	"collab-editor/internal/delivery/websocket"
	"collab-editor/internal/repository"
	"collab-editor/internal/usecase"
	"log"
	"net/http"
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
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
