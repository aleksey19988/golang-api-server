package api

import (
	"api_server/internal/domain"
	"api_server/internal/repository/memory"
	"api_server/internal/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	repo    = memory.NewMockMemoryUserRepository()
	svc     = service.NewUserService(repo)
	handler = func(s *service.UserService) *Handler {
		h := &Handler{
			userService: s,
		}
		_, err := h.userService.CreateUser("Test name 1", "test_1@example.com", 25)
		if err != nil {
			panic(err)
		}
		_, err = h.userService.CreateUser("Test name 2", "test_2@example.com", 50)
		if err != nil {
			panic(err)
		}

		return h
	}(svc)
)

func TestHandler_ParseUserId(t *testing.T) {
	h := &Handler{}
	tests := []struct {
		input       string
		expected    uint
		expectError bool
	}{
		{"123", 123, false},
		{"0", 0, false},
		{"abc", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		got, err := h.ParseUserId(tt.input)
		if (err != nil) != tt.expectError {
			t.Errorf("ParseUserId(%q) error = %v, wantErr %v", tt.input, err, tt.expectError)
		}
		if got != tt.expected {
			t.Errorf("ParseUserId(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}

func TestHandler_Ping(t *testing.T) {
	r := gin.Default()
	// Создаём тестовый роутер
	r.GET("/ping", handler.Ping)

	// Создаём тестовый HTTP-запрос
	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	// Выполняем запрос
	r.ServeHTTP(w, req)

	// Проверяем статус ответа
	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	// Проверяем тело ответа
	if w.Body.String() != "pong" {
		t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), "pong")
	}
}

func TestHandler_GetUsers(t *testing.T) {
	r := gin.Default()
	r.GET("/users", handler.GetUsers)
	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
	var gotUsers []domain.User
	err = json.Unmarshal(w.Body.Bytes(), &gotUsers)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	if len(gotUsers) != 2 {
		t.Errorf("handler returned wrong number of users: got %v want %v", len(gotUsers), 1)
	}

	gotUser := gotUsers[0]
	wantUser := domain.User{
		Name:  "Test name 1",
		Age:   25,
		Email: "test_1@example.com",
	}

	if gotUser.Name != wantUser.Name || gotUser.Email != wantUser.Email || gotUser.Age != wantUser.Age {
		t.Errorf("handler returned unexpected body: got %v want %v", gotUser, wantUser)
	}
}

func TestHandler_GetUserSuccess(t *testing.T) {
	r := gin.Default()
	// Создаём тестовый роутер
	r.GET("/user/:id", handler.GetUser)

	req, err := http.NewRequest(http.MethodGet, "/user/1", nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
	var gotUser domain.User
	err = json.Unmarshal(w.Body.Bytes(), &gotUser)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	wantUser := domain.User{
		Name:  "Test name 1",
		Age:   25,
		Email: "test_1@example.com",
	}
	if gotUser.Name != wantUser.Name || gotUser.Email != wantUser.Email || gotUser.Age != wantUser.Age {
		t.Errorf("handler returned unexpected body: got %v want %v", gotUser, wantUser)
	}
}

func TestHandler_GetUserError(t *testing.T) {
	r := gin.Default()
	// Проверяем несуществующий id
	r.GET("/user/:id", handler.GetUser)

	req, err := http.NewRequest(http.MethodGet, "/user/15", nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusNotFound)
	}

	// Проверяем передачу строки вместо числа
	r2 := gin.Default()
	r2.GET("/user/:id", handler.GetUser)
	req, err = http.NewRequest(http.MethodGet, "/user/test", nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_CreateUserSuccess(t *testing.T) {
	// Успешное создание пользователя
	r := gin.Default()
	r.POST("/user", handler.CreateUser)
	jsonBody := `{"name": "Test Name 3", "age": 75, "email": "test_3@example.com"}`
	req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(jsonBody))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusCreated)
	}
	var gotUser domain.User
	err = json.Unmarshal(w.Body.Bytes(), &gotUser)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	wantUser := domain.User{
		Name:  "Test Name 3",
		Age:   75,
		Email: "test_3@example.com",
	}
	if gotUser.Name != wantUser.Name || gotUser.Email != wantUser.Email || gotUser.Age != wantUser.Age {
		t.Errorf("handler returned unexpected body: got %v want %v", gotUser, wantUser)
	}
}

func TestHandler_CreateUserNotAllParams(t *testing.T) {
	//Передача не всех параметров для создания
	r := gin.Default()
	r.POST("/user", handler.CreateUser)
	jsonBody := `{"name": "Test Name"}`
	req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(jsonBody))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_CreateUserNotJsonRequest(t *testing.T) {
	//Передача невалидной строки
	r := gin.Default()
	r.POST("/user", handler.CreateUser)
	jsonBody := `test string`
	req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(jsonBody))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusInternalServerError)
	}
}

func TestHandler_UpdateUserName(t *testing.T) {
	r := gin.Default()
	r.PATCH("/user/:id", handler.UpdateUser)
	jsonBody := `{"name": "Updated test Name"}`
	req, err := http.NewRequest(http.MethodPatch, "/user/2", strings.NewReader(jsonBody))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusCreated)
	}
	var gotUser domain.User
	err = json.Unmarshal(w.Body.Bytes(), &gotUser)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	wantUser := domain.User{
		Name:  "Updated test Name",
		Age:   50,
		Email: "test_2@example.com",
	}
	if gotUser.Name != wantUser.Name || gotUser.Email != wantUser.Email || gotUser.Age != wantUser.Age {
		t.Errorf("handler returned unexpected body: got %v want %v", gotUser, wantUser)
	}
}

func TestHandler_UpdateUserAge(t *testing.T) {
	r := gin.Default()
	r.PATCH("/user/:id", handler.UpdateUser)
	jsonBody := `{"age": 20}`
	req, err := http.NewRequest(http.MethodPatch, "/user/1", strings.NewReader(jsonBody))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
	var gotUser domain.User
	err = json.Unmarshal(w.Body.Bytes(), &gotUser)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	wantUser := domain.User{
		Name:  "Test name 1",
		Age:   20,
		Email: "test_1@example.com",
	}
	if gotUser.Name != wantUser.Name || gotUser.Email != wantUser.Email || gotUser.Age != wantUser.Age {
		t.Errorf("handler returned unexpected body: got %v want %v", gotUser, wantUser)
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	r := gin.Default()
	r.DELETE("/user/:id", handler.DeleteUser)
	req, err := http.NewRequest(http.MethodDelete, "/user/1", nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}
