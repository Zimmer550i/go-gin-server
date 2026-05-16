package controllers

import (
	"net/http"

	"go-server/utils"

	"github.com/gin-gonic/gin"
)

func HealthCheck(ctx *gin.Context) {
	utils.Success(ctx, http.StatusOK, "Server is running", nil)
}