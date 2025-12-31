package postgres

import (
	"database/sql"

	domain "github.com/ak-repo/go-chat-system/internal/domian"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, display_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.PasswordHash,
		user.DisplayName, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &user, err
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.Get(&user, query, email)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &user, err
}

func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE username = $1`
	err := r.db.Get(&user, query, username)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &user, err
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users 
		SET username = $2, email = $3, display_name = $4, avatar_url = $5, updated_at = $6
		WHERE id = $1
	`
	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.DisplayName,
		user.AvatarURL, user.UpdatedAt)
	return err
}
