package handler

import (
	"net/http"

	"github.com/ak-repo/code-board/pkg/dto"
	"github.com/ak-repo/code-board/pkg/utils"
	"github.com/ak-repo/code-board/services/trello-service/internals/service"
	"github.com/gin-gonic/gin"
)

// BoardHandler defines all operations related to Boards
type BoardHandler interface {
	// Create a new board
	CreateBoard(ctx *gin.Context)

	// // Get all boards for a user (owner)
	// GetBoardsByUser(userID uint)

	// // Get single board by ID
	GetBoardByID(ctx *gin.Context)

	// // Update board details (e.g. name)
	// UpdateBoard(ctx *gin.Context)

	// // Delete a board and all related lists/cards
	// DeleteBoard(ctx *gin.Context)
}

type boardHandler struct {
	service service.BoardService
}

func NewBoardHandler(service service.BoardService) BoardHandler {

	return &boardHandler{service: service}
}

func (h *boardHandler) CreateBoard(ctx *gin.Context) {

	// userID, err := utils.GetUserID(ctx)
	// if err != nil {
	// 	utils.ErrorResponse(ctx, http.StatusBadRequest, "customer", "user id not found", err)
	// 	return
	// }
	req := &dto.CreateBoardRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	req.OwnerID = uint(1)

	board, err := h.service.CreateBoard(req)
	if err != nil {

		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "board creation failed")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "board created", board)

}

func (h *boardHandler) GetBoardByID(ctx *gin.Context) {

	response, err := h.service.GetBoardByID(1)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "no")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "list created", response)

}
