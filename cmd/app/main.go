package main

import (
	"api_server/internal/api"
	"api_server/internal/repository/memory"
	"api_server/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"os"
)

const defaultPort = ":8080"

// @title           Example user API
// @version         1.0
// @description     Это учебный проект для практики написания на Go

// @contact.name   Aleksey Kononenko
// @contact.email  a.knnnk@mail.ru

// @host      localhost:8080
// @BasePath  /
func main() {
	r := gin.Default()

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err.Error())
	}

	repo := memory.NewMemoryUserRepository()
	s := service.NewUserService(repo)
	handler := api.NewHandler(s)

	r.GET("/ping", handler.Ping)
	r.GET("/user/:id", handler.GetUser)
	r.POST("/user", handler.CreateUser)
	r.PATCH("/user/:id", handler.UpdateUser)
	r.DELETE("/user/:id", handler.DeleteUser)

	r.GET("/users", handler.GetUsers)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
