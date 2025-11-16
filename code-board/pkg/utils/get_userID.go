package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetUserID(ctx *gin.Context) (uint, error) {
	userIDAny, exists := ctx.Get("userID")
	if !exists {
		return 0, errors.New("userID not found in context")
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		return 0, errors.New("userID invalid type in context")
	}
	return userID, nil
}

//   userID, err := utils.GetUserID(ctx)
//     if err != nil {
//         utils.RenderError(ctx, http.StatusBadRequest, "customer", "user id not found", err)
//         return
//     }
