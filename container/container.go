package container

import (
	"go-server/controllers"
	"go-server/databases"
	"go-server/services"
)

type Container struct {
	UserController *controllers.UserController
}

func NewContainer() *Container {
	userRepository := databases.NewInMemoryUserStore()
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	return &Container{
		UserController: userController,
	}
}
