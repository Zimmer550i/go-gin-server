package container

import (
	"fmt"
	"go-server/controllers"
	"go-server/databases"
	"go-server/repositories"
	"go-server/services"
	"os"
)

type Container struct {
	UserController *controllers.UserController
	closers        []func() error
}

func NewContainer() (*Container, error) {
	userRepository, closeUserRepository, err := newUserRepository()
	if err != nil {
		return nil, err
	}

	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	return &Container{
		UserController: userController,
		closers:        []func() error{closeUserRepository},
	}, nil
}

func (c *Container) Close() error {
	for _, close := range c.closers {
		if close == nil {
			continue
		}

		if err := close(); err != nil {
			return err
		}
	}

	return nil
}

func newUserRepository() (repositories.UserRepository, func() error, error) {
	switch os.Getenv("USER_STORE") {
	case "", "memory", "in-memory":
		println("Using In-Memory Database")
		return databases.NewInMemoryUserStore(), nil, nil
	case "postgres", "pg":
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			return nil, nil, fmt.Errorf("DATABASE_URL is required when USER_STORE=postgres")
		}

		store, err := databases.NewPostgresUserStore(databaseURL)
		if err != nil {
			return nil, nil, err
		}

		println("Using PostgreSql Database")
		return store, store.Close, nil
	default:
		return nil, nil, fmt.Errorf("unsupported USER_STORE %q", os.Getenv("USER_STORE"))
	}
}
