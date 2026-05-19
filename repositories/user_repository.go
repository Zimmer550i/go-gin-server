package repositories

import (
	"errors"
	"go-server/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	FindAll() ([]models.User, error)
	FindByID(id int) (models.User, error)
	Create(user models.User) (models.User, error)
	DeleteByID(id int) (models.User, error)
}
