package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pesos228/bug-tracker/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func (u *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	page := getQueryInt(r.URL.Query(), "page", 1)
	pageSize := getQueryInt(r.URL.Query(), "pageSize", 10)
	fullName := getQueryString(r.URL.Query(), "fullName", "")

	users, err := u.userService.Search(r.Context(), &service.SearchUsersParams{
		Page:     page,
		PageSize: pageSize,
		FullName: fullName,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, fmt.Sprintf("Failed to Encode DTO: %s", err.Error()), http.StatusInternalServerError)
	}
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
