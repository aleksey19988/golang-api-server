package api

import (
	"api_server/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

type Handler struct {
	userService *service.UserService
}

type UpdateUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age" validate:"min=14"`
}

type CreateUserRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required,min=14"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ErrorUserResponse struct {
	Message string `json:"message"`
}

func NewHandler(s *service.UserService) *Handler {
	return &Handler{userService: s}
}

// Ping godoc
// @Summary      Проверка соединения
// @Tags         ping
// @Produce      html
// @Success      200 {string}  string    "ok"
// @Router       /ping [get]
func (h *Handler) Ping(c *gin.Context) {
	c.String(200, "pong")
}

// GetUser godoc
// @Summary      Получение данных о пользователе по его ID
// @Tags         user
// @Produce      json
// @Param        id   query      int  true "ID пользователя"
// @Success      200  {object}  domain.User
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /user/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	id, err := h.ParseUserId(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil && errors.Is(err, service.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
	return
}

// CreateUser godoc
// @Summary Создание нового пользователя
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request   body      CreateUserRequest  true  "JSON"
// @Success      201       {object}  domain.User
// @Failure      400       {object}  ErrorResponse
// @Failure      500       {object}  ErrorResponse
// @Router       /user [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var request CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.userService.CreateUser(request.Name, request.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, u)
}

// UpdateUser godoc
// @Summary Обновление пользователя
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id        query     int  true "ID пользователя"
// @Param        request   body      UpdateUserRequest  true  "JSON"
// @Success      200       {object}  domain.User
// @Failure      400       {object}  ErrorResponse
// @Failure      500       {object}  ErrorResponse
// @Failure      404       {object}  ErrorResponse
// @Router       /user/{id} [patch]
func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := h.ParseUserId(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil && errors.Is(err, service.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	userName := user.Name
	if request.Name != userName && request.Name != "" {
		userName = request.Name
	}

	userAge := user.Age
	if request.Age != userAge && request.Age > 14 {
		userAge = request.Age
	}

	updatedUser, err := h.userService.UpdateUser(user.ID, userName, userAge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser godoc
// @Summary Удаление пользователя
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id        query     int  true "ID пользователя"
// @Success      200
// @Failure      400       {object}  ErrorResponse
// @Failure      500       {object}  ErrorResponse
// @Failure      404       {object}  ErrorResponse
// @Router       /user/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := h.ParseUserId(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	err = h.userService.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// GetUsers godoc
// @Summary      Получение списка пользователей
// @Tags         users
// @Produce      json
// @Success      200  {object}  []domain.User
// @Router       /users [get]
func (h *Handler) GetUsers(c *gin.Context) {
	users := h.userService.GetUsers()
	c.JSON(http.StatusOK, users)
}

func (h *Handler) ParseUserId(idStr string) (int, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}
