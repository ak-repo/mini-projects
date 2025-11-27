package grpcserver

import (
	"errors"
	"log"
	"sync"

	"github.com/ak-repo/chat-application/gen/chatpb"
)

// ErrNoClient is returned when a recipient is not connected to this instance.
var ErrNoClient = errors.New("client not connected locally")

type StreamWrapper struct {
	User   string
	Stream chatpb.ChatService_StreamServer
	// SendCh is a buffered channel to handle backpressure.
	SendCh chan *chatpb.StreamResponse
	// Shutdown signals the writer goroutine to exit.
	Shutdown chan struct{}
}

// Hub manages the active gRPC streams (connected clients) in memory.
type Hub struct {
	mu sync.RWMutex
	// clients maps UserID/Email (string) to their active StreamWrapper
	clients map[string]*StreamWrapper
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{clients: make(map[string]*StreamWrapper)}
}

// Register adds a new active stream to the hub.
func (h *Hub) Register(user string, w *StreamWrapper) {
	h.mu.Lock()
	defer h.mu.Unlock()
	// Close any existing stream for this user first
	if old, ok := h.clients[user]; ok {
		log.Printf("Closing stale stream for user: %s", user)
		close(old.Shutdown)
	}
	h.clients[user] = w
	log.Printf("User registered: %s. Total active streams: %d", user, len(h.clients))
}

// Unregister removes an active stream from the hub.
func (h *Hub) Unregister(user string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if w, ok := h.clients[user]; ok {
		close(w.Shutdown)
		delete(h.clients, user)
		log.Printf("User unregistered: %s. Total active streams: %d", user, len(h.clients))
	}
}

// Get retrieves a stream wrapper for a specific user.
func (h *Hub) Get(user string) (*StreamWrapper, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	w, ok := h.clients[user]
	return w, ok
}

// SendTo attempts to deliver a message directly to a locally connected client.
func (h *Hub) SendTo(user string, resp *chatpb.StreamResponse) error {
	w, ok := h.Get(user)
	if !ok {
		return ErrNoClient
	}
	select {
	case w.SendCh <- resp:
		return nil
	default:
		// Drop message if channel is full (simple backpressure)
		return errors.New("send channel blocked/full")
	}
}
