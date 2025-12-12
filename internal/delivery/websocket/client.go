package websocket

import (
	"collab-editor/internal/domain"
	"collab-editor/internal/usecase"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var jwtSecret = []byte("my-secret-key")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	Hub     *Hub
	Service *usecase.EditorService
	Conn    *websocket.Conn
	RoomID  string
	UserID  string
	send    chan interface{}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		var op domain.Operation
		err := c.Conn.ReadJSON(&op)
		if err != nil {
			break
		}

		op.ClientID = c.UserID

		err = c.Service.ProcessEdit(c.RoomID, op)
		if err != nil {
			log.Printf("Error processing edit: %v", err)
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for {
		msg, ok := <-c.send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		c.Conn.WriteJSON(msg)
	}
}

func ServeWs(hub *Hub, service *usecase.EditorService, w http.ResponseWriter, r *http.Request) {

	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	userID := claims["user_id"].(string)

	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		http.Error(w, "Room ID required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		Hub:     hub,
		Service: service,
		Conn:    conn,
		RoomID:  roomID,
		UserID:  userID,
		send:    make(chan interface{}, 256),
	}
	client.Hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}
