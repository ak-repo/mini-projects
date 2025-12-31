package domain

import "time"

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	DisplayName  string
	AvatarURL    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) error
}
