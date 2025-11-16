package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/ak-repo/code-board/pkg/dto"
	"github.com/ak-repo/code-board/pkg/utils"
	"github.com/ak-repo/code-board/services/user-auth-service/internals/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Register(ctx *gin.Context, role string)
	UserRegister(ctx *gin.Context)

	Login(ctx *gin.Context, role string)
	AdminLogin(ctx *gin.Context)
	UserLogin(ctx *gin.Context)
}

// constants for roles
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// implementation struct
type authHandler struct {
	service service.AuthService
}

// constructor
func NewAuthHandler(service service.AuthService) AuthHandler {
	return &authHandler{service: service}
}

// REGISTER
func (a *authHandler) UserRegister(ctx *gin.Context) {
	a.Register(ctx, RoleUser)
}

func (a *authHandler) Register(ctx *gin.Context, role string) {

	req := &dto.RegisterRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid inputes", err.Error())
		return
	}

	response, err := a.service.Register(req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "user registration failed", err.Error())
		return

	}

	utils.SuccessResponse(ctx, http.StatusCreated, "User registered successfully", response)
}

// LOGIN
func (a *authHandler) AdminLogin(ctx *gin.Context) {
	a.Login(ctx, RoleAdmin)
}

func (a *authHandler) UserLogin(ctx *gin.Context) {
	a.Login(ctx, RoleUser)
}

func (a *authHandler) Login(ctx *gin.Context, role string) {
	req := &dto.LoginRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid inputes", err.Error())
		return
	}

	response, err := a.service.Login(req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "login failed", err.Error())
		return

	}

	refreshExp := time.Hour * time.Duration(response.RefreshExp)
	log.Println("refreshExp:", refreshExp)

	ctx.SetCookie(
		"refresh",
		response.RefreshToken,
		int(refreshExp.Seconds()),
		"/",
		"localhost",
		false,
		true,
	)

	utils.SuccessResponse(ctx, http.StatusOK, "Login successful", response)
}
