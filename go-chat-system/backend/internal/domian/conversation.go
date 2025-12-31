package domain

import "time"

type ConversationType string

const (
	ConversationTypeOneToOne ConversationType = "one_to_one"
	ConversationTypeGroup    ConversationType = "group"
)

type Conversation struct {
	ID        string
	Type      ConversationType
	Name      *string
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ConversationMember struct {
	ConversationID string
	UserID         string
	Role           string
	JoinedAt       time.Time
	LastReadAt     *time.Time
}

type ConversationRepository interface {
	Create(conversation *Conversation) error
	GetByID(id string) (*Conversation, error)
	GetByUsers(userIDs []string) (*Conversation, error)
	AddMember(conversationID, userID, role string) error
	RemoveMember(conversationID, userID string) error
	GetMembers(conversationID string) ([]*ConversationMember, error)
	IsMember(conversationID, userID string) (bool, error)
	UpdateLastRead(conversationID, userID string) error
}
