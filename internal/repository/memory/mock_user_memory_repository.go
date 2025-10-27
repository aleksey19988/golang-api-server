package memory

import (
	"api_server/internal/domain"
	"context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

func NewMockMemoryUserRepository() *UserRepository {
	// Один общий ин‑мемори инстанс для всех соединений
	dsn := "file::memory:?cache=shared"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	// Включаем FK‑ограничения для SQLite
	db.Exec("PRAGMA foreign_keys = ON;")

	// Опционально: ограничить пул соединений, чтобы не терять состояние
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Миграция схемы
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("migrate: %v", err)
	}

	return &UserRepository{
		db:  db,
		ctx: context.Background(),
	}
}

//func (r *MockUserRepository) GetAll() ([]domain.User, error) {
//	return r.users, nil
//}
//
//func (r *MockUserRepository) GetByID(id uint) (*domain.User, error) {
//	for i := range r.users {
//		if r.users[i].ID == id {
//			return &r.users[i], nil
//		}
//	}
//
//	return nil, service.ErrNotFound
//}
//
//func (r *MockUserRepository) GetByName(name string) (*domain.User, error) {
//	for i := range r.users {
//		if r.users[i].Name == name {
//			return &r.users[i], nil
//		}
//	}
//
//	return nil, service.ErrNotFound
//}
//
//func (r *MockUserRepository) Create(name, email string, age uint) (*domain.User, error) {
//	err := gorm.G[domain.User](r.db).Create(r.ctx, &domain.User{
//		Name:  name,
//		Age:   age,
//		Email: email,
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	users, err := r.GetAll()
//	if err != nil {
//		return nil, err
//	}
//
//	user, err := r.GetByID(uint(len(users)))
//
//	return user, nil
//}
//
//func (r *MockUserRepository) Update(id uint, name string, age uint) (*domain.User, error) {
//	user, err := r.GetByID(id)
//	if err != nil {
//		return nil, err
//	}
//	user.Name = name
//	user.Age = age
//
//	return user, nil
//}
//
//func (r *MockUserRepository) Delete(id uint) error {
//	for i := range r.users {
//		if r.users[i].ID == id {
//			r.users = append(r.users[:i], r.users[i+1:]...)
//			break
//		}
//	}
//	return nil
//}
