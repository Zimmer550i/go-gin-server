package routes

import (
	"go-server/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterHealthRoutes(r *gin.Engine) {
	r.GET("/", controllers.HealthCheck)
}