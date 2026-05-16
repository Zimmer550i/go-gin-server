package services

import (
	"go-server/dto"
	"go-server/models"
	"go-server/repositories"
)

type UserService interface {
	GetUsers() []models.User
	GetUserByID(id int) (models.User, error)
	CreateUser(req dto.CreateUserRequest) models.User
	DeleteUser(id int) (models.User, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) GetUsers() []models.User {
	return s.userRepository.FindAll()
}

func (s *userService) GetUserByID(id int) (models.User, error) {
	return s.userRepository.FindByID(id)
}

func (s *userService) CreateUser(req dto.CreateUserRequest) models.User {
	user := models.User{
		Name: req.Name,
		Age:  req.Age,
	}

	return s.userRepository.Create(user)
}

func (s *userService) DeleteUser(id int) (models.User, error) {
	return s.userRepository.DeleteByID(id)
}