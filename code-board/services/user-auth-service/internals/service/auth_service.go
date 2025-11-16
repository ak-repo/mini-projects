package service

import (
	"errors"
	"fmt"

	"github.com/ak-repo/code-board/pkg/config"
	"github.com/ak-repo/code-board/pkg/dto"
	"github.com/ak-repo/code-board/pkg/jwt"
	"github.com/ak-repo/code-board/pkg/utils"
	"github.com/ak-repo/code-board/services/user-auth-service/internals/models"
	"github.com/ak-repo/code-board/services/user-auth-service/internals/repo"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	repo repo.AuthRepo
	cfg  *config.Config
}

func NewAuthService(repo repo.AuthRepo, cfg *config.Config) AuthService {
	return &authService{repo: repo, cfg: cfg}
}

func (a *authService) Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: hash,
		Username:     req.Username,
		Role:         "user",
		Status:       "active",
	}

	if err := a.repo.Register(user); err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return &dto.RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (a *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := a.repo.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CompareHashAndPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	jwtData := jwt.JWT{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}

	access, err := jwt.AccessTokenGenerator(jwtData, a.cfg)
	if err != nil {
		return nil, err
	}

	refresh, err := jwt.RefreshTokenGenerator(jwtData, a.cfg)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		AccessToken:  access,
		RefreshToken: refresh,
		AccessExp:    int64(a.cfg.JWT.AccessExpireMinutes),
		RefreshExp:   int64(a.cfg.JWT.RefreshExpireHours),
	}, nil

}
