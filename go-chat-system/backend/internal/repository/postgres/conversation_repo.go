package postgres

import (
	"database/sql"

	domain "github.com/ak-repo/go-chat-system/internal/domian"
	"github.com/jmoiron/sqlx"
)

type ConversationRepository struct {
	db *sqlx.DB
}

func NewConversationRepository(db *sqlx.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) Create(conversation *domain.Conversation) error {
	query := `
		INSERT INTO conversations (id, type, name, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, conversation.ID, conversation.Type, conversation.Name,
		conversation.CreatedBy, conversation.CreatedAt, conversation.UpdatedAt)
	return err
}

func (r *ConversationRepository) GetByID(id string) (*domain.Conversation, error) {
	var conv domain.Conversation
	query := `SELECT * FROM conversations WHERE id = $1`
	err := r.db.Get(&conv, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &conv, err
}

func (r *ConversationRepository) GetByUsers(userIDs []string) (*domain.Conversation, error) {
	query := `
		SELECT c.* FROM conversations c
		WHERE c.type = 'one_to_one' AND c.id IN (
			SELECT conversation_id FROM conversation_members
			WHERE user_id = ANY($1)
			GROUP BY conversation_id
			HAVING COUNT(DISTINCT user_id) = $2
		)
		LIMIT 1
	`
	var conv domain.Conversation
	err := r.db.Get(&conv, query, userIDs, len(userIDs))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &conv, err
}

func (r *ConversationRepository) AddMember(conversationID, userID, role string) error {
	query := `
		INSERT INTO conversation_members (conversation_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT DO NOTHING
	`
	_, err := r.db.Exec(query, conversationID, userID, role)
	return err
}

func (r *ConversationRepository) RemoveMember(conversationID, userID string) error {
	query := `DELETE FROM conversation_members WHERE conversation_id = $1 AND user_id = $2`
	_, err := r.db.Exec(query, conversationID, userID)
	return err
}

func (r *ConversationRepository) GetMembers(conversationID string) ([]*domain.ConversationMember, error) {
	var members []*domain.ConversationMember
	query := `SELECT * FROM conversation_members WHERE conversation_id = $1`
	err := r.db.Select(&members, query, conversationID)
	return members, err
}

func (r *ConversationRepository) IsMember(conversationID, userID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM conversation_members WHERE conversation_id = $1 AND user_id = $2`
	err := r.db.Get(&count, query, conversationID, userID)
	return count > 0, err
}

func (r *ConversationRepository) UpdateLastRead(conversationID, userID string) error {
	query := `
		UPDATE conversation_members 
		SET last_read_at = NOW() 
		WHERE conversation_id = $1 AND user_id = $2
	`
	_, err := r.db.Exec(query, conversationID, userID)
	return err
}
