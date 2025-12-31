package http

import (
	"net/http"
	"strconv"

	domain "github.com/ak-repo/go-chat-system/internal/domian"
	"github.com/gin-gonic/gin"
)

type CreateConversationRequest struct {
	Type      string   `json:"type" binding:"required,oneof=one_to_one group"`
	MemberIDs []string `json:"member_ids" binding:"required,min=1"`
	Name      *string  `json:"name"`
}

func (h *Handler) createConversation(c *gin.Context) {
	userID := c.GetString("user_id")

	var req CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var convType domain.ConversationType
	if req.Type == "one_to_one" {
		convType = domain.ConversationTypeOneToOne
		if len(req.MemberIDs) != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "One-to-one conversation requires exactly one other member"})
			return
		}
	} else {
		convType = domain.ConversationTypeGroup
	}

	conv, err := h.chatService.CreateConversation(userID, req.MemberIDs, convType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         conv.ID,
		"type":       conv.Type,
		"created_at": conv.CreatedAt,
	})
}

func (h *Handler) getMessages(c *gin.Context) {
	userID := c.GetString("user_id")
	conversationID := c.Param("id")

	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	if limit > 100 {
		limit = 100
	}

	messages, err := h.chatService.GetConversationMessages(userID, conversationID, limit, offset)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
