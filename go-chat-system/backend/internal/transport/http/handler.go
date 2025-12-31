package http

import (
	"net/http"

	"github.com/ak-repo/go-chat-system/internal/config"
	"github.com/ak-repo/go-chat-system/internal/service"
	"github.com/ak-repo/go-chat-system/internal/transport/websocket"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService *service.AuthService
	chatService *service.ChatService
	hub         *websocket.Hub
	config      *config.Config
}

func NewHandler(
	authService *service.AuthService,
	chatService *service.ChatService,
	hub *websocket.Hub,
	config *config.Config,
) http.Handler {
	h := &Handler{
		authService: authService,
		chatService: chatService,
		hub:         hub,
		config:      config,
	}

	router := gin.Default()

	// Middleware
	router.Use(corsMiddleware())

	// Public routes
	router.POST("/api/auth/register", h.register)
	router.POST("/api/auth/login", h.login)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(h.authMiddleware())
	{
		protected.GET("/ws", h.websocketHandler)
		protected.POST("/conversations", h.createConversation)
		protected.GET("/conversations/:id/messages", h.getMessages)
		protected.GET("/users/me", h.getCurrentUser)
	}

	return router
}

// internal/transport/http/auth.go
type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	DisplayName string `json:"display_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

func (h *Handler) register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(req.Username, req.Email, req.Password, req.DisplayName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User: gin.H{
			"id":           user.ID,
			"username":     user.Username,
			"email":        user.Email,
			"display_name": user.DisplayName,
		},
	})
}

func (h *Handler) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User: gin.H{
			"id":           user.ID,
			"username":     user.Username,
			"email":        user.Email,
			"display_name": user.DisplayName,
		},
	})
}

func (h *Handler) getCurrentUser(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, user)
}

// internal/transport/http/chat.go
