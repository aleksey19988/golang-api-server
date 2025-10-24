package memory

import (
	"api_server/internal/domain"
	"api_server/internal/service"
)

type MockMemoryUserRepository struct {
	users []domain.User
}

func NewMockMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: []domain.User{
			{ID: 1, Name: "Алексей", Age: 27},
			{ID: 2, Name: "Валерия", Age: 23},
		},
	}
}
func (r *MockMemoryUserRepository) GetAll() []domain.User {
	return r.users
}

func (r *MockMemoryUserRepository) GetByID(id int) (*domain.User, error) {
	for i := range r.users {
		if r.users[i].ID == id {
			return &r.users[i], nil
		}
	}

	return nil, service.ErrNotFound
}

func (r *MockMemoryUserRepository) GetByName(name string) (*domain.User, error) {
	for i := range r.users {
		if r.users[i].Name == name {
			return &r.users[i], nil
		}
	}

	return nil, service.ErrNotFound
}

func (r *MockMemoryUserRepository) Create(name string, age int) (*domain.User, error) {
	user := domain.User{
		ID:   len(r.users) + 1,
		Name: name,
		Age:  age,
	}
	r.users = append(r.users, user)

	return &user, nil
}

func (r *MockMemoryUserRepository) Update(id int, name string, age int) (*domain.User, error) {
	user, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	user.Name = name
	user.Age = age

	return user, nil
}

func (r *MockMemoryUserRepository) Delete(id int) error {
	for i := range r.users {
		if r.users[i].ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
		}
	}

	return nil
}
