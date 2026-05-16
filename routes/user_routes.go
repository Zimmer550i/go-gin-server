package routes

import (
	"go-server/container"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, c *container.Container) {
	userRoutes := r.Group("/users")
	{
		userRoutes.GET("", c.UserController.GetUsers)
		userRoutes.POST("", c.UserController.CreateUser)
		userRoutes.GET("/:id", c.UserController.GetUserByID)
		userRoutes.DELETE("/:id", c.UserController.DeleteUser)
	}
}