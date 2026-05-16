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

type InMemoryUserRepository struct {
	users      []models.User
	nextID     int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: []models.User{
			{
				ID:   1,
				Name: "Wasiul",
				Age:  20,
			},
			{
				ID:   2,
				Name: "Islam",
				Age:  22,
			},
		},
		nextID: 3,
	}
}

func (r *InMemoryUserRepository) FindAll() []models.User {
	return r.users
}

func (r *InMemoryUserRepository) FindByID(id int) (models.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, ErrUserNotFound
}

func (r *InMemoryUserRepository) Create(user models.User) models.User {
	user.ID = r.nextID
	r.nextID++

	r.users = append(r.users, user)

	return user
}

func (r *InMemoryUserRepository) DeleteByID(id int) (models.User, error) {
	for index, user := range r.users {
		if user.ID == id {
			r.users = append(r.users[:index], r.users[index+1:]...)
			return user, nil
		}
	}

	return models.User{}, ErrUserNotFound
}