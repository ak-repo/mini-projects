package postgres

import (
	"database/sql"

	domain "github.com/ak-repo/go-chat-system/internal/domian"
	"github.com/jmoiron/sqlx"
)

type MessageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *domain.Message) error {
	query := `
		INSERT INTO messages (id, conversation_id, sender_id, content, message_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, message.ID, message.ConversationID, message.SenderID,
		message.Content, message.MessageType, message.CreatedAt, message.UpdatedAt)
	return err
}

func (r *MessageRepository) GetByID(id string) (*domain.Message, error) {
	var message domain.Message
	query := `SELECT * FROM messages WHERE id = $1`
	err := r.db.Get(&message, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &message, err
}

func (r *MessageRepository) GetByConversation(conversationID string, limit, offset int) ([]*domain.Message, error) {
	var messages []*domain.Message
	query := `
		SELECT * FROM messages 
		WHERE conversation_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`
	err := r.db.Select(&messages, query, conversationID, limit, offset)
	return messages, err
}

func (r *MessageRepository) UpdateDeliveryStatus(messageID, userID string, status domain.MessageStatus) error {
	query := `
		INSERT INTO message_delivery (message_id, user_id, status, timestamp)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (message_id, user_id) 
		DO UPDATE SET status = $3, timestamp = NOW()
	`
	_, err := r.db.Exec(query, messageID, userID, status)
	return err
}

func (r *MessageRepository) GetDeliveryStatus(messageID string) ([]*domain.MessageDelivery, error) {
	var deliveries []*domain.MessageDelivery
	query := `SELECT * FROM message_delivery WHERE message_id = $1`
	err := r.db.Select(&deliveries, query, messageID)
	return deliveries, err
}
