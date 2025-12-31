package domain

import "time"

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
	MessageTypeAudio MessageType = "audio"
	MessageTypeVideo MessageType = "video"
)

type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
)

type Message struct {
	ID             string
	ConversationID string
	SenderID       string
	Content        string
	MessageType    MessageType
	SequenceID     int64
	ReplyToID      *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type MessageDelivery struct {
	MessageID string
	UserID    string
	Status    MessageStatus
	Timestamp time.Time
}

type MessageRepository interface {
	Create(message *Message) error
	GetByID(id string) (*Message, error)
	GetByConversation(conversationID string, limit, offset int) ([]*Message, error)
	UpdateDeliveryStatus(messageID, userID string, status MessageStatus) error
	GetDeliveryStatus(messageID string) ([]*MessageDelivery, error)
}
