package repository

import "api_server/internal/domain"

type UserRepositoryInterface interface {
	GetAll() []domain.User
	GetByID(id int) (*domain.User, error)
	GetByName(name string) (*domain.User, error)
	Create(name string, age int) (*domain.User, error)
	Update(ID int, name string, age int) (*domain.User, error)
	Delete(id int) error
}
