package utils

import "github.com/gin-gonic/gin"

// SuccessResponse standard format for success results
func SuccessResponse(ctx *gin.Context, statusCode int, message string, data interface{}) {
	ctx.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// ErrorResponse standard format for errors
func ErrorResponse(ctx *gin.Context, statusCode int, errMessage string, details ...interface{}) {
	ctx.JSON(statusCode, gin.H{
		"success": false,
		"error": gin.H{
			"code":    statusCode,
			"message": errMessage,
			"details": details, // optional additional info (e.g., validation errors)
		},
	})
}
