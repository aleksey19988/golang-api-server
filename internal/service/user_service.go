package service

import (
	"api_server/internal/domain"
	"api_server/internal/repository"
)

const MinAge = 14

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUsers() ([]domain.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetUserByID(ID uint) (*domain.User, error) {
	return s.repo.GetByID(ID)
}

func (s *UserService) GetUserByName(name string) (*domain.User, error) {
	return s.repo.GetByName(name)
}

func (s *UserService) CreateUser(name, email string, age uint) (*domain.User, error) {
	if age < MinAge {
		return nil, ErrInvalidAge
	}
	return s.repo.Create(name, email, age)
}

func (s *UserService) UpdateUser(ID uint, name string, age uint) (*domain.User, error) {
	return s.repo.Update(ID, name, age)
}

func (s *UserService) DeleteUser(ID uint) error {
	return s.repo.Delete(ID)
}
