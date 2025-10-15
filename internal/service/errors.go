package service

import "errors"

var (
	ErrNotFound         = errors.New("Пользователь не найден")
	ErrInvalidAge       = errors.New("Возраст пользователя не может быть менее 14 лет")
	ErrIDNotTransmitted = errors.New("ID пользователя не передан")
	ErrIDNotValid       = errors.New("Некорректный ID пользователя")
)
