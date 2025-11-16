package repo

import (
	"github.com/ak-repo/code-board/services/trello-service/internals/models"
	"gorm.io/gorm"
)

type BoardRepo interface {
	// Create a new board
	CreateBoard(board *models.Board) error
	// Get all boards for a user (owner)
	GetBoardsByUser(userID uint) ([]models.Board, error)

	// Get single board by ID
	GetBoardByID(id uint) (*models.Board, error)

	// Update board details (e.g. name)
	UpdateBoard(updatedData *models.Board) error

	// Delete a board and all related lists/cards
	DeleteBoard(id uint) error
}

type boardRepo struct {
	db *gorm.DB
}

func NewBoardRepo(db *gorm.DB) BoardRepo {

	return &boardRepo{db: db}
}

func (r *boardRepo) CreateBoard(board *models.Board) error {

	return r.db.Create(board).Error

}

func (r *boardRepo) GetBoardsByUser(userID uint) ([]models.Board, error) {
	boards := []models.Board{}
	err := r.db.Where("owner_id = ?", userID).Find(&boards).Error
	return boards, err

}

func (r *boardRepo) GetBoardByID(id uint) (*models.Board, error) {

	board := &models.Board{}
	err := r.db.Preload("Lists").First(board, id).Error
	return board, err

}

func (r *boardRepo) UpdateBoard(updatedData *models.Board) error {

	return r.db.Save(updatedData).Error

}

func (r *boardRepo) DeleteBoard(id uint) error {
	return r.db.Delete(&models.Board{}, id).Error
}
