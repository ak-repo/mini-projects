package dto

import "time"

// CreateBoardRequest - payload to create a board
type CreateBoardRequest struct {
	OwnerID uint   `json:"owner_id"`
	Name    string `json:"name" binding:"required"`
}

// UpdateBoardRequest - payload to update board details
type UpdateBoardRequest struct {
	Name    string `json:"name" binding:"required"`
	OwnerID uint   `json:"owner_id"`
}

// BoardResponse - returned when fetching boards
type BoardResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	OwnerID   uint           `json:"owner_id"`
	Lists     []ListResponse `json:"lists,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// ======================================
// LIST DTOs
// ======================================

// CreateListRequest - create list inside a board
type CreateListRequest struct {
	Name     string `json:"name" binding:"required"`
	BoardID  uint   `json:"board_id" binding:"required"`
	Position int    `json:"position"`
}

// UpdateListRequest - update list name or position
type UpdateListRequest struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
}

// ListResponse - data returned for a list
type ListResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	BoardID   uint           `json:"board_id"`
	Position  int            `json:"position"`
	Cards     []CardResponse `json:"cards,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// ======================================
// CARD DTOs
// ======================================

// CreateCardRequest - payload for creating card
type CreateCardRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	ListID      uint   `json:"list_id" binding:"required"`
	Position    int    `json:"position"`
}

// UpdateCardRequest - payload for updating a card
type UpdateCardRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Position    int    `json:"position"`
}

// CardResponse - data returned for card
type CardResponse struct {
	ID          uint               `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	ListID      uint               `json:"list_id"`
	Position    int                `json:"position"`
	Activities  []ActivityResponse `json:"activities,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// ======================================
// ACTIVITY DTOs
// ======================================

// CreateActivityRequest - log a user action
type CreateActivityRequest struct {
	CardID uint   `json:"card_id" binding:"required"`
	UserID uint   `json:"user_id" binding:"required"`
	Action string `json:"action" binding:"required"`
}

// ActivityResponse - data returned for activity
type ActivityResponse struct {
	ID        uint      `json:"id"`
	CardID    uint      `json:"card_id"`
	UserID    uint      `json:"user_id"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}
