package main

import (
	_ "api_server/docs"
	"api_server/internal/api"
	"api_server/internal/repository"
	"api_server/internal/repository/memory"
	"api_server/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
)

// @title           Example user API
// @version         1.0
// @description     Это учебный проект для практики написания на Go

// @contact.name   Aleksey Kononenko
// @contact.email  a.knnnk@mail.ru

// @host      localhost:8080
// @BasePath  /
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	DB := repository.NewDB(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)
	repo, err := memory.NewUserRepository(DB)
	if err != nil {
		panic(err)
	}
	s := service.NewUserService(repo)
	handler := api.NewHandler(s)

	r := gin.Default()
	r.GET("/ping", handler.Ping)
	r.GET("/user/:id", handler.GetUser)
	r.POST("/user", handler.CreateUser)
	r.PATCH("/user/:id", handler.UpdateUser)
	r.DELETE("/user/:id", handler.DeleteUser)

	r.GET("/users", handler.GetUsers)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = r.Run(":" + os.Getenv("API_PORT"))
	if err != nil {
		panic(err)
	}
}
