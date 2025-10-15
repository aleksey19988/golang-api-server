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

func (s *UserService) GetUsers() []domain.User {
	return s.repo.GetAll()
}

func (s *UserService) GetUserByID(ID int) (*domain.User, error) {
	return s.repo.GetByID(ID)
}

func (s *UserService) GetUserByName(name string) (*domain.User, error) {
	return s.repo.GetByName(name)
}

func (s *UserService) CreateUser(name string, age int) (*domain.User, error) {
	if age < MinAge {
		return nil, ErrInvalidAge
	}
	return s.repo.Create(name, age)
}

func (s *UserService) UpdateUser(ID int, name string, age int) (*domain.User, error) {
	return s.repo.Update(ID, name, age)
}

func (s *UserService) DeleteUser(ID int) error {
	return s.repo.Delete(ID)
}
