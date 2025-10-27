package main

import (
	_ "api_server/docs"
	"api_server/internal/api"
	"api_server/internal/repository/memory"
	"api_server/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Example user API
// @version         1.0
// @description     Это учебный проект для практики написания на Go

// @contact.name   Aleksey Kononenko
// @contact.email  a.knnnk@mail.ru

// @host      localhost:8080
// @BasePath  /
func main() {
	r := gin.Default()

	repo := memory.NewUserRepository()
	s := service.NewUserService(repo)
	handler := api.NewHandler(s)

	r.GET("/ping", handler.Ping)
	r.GET("/user/:id", handler.GetUser)
	r.POST("/user", handler.CreateUser)
	r.PATCH("/user/:id", handler.UpdateUser)
	r.DELETE("/user/:id", handler.DeleteUser)

	r.GET("/users", handler.GetUsers)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := r.Run(":8085")
	if err != nil {
		panic(err)
	}
}
