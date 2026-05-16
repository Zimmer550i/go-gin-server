package routes

import (
	"go-server/controllers"
	"go-server/repositories"
	"go-server/services"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	userRepository := repositories.NewInMemoryUserRepository()
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	userRoutes := r.Group("/users")
	{
		userRoutes.GET("", userController.GetUsers)
		userRoutes.POST("", userController.CreateUser)
		userRoutes.GET("/:id", userController.GetUserByID)
		userRoutes.DELETE("/:id", userController.DeleteUser)
	}
}