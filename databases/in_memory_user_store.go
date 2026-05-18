package databases

import (
	"go-server/models"
	"go-server/repositories"
)

type InMemoryUserStore struct {
	users  []models.User
	nextID int
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
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

func (r *InMemoryUserStore) FindAll() []models.User {
	return r.users
}

func (r *InMemoryUserStore) FindByID(id int) (models.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, repositories.ErrUserNotFound
}

func (r *InMemoryUserStore) Create(user models.User) models.User {
	user.ID = r.nextID
	r.nextID++

	r.users = append(r.users, user)

	return user
}

func (r *InMemoryUserStore) DeleteByID(id int) (models.User, error) {
	for index, user := range r.users {
		if user.ID == id {
			r.users = append(r.users[:index], r.users[index+1:]...)
			return user, nil
		}
	}

	return models.User{}, repositories.ErrUserNotFound
}
