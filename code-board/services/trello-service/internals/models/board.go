package models

import (
	"time"

	"gorm.io/gorm"
)

// ----------------------
// 🪣 Board Table
// ----------------------
// Table Name: boards
type Board struct {
	gorm.Model        // -> id, created_at, updated_at, deleted_at
	Name       string `json:"name" gorm:"column:name"`                   // -> name
	OwnerID    uint   `json:"owner_id" gorm:"column:owner_id"`           // -> owner_id
	Lists      []List `json:"lists" gorm:"constraint:OnDelete:CASCADE;"` // relation
}

// ----------------------
// 📋 List Table
// ----------------------
// Table Name: lists
type List struct {
	gorm.Model
	Name     string `json:"name" gorm:"column:name"`         // -> name
	BoardID  uint   `json:"board_id" gorm:"column:board_id"` // -> board_id (FK to boards.id)
	Cards    []Card `json:"cards" gorm:"constraint:OnDelete:CASCADE;"`
	Position int    `json:"position" gorm:"column:position"` // -> position
}

// ----------------------
// 🧾 Card Table
// ----------------------
// Table Name: cards
type Card struct {
	gorm.Model
	Title       string     `json:"title" gorm:"column:title"`             // -> title
	Description string     `json:"description" gorm:"column:description"` // -> description
	ListID      uint       `json:"list_id" gorm:"column:list_id"`         // -> list_id (FK to lists.id)
	Position    int        `json:"position" gorm:"column:position"`       // -> position
	Activities  []Activity `json:"activities" gorm:"constraint:OnDelete:CASCADE;"`
}

// ----------------------
// ⚙️ Activity Table
// ----------------------
// Table Name: activities
type Activity struct {
	gorm.Model
	CardID    uint      `json:"card_id" gorm:"column:card_id"`     // -> card_id (FK to cards.id)
	UserID    uint      `json:"user_id" gorm:"column:user_id"`     // -> user_id
	Action    string    `json:"action" gorm:"column:action"`       // -> action
	Timestamp time.Time `json:"timestamp" gorm:"column:timestamp"` // -> timestamp
}
