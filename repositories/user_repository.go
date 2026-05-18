package repositories

import (
	"errors"
	"go-server/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	FindAll() []models.User
	FindByID(id int) (models.User, error)
	Create(user models.User) models.User
	DeleteByID(id int) (models.User, error)
}