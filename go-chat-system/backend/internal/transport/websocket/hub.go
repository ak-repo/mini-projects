package websocket

import (
	"sync"

	"github.com/ak-repo/go-chat-system/internal/service"
)

type Hub struct {
	clients      map[string]*Client
	broadcast    chan *Message
	register     chan *Client
	unregister   chan *Client
	mu           sync.RWMutex
	chatService  *service.ChatService
	presenceServ *service.PresenceService
}

func NewHub(chatService *service.ChatService, presenceServ *service.PresenceService) *Hub {
	return &Hub{
		clients:      make(map[string]*Client),
		broadcast:    make(chan *Message, 256),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		chatService:  chatService,
		presenceServ: presenceServ,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
			h.presenceServ.SetOnline(client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.mu.Unlock()
			h.presenceServ.SetOffline(client.userID)

		case message := <-h.broadcast:
			h.BroadcastMessage(message)
		}
	}
}

func (h *Hub) BroadcastMessage(msg *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, recipientID := range msg.Recipients {
		if client, ok := h.clients[recipientID]; ok {
			select {
			case client.send <- msg:
			default:
				// Client buffer full, skip
			}
		}
	}
}

func (h *Hub) GetOnlineUsers(userIDs []string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	online := []string{}
	for _, userID := range userIDs {
		if _, ok := h.clients[userID]; ok {
			online = append(online, userID)
		}
	}
	return online
}

// Add to hub.go for gRPC integration
func (h *Hub) IsClientConnected(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}
