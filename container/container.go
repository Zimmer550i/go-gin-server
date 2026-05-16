package container

import (
	"go-server/controllers"
	"go-server/repositories"
	"go-server/services"
)

type Container struct {
	UserController *controllers.UserController
}

func NewContainer() *Container {
	userRepository := repositories.NewInMemoryUserRepository()
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	return &Container{
		UserController: userController,
	}
}