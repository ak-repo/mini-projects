package app

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"time"

	chat "github.com/ak-repo/chat-application/gen/chatpb"

	grpcserver "github.com/ak-repo/chat-application/internal/adapter/grpc"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Repo defines the persistence contract for the chat service.
type Repo interface {
	SaveMessage(ctx context.Context, m *chat.ServerMessage) error
	ListHistory(ctx context.Context, chatID string, limit, offset int) ([]*chat.ServerMessage, error)
}

// Server implements the chatpb.ChatServiceServer interface.
type Server struct {
	chat.UnimplementedChatServiceServer
	hub  *grpcserver.Hub
	repo Repo
	rdb  *redis.Client
}

// NewServer creates a new gRPC chat server instance.
func NewServer(repo Repo, rdb *redis.Client) *Server {
	return &Server{
		hub:  grpcserver.NewHub(),
		repo: repo,
		rdb:  rdb,
	}
}

// StartRedisSubscriber listens on Redis for messages published by other chat instances.
// This is the core of horizontal scaling.
func (s *Server) StartRedisSubscriber(ctx context.Context) {
	log.Println("Starting Redis subscriber...")
	pubsub := s.rdb.Subscribe(ctx, "chat:messages")
	ch := pubsub.Channel()

	for msg := range ch {
		var sm chat.ServerMessage
		if err := json.Unmarshal([]byte(msg.Payload), &sm); err != nil {
			log.Println("Redis unmarshal error:", err)
			continue
		}

		// Attempt to deliver the message to the recipient if they are connected to THIS instance.
		resp := &chat.StreamResponse{
			Payload: &chat.StreamResponse_ServerMessage{ServerMessage: &sm},
		}
		if err := s.hub.SendTo(sm.To, resp); err != nil {
			// If ErrNoClient, the recipient is either connected to another instance or offline.
			// The message has already been persisted and published, so we ignore the error here.
			if !errors.Is(err, grpcserver.ErrNoClient) {
				log.Printf("Error sending message locally to %s: %v", sm.To, err)
			}
		}
	}
}

// verifyToken is a placeholder for actual JWT validation using pkg/jwt.
func (s *Server) verifyToken(token string) (string, error) {
	// TODO: Integrate with your pkg/jwt to validate token and extract UserID/Email.
	// For this runnable example, we just ensure the token is not empty.
	if token == "" {
		return "", errors.New("missing or invalid token")
	}
	// Assuming the valid token payload contains the user's ID/Email (e.g., "user@example.com")
	return token, nil
}

// Stream handles the core bidirectional chat stream.
func (s *Server) Stream(stream chat.ChatService_StreamServer) error {
	// 1. Authentication (Expected first message)
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	auth := req.GetAuth()
	if auth == nil {
		return errors.New("first message must be Auth payload")
	}
	user, err := s.verifyToken(auth.Token)
	if err != nil {
		return errors.New("auth failed: " + err.Error())
	}

	// 2. Setup Stream Wrapper and Registration
	sw := &grpcserver.StreamWrapper{
		User:     user,
		Stream:   stream,
		SendCh:   make(chan *chat.StreamResponse, 64), // Backpressure buffer
		Shutdown: make(chan struct{}),
	}
	s.hub.Register(user, sw)
	defer s.hub.Unregister(user) // Ensure cleanup

	// 3. Writer Goroutine (Service -> Gateway -> WS Client)
	go func() {
		for {
			select {
			case resp := <-sw.SendCh:
				if err := stream.Send(resp); err != nil {
					log.Printf("gRPC stream send error to %s: %v", user, err)
					return // End writer on send failure
				}
			case <-sw.Shutdown:
				return // Graceful shutdown
			}
		}
	}()

	// 4. Reader Loop (WS Client -> Gateway -> Service)
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil // Client closed the stream gracefully
		}
		if err != nil {
			log.Printf("gRPC stream receive error from %s: %v", user, err)
			return err
		}

		if cm := in.GetClientMessage(); cm != nil {
			// 4a. Process ClientMessage
			sm := &chat.ServerMessage{
				ServerId:  uuid.NewString(), // Generate server-side unique ID
				ClientId:  cm.ClientId,
				From:      user, // Use verified user ID as 'From' (prevents spoofing)
				To:        cm.To,
				ChatId:    cm.ChatId,
				Text:      cm.Text,
				FileId:    cm.FileId,
				MsgType:   cm.MsgType,
				CreatedAt: time.Now().UnixMilli(), // Use server time
			}

			// 4b. Persist message
			if err := s.repo.SaveMessage(stream.Context(), sm); err != nil {
				log.Println("Database save message error:", err)
			}

			// 4c. Broadcast to all instances via Redis
			b, _ := json.Marshal(sm)
			if err := s.rdb.Publish(stream.Context(), "chat:messages", b).Err(); err != nil {
				log.Println("Redis publish error:", err)
			}

			// 4d. Send acknowledgement back to sender (locally)
			resp := &chat.StreamResponse{
				Payload: &chat.StreamResponse_ServerMessage{ServerMessage: sm},
			}
			_ = s.hub.SendTo(user, resp) // Best effort send
		}
	}
}

// ListHistory handles the unary RPC for message history retrieval.
func (s *Server) ListHistory(ctx context.Context, req *chat.ListHistoryRequest) (*chat.ListHistoryResponse, error) {
	rows, err := s.repo.ListHistory(ctx, req.ChatId, int(req.Limit), int(req.Offset))
	if err != nil {
		log.Println("ListHistory DB error:", err)
		return nil, errors.New("failed to retrieve history")
	}
	return &chat.ListHistoryResponse{Messages: rows}, nil
}
