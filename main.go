package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Response struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	STATUS_SUCCESS = "success"
	SYATUS_ERROR   = "error"
)

var users = map[int]User{
	1: {
		ID:   1,
		Name: "Aleksey Kononenko",
		Age:  27,
	},
}

func main() {
	pingHandler := func(w http.ResponseWriter, req *http.Request) {
		_, err := io.WriteString(w, "pong\n")
		if err != nil {
			return
		}
	}

	usersHandler := func(w http.ResponseWriter, req *http.Request) {
		usersJson, err := json.Marshal(users)
		if err != nil {
			return
		}
		_, err = io.WriteString(w, string(usersJson))
		if err != nil {
			return
		}
	}

	userHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := strings.Split(req.URL.Path, "/")
		ID := []rune(path[len(path)-1])
		if unicode.IsNumber(ID[0]) {
			userID, err := strconv.Atoi(string(ID[0]))
			if err != nil {
				writeErrorResponse(w, "Ошибка при проверке ID пользователя", http.StatusInternalServerError)
				return
			}
			user, isExists := users[userID]
			if !isExists {
				writeErrorResponse(w, "Пользователь с id "+string(ID)+" не найден", http.StatusNotFound)
				return
			}

			if req.Method == http.MethodGet {
				userJSON, err := json.Marshal(user)
				if err != nil {
					writeErrorResponse(w, "Ошибка при получении данных пользователя", http.StatusInternalServerError)
				}
				writeSuccessResponse(w, string(userJSON), http.StatusOK)
			} else if req.Method == http.MethodPost {
				w.Write([]byte("Добавление пользователя"))
			}
			return
		} else {
			writeErrorResponse(w, "Некорректный id пользователя", http.StatusBadRequest)
			return
		}

	}

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/user/", userHandler)
	log.Fatal(http.ListenAndServe(":8085", nil))
}

func writeErrorResponse(
	w http.ResponseWriter,
	message string,
	code int,
) {
	w.WriteHeader(code)
	response := Response{
		Status:  SYATUS_ERROR,
		Code:    code,
		Message: message,
	}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

func writeSuccessResponse(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	response := Response{
		Status:  SYATUS_ERROR,
		Code:    code,
		Message: message,
	}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}
