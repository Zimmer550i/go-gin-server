package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(ctx *gin.Context, statusCode int, message string, data interface{}) {
	ctx.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(ctx *gin.Context, statusCode int, message string, err interface{}) {
	ctx.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func UnauthorizedAccess(ctx *gin.Context, message ...string) {
	responseMessage := "You don't have permission to access this API."

	if len(message) > 0 && message[0] != "" {
		responseMessage = responseMessage + " " + message[0]
	}

	ctx.JSON(http.StatusUnauthorized, APIResponse{
		Success: false,
		Message: responseMessage,
	})
}

func BadRequest(ctx *gin.Context, message string, err interface{}) {
	Error(ctx, http.StatusBadRequest, message, err)
}

func NotFound(ctx *gin.Context, message string, err interface{}) {
	Error(ctx, http.StatusNotFound, message, err)
}
