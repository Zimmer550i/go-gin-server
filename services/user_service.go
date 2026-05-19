package services

import (
	"go-server/dto"
	"go-server/models"
	"go-server/repositories"
)

type UserService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetUsers() ([]models.User, error) {
	return s.userRepository.FindAll()
}

func (s *UserService) GetUserByID(id int) (models.User, error) {
	return s.userRepository.FindByID(id)
}

func (s *UserService) CreateUser(req dto.CreateUserRequest) (models.User, error) {
	user := models.User{
		Name: req.Name,
		Age:  req.Age,
	}

	return s.userRepository.Create(user)
}

func (s *UserService) DeleteUser(id int) (models.User, error) {
	return s.userRepository.DeleteByID(id)
}
