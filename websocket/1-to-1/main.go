package main

import (
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

var (
	users = make(map[string]*websocket.Conn)
	mu    sync.Mutex
)

type Message struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

/*
WebSocket handler
*/
func chatHandler(ws *websocket.Conn) {
	// 1. Get userID from query param
	userID := ws.Request().URL.Query().Get("user_id")
	if userID == "" {
		log.Println("user_id missing")
		ws.Close()
		return
	}

	// 2. Register user
	mu.Lock()
	users[userID] = ws
	mu.Unlock()

	log.Println("User connected:", userID)

	// 3. Cleanup on disconnect
	defer func() {
		mu.Lock()
		delete(users, userID)
		mu.Unlock()
		ws.Close()
		log.Println("User disconnected:", userID)
	}()

	// 4. Read loop
	for {
		var msg Message

		// Receive message from this user
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			log.Println("Receive error:", err)
			return
		}

		// 5. Route message to recipient
		mu.Lock()
		receiverConn, ok := users[msg.To]
		mu.Unlock()

		if !ok {
			log.Println("User not connected:", msg.To)
			continue
		}

		// Send message to recipient
		err = websocket.JSON.Send(receiverConn, map[string]string{
			"from":    userID,
			"message": msg.Message,
		})
		if err != nil {
			log.Println("Send error:", err)
		}
	}
}

func main() {
	http.Handle("/ws", websocket.Handler(chatHandler))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
