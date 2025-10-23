package main

import (
	"api_server/internal/api"
	"api_server/internal/repository/memory"
	"api_server/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const defaultPort = ":8080"

func main() {
	r := gin.Default()

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err.Error())
	}

	repo := memory.NewMemoryUserRepository()
	s := service.NewUserService(repo)
	handler := api.NewHandler(s)

	r.GET("/ping", handler.PingHandler)
	r.GET("/user/:id", handler.GetUserHandler)
	r.POST("/user", handler.CreateUserHandler)
	r.PATCH("/user/:id", handler.UpdateUserHandler)
	r.DELETE("/user/:id", handler.DeleteUserHandler)

	r.GET("/users", handler.GetUsersHandler)

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Printf("$port is not set\n")
		port = defaultPort
	}

	err = r.Run(port)
	if err != nil {
		panic(err)
	}
}
