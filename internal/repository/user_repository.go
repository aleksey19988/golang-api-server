package repository

import "api_server/internal/domain"

type UserRepositoryInterface interface {
	GetAll() ([]domain.User, error)
	GetByID(id uint) (*domain.User, error)
	GetByName(name string) (*domain.User, error)
	Create(name, email string, age uint) (*domain.User, error)
	Update(ID uint, name string, age uint) (*domain.User, error)
	Delete(id uint) error
}
