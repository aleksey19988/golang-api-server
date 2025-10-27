package memory

import (
	"api_server/internal/domain"
	"api_server/internal/service"
	"context"
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepository struct {
	db    *gorm.DB
	ctx   context.Context
	users []domain.User
}

func NewUserRepository() *UserRepository {
	dsn := "host=localhost user=postgres dbname=golang_api password=KrA2/xW/ sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		panic(err)
	}

	return &UserRepository{
		db:  db,
		ctx: context.Background(),
	}
}
func (r *UserRepository) GetAll() ([]domain.User, error) {
	users, err := gorm.G[domain.User](r.db).Find(r.ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetByID(id uint) (*domain.User, error) {
	user, err := gorm.G[domain.User](r.db).Where("id = ?", id).First(r.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByName(name string) (*domain.User, error) {
	user, err := gorm.G[domain.User](r.db).Where("name = ?", name).First(r.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(name, email string, age uint) (*domain.User, error) {
	err := gorm.G[domain.User](r.db).Create(r.ctx, &domain.User{
		Name:  name,
		Age:   age,
		Email: email,
	})

	if err != nil {
		return nil, err
	}

	users, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	user, err := r.GetByID(uint(len(users)))

	return user, nil
}

func (r *UserRepository) Update(id uint, name string, age uint) (*domain.User, error) {
	_, err := gorm.G[domain.User](r.db).Where("id = ?", id).Update(r.ctx, "Name", name)
	if err != nil {
		return nil, err
	}

	_, err = gorm.G[domain.User](r.db).Where("id = ?", id).Update(r.ctx, "Age", age)
	if err != nil {
		return nil, err
	}

	user, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Delete(id uint) error {
	_, err := gorm.G[domain.User](r.db).Where("id = ?", id).Delete(r.ctx)
	if err != nil {
		return err
	}

	return nil
}
