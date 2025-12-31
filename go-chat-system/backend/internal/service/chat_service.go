package service

import (
	"time"

	domain "github.com/ak-repo/go-chat-system/internal/domian"
)

type ChatService struct {
	messageRepo      domain.MessageRepository
	conversationRepo domain.ConversationRepository
	userRepo         domain.UserRepository
}

func NewChatService(
	messageRepo domain.MessageRepository,
	conversationRepo domain.ConversationRepository,
	userRepo domain.UserRepository,
) *ChatService {
	return &ChatService{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
		userRepo:         userRepo,
	}
}

func (s *ChatService) SendMessage(senderID, conversationID, content string, msgType domain.MessageType) (*domain.Message, error) {
	// Verify sender is member
	isMember, err := s.conversationRepo.IsMember(conversationID, senderID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, domain.ErrPermissionDenied
	}

	message := &domain.Message{
		ID:             generateID(),
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		MessageType:    msgType,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *ChatService) GetConversationMessages(userID, conversationID string, limit, offset int) ([]*domain.Message, error) {
	// Verify user is member
	isMember, err := s.conversationRepo.IsMember(conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, domain.ErrPermissionDenied
	}

	return s.messageRepo.GetByConversation(conversationID, limit, offset)
}

func (s *ChatService) MarkAsDelivered(messageID, userID string) error {
	return s.messageRepo.UpdateDeliveryStatus(messageID, userID, domain.MessageStatusDelivered)
}

func (s *ChatService) MarkAsRead(messageID, userID string) error {
	return s.messageRepo.UpdateDeliveryStatus(messageID, userID, domain.MessageStatusRead)
}

func (s *ChatService) CreateConversation(creatorID string, memberIDs []string, convType domain.ConversationType) (*domain.Conversation, error) {
	// Check if 1-to-1 conversation already exists
	if convType == domain.ConversationTypeOneToOne {
		existing, _ := s.conversationRepo.GetByUsers(append(memberIDs, creatorID))
		if existing != nil {
			return existing, nil
		}
	}

	conv := &domain.Conversation{
		ID:        generateID(),
		Type:      convType,
		CreatedBy: creatorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.conversationRepo.Create(conv); err != nil {
		return nil, err
	}

	// Add creator
	if err := s.conversationRepo.AddMember(conv.ID, creatorID, "admin"); err != nil {
		return nil, err
	}

	// Add members
	for _, memberID := range memberIDs {
		if err := s.conversationRepo.AddMember(conv.ID, memberID, "member"); err != nil {
			return nil, err
		}
	}

	return conv, nil
}

// Helper method for chat service
func (s *ChatService) GetConversationMembers(conversationID string) ([]*domain.ConversationMember, error) {
	return s.conversationRepo.GetMembers(conversationID)
}
