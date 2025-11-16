package repo

import (
	"github.com/ak-repo/code-board/services/trello-service/internals/models"
	"gorm.io/gorm"
)

type ListRepo interface {
	// Create a new List
	CreateList(List *models.List) error
	// Get all Lists for a Board
	GetListsByBoard(boardID uint) ([]models.List, error)

	// Get single List by ID
	GetListByID(id uint) (*models.List, error)

	// Update List details (e.g. name)
	UpdateList(updatedData *models.List) error

	// Delete a List and all related lists/cards
	DeleteList(id uint) error
}

type listRepo struct {
	db *gorm.DB
}

func NewListRepo(db *gorm.DB) ListRepo {

	return &listRepo{db: db}
}

func (r *listRepo) CreateList(List *models.List) error {

	return r.db.Create(List).Error

}

func (r *listRepo) GetListsByBoard(boardID uint) ([]models.List, error) {
	Lists := []models.List{}
	err := r.db.Where("board_id = ?", boardID).Find(&Lists).Error
	return Lists, err

}

func (r *listRepo) GetListByID(id uint) (*models.List, error) {

	List := &models.List{}
	err := r.db.First(List, id).Error
	return List, err

}

func (r *listRepo) UpdateList(updatedData *models.List) error {

	return r.db.Save(updatedData).Error

}

func (r *listRepo) DeleteList(id uint) error {
	return r.db.Delete(&models.List{}, id).Error
}
