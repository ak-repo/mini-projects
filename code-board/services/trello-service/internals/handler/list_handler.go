package handler

import (
	"net/http"

	"github.com/ak-repo/code-board/pkg/dto"
	"github.com/ak-repo/code-board/pkg/utils"
	"github.com/ak-repo/code-board/services/trello-service/internals/service"
	"github.com/gin-gonic/gin"
)

type ListHandler interface {
	// Create a new List
	CreateList(ctx *gin.Context)
	// Get all Lists for a Board
	GetListsByBoard(ctx *gin.Context)

	// // Get single List by ID
	// GetListByID(ctx *gin.Context)

	// // Update List details (e.g. name)
	// UpdateList(ctx *gin.Context)

	// // Delete a List and all related lists/cards
	// DeleteList(ctx *gin.Context)
}

type listHandler struct {
	service service.ListService
}

func NewListHandler(service service.ListService) ListHandler {

	return &listHandler{service: service}
}

func (h *listHandler) CreateList(ctx *gin.Context) {

	// userID, err := utils.GetUserID(ctx)
	// if err != nil {
	// 	utils.ErrorResponse(ctx, http.StatusBadRequest, "customer", "user id not found", err)
	// 	return
	// }
	req := &dto.CreateListRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.CreateList(req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "board creation failed")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "list created", nil)

}

func (h *listHandler) GetListsByBoard(ctx *gin.Context) {

	response, err := h.service.GetListByID(1)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "no")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "list created", response)

}
