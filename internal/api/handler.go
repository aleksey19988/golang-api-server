package api

import (
	"api_server/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	userService *service.UserService
}

type UpdateUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type CreateUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
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

func (h *Handler) PingHandler(c *gin.Context) {
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
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
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
	c.JSON(http.StatusOK, user)
	return
}

// CreateUser godoc
// @Summary Создание нового пользователя
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request   body      CreateUserRequest  true  "JSON"
// @Success      200       {object}  domain.User
// @Failure      400       {object}  ErrorResponse
// @Failure      500       {object}  ErrorResponse
// @Failure      404       {object}  ErrorResponse
// @Router       /user [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var request CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	u, err := h.userService.CreateUser(request.Name, request.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, u)
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
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	updatedUser, err := h.userService.UpdateUser(user.ID, request.Name, request.Age)
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
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
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
