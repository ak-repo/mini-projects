package postgres

import (
	"context"
	"time"

	"github.com/ak-repo/chat-application/gen/chatpb"
	"github.com/ak-repo/chat-application/internal/app"
	"github.com/jackc/pgx/v5/pgxpool"
)

type chatRepo struct {
	pool *pgxpool.Pool
}

func NewChatRepo(pool *pgxpool.Pool) app.Repo {
	return &chatRepo{pool: pool}
}

// SaveMessage persists a ServerMessage to the database.
func (r *chatRepo) SaveMessage(ctx context.Context, m *chatpb.ServerMessage) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO messages (id, client_id, chat_id, sender, receiver, text, file_id, msg_type, created_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		m.ServerId, m.ClientId, m.ChatId, m.From, m.To, m.Text, m.FileId, m.MsgType, time.UnixMilli(m.CreatedAt).UTC(),
	)
	return err
}

// ListHistory retrieves paginated message history for a given chat ID.
// ListHistory retrieves paginated message history for a given chat ID.
func (r *chatRepo) ListHistory(ctx context.Context, chatID string, limit, offset int) ([]*chatpb.ServerMessage, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT 
			id, client_id, sender, receiver, text, file_id, msg_type, 
			EXTRACT(EPOCH FROM created_at)*1000 AS created_at -- convert to epoch milliseconds
		 FROM messages 
		 WHERE chat_id=$1 
		 ORDER BY created_at DESC 
		 LIMIT $2 OFFSET $3`,
		chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*chatpb.ServerMessage
	for rows.Next() {
		var m chatpb.ServerMessage
		var created float64 // Float to capture the result of EXTRACT(EPOCH)

		// Note: receiver can be NULL in the DB, so we scan into a pointer or string
		var receiver *string

		if err := rows.Scan(&m.ServerId, &m.ClientId, &m.From, &receiver, &m.Text, &m.FileId, &m.MsgType, &created); err != nil {
			return nil, err
		}

		// Handle optional receiver
		if receiver != nil {
			m.To = *receiver
		}

		m.CreatedAt = int64(created)
		out = append(out, &m)
	}
	return out, nil
}
