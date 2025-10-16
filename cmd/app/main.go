package main

import (
	"api_server/internal/api"
	"api_server/internal/repository/memory"
	"api_server/internal/service"
	"log"
	"net/http"
)

func NewRouter(userService *service.UserService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", api.PingHandler)
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		api.UsersHandler(w, r, userService)
	})
	mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		api.UserHandler(w, r, userService)
	})
	return mux
}

func main() {
	repo := memory.NewMemoryUserRepository()
	s := service.NewUserService(repo)
	router := NewRouter(s)

	log.Fatal(http.ListenAndServe(":8084", router))
}
