package repo

import (
	"github.com/ak-repo/code-board/services/user-auth-service/internals/models"
	"gorm.io/gorm"
)

type AuthRepo interface {
	Register(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID uint) (*models.User, error)
}

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) AuthRepo {
	return &authRepo{db: db}
}

func (a *authRepo) Register(user *models.User) error {
	return a.db.Create(&user).Error
}

func (a *authRepo) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := a.db.Where("email = ?", email).First(user).Error
	return user, err
}

func (a *authRepo) GetUserByID(userID uint) (*models.User, error) {
	user := &models.User{}
	err := a.db.First(user, userID).Error
	return user, err
}
