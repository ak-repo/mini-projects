package main

import (
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

var (
	users  = make(map[string]*websocket.Conn)
	groups = make(map[string]map[string]bool)
	group1 = "group1"

	mu sync.Mutex
)

type Message struct {
	GroupID string `json:"group_id"`
	Message string `json:"message"`
}

func chatHandler(ws *websocket.Conn) {
	// 1. Extract user_id
	userID := ws.Request().URL.Query().Get("user_id")
	if userID == "" {
		ws.Close()
		return
	}

	// 2. Register user
	mu.Lock()
	users[userID] = ws

	// For demo: auto-join a group called "group1"
	if groups[group1] == nil {
		groups[group1] = make(map[string]bool)
	}
	groups[group1][userID] = true
	mu.Unlock()

	log.Println("User connected:", userID)

	// 3. Cleanup on disconnect
	defer func() {
		mu.Lock()
		delete(users, userID)
		delete(groups[group1], userID)
		mu.Unlock()
		ws.Close()
		log.Println("User disconnected:", userID)
	}()

	// 4. Read loop
	for {
		var msg Message

		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			log.Println("Receive error:", err)
			return
		}

		// 5. Send message to all group members
		mu.Lock()
		members, ok := groups[msg.GroupID]
		mu.Unlock()

		if !ok {
			log.Println("Group not found:", msg.GroupID)
			continue
		}

		for memberID := range members {
			// Do not send message back to sender
			if memberID == userID {
				continue
			}

			mu.Lock()
			conn, ok := users[memberID]
			mu.Unlock()

			if !ok {
				continue
			}

			// Send message
			err := websocket.JSON.Send(conn, map[string]string{
				"group":   msg.GroupID,
				"from":    userID,
				"message": msg.Message,
			})
			if err != nil {
				log.Println("Send error:", err)
			}
		}
	}
}

func main() {
	http.Handle("/ws", websocket.Handler(chatHandler))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
