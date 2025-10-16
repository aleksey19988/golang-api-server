package api

import (
	"api_server/internal/service"
	"encoding/json"
	"net/http"
)

func UsersHandler(w http.ResponseWriter, req *http.Request, s *service.UserService) {
	w.Header().Set("Content-Type", "application/json")
	if req.Method == http.MethodPost {
		createUserHandler(w, req, s)
	} else if req.Method == http.MethodGet {
		users := s.GetUsers()
		usersJson, err := json.Marshal(users)

		if err != nil {
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeSuccessResponse(w, string(usersJson), http.StatusOK)
	} else {
		writeErrorResponse(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}
}

func UserHandler(w http.ResponseWriter, req *http.Request, s *service.UserService) {
	w.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case http.MethodGet:
		getUserHandler(w, req, s)
	case http.MethodPatch:
		updateUserHandler(w, req, s)
	case http.MethodDelete:
		deleteUserHandler(w, req, s)
	default:
		writeErrorResponse(w, "Метод не разрешён", http.StatusMethodNotAllowed)
	}
}

func getUserHandler(w http.ResponseWriter, req *http.Request, s *service.UserService) {
	ID, err := GetUserIDFromUrl(req.URL.Path)
	if err != nil {
		writeErrorResponse(w, err.Error(), GetStatusCodeByError(err))
		return
	}
	user, err := s.GetUserByID(ID)
	if err != nil {
		writeErrorResponse(w, err.Error(), GetStatusCodeByError(err))
		return
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		writeErrorResponse(w, err.Error(), GetStatusCodeByError(err))
	}
	writeSuccessResponse(w, string(userJSON), http.StatusOK)
	return
}

func createUserHandler(w http.ResponseWriter, req *http.Request, s *service.UserService) {
	userData, err := parseUserData(req)
	if err != nil {
		writeErrorResponse(w, "Ошибка при проверке данных пользователя: "+err.Error(), GetStatusCodeByError(err))
		return
	}
	user, err := s.GetUserByName(userData.Name)
	if err == nil {
		writeErrorResponse(w, "Пользователь с именем '"+user.Name+"' уже существует", GetStatusCodeByError(err))
		return
	}
	user, err = s.CreateUser(userData.Name, userData.Age)
	if err != nil {
		writeErrorResponse(w, err.Error(), GetStatusCodeByError(err))
		return
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		writeErrorResponse(w, "Ошибка при получении данных созданного пользователя", GetStatusCodeByError(err))
		return
	}
	writeSuccessResponse(w, string(userJSON), http.StatusCreated)
}

func updateUserHandler(w http.ResponseWriter, req *http.Request, s *service.UserService) {
	ID, err := GetUserIDFromUrl(req.URL.Path)
	if err != nil {
		writeErrorResponse(w, err.Error(), GetStatusCodeByError(err))
	}
	user, err := s.GetUserByID(ID)
	if err != nil {
		writeErrorResponse(w, "Ошибка при поиске пользователя:"+err.Error(), GetStatusCodeByError(err))
		return
	}
	userData, err := parseUserData(req)
	if err != nil {
		writeErrorResponse(w, "Ошибка при проверке данных пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}
	user, err = s.UpdateUser(ID, userData.Name, userData.Age)
	if err != nil {
		writeErrorResponse(w, err.Error(), http.StatusNotFound)
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		writeErrorResponse(w, "Ошибка при получении данных пользователя: "+err.Error(), http.StatusInternalServerError)
	}
	writeSuccessResponse(w, string(userJSON), http.StatusOK)
	return
}

func deleteUserHandler(w http.ResponseWriter, req *http.Request, s *service.UserService) {
	ID, err := GetUserIDFromUrl(req.URL.Path)
	if err != nil {
		writeErrorResponse(w, err.Error(), http.StatusBadRequest)
	}
	err = s.DeleteUser(ID)
	if err != nil {
		return
	}
	writeSuccessResponse(w, "Пользователь успешно удалён", http.StatusOK)
	return
}
