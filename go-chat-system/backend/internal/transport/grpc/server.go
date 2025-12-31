package grpc

import (
	"context"

	"github.com/ak-repo/go-chat-system/internal/service"
	"github.com/ak-repo/go-chat-system/internal/transport/websocket"
	"github.com/ak-repo/go-chat-system/pkg/pb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedChatServiceServer
	chatService *service.ChatService
	hub         *websocket.Hub
}

func NewServer(chatService *service.ChatService, hub *websocket.Hub) *grpc.Server {
	grpcServer := grpc.NewServer()

	srv := &Server{
		chatService: chatService,
		hub:         hub,
	}

	pb.RegisterChatServiceServer(grpcServer, srv)

	return grpcServer
}

func (s *Server) BroadcastMessage(ctx context.Context, req *pb.BroadcastRequest) (*pb.BroadcastResponse, error) {
	// This would be called by other nodes in a multi-node setup
	msg := &websocket.Message{
		Type: "new_message",
		Payload: map[string]interface{}{
			"id":              req.MessageId,
			"conversation_id": req.ConversationId,
			"sender_id":       req.SenderId,
			"content":         req.Content,
			"message_type":    req.MessageType,
			"created_at":      req.CreatedAt,
		},
		Recipients: req.RecipientIds,
	}

	s.hub.BroadcastMessage(msg)

	// Count how many recipients are online on this node
	onlineCount := 0
	for _, recipientID := range req.RecipientIds {
		if s.hub.IsClientConnected(recipientID) {
			onlineCount++
		}
	}

	return &pb.BroadcastResponse{
		Success:        true,
		DeliveredCount: int32(onlineCount),
	}, nil
}

func (s *Server) GetOnlineUsers(ctx context.Context, req *pb.OnlineUsersRequest) (*pb.OnlineUsersResponse, error) {
	onlineUsers := s.hub.GetOnlineUsers(req.UserIds)

	return &pb.OnlineUsersResponse{
		OnlineUserIds: onlineUsers,
	}, nil
}

func (s *Server) NotifyPresence(ctx context.Context, req *pb.PresenceRequest) (*pb.PresenceResponse, error) {
	// Broadcast presence update to connected clients
	msg := &websocket.Message{
		Type: "presence_update",
		Payload: map[string]interface{}{
			"user_id":   req.UserId,
			"status":    req.Status,
			"timestamp": req.Timestamp,
		},
		Recipients: []string{}, // Broadcast to all would need implementation
	}

	s.hub.BroadcastMessage(msg)

	return &pb.PresenceResponse{
		Success: true,
	}, nil
}
