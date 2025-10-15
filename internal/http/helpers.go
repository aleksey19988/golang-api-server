package http

import (
	"api_server/internal/domain"
	"api_server/internal/service"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

func parseUserData(req *http.Request) (domain.User, error) {
	user := domain.User{}

	name := req.FormValue("name")
	if name == "" {
		return user, errors.New("Не передано имя пользователя")
	}

	age := req.FormValue("age")
	if age == "" {
		return user, errors.New("Не передан возраст пользователя")
	}

	user.Name = name

	userAge, err := strconv.Atoi(age)
	if err != nil {
		return user, errors.New("Ошибка при преобразовании возраста пользователя")
	}

	user.Age = userAge

	return user, nil
}

func GetUserIDFromUrl(url string) (int, error) {
	path := strings.Split(url, "/")
	ID := []rune(path[len(path)-1])

	if len(ID) == 0 {
		return 0, service.ErrIDNotTransmitted
	}

	if unicode.IsNumber(ID[0]) {
		userID, err := strconv.Atoi(string(ID))
		if err != nil {
			return 0, err
		}

		return userID, nil
	}

	return 0, service.ErrIDNotValid
}

func GetStatusCodeByError(err error) int {
	switch {
	case errors.Is(err, service.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, service.ErrInvalidAge),
		errors.Is(err, service.ErrIDNotTransmitted),
		errors.Is(err, service.ErrIDNotValid):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
